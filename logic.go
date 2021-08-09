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

func (c Coord) outOfBounds(width, height int) bool {
	if c.X < 0 {
		return true
	} else if c.X > width-1 {
		return true
	}
	if c.Y < 0 {
		return true
	} else if c.Y > height-1 {
		return true
	}
	return false
}

type Scored map[string]float64

// Remap scores to 0 - 1 range
func (moves Scored) zeroToOne() Scored {
	minScore := 0.0
	maxScore := 99999999.9
	for _, score := range moves {
		minScore = min(minScore, score)
		maxScore = max(maxScore, score)
	}
	for move, score := range moves {
		moves[move] = remap(score, minScore, maxScore, 0.0, 1.0)
	}
	return moves
}

func (m Scored) best() string {
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

func (m Scored) safeMoves() []string {
	var safeMoves []string
	for move, score := range m {
		if score > 0.0 {
			safeMoves = append(safeMoves, move)
		}
	}
	return safeMoves
}

type WeightedScore struct {
	weight float64
	scored Scored
}

func combineMoves(scores []WeightedScore) Scored {
	moves := map[string]float64{
		"up":    1.0,
		"down":  1.0,
		"left":  1.0,
		"right": 1.0,
	}
	for _, ws := range scores {
		for move, score := range ws.scored {
			moves[move] = moves[move] * (score * ws.weight)
		}
	}
	return moves
}

// Don't let your Battlesnake collide with itself (tail chasing ok though)
func avoidSelf(state *GameState) Scored {
	moves := map[string]float64{
		"up":    1.0,
		"down":  1.0,
		"left":  1.0,
		"right": 1.0,
	}
	allMoves := []string{"up", "down", "left", "right"}

	for _, move := range allMoves {
		nextHeadPos := newHead(state.You.Head, move)
		for i, coord := range state.You.Body {

			// it's ok to tail chase, but not just after eating
			isTail := len(state.You.Body) == i+1
			if isTail && state.You.Health < 100 {
				continue
			}

			if samePos(nextHeadPos, coord) {
				moves[move] = 0.0
				break
			}
		}
	}
	// log.Printf("avoidSelf: up: %f down: %f left: %f right: %f", moves["up"], moves["down"], moves["left"], moves["right"])
	return moves
}

// Don't hit walls.
func avoidWalls(state *GameState) Scored {
	moves := map[string]float64{
		"up":    1.0,
		"down":  1.0,
		"left":  1.0,
		"right": 1.0,
	}
	if state.You.Head.X == 0 {
		moves["left"] = 0.0
	} else if state.You.Head.X == state.Board.Width-1 {
		moves["right"] = 0.0
	}
	if state.You.Head.Y == 0 {
		moves["down"] = 0.0
	} else if state.You.Head.Y == state.Board.Height-1 {
		moves["up"] = 0.0
	}

	return moves
}

func avoidOthers(state *GameState, moves Scored) Scored {
	// Don't collide with others.
	for _, move := range moves.safeMoves() {
		nextHeadPos := newHead(state.You.Head, move)
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

				if samePos(nextHeadPos, coord) {
					moves[move] = 0.0
					break
				}
			}
		}
	}
	return moves
}

// Score moves based on exciting head2head possibilities
func h2h(state *GameState, moves Scored) Scored {
	// start at 0.5 for all options, as far as h2h goes, this is boring and safe
	for move, score := range moves {
		moves[move] = score * 0.5
	}
	// Avoid head to head
	allMoves := []string{"up", "down", "left", "right"}
	for _, move := range allMoves {
		if moves[move] <= 0.0 {
			continue
		}
		nextHeadPos := newHead(state.You.Head, move)
		for _, other := range state.Board.Snakes {
			if state.You.ID == other.ID {
				continue
			}
			// avoid head to head
			// TODO: account for which snake is longer
			for _, otherMove := range allMoves {
				otherHeadPos := newHead(other.Head, otherMove)
				if samePos(nextHeadPos, otherHeadPos) {
					if state.You.Length > other.Length {
						// this could be bad because the other snake doesn't have
						// to move here, putting us in a trapped position
						// but as far as h2h scores go, it's the best
						moves[move] = 1.0
						// what about 3-way possible collision?
					} else if state.You.Length == other.Length {
						// not great since you both die... but still very
						// exiting
						moves[move] = 0.01
					} else {
						// very bad, run away
						moves[move] = 0.0
					}
				}
			}
		}
	}
	return moves
}

