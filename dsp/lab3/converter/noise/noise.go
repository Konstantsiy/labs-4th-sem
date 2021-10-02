package noise

import (
	"math"
	"math/rand"
)

func GenerateNoise(bm [][]byte, percent int) [][]byte {
	h, w := len(bm), len(bm[0])
	noises := int(math.Round(float64(percent) / float64(h * w) * 100))
	used := map[int]struct{}{}

	for i := 0; i < noises; i++ {
		index := 0
		for {
			index = rand.Intn(h * w)
			if _, ok := used[index]; !ok {
				used[index] = struct{}{}
				break
			}
		}
		hNoise, wNoise := index / h, index % h
		if bm[hNoise][wNoise] == 0 {
			bm[hNoise][wNoise] = 1
		} else {
			bm[hNoise][wNoise] = 0
		}
	}

	return bm
}
