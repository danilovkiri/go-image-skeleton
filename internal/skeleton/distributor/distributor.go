// Package distributor implements task distribution functionality
package distributor

import (
	"github.com/danilovkiri/go-image-skeleton/internal/skeleton/neighbors"
)

// Distributor implements attributes and methods for task distribution
type Distributor struct {
	ChOut          chan *neighbors.Neighbours
	BinImageMatrix *[][]int
}

// Do iterates over binary image matrix excluding its edges of length 1, created 3x3 neighbor matrices and sends them to queue
func (w *Distributor) Do() {
	// in a matrix representation, the first shape length corresponds to a Y axis in an image
	for y := 1; y < len(*w.BinImageMatrix)-1; y++ {
		// in a matrix representation, the second shape length corresponds to an X axis in an image
		for x := 1; x < len((*w.BinImageMatrix)[0])-1; x++ {
			p1 := neighbors.Point{X: x, Y: y, Color: (*w.BinImageMatrix)[y][x]}
			// skip neighbor matrix when the point of interest (P1) is not colored
			if p1.Color == 0 {
				continue
			}
			// enumerate the neighbors CW
			p2 := neighbors.Point{X: x, Y: y - 1, Color: (*w.BinImageMatrix)[y-1][x]}
			p3 := neighbors.Point{X: x + 1, Y: y - 1, Color: (*w.BinImageMatrix)[y-1][x+1]}
			p4 := neighbors.Point{X: x + 1, Y: y, Color: (*w.BinImageMatrix)[y][x+1]}
			p5 := neighbors.Point{X: x + 1, Y: y + 1, Color: (*w.BinImageMatrix)[y+1][x+1]}
			p6 := neighbors.Point{X: x, Y: y + 1, Color: (*w.BinImageMatrix)[y+1][x]}
			p7 := neighbors.Point{X: x - 1, Y: y + 1, Color: (*w.BinImageMatrix)[y+1][x-1]}
			p8 := neighbors.Point{X: x - 1, Y: y, Color: (*w.BinImageMatrix)[y][x-1]}
			p9 := neighbors.Point{X: x - 1, Y: y - 1, Color: (*w.BinImageMatrix)[y-1][x-1]}
			neighbours := &neighbors.Neighbours{
				P1: p1,
				P2: p2,
				P3: p3,
				P4: p4,
				P5: p5,
				P6: p6,
				P7: p7,
				P8: p8,
				P9: p9,
			}
			// send a neighbor set to the distribution channel
			w.ChOut <- neighbours
		}
	}
}
