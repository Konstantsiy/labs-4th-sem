package filter

import (
	"image"
	"image/color"
	"math"
)

var (
	kernelX = [3][3]int8{
		{1, 0, -1},
		{1, 0, -1},
		{1, 0, -1},
	}

	kernelY = [3][3]int8{
		{-1, -1, -1},
		{0, 0, 0},
		{1, 1, 1},
	}
)

func ToGrayscale(img image.Image) image.Image {
	max := img.Bounds().Max
	min := img.Bounds().Min

	var filtered = image.NewGray(image.Rect(max.X, max.Y, min.X, min.Y))

	for x := 0; x < max.X; x++ {
		for y := 0; y < max.Y; y++ {
			//grayColor := color.GrayModel.Convert(img.At(x, y))
			filtered.Set(x, y, img.At(x, y))
		}
	}

	return filtered
}

func getGrayPixel(c color.Color) uint8 {
	p, _, _, _ := c.RGBA()
	ret := uint8(p)
	return ret
}

func ApplySobel(img image.Image) image.Image {
	img = ToGrayscale(img)
	bounds := img.Bounds()

	var pixel color.Color
	var filtered = image.NewGray(image.Rect(bounds.Max.X-2, bounds.Max.Y-2, bounds.Min.X, bounds.Min.Y))

	for x := 1; x < bounds.Max.X - 1; x++ {
		for y := 1; y < bounds.Max.Y - 1; y++ {

			//-----------------------------
			var fX, fY int
			curX, curY := x - 1, y - 1
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					pixel := getGrayPixel(img.At(curX, curY))
					fX += int(kernelX[i][j]) * int(pixel)
					fY += int(kernelY[i][j]) * int(pixel)
					if x > 1 {
						curX = curX + j - 1
					}
				}
				curY++
			}
			//-----------------------------

			//fX, fY := applyKernels(img, x, y)

			v := uint32(math.Ceil(float64(getMax(fX, fY))))
			//v := uint32(math.Ceil(math.Sqrt(float64((fX * fX) + (fY * fY)))))
			pixel = color.Gray{Y: uint8(v)}
			filtered.SetGray(x, y, pixel.(color.Gray))
		}
	}

	return filtered
}

func getMax(f1, f2 int) int {
	if f1 > f2 {
		return f1
	}
	return f2
}

func applyKernels(img image.Image, x, y int) (uint32, uint32) {
	var fX, fY int

	curX := x - 1
	curY := y - 1

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			pixel := getGrayPixel(img.At(curX, curY))
			fX += int(kernelX[i][j]) * int(pixel)
			fY += int(kernelY[i][j]) * int(pixel)
			if curX > 0 {
				curX = curX + j - 1
			}
		}
		curY++
	}

	uFX, uFY := uint32(math.Abs(float64(fX))), uint32(math.Abs(float64(fY)))
	return uFX, uFY
}


func ApplySobel1(img image.Image) image.Image {
	greyImg := ToGrayscale(img)
	bounds := greyImg.Bounds()
	var filtered = image.NewGray(image.Rect(bounds.Max.X-2, bounds.Max.Y-2, bounds.Min.X, bounds.Min.Y))

	for x := 1; x < bounds.Max.X-1; x++ {
		for y := 1; y < bounds.Max.Y-1; y++ {
			var fX, fY int
			var pixel color.Color
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					xn := x+i-1
					yn := y+j-1
					fX += int(kernelX[i][j]) * int(getGrayPixel(img.At(xn, yn)))
					fY += int(kernelY[i][j]) * int(getGrayPixel(img.At(xn, yn)))
				}
			}
			ufX, ufY := uint32(math.Abs(float64(fX))), uint32(math.Abs(float64(fY)))
			v := uint32(math.Ceil(math.Sqrt(float64((ufX * ufX) + (ufY * ufY)))))
			pixel = color.Gray{Y: uint8(v)}
			filtered.SetGray(x, y, pixel.(color.Gray))
		}
	}

	return filtered
}
