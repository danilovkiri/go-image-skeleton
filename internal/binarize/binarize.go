// Package binarize provides image binarization functionality
package binarize

import (
	"errors"
	"github.com/danilovkiri/go-image-skeleton/internal/integral"
	"image"
	"image/color"
	_ "image/png"
	"os"
)

// Binarizer implements attributes and methods for image binarization
type Binarizer struct {
	bounds   image.Rectangle
	img      *image.Image
	grayImg  *image.Gray
	intImage *integral.IntegralImage
	BinImage *image.Gray
}

// NewBinarizer initializes a new Binarizer instance
func NewBinarizer(path string) (*Binarizer, error) {
	infile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer infile.Close()
	img, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}
	binarizer := &Binarizer{
		bounds:   img.Bounds(),
		img:      &img,
		grayImg:  nil,
		intImage: nil,
		BinImage: nil,
	}
	return binarizer, nil
}

// RgbaToGray converts an image.Image to an image.Gray
func (bn *Binarizer) RgbaToGray() error {
	if bn.img == nil {
		return errors.New("original image is not initialized")
	}
	grayImg := image.NewGray(bn.bounds)
	for x := 0; x < bn.bounds.Max.X; x++ {
		for y := 0; y < bn.bounds.Max.Y; y++ {
			var rgba = (*bn.img).At(x, y)
			grayImg.Set(x, y, rgba)
		}
	}
	bn.grayImg = grayImg
	return nil
}

// BradleyBinarize performs Bradley binarization of an image.Gray
func (bn *Binarizer) BradleyBinarize() error {
	if bn.intImage == nil {
		return errors.New("integral image is not initialized")
	}
	binImg := image.NewGray(bn.bounds)
	s := bn.bounds.Max.X / 8
	s2 := s / 2
	t := 0.15
	for x := 0; x < bn.bounds.Max.X; x++ {
		for y := 0; y < bn.bounds.Max.Y; y++ {
			x1 := x - s2
			x2 := x + s2
			y1 := y - s2
			y2 := y + s2
			if x1 < 0 {
				x1 = 0
			}
			if x2 >= bn.bounds.Max.X {
				x2 = bn.bounds.Max.X - 1
			}
			if y1 < 0 {
				y1 = 0
			}
			if y2 >= bn.bounds.Max.Y {
				y2 = bn.bounds.Max.Y - 1
			}
			intensityAvrIntegralSubImage := bn.intImage.GetSum(x1, y1, x2, y2) / float64((x2-x1)*(y2-y1))
			// we invert here since intensity of 255 represents white and of 0 — black
			if float64(bn.grayImg.GrayAt(x, y).Y) < intensityAvrIntegralSubImage*(1-t) {
				binImg.SetGray(x, y, color.Gray{Y: 0})
			} else {
				binImg.SetGray(x, y, color.Gray{Y: 255})
			}
		}
	}
	bn.BinImage = binImg
	return nil
}

// GetIntegralImage creates an integral image matrix
func (bn *Binarizer) GetIntegralImage() error {
	if bn.grayImg == nil {
		return errors.New("gray image is not initialized")
	}
	cx := (*bn.img).Bounds().Max.X
	cy := (*bn.img).Bounds().Max.Y
	intImage := integral.NewIntegralImage(cx, cy, bn.grayImg)
	intImage.Calculate()
	bn.intImage = intImage
	return nil
}
