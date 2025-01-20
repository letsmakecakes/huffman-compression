package models

// FrequencyMap is a map that stores the frequency of each type (character).
type FrequencyMap map[byte]uint64

// Header represents the file header for compressed files,
// containing metadata such as version, character frequencies, and data size.
type Header struct {
	Version  uint8        // Version of the compressed file format
	FreqMap  FrequencyMap // Map of character frequencies
	DataSize uint64       // Size of the original uncompressed data
}
