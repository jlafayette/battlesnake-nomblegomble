package main

type VState struct {
	Snakes map[string]bool
}

type Grid struct {
	w       int
	h       int
	q       chan Coord
	qlen    int
	visited [][]bool
}

func NewGrid(w, h int) Grid {
	v := make([][]bool, h)
	for i := range v {
		v[i] = make([]bool, h)
	}
	q := make(chan Coord, w*h)
	return Grid{w: w, h: h, q: q, visited: v}
}

func (g *Grid) Clear() {
	for x := range g.visited {
		for y := range g.visited[x] {
			g.visited[x][y] = false
		}
	}
}

func (g *Grid) push(c Coord) {
	g.q <- c
	// log.Printf("push %v", c)
	g.qlen += 1
}

func (g *Grid) pop() (Coord, bool) {
	if g.qlen <= 0 {
		// log.Printf("pop empty")
		return Coord{X: 0, Y: 0}, true
	}
	g.qlen -= 1
	c := <-g.q
	// log.Printf("pop %v", c)
	return c, false
}

func (g *Grid) Area(state *GameState, move string) int {
	// super simple time
	// for each move, find the area of connected cells
	// for move in moves
	startingCoord := newHead(state.You.Head, move)
	if startingCoord.outOfBounds(state.Board.Width, state.Board.Height) {
		return 0
	}
	g.push(startingCoord)
	area := 0
	for {
		// pop item
		coord, empty := g.pop()
		if empty {
			break
		}
		if g.visited[coord.X][coord.Y] {
			continue
		}

		// mark as visited
		g.visited[coord.X][coord.Y] = true

		// record area
		area += 1

		// find neighbors
		for _, n := range g.findNeighbors(state, coord) {
			// add neighbors to queue
			g.push(n)
		}
	}
	g.Clear()
	return area
}

func (g *Grid) findNeighbors(state *GameState, coord Coord) []Coord {
	n := []Coord{}
	for _, candidate := range []Coord{
		{X: coord.X - 1, Y: coord.Y},
		{X: coord.X + 1, Y: coord.Y},
		{X: coord.X, Y: coord.Y - 1},
		{X: coord.X, Y: coord.Y + 1},
	} {
		// walls
		if candidate.X < 0 {
			continue
		} else if candidate.X >= g.w {
			continue
		}
		if candidate.Y < 0 {
			continue
		} else if candidate.Y >= g.h {
			continue
		}
		// log.Printf("%v is not a wall", candidate)

		// self and other snakes
		breakAll := false
		for _, snake := range state.Board.Snakes {
			for _, bc := range snake.Body {
				if samePos(candidate, bc) {
					breakAll = true
					break
				}
			}
			if breakAll {
				break
			}
		}
		if breakAll {
			continue
		}
		n = append(n, candidate)
	}
	return n
}
