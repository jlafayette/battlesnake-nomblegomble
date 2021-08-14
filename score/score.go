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

// -- Score functions

type Scored map[string]float64

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
