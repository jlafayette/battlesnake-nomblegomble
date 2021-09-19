package tree

import (
	"fmt"
	"testing"

	"github.com/jlafayette/battlesnake-go/wire"
)

func TestSpace1(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "e7e26081-2012-4184-b755-dcddf9d027b6",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 230,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{10, 8}, {10, 7}, {10, 5}, {2, 6}, {7, 6}, {6, 1}, {0, 8}, {3, 7}, {1, 0}, {10, 9}, {10, 2}, {10, 6}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_TJg73JVSFFj3hP3YSSXprHq4",
					Name:    "Nessegrev-Lefty",
					Health:  86,
					Head:    wire.Coord{6, 10},
					Body:    []wire.Coord{{6, 10}, {7, 10}, {8, 10}, {9, 10}, {9, 9}, {9, 8}, {8, 8}, {7, 8}, {6, 8}, {5, 8}, {4, 8}},
					Length:  11,
					Latency: "260",
					Shout:   "",
				},
				{
					ID:      "gs_GHJ4JcWR3hv4S6tJ7fWQDXgd",
					Name:    "nomblegomble",
					Health:  56,
					Head:    wire.Coord{0, 2},
					Body:    []wire.Coord{{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {4, 3}, {4, 4}, {3, 4}, {2, 4}, {1, 4}, {0, 4}, {0, 5}, {1, 5}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, {5, 4}, {6, 4}, {7, 4}, {8, 4}, {9, 4}},
					Length:  22,
					Latency: "25",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_GHJ4JcWR3hv4S6tJ7fWQDXgd",
			Name:    "nomblegomble",
			Health:  56,
			Head:    wire.Coord{0, 2},
			Body:    []wire.Coord{{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {4, 3}, {4, 4}, {3, 4}, {2, 4}, {1, 4}, {0, 4}, {0, 5}, {1, 5}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, {5, 4}, {6, 4}, {7, 4}, {8, 4}, {9, 4}},
			Length:  22,
			Latency: "25",
			Shout:   "",
		},
	}

	treeState := NewState(&state, 1)
	move := treeState.FindBestMove(true)

	if move == Up {
		t.Errorf("snake moved into too small of space, %v", move)
	}
	if move == Right {
		t.Errorf("snake moved into self, %v", move)
	}
	if move == Left {
		t.Errorf("snake moved into wall, %v", move)
	}
}

func TestFood1(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "9bc1403e-f1bb-4b7f-9828-57ab053ae291",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 160,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{1, 10}, {8, 4}, {6, 8}, {10, 8}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_kq96m7hDX7MWDbRRMgHjDg68",
					Name:    "trentren-vilu",
					Health:  93,
					Head:    wire.Coord{4, 2},
					Body:    []wire.Coord{{4, 2}, {4, 1}, {4, 0}, {3, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 2}, {2, 2}, {1, 2}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {1, 6}, {1, 5}, {2, 5}, {3, 5}, {4, 5}, {4, 4}, {4, 3}},
					Length:  22,
					Latency: "45",
					Shout:   "",
				},
				{
					ID:      "gs_6CjkCfdhJmq9QckqmYRFMfgJ",
					Name:    "nomblegomble",
					Health:  7,
					Head:    wire.Coord{7, 7},
					Body:    []wire.Coord{{7, 7}, {8, 7}, {9, 7}, {10, 7}, {10, 6}, {9, 6}},
					Length:  6,
					Latency: "26",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_6CjkCfdhJmq9QckqmYRFMfgJ",
			Name:    "nomblegomble",
			Health:  7,
			Head:    wire.Coord{7, 7},
			Body:    []wire.Coord{{7, 7}, {8, 7}, {9, 7}, {10, 7}, {10, 6}, {9, 6}},
			Length:  6,
			Latency: "26",
			Shout:   "",
		},
	}

	treeState := NewState(&state, 2)
	move := treeState.FindBestMove(true)
	fmt.Printf("got move: %v\n", move)

	if move == Down {
		t.Errorf("snake moved away from food, %v", move)
	}
	if move == Right {
		t.Errorf("snake moved into self, %v", move)
	}
}
