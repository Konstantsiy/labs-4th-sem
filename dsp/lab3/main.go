package main

import (
	"fmt"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab3/converter"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab3/converter/noise"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab3/network"
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab3/util"
	"os"
)

func main() {
	filename, size, noisePercent := util.PrepareVars()

	curDir, _ := os.Getwd()
	pathFiles := curDir+"/dsp/lab3/files/"
	pathImages := curDir+"/dsp/lab3/images/"
	filename += ".txt"

	var dataset [][]byte
	for i := 1; i <= 3; i++ {
		fPath := fmt.Sprintf("%s%d%s", pathFiles, i, ".txt")
		fmt.Printf("file: %s\n", fPath)
		v := converter.TxtToVector(fPath, size)
		dataset = append(dataset, v)
	}

	hnn := network.NewHopfieldNN()
	hnn.Learn(dataset)

	testBm := converter.TxtToBinaryMap(pathFiles+filename, size)
	testBm = noise.GenerateNoise(testBm, 10)
	testV := converter.BinMapToVector(testBm)

	fmt.Printf("result: %t\n", hnn.Check(testV))


	img := converter.TxtToImage(pathFiles+filename, size, noisePercent)
	util.SaveJPG(img, pathImages, filename, "bin")
}
