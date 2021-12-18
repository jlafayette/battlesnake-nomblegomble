package main

// This file hold the main logic for the endpoints - mostly the move one.

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jlafayette/battlesnake-go/tree"
	"github.com/jlafayette/battlesnake-go/wire"
)

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() wire.BattlesnakeInfoResponse {
	log.Println("INFO")
	return wire.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "jlafayette",
		Color:      "#03d3fc",
		Head:       "beluga",
		Tail:       "bolt",
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state wire.GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state wire.GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
func move(state wire.GameState) wire.BattlesnakeMoveResponse {

	// Depth of 3 is the upper limit. This is for the public version. In the
	// Fall 2021 and Winter 2021 tournaments this was set to 20, but most of
	// the time it returned because of running out of time at a depth of 5 or 6.
	depth := 3
	treeState := tree.NewState(&state, depth)
	move, lvl := treeState.FindBestMove(true)

	move_str := move.String()
	if move == tree.Dead || move == tree.NoMove {
		fmt.Printf("ERROR: best move should never be 'dead' or 'nomove'")
		move_str = "up"
	}

	log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, move_str)
	return wire.BattlesnakeMoveResponse{
		Move: move_str, Shout: strconv.Itoa(lvl),
	}
}
