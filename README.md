# AprioriGO: High-Performance Association Rule Mining

[![Build Status](https://github.com/RiceaRaul/AprioriGO/actions/workflows/go.yml/badge.svg)](https://github.com/RiceaRaul/AprioriGO/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/RiceaRaul/AprioriGO)](https://goreportcard.com/report/github.com/RiceaRaul/AprioriGO)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Features

- **Modular Architecture**: Clean, maintainable codebase with separation of concerns
- **Multiple Metrics**: Support, confidence, lift, leverage, and conviction metrics
- **Comprehensive Benchmarking**: Tools to find optimal parameters for your dataset
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **CSV Integration**: Easy loading from and saving to CSV files

## Installation

### Prerequisites

- Go 1.18 or higher

### Building from Source

```bash
# Clone the repository
git clone https://github.com/RiceaRaul/AprioriGO.git
cd AprioriGO

# Build the main executable
go build -o apriori ./cmd/apriori

# Build benchmark utilities (optional)
go build -o benchmark ./cmd/benchmark
go build -o visualize ./cmd/visualize
```

## Quick Start

```bash
# Run with default parameters
./apriori your_data.csv

# Run with custom parameters
./apriori your_data.csv 0.01 0.3 4
```

Parameters:
- `your_data.csv`: Path to the CSV file with columns for Basket and Item
- `0.01`: Minimum support threshold (default: 0.01)
- `0.3`: Minimum confidence threshold (default: 0.2)
- `4`: Maximum itemset length (default: 5)

## Input Data Format

The algorithm expects a CSV file with at least two columns:
- First column: Basket/transaction ID
- Second column: Item name

Example:
```csv
basket_id,item
1001,bread
1001,milk
1001,eggs
1002,bread
1002,butter
```

## Output Files

Two CSV files are generated:

1. `frequent_itemsets.csv`:
   - support: The support value
   - itemsets: The set of items
   - length: Number of items in the set

2. `association_rules.csv`:
   - antecedents: The items on the left side of the rule
   - consequents: The items on the right side of the rule
   - support: Support of the entire itemset
   - confidence: Confidence of the rule
   - lift: Lift metric
   - leverage: Leverage metric
   - conviction: Conviction metric

## Advanced Usage

### Finding Optimal Parameters

For optimal performance with your specific dataset, use the benchmarking tools:

```bash
# Windows
benchmark.bat your_data.csv

# Linux/macOS
./batch_run.sh your_data.csv
```

This will:
1. Run with default parameters
2. Benchmark multiple parameter combinations
3. Analyze and visualize the results
4. Recommend optimal parameters

### Performance Considerations

- **Memory usage** scales with the number of frequent itemsets found. Lower support thresholds result in more itemsets and higher memory usage.
- **Execution time** is most affected by minimum support and maximum itemset length. Use the benchmark tool to find the sweet spot.
- For extremely large datasets, start with a higher support threshold and gradually decrease it.

## Project Structure

```
apriori/
├── cmd/                    # Application entry points
│   ├── apriori/            # Main application
│   ├── benchmark/          # Benchmarking tool
│   └── visualize/          # Results visualization
├── internal/               # Internal packages
│   ├── models/             # Data structures
│   ├── loader/             # Data loading
│   ├── algorithm/          # Algorithm implementation
│   └── output/             # Results output
├── scripts/                # Helper scripts
│   ├── batch_run.bat       # Windows benchmark script
│   └── batch_run.sh        # Linux/Mac benchmark script
├── LICENSE                 # License information
└── README.md               # This file
```

You can view the full source code at [github.com/RiceaRaul/AprioriGO](https://github.com/RiceaRaul/AprioriGO)

## Algorithm Details

This implementation uses the classic Apriori algorithm with the following steps:

1. Find all frequent 1-itemsets
2. Generate candidate k-itemsets from frequent (k-1)-itemsets
3. Count support for each candidate and retain those above the threshold
4. Repeat steps 2-3 until no more frequent itemsets are found
5. Generate association rules from the frequent itemsets

### Metrics Calculated

- **Support**: Fraction of transactions containing the itemset
- **Confidence**: Probability of finding the consequent given the antecedent
- **Lift**: Ratio of observed support to expected support if items were independent
- **Leverage**: Difference between observed and expected support
- **Conviction**: Measure of implication strength

## Performance Benchmarks

Real benchmark results from our comprehensive testing:

### Analysis by Support Value

| Support | Average Time (ms) | Average Itemsets | Average Rules | Average Memory (MB) |
|---------|------------------|------------------|---------------|---------------------|
| 0.001   | 2,293.0          | 3,895.0          | 1,985.3       | 3.67                |
| 0.005   | 771.3            | 681.8            | 339.0         | 3.00                |
| 0.01    | 360.4            | 252.8            | 102.5         | 2.00                |
| 0.02    | 167.8            | 103.0            | 33.8          | 1.56                |
| 0.05    | 46.5             | 28.0             | 2.8           | 1.36                |

### Analysis by Confidence Value

| Confidence | Average Time (ms) | Average Itemsets | Average Rules | Average Memory (MB) |
|------------|------------------|------------------|---------------|---------------------|
| 0.1        | 546.9            | 669.6            | 926.2         | 2.19                |
| 0.2        | 560.2            | 669.6            | 467.9         | 2.20                |
| 0.3        | 558.1            | 669.6            | 209.0         | 2.17                |
| 0.5        | 555.6            | 669.6            | 28.7          | 2.17                |
| 0.7        | 548.7            | 669.6            | 2.4           | 2.12                |

### Analysis by Max Length

| Max Length | Average Time (ms) | Average Itemsets | Average Rules | Average Memory (MB) |
|------------|------------------|------------------|---------------|---------------------|
| 2          | 356.0            | 620.2            | 139.6         | 2.15                |
| 3          | 1,057.9          | 1,345.0          | 820.2         | 2.38                |
| 4          | 361.5            | 278.2            | 135.6         | 2.00                |
| 5          | 363.6            | 278.2            | 135.6         | 2.09                |

### Top 5 Fastest Configurations

| Support | Confidence | Max Length | Time (ms) | Itemsets | Rules | Memory (MB) |
|---------|------------|------------|-----------|----------|-------|-------------|
| 0.05    | 0.1        | 2          | 44        | 28       | 6     | 1.36        |
| 0.05    | 0.3        | 4          | 45        | 28       | 2     | 1.36        |
| 0.05    | 0.7        | 3          | 45        | 28       | 0     | 1.36        |
| 0.05    | 0.7        | 5          | 45        | 28       | 0     | 1.36        |
| 0.05    | 0.1        | 4          | 46        | 28       | 6     | 1.36        |

### Optimal Configurations (Best Time/Rules Ratio)

| Support | Confidence | Max Length | Time (ms) | Itemsets | Rules | Ratio (ms/rule) |
|---------|------------|------------|-----------|----------|-------|-----------------|
| 0.001   | 0.1        | 3          | 3,813     | 5,612    | 9,754 | 0.39            |
| 0.001   | 0.1        | 2          | 684       | 2,178    | 1,449 | 0.47            |
| 0.005   | 0.1        | 3          | 815       | 725      | 1,103 | 0.74            |
| 0.001   | 0.2        | 3          | 3,928     | 5,612    | 5,094 | 0.77            |
| 0.005   | 0.1        | 4          | 860       | 725      | 1,103 | 0.78            |

*Testing system: AMD Ryzen™ 7 8845HS, 16GB RAM, Go 1.22.0*

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request to the [AprioriGO repository](https://github.com/RiceaRaul/AprioriGO)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- R. Agrawal and R. Srikant for the original Apriori algorithm paper
- The Go community for the excellent language and tools

---

*This implementation is designed for educational and research purposes. For critical production applications, consider additional testing and validation.*
