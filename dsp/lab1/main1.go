package main

import (
	"fmt"
	"github.com/anthonynsimon/bild/histogram"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"log"
	"os"
)

func calcHistogramComponents(path, filename string) error {
	file, err := os.Open(path+filename)
	if err != nil {
		return err
	}

	m, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := m.Bounds()
	var hist [16][4]int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			hist[r>>12][0]++
			hist[g>>12][1]++
			hist[b>>12][2]++
			hist[a>>12][3]++
		}
	}

	fmt.Printf("%-14s %6s %6s %6s %6s\n", "bin", "red", "green", "blue", "alpha")
	for i, x := range hist {
		fmt.Printf("0x%04x-0x%04x: %6d %6d %6d %6d\n", i<<12, (i+1)<<12-1, x[0], x[1], x[2], x[3])
	}

	return nil
}

func drawHistogram(path, filename string) error {
	img, err := imgio.Open(path+filename)
	if err != nil {
		return err
	}

	hist := histogram.NewRGBAHistogram(img)
	result := hist.Image()

	version := uuid.New().String()[0:4]

	f, err := os.Create(path+ filename[0:len(filename)-4]+"_"+version+filename[len(filename)-4:])
	if err != nil {
		return err
	}
	defer f.Close()

	err = jpeg.Encode(f, result, &jpeg.Options{Quality: 99})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	curDir, _ := os.Getwd()
	path := curDir+"/dsp/lab1/images/"

	go func(){
		log.Print("start calc")
		err := drawHistogram(path, "1.jpg")
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := calcHistogramComponents(path, "1.jpg")
	if err != nil {
		log.Fatal(err)
	}
}


