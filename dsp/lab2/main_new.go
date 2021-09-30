package main

import (
	"fmt"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab1/util"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/segment"
	"github.com/sirupsen/logrus"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
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

type Coordinate struct {
	H int
	W int
}

type Coordinates []Coordinate

func moment(i, j int, wMean, hMean float64, coordinates Coordinates) float64 {
	var result float64
	for _, c := range coordinates {
		result += math.Pow(float64(c.W)-wMean, float64(i)) * math.Pow(float64(c.H)-hMean, float64(j))
	}
	return result
}

type RGB struct {
	R float64
	G float64
	B float64
}

type Photometrics struct {
	AverageRGB RGB
	AverageGray float64
	DispRGB RGB
	DispGray float64
}

func CalcPhotometrics(img image.Image, coordinates Coordinates) Photometrics {
	var rList, gList, bList []uint32
	var grayList []float64
	var rSum, gSum, bSum uint32
	var graySum float64
	for _, c := range coordinates {
		r, g, b, _ := img.At(c.H, c.W).RGBA()
		rSum += r
		gSum += g
		bSum += b

		gray := 0.3 * float64(r) + 0.59 * float64(g) + 0.11 * float64(b)
		graySum += gray

		rList = append(rList, r)
		gList = append(gList, g)
		bList = append(bList, b)
		grayList = append(grayList, gray)
	}
	rAve := float64(rSum) / float64(len(rList))
	gAve := float64(gSum) / float64(len(gList))
	bAve := float64(bSum) / float64(len(bList))
	grayAve := graySum / float64(len(grayList))

	rMin, rMax := GetMinMaxUint32(rList)
	rDisp := rMax - rMin
	gMin, gMax := GetMinMaxUint32(gList)
	gDisp := gMax - gMin
	bMin, bMax := GetMinMaxUint32(bList)
	bDisp := bMax - bMin
	grayMin, grayMax := GetMinMaxFloat64(grayList)
	grayDisp := grayMax - grayMin

	return Photometrics{
		AverageRGB: RGB{R: rAve, G: gAve, B: bAve},
		AverageGray: grayAve,
		DispRGB: RGB{R: float64(rDisp), G: float64(gDisp), B: float64(bDisp)},
		DispGray: grayDisp,
	}
}

func GetMinMaxFloat64(arr []float64) (float64, float64) {
	min, max := arr[0], arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}
	return min, max
}

func GetMinMaxUint32(arr []uint32) (uint32, uint32) {
	min, max := arr[0], arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}
	return min, max
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

type Characteristic struct {
	Square int
	Perimeter int
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

func Bin(img image.Image, level uint8) *image.Gray {
	src := clone.AsRGBA(img)
	bounds := src.Bounds()

	dst := image.NewGray(bounds)

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			srcPos := y*src.Stride + x*4
			dstPos := y*dst.Stride + x

			c := src.Pix[srcPos : srcPos+4]
			r := util.Rank(color.RGBA{c[0], c[1], c[2], c[3]})

			// transparent pixel is always white
			if c[0] == 0 && c[1] == 0 && c[2] == 0 && c[3] == 0 {
				dst.Pix[dstPos] = 0xFF
				continue
			}

			if uint8(r) >= level {
				dst.Pix[dstPos] = 0xFF
			} else {
				dst.Pix[dstPos] = 0x00
			}
		}
	}

	return dst
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

