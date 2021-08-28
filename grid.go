package main

import (
	"github.com/jlafayette/battlesnake-go/q"
	"github.com/jlafayette/battlesnake-go/t"
)

type GridSquare struct {
	x       int
	y       int
	isFood  bool
	visited map[int]int
}

type GridSnake struct {
	isMe   bool
	index  int       // the index of snake in *state
	myfood int       // food eaten so far (this only works on my snake)
	body   []t.Coord // will update as turns go by and snake may eat food
}

func (s *GridSnake) nom() {
	s.myfood += 1
}

func (s *GridSnake) myLength() int {
	return len(s.body) + s.myfood
}

func (s *GridSnake) otherLength(p2 t.Coord, food *[]t.Coord) int {
	// get all the food in the area between start and pos
	p1 := s.body[0]

	xMin := minInt(p1.X, p2.X)
	xMax := maxInt(p1.X, p2.X)
	yMin := minInt(p1.Y, p2.Y)
	yMax := maxInt(p1.Y, p2.Y)

	foodCount := 0
	for _, f := range *food {
		if f.X >= xMin && f.X <= xMax && f.Y >= yMin && f.Y <= yMax {
			foodCount += 1
		}
	}
	return len(s.body) + foodCount
}

func NewGridSnake(i int, b t.Battlesnake) *GridSnake {
	gs := GridSnake{index: i}
	for _, c := range b.Body {
		gs.body = append(gs.body, t.Coord{c.X, c.Y})
	}
	return &gs
}

type Grid struct {
	w       int
	h       int
	myIndex int
	q       q.Fifo
	snakes  map[int]*GridSnake
	squares [][]GridSquare
}

type Area struct {
	Space   int
	Trapped bool
	Target  t.Coord
}

func NewGrid(state *t.GameState) Grid {
	v := make([][]GridSquare, state.Board.Width)
	for i := range v {
		v[i] = make([]GridSquare, state.Board.Height)
	}
	for x := range v {
		for y := range v[x] {
			isFood := false
			for _, food := range state.Board.Food {
				if food.SamePos(t.Coord{X: x, Y: y}) {
					isFood = true
					break
				}
			}
			v[x][y].x = x
			v[x][y].y = x
			v[x][y].isFood = isFood
			v[x][y].visited = make(map[int]int)
		}
	}
	q := q.NewFifo(state.Board.Width * state.Board.Height)
	grid := Grid{w: state.Board.Width, h: state.Board.Height, q: q, squares: v}
	grid.ResetSnakes(state)
	return grid
}

func (g *Grid) ResetSnakes(state *t.GameState) {
	g.snakes = map[int]*GridSnake{}
	for i, srcSnake := range state.Board.Snakes {
		g.snakes[i] = NewGridSnake(i, srcSnake)
		if srcSnake.ID == state.You.ID {
			g.snakes[i].isMe = true
			g.myIndex = i
		}
	}
}

