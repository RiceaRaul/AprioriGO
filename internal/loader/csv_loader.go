package loader

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/RiceaRaul/AprioriGO/internal/models"
)

// LoadFromCSV loads transactions from a CSV file with basket and item columns
func LoadFromCSV(filePath string) (*models.Dataset, error) {
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
	dataset := &models.Dataset{
		Transactions: make([]models.Transaction, 0, len(basketMap)),
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
		transaction := make(models.Transaction, 0, len(uniqueItems))
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
