package main

import (
	"fmt"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab1/util"
	"github.com/anthonynsimon/bild/imgio"
	"image"
	"image/draw"
	"log"
	"os"
	"strconv"
)

const (
	ColorBlack = 0x00
	ColorWhite = 0xFF
	ColorYellow = 0xE0
)

var (
	colorBlack = [3]uint8{0, 0, 0}
	colorWhite = [3]uint8{255, 255, 255}
	colorsRGB = [7][3]uint8{
		{255, 0, 0},
		{255, 69, 0},
		{0, 255, 0},
		{0, 191, 255},
		{0, 0, 205},
		{139, 0, 139},
		{255, 105, 180},
	}
)

func prepareVars() (string, uint8, error) {
	args := os.Args

	if len(args) < 3 {
		return "", 0, fmt.Errorf("need minimum 3 arguments")
	}

	filename := args[1]

	level, err := strconv.ParseUint(args[2], 10, 8)
	if err != nil {
		return "", 0, err
	}

	return filename, uint8(level), nil
}

func AsRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	res := image.NewRGBA(bounds)
	draw.Draw(res, bounds, src, bounds.Min, draw.Src)
	return res
}

func fillRGBPixel(pix []uint8, pos int, colorRGB [3]uint8) {
	for i := 0; i < 3; i++ {
		 pix[pos+i] = colorRGB[i]
	}
}

func BinarizeImage(img image.Image, level uint8) *image.Gray {
	//result := image.NewGray(img.Bounds())
	//for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
	//	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
	//		if uint8(util.Rank(img.At(x, y))) <= level {
	//			result.Set(x, y, color.Black)
	//		} else {
	//			result.Set(x, y, color.White)
	//		}
	//	}
	//}

	//result := image.NewGray(img.Bounds())
	//for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
	//	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
	//		result.Set(x, y, img.At(x, y))
	//	}
	//}

	src := AsRGBA(img)
	bounds := img.Bounds()
	result := image.NewGray(bounds)

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			srcPos := y * src.Stride + x * 4
			resPos := y * result.Stride + x

			c := src.Pix[srcPos : srcPos+4]

			r := float64(c[0])*0.3 + float64(c[1])*0.6 + float64(c[2])*0.1

			if uint8(r) >= level {
				result.Pix[resPos] = ColorWhite
			} else {
				result.Pix[resPos] = ColorBlack
			}
		}
	}

	return result
}

func MarkupImage(img image.Gray) image.Gray {
	var A, B, C uint8
	curColorPos := 0
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			curPos := y*img.Stride+x
			kn := x - 1
			if kn <= 0 {
				kn = 1
				B = ColorBlack
			} else {
				B = img.Pix[kn * img.Stride + x]
			}
			km := y - 1
			if km <= 0 {
				km = 1
				C = ColorBlack
			} else {
				C = img.Pix[km*img.Stride+x]
			}
			A = img.Pix[curPos]
			if A == 0 {
				continue
			} else if B == ColorBlack && C == ColorBlack {
				fillRGBPixel(img.Pix, curPos, colorsRGB[curColorPos])
				//for i := 0; i < 3; i++ {
				//	img.Pix[curPos+i] = colorsRGB[curColorPos][i]
				//}
				curColorPos++
			} else if B != ColorBlack && C == ColorBlack {
				img.Pix[curPos] = B
			} else if B == ColorBlack && C != ColorBlack {
				img.Pix[curPos] = C
			} else if B != ColorBlack && C != ColorBlack {
				if B == C {
					img.Pix[curPos] = B
				} else {
					img.Pix[curPos] = B

				}
			}
		}
	}

	return img
}

func main() {
	filename, level, err := prepareVars()
	if err != nil {
		log.Fatal(err)
	}

	curDir, _ := os.Getwd()
	path := curDir+"/dsp/lab2/images/"
	filename += ".jpg"

	img, err := imgio.Open(path+filename)
	if err != nil {
		log.Fatal(err)
	}

	binImg := BinarizeImage(img, level)
	err = util.SaveFile(binImg, path, filename, "bin")
	if err != nil {
		log.Fatal(err)
	}

	//markImg := MarkupImage(*binImg)
	//err = util.SaveFile(&markImg, path, filename, "markup")
	//if err != nil {
	//	log.Fatal(err)
	//}


	//err = hist.DrawHistogram(img, path, filename, "source")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//imgGamma := gamma.AddGamma(img, c, y)
	//err = utils.SaveFile(imgGamma, path, filename, "gamma")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = hist.DrawHistogram(imgGamma, path, filename, "gamma")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//imgSobel := filter.ApplySobel(img)
	//err = utils.SaveFile(imgSobel, path, filename, "sobel")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = hist.DrawHistogram(imgSobel, path, filename, "sobel")
	//if err != nil {
	//	log.Fatal(err)
	//}
}


