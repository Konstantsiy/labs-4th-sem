package main

import (
	"fmt"
	"github.com/Konstantsiy/labs-4th-sem/mod/lab1/generator"
	"log"
	"math"
	"os"
	"strconv"
)

func prepareVars() (float64, float64, float64, float64, error) {
	args := os.Args

	if len(args) < 3 {
		return 0, 0, 0, 0, fmt.Errorf("not enough arguments")
	}

	r0, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	n, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	a, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	m, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return r0, n, a, m, nil
}

func main() {
	r0 := 1
	n := 1000000
	a := 134279
	m := 337109

	seq, err := generator.GenerateSequence(r0, n, a, m)
	if err != nil {
		log.Fatal(err)
	}

	mathMean := generator.CalcEstMathExpectation(seq)
	fmt.Printf("mean math expectation: %.4f\n", mathMean)

	dispMean := generator.CalcEstVariance(seq, mathMean)
	fmt.Printf("disp mean expectation: %.4f\n", dispMean)

	d := math.Sqrt(dispMean)
	fmt.Printf("average square deviation: %.4f\n", d)

	period, _ := generator.CalcPeriodLength(r0, n, a, m)
	fmt.Printf("period length: %d\n", period)

	aperiod := generator.CalcAperiodLength(seq, period)
	fmt.Printf("aperiod length: %d\n", aperiod)

	//err = hist.DrawHistogram(seq)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
