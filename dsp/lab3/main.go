package main

import (
	"github.com/Konstantsiy/labs-4th-sem/dsp/lab3/converter"
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
		v := converter.TxtToVector(pathFiles+filename, size)
		dataset = append(dataset, v)
	}

	hnn := network.NewHopfieldNN()
	hnn.Learn(dataset)



	img := converter.TxtToImage(pathFiles+filename, size, noisePercent)
	util.SaveJPG(img, pathImages, filename, "bin")
}
