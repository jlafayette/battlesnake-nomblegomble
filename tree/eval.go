package tree

import (
	"fmt"
	"strings"
)

type Board struct {
	Width        int
	Height       int
	snakeCount   int
	Turn         Turn
	Cells        []*Cell
	lengths1     map[SnakeIndex]int          // original lengths
	lengths      map[SnakeIndex]int          // current lengths
	areas        map[SnakeIndex]float64      // this is a float so hazards can count less
	foodTrackers map[SnakeIndex]*foodTracker // foods in area weighted by distance
	ate          map[SnakeIndex]bool
	dead         map[SnakeIndex]bool
	health       map[SnakeIndex]int
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
	areas := make(map[SnakeIndex]float64, snakeNumber)
	dead := make(map[SnakeIndex]bool, snakeNumber)
	health := make(map[SnakeIndex]int, snakeNumber)
	trackers := make(map[SnakeIndex]*foodTracker, snakeNumber)
	for i := 0; i < snakeNumber; i++ {
		trackers[SnakeIndex(i)] = newFoodTracker()
	}

	b := &Board{
		Width:        w,
		Height:       h,
		snakeCount:   snakeNumber,
		Turn:         0,
		Cells:        cells,
		lengths:      lengths,
		lengths1:     lengths1,
		areas:        areas,
		foodTrackers: trackers,
		ate:          ate,
		dead:         dead,
		health:       health,
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
	b.snakeCount = len(snakes)
	b.Turn = 0
	// clear areas
	for i := range b.areas {
		b.areas[i] = 0
	}
	// clear ate
	for i := range b.ate {
		b.ate[i] = false
	}
	// clear food tracker
	for i := range b.foodTrackers {
		b.foodTrackers[i].reset()
	}
	// clear cells
	for _, cell := range b.Cells {
		cell.Clear()
	}
	for _, snake := range snakes {
		b.dead[SnakeIndex(snake.Index)] = snake.Dead
		if snake.Dead {
			// Lengths shouldn't ever be zero in case we need to divide by them
			// besides, 3 is the smallest a snake can be
			b.lengths[SnakeIndex(snake.Index)] = 3
			b.lengths1[SnakeIndex(snake.Index)] = 3
			b.health[SnakeIndex(snake.Index)] = 0
			continue
		}
		b.lengths[SnakeIndex(snake.Index)] = snake.Length
		b.lengths1[SnakeIndex(snake.Index)] = snake.Length
		b.health[SnakeIndex(snake.Index)] = snake.Health
		for bodyIndex, bCoord := range snake.Body {
			cell, _ := b.getCell(bCoord.X, bCoord.Y)
			cell.SetSnake(SnakeIndex(snake.Index), BodyIndex(bodyIndex), snake.Length, snake.Health, Turn(0))
		}
	}
	for _, foodCoord := range foods {
		cell, _ := b.getCell(foodCoord.X, foodCoord.Y)
		cell.SetFood()
	}
	for _, hazardCoord := range hazards {
		cell, _ := b.getCell(hazardCoord.X, hazardCoord.Y)
		cell.SetHazard()
	}

	// So contents set by SetSnake are copied to prevContents
	for _, cell := range b.Cells {
		cell.NewTurn()
	}
}

func (b *Board) _checkNeighbor(x, y int, cell *Cell, nx, ny int) bool {
	// Check for nearby heads
	nCell, ok := b.getCell(nx, ny)
	if ok && nCell.IsHead() {
		id := nCell.SnakeId()
		cell.NewHeadFrom(nCell, b.Turn)
		if cell.IsFood() {
			b.ate[id] = true
			b.foodTrackers[id].add(Coord{x, y}, b.Turn)
		}
		b.areas[id] += cell.Area()
		return true
	}
	return false
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
			head := b._checkNeighbor(x, y, cell, x-1, y)
			if head {
				done = false
			}
			head = b._checkNeighbor(x, y, cell, x+1, y)
			if head {
				done = false
			}
			head = b._checkNeighbor(x, y, cell, x, y+1)
			if head {
				done = false
			}
			head = b._checkNeighbor(x, y, cell, x, y-1)
			if head {
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
	// for snakeIndex, length := range b.lengths {
	// 	origLen := b.lengths1[snakeIndex]
	// 	food := length - origLen
	// 	results[snakeIndex].Food = int(food)
	// }
	for id := range b.foodTrackers {
		results[id].Food = b.foodTrackers[id].score
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
	Food float64
	// TODO: health cost to nearest food (acount for hazards)
}

func (e EvalResult) String() string {
	return fmt.Sprintf("area: %d, food: %.1f", e.Area, e.Food)
}

// For debugging
func PrintResults(r map[SnakeIndex]*EvalResult) {
	for k, v := range r {
		fmt.Printf("%d: %v\n", k, v)
	}
}

func (b *Board) Eval(myIndex SnakeIndex) []float64 {
	scores := make([]float64, 4)
	results := b.fill()

	for indexI := 0; indexI < 4; indexI++ {
		index := SnakeIndex(indexI)
		score := 0.0

		// are we dead?
		// -9999
		iDeadScore := 0.0
		dead, ok := b.dead[index]
		if !ok {
			scores[index] = 0.0
			continue
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
		health, ok := b.health[index]
		if !ok {
			panic("no health?!")
		}
		healthScore := float64(health) * 2
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
		if myLength <= 0 {
			panic("my length is less than or equal to zero!")
		}
		otherLongest := 0
		for i, l := range b.lengths {
			if i == index {
				continue
			}
			otherLongest = max(otherLongest, l)
		}
		longestScore := remap(float64(clamp(myLength-otherLongest, -10, 10)), -10, 10, -120, 120)
		score += longestScore

		// how much space do we have relative to other snakes?
		// mySpace - otherSmallest
		// -100,100
		// Only be concerned with the smallest of the other areas. Better to
		// trap one other snake than to reduce the average space of the others.
		myArea := 0.0
		otherSmallestArea := HIGHEST
		otherSmallestFound := false
		for k, v := range results {
			if k == index {
				myArea = minf(float64(v.Area)/float64(myLength), 3)
			} else {
				otherSmallestArea = minf(minf(float64(v.Area)/float64(b.lengths[k]), 3), otherSmallestArea)
				otherSmallestFound = true
			}
		}
		rawArea := myArea
		if otherSmallestFound {
			rawArea = myArea - otherSmallestArea
		}
		areaScore := remap(float64(rawArea), -3, 3, -100, 100)
		score += areaScore

		// Food
		// this is all calculated in the foodTracker struct for my snake
		rawFoodScore := float64(results[index].Food)
		foodScore := rawFoodScore * 2.0
		score += foodScore

		// fmt.Printf("%d score: %.1f iDead: %.1f othersDead: %.1f health: %.1f food/score: %.1f/%.1f length: %.1f area me/otherSmallest/raw/score: %.1f/%.1f/%.1f/%.1f\n", index, score, iDeadScore, othersDeadScore, healthScore, rawFoodScore, foodScore, longestScore, myArea, otherSmallestArea, rawArea, areaScore)

		scores[index] = score
	}
	if len(scores) != 4 {
		panic(fmt.Sprintf("expected 4 scores, got %d", len(scores)))
	}
	return scores
}
