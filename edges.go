package main

import (
	"github.com/jlafayette/battlesnake-go/score"
	"github.com/jlafayette/battlesnake-go/t"
)

func avoidEdges(state *t.GameState, moves *score.Moves) {
	alreadyOnEdge := state.You.Head.OnEdge(state.Board.Width, state.Board.Height)
	if alreadyOnEdge {
		return
	}

	for _, move := range moves.SafeMoves() {
		newCoord := state.You.Head.Moved(move.Str)
		if newCoord.OnEdge(state.Board.Width, state.Board.Height) {
			move.ToEdge = true
		}
	}
}
