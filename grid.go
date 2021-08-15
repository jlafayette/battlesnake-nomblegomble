package main

type GridSquare struct {
	x       int
	y       int
	isFood  bool
	visited map[int]int
}

func (g *GridSquare) clear() {
	g.visited = make(map[int]int)
	g.isFood = false
}

type GridSnake struct {
	isMe  bool
	index int     // the index of snake in *state
	food  int     // how much food the snake may have eaten
	body  []Coord // will update as turns go by and snake may eat food
}

func (s *GridSnake) nom() {
	s.food += 1
}

func (s *GridSnake) length() int {
	return len(s.body) + s.food
}

func NewGridSnake(i int, b Battlesnake) *GridSnake {
	gs := GridSnake{index: i}
	for _, c := range b.Body {
		gs.body = append(gs.body, Coord{c.X, c.Y})
	}
	return &gs
}

type QItem struct {
	X    int
	Y    int
	turn int
}

type Grid struct {
	w       int
	h       int
	myIndex int
	q       chan QItem
	qlen    int
	snakes  map[int]*GridSnake
	squares [][]GridSquare
}

func NewGrid(state *GameState) Grid {
	v := make([][]GridSquare, state.Board.Height)
	for i := range v {
		v[i] = make([]GridSquare, state.Board.Height)
	}
	for x := range v {
		for y := range v[x] {
			isFood := false
			for _, food := range state.Board.Food {
				if samePos(Coord{x, y}, food) {
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
	q := make(chan QItem, state.Board.Width*state.Board.Height)
	grid := Grid{w: state.Board.Width, h: state.Board.Height, q: q, squares: v}
	grid.ResetSnakes(state)
	return grid
}

func (g *Grid) Reset(state *GameState) {
	g.Clear()
	g.ResetSnakes(state)
}

func (g *Grid) ResetSnakes(state *GameState) {
	g.snakes = map[int]*GridSnake{}
	for i, srcSnake := range state.Board.Snakes {
		g.snakes[i] = NewGridSnake(i, srcSnake)
		if srcSnake.ID == state.You.ID {
			g.snakes[i].isMe = true
			g.myIndex = i
		}
	}
}

func (g *Grid) Clear() {
	for x := range g.squares {
		for y := range g.squares[x] {
			g.squares[x][y].clear()
		}
	}
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
		if snake.isMe {
			continue
		}
		g.push(QItem{X: snake.body[0].X, Y: snake.body[0].Y, turn: 0}) // add starting coord for each other snake (you already moved once)
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
			if g.squares[x][y].isFood {
				snake.nom()
			}

			// need to record more info
			// otherSnake has visited on turn x
			// coord has food (distance)
			// mySnake has visited on turn x
			// if food, and orig = otherSnake, then other snake can grow
			// each turn, the tail of each snake is no longer collidable (unless othersnake can grow=true for that turn)

			// find neighbors
			for _, n := range g.findNeighbors(state, item) {
				// add neighbors to queue
				_, ok := g.squares[n.X][n.Y].visited[snakeIdx]
				if !ok {
					g.push(n)
				}
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
		// log.Printf("pop %v", item)

		// starting at mySnake + move, traverse neighbors until reaching otherSnake visited on same turn
		// 		do some fancy h2h logic here, can beat other snake?
		// if other snake has visited first, then this is same as a wall, hard stop
		_, iVisited := g.squares[x][y].visited[g.myIndex]
		if iVisited {
			continue
		}
		isOk := true
		for otherIdx := range g.snakes {
			otherTurn, visited := g.squares[x][y].visited[otherIdx]
			if visited {
				willDieInH2H := g.snakes[otherIdx].length() >= g.snakes[g.myIndex].length()
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
		g.squares[x][y].visited[g.myIndex] = item.turn
		if g.squares[x][y].isFood {
			me, ok := g.snakes[g.myIndex]
			if ok {
				me.nom()
				// fmt.Printf("I ate food and I'm %d long now\n", me.length())
			}
		}

		// need to record more info
		// otherSnake has visited on turn x
		// coord has food (distance)
		// mySnake has visited on turn x
		// if food, and orig = otherSnake, then other snake can grow
		// each turn, the tail of each snake is no longer collidable (unless othersnake can grow=true for that turn)

		// find neighbors
		for _, n := range g.findNeighbors(state, item) {
			// add neighbors to queue
			_, iVisited := g.squares[n.X][n.Y].visited[g.myIndex]
			if !iVisited {
				g.push(n)
			}
		}
	}

	return area
}

func (g *Grid) findNeighbors(state *GameState, item QItem) []QItem {
	// fmt.Println("finding neighbors")
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
			for i, bc := range gsnake.body {
				// Subtract tails from body as turns go by
				// (should account for eating)
				if i > gsnake.length()-1-turn {
					// fmt.Printf("break... because %d > %d\n", i, gsnake.length()-1-turn)
					break
				}
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

func GetArea(state *GameState, move string) int {
	grid := NewGrid(state)
	return grid.Area(state, move)
}
