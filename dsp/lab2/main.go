package main

import (
	"fmt"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab1/util"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/segment"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
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
)

var colors = [][3]byte{
	{255, 0, 0},
	{0, 255, 0},
}

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
}

type Coordinates1 []Coordinate

func FindObjects(binMap [][]byte) (map[byte]Coordinates1, [][]byte) {
	height, width := len(binMap), len(binMap[0])
	objects := make(map[byte]Coordinates1)
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
					objects[cur] = Coordinates1{}
				}
				objects[cur] = append(objects[cur], Coordinate{H: i, W: j})
			} else if B != 0 && C == 0 {
				binMap[i][j] = B
				if _, ok := objects[B]; !ok {
					objects[B] = Coordinates1{}
				}
				objects[B] = append(objects[B], Coordinate{H: i, W: j})
			} else if B == 0 && C != 0 {
				binMap[i][j] = C
				if _, ok := objects[C]; !ok {
					objects[C] = Coordinates1{}
				}
				objects[C] = append(objects[C], Coordinate{H: i, W: j})
			} else if B != 0 && C != 0 {
				binMap[i][j] = B
				if _, ok := objects[B]; !ok {
					objects[B] = Coordinates1{}
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

func moment(i, j int, wMean, hMean float64, coordinates Coordinates1) float64 {
	var result float64
	for _, c := range coordinates {
		result += math.Pow(float64(c.W)-wMean, float64(i)) * math.Pow(float64(c.H)-hMean, float64(j))
	}
	return result
}

func CalcCharacteristics(bm [][]byte, coordinates Coordinates1) (int, int, float64, float64, float64) {
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

func getPixel(img image.RGBA, x, y int) []uint8 {
	pos := y * img.Stride + x * 4
	return img.Pix[pos:pos+4]
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

func fill(bm [][]byte, x, y int, c byte, objects map[byte]Coordinates1) {
	if bm[x][y] == 1 {
		bm[x][y] = c
		if _, ok := objects[c]; !ok {
			objects[c] = Coordinates1{}
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

func FindObjectsRec(bm [][]byte) (map[byte]Coordinates1, [][]byte) {
	objects := make(map[byte]Coordinates1)
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

func CalcPerim(bm [][]byte, coordinates Coordinates1) int {
	n := 0
	for _, c := range coordinates {
		if isBoundary(bm, c.H, c.W) {
			n++
		}
	}

	return n
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

	mx := GetBinMap(*imgGray)
	o, _ := FindObjectsRec(mx)

	//for k, v := range objects {
	//	s, p, c, e, o := CalcCharacteristics(bm, v)
	//	fmt.Printf("k: %d \tsquare: %d \tperimeter: %d \tcompact: %.4f \telongation: %.4f \torientation: %.4f\n", k, s, p, c, e, o)
	//}
	//
	//var cors []Cor2C
	//for _, coordinates := range objects {
	//	for i, c := range coordinates {
	//		cors = append(cors, Cor2C{Cor: Coordinate{c.H, c.W}, C: i})
	//	}
	//}
	//
	////cls := InitClusters(4, len(bm[0]), len(bm))
	//
	////KMeans(cls, cors)
	//
	//var n byte = 0
	//var curC int
	//for _, c := range cors {
	//	if c.C != curC {
	//		n++
	//	}
	//	bm[c.Cor.H][c.Cor.W] = n
	//}
	//
	//for i := 0; i < len(bm); i++ {
	//	fmt.Println(bm[i])
	//}


	//var mx = [][]byte{
	//	{0,0,0,0,0,0,0,1,1,1},
	//	{0,0,1,1,1,0,0,1,1,1},
	//	{0,0,1,1,1,0,0,1,1,1},
	//	{0,0,1,1,1,0,0,1,1,1},
	//	{0,0,0,0,0,0,0,0,0,0},
	//	{0,0,0,0,0,0,0,0,0,0},
	//	{0,0,0,0,0,1,1,1,1,1},
	//	{0,0,0,0,0,1,1,1,1,1},
	//	{0,0,0,0,0,1,1,1,1,1},
	//	{0,0,0,0,0,0,0,1,1,1},
	//	//{0,0,1,1,1,1,1,0,0,0,0,0,0,0,0,0,0,0},
	//	//{0,0,1,1,1,1,1,1,1,0,0,0,0,0,1,1,1,1},
	//	//{0,0,1,1,1,1,1,1,1,0,0,0,0,0,1,1,1,1},
	//	//{0,0,1,1,1,1,1,1,1,0,0,0,0,0,1,1,1,1},
	//	//{0,0,0,0,1,1,0,0,0,0,0,0,0,0,0,0,0,0},
	//	//{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
	//	//{0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1},
	//	//{0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1},
	//	//{0,0,0,0,0,0,0,0,1,1,1,1,1,1,1,1,1,1},
	//}

	//o, _ := FindObjectsRec(mx)
	//var points []Point
	//for _, cors := range o {
	//	for _, c := range cors {
	//		points = append(points, Point{H: c.H, W: c.W, C: 0})
	//	}
	//}

	var d clusters.Observations
	for _, v :=  range o {
		for _, c := range v {
			d = append(d, clusters.Coordinates{
				float64(c.W),
				float64(c.H),
			})
		}
	}

	km := kmeans.New()

	var t byte = 1
	cls, _ := km.Partition(d, 8)
	for _, c := range cls {
		//fmt.Printf("%d Centered at (%.f, %.f)\n", i+1, c.Center[0], c.Center[1])
		//fmt.Printf("Matching data points: %+v\n", c.Observations)
		//
		//fmt.Printf("Matching round points: ")
		//for _, o := range c.Observations {
		//	fmt.Printf("[%d %d] ", int(math.Round(o.Coordinates()[0])), int(math.Round(o.Coordinates()[1])))
		//}
		//fmt.Println("")

		for _, o := range c.Observations {
			h := o.Coordinates()[0]
			w := o.Coordinates()[1]
			mx[int(math.Round(w))][int(math.Round(h))] = t*10
		}
		t++
	}

	for j := 0; j < len(mx); j++ {
		fmt.Println(mx[j])
	}
}

//type Point struct {
//	H int
//	W int
//	C int
//}
//func (p *Point) FindSameMarkPoints(points []Point) []Point {
//	var result []Point
//	for _, point := range points {
//		if p.C == point.C {
//			result = append(result, point)
//		}
//	}
//	return result
//}
//
//func (p *Point) FindNearestCluster(clusters []Point) {
//	dests := make(map[int]float64)
//	for i, c := range clusters {
//		dests[i] = p.calcDestTo(c)
//	}
//	minI := getMin(dests)
//	p.C = clusters[minI].C
//}
//
//func (p *Point) calcDestTo(anotherPoint Point) float64 {
//	return math.Sqrt(math.Pow(float64(p.H-anotherPoint.H), 2) + math.Pow(float64(p.W-anotherPoint.W), 2))
//}
//
//func getMin(dests map[int]float64) int {
//	min := dests[0]
//	minI := 0
//	for i, v := range dests {
//		if v < min {
//			minI = i
//		}
//	}
//	return minI
//}
//
//func CalculateCenterMass(points []Point) Point {
//	var sumW, sumH int
//	for _, p := range points {
//		sumW += p.W
//		sumH += p.H
//	}
//	aveW, aveH := sumW / len(points), sumH / len(points)
//	return Point{H: aveH, W: aveW}
//}
//
//func GenerateClusters(k, w, h int) []Point {
//	var points []Point
//	rand.Seed(time.Now().UnixNano())
//	for i := 0; i < k; i++ {
//		rW := rand.Intn(w)
//		rH := rand.Intn(h)
//		points = append(points, Point{H: rH, W: rW, C: (i+1) * 10})
//	}
//	return points
//}
//
//func KMeansProcess(pointsRes []Point, countK, w, h, iters int) []Point {
//	clusters := GenerateClusters(countK, w, h)
//	var points []Point
//	for i := 0; i < 1; i++ {
//		fmt.Printf("%d iteration\n", i+1)
//
//		points = points[:0]
//
//		for _, p := range pointsRes {
//			p.FindNearestCluster(clusters)
//			points = append(points, p)
//		}
//
//		fmt.Printf("%+v", points)
//
//
//		c_points := make(map[Point][]Point)
//
//		for _, c := range clusters {
//			ps := c.FindSameMarkPoints(points)
//			fmt.Printf("c: %v, c_ps: %+v\n", c, ps)
//			c_points[c] = ps
//		}
//
//		for c, p := range c_points {
//			fmt.Printf("cluster: %v, points: %+v\n", c, p)
//		}
//
//		var j = 0
//		for k, v := range c_points {
//			new_c := CalculateCenterMass(v)
//			if k != new_c {
//				clusters[j] = new_c
//				j++
//			}
//		}
//
//	}
//	return points
//}