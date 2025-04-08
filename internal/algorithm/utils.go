package algorithm

import (
	"github.com/RiceaRaul/AprioriGO/internal/models"
)

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
func containsItem(transaction models.Transaction, item string) bool {
	for _, t := range transaction {
		if t == item {
			return true
		}
	}
	return false
}

// isSubset checks if items is a subset of transaction
func isSubset(items []string, transaction models.Transaction) bool {
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
