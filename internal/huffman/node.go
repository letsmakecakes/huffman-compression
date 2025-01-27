package huffman

// Node represents a single node in a Huffman tree.
type Node struct {
	Char        byte   // The character stored in the node
	Frequency   uint64 // The frequency of the character
	Left, Right *Node  // Pointers to the left and right child nodes
}

// PriorityQueue implements a priority queue for Nodes based on their frequency.
type PriorityQueue []*Node

// Len returns the length of the priority queue.
func (pq PriorityQueue) Len() int {
	return len(pq)
}

// Less compares the frequency of two nodes and returns true if the first node has a lower frequency.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Frequency < pq[j].Frequency
}

// Swap swaps the positions of two nodes in the priority queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push adds a new node to the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

// Pop removes and returns the node with the highest priority (the lowest frequency)
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
