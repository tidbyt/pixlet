package render

type Path interface {
	Length() int
	Size() (int, int)
	Point(i int) (int, int)
}

type PathPoint struct {
	X int
	Y int
}

// CircularPath draws a circle
// XXX: from where, in what direction?
type CircularPath struct {
	Radius int
	points []PathPoint
}

func (cp *CircularPath) computePoints() {
	R := cp.Radius - 1
	x, y := R, 0
	xc, yc, re := 1-2*R, 0, 0

	cp.points = make([]PathPoint, 0, 5*R)

	// Compute the lower NE octant
	for y <= x {
		cp.points = append(cp.points, PathPoint{X: x, Y: y})
		y += 1
		re += yc
		yc += 2
		if 2*re+xc > 0 {
			x -= 1
			re += xc
			xc += 1
		}
	}

	// Then the complete NE quadrant by reflecting the octant,
	// except the last point if there's an odd number of them (the
	// case where the final point is shared).
	octLen := len(cp.points)
	for i := octLen - 1; i >= 0; i-- {
		if cp.points[i].X == cp.points[i].Y {
			continue
		}
		cp.points = append(
			cp.points,
			PathPoint{
				X: cp.points[i].Y,
				Y: cp.points[i].X,
			},
		)
	}

	// The NE quadrant can now be mirrored to the other 3.
	quadLen := len(cp.points)

	// NW
	for _, ne := range cp.points[0:quadLen] {
		cp.points = append(cp.points, PathPoint{X: -ne.Y - 1, Y: ne.X})
	}

	// SW
	for _, ne := range cp.points[0:quadLen] {
		cp.points = append(cp.points, PathPoint{X: -ne.X - 1, Y: -ne.Y - 1})
	}

	// SE
	for _, ne := range cp.points[0:quadLen] {
		cp.points = append(cp.points, PathPoint{X: ne.Y, Y: -ne.X - 1})
	}
}

func (cp *CircularPath) Length() int {
	if len(cp.points) == 0 {
		cp.computePoints()
	}
	return len(cp.points)
}

func (cp *CircularPath) Size() (int, int) {
	return cp.Radius * 2, cp.Radius * 2
}

func (cp *CircularPath) Point(i int) (int, int) {
	if len(cp.points) == 0 {
		cp.computePoints()
	}
	i = ModInt(i, len(cp.points))
	return cp.Radius + cp.points[i].X, cp.Radius + cp.points[i].Y
}

// PolyLine draws straight lines passing through a list of vertices
type PolyLine struct {
	Vertices []PathPoint
	path     []PathPoint
}

func (pl *PolyLine) Length() int {
	if len(pl.path) == 0 {
		pl.compute()
	}

	return len(pl.path)
}

func (pl *PolyLine) Size() (int, int) {
	return 0, 0
}

func (pl *PolyLine) Point(i int) (int, int) {
	if len(pl.path) == 0 {
		pl.compute()
	}
	return pl.path[i].X, pl.path[i].Y
}

func (pl *PolyLine) compute() {
	pl.path = []PathPoint{}
	for i := 0; i < len(pl.Vertices)-1; i++ {
		pl.addLineSegment(
			pl.Vertices[i].X,
			pl.Vertices[i].Y,
			pl.Vertices[i+1].X,
			pl.Vertices[i+1].Y,
		)
	}
}

func (pl *PolyLine) addLineSegment(x0, y0, x1, y1 int) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}

	sx := -1
	if x0 < x1 {
		sx = 1
	}

	dy := y1 - y0
	if dy > 0 {
		dy = -dy
	}

	sy := -1
	if y0 < y1 {
		sy = 1
	}

	if dx == 0 {
		for ; y0 != y1; y0 += sy {
			pl.path = append(pl.path, PathPoint{X: x0, Y: y0})
		}
		pl.path = append(pl.path, PathPoint{X: x0, Y: y0})
		return
	}

	if dy == 0 {
		for ; x0 != x1; x0 += sx {
			pl.path = append(pl.path, PathPoint{X: x0, Y: y0})
		}
		pl.path = append(pl.path, PathPoint{X: x0, Y: y0})
		return
	}

	err := dx + dy

	pl.path = append(pl.path, PathPoint{X: x0, Y: y0})
	for x0 != x1 || y0 != y1 {
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
		pl.path = append(pl.path, PathPoint{X: x0, Y: y0})
	}
}
