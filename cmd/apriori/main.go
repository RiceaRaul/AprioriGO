package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RiceaRaul/apriori/internal/algorithm"
	"github.com/RiceaRaul/apriori/internal/loader"
	"github.com/RiceaRaul/apriori/internal/output"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: apriori <csv_file> [min_support] [min_confidence] [max_length]")
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
	dataset, err := loader.LoadFromCSV(inputFile)
	if err != nil {
		log.Fatalf("Error loading dataset: %v", err)
	}

	fmt.Printf("Dataset loaded in %v\n", time.Since(startLoadTime))
	fmt.Printf("Found %d transactions and %d unique items\n",
		len(dataset.Transactions), len(dataset.UniqueItems))

	// Find frequent itemsets
	fmt.Println("Finding frequent itemsets...")
	startItemsetTime := time.Now()
	frequentItemsets := algorithm.FindFrequentItemsets(dataset, minSupport, maxLen)
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
	rules := algorithm.GenerateAssociationRules(frequentItemsets, minConfidence)
	ruleTime := time.Since(startRuleTime)

	fmt.Printf("Generated %d association rules in %v\n", len(rules), ruleTime)

	// Save results
	itemsetsFile := "frequent_itemsets.csv"
	rulesFile := "association_rules.csv"

	fmt.Println("Saving results to files...")
	if err := output.SaveItemsetsToCSV(frequentItemsets, itemsetsFile); err != nil {
		log.Fatalf("Error saving itemsets: %v", err)
	}

	if err := output.SaveRulesToCSV(rules, rulesFile); err != nil {
		log.Fatalf("Error saving rules: %v", err)
	}

	fmt.Printf("Frequent itemsets saved to %s\n", itemsetsFile)
	fmt.Printf("Association rules saved to %s\n", rulesFile)
	fmt.Printf("Total execution time: %v\n", time.Since(startLoadTime))
}
