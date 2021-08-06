package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"
)

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "jlafayette",
		Color:      "#6ad7e5",
		Head:       "beluga",
		Tail:       "bolt",
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

func newHead(head Coord, move string) Coord {
	switch move {
	case "up":
		return Coord{head.X, head.Y + 1}
	case "down":
		return Coord{head.X, head.Y - 1}
	case "left":
		return Coord{head.X - 1, head.Y}
	case "right":
		return Coord{head.X + 1, head.Y}
	}
	panic("invalid move")
	// return Coord{head.X, head.Y} // invalid move - shouldn't happen
}

func samePos(a, b Coord) bool {
	return a.X == b.X && a.Y == b.Y
}

type Moves struct {
	All      []string
	Weighted map[string]float64
}

func NewMoves() Moves {
	return Moves{
		All: []string{"up", "down", "left", "right"},
		Weighted: map[string]float64{
			"up":    1.0,
			"down":  1.0,
			"left":  1.0,
			"right": 1.0,
		},
	}
}

func (m *Moves) best() string {
	nextMove := "up" // a default in case nothing is better
	var bestScore float64
	bestScore = 0.0
	for move, score := range m.Weighted {
		if score > bestScore {
			nextMove = move
			bestScore = score
		}
	}
	return nextMove
}

func (m *Moves) safe() []string {
	safeMoves := []string{}
	for move, score := range m.Weighted {
		if score > 0.0 {
			safeMoves = append(safeMoves, move)
		}
	}
	return safeMoves
}

func (m *Moves) avoidWall(move string) {
	log.Printf("wall avoidance: setting %s to 0.0", move)
	m.Weighted[move] = 0.0
}

func (m *Moves) avoidSelf(move string) {
	log.Printf("self avoidance: setting %s to 0.0", move)
	m.Weighted[move] = 0.0
}

func (m *Moves) avoidOther(move string) {
	log.Printf("other snake avoidance: setting %s to 0.0", move)
	m.Weighted[move] = 0.0
}

func (m *Moves) avoidHead2Head(move string) {
	m.Weighted[move] -= 0.1
	log.Printf("head2head avoidance: setting %s to %f", move, m.Weighted[move])
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state GameState) BattlesnakeMoveResponse {

	possibleMoves := NewMoves()

	// Don't let your Battlesnake move back in on it's own neck
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of body piece directly behind your head (your "neck")
	if myNeck.X < myHead.X {
		possibleMoves.avoidSelf("left")
	} else if myNeck.X > myHead.X {
		possibleMoves.avoidSelf("right")
	} else if myNeck.Y < myHead.Y {
		possibleMoves.avoidSelf("down")
	} else if myNeck.Y > myHead.Y {
		possibleMoves.avoidSelf("up")
	}

	// Don't hit walls.
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height
	if myHead.X == 0 {
		possibleMoves.avoidWall("left")
	} else if myHead.X == boardWidth-1 {
		possibleMoves.avoidWall("right")
	}
	if myHead.Y == 0 {
		possibleMoves.avoidWall("down")
	} else if myHead.Y == boardHeight-1 {
		possibleMoves.avoidWall("up")
	}

	// Don't hit yourself.
	mybody := state.You.Body
	for _, move := range possibleMoves.All {
		nextHeadPos := newHead(myHead, move)
		for i, coord := range mybody {

			// it's ok to tail chase, but not just after eating
			isTail := len(mybody) == i+1
			if isTail && state.You.Health < 100 {
				continue
			}

			if samePos(nextHeadPos, coord) {
				possibleMoves.avoidSelf(move)
				break
			}
		}
	}

	// Don't collide with others.
	for _, move := range possibleMoves.safe() {
		nextHeadPos := newHead(myHead, move)
		for _, other := range state.Board.Snakes {
			if state.You.ID == other.ID {
				continue
			}
			for _, coord := range other.Body {
				if samePos(nextHeadPos, coord) {
					possibleMoves.avoidOther(move)
					break
				}
			}
			// avoid head to head
			// TODO: account for which snake is longer
			for _, otherMove := range possibleMoves.All {
				otherHeadPos := newHead(other.Head, otherMove)
				if samePos(nextHeadPos, otherHeadPos) {
					possibleMoves.avoidHead2Head(move)
				}
			}
		}
	}

	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.

	// Finally, choose a move from the available safe moves.
	nextMove := possibleMoves.best()
	// log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}
