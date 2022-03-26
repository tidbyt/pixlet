package animation

type Vec2f struct {
	X float64 `starlark:"x,required"`
	Y float64 `starlark:"y,required"`
}

func (lhs Vec2f) Lerp(rhs Vec2f, progress float64) Vec2f {
	return Vec2f{
		Lerp(lhs.X, rhs.X, progress),
		Lerp(lhs.Y, rhs.Y, progress),
	}
}
