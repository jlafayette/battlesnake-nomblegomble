package main

// This file hold the main logic for the endpoints - mostly the move one.

import (
	"log"

	"github.com/jlafayette/battlesnake-go/score"
	"github.com/jlafayette/battlesnake-go/t"
)

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() t.BattlesnakeInfoResponse {
	log.Println("INFO")
	return t.BattlesnakeInfoResponse{
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
func start(state t.GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state t.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// Don't let your Battlesnake collide with itself (tail chasing ok though)
func avoidSelf(state *t.GameState, moves *score.Moves) {
	for _, move := range moves.Iter() {
		nextHeadPos := state.You.Head.Moved(move.Str)
		for i, coord := range state.You.Body {

			// it's ok to tail chase, but not just after eating
			isTail := len(state.You.Body) == i+1
			if isTail && state.You.Health < 100 {
				continue
			}

			if nextHeadPos.SamePos(coord) {
				move.Death = true
				break
			}
		}
	}
}

// Don't hit walls.
func avoidWalls(state *t.GameState, head t.Coord, moves *score.Moves) {
	if head.X == 0 {
		moves.Left.Death = true
	} else if head.X == state.Board.Width-1 {
		moves.Right.Death = true
	}
	if head.Y == 0 {
		moves.Down.Death = true
	} else if head.Y == state.Board.Height-1 {
		moves.Up.Death = true
	}
}

// Avoid moves that collide with others snakes.
func avoidOthers(state *t.GameState, moves *score.Moves) {
	for _, move := range moves.SafeMoves() {
		nextHeadPos := state.You.Head.Moved(move.Str)
		for _, other := range state.Board.Snakes {
			// if state.You.ID == other.ID {
			// 	continue
			// }
			for i, coord := range other.Body {
				// it's ok to tail chase, but not just after other snake has eaten
				isTail := len(other.Body) == i+1
				if isTail && other.Health < 100 {
					continue
				}

				if nextHeadPos.SamePos(coord) {
					move.Death = true
					break
				}
			}
		}
	}
}

// Score moves based on exciting head2head possibilities
func h2h(state *t.GameState, moves *score.Moves) {

	// Avoid head to head
	allMoves := []string{"up", "down", "left", "right"}
	for _, move := range moves.SafeMoves() {
		nextHeadPos := state.You.Head.Moved(move.Str)
		for _, other := range state.Board.Snakes {
			if state.You.ID == other.ID {
				continue
			}

			// All safe moves of other snake
			otherMoves := score.NewMoves(other.Length)
			avoidWalls(state, other.Head, otherMoves)
			avoidOthers(state, otherMoves)
			optionCount := otherMoves.SafeCount()

			// avoid head to head
			for _, otherMove := range allMoves {
				otherHeadPos := other.Head.Moved(otherMove)
				if nextHeadPos.SamePos(otherHeadPos) {
					// Check for food on this square
					for _, f := range state.Board.Food {
						if nextHeadPos.SamePos(f) {
							move.H2h.IsFood = true
						}
					}

					// If there is already a H2H with a longer snake, then keep that one. This is
					// a possible 3-way h2h and the longest other snake is the one to worry about.
					if move.H2h.ID != "" {
						if move.H2h.Len > int(other.Length) {
							continue
						}
					}

					move.H2h.ID = other.ID
					move.H2h.Len = int(other.Length)
					move.H2h.OptionCount = optionCount
					if state.You.Length > other.Length {
						// this could be bad because the other snake doesn't have
						// to move here, putting us in a trapped position
						// but as far as h2h scores go, it's the best
						move.H2h.Outcome = score.Win
						// moves[move] = 1.0
						// what about 3-way possible collision?
					} else if state.You.Length == other.Length {
						// not great since you both die... but still very
						// exiting
						move.H2h.Outcome = score.Tie
						// moves[move] = 0.1
					} else {
						// very bad, run away (but still better than the wall)
						move.H2h.Outcome = score.Lose
						// moves[move] = 0.01
					}
				}
			}
		}
	}
}

// Find food.
func foooood(state *t.GameState, moves *score.Moves) {
	for _, move := range moves.SafeMoves() {
		pos := state.You.Head.Moved(move.Str)
		for _, food := range state.Board.Food {
			if pos.SamePos(food) {
				move.Food.LegacyScore += 1.0
				// log.Printf("food: increased %s weight by: %f", move, 1.0)
				break
			}
			d1 := state.You.Head.Distance(food)
			d2 := pos.Distance(food)

			if d2 < d1 {
				amount := 0.0
				switch d2 {
				case 1:
					amount = 0.50
				case 2:
					amount = 0.25
				case 3:
					amount = 0.12
				case 4:
					amount = 0.05
				case 5:
					amount = 0.01
				default:
					amount = 0.001
				}
				move.Food.LegacyScore += amount
				// log.Printf("food: increased %s weight by: %f", move, amount)
			}
		}
	}
}

func gimmeSomeSpace(state *t.GameState, moves *score.Moves) {
	// Seek out larger spaces
	// From the head of each snake, do a breadth first search of possible moves
	// grid := NewGrid(state)

	// area of snake len is ok
	// anything less is bad news
	for _, move := range moves.SafeMoves() {
		// log.Printf("checking area for move: %s", move)
		// grid.Reset(state)  // TODO: Figure out why this doesn't reset properly.
		area := GetArea(state, move.Str)
		// log.Printf("move: %s, area: %d", move, area)

		// area 1 len 10  0.1
		// area 11 len 10 1.1
		move.Space = score.SpaceInfo{Area: area.Space, Trapped: area.Trapped, TargetX: area.Target.X, TargetY: area.Target.Y}
		// log.Printf("area: set %s weight to: %.2f", move, amount)
	}
	largestArea := 0
	for _, move := range moves.SafeMoves() {
		largestArea = maxInt(largestArea, move.Space.Area)
	}
	moves.Trapped = largestArea < len(state.You.Body)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
func move(state t.GameState) t.BattlesnakeMoveResponse {

	// 4 possible moves
	// Each computation stage will add information to the score.Moves struct.
	// This will calcuation the final move once all the data has been added to
	// it with moves.Choice()
	moves := score.NewMoves(state.You.Length)

	// avoid self
	avoidSelf(&state, moves)

	// avoid walls
	avoidWalls(&state, state.You.Head, moves)

	// avoid others (body)
	// the result is the combination of avoid self, avoid walls, and avoid other snakes
	avoidOthers(&state, moves)

	// head2head
	h2h(&state, moves)

	// prefer larger areas (don't get boxed in)
	gimmeSomeSpace(&state, moves)
	if moves.Trapped {
		Escape(&state, moves)
	}

	// seek food
	foooood(&state, moves)

	// Detect when moving to the edge of the board, this is often a bad move
	// since it cuts the space in half and removes a possible escape route.
	// Often it's better to turn one move before the edge, so if things go
	// bad you can double back.
	avoidEdges(&state, moves)

	nextMove := moves.Choice()

	log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	return t.BattlesnakeMoveResponse{
		Move: nextMove,
	}
}
