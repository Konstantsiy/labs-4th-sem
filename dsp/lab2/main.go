package main

import (
	"fmt"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab1/util"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/segment"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	ColorBlack = 0x00
	ColorWhite = 0xFF
	ColorYellow = 0xE0
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

func BinarizeImage1(img image.Image, level uint8) (*image.Gray, [][]byte) {
	src := AsRGBA(img)
	bounds := img.Bounds()
	result := image.NewGray(bounds)

	binMap := make([][]byte, bounds.Dy())

	for y := 0; y < bounds.Dy(); y++ {
		binMap[y] = make([]byte, bounds.Dx())
		for x := 0; x < bounds.Dx(); x++ {
			srcPos := y * src.Stride + x * 4
			resPos := y * result.Stride + x

			c := src.Pix[srcPos : srcPos+4]

			r := float64(c[0])*0.3 + float64(c[1])*0.6 + float64(c[2])*0.1

			if uint8(r) >= level {
				result.Pix[resPos] = ColorWhite
				binMap[y][x] = 1
			} else {
				result.Pix[resPos] = ColorBlack
				binMap[y][x] = 0
			}
		}
	}

	return result, binMap
}

type Coordinate struct {
	H int
	W int
	C bool
}

type Coordinates []Coordinate

func FindObjects(binMap [][]byte) (map[byte]Coordinates, [][]byte) {
	height, width := len(binMap), len(binMap[0])
	objects := make(map[byte]Coordinates)
	var cur byte
	var A, B, C byte
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			kn := j - 1
			if kn <= 0 {
				kn = 1
				B = 0
			} else {
				B = binMap[i][kn]
			}
			km := i - 1
			if km <= 0 {
				km = 1
				C = 0
			} else {
				C = binMap[km][j]
			}
			A = binMap[i][j]
			if A == 0 {
			} else if B == 0 && C == 0 {
				if len(objects) == 0 {
					cur = A
				} else {
					var m byte
					for k, _ := range objects {
						m = k
					}
					cur = m + 1
				}
				binMap[i][j] = cur
				if _, ok := objects[cur]; !ok {
					objects[cur] = Coordinates{}
				}
				objects[cur] = append(objects[cur], Coordinate{H: i, W: j})
			} else if B != 0 && C == 0 {
				binMap[i][j] = B
				if _, ok := objects[B]; !ok {
					objects[B] = Coordinates{}
				}
				objects[B] = append(objects[B], Coordinate{H: i, W: j})
			} else if B == 0 && C != 0 {
				binMap[i][j] = C
				if _, ok := objects[C]; !ok {
					objects[C] = Coordinates{}
				}
				objects[C] = append(objects[C], Coordinate{H: i, W: j})
			} else if B != 0 && C != 0 {
				binMap[i][j] = B
				if _, ok := objects[B]; !ok {
					objects[B] = Coordinates{}
				}
				objects[B] = append(objects[B], Coordinate{H: i, W: j})
				if B != C {
					if _, ok := objects[C]; ok {
						for _, cor := range objects[C] {
							binMap[cor.H][cor.W] = B
						}
						for _, cor := range objects[C] {
							objects[B] = append(objects[B], cor)
						}
						delete(objects, C)
					}
				}
			}
		}
	}
	return objects, binMap
}

func moment(i, j int, wMean, hMean float64, coordinates Coordinates) float64 {
	var result float64
	for _, c := range coordinates {
		result += math.Pow(float64(c.W)-wMean, float64(i)) * math.Pow(float64(c.H)-hMean, float64(j))
	}
	return result
}

func CalcCharacteristics(bm [][]byte, coordinates Coordinates) (int, int, float64, float64, float64) {
	square := len(coordinates)
	perimeter := CalcPerim(bm, coordinates)
	compact := math.Pow(float64(perimeter), 2) / float64(square)

	sumW, sumH := 0, 0
	for _, c := range coordinates {
		sumH += c.H
		sumW += c.W
	}
	hMean := float64(sumH) / float64(square)
	wMean := float64(sumW) / float64(square)

	m02 := moment(0, 2, wMean, hMean, coordinates)
	m20 := moment(2, 0, wMean, hMean, coordinates)
	m11 := moment(1, 1, wMean, hMean, coordinates)

	nominator := m20 + m02 + math.Sqrt(math.Pow(m20-m02, 2)+4*math.Pow(m11, 2))
	denominator := m20 + m02 - math.Sqrt(math.Pow(m20-m02, 2)+4*math.Pow(m11, 2))
	elongation := nominator / denominator

	orientation := 0.5 * math.Atan((2 * m11) / (m20 - m02))

	return square, perimeter, compact, elongation, orientation
}

func bc(c []uint8) bool {
	return !(c[0] == 255 || c[1] == 255 || c[2] == 255)
}

