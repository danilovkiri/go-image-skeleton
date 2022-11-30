// Package processor implements neighbor set processing functionality
package processor

import "github.com/danilovkiri/go-image-skeleton/internal/skeleton/neighbors"

// Processor implements attributes and methods for task processing
type Processor struct {
	ChIn      chan *neighbors.Neighbours
	ChOut     chan *neighbors.Point
	FirstIter bool
}

// Do analyzes the neighbor set and sends the resulting point with its removal status to the collector channel
func (p *Processor) Do() error {
	// iterate over the tasks in the distributor channel until it is closed
	for neighborSet := range p.ChIn {
		nNonZeros := neighborSet.GetNonZeros()                                           // condition a
		nTransitions := neighborSet.GetTransitions()                                     // condition b
		product246 := neighborSet.P2.Color * neighborSet.P4.Color * neighborSet.P6.Color // condition c
		product468 := neighborSet.P4.Color * neighborSet.P6.Color * neighborSet.P8.Color // condition d
		product248 := neighborSet.P2.Color * neighborSet.P4.Color * neighborSet.P8.Color // condition c'
		product268 := neighborSet.P2.Color * neighborSet.P6.Color * neighborSet.P8.Color // condition d'
		var diff int
		switch p.FirstIter {
		case true:
			// if it is the first sub-iteration, the use conditions a, b, c and d
			if nNonZeros >= 2 && nNonZeros <= 6 && nTransitions == 1 && product246 == 0 && product468 == 0 {
				diff = 1
			} else {
				diff = 0
			}
		case false:
			// if it is the secod sub-iteration, the use conditions a, b, c' and d'
			if nNonZeros >= 2 && nNonZeros <= 6 && nTransitions == 1 && product248 == 0 && product268 == 0 {
				diff = 1
			} else {
				diff = 0
			}
		}
		p.ChOut <- &neighbors.Point{
			X:     neighborSet.P1.X,
			Y:     neighborSet.P1.Y,
			Color: diff,
		}
	}
	return nil
}
