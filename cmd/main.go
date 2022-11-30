package main

import (
	"flag"
	"github.com/danilovkiri/go-image-skeleton/internal/binarize"
	"github.com/danilovkiri/go-image-skeleton/internal/logger"
	"github.com/danilovkiri/go-image-skeleton/internal/skeleton"
	"github.com/rs/zerolog"
	"image/png"
	"log"
	"os"
)

func main() {
	// parse CLI
	debug := flag.Bool("debug", false, "sets log level to debug")
	nWorkers := flag.Int("workers", 1, "sets number of workers for parallel computing")
	fileIn := flag.String("input", "defaultInput.png", "sets input image path")
	fileOut := flag.String("output", "defaultOutput.png", "sets output image path")
	flag.Parse()
	if *fileIn == "defaultInput.png" {
		log.Panic("--input is required")
	}

	//initialize logger
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	loggerInstance := logger.InitLog()

	bn, err := binarize.NewBinarizer(*fileIn)
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("binarizer was not instantiated")
	}
	err = bn.RgbaToGray()
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("image was not converted to grayscale")
	}
	err = bn.GetIntegralImage()
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("integral image matrix was not computed")
	}
	err = bn.BradleyBinarize()
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("image binarization was not performed")
	}

	skeletonInstance := skeleton.NewSkeleton(bn.BinImage, loggerInstance)
	skeletonInstance.Calculate(*nWorkers)

	f, _ := os.Create(*fileOut)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			loggerInstance.Fatal().Err(err).Msg("result image file closure failed")
		}
	}(f)
	//err = png.Encode(f, bn.BinImage)
	err = png.Encode(f, skeletonInstance.GetImage())
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("result image encoding failed")
	}
}
