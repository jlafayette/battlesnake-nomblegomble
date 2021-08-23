package main

import "github.com/jlafayette/battlesnake-go/score"

func avoidEdges(state *GameState, moves *score.Moves) {
	alreadyOnEdge := state.You.Head.onEdge(state.Board.Width, state.Board.Height)
	if alreadyOnEdge {
		return
	}

	for _, move := range moves.SafeMoves() {
		newCoord := newHead(state.You.Head, move.Str)
		if newCoord.onEdge(state.Board.Width, state.Board.Height) {
			move.ToEdge = true
		}
	}
}
