package converter

import (
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab3/converter/noise"
	"image"
	"image/color"
	"io/ioutil"
	"strconv"
	"strings"
)

func BinMapToVector(bm [][]byte) []byte {
	var v []byte
	for i := 0; i < len(bm); i++ {
		for j := 0; j < len(bm[0]); j++ {
			v = append(v, bm[i][j])
		}
	}
	return v
}

func CompareVectors(v1, v2 []byte) bool {
	for i := 0; i < len(v1); i++ {
		if v1[i] != v2[i] {
			return false
		}
	}
	return true
}

func MatrixDiagonalToZero(mx [][]byte) [][]byte {
	for i := 0; i < len(mx); i++ {
		mx[i][i] = 0
	}
	return mx
}

func MultiplyMatrixByVector(mx [][]byte, v []byte) []byte {
	vRes := make([]byte, len(v))
	for i := 0; i < len(mx); i++ {
		var s byte
		for j := 0; j < len(v); j++ {
			s += v[i] * mx[i][j]
		}
		vRes[i] = s
	}
	return vRes
}

func AddMatrix(mx1, mx2 [][]byte) [][]byte {
	for i := 0; i < len(mx1); i++ {
		for j := 0; j < len(mx1[0]); j++ {
			mx1[i][j] += mx2[i][j]
		}
	}
	return mx1
}

func VectorToMatrix(v []byte) [][]byte {
	mx := make([][]byte, len(v))
	v1 := v[:]
	for i := 0; i < len(v); i++ {
		mx[i] = make([]byte, len(v))
		for j := 0; j < len(v); j++ {
			mx[i][j] = v[i] * v1[j]
		}
	}
	return mx
}

func TxtToBinaryMap(path string, size int) [][]byte {
	bytes, _ := ioutil.ReadFile(path)
	lines := strings.Split(string(bytes), "\n")
	bm := make([][]byte, size)
	for i := 0; i < 10; i++ {
		bm[i] = make([]byte, size)
		for j := 0; j < 10; j++ {
			b, _ := strconv.ParseInt(string(lines[i][j]), 10, 8)
			bm[i][j] = byte(b)
		}
	}
	return bm
}

func TxtToVector(path string, size int) []byte {
	bm := TxtToBinaryMap(path, size)
	var v []byte
	for i := 0; i < len(bm); i++ {
		for j := 0; j < len(bm[0]); j++ {
			v =  append(v, bm[i][j])
		}
	}
	return v
}

func BinaryMapToImage(bm [][]byte) image.Image {
	leftUp, rightDown := image.Point{}, image.Point{X: len(bm), Y: len(bm[0])}
	img := image.NewRGBA(image.Rectangle{Min: leftUp, Max: rightDown})
	for h := 0; h < len(bm); h++ {
		for w := 0; w < len(bm[0]); w++ {
			if bm[h][w] != 0 {
				img.Set(w, h, color.Black)
			} else {
				img.Set(w, h, color.White)
			}
		}
	}
	return img
}

func TxtToImage(path string, size, percent int) image.Image {
	bm := TxtToBinaryMap(path, size)
	noise.GenerateNoise(bm, percent)
	img := BinaryMapToImage(bm)
	return img
}