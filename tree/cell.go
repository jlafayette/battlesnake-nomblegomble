package tree

import "strconv"

type Turn int
type SnakeIndex int
type BodyIndex int

type Contents int

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
	length    int        // for heads
	snakeId   SnakeIndex // index of current snake
	bodyIndex BodyIndex  // 0 is head len(body)-1 is the tail
	contents  Contents   // head,body,tail, empty, food

	prevContents Contents
}

func NewCell() *Cell {
	return &Cell{
		length:       0,
		snakeId:      -1,
		bodyIndex:    -1,
		contents:     Empty,
		prevContents: Empty,
	}
}

func (c *Cell) Clear() {
	c.length = 0
	c.snakeId = -1
	c.bodyIndex = -1
	c.contents = Empty
	c.prevContents = Empty
}

func (c *Cell) Blocked() bool {
	return c.prevContents == Head || c.prevContents == Body || c.prevContents == DoubleTail
}

func (c *Cell) IsHead() bool {
	return c.prevContents == Head
}

func (c *Cell) IsTail() bool {
	return c.prevContents == Tail
}

func (c *Cell) IsFood() bool {
	return c.prevContents == Food
}

func (c *Cell) IsSnake() bool {
	return c.prevContents == Head || c.prevContents == Body || c.prevContents == Tail || c.prevContents == DoubleTail
}

func (c *Cell) UpdateSnake(ate bool, length int) {
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
		// New head updates happen already before updating the rest of the snake
		// So check if the contents have been switch to head
		if c.contents == Head {
			c.bodyIndex = 0
		} else {
			c.bodyIndex = -1
			c.contents = Empty
		}
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

func (c *Cell) SetSnake(snakeIndex SnakeIndex, bodyIndex BodyIndex, length int, turn Turn) {
	isHead := bodyIndex == 0
	isTail := int(bodyIndex) == length-1
	if isHead {
		c.contents = Head
		c.length = length
	} else if isTail {
		c.contents = Tail
	} else {
		c.contents = Body
	}
	// TODO: detect double tail properly

	c.bodyIndex = bodyIndex
	c.snakeId = snakeIndex
}

func (c *Cell) SetFood() {
	c.contents = Food
}

func (c *Cell) NewHeadFrom(nc *Cell, turn Turn) {
	if c.contents == Head {
		// resolve H2H
		if c.snakeId == nc.snakeId {
			// ignore h2h if it's the same snake, the first one can take the cell
			return
		}
		// if they have the same length, they both die
		if c.length == nc.length {
			c.contents = H2H
			c.bodyIndex = -1
			c.snakeId = -1
		} else if c.length > nc.length {
			// the head that was marked here earlier won, so prev info can stay
		} else {
			// new nc head wins
			c.contents = Head
			c.length = nc.length
			c.snakeId = nc.snakeId
			c.bodyIndex = 0
			if c.prevContents == Food {
				c.length += 1
			}
		}
		return
	}

	// Default (no h2h)
	c.contents = Head
	c.length = nc.length
	if c.prevContents == Food {
		c.length += 1
	}
	c.snakeId = nc.snakeId
	c.bodyIndex = 0
}

func (c *Cell) NewTurn() {
	c.prevContents = c.contents
	if c.prevContents != Head {
		// Non-heads don't need length, so zero it out to help with debugging if things get screwed up
		c.length = 0
	}
}

func (c *Cell) String() string {
	switch c.prevContents {
	case Food:
		return "  fff"
	case Head:
		return "  " + strconv.Itoa(int(c.snakeId)) + "H" + strconv.Itoa(int(c.bodyIndex))
	case Body:
		return "  " + strconv.Itoa(int(c.snakeId)) + "B" + strconv.Itoa(int(c.bodyIndex))
	case Tail:
		return "  " + strconv.Itoa(int(c.snakeId)) + "T" + strconv.Itoa(int(c.bodyIndex))
	case DoubleTail:
		return "  " + strconv.Itoa(int(c.snakeId)) + "D" + strconv.Itoa(int(c.bodyIndex))
	case H2H:
		return "  xXx"
	case Empty:
		return "  [ ]"
	}
	return "     "
}
