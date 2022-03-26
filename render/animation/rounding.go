package animation

import (
	"math"
)

type Rounding interface {
	Apply(v float64) float64
}

type Round struct{}

func (self Round) Apply(v float64) float64 {
	return math.Round(v)
}

type RoundFloor struct{}

func (self RoundFloor) Apply(v float64) float64 {
	return math.Floor(v)
}

type RoundCeil struct{}

func (self RoundCeil) Apply(v float64) float64 {
	return math.Ceil(v)
}

type RoundNone struct{}

func (self RoundNone) Apply(v float64) float64 {
	return v
}

var DefaultRounding = Round{}
