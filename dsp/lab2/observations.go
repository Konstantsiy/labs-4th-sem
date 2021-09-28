package main

import (
	"fmt"
	"math"
)

// Coordinates2 is a slice of float64
type Coordinates2 []float64

// Observation1 is a data point (float64 between 0.0 and 1.0) in n dimensions
type Observation1 interface {
	Coordinates() Coordinates2
	Distance(point Coordinates2) float64
}

// Observations is a slice of observations
type Observations []Observation1

// Coordinates2 implements the Observation1 interface for a plain set of float64
// coordinates
func (c Coordinates2) Coordinates() Coordinates2 {
	return Coordinates2(c)
}

// Distance returns the euclidean distance between two coordinates
func (c Coordinates2) Distance(p2 Coordinates2) float64 {
	var r float64
	for i, v := range c {
		r += math.Pow(v-p2[i], 2)
	}
	return r
}

// Center returns the center coordinates of a set of Observations
func (c Observations) Center() (Coordinates2, error) {
	var l = len(c)
	if l == 0 {
		return nil, fmt.Errorf("there is no mean for an empty set of points")
	}

	cc := make([]float64, len(c[0].Coordinates()))
	for _, point := range c {
		for j, v := range point.Coordinates() {
			cc[j] += v
		}
	}

	var mean Coordinates2
	for _, v := range cc {
		mean = append(mean, v/float64(l))
	}
	return mean, nil
}

// AverageDistance returns the average distance between o and all observations
func AverageDistance(o Observation1, observations Observations) float64 {
	var d float64
	var l int

	for _, observation := range observations {
		dist := o.Distance(observation.Coordinates())
		if dist == 0 {
			continue
		}

		l++
		d += dist
	}

	if l == 0 {
		return 0
	}
	return d / float64(l)
}