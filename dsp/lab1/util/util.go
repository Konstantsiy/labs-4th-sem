package util

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

func SaveJPG(img image.Image, path, filename, postfix string) error {
	filename += ".jpg"
	f, err := os.Create(path+ filename[0:len(filename)-4]+"_"+postfix+filename[len(filename)-4:])
	if err != nil {
		return err
	}

	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 99})
	if err != nil {
		return err
	}
	return nil
}

func SavePNG(img image.Image, path, filename, postfix string) error {
	filename += ".png"
	f, err := os.Create(path+ filename[0:len(filename)-4]+"_"+postfix+filename[len(filename)-4:])
	if err != nil {
		return err
	}

	var enc png.Encoder
	enc.CompressionLevel = 90
	err = enc.Encode(f, img)
	if err != nil {
		return err
	}
	return nil
}

func Rank(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return float64(r)*0.3 + float64(g)*0.6 + float64(b)*0.1
}