func getPixel(img image.RGBA, x, y int) []uint8 {
	pos := y * img.Stride + x * 4
	return img.Pix[pos:pos+4]
}

func isBinaryNoise(src *image.RGBA, x, y int, width, height int) bool {
	if x == 0 {
		if y == 0 {
			return bc(getPixel(*src, x, y)) && bc(getPixel(*src, x+1, y)) && bc(getPixel(*src, x, y+1))
		} else if y >= height-1 {
			return bc(getPixel(*src, x, y)) && bc(getPixel(*src, x+1, y)) && bc(getPixel(*src, x, y-1))
		} else {
			return bc(getPixel(*src, x, y)) && bc(getPixel(*src, x+1, y)) && bc(getPixel(*src, x, y+1)) && bc(getPixel(*src, x, y-1))
		}
	} else if x >= width-1 {
		if y == 0 {
			return bc(getPixel(*src, x, y)) && bc(getPixel(*src, x-1, y)) && bc(getPixel(*src, x, y+1))
		} else if y >= height-1 {
			return bc(getPixel(*src, x, y)) && bc(getPixel(*src, x-1, y)) && bc(getPixel(*src, x, y-1))
		} else {
			return bc(getPixel(*src, x, y)) && bc(getPixel(*src, x-1, y)) && bc(getPixel(*src, x, y+1)) && bc(getPixel(*src, x, y-1))
		}
	} else {
		if y == 0 {
			return bc(getPixel(*src, x, y)) && bc(getPixel(*src, x-1, y)) && bc(getPixel(*src, x+1, y)) && bc(getPixel(*src, x, y+1))
		} else if y >= height-1{
			return bc(getPixel(*src, x, y)) && bc(getPixel(*src, x-1, y)) && bc(getPixel(*src, x+1, y)) && bc(getPixel(*src, x, y-1))
		} else {
			return bc(getPixel(*src, x, y)) && bc(getPixel(*src, x-1, y)) && bc(getPixel(*src, x+1, y)) && bc(getPixel(*src, x, y+1)) && bc(getPixel(*src, x, y-1))
		}
	}
}

func ClearBinaryNoise(img *image.RGBA) *image.RGBA {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			pos := y * img.Stride + x * 4
			if !isBinaryNoise(img, x, y, width, height) {
				img.Pix[pos+0], img.Pix[pos+1], img.Pix[pos+2] = 255, 255, 255
			} else {
				img.Pix[pos+0], img.Pix[pos+1], img.Pix[pos+2] = 0, 0, 0
			}
		}
	}
	return img
}

func BinarizeImageWithLevel(img image.Image, level uint8) *image.RGBA {
	src := AsRGBA(img)
	for y := 0; y < src.Bounds().Dy(); y++ {
		for x := 0; x < src.Bounds().Dx(); x++ {
			pos := y * src.Stride + x * 4
			c := src.Pix[pos : pos+4]
			r := float64(c[0])*0.3 + float64(c[1])*0.6 + float64(c[2])*0.1
			if uint8(r) >= level {
				src.Pix[pos+0], src.Pix[pos+1], src.Pix[pos+2] = 255, 255, 255
			} else {
				src.Pix[pos+0], src.Pix[pos+1], src.Pix[pos+2] = 0, 0, 0
			}
		}
	}
	return src
}

func Binarization(img image.Image, level uint8) (*image.Gray, [][]byte) {
	src := clone.AsRGBA(img)
	bounds := src.Bounds()

	dst := image.NewGray(bounds)

	binMap := make([][]byte, bounds.Dy())

	for y := 0; y < bounds.Dy(); y++ {
		binMap[y] = make([]byte, bounds.Dx())
		for x := 0; x < bounds.Dx(); x++ {
			srcPos := y*src.Stride + x*4
			dstPos := y*dst.Stride + x

			c := src.Pix[srcPos : srcPos+4]
			r := util.Rank(color.RGBA{c[0], c[1], c[2], c[3]})

			// transparent pixel is always white
			if c[0] == 0 && c[1] == 0 && c[2] == 0 && c[3] == 0 {
				dst.Pix[dstPos] = 0xFF
				binMap[y][x] = 0
				continue
			}

			if uint8(r) >= level {
				dst.Pix[dstPos] = 0xFF
				binMap[y][x] = 1
			} else {
				dst.Pix[dstPos] = 0x00
				binMap[y][x] = 0
			}
		}
	}

	return dst, binMap
}

func GetBinMap(img image.Gray) [][]byte {
	bounds := img.Bounds()
	binMap := make([][]byte, bounds.Dy())

	for y := 0; y < bounds.Dy(); y++ {
		binMap[y] = make([]byte, bounds.Dx())
		for x := 0; x < bounds.Dx(); x++ {
			pos := y * img.Stride + x

			if img.Pix[pos] == 0xFF {
				binMap[y][x] = 1
			} else if img.Pix[pos] == 0x00 {
				binMap[y][x] = 0
			}
		}
	}

	return binMap
}

