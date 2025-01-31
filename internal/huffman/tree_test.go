package huffman

import (
	"huffman-compression/internal/models"
	"reflect"
	"sort"
	"testing"
)

// TestCase represents a test case for Huffman tree construction
type TestCase struct {
	name        string
	frequencies models.FrequencyMap
	expected    map[byte]string
}

func TestBuildTree(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t2 *testing.T)
	}{
		{"Basic Frequencies", testBasicFrequencies},
		{"Single Character", testSingleCharacter},
		{"Empty Frequencies", testEmptyFrequencies},
		{"Code Lengths", testCodeLengths},
		{"Specific Codes", testSpecificCodes},
		{"Tree Structure", testTreeStructure},
		{"Frequency Preservation", testFrequencyPreservation},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

func testBasicFrequencies(t *testing.T) {
	frequencies := models.FrequencyMap{
		'a': 5, 'b': 2, 'c': 1, 'd': 3,
	}

	tree := BuildTree(frequencies)

	// Verify code table size
	assertCodeTableSize(t, tree, len(frequencies))

	// Verify prefix property
	assertNoPrefixViolations(t, tree.CodeTable)
}

func testSingleCharacter(t *testing.T) {
	tree := BuildTree(models.FrequencyMap{'a': 5})

	assertCodeTableSize(t, tree, 1)
	assertCodeValue(t, tree, 'a', "0")
}

func testEmptyFrequencies(t *testing.T) {
	tree := BuildTree(models.FrequencyMap{})

	if tree.Root != nil {
		t.Error("Expected nil root for empty frequencies")
	}
	assertCodeTableSize(t, tree, 0)
}

func testCodeLengths(t *testing.T) {
	frequencies := models.FrequencyMap{
		'a': 64, 'b': 32, 'c': 16,
		'd': 8, 'e': 4, 'f': 2,
	}

	tree := BuildTree(frequencies)
	assertCodeLengthsAreOptimal(t, tree, frequencies)
}

func testSpecificCodes(t *testing.T) {
	testCases := []TestCase{
		{
			name:        "Simple case",
			frequencies: models.FrequencyMap{'a': 2, 'b': 1},
			expected:    map[byte]string{'a': "0", 'b': "1"},
		},
		{
			name:        "Three characters",
			frequencies: models.FrequencyMap{'a': 4, 'b': 2, 'c': 1},
			expected:    map[byte]string{'a': "0", 'b': "10", 'c': "11"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree := BuildTree(tc.frequencies)
			assertCodeTableEquals(t, tree.CodeTable, tc.expected)
		})
	}
}

func testTreeStructure(t *testing.T) {
	frequencies := models.FrequencyMap{'a': 5, 'b': 2, 'c': 1}
	tree := BuildTree(frequencies)
	verifyTreeStructure(t, tree.Root, frequencies)
}

func testFrequencyPreservation(t *testing.T) {
	frequencies := models.FrequencyMap{'a': 10, 'b': 5, 'c': 3}
	tree := BuildTree(frequencies)
	assertFrequencyPreservation(t, tree, frequencies)
}

// Benchmark tests
func BenchmarkBuildTree(b *testing.B) {
	frequencies := models.FrequencyMap{
		'a': 10000, 'b': 5000, 'c': 3000, 'd': 2000,
		'e': 1000, 'f': 500, 'g': 100, 'h': 50,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildTree(frequencies)
	}
}

// Helper functions
func assertCodeTableSize(t *testing.T, tree *Tree, expectedSize int) {
	t.Helper()
	if len(tree.CodeTable) != expectedSize {
		t.Errorf("code table size mismatch. Expected %d, got %d", expectedSize, len(tree.CodeTable))
	}
}

func assertCodeValue(t *testing.T, tree *Tree, char byte, expectedCode string) {
	t.Helper()
	if code, exists := tree.CodeTable[char]; !exists || code != expectedCode {
		t.Errorf("Expected code '%s' for character '%c', got '%s'", expectedCode, char, code)
	}
}

func assertNoPrefixViolations(t *testing.T, codeTable map[byte]string) {
	t.Helper()
	codes := extractCodes(codeTable)
	for i, code1 := range codes {
		for j, code2 := range codes {
			if i != j && isPrefixOf(code1, code2) {
				t.Errorf("Code %s is a prefix of %s", code1, code2)
			}
		}
	}
}

func assertCodeLengthsAreOptimal(t *testing.T, tree *Tree, frequencies models.FrequencyMap) {
	t.Helper()
	codes := getSortedCodes(tree.CodeTable, frequencies)

	for i := 1; i < len(codes); i++ {
		if len(codes[i-1].code) > len(codes[i].code) {
			t.Errorf("Higher frequency character %c has longer code than lower frequency character %c",
				codes[i-1].char, codes[i].char)
		}
	}
}

func assertCodeTableEquals(t *testing.T, got, want map[byte]string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Code table mismatch.\nExpected: %v\nGot: %v", want, got)
	}
}

func assertFrequencyPreservation(t *testing.T, tree *Tree, frequencies models.FrequencyMap) {
	t.Helper()
	var totalFreq uint64
	for _, freq := range frequencies {
		totalFreq += freq
	}

	if tree.Root.Frequency != totalFreq {
		t.Errorf("Root frequency mismatch. Expected %d, got %d", totalFreq, tree.Root.Frequency)
	}
}

type charCode struct {
	char byte
	code string
}

func extractCodes(codeTable map[byte]string) []string {
	codes := make([]string, 0, len(codeTable))
	for _, code := range codeTable {
		codes = append(codes, code)
	}
	return codes
}

func getSortedCodes(codeTable map[byte]string, frequencies models.FrequencyMap) []charCode {
	codes := make([]charCode, 0, len(codeTable))
	for char, code := range codeTable {
		codes = append(codes, charCode{char, code})
	}

	sort.Slice(codes, func(i, j int) bool {
		return frequencies[codes[i].char] > frequencies[codes[j].char]
	})

	return codes
}

func isPrefixOf(s1, s2 string) bool {
	return len(s1) <= len(s2) && s2[:len(s1)] == s1
}

func verifyTreeStructure(t *testing.T, node *Node, frequencies models.FrequencyMap) uint64 {
	t.Helper()
	if node == nil {
		return 0
	}

	if node.isLeaf() {
		expectedFreq, exists := frequencies[node.Char]
		if !exists {
			t.Errorf("Unexpected character %c in tree", node.Char)
		}
		if node.Frequency != expectedFreq {
			t.Errorf("Frequency mismatch for character %c. Expected %d, got %d", node.Char, expectedFreq, node.Frequency)
		}
		return node.Frequency
	}

	leftFreq := verifyTreeStructure(t, node.Left, frequencies)
	rightFreq := verifyTreeStructure(t, node.Right, frequencies)

	if node.Frequency != leftFreq+rightFreq {
		t.Errorf("Internal node frequency mismatch. Expectd %d, got %d", leftFreq+rightFreq, node.Frequency)
	}

	return node.Frequency
}
