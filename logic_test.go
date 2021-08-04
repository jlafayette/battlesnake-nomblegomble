package main

import (
	"testing"
)

func TestNeckAvoidance(t *testing.T) {
	tests := []struct {
		input Battlesnake
		noGo  string
	}{
		{input: Battlesnake{
			// Length 3, facing right
			Head: Coord{X: 2, Y: 0},
			Body: []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		}, noGo: "left"},
		{input: Battlesnake{
			// Length 3, facing left
			Head: Coord{X: 7, Y: 0},
			Body: []Coord{{X: 7, Y: 0}, {X: 8, Y: 0}, {X: 9, Y: 0}},
		}, noGo: "right"},
		{input: Battlesnake{
			// Length 3, facing up
			Head: Coord{X: 5, Y: 10},
			Body: []Coord{{X: 5, Y: 9}, {X: 5, Y: 8}, {X: 5, Y: 7}},
		}, noGo: "down"},
		{input: Battlesnake{
			// Length 3, facing down
			Head: Coord{X: 5, Y: 1},
			Body: []Coord{{X: 5, Y: 2}, {X: 5, Y: 3}, {X: 5, Y: 4}},
		}, noGo: "up"},
	}

	for _, tc := range tests {
		state := GameState{
			Board: Board{
				Snakes: []Battlesnake{tc.input},
			},
			You: tc.input,
		}

		nextMove := move(state)

		if nextMove.Move == tc.noGo {
			t.Errorf("snake moved onto its own neck, %s", nextMove.Move)
		}
	}
}
