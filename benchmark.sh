#!/bin/bash

echo "Apriori Algorithm Benchmark Suite"
echo "==============================="

DATA_FILE=$1
if [ -z "$DATA_FILE" ]; then
    echo "ERROR: Please provide a data file path"
    echo "Usage: ./batch_run.sh [data_file_path]"
    exit 1
fi

if [ ! -f "$DATA_FILE" ]; then
    echo "ERROR: Data file not found: $DATA_FILE"
    exit 1
fi

echo "Starting benchmark with data file: $DATA_FILE"
echo

echo "Step 1: Running benchmark with multiple parameter combinations..."
./apriori "$DATA_FILE" 0.01 0.3 3
echo
echo "Basic run completed."
echo

echo "Step 2: Running comprehensive benchmark..."
./benchmark "$DATA_FILE" benchmark_results.csv
echo
echo "Benchmark completed. Results saved to benchmark_results.csv"
echo

echo "Step 3: Analyzing benchmark results..."
./visualize benchmark_results.csv
echo

echo "Step 4: Recommended optimal parameters:"
echo "(These parameters are suggested based on your specific dataset)"
echo
echo "  Support:    See the \"OPTIMAL CONFIGURATIONS\" section above"
echo "  Confidence: See the \"OPTIMAL CONFIGURATIONS\" section above"
echo "  Max Length: See the \"OPTIMAL CONFIGURATIONS\" section above"
echo

echo "==============================="
echo "Benchmark suite completed!"
echo
echo "To run with optimal parameters:"
echo "  ./apriori $DATA_FILE [optimal_support] [optimal_confidence] [optimal_max_length]"
echo

# No direct equivalent to 'pause' in bash, but you can use this:
read -p "Press Enter to continue..."