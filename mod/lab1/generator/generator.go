package generator

import (
	"fmt"
	"gonum.org/v1/plot/plotter"
	"math"
)

func GenerateSequence(r0, n , a, m int) (plotter.Values, error) {
	if m <= a {
		return nil, fmt.Errorf("m must be greater than a")
	}
	if m < 0 || a < 0 {
		return nil, fmt.Errorf("m and a must be non-negative integers")
	}

	//tData := make([][]string, n)
	var values plotter.Values

	for i := 0; i < n; i++ {
		aR0 := r0 * a
		Rn := aR0 % m
		xn := float64(Rn) / float64(m)
		values = append(values, xn)
		//tData = append(tData, []string{strconv.Itoa(i + 1), strconv.Itoa(r0), strconv.Itoa(aR0), strconv.Itoa(Rn), fmt.Sprintf("%.1f", xn)})
		r0 = Rn
	}

	//table := tablewriter.NewWriter(os.Stdout)
	//table.SetHeader([]string{"n", "R(n-1)", "a * R(n-1)", "Rn", "Xn"})
	//for _, data := range tData {
	//	table.Append(data)
	//}
	//table.Render()

	return values, nil
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

func CalcPeriodLength(r0, n, a, m int) (int, error) {
	seq, err := GenerateSequence(r0, n, a, m)
	if err != nil {
		return 0, err
	}

	last := seq.Value(len(seq)-1)
	for i := 0; i < len(seq); i++ {
		if seq.Value(i) == last {
			for j := i + 1; j < len(seq); j++ {
				if seq.Value(j) == seq[i] {
					return j - i, nil
				}
			}
		}
	}

	return len(seq), nil
}

func CalcAperiodLength(seq plotter.Values, period int) int {
	i := 0
	for i + period < len(seq) {
		if seq.Value(i) != seq.Value(i + period) {
			i++
		} else {
			return i + period
		}
	}
	return i + period
}

//func CalcEstMeanSquareDeviation(vrc float64, N int) float64 {
//	return vrc * math.Sqrt(float64(N))
//}
//
//func getMinMax(arr []float64) (min float64, max float64) {
//	min = arr[0]
//	max = min
//
//	for _, v := range arr {
//		if v < min {
//			min = v
//		}
//		if v > max {
//			max = v
//		}
//	}
//
//	return min, max
//}
//
//func СalcIntervalLength(arr []float64, min, max float64) float64 {
//	k := 0
//
//	for _, v := range arr {
//		if v >= min && v <= max {
//			k++
//		}
//	}
//
//	return (max - min) / float64(k)
//}

