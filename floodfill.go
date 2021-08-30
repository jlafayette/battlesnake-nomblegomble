package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jlafayette/battlesnake-go/t"
)

type Turn int32
type SnakeIndex int32
type BodyIndex int32

type Contents int32

const (
	Empty Contents = iota
	Head
	Body
	Tail
	DoubleTail
	Food
	H2H
)

type Cell struct {
	// Keeping track of state.
	// State is a bit weird since snakes can overlap on h2h. But... they can't go through
	// The body of other snakes
	X         int
	Y         int
	snakeId   SnakeIndex // index of current snake
	bodyIndex BodyIndex  // 0 is head len(body)-1 is the tail
	contents  Contents   // head,body,tail, empty, food

	prevContents Contents

	// Keeping track of when snakes visited
	visited map[SnakeIndex]Turn // key is index of snake, value is turn visited on
	h2h     map[Turn]SnakeIndex // key is turn H2H occurs on, value is snake index that won (-1 for both die)
}

func NewCell(x, y int) *Cell {
	visited := make(map[SnakeIndex]Turn)
	h2h := make(map[Turn]SnakeIndex)
	return &Cell{
		X:        x,
		Y:        y,
		contents: Empty,
		visited:  visited,
		h2h:      h2h,
	}
}

func (c *Cell) Blocked() bool {
	return c.prevContents == Head || c.prevContents == Body || c.prevContents == DoubleTail
}

func (c *Cell) IsHead() bool {
	return c.prevContents == Head
}

func (c *Cell) IsFood() bool {
	return c.prevContents == Food
}

func (c *Cell) IsSnake() bool {
	return c.prevContents == Head || c.prevContents == Body || c.prevContents == Tail || c.prevContents == DoubleTail
}

func (c *Cell) UpdateSnake(ate bool, length int32) {
	if !c.IsSnake() {
		return
	}
	c.bodyIndex += 1
	switch c.prevContents {
	case Head:
		c.contents = Body
	case Body:
		// check if we are the tail now...
		if ate {
			if c.bodyIndex == BodyIndex(length-2) {
				c.contents = DoubleTail
			}
		} else {
			if c.bodyIndex == BodyIndex(length-1) {
				c.contents = Tail
			}
		}
	case Tail:
		c.bodyIndex = -1
		c.contents = Empty
	case DoubleTail:
		if ate {
			c.contents = DoubleTail
		} else {
			c.contents = Tail
		}
	}
}

func (c *Cell) SnakeId() SnakeIndex {
	return c.snakeId
}

func (c *Cell) SetSnake(snakeIndex SnakeIndex, bodyIndex BodyIndex, length int32, turn Turn) {
	isHead := bodyIndex == 0
	isTail := int32(bodyIndex) == length-1
	if isHead {
		c.contents = Head
	} else if isTail {
		c.contents = Tail
	} else {
		c.contents = Body
	}
	// TODO: detect double tail properly

	c.bodyIndex = bodyIndex
	c.snakeId = snakeIndex
	c.visited[snakeIndex] = turn
}

func (c *Cell) SetFood() {
	c.contents = Food
}

func (c *Cell) NewHeadFrom(nc *Cell, turn Turn) {
	c.contents = Head
	c.snakeId = nc.snakeId
	c.bodyIndex = 0
	c.visited[c.snakeId] = turn
}

func (c *Cell) NewTurn() {
	c.prevContents = c.contents
}

func (c *Cell) String() string {
	switch c.prevContents {
	case Food:
		return "  fff"
	case Head:
		return "  H" + strconv.Itoa(int(c.snakeId)) + strconv.Itoa(int(c.bodyIndex))
	case Body:
		return "  B" + strconv.Itoa(int(c.snakeId)) + strconv.Itoa(int(c.bodyIndex))
	case Tail:
		return "  T" + strconv.Itoa(int(c.snakeId)) + strconv.Itoa(int(c.bodyIndex))
	case DoubleTail:
		return "  D" + strconv.Itoa(int(c.snakeId)) + strconv.Itoa(int(c.bodyIndex))
	case H2H:
		return "  XXX"
	case Empty:
		return "  [ ]"
	}
	return "     "
}

type Board struct {
	Width    int
	Height   int
	Turn     Turn
	Cells    []*Cell
	lengths1 map[SnakeIndex]int32
	lengths  map[SnakeIndex]int32
	areas    map[SnakeIndex]int32
	ate      map[SnakeIndex]bool
}

func (b *Board) getCell(x, y int) (*Cell, bool) {
	if x < 0 || x > b.Width-1 {
		return nil, false
	}
	if y < 0 || y > b.Height-1 {
		return nil, false
	}
	index := x + y*b.Width
	return b.Cells[index], true
}

