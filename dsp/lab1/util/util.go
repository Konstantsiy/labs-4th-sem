package util

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func SaveFile(img image.Image, path, filename, postfix string) error {
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

func Rank(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return float64(r)*0.3 + float64(g)*0.6 + float64(b)*0.1
}