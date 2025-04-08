package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

type Transaction []string

type FrequentItemset struct {
	Items   []string
	Support float64
	Length  int
}

type AssociationRule struct {
	Antecedent       []string
	Consequent       []string
	Support          float64
	Confidence       float64
	Lift             float64
	LeverageMetric   float64
	ConvictionMetric float64
}

type Dataset struct {
	Transactions []Transaction
	UniqueItems  []string
	ItemsMap     map[string]bool
}

func loadTransactionsFromCSV(filePath string) (*Dataset, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	// Group by basket
	basketMap := make(map[string][]string)
	headerRow := true // Assume first row is header

	for i, record := range records {
		// Skip header row
		if i == 0 && headerRow {
			// Check if first row looks like a header
			if len(record) >= 2 && (strings.Contains(strings.ToLower(record[0]), "basket") ||
				strings.Contains(strings.ToLower(record[1]), "item")) {
				continue
			} else {
				headerRow = false // Not a header, process this row
			}
		}

		if len(record) < 2 {
			fmt.Printf("Skipping invalid row %d: fewer than 2 columns\n", i+1)
			continue
		}

		basket := strings.TrimSpace(record[0])
		item := strings.TrimSpace(record[1])

		if basket == "" || item == "" {
			continue
		}

		basketMap[basket] = append(basketMap[basket], item)
	}

	// Convert to transactions
	dataset := &Dataset{
		Transactions: make([]Transaction, 0, len(basketMap)),
		ItemsMap:     make(map[string]bool),
	}

	for _, items := range basketMap {
		// Remove duplicates within a basket
		uniqueItems := make(map[string]bool)
		for _, item := range items {
			uniqueItems[item] = true
			dataset.ItemsMap[item] = true
		}

		// Create transaction with unique items
		transaction := make(Transaction, 0, len(uniqueItems))
		for item := range uniqueItems {
			transaction = append(transaction, item)
		}

		dataset.Transactions = append(dataset.Transactions, transaction)
	}

	// Create slice of unique items
	dataset.UniqueItems = make([]string, 0, len(dataset.ItemsMap))
	for item := range dataset.ItemsMap {
		dataset.UniqueItems = append(dataset.UniqueItems, item)
	}

	sort.Strings(dataset.UniqueItems)

	return dataset, nil
}

// findFrequentItemsets finds frequent itemsets using the Apriori algorithm
func findFrequentItemsets(dataset *Dataset, minSupport float64, maxLen int) []FrequentItemset {
	transactionCount := float64(len(dataset.Transactions))
	result := make([]FrequentItemset, 0)

	// Find frequent 1-itemsets
	L1 := make([]FrequentItemset, 0)
	for _, item := range dataset.UniqueItems {
		count := 0
		for _, transaction := range dataset.Transactions {
			if containsItem(transaction, item) {
				count++
			}
		}

		support := float64(count) / transactionCount
		if support >= minSupport {
			L1 = append(L1, FrequentItemset{
				Items:   []string{item},
				Support: support,
				Length:  1,
			})
		}
	}

	result = append(result, L1...)

	Lk_1 := L1
	for k := 2; k <= maxLen; k++ {
		Ck := generateCandidates(Lk_1, k)

		Lk := make([]FrequentItemset, 0)
		for _, candidate := range Ck {
			count := 0
			for _, transaction := range dataset.Transactions {
				if isSubset(candidate.Items, transaction) {
					count++
				}
			}

			support := float64(count) / transactionCount
			if support >= minSupport {
				Lk = append(Lk, FrequentItemset{
					Items:   candidate.Items,
					Support: support,
					Length:  k,
				})
			}
		}

		if len(Lk) == 0 {
			break
		}

		result = append(result, Lk...)
		Lk_1 = Lk
	}

	return result
}

func generateCandidates(itemsets []FrequentItemset, k int) []FrequentItemset {
	candidates := make([]FrequentItemset, 0)

	for i := 0; i < len(itemsets); i++ {
		for j := i + 1; j < len(itemsets); j++ {
			if k > 2 {
				canJoin := true
				for l := 0; l < k-2; l++ {
					if itemsets[i].Items[l] != itemsets[j].Items[l] {
						canJoin = false
						break
					}
				}

				if !canJoin {
					continue
				}
			}

			candidate := make([]string, k)
			copy(candidate, itemsets[i].Items)
			candidate[k-1] = itemsets[j].Items[k-2]

			sort.Strings(candidate)

			if k > 2 {
				isValid := true
				for l := 0; l < k; l++ {
					subset := make([]string, k-1)
					copy(subset[:l], candidate[:l])
					if l < k-1 {
						copy(subset[l:], candidate[l+1:])
					}

					isFrequent := false
					for _, itemset := range itemsets {
						if slicesEqual(itemset.Items, subset) {
							isFrequent = true
							break
						}
					}

					if !isFrequent {
						isValid = false
						break
					}
				}

				if !isValid {
					continue
				}
			}

			candidates = append(candidates, FrequentItemset{
				Items:  candidate,
				Length: k,
			})
		}
	}

	return candidates
}

