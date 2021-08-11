package main

type GridSquare struct {
	visited map[int]int
	x       int
	y       int
}

func (g *GridSquare) clear() {
	g.visited = make(map[int]int)
}

type GridSnake struct {
	index int     // the index of snake in *state
	body  []Coord // will update as turns go by and snake may eat food
}

func NewGridSnake(i int, b Battlesnake) GridSnake {
	gs := GridSnake{index: i}
	for _, c := range b.Body {
		gs.body = append(gs.body, Coord{c.X, c.Y})
	}
	return gs
}

type QItem struct {
	X    int
	Y    int
	turn int
}

type Grid struct {
	w       int
	h       int
	q       chan QItem
	qlen    int
	snakes  map[int]GridSnake
	squares [][]GridSquare
}

func NewGrid(state *GameState) Grid {
	v := make([][]GridSquare, state.Board.Height)
	for i := range v {
		v[i] = make([]GridSquare, state.Board.Height)
	}
	for x := range v {
		for y := range v[x] {
			v[x][y].x = x
			v[x][y].y = x
			v[x][y].visited = make(map[int]int)
		}
	}
	q := make(chan QItem, state.Board.Width*state.Board.Height)
	grid := Grid{w: state.Board.Width, h: state.Board.Height, q: q, squares: v}
	grid.ResetSnakes(state)
	return grid
}

func (g *Grid) ResetSnakes(state *GameState) {
	g.snakes = map[int]GridSnake{}
	for i, srcSnake := range state.Board.Snakes {
		g.snakes[i] = NewGridSnake(i, srcSnake)
	}
}

func (g *Grid) Clear() {
	for x := range g.squares {
		for y := range g.squares[x] {
			g.squares[x][y].clear()
		}
	}
	// Clear snakes
}

func (g *Grid) push(item QItem) {
	g.q <- item
	// log.Printf("push %v", item)
	g.qlen += 1
}

func (g *Grid) pop() (QItem, bool) {
	if g.qlen <= 0 {
		// log.Printf("pop empty")
		return QItem{}, true
	}
	g.qlen -= 1
	item := <-g.q
	// log.Printf("pop %v", item)
	return item, false
}

func (g *Grid) Area(state *GameState, move string) int {
	// super simple time
	// for each move, find the area of connected cells
	// for move in moves
	myStartingCoord := newHead(state.You.Head, move)
	if myStartingCoord.outOfBounds(state.Board.Width, state.Board.Height) {
		return 0
	}

	for snakeIdx, snake := range g.snakes {
		if snakeIdx == 0 {
			continue
		}
		g.push(QItem{X: snake.body[0].X, Y: snake.body[0].Y, turn: 0}) // add starting coord for each other snake (you already moved once)
		// area := 0
		for {
			// pop item
			item, empty := g.pop()
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
			g.squares[x][y].visited[snakeIdx] = item.turn

			// need to record more info
			// otherSnake has visited on turn x
			// coord has food (distance)
			// mySnake has visited on turn x
			// if food, and orig = otherSnake, then other snake can grow
			// each turn, the tail of each snake is no longer collidable (unless othersnake can grow=true for that turn)

			// record area
			// area += 1

			// find neighbors
			for _, n := range g.findNeighbors(state, item) {
				// add neighbors to queue
				g.push(n)
			}
		}
	}

	// do a second pass to count up actual area for your snake
	// add your starting coord
	g.push(QItem{
		X:    myStartingCoord.X,
		Y:    myStartingCoord.Y,
		turn: 1, // (turn 1 since already moved once)
	})
	area := 0
	for {
		// pop item
		item, empty := g.pop()
		x := item.X
		y := item.Y
		if empty {
			break
		}

		// starting at mySnake + move, traverse neighbors until reaching otherSnake visited on same turn
		// 		do some fancy h2h logic here, can beat other snake?
		// if other snake has visited first, then this is same as a wall, hard stop
		_, iVisited := g.squares[x][y].visited[0]
		if iVisited {
			continue
		}
		isOk := true
		for otherIdx := range g.snakes {
			otherTurn, visited := g.squares[x][y].visited[otherIdx]
			if visited {
				willDieInH2H := len(g.snakes[otherIdx].body) >= len(g.snakes[0].body)
				if otherTurn == item.turn && willDieInH2H {
					// fmt.Println("turns are equal and will die in H2H")
					isOk = false
				} else if otherTurn < item.turn {
					// log.Printf("they got there first, %v us: %d, them: %d", item, item.turn, otherTurn)
					isOk = false
				}
			}
		}
		if isOk {
			// record area
			area += 1
		}

		// mark as visited
		g.squares[x][y].visited[0] = item.turn

		// need to record more info
		// otherSnake has visited on turn x
		// coord has food (distance)
		// mySnake has visited on turn x
		// if food, and orig = otherSnake, then other snake can grow
		// each turn, the tail of each snake is no longer collidable (unless othersnake can grow=true for that turn)

		// find neighbors
		for _, n := range g.findNeighbors(state, item) {
			// add neighbors to queue
			g.push(n)
		}
	}

	g.Clear()
	g.ResetSnakes(state)
	return area
}

func (g *Grid) findNeighbors(state *GameState, item QItem) []QItem {
	n := []QItem{}
	turn := item.turn + 1
	for _, candidate := range []QItem{
		{X: item.X - 1, Y: item.Y, turn: turn},
		{X: item.X + 1, Y: item.Y, turn: turn},
		{X: item.X, Y: item.Y - 1, turn: turn},
		{X: item.X, Y: item.Y + 1, turn: turn},
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
			// todo: subtract tails from body
			for _, bc := range gsnake.body {
				if samePos(Coord{candidate.X, candidate.Y}, bc) {
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
