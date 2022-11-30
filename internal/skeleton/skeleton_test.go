package skeleton

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"os"
	"testing"
)

func getResource() *image.Gray {
	grayImg := image.NewGray(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: 10, Y: 15},
	})
	matrix := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 1, 0, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 0, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 0, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 0, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 0, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 0, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 0, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 0, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	for i := 0; i < len(matrix); i++ {
		for k := 0; k < len(matrix[0]); k++ {
			var c uint8
			switch matrix[i][k] {
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

func TestNewSkeleton(t *testing.T) {
	grayImg := getResource()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	_ = NewSkeleton(grayImg, &logger)
}

func TestSkeleton_Calculate(t *testing.T) {
	grayImg := getResource()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	skeleton := NewSkeleton(grayImg, &logger)
	skeleton.Calculate(8)
	expected := [][]int{[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, []int{0, 0, 0, 1, 0, 0, 1, 0, 0, 0}, []int{0, 0, 0, 1, 0, 0, 1, 0, 0, 0}, []int{0, 0, 0, 1, 0, 0, 1, 0, 0, 0}, []int{0, 0, 0, 1, 0, 0, 1, 0, 0, 0}, []int{0, 0, 0, 1, 1, 1, 1, 0, 0, 0}, []int{0, 0, 0, 1, 0, 0, 1, 0, 0, 0}, []int{0, 0, 0, 1, 0, 0, 1, 0, 0, 0}, []int{0, 0, 0, 1, 0, 0, 1, 0, 0, 0}, []int{0, 0, 0, 0, 0, 0, 1, 0, 0, 0}, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	assert.Equal(t, expected, skeleton.BinImageMatrix)
}
