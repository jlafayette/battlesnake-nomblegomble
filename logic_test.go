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

func contains(moves []string, move string) bool {
	for _, m := range moves {
		if move == m {
			return true
		}
	}
	return false
}

func TestWallAvoidance(t *testing.T) {
	tests := []struct {
		input    Battlesnake
		intoNeck string
		intoWall []string
	}{
		{
			input: Battlesnake{
				// Lower left corner
				Head: Coord{X: 0, Y: 0},
				Body: []Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			},
			intoNeck: "right",
			intoWall: []string{"left", "down"},
		},
		{
			input: Battlesnake{
				// top right corner
				Head: Coord{X: 11, Y: 11},
				Body: []Coord{{X: 11, Y: 11}, {X: 10, Y: 11}, {X: 9, Y: 11}},
			},
			intoNeck: "left",
			intoWall: []string{"up", "right"},
		},
		{
			input: Battlesnake{
				// bottom right corner (facing down)
				Head: Coord{X: 11, Y: 0},
				Body: []Coord{{X: 11, Y: 0}, {X: 11, Y: 10}, {X: 11, Y: 9}},
			},
			intoNeck: "up",
			intoWall: []string{"down", "right"},
		},
		{
			input: Battlesnake{
				// top left corner (facing up)
				Head: Coord{X: 0, Y: 11},
				Body: []Coord{{X: 0, Y: 11}, {X: 0, Y: 10}, {X: 0, Y: 9}},
			},
			intoNeck: "down",
			intoWall: []string{"left", "up"},
		},
	}

	for _, tc := range tests {
		state := GameState{
			Board: Board{
				Width:  12,
				Height: 12,
				Snakes: []Battlesnake{tc.input},
			},
			You: tc.input,
		}

		nextMove := move(state)

		if nextMove.Move == tc.intoNeck {
			t.Errorf("snake moved onto its own neck, %s", nextMove.Move)
		}
		if contains(tc.intoWall, nextMove.Move) {
			t.Errorf("snake moved into a wall, %s", nextMove.Move)
		}
	}
}

func TestSelfAvoidance(t *testing.T) {
	tests := []struct {
		input    Battlesnake
		intoSelf []string
	}{
		{
			input: Battlesnake{
				Head: Coord{X: 5, Y: 5},
				Body: []Coord{{X: 5, Y: 5}, {X: 5, Y: 4}, {X: 6, Y: 4}, {X: 6, Y: 5}, {X: 6, Y: 6}, {X: 5, Y: 6}, {X: 4, Y: 6}},
			},
			intoSelf: []string{"up", "right", "down"},
		},
	}

	for _, tc := range tests {
		state := GameState{
			Board: Board{
				Width:  12,
				Height: 12,
				Snakes: []Battlesnake{tc.input},
			},
			You: tc.input,
		}

		nextMove := move(state)

		if contains(tc.intoSelf, nextMove.Move) {
			t.Errorf("snake moved into self, %s", nextMove.Move)
		}
	}
}
