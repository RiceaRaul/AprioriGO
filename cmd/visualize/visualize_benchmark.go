package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Result represents a benchmark result
type Result struct {
	MinSupport    float64
	MinConfidence float64
	MaxLength     int
	ItemsetTime   int64
	RuleTime      int64
	TotalTime     int64
	ItemsetCount  int
	RuleCount     int
	MemoryUsage   float64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: visualize <benchmark_results.csv>")
		os.Exit(1)
	}

	// Get input file
	inputFile := os.Args[1]

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		log.Fatalf("Input file %s does not exist", inputFile)
	}

	// Load benchmark results
	results, err := loadBenchmarkResults(inputFile)
	if err != nil {
		log.Fatalf("Error loading results: %v", err)
	}

	// Visualize results
	visualizeResults(results)
}

func loadBenchmarkResults(filePath string) ([]Result, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data found in CSV")
	}

	// Skip header
	records = records[1:]

	results := make([]Result, 0, len(records))
	for i, record := range records {
		if len(record) < 9 {
			return nil, fmt.Errorf("invalid record at line %d: not enough fields", i+2)
		}

		minSupport, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min_support at line %d: %v", i+2, err)
		}

		minConfidence, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min_confidence at line %d: %v", i+2, err)
		}

		maxLength, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("invalid max_length at line %d: %v", i+2, err)
		}

		itemsetTime, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid itemset_time at line %d: %v", i+2, err)
		}

		ruleTime, err := strconv.ParseInt(record[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid rule_time at line %d: %v", i+2, err)
		}

		totalTime, err := strconv.ParseInt(record[5], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid total_time at line %d: %v", i+2, err)
		}

		itemsetCount, err := strconv.Atoi(record[6])
		if err != nil {
			return nil, fmt.Errorf("invalid itemset_count at line %d: %v", i+2, err)
		}

		ruleCount, err := strconv.Atoi(record[7])
		if err != nil {
			return nil, fmt.Errorf("invalid rule_count at line %d: %v", i+2, err)
		}

		memoryUsage, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid memory_usage at line %d: %v", i+2, err)
		}

		results = append(results, Result{
			MinSupport:    minSupport,
			MinConfidence: minConfidence,
			MaxLength:     maxLength,
			ItemsetTime:   itemsetTime,
			RuleTime:      ruleTime,
			TotalTime:     totalTime,
			ItemsetCount:  itemsetCount,
			RuleCount:     ruleCount,
			MemoryUsage:   memoryUsage,
		})
	}

	return results, nil
}