func fill(bm [][]byte, x, y int, c byte, objects map[byte]Coordinates) {
	if bm[x][y] == 1 {
		bm[x][y] = c
		if _, ok := objects[c]; !ok {
			objects[c] = Coordinates{}
		}
		objects[c] = append(objects[c], Coordinate{H: x, W: y})
		if x > 0 {
			fill(bm, x - 1, y, c, objects)
		}
		if x < len(bm) - 1 {
			fill(bm, x + 1, y, c, objects)
		}
		if y > 0 {
			fill(bm, x, y - 1, c, objects)
		}
		if y < len(bm[0]) - 1 {
			fill(bm, x, y + 1, c, objects)
		}
	}
}

func FindObjectsRec(bm [][]byte) (map[byte]Coordinates, [][]byte) {
	objects := make(map[byte]Coordinates)
	var c byte = 1
	for i := 0; i < len(bm); i++ {
		for j := 0; j < len(bm[0]); j++ {
			c++
			fill(bm, i, j, c, objects)
		}
	}
	return objects, bm
}

func isBoundary(bm [][]byte, h, w int) bool {
	if h == 0 || h == len(bm)-1 || w == 0 || w == len(bm[0])-1 {
		return true
	}
	return bm[h+1][w] == 0 || bm[h-1][w] == 0 || bm[h][w+1] == 0 || bm[h][w-1] == 0
}

func CalcPerim(bm [][]byte, coordinates Coordinates) int {
	n := 0
	for _, c := range coordinates {
		if isBoundary(bm, c.H, c.W) {
			n++
		}
	}

	return n
}

func CalcPhotometricParams(img image.Image) (float64, float64, float64) {
	var rs, gs, bs []uint32
	var rSum, gSum, bSum uint32
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()

			rs = append(rs, r)
			gs = append(gs, g)
			bs = append(bs, b)

			rSum += r
			gSum += g
			bSum += b
		}
	}
	rAve := float64(rSum) / float64(len(rs))
	gAve := float64(gSum) / float64(len(gs))
	bAve := float64(bSum) / float64(len(bs))

	return rAve, gAve, bAve
}

func main() {
	filename, level, err := prepareVars()
	if err != nil {
		log.Fatal(err)
	}

	curDir, _ := os.Getwd()
	path := curDir+"/dsp/lab2/images/"

	img, err := imgio.Open(path+filename+".jpg")
	if err != nil {
		log.Fatal(err)
	}

	binImg := BinarizeImageWithLevel(img, level)
	err = util.SavePNG(binImg, path, filename, "bin_1")
	if err != nil {
		log.Fatal(err)
	}

	img = blur.Gaussian(img, 3.3)
	imgGray := segment.Threshold(img, level)
	err = util.SavePNG(imgGray, path, filename, "bin_2")
	if err != nil {
		log.Fatal(err)
	}

	bm := GetBinMap(*imgGray)
	objs, _ := FindObjectsRec(bm)

	for k, v := range objs {
		s, p, c, e, o := CalcCharacteristics(bm, v)
		fmt.Printf("k: %d \tsquare: %d \tperimeter: %d \tcompact: %.4f \telongation: %.4f \torientation: %.4f\n", k, s, p, c, e, o)
	}

	//bm := GetBinMap(*imgGray)
	//bm = FindObjectsRec(bm)
	//for i := 0; i < len(bm); i++ {
	//	fmt.Println(bm[i])
	//}






	//var mx = [][]byte{
	//	{0,0,1,1,1,1,1,0,0,0,0,0,0,0,0,0,0,0},
	//	{0,0,1,1,1,1,1,1,1,1,0,0,0,0,1,1,1,1},
	//	{0,0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
	//	{0,0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
	//	{0,0,0,0,1,1,1,1,1,1,1,1,1,0,0,0,0,0},
	//	{0,0,0,0,0,0,0,0,1,1,1,1,1,0,0,0,0,0},
	//	{0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1},
	//	{0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1},
	//	{0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1},
	//}
	//
	//o, mx := FindObjectsRec(mx)
	//for i := 0; i < len(mx); i++ {
	//	fmt.Println(mx[i])
	//}
	//
	//for k, v := range o {
	//	fmt.Printf("k: %d v's: %+v\n", k, v)
	//}
	//
	//for k, v := range o {
	//	p := CalcPerim(mx, v)
	//	fmt.Printf("k: %d, p: %d\n", k, p)
	//	//s, p, c, e := CalcCharacteristics(v)
	//	//fmt.Printf("k: %d\tsquare: %d\tperimeter: %d\tcompact: %.4f\telongation: %.4f\n", k, s, p, c, e)
	//}
}


