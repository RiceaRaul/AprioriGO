package output

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/RiceaRaul/apriori/internal/models"
)

// SaveRulesToCSV saves association rules to a CSV file
func SaveRulesToCSV(rules []models.AssociationRule, filePath string) error {
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

// SaveItemsetsToCSV saves frequent itemsets to a CSV file
func SaveItemsetsToCSV(itemsets []models.FrequentItemset, filePath string) error {
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
