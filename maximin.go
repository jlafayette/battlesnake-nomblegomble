package main

import (
	"fmt"

	"github.com/jlafayette/battlesnake-go/score"
	"github.com/jlafayette/battlesnake-go/t"
)

func minR(a, b scoredMove) scoredMove {
	if a.score < b.score {
		fmt.Printf("%v is more min than %v\n", a, b)
		return a
	}
	fmt.Printf("%v is more min than %v\n", b, a)
	return b
}

func maxR(a, b scoredMove) scoredMove {
	if a.score > b.score {
		fmt.Printf("%v is more max than %v\n", a, b)
		return a
	}
	fmt.Printf("%v is more max than %v\n", b, a)
	return b
}

type scoredMove struct {
	score int
	move  Move
}

func (s scoredMove) String() string {
	return fmt.Sprintf("move %v score: %d", s.move, s.score)
}

func scoreResult(myIndex SnakeIndex, result map[SnakeIndex]*FloodFillResult) int {
	mySpace := 0
	myFood := 0
	othersSpace := 0
	othersFood := 0
	count := 0
	for idx, r := range result {
		count += 1
		if idx == myIndex {
			mySpace = r.Area
			myFood = r.Food
			continue
		}
		count += 1
		othersSpace += r.Area
		othersFood += r.Food
	}
	if count >= 1 {
		// TODO: Add food calulation here
		space := mySpace - (othersSpace / count)
		fmt.Printf("score: %d\n", space)
		return space
	} else {
		score := mySpace + myFood
		fmt.Printf("score: %d\n", score)
		return score
	}
}

func Maximin(state *t.GameState, moves *score.Moves) {
	board := NewBoard(state)

	var myIndex SnakeIndex
	for index, snake := range state.Board.Snakes {
		if snake.ID == state.You.ID {
			myIndex = SnakeIndex(index)
			break
		}
	}

	// maybe options could be like a tree? where each move node has parent and children
	// how do we pop moves off of the board? Need to know where the tail was and where dead snake was
	// might have to remember back X number of moves

	options := board.Moves(myIndex)
	fmt.Printf("options: %v\n", options)
	max := scoredMove{score: -999999, move: Down}
	for _, minmove := range options {
		min := scoredMove{score: 999999, move: minmove.You.Move}
		for _, moveset := range minmove.Other {
			board.Push(minmove.You, moveset)
			result := board.Fill()
			board.Restore()
			newScoredMove := scoredMove{score: scoreResult(myIndex, result), move: minmove.You.Move}
			min = minR(min, newScoredMove)
			board.Pop(minmove.You, moveset)
		}
		max = maxR(max, min)
	}

	switch max.move {
	case Up:
		moves.Up.Maximin = 1.0
	case Down:
		moves.Down.Maximin = 1.0
	case Left:
		moves.Left.Maximin = 1.0
	case Right:
		moves.Right.Maximin = 1.0
	}

	// fastest way to integrate the new space/food floodfill is to look one move ahead
	// max = -9999
	// for m in mysafe moves
	//		min = 9999
	// 		for combo in others
	//			pos = m + combo
	//			score pos
	// 			min = min(score, min)
	//		max = max(min, max)
	// do the max move

	// board needs to be able to push and pop moves onto it, then solve the position
	// push move {0: left, 1: up}
	// if at depth limit, solve
	// pop move reverse logic
}