// generateAssociationRules generates association rules from frequent itemsets
func generateAssociationRules(itemsets []FrequentItemset, minConfidence float64) []AssociationRule {
	rules := make([]AssociationRule, 0)
	itemsetMap := make(map[string]float64)

	// Create a map for quick lookup of itemset support
	for _, itemset := range itemsets {
		itemsetMap[strings.Join(itemset.Items, ",")] = itemset.Support
	}

	// Generate rules for each itemset with length > 1
	for _, itemset := range itemsets {
		if itemset.Length <= 1 {
			continue
		}

		// Generate all possible non-empty subsets as antecedents
		antecedents := generateAllSubsets(itemset.Items)

		for _, antecedent := range antecedents {
			// Skip if antecedent is empty or the same as the itemset
			if len(antecedent) == 0 || len(antecedent) == len(itemset.Items) {
				continue
			}

			// Generate consequent
			consequent := difference(itemset.Items, antecedent)

			// Get antecedent support
			antecedentKey := strings.Join(antecedent, ",")
			antecedentSupport, exists := itemsetMap[antecedentKey]
			if !exists {
				continue // Should not happen with proper subsets
			}

			// Calculate confidence
			confidence := itemset.Support / antecedentSupport

			if confidence >= minConfidence {
				// Calculate additional metrics
				consequentKey := strings.Join(consequent, ",")
				consequentSupport, exists := itemsetMap[consequentKey]
				if !exists {
					continue // Should not happen with proper subsets
				}

				lift := confidence / consequentSupport
				leverage := itemset.Support - (antecedentSupport * consequentSupport)

				var conviction float64
				if consequentSupport == 1.0 || confidence == 1.0 {
					conviction = math.Inf(1)
				} else {
					conviction = (1.0 - consequentSupport) / (1.0 - confidence)
				}

				rules = append(rules, AssociationRule{
					Antecedent:       antecedent,
					Consequent:       consequent,
					Support:          itemset.Support,
					Confidence:       confidence,
					Lift:             lift,
					LeverageMetric:   leverage,
					ConvictionMetric: conviction,
				})
			}
		}
	}

	return rules
}

// generateAllSubsets generates all non-empty subsets of a set
func generateAllSubsets(set []string) [][]string {
	n := len(set)
	count := 1 << uint(n) // 2^n
	result := make([][]string, 0, count-1)

	// For each possible subset (excluding empty set)
	for i := 1; i < count; i++ {
		subset := make([]string, 0)
		for j := 0; j < n; j++ {
			if (i & (1 << uint(j))) > 0 {
				subset = append(subset, set[j])
			}
		}
		result = append(result, subset)
	}

	return result
}

// containsItem checks if a transaction contains an item
func containsItem(transaction Transaction, item string) bool {
	for _, t := range transaction {
		if t == item {
			return true
		}
	}
	return false
}

// isSubset checks if items is a subset of transaction
func isSubset(items []string, transaction Transaction) bool {
	for _, item := range items {
		if !containsItem(transaction, item) {
			return false
		}
	}
	return true
}

// slicesEqual checks if two string slices are equal
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// difference returns the elements in a that are not in b
func difference(a, b []string) []string {
	result := make([]string, 0)

	// Create a map for quick lookup
	bMap := make(map[string]bool)
	for _, item := range b {
		bMap[item] = true
	}

	// Add items from a that are not in b
	for _, item := range a {
		if !bMap[item] {
			result = append(result, item)
		}
	}

	return result
}

// saveRulesToCSV saves association rules to a CSV file
func saveRulesToCSV(rules []AssociationRule, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"antecedents", "consequents", "support", "confidence", "lift", "leverage", "conviction"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Write rules
	for _, rule := range rules {
		antecedentStr := "{" + strings.Join(rule.Antecedent, ",") + "}"
		consequentStr := "{" + strings.Join(rule.Consequent, ",") + "}"

		conviction := fmt.Sprintf("%.6f", rule.ConvictionMetric)
		if math.IsInf(rule.ConvictionMetric, 1) {
			conviction = "inf"
		}

		record := []string{
			antecedentStr,
			consequentStr,
			fmt.Sprintf("%.6f", rule.Support),
			fmt.Sprintf("%.6f", rule.Confidence),
			fmt.Sprintf("%.6f", rule.Lift),
			fmt.Sprintf("%.6f", rule.LeverageMetric),
			conviction,
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing rule: %v", err)
		}
	}

	return nil
}

