package eliptic_curve

type Point struct {
	x, y, z int64
}

func NewPoint (x, y, z int64) *Point {
	return &Point{x, y, z}
}
