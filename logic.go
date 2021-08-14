package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"

	"github.com/jlafayette/battlesnake-go/score"
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

// Don't let your Battlesnake collide with itself (tail chasing ok though)
func avoidSelf(state *GameState, moves *score.Moves) {
	// allMoves := []string{"up", "down", "left", "right"}

	for _, move := range moves.Iter() {
		nextHeadPos := newHead(state.You.Head, move.Str)
		for i, coord := range state.You.Body {

			// it's ok to tail chase, but not just after eating
			isTail := len(state.You.Body) == i+1
			if isTail && state.You.Health < 100 {
				continue
			}

			if samePos(nextHeadPos, coord) {
				move.Death = true
				break
			}
		}
	}
	// log.Printf("avoidSelf: up: %f down: %f left: %f right: %f", moves["up"], moves["down"], moves["left"], moves["right"])
	// return moves
}

// Don't hit walls.
func avoidWalls(state *GameState, moves *score.Moves) {
	if state.You.Head.X == 0 {
		moves.Left.Death = true
	} else if state.You.Head.X == state.Board.Width-1 {
		moves.Right.Death = true
	}
	if state.You.Head.Y == 0 {
		moves.Down.Death = true
	} else if state.You.Head.Y == state.Board.Height-1 {
		moves.Up.Death = true
	}
}

func avoidOthers(state *GameState, moves *score.Moves) {
	// moves := prevMoves.Copy()
	// Don't collide with others.
	for _, move := range moves.SafeMoves() {
		nextHeadPos := newHead(state.You.Head, move.Str)
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
					move.Death = true
					break
				}
			}
		}
	}
}

// Score moves based on exciting head2head possibilities
func h2h(state *GameState, moves *score.Moves) {

	// Avoid head to head
	allMoves := []string{"up", "down", "left", "right"}
	for _, move := range moves.SafeMoves() {
		nextHeadPos := newHead(state.You.Head, move.Str)
		for _, other := range state.Board.Snakes {
			if state.You.ID == other.ID {
				continue
			}
			// avoid head to head
			for _, otherMove := range allMoves {
				otherHeadPos := newHead(other.Head, otherMove)
				if samePos(nextHeadPos, otherHeadPos) {
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
func foooood(state *GameState, moves *score.Moves) {
	// TODO: tune based on hunger and strategy

	// Need food in <Health> moves
	// Find closest food
	// If closest food route is in space score highest space - <var>
	// highest := 0.0
	// for _, score := range scoresSoFar {
	// 	highest = max(highest, score)
	// }
	// threshhold := highest * 0.5 // if it's higher than this, go for it

	// how do we not cancel ourselves out here?

	// sort food by distance and direction and other snakes near it?
	// sort.Slice(state.Board.Food, func(i, j int) bool {
	// 	di := distance(state.You.Head, state.Board.Food[i])
	// 	dj := distance(state.You.Head, state.Board.Food[j])
	// 	return di < dj
	// })

	for _, move := range moves.SafeMoves() {
		pos := newHead(state.You.Head, move.Str)
		for _, food := range state.Board.Food {
			if samePos(pos, food) {
				move.Food.LegacyScore += 1.0
				// log.Printf("food: increased %s weight by: %f", move, 1.0)
				break
			}
			d1 := distance(state.You.Head, food)
			d2 := distance(pos, food)

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

func gimmeSomeSpace(state *GameState, moves *score.Moves) {
	// Seek out larger spaces
	// From the head of each snake, do a breadth first search of possible moves
	grid := NewGrid(state)

	// area of snake len is ok
	// anything less is bad news
	for _, move := range moves.SafeMoves() {
		// log.Printf("checking area for move: %s", move)
		area := grid.Area(state, move.Str)
		// log.Printf("move: %s, area: %d", move, area)

		// area 1 len 10  0.1
		// area 11 len 10 1.1
		move.Space = area
		// log.Printf("area: set %s weight to: %.2f", move, amount)
	}
	// log.Printf("m          : %v", moves)
	// m := moves.ZeroToOne()
	// log.Print("m zeroToOne: ", m)
	// return m
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

	moves := score.NewMoves(state.You.Length)

	// avoid self
	avoidSelf(&state, moves)
	// log.Print("avoidSelfScore: ", avoidSelfScore)
	// avoid walls
	avoidWalls(&state, moves)
	// log.Printf("avoidWallsScore: %v", avoidWallsScore)
	// moves := score.CombineMoves([]score.WeightedScore{
	// 	score.NewWeightedScore(true, 1.0, avoidSelfScore),
	// 	score.NewWeightedScore(true, 1.0, avoidWallsScore),
	// })
	// log.Print("combinedSelfAndWall: ", moves)

	// avoid others (body)
	// the result is the combination of avoid self, avoid walls, and avoid other snakes
	avoidOthers(&state, moves)
	// log.Print("avoidInstantDeath: ", avoidInstantDeath)

	// head2head
	h2h(&state, moves)
	// log.Print("h2hScore: ", h2hScore)

	// prefer larger areas (don't get boxed in)
	gimmeSomeSpace(&state, moves)
	// log.Print("spaceScore: ", spaceScore)

	// seek food
	// soFarWeightedScores := []score.WeightedScore{
	// 	score.NewWeightedScore(true, 1.0, avoidInstantDeath),
	// 	score.NewWeightedScore(true, 1.0, h2hScore),
	// }
	// if !spaceScore.IsEmpty() {
	// 	soFarWeightedScores = append(soFarWeightedScores, score.NewWeightedScore(true, 1.0, spaceScore))
	// }
	// scoreSoFar := score.CombineMoves(soFarWeightedScores)
	// log.Print("scoreSoFar: ", scoreSoFar)
	foooood(&state, moves)
	// log.Print("foodScore: ", foodScore)

	// Determine how hungry the snake is
	// foodWeight := 0.0
	// longEnough := true
	// for _, snake := range state.Board.Snakes {
	// 	if state.You.Length < snake.Length+4 {
	// 		longEnough = false
	// 		break
	// 	}
	// }
	// if state.You.Health < 50 {
	// 	foodWeight = 0.5
	// }
	// if state.You.Health < 25 {
	// 	foodWeight = 0.75
	// } else if state.You.Health < 10 {
	// 	foodWeight = 1.0
	// }
	// if !longEnough {
	// 	foodWeight = max(foodWeight+0.25, 1.0)
	// }

	// weightedScores := []score.WeightedScore{
	// 	score.NewWeightedScore(true, 1.0, avoidInstantDeath),
	// 	score.NewWeightedScore(true, 1.0, h2hScore),
	// 	// score.NewWeightedScore(false, foodWeight, foodScore),
	// }
	// if !spaceScore.IsEmpty() {
	// 	weightedScores = append(weightedScores, score.NewWeightedScore(true, 1.0, spaceScore))
	// }
	// weightedScores = append(weightedScores, score.NewWeightedScore(false, foodWeight, foodScore))
	// // log.Printf("%#v", weightedScores)

	// finalWeightedScore := score.CombineMoves(weightedScores)
	// log.Print("finalWeightedScore: ", finalWeightedScore)
	// nextMove := finalWeightedScore.Best()
	// log.Println(moves)
	nextMove := moves.Choice()

	log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}
