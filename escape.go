package main

import (
	"github.com/jlafayette/battlesnake-go/score"
	"github.com/jlafayette/battlesnake-go/t"
)

// Define enum
// Moving away, moving diagonally towards (row), moving diagonally towards (column), moving directly towards
// Get last move dir
// get candiate dir
// then we can reason about if it's ok to turn
type Direction int

const (
	AwayV Direction = iota
	AwayH
	TowardsV
	TowardsH
)

type EscapeDir struct {
	Dir    Direction
	Offset int
}

func (e *EscapeDir) MovingAway() bool {
	return e.Dir == AwayV || e.Dir == AwayH
}

func (e *EscapeDir) MovingTowards() bool {
	return e.Dir == TowardsV || e.Dir == TowardsH
}

func getEscapeDir(c1, c2, goal t.Coord) EscapeDir {
	d1 := c1.Distance(goal)
	d2 := c2.Distance(goal)
	away := d2 > d1
	horizontal := c1.X == c2.X
	if horizontal {
		offset := Abs(c2.Y - goal.Y)
		if away {
			return EscapeDir{Dir: AwayH, Offset: offset}
		} else {
			return EscapeDir{Dir: TowardsH, Offset: offset}
		}
	} else {
		offset := Abs(c2.X - goal.X)
		if away {
			return EscapeDir{Dir: AwayV, Offset: offset}
		} else {
			return EscapeDir{Dir: TowardsV, Offset: offset}
		}
	}
}

func Escape(state *t.GameState, moves *score.Moves) {
	if !moves.Trapped {
		return
	}
	// Find escape route (included in move)
	// Flood fill from head, find body position closest to tail
	for _, move := range moves.Iter() {
		if !move.Space.Trapped {
			move.Space.EscapeScore = 1.0 // ? not sure if this would ever happen ?
			continue
		}
		goal := t.Coord{move.Space.TargetX, move.Space.TargetY}
		neck := state.You.Body[1]
		head := state.You.Head
		head2 := state.You.Head.Moved(move.Str)

		// log.Printf("move: %s goal: %v", move.Str, goal)

		prevDir := getEscapeDir(neck, head, goal)
		newDir := getEscapeDir(head, head2, goal)

		// Is this needed?
		// Calculate available space (remove dead ends)
		// This just requires 2x2 chunks

		// If you are moving away from the goal,
		// keep going until you must turn (dead end or corner)

		// NOTE: Dead ends are detected by lack of space, the escape plan just
		//       scores everything and relies on the other scores to shoot down
		//       an obvious bad move. Hopefully that works well enough...

		// If you have a choice prefer moving away from the goal
		if newDir.MovingAway() {
			move.Space.EscapeScore = 1.0
			continue
		}

		// If going directly towards the goal (same row/column)
		// try to turn as soon as possible
		if newDir.MovingTowards() && newDir.Offset == 0 {
			move.Space.EscapeScore = 0.1
			continue
		}

		// If moving towards the goal (diagonally), it's fine to turn if the number
		// of space in front of you is even, keep going if they are odd.
		sameDir := prevDir.Dir == newDir.Dir
		if newDir.MovingTowards() && sameDir {
			move.Space.EscapeScore = 0.75
			continue
		}

		// default...
		move.Space.EscapeScore = 0.35

		// Not sure if we need this part
		// If the space in front of you is 1 and not a dead end, go into it
		// If the space in front of you is 2 or even, then it's fine to turn
	}

}