func NewBoard(state *t.GameState) *Board {
	w := state.Board.Width
	h := state.Board.Height
	snakeNumber := len(state.Board.Snakes)
	cells := make([]*Cell, 0, w*h)
	lengths := make(map[SnakeIndex]int32, snakeNumber)
	lengths1 := make(map[SnakeIndex]int32, snakeNumber)
	ate := make(map[SnakeIndex]bool, snakeNumber)
	areas := make(map[SnakeIndex]int32, snakeNumber)
	b := &Board{
		Width:    w,
		Height:   h,
		Turn:     0,
		Cells:    cells,
		lengths:  lengths,
		lengths1: lengths1,
		areas:    areas,
		ate:      ate,
	}

	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			b.Cells = append(b.Cells, NewCell(x, y))
		}
	}
	for snakeIndex, snake := range state.Board.Snakes {
		b.lengths[SnakeIndex(snakeIndex)] = snake.Length
		b.lengths1[SnakeIndex(snakeIndex)] = snake.Length
		for bodyIndex, bCoord := range snake.Body {
			cell, _ := b.getCell(bCoord.X, bCoord.Y)
			cell.SetSnake(SnakeIndex(snakeIndex), BodyIndex(bodyIndex), snake.Length, Turn(0))
		}
	}
	for _, foodCoord := range state.Board.Food {
		cell, _ := b.getCell(foodCoord.X, foodCoord.Y)
		cell.SetFood()
	}
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			cell, _ := b.getCell(x, y)
			cell.NewTurn()
		}
	}

	return b
}

func (b *Board) Update() bool {
	done := true
	b.Turn += 1

	for k := range b.ate {
		b.ate[k] = false
	}

	// First pass, mark all the squares where the new heads are
	// For each snake, if there are no viable moves (other snake got there first, isBody, outOfBounds)
	// If all snakes have no more moves, then done=true
	// New head square, SnakeId, SnakeBody=0, SnakeHead=true
	// Record which snakes should grow (ate food)
	// Update Visited
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			cell, exists := b.getCell(x, y)
			if !exists {
				continue
			}
			// check for early return
			if cell.Blocked() {
				continue
			}

			// Check for nearby heads
			nCell, ok := b.getCell(x-1, y)
			if ok && nCell.IsHead() {
				id := nCell.SnakeId()
				cell.NewHeadFrom(nCell, b.Turn)
				if cell.IsFood() {
					b.ate[id] = true
				}
				b.areas[id] += 1
				done = false
			}
			nCell, ok = b.getCell(x+1, y)
			if ok && nCell.IsHead() {
				id := nCell.SnakeId()
				cell.NewHeadFrom(nCell, b.Turn)
				if cell.IsFood() {
					b.ate[id] = true
				}
				b.areas[id] += 1
				done = false
			}
			nCell, ok = b.getCell(x, y+1)
			if ok && nCell.IsHead() {
				id := nCell.SnakeId()
				cell.NewHeadFrom(nCell, b.Turn)
				if cell.IsFood() {
					b.ate[id] = true
				}
				b.areas[id] += 1
				done = false
			}
			nCell, ok = b.getCell(x, y-1)
			if ok && nCell.IsHead() {
				id := nCell.SnakeId()
				cell.NewHeadFrom(nCell, b.Turn)
				if cell.IsFood() {
					b.ate[id] = true
				}
				b.areas[id] += 1
				done = false
			}
		}
	}

	// update lengths
	for snakeid, ate := range b.ate {
		if ate {
			b.lengths[snakeid] += 1
		}
	}

	// Update Body index
	// Body index goes down by one
	// Tail becomes clear unless it's a double tail
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			cell, _ := b.getCell(x, y)
			if cell.IsSnake() {
				ate := b.ate[cell.SnakeId()]
				length := b.lengths[cell.SnakeId()]
				cell.UpdateSnake(ate, length)
			}
		}
	}

	// Update H2H
	// Because space is already awarded, this should take away space from the smaller snake
	// or from both if they both die

	// Resolve the turn for each cell
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			cell, ok := b.getCell(x, y)
			if !ok {
				continue
			}
			cell.NewTurn()
		}
	}

	return done
}

func (b *Board) Fill() map[SnakeIndex]*FloodFillResult {
	results := make(map[SnakeIndex]*FloodFillResult)
	for snakeIndex := range b.lengths {
		results[snakeIndex] = &FloodFillResult{}
	}

	// Loop until end condition is set
	// TODO: allow revisiting the same square over and over
	//       each time it should count for an extra space (like going in a
	//       circle forever is fine)
	//       but... we need to find an alternate way of terminating the loop
	//       option1: check every x turns if any progress is being made/enough area is stacked out
	//       option2: check every update for some terminating condition
	done := false
	for !done {
		fmt.Println(b.String())
		done = b.Update()
		if int(b.Turn) >= b.Width*b.Height {
			done = true
		}
	}
	fmt.Printf("Ended updates at turn: %d\n", b.Turn)
	fmt.Println(b.String())

	// Go over the board and get all the results
	for id, area := range b.areas {
		results[id].Area = int(area)
	}
	for snakeIndex, length := range b.lengths {
		origLen := b.lengths1[snakeIndex]
		food := length - origLen
		results[snakeIndex].Food = int(food)
	}
	return results
}

func (b *Board) String() string {

	var sb strings.Builder
	for y := b.Height - 1; y >= 0; y-- {
		for x := 0; x < b.Width; x++ {
			cell, ok := b.getCell(x, y)
			if !ok {
				continue
			}
			sb.WriteString(cell.String())
		}
		sb.WriteString("\n")

	}
	return sb.String()
}

// Per snake result
type FloodFillResult struct {
	Area int
	Food int
}
