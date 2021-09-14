package utils

import (
	"image"
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