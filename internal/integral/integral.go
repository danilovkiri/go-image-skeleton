// Package integral provides integral image matrix calculation functionality
package integral

import "image"

// IntegralImage defines attributes and methods for integral image matrix calculation
type IntegralImage struct {
	Pix     []float64
	Width   int
	Height  int
	GrayImg *image.Gray
}

// NewIntegralImage initializes a new IntegralImage instance
func NewIntegralImage(x, y int, grayImage *image.Gray) *IntegralImage {
	return &IntegralImage{
		Pix:     make([]float64, x*y),
		Width:   x,
		Height:  y,
		GrayImg: grayImage,
	}
}

// Calculate calculates the integral image matrix
func (i *IntegralImage) Calculate() {
	offset := 0
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			a := float64(i.GrayImg.Pix[offset])
			b := i.get(x-1, y)
			c := i.get(x, y-1)
			d := i.get(x-1, y-1)
			i.Pix[offset] = a + b + c - d
			offset++
		}
	}
}

// get retrieves integral image matrix value by its indices
func (i *IntegralImage) get(x, y int) float64 {
	if x < 0 || y < 0 {
		return 0
	}
	return float64(i.Pix[(y*i.Width)+x])
}

// GetSum calculates a sum of intensities in the original image inside a specified area by its indices
func (i *IntegralImage) GetSum(x1, y1, x2, y2 int) float64 {
	a := i.get(x1-1, y1-1)
	b := i.get(x2, y1-1)
	c := i.get(x1-1, y2)
	d := i.get(x2, y2)
	return a + d - b - c
}
