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
		Color:      "#6ad7e5",
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

	// Depth of 20 is the upper limit, it will return when the timeout is
	// reached most of the time.
	treeState := tree.NewState(&state, 20)
	move, lvl := treeState.FindBestMove(true)

	move_str := move.String()
	if move_str == "dead" {
		fmt.Printf("ERROR: best move should never be dead")
		move_str = "up"
	}

	log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, move_str)
	return wire.BattlesnakeMoveResponse{
		Move: move_str, Shout: strconv.Itoa(lvl),
	}
}
