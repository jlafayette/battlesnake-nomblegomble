package tree

type Coord struct {
	X int
	Y int
}

func (a Coord) Equals(b Coord) bool {
	return a.X == b.X && a.Y == b.Y
}

func (c Coord) Move(m Move) Coord {
	return m.Next(c)
}

func (c Coord) InBounds(width, height int) bool {
	return c.X >= 0 && c.X < width && c.Y >= 0 && c.Y < height
}

func remove(slice []Coord, s int) []Coord {
	return append(slice[:s], slice[s+1:]...)
}
