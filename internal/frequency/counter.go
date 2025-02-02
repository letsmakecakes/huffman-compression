package frequency

import (
	"bufio"
	"huffman-compression/internal/models"
	"os"
	"sync"
)

const (
	chunkSize   = 8192 // 8KB chunks for reading
	workerCount = 4    // Number of concurrent workers
)

// Counter handles character frequency counting.
type Counter struct {
	frequencies models.FrequencyMap
	mu          sync.Mutex
}

// New creates a new frequency counter.
func New() *Counter {
	return &Counter{
		frequencies: make(models.FrequencyMap),
	}
}

// Count reads the file and counts character frequencies using worker goroutines.
func (c *Counter) Count(filename string) (models.FrequencyMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Channels for work distribution and results
	jobs := make(chan []byte, workerCount)
	results := make(chan models.FrequencyMap, workerCount)

	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go c.worker(jobs, results, &wg)
	}

	// Read the file in chunks and send to workers
	go func() {
		defer close(jobs)
		buffer := make([]byte, chunkSize)
		for {
			n, err := reader.Read(buffer)
			if n > 0 {
				jobs <- buffer[:n] // Avoids unnecessary copying
			}
			if err != nil {
				break
			}
		}
	}()

	// Wait for all workers to finish processing
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and merge results {
	for freqMap := range results {
		c.mergeFrequencies(freqMap)
	}

	return c.frequencies, nil
}

// worker processes chunks of data and counts frequencies.
func (c *Counter) worker(jobs <-chan []byte, results chan<- models.FrequencyMap, wg *sync.WaitGroup) {
	defer wg.Done()

	for chunk := range jobs {
		freqMap := make(models.FrequencyMap)
		for _, b := range chunk {
			freqMap[b]++
		}
		results <- freqMap
	}
}

// mergeFrequencies combines frequency maps in a thread-safe manner.
func (c *Counter) mergeFrequencies(freqMap models.FrequencyMap) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for char, count := range freqMap {
		c.frequencies[char] += count
	}
}
