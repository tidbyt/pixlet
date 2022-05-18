package animation

type FillMode interface {
	Value() float64
}

type FillModeForwards struct{}

func (self FillModeForwards) Value() float64 {
	return 1.0
}

type FillModeBackwards struct{}

func (self FillModeBackwards) Value() float64 {
	return 0.0
}

var DefaultFillMode = FillModeForwards{}
