package main

// This file hold the main logic for the endpoints - mostly the move one.

import (
	"fmt"
	"log"

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

	// At 4 moves, things are taking waaaay too long sometimes:
	// === RUN   TestAvoidBadHead2Head
	// Ran 15,926,417 loops
	// 2021/09/18 22:06:14 3509d89e-8809-46c9-b46c-164158eaac26 MOVE 3: right
	// --- PASS: TestAvoidBadHead2Head (217.13s)
	// And a later one timed out the test at 10 minutes

	// depth := 2
	// if len(state.Board.Snakes) <= 2 {
	// 	depth = 3
	// }

	treeState := tree.NewState(&state, 20)
	move := treeState.FindBestMove(true)

	move_str := move.String()
	if move_str == "dead" {
		fmt.Printf("ERROR: best move should never be dead")
		move_str = "up"
	}

	log.Printf("%s MOVE %d: %s/%s\n", state.Game.ID, state.Turn, move.String(), move_str)
	return wire.BattlesnakeMoveResponse{
		Move: move_str,
	}
}
