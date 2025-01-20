package frequency

import (
	"huffman-compression/internal/models"
	"io"
)

// Counter is a struct that holds a map of character frequencies.
type Counter struct {
	frequencies models.FrequencyMap
}

// NewCounter initializes and returns a new Counter instance.
func NewCounter() *Counter {
	return &Counter{
		frequencies: make(models.FrequencyMap),
	}
}

// Count reads from the provided io.Reader and updates the character frequencies.
func (c *Counter) Count(reader io.Reader) error {
	const bufferSize = 32 * 1024 // 32KB buffer efficient reading
	buf := make([]byte, bufferSize)

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		for i := 0; i < n; i++ {
			c.frequencies[buf[i]]++
		}
	}

	return nil
}

// GetFrequencies returns the map of character frequencies.
func (c *Counter) GetFrequencies() models.FrequencyMap {
	return c.frequencies
}
