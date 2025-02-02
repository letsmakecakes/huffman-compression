package frequency

import (
	"bytes"
	"huffman-compression/internal/models"
	"os"
	"reflect"
	"testing"
)

// createTestFile is a helper function to create temporary test files.
func createTestFile(filename string, data []byte, t *testing.T) {
	t.Helper() // Marks this function as a test helper
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		t.Fatalf("failed to create test file %s: %v", filename, err)
	}
	t.Cleanup(func() { os.Remove(filename) }) // Ensures cleanup after test
}

// TestCounter_Count checks if the frequency counter works correctly.
func TestCounter_Count(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     models.FrequencyMap
		wantErr  bool
		fileData []byte
	}{
		{
			name:     "basic test with simple string",
			input:    "test.txt",
			want:     models.FrequencyMap{'h': 1, 'e': 1, 'l': 2, 'o': 1},
			fileData: []byte("hello"),
		},
		{
			name:     "empty file",
			input:    "empty.txt",
			want:     models.FrequencyMap{},
			fileData: []byte(""),
		},
		{
			name:     "repeated characters",
			input:    "repeated.txt",
			want:     models.FrequencyMap{'a': 5, 'b': 3, 'c': 1},
			fileData: []byte("aaaaabbbc"),
		},
		{
			name:     "large chunk test",
			input:    "large.txt",
			want:     models.FrequencyMap{'x': 8192},
			fileData: bytes.Repeat([]byte{'x'}, 8192),
		},
		{
			name:    "non-existent file",
			input:   "nonexistent.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test file if needed
			if !tt.wantErr && tt.fileData != nil {
				createTestFile(tt.input, tt.fileData, t)
			}

			// Create counter-instance
			c := New()

			// Run the test
			got, err := c.Count(tt.input)

			// Check error expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("Counter.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Skip result comparison if expecting an error
			if tt.wantErr {
				return
			}

			// Compare results
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Counter.Count() = %v, want = %v", got, tt.want)
			}
		})
	}
}

// TestCounter_MergeFrequencies checks if frequency maps merge correctly.
func TestCounter_MergeFrequencies(t *testing.T) {
	tests := []struct {
		name     string
		initial  models.FrequencyMap
		toMerge  models.FrequencyMap
		expected models.FrequencyMap
	}{
		{"merge empty maps", models.FrequencyMap{}, models.FrequencyMap{}, models.FrequencyMap{}},
		{"merge with empty map",
			models.FrequencyMap{'a': 1, 'b': 2}, models.FrequencyMap{},
			models.FrequencyMap{'a': 1, 'b': 2}},
		{"merge overlapping maps",
			models.FrequencyMap{'a': 1, 'b': 2},
			models.FrequencyMap{'b': 3, 'c': 4},
			models.FrequencyMap{'a': 1, 'b': 5, 'c': 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Counter{frequencies: tt.initial}
			c.mergeFrequencies(tt.toMerge)

			if !reflect.DeepEqual(c.frequencies, tt.expected) {
				t.Errorf("mergeFrequencies() got = %v, want %v", c.frequencies, tt.expected)
			}
		})
	}
}

// BenchmarkCounter_Count measures performance of the frequency counter.
func BenchmarkCounter_Count(b *testing.B) {
	filename := "benchmark_test.txt"
	size := 1024 * 1024 // 1MB
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}

	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		b.Errorf("failed to create benchmark file: %v", err)
		return
	}
	defer os.Remove(filename)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c := New()
		_, err := c.Count(filename)
		if err != nil {
			b.Fatalf("Counter.Count() error = %v", err)
		}
	}
}