func BinMapToImage(bm [][]byte, img image.Gray) image.Image {
	src := AsRGBA(&img)
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if bm[y][x] != 0 {
				colorID := bm[y][x]
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

type Cluster struct {
	Center Point
	Points []Point
}

func (cluster *Cluster) repositionCenter() {
	var x, y float64
	var clusterCount = len(cluster.Points)
	//fmt.Printf("cluster points count: %d\t", clusterCount)

	for i := 0; i < clusterCount; i++ {
		x += cluster.Points[i].X
		y += cluster.Points[i].Y
	}
	//fmt.Printf("old center: %v\t", cluster.Center)
	cluster.Points = []Point{}
	cluster.Center = Point{x / float64(clusterCount), y / float64(clusterCount)}
	//fmt.Printf("new center: %v\n", cluster.Center)
}

type Point struct {
	X float64
	Y float64
}

func (p Point) Distance(p2 Point) float64 {
	return math.Sqrt(math.Pow(p.X-p2.X, 2) + math.Pow(p.Y-p2.Y, 2))
}

func initClusters(dataset []Point, k int) []Cluster {
	rand.Seed(time.Now().UnixNano())
	var clusters []Cluster

	for i := 0; i < k; i++ {
		center := dataset[rand.Intn(len(dataset))]
		clusters = append(clusters, Cluster{Center: center, Points: []Point{}})
	}

	return clusters
}

func repositionCenters(clusters []Cluster) {
	for i := 0; i < len(clusters); i++ {
		clusters[i].repositionCenter()
	}
}

func RunKMeans(dataset []Point, k int) []Cluster {
	pointsClusterIndex := make([]int, len(dataset))
	clusters := initClusters(dataset, k)
	maxIter := 30
	counter := 0
	for hasChanged := true; hasChanged; {
		hasChanged = false
		for i := 0; i < len(dataset); i++ {
			var minDist float64
			var updatedClusterIndex int
			for j := 0; j < len(clusters); j++ {
				tmpDist := dataset[i].Distance(clusters[j].Center)
				if minDist == 0 || tmpDist < minDist {
					minDist = tmpDist
					updatedClusterIndex = j
				}
			}
			clusters[updatedClusterIndex].Points = append(clusters[updatedClusterIndex].Points, dataset[i])
			if pointsClusterIndex[i] != updatedClusterIndex {
				pointsClusterIndex[i] = updatedClusterIndex
				hasChanged = true
			}
		}
		if hasChanged {
			repositionCenters(clusters)
		}
		counter++
		if counter >= maxIter {
			logrus.Info("exceeded the maximum number of iterations")
			break
		}
	}
	return clusters
}

func main() {
	filename, level, k, _ := prepareVars() // парсинг аргументов: имя файла, уровень бинаризации, кол-во кластеров

	curDir, _ := os.Getwd() // получение текущей папки проекта
	path := curDir+"/dsp/lab2/images/"

	src, _ := imgio.Open(path+filename+".jpg")

	binImg := BinarizeImageWithLevel(src, level) // бинаризация изображения по заданному уровню (default = 200)
	util.SavePNG(binImg, path, filename, "bin_1")

	// избавление от шума + бинаризация (сравнение с простой бинаризацией)
	img := blur.Gaussian(src, 3.3) // применение размытия по Гауссу к исходному изображению
	imgGray := segment.Threshold(img, level)
	util.SavePNG(imgGray, path, filename, "bin_2")

	bm := GetBinMap(*imgGray) // получение бинарной матрицы изображения

	objects, _ := FindObjectsRec(bm) // рекурсивный поиск объектов на бинарной матрице

	// вычисдение геометрических характеристик и фотометрических признаков объектов
	var obj_chars []ObjectCharacteristic
	for k, v := range objects {
		s, p, c, e, o := CalcCharacteristics(bm, v)
		if s == 1 || p == 1 {
			continue
		}
		fmt.Println("--------------------")
		fmt.Printf("object id: %d\n", k)
		fmt.Printf("geometric:\tsquare: %d \tperimeter: %d \tcompact: %.4f \telongation: %.4f \torientation: %.4f\n", s, p, c, e, o)

		obj_chars = append(obj_chars, ObjectCharacteristic{ObjectID: k, Ch: Characteristic{Square: s, Perimeter: p}})
		ph := CalcPhotometrics(src, v)
		fmt.Printf("photometric:\taverage: (%.2f, %.2f, %.2f)\tgray average: %.2f\tdelta: (%.2f, %.2f, %.2f)\tgrey disp: %.2f\n",
			ph.AverageRGB.R, ph.AverageRGB.G, ph.AverageRGB.B, ph.AverageGray,
			ph.DispRGB.R, ph.DispRGB.G, ph.DispRGB.B, ph.DispGray)
		fmt.Println("--------------------")
	}

	var dataset []Point
	for _, o_ch := range obj_chars {
		dataset = append(dataset, Point{float64(o_ch.Ch.Square), float64(o_ch.Ch.Perimeter)})
	}

	clusters := RunKMeans(dataset, k)

	for i, cl := range clusters { // отображение координат центров кластеров
		fmt.Printf("%d centered at (%.f, %.f)\n", i+1, cl.Center.X, cl.Center.Y)
	}

	fmt.Printf("count of objects: %d\n", len(objects))

	// привязываем объект к соответствующему сластеру, чтобы раскрасить все точки,
	// принадлежащие данному объекту, в уникальный цвет
	objects_colors := make(map[byte]int)
	for _, ob_ch := range obj_chars {
		for cl_i, cl := range clusters {
			for _, point := range cl.Points {
				sq := int(math.Round(point.X))
				per := int(math.Round(point.Y))
				if ob_ch.Ch.Square == sq && ob_ch.Ch.Perimeter == per {
					objects_colors[ob_ch.ObjectID] = cl_i+1
					break
				}
			}
		}
	}

	// заполнение бинарной матрицы, исходя из раскраски объектов (их принадлежности опрделенному кластеру)
	for obj_k, cors := range objects {
		if color_i, ok := objects_colors[obj_k]; ok {
			for _, c := range cors {
				bm[c.H][c.W] = byte(color_i)
			}
		}
	}

	// нанесение бинарной матрицы на ранее обработанное черно-белое изображение
	imgRes := BinMapToImage(bm, *imgGray)
	util.SavePNG(imgRes, path, filename, "bin_3")
}