// Find food.
func foooood(state *GameState, moves Scored) Scored {
	// TODO: tune based on hunger and strategy
	safeMoves := moves.safeMoves()

	maxDistance := distance(Coord{X: 0, Y: 0}, Coord{X: state.Board.Width - 1, Y: state.Board.Height - 1})
	for _, move := range safeMoves {
		pos := newHead(state.You.Head, move)
		for _, food := range state.Board.Food {
			if samePos(pos, food) {
				moves[move] += 0.1
				// log.Printf("food: increased %s weight by: %f", move, 1.0)
				break
			}
			d1 := distance(state.You.Head, food)
			d2 := distance(pos, food)
			if d2 < d1 {
				amount := 0.01 * float64(maxDistance-d2)
				moves[move] += amount
				// log.Printf("food: increased %s weight by: %f", move, amount)
			}
		}
	}

	return moves.zeroToOne()
}

func gimmeSomeSpace(state *GameState, moves Scored) Scored {
	// Seek out larger spaces
	// From the head of each snake, do a breadth first search of possible moves
	grid := NewGrid(state.Board.Width, state.Board.Height)
	safeMoves := moves.safeMoves()
	if len(safeMoves) > 1 {
		for _, move := range safeMoves {
			// log.Printf("checking area for move: %s", move)
			area := grid.Area(state, move)
			amount := 0.01 * float64(area)
			moves[move] += amount
			// log.Printf("area: increased %s weight by: %f", move, amount)
		}
	}

	return moves.zeroToOne()
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state GameState) BattlesnakeMoveResponse {

	// 4 possible moves
	// Each computation stage will return the four moves with a score from 0 to 1.
	// A strategy is run which outputs a weight for each stage from 0 to 1. This
	//   represents how important that stage is given the current context.
	// The final score for the moves is determined multiplying all the weighted
	//   scores for a move together (min=0, max=1).
	// The move with the heighest final score wins.

	// avoid self
	avoidSelfScore := avoidSelf(&state)
	// avoid walls
	avoidWallsScore := avoidWalls(&state)
	moves := combineMoves([]WeightedScore{{1.0, avoidSelfScore}, {1.0, avoidWallsScore}})

	// avoid others (body)
	// the result is the combination of avoid self, avoid walls, and avoid other snakes
	avoidInstantDeath := avoidOthers(&state, moves)
	log.Printf("avoidInstantDeath: %v", avoidInstantDeath)

	// head2head
	h2hScore := h2h(&state, avoidInstantDeath)

	// prefer larger areas (don't get boxed in)
	spaceScore := gimmeSomeSpace(&state, avoidInstantDeath)

	// seek food
	foodScore := foooood(&state, avoidInstantDeath)
	log.Printf("foodScore: %v", foodScore)

	foodWeight := 1.0
	// longEnough := true
	// for _, snake := range state.Board.Snakes {
	// 	if state.You.Length < snake.Length+4 {
	// 		longEnough = false
	// 		break
	// 	}
	// }
	// foodWeight = 1.0 - remap(float64(state.You.Health), 1.0, 50.0, 0.0, 1.0)
	// if !longEnough {
	// 	foodWeight = foodWeight + 1.0
	// }
	// foodWeight = remap(foodWeight, 0.0, 2.0, 0.0, 1.0)

	finalWeightedScore := combineMoves([]WeightedScore{
		{1.0, avoidInstantDeath},
		{1.0, h2hScore},
		{1.0, spaceScore},
		{foodWeight, foodScore}, // TODO: tune based on other snakes length and hunger
	})
	log.Printf("finalWeightedScore: %v", finalWeightedScore)
	nextMove := finalWeightedScore.best()

	log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}