// Use a flood-fill to determine the area after moving in the givin direction.
func (g *Grid) Area(state *t.GameState, move string) Area {
	myStartingCoord := state.You.Head.Moved(move)
	if myStartingCoord.OutOfBounds(state.Board.Width, state.Board.Height) {
		return Area{Space: 0, Trapped: false}
	}

	for snakeIdx, snake := range g.snakes {
		if snake.isMe {
			continue
		}
		// add starting coord for each other snake (you already moved once)
		g.q.Push(q.Item{X: snake.body[0].X, Y: snake.body[0].Y, Turn: 0})
		for {
			// pop item
			item, empty := g.q.Pop()
			x := item.X
			y := item.Y
			if empty {
				break
			}
			_, ok := g.squares[x][y].visited[snakeIdx]
			if ok {
				continue
			}

			// mark as visited
			g.squares[x][y].visited[snakeIdx] = item.Turn
			if g.squares[x][y].isFood {
				snake.nom()
			}

			// find neighbors
			for _, n := range g.findNeighbors(state, item) {
				// add neighbors to queue
				_, ok := g.squares[n.X][n.Y].visited[snakeIdx]
				if !ok {
					g.q.Push(n)
				}
			}
		}
	}

	// do a second pass to count up actual area for your snake
	// add your starting coord
	g.q.Push(q.Item{
		X:    myStartingCoord.X,
		Y:    myStartingCoord.Y,
		Turn: 1, // (turn 1 since already moved once)
	})
	area := 0
	lastOkTurn := 0
	escapeCoord := t.Coord{myStartingCoord.X, myStartingCoord.Y}
	escapeOpenIn := 999
	for {
		// pop item
		item, empty := g.q.Pop()
		x := item.X
		y := item.Y
		if empty {
			break
		}
		// log.Printf("pop %v", item)

		// starting at mySnake + move, traverse neighbors until reaching otherSnake visited on same turn
		// 		do some fancy h2h logic here, can beat other snake?
		// if other snake has visited first, then this is same as a wall, hard stop
		_, iVisited := g.squares[x][y].visited[g.myIndex]
		if iVisited {
			continue
		}
		isOk := true
		for otherIdx, otherSnake := range g.snakes {
			otherTurn, visited := g.squares[x][y].visited[otherIdx]
			if visited {
				otherLen := otherSnake.otherLength(t.Coord{item.X, item.Y}, &state.Board.Food)
				myLen := g.snakes[g.myIndex].myLength()
				willDieInH2H := otherLen >= myLen
				if otherTurn == item.Turn && willDieInH2H {
					// log.Printf("turns are equal and will die in H2H, %v Turns (us: %d, them: %d), Length (us: %d, them: %d)", item, item.turn, otherTurn, myLen, otherLen)
					isOk = false
				} else if otherTurn < item.Turn {
					// log.Printf("they got there first, {X:%d, Y:%d} us: %d, them: %d", item.X, item.Y, item.turn, otherTurn)
					isOk = false
				}
			}
		}
		// you can't magically jump turns, each ok square must lead to the next one.
		// (not sure if this is a perfect system, see TestGrid2)
		if isOk && lastOkTurn >= item.Turn-1 {
			// record area
			area += 1
			lastOkTurn = item.Turn
			// log.Printf("area ok: %#v", item)
		}

		// mark as visited
		g.squares[x][y].visited[g.myIndex] = item.Turn

		// track food eaten so far
		if g.squares[x][y].isFood {
			me, ok := g.snakes[g.myIndex]
			if ok {
				me.nom()
				// fmt.Printf("I ate food and I'm %d long now\n", me.myLength())
			}
		}

		// find neighbors
		for _, n := range g.findNeighbors(state, item) {
			// add neighbors to queue
			_, iVisited := g.squares[n.X][n.Y].visited[g.myIndex]
			if !iVisited && !n.Collide {
				g.q.Push(n)
			}
			if n.Collide {
				// log.Printf("collide: %#v", n)
				if escapeOpenIn > n.OpenIn {
					escapeOpenIn = n.OpenIn
					escapeCoord = t.Coord{n.X, n.Y}
				}
			}
		}
	}
	trapped := area < len(state.You.Body)
	return Area{Space: area, Trapped: trapped, Target: escapeCoord}
}

func (g *Grid) findNeighbors(state *t.GameState, item q.Item) []q.Item {
	// fmt.Println("finding neighbors")
	n := []q.Item{}
	turn := item.Turn + 1
	for _, candidate := range []q.Item{
		{X: item.X - 1, Y: item.Y, Turn: turn, OpenIn: 0},
		{X: item.X + 1, Y: item.Y, Turn: turn, OpenIn: 0},
		{X: item.X, Y: item.Y - 1, Turn: turn, OpenIn: 0},
		{X: item.X, Y: item.Y + 1, Turn: turn, OpenIn: 0},
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

		// self and other snakes
		breakAll := false
		for _, gsnake := range g.snakes {
			snakeLen := gsnake.otherLength(t.Coord{candidate.X, candidate.Y}, &state.Board.Food)
			for i, bc := range gsnake.body {
				// Subtract tails from body as turns go by
				// (should account for eating)
				if i > snakeLen-1-turn {
					// fmt.Printf("break... because %d > %d\n", i, snakeLen-1-turn)
					break
				}
				if bc.SamePos(t.Coord{candidate.X, candidate.Y}) {
					candidate.OpenIn = snakeLen - i
					candidate.Collide = true
					breakAll = true
					break
				}
			}
			if breakAll {
				break
			}
		}
		// if breakAll {
		// 	continue
		// }
		n = append(n, candidate)
	}
	return n
}

func GetArea(state *t.GameState, move string) Area {
	grid := NewGrid(state)
	return grid.Area(state, move)
}
