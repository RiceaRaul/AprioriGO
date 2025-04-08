package algorithm

import (
	"math"
	"sort"
	"strings"

	"github.com/RiceaRaul/AprioriGO/internal/models"
)

// FindFrequentItemsets finds frequent itemsets using the Apriori algorithm
func FindFrequentItemsets(dataset *models.Dataset, minSupport float64, maxLen int) []models.FrequentItemset {
	transactionCount := float64(len(dataset.Transactions))
	result := make([]models.FrequentItemset, 0)

	// Find frequent 1-itemsets
	L1 := make([]models.FrequentItemset, 0)
	for _, item := range dataset.UniqueItems {
		count := 0
		for _, transaction := range dataset.Transactions {
			if containsItem(transaction, item) {
				count++
			}
		}

		support := float64(count) / transactionCount
		if support >= minSupport {
			L1 = append(L1, models.FrequentItemset{
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

		Lk := make([]models.FrequentItemset, 0)
		for _, candidate := range Ck {
			count := 0
			for _, transaction := range dataset.Transactions {
				if isSubset(candidate.Items, transaction) {
					count++
				}
			}

			support := float64(count) / transactionCount
			if support >= minSupport {
				Lk = append(Lk, models.FrequentItemset{
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

// generateCandidates generates candidate itemsets of size k from frequent itemsets of size k-1
func generateCandidates(itemsets []models.FrequentItemset, k int) []models.FrequentItemset {
	candidates := make([]models.FrequentItemset, 0)

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

			candidates = append(candidates, models.FrequentItemset{
				Items:  candidate,
				Length: k,
			})
		}
	}

	return candidates
}

// GenerateAssociationRules generates association rules from frequent itemsets
func GenerateAssociationRules(itemsets []models.FrequentItemset, minConfidence float64) []models.AssociationRule {
	rules := make([]models.AssociationRule, 0)
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

				rules = append(rules, models.AssociationRule{
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
