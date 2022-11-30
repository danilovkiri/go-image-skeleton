// Package collector implements processed tasks collecting functionality
package collector

import "github.com/danilovkiri/go-image-skeleton/internal/skeleton/neighbors"

// Collector implements attributes and methods for task collecting
type Collector struct {
	ChIn       chan *neighbors.Point
	DiffMatrix [][]int
}

// Do creates a diff matrix filling it with points for further removal from binary image matrix
func (c *Collector) Do() {
	for i := range c.ChIn {
		c.DiffMatrix[i.Y][i.X] = i.Color
	}
}
