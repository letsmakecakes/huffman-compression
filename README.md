# Huffman Compression Tool

A high-performance text file compression tool implemented in Go using Huffman coding. This tool provides lossless compression for text files by utilizing variable-length prefix coding based on character frequencies.

## Features

- Efficient text file compression using Huffman coding
- Lossless compression and decompression
- High-performance I/O with buffered operations
- Memory-efficient bit manipulation
- Comprehensive error handling
- Unit tested components
- Command-line interface
- Progress indication for large files

## Installation

Ensure you have Go 1.21 or later installed on your system.

```bash
# Clone the repository
git clone https://github.com/letsmakecakes/huffman-compression
cd huffman-compression

# Build the project
make build

# Run tests
make test
```

## Usage

The tool provides two main operations: compression and decompression.

### Basic Commands

```bash
# Compress a file
./huffman -input input.txt -output compressed.huf -mode compress

# Decompress a file
./huffman -input compressed.huf -output decompressed.txt -mode decompress
```

### Command Line Options

- `-input`: Path to the input file (required)
- `-output`: Path to the output file (required)
- `-mode`: Operation mode - "compress" or "decompress" (default: "compress")
- `-verbose`: Enable verbose logging (optional)

## Project Structure

```
huffman-compression/
├── cmd/                  # Command-line application
├── internal/             # Internal packages
│   ├── codec/           # Encoding/decoding logic
│   ├── frequency/       # Character frequency analysis
│   ├── huffman/         # Huffman tree implementation
│   ├── io/              # File I/O operations
│   └── models/          # Shared types and interfaces
├── pkg/                 # Public packages
│   └── bitutils/        # Bit manipulation utilities
└── test/                # Test files and integration tests
```

## How It Works

1. **Frequency Analysis**: The tool first analyzes the input file to count the frequency of each character.

2. **Tree Construction**: A Huffman tree is built using these frequencies, with more frequent characters placed closer to the root for shorter codes.

3. **Code Generation**: Variable-length prefix codes are generated for each character based on their position in the tree.

4. **Compression**:
    - The Huffman tree structure is written to the output file header
    - The input text is encoded using the generated codes
    - Bits are packed efficiently into bytes for storage

5. **Decompression**:
    - The header is read to reconstruct the Huffman tree
    - The compressed data is read bit by bit
    - The original text is reconstructed using the tree

## Performance

The tool implements several optimizations for better performance:

- Buffered I/O operations
- Efficient bit manipulation
- CPU cache-friendly data structures
- Minimal memory allocations
- Parallel processing for large files

## Example

```go
package main

import "fmt"

func main() {
    // Sample compression
    inputFile := "les-miserables.txt"
    compressedFile := "compressed.huf"
    
    // Compress
    if err := compress(inputFile, compressedFile); err != nil {
        fmt.Printf("Compression error: %v\n", err)
        return
    }
    
    // Show compression ratio
    PrintCompressionStats(inputFile, compressedFile)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. Make sure to:

1. Add tests for any new features
2. Follow the existing code style
3. Update documentation as needed
4. Add comments for non-obvious code sections

## Testing

Run the test suite:

```bash
make test
```

Run benchmarks:

```bash
make bench
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

This implementation is inspired by:
- David Huffman's original 1952 paper "A Method for the Construction of Minimum-Redundancy Codes"
- Modern performance optimization techniques for Go applications
- Best practices from the Go community

## Author

Adwaith Rajeev

## Acknowledgments

Special thanks to:
- The Go team for their excellent standard library
- The open source community for their valuable feedback
- Everyone who has contributed to the project

## Project Status

Active development - Bug reports and feature requests are welcome!