package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/RiceaRaul/AprioriGO/internal/algorithm"
	"github.com/RiceaRaul/AprioriGO/internal/loader"
	"github.com/RiceaRaul/AprioriGO/internal/models"
)

type BenchmarkResult struct {
	MinSupport    float64
	MinConfidence float64
	MaxLength     int
	LoadTime      time.Duration
	ItemsetTime   time.Duration
	RuleTime      time.Duration
	TotalTime     time.Duration
	ItemsetCount  int
	RuleCount     int
	Memory        uint64 // in bytes
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: benchmark <csv_file> [output_file]")
		fmt.Println("  - csv_file: Path to the CSV file with transaction data")
		fmt.Println("  - output_file: Optional path to save benchmark results (default: benchmark_results.csv)")
		os.Exit(1)
	}

	// Get input file
	inputFile := os.Args[1]

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		log.Fatalf("Input file %s does not exist", inputFile)
	}

	// Set output file
	outputFile := "benchmark_results.csv"
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}

	// Create CPU profile if needed (uncomment to enable)
	// cpuProfile, err := os.Create("cpu_profile.prof")
	// if err != nil {
	// 	log.Fatal("Could not create CPU profile: ", err)
	// }
	// defer cpuProfile.Close()
	// if err := pprof.StartCPUProfile(cpuProfile); err != nil {
	// 	log.Fatal("Could not start CPU profile: ", err)
	// }
	// defer pprof.StopCPUProfile()

	// Parameter combinations to test
	minSupports := []float64{0.001, 0.005, 0.01, 0.02, 0.05}
	minConfidences := []float64{0.1, 0.2, 0.3, 0.5, 0.7}
	maxLengths := []int{2, 3, 4, 5}

	// Prepare results slice
	results := make([]BenchmarkResult, 0)

	// Load dataset once
	fmt.Println("Loading dataset...")
	dataset, err := loader.LoadFromCSV(inputFile)
	if err != nil {
		log.Fatalf("Error loading dataset: %v", err)
	}
	fmt.Printf("Dataset loaded with %d transactions and %d unique items\n\n",
		len(dataset.Transactions), len(dataset.UniqueItems))

	// Format for output
	fmt.Printf("%-10s %-10s %-10s %-15s %-15s %-15s %-10s %-10s\n",
		"Support", "Confidence", "MaxLen", "Itemset Time", "Rule Time", "Total Time", "Itemsets", "Rules")
	fmt.Println(strings.Repeat("-", 100))

	// Run the benchmark for each parameter combination
	for _, minSupport := range minSupports {
		for _, minConfidence := range minConfidences {
			for _, maxLength := range maxLengths {
				// Skip combinations that are likely to be too slow or memory-intensive
				if minSupport < 0.005 && maxLength > 3 {
					continue
				}

				fmt.Printf("Testing: support=%.4f, confidence=%.4f, maxLength=%d\n",
					minSupport, minConfidence, maxLength)

				// Run the benchmark
				result := runBenchmark(dataset, minSupport, minConfidence, maxLength)
				results = append(results, result)

				// Format output
				fmt.Printf("%-10.4f %-10.4f %-10d %-15s %-15s %-15s %-10d %-10d\n",
					minSupport, minConfidence, maxLength,
					formatDuration(result.ItemsetTime),
					formatDuration(result.RuleTime),
					formatDuration(result.TotalTime),
					result.ItemsetCount,
					result.RuleCount)

				// Force garbage collection to prevent memory buildup
				runtime.GC()
			}
		}
	}

	// Save results to CSV
	if err := saveResultsToCSV(results, outputFile); err != nil {
		log.Fatalf("Error saving results: %v", err)
	}

	fmt.Printf("\nBenchmark completed. Results saved to %s\n", outputFile)

	// Create memory profile
	memProfile, err := os.Create("memory_profile.prof")
	if err != nil {
		log.Fatal("Could not create memory profile: ", err)
	}
	defer memProfile.Close()
	runtime.GC() // Run garbage collection before taking memory profile
	if err := pprof.WriteHeapProfile(memProfile); err != nil {
		log.Fatal("Could not write memory profile: ", err)
	}
}

func runBenchmark(dataset *models.Dataset, minSupport, minConfidence float64, maxLength int) BenchmarkResult {
	startTotal := time.Now()
	var itemsetCount, ruleCount int
	var itemsetTime, ruleTime time.Duration
	var memStats runtime.MemStats

	// Find frequent itemsets
	startItemset := time.Now()
	frequentItemsets := algorithm.FindFrequentItemsets(dataset, minSupport, maxLength)
	itemsetTime = time.Since(startItemset)
	itemsetCount = len(frequentItemsets)

	// Generate association rules
	startRule := time.Now()
	rules := algorithm.GenerateAssociationRules(frequentItemsets, minConfidence)
	ruleTime = time.Since(startRule)
	ruleCount = len(rules)

	// Get memory usage
	runtime.ReadMemStats(&memStats)

	return BenchmarkResult{
		MinSupport:    minSupport,
		MinConfidence: minConfidence,
		MaxLength:     maxLength,
		LoadTime:      0, // Dataset already loaded
		ItemsetTime:   itemsetTime,
		RuleTime:      ruleTime,
		TotalTime:     time.Since(startTotal),
		ItemsetCount:  itemsetCount,
		RuleCount:     ruleCount,
		Memory:        memStats.Alloc,
	}
}

func saveResultsToCSV(results []BenchmarkResult, outputFile string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(outputFile)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	}

	// Create output file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"min_support",
		"min_confidence",
		"max_length",
		"itemset_time_ms",
		"rule_time_ms",
		"total_time_ms",
		"itemset_count",
		"rule_count",
		"memory_usage_mb",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Write results
	for _, result := range results {
		row := []string{
			fmt.Sprintf("%.6f", result.MinSupport),
			fmt.Sprintf("%.6f", result.MinConfidence),
			fmt.Sprintf("%d", result.MaxLength),
			fmt.Sprintf("%d", result.ItemsetTime.Milliseconds()),
			fmt.Sprintf("%d", result.RuleTime.Milliseconds()),
			fmt.Sprintf("%d", result.TotalTime.Milliseconds()),
			fmt.Sprintf("%d", result.ItemsetCount),
			fmt.Sprintf("%d", result.RuleCount),
			fmt.Sprintf("%.2f", float64(result.Memory)/(1024*1024)), // Convert to MB
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing result: %v", err)
		}
	}

	return nil
}

func formatDuration(d time.Duration) string {
	if d.Seconds() < 1 {
		return fmt.Sprintf("%dms", d.Milliseconds())
	} else if d.Minutes() < 1 {
		return fmt.Sprintf("%.2fs", d.Seconds())
	} else {
		return fmt.Sprintf("%.1fm %.1fs", d.Minutes(), d.Seconds()-float64(int(d.Minutes()))*60)
	}
}
