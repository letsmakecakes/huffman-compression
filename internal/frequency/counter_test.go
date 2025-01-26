package frequency

import (
	"bytes"
	"testing"
)

type testCase struct {
	name     string
	input    string
	expected map[byte]uint64
}

func TestCounter_Count(t *testing.T) {
	testCases := []testCase{
		{
			name:  "Simple text",
			input: "hello world",
			expected: map[byte]uint64{
				'h': 1, 'e': 1, 'l': 3, 'o': 2, ' ': 1, 'w': 1, 'r': 1, 'd': 1,
			},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: map[byte]uint64{},
		},
		{
			name:  "Repeated characters",
			input: "aaabbbccc",
			expected: map[byte]uint64{
				'a': 3, 'b': 3, 'c': 3,
			},
		},
		{
			name:  "Mixed characters",
			input: "Hello, World! 123",
			expected: map[byte]uint64{
				'H': 1, 'e': 1, 'l': 3, 'o': 2, ',': 1, ' ': 2, 'W': 1, 'r': 1, 'd': 1, '!': 1, '1': 1, '2': 1, '3': 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validateFrequencies(t, tc)
		})
	}
}

func validateFrequencies(t *testing.T, tc testCase) {
	t.Helper()

	counter := NewCounter()
	reader := bytes.NewBufferString(tc.input)

	if err := counter.Count(reader); err != nil {
		t.Fatalf("Failed to count frequencies: %v", err)
	}

	frequencies := counter.GetFrequencies()
	validateFrequencyCount(t, frequencies, tc.expected)
	validateCharacterFrequencies(t, frequencies, tc.expected)
}

func validateFrequencyCount(t *testing.T, got, want map[byte]uint64) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("Frequency count mismatch: want %d, got %d", len(want), len(got))
	}
}

func validateCharacterFrequencies(t *testing.T, got, want map[byte]uint64) {
	t.Helper()

	for char, expectedCount := range want {
		actualCount, exists := got[char]
		if !exists {
			t.Errorf("Character %q not found in frequencies", char)
			continue
		}

		if actualCount != expectedCount {
			t.Errorf("Frequency mismatch for %q: want %d, got %d", char, expectedCount, actualCount)
		}
	}
}

type benchmarkCase struct {
	name  string
	input string
}

func BenchmarkCounter_Count(b *testing.B) {
	testData := []benchmarkCase{
		{
			name:  "Small text",
			input: "hello world",
		},
		{
			name:  "Medium text",
			input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		},
		{
			name:  "Large text",
			input: string(bytes.Repeat([]byte("abcdefghijklmnopqrstuvwxyz "), 1000)),
		},
	}

	for _, bc := range testData {
		b.Run(bc.name, func(b *testing.B) {
			runCountBenchmark(b, bc.input)
		})
	}
}

func runCountBenchmark(b *testing.B, input string) {
	b.Helper()

	for i := 0; i < b.N; i++ {
		counter := NewCounter()
		reader := bytes.NewBufferString(input)

		if err := counter.Count(reader); err != nil {
			b.Fatalf("Failed to count frequencies: %v", err)
		}
	}
}

func TestCounter_MultipleReads(t *testing.T) {
	const input = "Multiple reads test"
	counter := NewCounter()

	// Split input for multiple read operations
	chunks := []string{
		input[:10],
		input[10:],
	}

	for i, chunk := range chunks {
		reader := bytes.NewBufferString(chunk)
		if err := counter.Count(reader); err != nil {
			t.Fatalf("Failed to process chunk %d: %v", i+1, err)
		}
	}

	validateTotalFrequencies(t, counter.GetFrequencies(), input)
}

func validateTotalFrequencies(t *testing.T, frequencies map[byte]uint64, input string) {
	t.Helper()

	var total uint64
	for _, count := range frequencies {
		total += count
	}

	expected := uint64(len(input))
	if total != expected {
		t.Errorf("Total frequency mismatch: want %d, got %d", expected, total)
	}
}
