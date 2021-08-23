package score

import (
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

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
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
	IsFood      bool       // Is the tile food?
	Len         int        // Length of other snake (could open up if h2h is a win)
	OptionCount int        // number of safe moves the other snake has
	ID          string     // ID of the other snake
	Outcome     H2HOutcome // What would be the result? (Na means no h2h)
}

type FoodInfo struct {
	LongEnough  bool
	Health      int
	LegacyScore float64
}

type SpaceInfo struct {
	Area        int
	Trapped     bool
	TargetX     int
	TargetY     int
	EscapeScore float64
	Food        int
}

type Score struct {
	Death  bool
	Mylen  int
	Space  SpaceInfo
	H2h    H2H
	Food   FoodInfo
	result float64
	Str    string
	ToEdge bool
}

type Moves struct {
	Trapped bool
	Up      *Score
	Down    *Score
	Left    *Score
	Right   *Score
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

func (m *Moves) SafeCount() int {
	count := 0
	for _, move := range m.Iter() {
		if !move.Death {
			count += 1
		}
	}
	return count
}

func (m Moves) maxSpace() int {
	return maxInt(0, maxInt(m.Up.Space.Area, maxInt(m.Down.Space.Area, maxInt(m.Left.Space.Area, m.Right.Space.Area))))
}

func (m Moves) maxLegacyFood() float64 {
	return max(0, max(m.Up.Food.LegacyScore, max(m.Down.Food.LegacyScore, max(m.Left.Food.LegacyScore, m.Right.Food.LegacyScore))))
}

func (m Moves) h2hDeathCount() int {
	h2hCount := 0
	for _, score := range m.SafeMoves() {
		if score.H2h.Outcome != Na && score.H2h.Outcome != Win {
			h2hCount += 1
		}
	}
	return h2hCount
}

func (m Moves) Choice() string {

	// More than one tied or losing h2h?
	// This is useful to try and avoid food in this case.
	h2hDeathCount := m.h2hDeathCount()
	// log.Printf("h2hDeathCount: %d", h2hDeathCount)

	for _, score := range m.Iter() {
		// Death
		if score.Death {
			score.result = 0.0
			// log.Printf("%s = death", score.Str)
			continue
		}

		// H2H
		var h2h float64
		ignoreFood := false
		ignoreSpace := false
		switch score.H2h.Outcome {
		case Na:
			h2h = 0.0
		case Win:
			h2h = (1.0 / float64(score.H2h.OptionCount))
		case Tie:
			// If there are multiple h2h and one of them is food, chances are
			// the other snake will go for the food, so it's a better bet to
			// go the other way.
			if h2hDeathCount > 1 {
				if score.H2h.IsFood {
					h2h = 0.05
					ignoreFood = true // don't go for the food!
				} else {
					h2h = 0.4 // prefer the non food (mostly to overcome possible area difference)
				}
			} else {
				h2h = 0.1
			}
		case Lose:
			h2h = 0.01
			ignoreFood = true
			// Area has a bug where since you move first, a losing h2h is skipped
			// over, resulting in way too much space where is should be zero.
			ignoreSpace = true
		}
		score.result += h2h

		// Space
		space := remap(float64(score.Space.Area), 0.0, float64(score.Mylen), 0.0, 1.0)
		// If you have twice your body length to work with, extra doesn't really matter
		space = min(space, 2.0)
		if ignoreSpace {
			space = 0.0
		}
		score.result += space
		if score.Space.Area < score.Mylen {
			ignoreFood = true
		}

		// Escape plan
		escapeScore := 0.0
		if m.Trapped {
			escapeScore = score.Space.EscapeScore
			ignoreFood = true
		}
		score.result += (escapeScore * space)

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
		foodScore := food * foodWeight * min(space, 1.0)
		if ignoreFood {
			foodScore = 0.0
		}
		score.result += foodScore
		// log.Printf("%s food: %.2f, foodWeight: %.2f, foodScore: %.2f", score.Str, food, foodWeight, foodScore)

		if score.ToEdge {
			score.result -= 0.05 // small penalty to moving to the edge
		}

		log.Printf("%s scores | h2h: %.2f, area/space/escape: %d/%.2f/%.2f, food: %.2f, toEdge: %v", score.Str, h2h, score.Space.Area, space, escapeScore, foodScore, score.ToEdge)
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
