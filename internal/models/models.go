package models

// Transaction represents a set of items in a basket
type Transaction []string

// FrequentItemset represents a set of items that appear together with their support
type FrequentItemset struct {
	Items   []string
	Support float64
	Length  int
}

// AssociationRule represents a rule with antecedent -> consequent with metrics
type AssociationRule struct {
	Antecedent       []string
	Consequent       []string
	Support          float64
	Confidence       float64
	Lift             float64
	LeverageMetric   float64
	ConvictionMetric float64
}

// Dataset contains the transaction data and metadata
type Dataset struct {
	Transactions []Transaction
	UniqueItems  []string
	ItemsMap     map[string]bool
}