// saveItemsetsToCSV saves frequent itemsets to a CSV file
func saveItemsetsToCSV(itemsets []FrequentItemset, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"support", "itemsets", "length"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Write itemsets
	for _, itemset := range itemsets {
		itemsetStr := "{" + strings.Join(itemset.Items, ",") + "}"

		record := []string{
			fmt.Sprintf("%.6f", itemset.Support),
			itemsetStr,
			fmt.Sprintf("%d", itemset.Length),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing itemset: %v", err)
		}
	}

	return nil
}

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <csv_file> [min_support] [min_confidence] [max_length]")
		fmt.Println("  - csv_file: Path to the CSV file with columns for Basket and Item")
		fmt.Println("  - min_support: Minimum support threshold (default: 0.01)")
		fmt.Println("  - min_confidence: Minimum confidence threshold (default: 0.2)")
		fmt.Println("  - max_length: Maximum itemset length (default: 5)")
		os.Exit(1)
	}

	// Get input file
	inputFile := os.Args[1]

	// Set parameters with defaults
	minSupport := 0.01
	minConfidence := 0.2
	maxLen := 5

	// Override from command line if provided
	if len(os.Args) > 2 {
		_, err := fmt.Sscanf(os.Args[2], "%f", &minSupport)
		if err != nil {
			log.Fatalf("Invalid min_support value: %v", err)
		}
	}

	if len(os.Args) > 3 {
		_, err := fmt.Sscanf(os.Args[3], "%f", &minConfidence)
		if err != nil {
			log.Fatalf("Invalid min_confidence value: %v", err)
		}
	}

	if len(os.Args) > 4 {
		_, err := fmt.Sscanf(os.Args[4], "%d", &maxLen)
		if err != nil {
			log.Fatalf("Invalid max_length value: %v", err)
		}
	}

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		log.Fatalf("Input file %s does not exist", inputFile)
	}

	// Start execution
	fmt.Println("Starting Apriori algorithm...")
	fmt.Printf("Input file: %s\n", inputFile)
	fmt.Printf("Parameters: minSupport=%.4f, minConfidence=%.4f, maxLen=%d\n",
		minSupport, minConfidence, maxLen)

	// Load data
	fmt.Println("Loading and transforming dataset...")
	startLoadTime := time.Now()
	dataset, err := loadTransactionsFromCSV(inputFile)
	if err != nil {
		log.Fatalf("Error loading dataset: %v", err)
	}

	fmt.Printf("Dataset loaded in %v\n", time.Since(startLoadTime))
	fmt.Printf("Found %d transactions and %d unique items\n",
		len(dataset.Transactions), len(dataset.UniqueItems))

	// Find frequent itemsets
	fmt.Println("Finding frequent itemsets...")
	startItemsetTime := time.Now()
	frequentItemsets := findFrequentItemsets(dataset, minSupport, maxLen)
	itemsetTime := time.Since(startItemsetTime)

	fmt.Printf("Found %d frequent itemsets in %v\n", len(frequentItemsets), itemsetTime)

	// Print frequent itemsets by length
	lengths := make(map[int]int)
	for _, itemset := range frequentItemsets {
		lengths[itemset.Length]++
	}

	for k, v := range lengths {
		fmt.Printf("  Length %d: %d itemsets\n", k, v)
	}

	// Generate association rules
	fmt.Println("Generating association rules...")
	startRuleTime := time.Now()
	rules := generateAssociationRules(frequentItemsets, minConfidence)
	ruleTime := time.Since(startRuleTime)

	fmt.Printf("Generated %d association rules in %v\n", len(rules), ruleTime)

	// Save results
	itemsetsFile := "frequent_itemsets.csv"
	rulesFile := "association_rules.csv"

	fmt.Println("Saving results to files...")
	if err := saveItemsetsToCSV(frequentItemsets, itemsetsFile); err != nil {
		log.Fatalf("Error saving itemsets: %v", err)
	}

	if err := saveRulesToCSV(rules, rulesFile); err != nil {
		log.Fatalf("Error saving rules: %v", err)
	}

	fmt.Printf("Frequent itemsets saved to %s\n", itemsetsFile)
	fmt.Printf("Association rules saved to %s\n", rulesFile)
	fmt.Printf("Total execution time: %v\n", time.Since(startLoadTime))
}
