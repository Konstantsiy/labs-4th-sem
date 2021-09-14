package generator

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"math"
	"os"
	"strconv"
)

func GenerateSequence(R0, n , a, m int) ([]float64, error) {
	if m <= a {
		return nil, fmt.Errorf("m must be greater than a")
	}
	if m < 0 || a < 0 {
		return nil, fmt.Errorf("m and a must be non-negative integers")
	}

	tData := make([][]string, n)

	var res []float64
	for i := 0; i < n; i++ {
		aR0 := R0 * a
		Rn := aR0 % m
		xn := float64(Rn) / float64(m)
		res = append(res, xn)
		tData = append(tData, []string{strconv.Itoa(i + 1), strconv.Itoa(R0), strconv.Itoa(aR0), strconv.Itoa(Rn), fmt.Sprintf("%.1f", xn)})
		R0 = Rn
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"n", "R(n-1)", "a * R(n-1)", "Rn", "Xn"})
	for _, data := range tData {
		table.Append(data)
	}
	table.Render()

	return res, nil
}

func CalcEstMathExpectation(arr []float64) float64 {
	var sum float64 = 0
	for _, v := range arr {
		sum += v
	}

	return sum / float64(len(arr))
}

func CalcEstVariance(arr []float64, exp float64) float64 {
	var sum float64 = 0

	for _, v := range arr {
		sum += math.Pow(v - exp, 2)
	}

	return sum / float64(len(arr) - 1)
}

func CalcEstMeanSquareDeviation(vrc float64, N int) float64 {
	return vrc * math.Sqrt(float64(N))
}

func getMinMax(arr []float64) (min float64, max float64) {
	min = arr[0]
	max = min

	for _, v := range arr {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	return min, max
}

func calcIntervalLength(arr []float64, min, max float64) float64 {
	k := 0

	for _, v := range arr {
		if v >= min && v <= max {
			k++
		}
	}

	return (max - min) / float64(k)
}

func BuildHistogram() {

}

