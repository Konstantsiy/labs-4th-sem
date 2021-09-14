package main

import (
	"fmt"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab1/filter"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab1/gamma"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab1/hist"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab1/utils"
	"github.com/anthonynsimon/bild/imgio"
	"log"
	"os"
	"strconv"
)

func prepareVars() (float64, float64, error) {
	args := os.Args

	if len(args) < 3 {
		return 0, 0, fmt.Errorf("need minimum 3 arguments")
	}

	c, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return 0, 0, err
	}

	y, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return 0, 0, err
	}

	return c, y, nil
}

func main() {
	c, y, err := prepareVars()
	if err != nil {
		log.Fatal(err)
	}

	curDir, _ := os.Getwd()
	path := curDir+"/dsp/lab1/images/"
	filename := "1.jpg"


	img, err := imgio.Open(path+filename)
	if err != nil {
		log.Fatal(err)
	}

	err = hist.DrawHistogram(img, path, filename, "source")
	if err != nil {
		log.Fatal(err)
	}

	imgGamma := gamma.AddGamma(img, c, y)
	err = utils.SaveFile(imgGamma, path, filename, "gamma")
	if err != nil {
		log.Fatal(err)
	}

	err = hist.DrawHistogram(imgGamma, path, filename, "gamma")
	if err != nil {
		log.Fatal(err)
	}

	imgSobel := filter.ApplySobel(img)
	err = utils.SaveFile(imgSobel, path, filename, "sobel")
	if err != nil {
		log.Fatal(err)
	}

	err = hist.DrawHistogram(imgSobel, path, filename, "sobel")
	if err != nil {
		log.Fatal(err)
	}
}


