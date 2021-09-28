package main

import (
	"fmt"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab1/util"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/segment"
	"log"

	//"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/clone"
	//"github.com/anthonynsimon/bild/imgio"
	//"github.com/anthonynsimon/bild/segment"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	//"math/rand"

	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"strconv"
)

const (
	ColorBlack = 0x00
	ColorWhite = 0xFF
)

var colors = [][3]byte{
	{255, 0, 0},
	{0, 255, 0},
	{0, 0, 255},
	{255, 255, 0},
	{0, 255, 255},
	{255, 0, 255},
	{100, 0, 0},
	{0, 100, 0},
	{0, 0, 100},
	{100, 100, 0},
	{0, 100, 100},
	{100, 0, 100},
	{175, 0, 0},
	{0, 175, 0},
	{0, 0, 175},
	{175, 175, 0},
	{175, 0, 175},
	{0, 175, 175},
}

func prepareVars() (string, uint8, int, error) {
	args := os.Args

	filename := args[1]
	level, _ := strconv.ParseUint(args[2], 10, 8)
	k, _ := strconv.Atoi(args[3])

	return filename, uint8(level), k, nil
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

type Characteristic struct {
	Square int
	Perimeter int
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

func BinMapToImage(bm [][]byte, img image.Gray) image.Image {
	src := AsRGBA(&img)
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if bm[y][x] != 0 {
				colorID := bm[y][x] / 10
				c := colors[colorID]
				pos := y * src.Stride + x * 4

				src.Pix[pos+0] = c[0]
				src.Pix[pos+1] = c[1]
				src.Pix[pos+2] = c[2]
			}
		}
	}
	return src
}

type ObjectCharacteristic struct {
	Ch       Characteristic
	ObjectID byte
}

func main() {
	filename, level, k, err := prepareVars()
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
	objects, _ := FindObjectsRec(bm)

	var obj_chars []ObjectCharacteristic
	for k, v := range objects {
		s, p, c, e, o := CalcCharacteristics(bm, v)
		fmt.Printf("k: %d \tsquare: %d \tperimeter: %d \tcompact: %.4f \telongation: %.4f \torientation: %.4f\n", k, s, p, c, e, o)
		obj_chars = append(obj_chars, ObjectCharacteristic{ObjectID: k, Ch: Characteristic{Square: s, Perimeter: p}})
	}

	var d clusters.Observations
	for _, ch := range obj_chars {
		d = append(d, clusters.Coordinates{float64(ch.Ch.Square), float64(ch.Ch.Perimeter)})
	}
	fmt.Println("")

	km := kmeans.New()

	cls, _ := km.Partition(d, k)
	for i, cl := range cls {
		fmt.Printf("%d centered at (%.f, %.f)\n", i+1, cl.Center[0], cl.Center[1])
	}

	objects_colors := make(map[byte]int)
	for _, ob_ch := range obj_chars {
		for cl_i, cl := range cls {
			for _, obs := range cl.Observations {
				sq := int(math.Round(obs.Coordinates()[0]))
				per := int(math.Round(obs.Coordinates()[1]))
				if ob_ch.Ch.Square == sq && ob_ch.Ch.Perimeter == per {
					objects_colors[ob_ch.ObjectID] = (cl_i+1)*10
					break
				}
			}
		}
	}

	for o_k, v := range objects_colors {
		fmt.Printf("k: %d, clolor: %d\n", o_k, v)
	}

	for obj_k, cors := range objects {
		if color_i, ok := objects_colors[obj_k]; ok {
			for _, c := range cors {
				bm[c.H][c.W] = byte(color_i)
			}
		}
	}

	imgRes := BinMapToImage(bm, *imgGray)
	err = util.SavePNG(imgRes, path, filename, "bin_3")
	if err != nil {
		log.Fatal(err)
	}
}