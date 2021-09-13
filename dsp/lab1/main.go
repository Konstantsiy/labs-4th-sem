package main

import (
	"labs/dsp/lab1/hist"
	"log"
	"os"
)

func main() {
	curDir, _ := os.Getwd()
	path := curDir+"/dsp/lab1/images/"

	go func(){
		err := hist.DrawHistogram(path, "1.jpg")
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := hist.CalcHistogramComponents(path, "1.jpg")
	if err != nil {
		log.Fatal(err)
	}
}


