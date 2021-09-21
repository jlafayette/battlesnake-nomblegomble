package tree

import (
	"fmt"
	"testing"

	"github.com/jlafayette/battlesnake-go/wire"
)

func TestSpace01(t *testing.T) {
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

func TestSpace02(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "19cd79a4-e053-4045-8eeb-388de9581ef4",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 123,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{6, 10}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_WKyKwcjyQdRj8prJxrmWYDH8",
					Name:    "Voxel",
					Health:  98,
					Head:    wire.Coord{6, 9},
					Body:    []wire.Coord{{6, 9}, {6, 8}, {6, 7}, {7, 7}, {7, 6}, {7, 5}, {7, 4}, {6, 4}, {6, 3}, {5, 3}, {5, 4}, {4, 4}, {4, 5}, {3, 5}, {3, 6}, {3, 7}, {3, 8}, {4, 8}},
					Length:  18,
					Latency: "72",
					Shout:   "",
				},
				{
					ID:      "gs_CXjSk4VBPdF6wRFvbcBPfQQd",
					Name:    "nomblegomble",
					Health:  77,
					Head:    wire.Coord{0, 9},
					Body:    []wire.Coord{{0, 9}, {0, 10}, {1, 10}, {2, 10}, {3, 10}, {3, 9}, {2, 9}, {2, 8}, {1, 8}, {1, 7}, {0, 7}},
					Length:  11,
					Latency: "28",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_CXjSk4VBPdF6wRFvbcBPfQQd",
			Name:    "nomblegomble",
			Health:  77,
			Head:    wire.Coord{0, 9},
			Body:    []wire.Coord{{0, 9}, {0, 10}, {1, 10}, {2, 10}, {3, 10}, {3, 9}, {2, 9}, {2, 8}, {1, 8}, {1, 7}, {0, 7}},
			Length:  11,
			Latency: "28",
			Shout:   "",
		},
	}

	// [0: up 1: down 2: dead 3: dead]
	// [0: left 1: left 2: dead 3: dead]
	// [0: left 1: right 2: dead 3: dead]
	// []
	// score: 92.0 iDead: 0.0 othersDead: 50.0 health: 77.0 food: 0.0/0.0 length: -35.0 area me/others/raw/score: 3.0/3.0/0.0/0.0

	// This was tricky because 1 and 2 depth passed, but 3 failed
	treeState := NewState(&state, 3)
	move := treeState.FindBestMove(true)

	if move == Right {
		t.Errorf("snake moved into too small of space, %v", move)
	}
	if move == Up {
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

	treeState := NewState(&state, 1)
	move := treeState.FindBestMove(true)
	fmt.Printf("got move: %v\n", move)

	if move == Down {
		t.Errorf("snake moved away from food, %v", move)
	}
	if move == Right {
		t.Errorf("snake moved into self, %v", move)
	}
}

