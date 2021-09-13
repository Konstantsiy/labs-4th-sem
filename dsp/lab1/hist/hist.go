package hist

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

func getImageSize(r *os.File) {
	im, _, err := image.DecodeConfig(r)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("filename: %s\twidht: %d\theight: %d\n", r.Name(), im.Width, im.Height)
}

func CalcHistogramComponents(path, filename string) error {
	file, err := os.Open(path+filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//getImageSize(file)

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

func DrawHistogram(path, filename string) error {
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
