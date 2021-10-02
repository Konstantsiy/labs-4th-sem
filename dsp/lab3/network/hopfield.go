package network

import "github.com/Konstantsiy/labs-4th-sem/dsp/lab3/converter"

type HopfieldNN struct {
	Dataset [][]byte
	State [][]byte
}

func NewHopfieldNN() *HopfieldNN {
	return &HopfieldNN{}
}

func (nn *HopfieldNN) Learn(dataset [][]byte) {
	nn.Dataset = dataset
	var W [][]byte
	for i, vector := range dataset {
		mx := converter.VectorToMatrix(vector)
		if i != 0 {
			W = converter.AddMatrix(W, mx)
		} else {
			W = mx
		}
	}
	W0 := converter.MatrixDiagonalToZero(W)
	nn.State = W0
}

func (nn *HopfieldNN) Check(v []byte) []byte {
	for _, dV := range nn.Dataset {
		if converter.CompareVectors(dV, v) {
			return dV
		}
	}
	return nil
}