func TestTailOk01(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "be32fe6c-b22d-4348-a01c-ab079df2a83a",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 328,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{1, 1}, {7, 1}, {2, 0}, {10, 0}, {7, 0}, {7, 9}, {3, 9}, {0, 1}, {8, 5}, {0, 9}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_rJQW9pJtHPKkR3F3ScGPMQc7",
					Name:    "nomblegomble",
					Health:  98,
					Head:    wire.Coord{9, 7},
					Body:    []wire.Coord{{9, 7}, {8, 7}, {7, 7}, {6, 7}, {5, 7}, {4, 7}, {3, 7}, {2, 7}, {1, 7}, {0, 7}, {0, 6}, {0, 5}, {1, 5}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, {6, 5}, {6, 4}, {7, 4}, {8, 4}, {9, 4}, {9, 3}, {8, 3}, {7, 3}, {6, 3}, {5, 3}, {5, 4}, {4, 4}, {3, 4}, {2, 4}, {2, 3}, {1, 3}, {0, 3}, {0, 2}, {1, 2}, {2, 2}},
					Length:  37,
					Latency: "21",
					Shout:   "",
				},
				{
					ID:      "gs_kWp8DCQxQVdFvPgVPBtJ6rxK",
					Name:    "Nessegrev-Lefty",
					Health:  95,
					Head:    wire.Coord{9, 9},
					Body:    []wire.Coord{{9, 9}, {9, 8}, {8, 8}, {7, 8}, {6, 8}, {5, 8}, {4, 8}, {3, 8}, {2, 8}, {1, 8}, {1, 9}, {1, 10}, {2, 10}, {3, 10}, {4, 10}, {5, 10}, {6, 10}, {7, 10}, {8, 10}, {9, 10}, {10, 10}, {10, 9}, {10, 8}, {10, 7}, {10, 6}, {9, 6}},
					Length:  26,
					Latency: "61",
					Shout:   "help me obiwan you're my only hope",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_rJQW9pJtHPKkR3F3ScGPMQc7",
			Name:    "nomblegomble",
			Health:  98,
			Head:    wire.Coord{9, 7},
			Body:    []wire.Coord{{9, 7}, {8, 7}, {7, 7}, {6, 7}, {5, 7}, {4, 7}, {3, 7}, {2, 7}, {1, 7}, {0, 7}, {0, 6}, {0, 5}, {1, 5}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, {6, 5}, {6, 4}, {7, 4}, {8, 4}, {9, 4}, {9, 3}, {8, 3}, {7, 3}, {6, 3}, {5, 3}, {5, 4}, {4, 4}, {3, 4}, {2, 4}, {2, 3}, {1, 3}, {0, 3}, {0, 2}, {1, 2}, {2, 2}},
			Length:  37,
			Latency: "21",
			Shout:   "",
		},
	}

	treeState := NewState(&state, 1)
	move := treeState.FindBestMove(true)

	if move == Up {
		t.Errorf("snake moved into other snake, %v", move)
	}
	if move == Left {
		t.Errorf("snake moved into self, %v", move)
	}
	if move == Right {
		t.Errorf("snake moved into other snake, %v", move)
	}
}

func TestFood2(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "aca4534f-c894-4e02-8aa3-a006cb1f4a54",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 251,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{6, 0}, {10, 2}, {5, 7}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_Wkf6xBfg7MDgMRPXPqKDkgDB",
					Name:    "Wild Heart ≡ƒºí",
					Health:  89,
					Head:    wire.Coord{5, 0},
					Body:    []wire.Coord{{5, 0}, {5, 1}, {5, 2}, {5, 3}, {5, 4}, {4, 4}, {4, 3}, {4, 2}, {3, 2}, {3, 3}, {3, 4}, {3, 5}, {4, 5}, {4, 6}, {4, 7}, {3, 7}, {2, 7}, {2, 6}, {2, 5}, {1, 5}, {1, 4}, {2, 4}, {2, 3}, {2, 2}, {2, 1}, {2, 0}},
					Length:  26,
					Latency: "103",
					Shout:   "≡ƒºí",
				},
				{
					ID:      "gs_4chD3GwkVwFdQbvK9tqMwb3V",
					Name:    "nomblegomble",
					Health:  7,
					Head:    wire.Coord{5, 10},
					Body:    []wire.Coord{{5, 10}, {4, 10}, {3, 10}, {2, 10}, {1, 10}, {0, 10}, {0, 9}},
					Length:  7,
					Latency: "57",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_4chD3GwkVwFdQbvK9tqMwb3V",
			Name:    "nomblegomble",
			Health:  7,
			Head:    wire.Coord{5, 10},
			Body:    []wire.Coord{{5, 10}, {4, 10}, {3, 10}, {2, 10}, {1, 10}, {0, 10}, {0, 9}},
			Length:  7,
			Latency: "57",
			Shout:   "",
		},
	}

	treeState := NewState(&state, 1)
	move := treeState.FindBestMove(true)

	if move == Right {
		t.Errorf("snake moved away from food, %v", move)
	}
	if move == Left {
		t.Errorf("snake moved into self, %v", move)
	}
	if move == Up {
		t.Errorf("snake moved into wall, %v", move)
	}
}

func TestFoodBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    wire.Battlesnake
		food     []wire.Coord
		expected Move
	}{
		{
			name: "eat when starving",
			input: wire.Battlesnake{
				Head:   wire.Coord{X: 5, Y: 5},
				Body:   []wire.Coord{{X: 5, Y: 5}, {X: 5, Y: 4}, {X: 6, Y: 4}},
				Health: 1,
				Length: 3,
				ID:     "my-id",
			},
			food:     []wire.Coord{{X: 6, Y: 5}},
			expected: Right,
		},
		{
			name: "go towards food when hungry",
			input: wire.Battlesnake{
				Head:   wire.Coord{X: 5, Y: 5},
				Body:   []wire.Coord{{X: 5, Y: 5}, {X: 5, Y: 4}, {X: 6, Y: 4}},
				Health: 20,
				Length: 3,
				ID:     "my-id",
			},
			food:     []wire.Coord{{X: 0, Y: 5}},
			expected: Left,
		},
	}

	for _, tc := range tests {
		state := wire.GameState{
			Board: wire.Board{
				Width:  12,
				Height: 12,
				Snakes: []wire.Battlesnake{tc.input},
				Food:   tc.food,
			},
			You: tc.input,
		}

		treeState := NewState(&state, 1)
		move := treeState.FindBestMove(true)

		if move != tc.expected {
			t.Errorf("%s: expected %s, got %s", tc.name, tc.expected, move)
		}
	}
}

// prefer a 50% chance to hitting the wall
func TestH2H01(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "7c1978a8-f0b5-4a58-9e0b-5df817230715",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 29,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{3, 0}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_MSJgRvk6tpfKQS9g93h9D8fD",
					Name:    "Secret Snake",
					Health:  89,
					Head:    wire.Coord{6, 3},
					Body:    []wire.Coord{{6, 3}, {7, 3}, {7, 2}, {8, 2}},
					Length:  4,
					Latency: "42",
					Shout:   "",
				},
				{
					ID:      "gs_JKHrcwtwQmP3XSPVwpFjyfC7",
					Name:    "trentren-vilu",
					Health:  99,
					Head:    wire.Coord{6, 9},
					Body:    []wire.Coord{{6, 9}, {5, 9}, {4, 9}, {3, 9}, {3, 8}, {4, 8}, {5, 8}, {5, 7}, {4, 7}},
					Length:  9,
					Latency: "43",
					Shout:   "",
				},
				{
					ID:      "gs_R3HxxD89jPQYjrrPmc8pYW7K",
					Name:    "nomblegomble",
					Health:  96,
					Head:    wire.Coord{7, 10},
					Body:    []wire.Coord{{7, 10}, {8, 10}, {9, 10}, {9, 9}, {9, 8}, {8, 8}},
					Length:  6,
					Latency: "68",
					Shout:   "",
				},
				{
					ID:      "gs_QGVk7YrR3QWPPXXdvYGhpdHG",
					Name:    "Serpentor",
					Health:  87,
					Head:    wire.Coord{7, 4},
					Body:    []wire.Coord{{7, 4}, {7, 5}, {7, 6}, {7, 7}, {8, 7}, {8, 6}},
					Length:  6,
					Latency: "92",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_R3HxxD89jPQYjrrPmc8pYW7K",
			Name:    "nomblegomble",
			Health:  96,
			Head:    wire.Coord{7, 10},
			Body:    []wire.Coord{{7, 10}, {8, 10}, {9, 10}, {9, 9}, {9, 8}, {8, 8}},
			Length:  6,
			Latency: "68",
			Shout:   "",
		},
	}

	// 1 is fine, >2 was returning 'dead' move since all H2H moves end in (possible death)
	// to fix this, add the 'lucky' move as a backup.
	treeState := NewState(&state, 2)
	move := treeState.FindBestMove(true)
	fmt.Printf("got move: %v\n", move)

	if move == Dead {
		t.Errorf("should never get '%v' as move", move)
	}
	if move == Right {
		t.Errorf("snake moved into self, %v", move)
	}
	if move == Up {
		t.Errorf("snake moved into wall, %v", move)
	}
}