func visualizeResults(results []Result) {
	// Print table header
	fmt.Println("\n===== BENCHMARK RESULTS SUMMARY =====")

	// Get unique parameter values for grouping
	supportsMap := make(map[float64]bool)
	confidencesMap := make(map[float64]bool)
	lengthsMap := make(map[int]bool)

	for _, r := range results {
		supportsMap[r.MinSupport] = true
		confidencesMap[r.MinConfidence] = true
		lengthsMap[r.MaxLength] = true
	}

	// 1. Analysis by support value
	fmt.Println("\n----- ANALYSIS BY SUPPORT VALUE -----")
	fmt.Printf("%-10s %-15s %-15s %-15s %-15s\n",
		"Support", "Avg Time (ms)", "Avg Itemsets", "Avg Rules", "Avg Memory (MB)")
	fmt.Println(strings.Repeat("-", 75))

	for support := range supportsMap {
		var count int
		var totalTime, totalItemsets, totalRules, totalMemory float64

		for _, r := range results {
			if r.MinSupport == support {
				count++
				totalTime += float64(r.TotalTime)
				totalItemsets += float64(r.ItemsetCount)
				totalRules += float64(r.RuleCount)
				totalMemory += r.MemoryUsage
			}
		}

		if count > 0 {
			fmt.Printf("%-10.4f %-15.1f %-15.1f %-15.1f %-15.2f\n",
				support,
				totalTime/float64(count),
				totalItemsets/float64(count),
				totalRules/float64(count),
				totalMemory/float64(count))
		}
	}

	// 2. Analysis by confidence value
	fmt.Println("\n----- ANALYSIS BY CONFIDENCE VALUE -----")
	fmt.Printf("%-10s %-15s %-15s %-15s %-15s\n",
		"Confidence", "Avg Time (ms)", "Avg Itemsets", "Avg Rules", "Avg Memory (MB)")
	fmt.Println(strings.Repeat("-", 75))

	for confidence := range confidencesMap {
		var count int
		var totalTime, totalItemsets, totalRules, totalMemory float64

		for _, r := range results {
			if r.MinConfidence == confidence {
				count++
				totalTime += float64(r.TotalTime)
				totalItemsets += float64(r.ItemsetCount)
				totalRules += float64(r.RuleCount)
				totalMemory += r.MemoryUsage
			}
		}

		if count > 0 {
			fmt.Printf("%-10.4f %-15.1f %-15.1f %-15.1f %-15.2f\n",
				confidence,
				totalTime/float64(count),
				totalItemsets/float64(count),
				totalRules/float64(count),
				totalMemory/float64(count))
		}
	}

	// 3. Analysis by max length
	fmt.Println("\n----- ANALYSIS BY MAX LENGTH -----")
	fmt.Printf("%-10s %-15s %-15s %-15s %-15s\n",
		"Max Length", "Avg Time (ms)", "Avg Itemsets", "Avg Rules", "Avg Memory (MB)")
	fmt.Println(strings.Repeat("-", 75))

	for length := range lengthsMap {
		var count int
		var totalTime, totalItemsets, totalRules, totalMemory float64

		for _, r := range results {
			if r.MaxLength == length {
				count++
				totalTime += float64(r.TotalTime)
				totalItemsets += float64(r.ItemsetCount)
				totalRules += float64(r.RuleCount)
				totalMemory += r.MemoryUsage
			}
		}

		if count > 0 {
			fmt.Printf("%-10d %-15.1f %-15.1f %-15.1f %-15.2f\n",
				length,
				totalTime/float64(count),
				totalItemsets/float64(count),
				totalRules/float64(count),
				totalMemory/float64(count))
		}
	}

	// 4. Find the top 5 fastest configurations
	fmt.Println("\n----- TOP 5 FASTEST CONFIGURATIONS -----")
	fmt.Printf("%-8s %-8s %-8s %-15s %-15s %-15s %-15s\n",
		"Support", "Conf", "MaxLen", "Time (ms)", "Itemsets", "Rules", "Memory (MB)")
	fmt.Println(strings.Repeat("-", 90))

	// Sort by execution time (bubble sort for simplicity)
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].TotalTime > results[j].TotalTime {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Print top 5 or fewer if less than 5 results
	count := 5
	if len(results) < 5 {
		count = len(results)
	}

	for i := 0; i < count; i++ {
		r := results[i]
		fmt.Printf("%-8.4f %-8.2f %-8d %-15d %-15d %-15d %-15.2f\n",
			r.MinSupport,
			r.MinConfidence,
			r.MaxLength,
			r.TotalTime,
			r.ItemsetCount,
			r.RuleCount,
			r.MemoryUsage)
	}

	// 5. Find optimal configurations (best balance between time and results)
	fmt.Println("\n----- OPTIMAL CONFIGURATIONS (TIME/RULES RATIO) -----")
	fmt.Printf("%-8s %-8s %-8s %-15s %-15s %-15s %-15s\n",
		"Support", "Conf", "MaxLen", "Time (ms)", "Itemsets", "Rules", "Ratio (ms/rule)")
	fmt.Println(strings.Repeat("-", 90))

	// Create a copy and sort by time/rules ratio (only consider cases with rules)
	optimalResults := make([]Result, 0)
	for _, r := range results {
		if r.RuleCount > 0 {
			optimalResults = append(optimalResults, r)
		}
	}

	// Sort by time/rules ratio (bubble sort for simplicity)
	for i := 0; i < len(optimalResults); i++ {
		for j := i + 1; j < len(optimalResults); j++ {
			ratio1 := float64(optimalResults[i].TotalTime) / float64(optimalResults[i].RuleCount)
			ratio2 := float64(optimalResults[j].TotalTime) / float64(optimalResults[j].RuleCount)
			if ratio1 > ratio2 {
				optimalResults[i], optimalResults[j] = optimalResults[j], optimalResults[i]
			}
		}
	}

	// Print top 5 or fewer if less than 5 results
	count = 5
	if len(optimalResults) < 5 {
		count = len(optimalResults)
	}

	for i := 0; i < count; i++ {
		r := optimalResults[i]
		ratio := float64(r.TotalTime) / float64(r.RuleCount)
		fmt.Printf("%-8.4f %-8.2f %-8d %-15d %-15d %-15d %-15.2f\n",
			r.MinSupport,
			r.MinConfidence,
			r.MaxLength,
			r.TotalTime,
			r.ItemsetCount,
			r.RuleCount,
			ratio)
	}
}
