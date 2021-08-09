package main

func samePos(a, b Coord) bool {
	return a.X == b.X && a.Y == b.Y
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// distance returns the manhatten distance between two coordinates.
func distance(p1, p2 Coord) int {
	// In a plane with p1 at (x1, y1) and p2 at (x2, y2), it is |x1 - x2| + |y1 - y2|
	return Abs(p1.X-p2.X) + Abs(p1.Y-p2.Y)
}
