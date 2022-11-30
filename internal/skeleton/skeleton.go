// Package skeleton provides pattern thinning functionality known via Zhang-Suen algorithm
package skeleton

import (
	"context"
	"fmt"
	"github.com/danilovkiri/go-image-skeleton/internal/skeleton/collector"
	"github.com/danilovkiri/go-image-skeleton/internal/skeleton/distributor"
	"github.com/danilovkiri/go-image-skeleton/internal/skeleton/neighbors"
	"github.com/danilovkiri/go-image-skeleton/internal/skeleton/processor"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"image"
	"image/color"
	"log"
	"sync"
)

// Skeleton implements attributes and methods for task processing
type Skeleton struct {
	binImage       *image.Gray
	BinImageMatrix [][]int
	logger         *zerolog.Logger
}

// GetImage retrieves a thinned image.Gray ready to be encoded
func (sk *Skeleton) GetImage() *image.Gray {
	grayImg := image.NewGray(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: len(sk.BinImageMatrix[0]), Y: len(sk.BinImageMatrix)},
	})
	for i := 0; i < len(sk.BinImageMatrix); i++ {
		for k := 0; k < len(sk.BinImageMatrix[0]); k++ {
			var c uint8
			switch sk.BinImageMatrix[i][k] {
			case 0:
				c = 255
			case 1:
				c = 0
			}
			grayImg.SetGray(k, i, color.Gray{Y: c})
		}
	}
	return grayImg
}

// NewSkeleton initializes a new Skeleton instance and created a binary image matrix
func NewSkeleton(binImage *image.Gray, logger *zerolog.Logger) *Skeleton {
	logger.Warn().Msg("Initializing skeleton instance")
	nRows := binImage.Bounds().Max.Y
	nColumns := binImage.Bounds().Max.X
	binImageMatrix := make([][]int, nRows)
	for i := 0; i < nRows; i++ {
		for j := 0; j < nColumns; j++ {
			c := int(binImage.GrayAt(j, i).Y)
			if c == 0 {
				c = 1
			} else {
				c = 0
			}
			binImageMatrix[i] = append(binImageMatrix[i], c)
		}
	}
	logger.Info().Msg(fmt.Sprintf("Binary image matrix created with shape %d x %d", len(binImageMatrix[0]), len(binImageMatrix)))
	return &Skeleton{
		binImage:       binImage,
		BinImageMatrix: binImageMatrix,
		logger:         logger,
	}
}

// Calculate performs parallel computation of each thinning iteration
func (sk *Skeleton) Calculate(nWorkers int) {
	sk.logger.Warn().Msg(fmt.Sprintf("Initializing skeleton calculation with %d workers", nWorkers))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for counter := 1; ; counter++ {
		sk.logger.Info().Msg(fmt.Sprintf("Performing subiteration %d ", counter))
		wg := sync.WaitGroup{}
		chanNeighbors := make(chan *neighbors.Neighbours)
		chanResult := make(chan *neighbors.Point)

		// put tasks into the distributor channel in the background
		dist := &distributor.Distributor{
			ChOut:          chanNeighbors,
			BinImageMatrix: &sk.BinImageMatrix,
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			sk.logger.Warn().Msg("Starting sending tasks")
			dist.Do()
			// we can close chanNeighbors here since the distributor is done
			sk.logger.Warn().Msg("Closing distributor channel")
			close(chanNeighbors)
		}()

		// collect processed tasks in the background
		diffMatrix := make([][]int, len(sk.BinImageMatrix))
		for i := 0; i < len(sk.BinImageMatrix); i++ {
			for k := 0; k < len(sk.BinImageMatrix[0]); k++ {
				diffMatrix[i] = append(diffMatrix[i], 0)
			}
		}
		coll := &collector.Collector{
			ChIn:       chanResult,
			DiffMatrix: diffMatrix,
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			sk.logger.Warn().Msg("Starting collecting tasks")
			coll.Do()
		}()

		// start workers to process tasks from distributor channel in the background
		sk.logger.Warn().Msg("Starting processing tasks")
		g, _ := errgroup.WithContext(ctx)
		for i := 0; i < nWorkers; i++ {
			w := &processor.Processor{ChIn: chanNeighbors, ChOut: chanResult, FirstIter: counter%2 == 1}
			g.Go(w.Do)
		}

		// wait for workers completion after all tasks processed and chanNeighbors is closed
		err := g.Wait()
		if err != nil {
			sk.logger.Fatal().Err(err).Msg("processor group exited with error")
			log.Fatal(err)
		}

		// close chanResult after all workers processed their tasks
		sk.logger.Warn().Msg("Closing collector channel")
		close(chanResult)

		// wait for distributor and collector exit since chanResult was closed
		wg.Wait()

		// check that diffMatrix contains any elements for removal
		nonZero := false
		for i := 0; i < len(coll.DiffMatrix); i++ {
			if !nonZero {
				for k := 0; k < len(coll.DiffMatrix[0]); k++ {
					if coll.DiffMatrix[i][k] == 1 {
						nonZero = true
						break
					}
				}
			} else {
				break
			}
		}
		if !nonZero {
			sk.logger.Warn().Msg(fmt.Sprintf("Zero diff was obtained at subiteration %d ", counter))
			break
		}

		// substitute the diffMatrix from the binImageMatrix
		for i := 0; i < len(sk.BinImageMatrix); i++ {
			for k := 0; k < len(sk.BinImageMatrix[0]); k++ {
				sk.BinImageMatrix[i][k] = sk.BinImageMatrix[i][k] - coll.DiffMatrix[i][k]
			}
		}
		if sk.logger.Debug().Enabled() {
			for i := 0; i < len(sk.BinImageMatrix); i++ {
				var entry []string
				for k := 0; k < len(sk.BinImageMatrix[i]); k++ {
					switch sk.BinImageMatrix[i][k] {
					case 1:
						entry = append(entry, "@")
					case 0:
						entry = append(entry, ".")
					}
				}
				sk.logger.Debug().Msg(fmt.Sprint(entry))
			}
		}
	}
}
