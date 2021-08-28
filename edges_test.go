package main

import (
	"testing"

	tt "github.com/jlafayette/battlesnake-go/t"
)

func TestEdgeAvoidanceGoneTooFar1(t *testing.T) {
	state := tt.GameState{
		Game: tt.Game{
			ID: "bf50aa4a-413a-4b33-8a42-abe37cdddeda",
			Ruleset: tt.Ruleset{
				Name:    "standard",
				Version: "v1.0.20",
			},
			Timeout: 500,
		},
		Turn: 252,
		Board: tt.Board{
			Height: 11,
			Width:  11,
			Food:   []tt.Coord{{0, 0}},
			Snakes: []tt.Battlesnake{
				{
					ID:      "gs_mxgfpBFHd4jcd4S8PSgcRBxS",
					Name:    "nomblegomble",
					Health:  90,
					Head:    tt.Coord{7, 9},
					Body:    []tt.Coord{{7, 9}, {8, 9}, {9, 9}, {9, 8}, {9, 7}, {9, 6}, {9, 5}, {9, 4}, {8, 4}, {8, 5}, {8, 6}, {7, 6}, {7, 7}, {6, 7}, {6, 8}, {6, 9}, {6, 10}, {5, 10}, {4, 10}, {3, 10}, {3, 9}, {4, 9}, {5, 9}, {5, 8}, {5, 7}, {5, 6}, {6, 6}, {6, 5}},
					Length:  28,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_RTD6hXbGw6kQw4rK7HMxGcXb",
					Name:    "Hot Soup",
					Health:  97,
					Head:    tt.Coord{3, 5},
					Body:    []tt.Coord{{3, 5}, {4, 5}, {5, 5}, {5, 4}, {4, 4}, {3, 4}, {3, 3}, {3, 2}, {3, 1}, {3, 0}, {2, 0}, {2, 1}, {1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}, {2, 6}},
					Length:  19,
					Latency: "26",
					Shout:   "",
				},
			},
		},
		You: tt.Battlesnake{
			ID:      "gs_mxgfpBFHd4jcd4S8PSgcRBxS",
			Name:    "nomblegomble",
			Health:  90,
			Head:    tt.Coord{7, 9},
			Body:    []tt.Coord{{7, 9}, {8, 9}, {9, 9}, {9, 8}, {9, 7}, {9, 6}, {9, 5}, {9, 4}, {8, 4}, {8, 5}, {8, 6}, {7, 6}, {7, 7}, {6, 7}, {6, 8}, {6, 9}, {6, 10}, {5, 10}, {4, 10}, {3, 10}, {3, 9}, {4, 9}, {5, 9}, {5, 8}, {5, 7}, {5, 6}, {6, 6}, {6, 5}},
			Length:  28,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved into too small of space, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}
