package main

import (
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab1/gamma"
	"github.com/anthonynsimon/bild/imgio"
	"image/jpeg"

	"os"
)





func main() {
	curDir, _ := os.Getwd()
	path := curDir+"/dsp/lab1/images/"

	img, err := imgio.Open(path+"1.jpg")
	if err != nil {
		panic(err)
	}

	res := gamma.AddGamma(img, 0.2, 0.22)

	f, _ := os.Create(path + "1_gamma.jpg")

	err = jpeg.Encode(f, res, &jpeg.Options{Quality: 99})
	if err != nil {
		panic(err)
	}

}


