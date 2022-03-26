package animation

type Direction interface {
	FrameCount(delay, duration int) int
	Progress(delay, duration int, fill float64, frameIdx int) float64
}

type DirectionImpl struct {
	Alternate bool
	Reverse   bool
}

func (self DirectionImpl) FrameCount(delay, duration int) int {
	if self.Alternate {
		return 2 * (delay + duration)
	}

	return delay + duration + delay
}

func (self DirectionImpl) Progress(delay, duration int, fill float64, frameIdx int) (progress float64) {
	idx1 := delay
	idx2 := delay + duration
	idx3 := delay + duration + delay
	idx4 := delay + duration + delay + duration

	if frameIdx < idx1 {
		progress = 0.0
	} else if frameIdx < idx2 {
		progress = float64(frameIdx-idx1) / float64(duration)
	} else if frameIdx < idx3 {
		progress = 1.0
	} else if self.Alternate && frameIdx < idx4 {
		progress = float64(frameIdx-idx3) / float64(duration)
		progress = 1.0 - progress
	} else {
		progress = fill
	}

	if self.Reverse {
		progress = 1.0 - progress
	}

	return
}

var DirectionNormal = DirectionImpl{Alternate: false, Reverse: false}
var DirectionReverse = DirectionImpl{Alternate: false, Reverse: true}
var DirectionAlternate = DirectionImpl{Alternate: true, Reverse: false}
var DirectionAlternateReverse = DirectionImpl{Alternate: true, Reverse: true}
var DefaultDirection = DirectionNormal
