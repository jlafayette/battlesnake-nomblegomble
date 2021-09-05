package main

import (
	"testing"

	tt "github.com/jlafayette/battlesnake-go/t"
)

func Test__NAME__(t *testing.T) {
	state := tt.GameState{
		Game: tt.Game{
			ID: "510f71b1-d22a-49e5-a65c-81ee04751204",
			Ruleset: tt.Ruleset{
				Name:    "standard",
				Version: "v1.0.20",
			},
			Timeout: 500,
		},
		Turn: 143,
		Board: tt.Board{
			Height: 11,
			Width:  11,
			Food:   []tt.Coord{{0, 9}, {4, 2}, {2, 0}, {7, 10}, {0, 10}, {5, 5}},
			Snakes: []tt.Battlesnake{
				{
					ID:      "gs_qVYrT9XMmbptkYkjFKH6JFYT",
					Name:    "nomblegomble",
					Health:  61,
					Head:    tt.Coord{1, 8},
					Body:    []tt.Coord{{1, 8}, {1, 7}, {1, 6}, {1, 5}, {1, 4}, {1, 3}, {1, 2}, {0, 2}, {0, 1}, {0, 0}, {1, 0}, {1, 1}, {2, 1}, {3, 1}, {3, 2}, {2, 2}, {2, 3}},
					Length:  17,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_kKbKFgxtRSHGG3CYcddBVvxC",
					Name:    "Mmm, tasty",
					Health:  81,
					Head:    tt.Coord{10, 9},
					Body:    []tt.Coord{{10, 9}, {10, 8}, {10, 7}, {10, 6}, {10, 5}, {10, 4}},
					Length:  6,
					Latency: "70",
					Shout:   "",
				},
				{
					ID:      "gs_pgTPv74Pc6x8Bb4cJTgSRhd8",
					Name:    "Shai-Hulud",
					Health:  92,
					Head:    tt.Coord{8, 7},
					Body:    []tt.Coord{{8, 7}, {8, 6}, {8, 5}, {9, 5}, {9, 4}, {9, 3}, {8, 3}, {7, 3}, {7, 2}, {6, 2}, {6, 3}, {6, 4}, {6, 5}, {7, 5}, {7, 6}, {6, 6}, {6, 7}},
					Length:  17,
					Latency: "23",
					Shout:   "Goin ' up",
				},
			},
		},
		You: tt.Battlesnake{
			ID:      "gs_qVYrT9XMmbptkYkjFKH6JFYT",
			Name:    "nomblegomble",
			Health:  61,
			Head:    tt.Coord{1, 8},
			Body:    []tt.Coord{{1, 8}, {1, 7}, {1, 6}, {1, 5}, {1, 4}, {1, 3}, {1, 2}, {0, 2}, {0, 1}, {0, 0}, {1, 0}, {1, 1}, {2, 1}, {3, 1}, {3, 2}, {2, 2}, {2, 3}},
			Length:  17,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake did not fight for extra space, %s", nextMove.Move)
	}
	if nextMove.Move == "down" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake did not fight for extra space, %s", nextMove.Move)
	}
}
