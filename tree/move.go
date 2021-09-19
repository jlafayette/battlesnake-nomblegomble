package tree

type Move uint8

const (
	Up Move = iota
	Down
	Left
	Right
	Dead
)

func (m Move) String() string {
	switch m {
	case Up:
		return "up"
	case Down:
		return "down"
	case Left:
		return "left"
	case Right:
		return "right"
	case Dead:
		return "dead"
	default:
		panic("invalid move")
	}
}

func (m Move) Next(coord Coord) Coord {
	switch m {
	case Up:
		return Coord{X: coord.X, Y: coord.Y + 1}
	case Down:
		return Coord{X: coord.X, Y: coord.Y - 1}
	case Left:
		return Coord{X: coord.X - 1, Y: coord.Y}
	case Right:
		return Coord{X: coord.X + 1, Y: coord.Y}
	case Dead:
		return Coord{X: coord.X, Y: coord.Y}
	default:
		panic("invalid move")
	}
}
