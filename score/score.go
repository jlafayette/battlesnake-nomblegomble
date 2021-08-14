package score

import (
	"fmt"
	"log"
)

// -- Utility functions

func remap(old, oldMin, oldMax, newMin, newMax float64) float64 {
	oldRange := oldMax - oldMin
	if oldRange == 0 {
		return newMin
	} else {
		newRange := newMax - newMin
		return (((old - oldMin) * newRange) / oldRange) + newMin
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// -- Score functions

type H2HOutcome int

const (
	Na H2HOutcome = iota
	Lose
	Tie
	Win
)

// Details like percentage, equal, how many other options you have or the other snake has
type H2H struct {
	IsFood  bool
	Outcome H2HOutcome
	// Percentage float64 // based on other moves the snake has 0-1
}

type FoodInfo struct {
	LongEnough  bool
	Health      int
	LegacyScore float64
}

type Score struct {
	Death  bool
	Mylen  int
	Space  int
	H2h    H2H
	food   int // food in the space
	Food   FoodInfo
	result float64
	Str    string
}

type Moves struct {
	Up    *Score
	Down  *Score
	Left  *Score
	Right *Score
}

func NewMoves(length int32) *Moves {
	return &Moves{
		Up:    &Score{Mylen: int(length), Str: "up"},
		Down:  &Score{Mylen: int(length), Str: "down"},
		Left:  &Score{Mylen: int(length), Str: "left"},
		Right: &Score{Mylen: int(length), Str: "right"},
	}
}

func (m *Moves) Iter() []*Score {
	return []*Score{m.Up, m.Down, m.Left, m.Right}
}

func (m *Moves) SafeMoves() []*Score {
	safe := []*Score{}
	for _, move := range m.Iter() {
		if !move.Death {
			safe = append(safe, move)
		}
	}
	return safe
}

func (m Moves) maxSpace() int {
	return maxInt(0, maxInt(m.Up.Space, maxInt(m.Down.Space, maxInt(m.Left.Space, m.Right.Space))))
}

func (m Moves) maxLegacyFood() float64 {
	return max(0, max(m.Up.Food.LegacyScore, max(m.Down.Food.LegacyScore, max(m.Left.Food.LegacyScore, m.Right.Food.LegacyScore))))
}

func (m Moves) Choice() string {

	for _, score := range m.Iter() {
		// Death
		if score.Death {
			score.result = 0.0
		}

		// H2H
		var h2h float64
		switch score.H2h.Outcome {
		case Na:
			h2h = 0.0
		case Win:
			h2h = 1.0
		case Tie:
			h2h = 0.1
		case Lose:
			h2h = 0.01
		}
		// TODO: work in food and other choices into H2H calculation
		// log.Printf("%s h2h score: %.2f", score.Str, h2h)
		score.result += h2h

		// Space
		// TODO: make space relative to mylen
		// TODO: turn off food if space test does not pass
		space := remap(float64(score.Space), 0.0, float64(m.maxSpace()), 0.0, 1.0)
		score.result += space
		// log.Printf("%s space score: %.2f", score.Str, space)

		// Food
		// TODO: replace with calculation based on food in space
		food := remap(score.Food.LegacyScore, 0.0, m.maxLegacyFood(), 0.0, 1.0)
		foodWeight := 0.0
		if score.Food.Health < 50 {
			foodWeight = 0.5
		}
		if score.Food.Health < 25 {
			foodWeight = 0.75
		} else if score.Food.Health < 10 {
			foodWeight = 1.0
		}
		if !score.Food.LongEnough {
			foodWeight = max(foodWeight+0.25, 1.0)
		}
		foodScore := food * foodWeight * space
		// log.Printf("%s food score: %.2f", score.Str, foodScore)
		score.result += foodScore
	}

	// Pick move based on result value.
	best := 0.0
	move := m.Up
	for _, score := range m.Iter() {
		if score.result > best {
			best = score.result
			move = score
		}
	}
	log.Printf("Final score: up: %.2f, down: %.2f, left: %.2f, right %.2f", m.Up.result, m.Down.result, m.Left.result, m.Right.result)
	return move.Str
}

type Scored map[string]float64

func NewScored(value float64) Scored {
	return Scored{
		"up":    value,
		"down":  value,
		"left":  value,
		"right": value,
	}
}

func (moves Scored) IsEmpty() bool {
	for _, score := range moves {
		if score > 0.0 {
			return false
		}
	}
	return true
}

func (moves Scored) Copy() Scored {
	newMoves := Scored{
		"up":    moves["up"],
		"down":  moves["down"],
		"left":  moves["left"],
		"right": moves["right"],
	}
	return newMoves
}

// Remap scores to 0 - 1 range
func (moves Scored) ZeroToOne() Scored {
	minScore := 0.0
	maxScore := 0.0
	for _, score := range moves {
		maxScore = max(maxScore, score)
	}
	for move, score := range moves {
		moves[move] = remap(score, minScore, maxScore, 0.0, 1.0)
	}
	return moves
}

func (moves Scored) String() string {
	return fmt.Sprintf("Scored[up:%.2f down:%.2f left:%.2f right:%.2f]", moves["up"], moves["down"], moves["left"], moves["right"])
}

func (m Scored) Best() string {
	// log.Printf("Finding best move from: %v", m)
	nextMove := "none"
	var bestScore float64
	bestScore = 0.0
	for move, score := range m {
		if score > bestScore {
			nextMove = move
			bestScore = score
		}
	}
	if nextMove == "none" {
		log.Printf("No good moves found! Defaulting to 'up'. Moves: %v", m)
		nextMove = "up"
	}
	return nextMove
}

func (m Scored) SafeMoves() []string {
	var safeMoves []string
	for move, score := range m {
		if score > 0.0 {
			safeMoves = append(safeMoves, move)
		}
	}
	return safeMoves
}

type WeightedScore struct {
	negate bool
	weight float64
	scored Scored
}

func NewWeightedScore(negate bool, weight float64, scored Scored) WeightedScore {
	return WeightedScore{negate: negate, weight: weight, scored: scored}
}

func CombineMoves(scores []WeightedScore) Scored {
	moves := map[string]float64{
		"up":    1.0,
		"down":  1.0,
		"left":  1.0,
		"right": 1.0,
	}
	// Add all the bonus scores first
	for _, ws := range scores {
		for move, score := range ws.scored {
			if !ws.negate {
				// This is more additive bonus
				moves[move] = moves[move] + (score * ws.weight)
			}
		}
	}
	// Then apply the multiplier scores
	for _, ws := range scores {
		for move, score := range ws.scored {
			if ws.negate {
				moves[move] = moves[move] * (score * ws.weight)
			}
		}
	}
	return moves
}

// Something to add ability to say skip a later score if in certain conditions
// it would be confounding factor.
type RollingScore struct {
}
