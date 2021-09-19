package tree

import (
	"fmt"
	"strings"
)

type Board struct {
	Width    int
	Height   int
	Turn     Turn
	Cells    []*Cell
	lengths1 map[SnakeIndex]int // original lengths
	lengths  map[SnakeIndex]int // current lengths
	areas    map[SnakeIndex]int
	ate      map[SnakeIndex]bool
	dead     map[SnakeIndex]bool
	health   map[SnakeIndex]int
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

func NewBoard(width, height int, snakes []*Snake, foods, hazards []Coord) *Board {
	w := width
	h := height
	snakeNumber := len(snakes)
	cells := make([]*Cell, 0, w*h)
	lengths := make(map[SnakeIndex]int, snakeNumber)
	lengths1 := make(map[SnakeIndex]int, snakeNumber)
	ate := make(map[SnakeIndex]bool, snakeNumber)
	areas := make(map[SnakeIndex]int, snakeNumber)
	dead := make(map[SnakeIndex]bool, snakeNumber)
	health := make(map[SnakeIndex]int, snakeNumber)
	b := &Board{
		Width:    w,
		Height:   h,
		Turn:     0,
		Cells:    cells,
		lengths:  lengths,
		lengths1: lengths1,
		areas:    areas,
		ate:      ate,
		dead:     dead,
		health:   health,
	}
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			b.Cells = append(b.Cells, NewCell())
		}
	}
	b.Load(snakes, foods, hazards)
	return b
}

// Load a position to evaluate
func (b *Board) Load(snakes []*Snake, foods, hazards []Coord) {
	b.Turn = 0
	// clear areas
	for i := range b.areas {
		b.areas[i] = 0
	}
	// clear ate
	for i := range b.ate {
		b.ate[i] = false
	}
	// clear cells
	for _, cell := range b.Cells {
		cell.Clear()
	}
	for _, snake := range snakes {
		b.dead[SnakeIndex(snake.Index)] = snake.Dead
		if snake.Dead {
			b.lengths[SnakeIndex(snake.Index)] = 0
			b.lengths1[SnakeIndex(snake.Index)] = 0
			b.health[SnakeIndex(snake.Index)] = 0
			continue
		}
		b.lengths[SnakeIndex(snake.Index)] = snake.Length
		b.lengths1[SnakeIndex(snake.Index)] = snake.Length
		b.health[SnakeIndex(snake.Index)] = snake.Health
		for bodyIndex, bCoord := range snake.Body {
			cell, _ := b.getCell(bCoord.X, bCoord.Y)
			cell.SetSnake(SnakeIndex(snake.Index), BodyIndex(bodyIndex), snake.Length, Turn(0))
		}
	}
	for _, foodCoord := range foods {
		cell, _ := b.getCell(foodCoord.X, foodCoord.Y)
		cell.SetFood()
	}
	// TODO: Add hazards

	// So contents set by SetSnake are copied to prevContents
	for _, cell := range b.Cells {
		cell.NewTurn()
	}
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

func (b *Board) fill() map[SnakeIndex]*EvalResult {
	results := make(map[SnakeIndex]*EvalResult)
	var longest int = 0
	for snakeIndex, l := range b.lengths {
		results[snakeIndex] = &EvalResult{}
		if l > longest {
			longest = l
		}
	}
	// fmt.Printf("longest: %d\n", longest)

	// Loop until end condition is set
	// TODO: allow revisiting the same square over and over
	//       each time it should count for an extra space (like going in a
	//       circle forever is fine)
	//       but... we need to find an alternate way of terminating the loop
	//       option1: check every x turns if any progress is being made/enough area is stacked out
	//       option2: check every update for some terminating condition
	done := false
	for !done {
		// fmt.Println(b.String())
		done = b.Update()
		if (int(b.Turn) >= b.Width*b.Height) || (int(b.Turn) >= longest*2) {
			// if int(b.Turn) >= b.Width*b.Height {
			done = true
		}
	}
	// fmt.Printf("Ended updates at turn: %d\n", b.Turn)
	// fmt.Println(b.String())

	// Collect the results
	for id, area := range b.areas {
		results[id].Area = int(area)
	}
	for snakeIndex, length := range b.lengths {
		origLen := b.lengths1[snakeIndex]
		food := length - origLen
		results[snakeIndex].Food = int(food)
	}
	// for k, v := range results {
	// 	fmt.Printf("Snake %d Area: %d\n", k, v.Area)
	// 	fmt.Printf("Snake %d Food: %d\n", k, v.Food)
	// }
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
type EvalResult struct {
	Area int
	Food int
	// TODO: health cost to nearest food (acount for hazards)
}

func (e EvalResult) String() string {
	return fmt.Sprintf("area: %d, food: %d", e.Area, e.Food)
}

// For debugging
func PrintResults(r map[SnakeIndex]*EvalResult) {
	for k, v := range r {
		fmt.Printf("%d: %v\n", k, v)
	}
}

func (b *Board) Eval(index SnakeIndex) float64 {
	score := 0.0

	// are we dead?
	// -9999
	iDeadScore := 0.0
	dead, ok := b.dead[index]
	if !ok {
		panic("given index is missing from dead")
	}
	if dead {
		iDeadScore = -9999.0
	}
	score += iDeadScore

	// how many other snakes are still alive?
	// 0,50
	snakeCount := 0
	aliveCount := 0
	for k, v := range b.dead {
		if k == index {
			continue
		}
		snakeCount += 1
		if !v {
			aliveCount += 1
		}
	}
	othersDeadScore := remap(float64(aliveCount), 0, float64(snakeCount), 0, 50)
	score += othersDeadScore

	// what's our health relative to other snakes?
	// (ok if there is nearby food)
	// -200,10
	health, ok := b.health[index]
	if !ok {
		panic("no health?!")
	}
	// TODO: account for distance to nearest food
	healthScore := remap(float64(100-health), 0, 100, -200, 10)
	score += healthScore

	// what's our length relative to other snakes?
	// 0, 50
	// my length
	// other longest
	// my length - longest
	// -20?, 20? --> we don't really care if it's too much positive, so we clamp it
	// -10, 10 -> -50, 50
	myLength, ok := b.lengths[index]
	if !ok {
		panic("no length!")
	}
	otherLongest := 0
	for i, l := range b.lengths {
		if i == index {
			continue
		}
		otherLongest = max(otherLongest, l)
	}
	longestScore := remap(float64(clamp(myLength-otherLongest, -10, 10)), -10, 10, -50, 50)
	score += longestScore

	// how much space do we have relative to other snakes?
	// mySpace - (otherSpace / number of opponents)
	// -100,100
	results := b.fill()
	myArea := 0
	// myFood := 0
	othersArea := 0
	othersFood := 0
	for k, v := range results {
		if k == index {
			myArea = v.Area
			// myFood = v.Food
		} else {
			othersArea += v.Area
			othersFood += v.Food
		}
	}
	rawArea := myArea
	if aliveCount > 0 {
		rawArea = myArea - (othersArea / aliveCount)
	}
	areaClamped := clamp(rawArea, 0, myLength*3)
	areaScore := remap(float64(areaClamped), 0, float64(myLength*3), -100, 100)
	score += areaScore

	// fmt.Printf("score: %.1f iDead: %.1f othersDead: %.1f health: %.1f length: %.1f area: %d/%.1f\n", score, iDeadScore, othersDeadScore, healthScore, longestScore, rawArea, areaScore)

	return score
}
