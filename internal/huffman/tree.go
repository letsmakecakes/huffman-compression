package huffman

import (
	"container/heap"
	"huffman-compression/internal/models"
)

// Tree represents a Huffman tree with a root node and a code table for encoding characters.
type Tree struct {
	Root      *Node
	CodeTable map[byte]string
}

// BuildTree constructs a Huffman tree from a frequency map of characters.
func BuildTree(frequencies models.FrequencyMap) *Tree {
	if len(frequencies) == 0 {
		return &Tree{
			Root:      nil,
			CodeTable: make(map[byte]string),
		}
	}

	// Handle a single character case
	if len(frequencies) == 1 {
		var char byte
		var freq uint64
		for c, f := range frequencies {
			char = c
			freq = f
			break // Only one character, so we can break after the first iteration
		}
		tree := &Tree{
			Root: &Node{
				Char:      char,
				Frequency: freq,
			},
			CodeTable: map[byte]string{char: "0"}, // Single character gets code "0"
		}
		return tree
	}

	pq := initializePriorityQueue(frequencies)
	root := constructTree(pq)

	tree := &Tree{
		Root:      root,
		CodeTable: make(map[byte]string, len(frequencies)),
	}

	tree.generateCodeTable()
	return tree
}

// initializePriorityQueue initializes the priority queue with leaf nodes.
func initializePriorityQueue(frequencies models.FrequencyMap) PriorityQueue {
	pq := make(PriorityQueue, 0, len(frequencies))
	heap.Init(&pq)

	for char, freq := range frequencies {
		node := &Node{
			Char:      char,
			Frequency: freq,
		}
		heap.Push(&pq, node)
	}

	return pq
}

// constructTree builds the Huffman tree from a priority queue
func constructTree(pq PriorityQueue) *Node {
	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*Node)
		right := heap.Pop(&pq).(*Node)

		parent := &Node{
			Frequency: left.Frequency + right.Frequency,
			Left:      left,
			Right:     right,
		}

		heap.Push(&pq, parent)
	}

	if pq.Len() == 0 {
		return nil
	}

	return heap.Pop(&pq).(*Node)
}

// generateCodeTable creates the encoding table for the Huffman tree
func (t *Tree) generateCodeTable() {
	if t.Root == nil {
		return
	}
	t.traverseTree(t.Root, "")
}

func (t *Tree) traverseTree(node *Node, code string) {
	if node == nil {
		return
	}

	if node.isLeaf() {
		t.CodeTable[node.Char] = code
		return
	}

	// Always assign '0' to higher frequency path (left)
	// and '1' to lower frequency path (right)
	if node.Left != nil && node.Right != nil &&
		node.Left.Frequency < node.Right.Frequency {
		// Swap left and right if left has lower frequency
		node.Left, node.Right = node.Right, node.Left
	}

	// Traverse left
	t.traverseTree(node.Left, code+"0")

	// Traverse right with '1'
	t.traverseTree(node.Right, code+"1")
}
