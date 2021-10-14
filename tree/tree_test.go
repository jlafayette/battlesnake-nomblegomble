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
	move, _ := treeState.FindBestMove(true)

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
	move, _ := treeState.FindBestMove(true)

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
	move, _ := treeState.FindBestMove(true)
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
	move, _ := treeState.FindBestMove(true)

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
	move, _ := treeState.FindBestMove(true)

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
			Game: wire.Game{
				Timeout: 500,
			},
			Board: wire.Board{
				Width:  12,
				Height: 12,
				Snakes: []wire.Battlesnake{tc.input},
				Food:   tc.food,
			},
			You: tc.input,
		}

		treeState := NewState(&state, 1)
		move, _ := treeState.FindBestMove(true)

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
	move, _ := treeState.FindBestMove(true)
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

func TestRespect01(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "31a26274-75d7-49d7-8dec-e34e3a348802",
			Ruleset: wire.Ruleset{
				Name:    "royale",
				Version: "v1.0.22",
			},
			Timeout: 600,
		},
		Turn: 82,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{10, 0}, {8, 9}, {10, 9}, {2, 0}, {7, 1}, {4, 10}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_jP3b937PT9XBMRGRp6JX3GmD",
					Name:    "ChoffesBattleSnakeV1",
					Health:  89,
					Head:    wire.Coord{3, 3},
					Body:    []wire.Coord{{3, 3}, {4, 3}, {4, 4}, {4, 5}, {3, 5}, {2, 5}, {2, 4}, {1, 4}, {1, 5}, {1, 6}, {1, 7}, {1, 8}},
					Length:  12,
					Latency: "168",
					Shout:   "",
				},
				{
					ID:      "gs_d7KD34x7mk4PBTP9RvdVSTFT",
					Name:    "king crimson",
					Health:  61,
					Head:    wire.Coord{6, 4},
					Body:    []wire.Coord{{6, 4}, {6, 3}, {5, 3}, {5, 4}, {5, 5}, {5, 6}, {6, 6}, {6, 7}},
					Length:  8,
					Latency: "175",
					Shout:   "",
				},
				{
					ID:      "gs_KCDr3MxmSC7Y6C3J6WdmW44H",
					Name:    "nomblegomble",
					Health:  29,
					Head:    wire.Coord{9, 5},
					Body:    []wire.Coord{{9, 5}, {9, 4}, {9, 3}, {9, 2}, {8, 2}, {7, 2}, {6, 2}, {5, 2}, {4, 2}},
					Length:  9,
					Latency: "31",
					Shout:   "",
				},
			},
			Hazards: []wire.Coord{{0, 10}, {1, 10}, {2, 10}, {3, 10}, {4, 10}, {5, 10}, {6, 10}, {7, 10}, {8, 10}, {9, 0}, {9, 1}, {9, 2}, {9, 3}, {9, 4}, {9, 5}, {9, 6}, {9, 7}, {9, 8}, {9, 9}, {9, 10}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}},
		},
		You: wire.Battlesnake{
			ID:      "gs_KCDr3MxmSC7Y6C3J6WdmW44H",
			Name:    "nomblegomble",
			Health:  29,
			Head:    wire.Coord{9, 5},
			Body:    []wire.Coord{{9, 5}, {9, 4}, {9, 3}, {9, 2}, {8, 2}, {7, 2}, {6, 2}, {5, 2}, {4, 2}},
			Length:  9,
			Latency: "31",
			Shout:   "",
		},
	}

	treeState := NewState(&state, 2)
	move, _ := treeState.FindBestMove(true)

	if move == Up {
		t.Errorf("snake moved needlessly through the sauce, %v", move)
	}
	if move == Right {
		t.Errorf("snake moved deeper into the sauce, %v", move)
	}
	if move == Down {
		t.Errorf("snake moved into self, %v", move)
	}
}

// Gotta get food when you need it, even in the sauce
func TestRespect02(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "e40a1d13-0c5f-45cf-a37c-eb6ff5ff39ee",
			Ruleset: wire.Ruleset{
				Name:    "royale",
				Version: "v1.0.22",
			},
			Timeout: 600,
		},
		Turn: 172,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{7, 0}, {6, 9}, {9, 5}, {5, 0}, {8, 7}, {3, 8}, {0, 3}, {10, 0}, {1, 0}, {8, 8}, {3, 7}, {10, 10}, {6, 5}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_dMhw4HRymGGycmVCDDJQr7FQ",
					Name:    "king crimson",
					Health:  64,
					Head:    wire.Coord{3, 5},
					Body:    []wire.Coord{{3, 5}, {3, 4}, {2, 4}, {2, 3}, {1, 3}, {1, 4}, {1, 5}, {1, 6}, {2, 6}, {3, 6}},
					Length:  10,
					Latency: "243",
					Shout:   "",
				},
				{
					ID:      "gs_3G66ck9qjChbvKgch6gqhfT8",
					Name:    "nomblegomble",
					Health:  23,
					Head:    wire.Coord{5, 1},
					Body:    []wire.Coord{{5, 1}, {6, 1}, {7, 1}, {8, 1}, {9, 1}, {9, 2}, {9, 3}, {8, 3}, {8, 2}, {7, 2}, {6, 2}, {6, 3}, {7, 3}},
					Length:  13,
					Latency: "38",
					Shout:   "",
				},
			},
			Hazards: []wire.Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {0, 9}, {0, 10}, {1, 0}, {1, 1}, {1, 9}, {1, 10}, {2, 0}, {2, 1}, {2, 9}, {2, 10}, {3, 0}, {3, 1}, {3, 9}, {3, 10}, {4, 0}, {4, 1}, {4, 9}, {4, 10}, {5, 0}, {5, 1}, {5, 9}, {5, 10}, {6, 0}, {6, 1}, {6, 9}, {6, 10}, {7, 0}, {7, 1}, {7, 9}, {7, 10}, {8, 0}, {8, 1}, {8, 9}, {8, 10}, {9, 0}, {9, 1}, {9, 9}, {9, 10}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}},
		},
		You: wire.Battlesnake{
			ID:      "gs_3G66ck9qjChbvKgch6gqhfT8",
			Name:    "nomblegomble",
			Health:  23,
			Head:    wire.Coord{5, 1},
			Body:    []wire.Coord{{5, 1}, {6, 1}, {7, 1}, {8, 1}, {9, 1}, {9, 2}, {9, 3}, {8, 3}, {8, 2}, {7, 2}, {6, 2}, {6, 3}, {7, 3}},
			Length:  13,
			Latency: "38",
			Shout:   "",
		},
	}

	treeState := NewState(&state, 3)
	move, _ := treeState.FindBestMove(true)

	if move == Up {
		t.Errorf("snake moved away from food, %v", move)
	}
	if move == Right {
		t.Errorf("snake moved into self, %v", move)
	}
	if move == Left {
		t.Errorf("snake moved away from food, %v", move)
	}
}

// Sauce area < normal area
func TestRespect03(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "7d56539b-ef87-4838-88e7-a12201c0057c",
			Ruleset: wire.Ruleset{
				Name:    "royale",
				Version: "v1.0.22",
			},
			Timeout: 600,
		},
		Turn: 110,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{1, 10}, {3, 9}, {0, 10}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_PdHwHTYfm6BTVV36tGt3JTwG",
					Name:    "nomblegomble",
					Health:  100,
					Head:    wire.Coord{4, 0},
					Body:    []wire.Coord{{4, 0}, {4, 1}, {4, 2}, {3, 2}, {3, 3}, {3, 4}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {0, 6}, {0, 7}, {1, 7}, {2, 7}, {2, 7}},
					Length:  15,
					Latency: "48",
					Shout:   "",
				},
				{
					ID:      "gs_4rmj8dKwfF4q88cRQqfvFjrb",
					Name:    "bsnekGo",
					Health:  93,
					Head:    wire.Coord{9, 3},
					Body:    []wire.Coord{{9, 3}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {8, 5}, {8, 4}, {8, 3}},
					Length:  8,
					Latency: "60",
					Shout:   "",
				},
				{
					ID:      "gs_crPB88ccjJ3t84cKbSgXYyRM",
					Name:    "trentren-vilu",
					Health:  96,
					Head:    wire.Coord{9, 9},
					Body:    []wire.Coord{{9, 9}, {8, 9}, {8, 8}, {7, 8}, {6, 8}, {5, 8}, {5, 9}, {5, 10}, {4, 10}, {4, 9}, {4, 8}, {4, 7}, {5, 7}},
					Length:  13,
					Latency: "48",
					Shout:   "",
				},
			},
			Hazards: []wire.Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {0, 9}, {0, 10}, {1, 0}, {1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}, {1, 7}, {1, 8}, {1, 9}, {1, 10}, {2, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {2, 7}, {2, 8}, {2, 9}, {2, 10}, {3, 10}, {4, 10}, {5, 10}, {6, 10}, {7, 10}, {8, 10}, {9, 10}, {10, 10}},
		},
		You: wire.Battlesnake{
			ID:      "gs_PdHwHTYfm6BTVV36tGt3JTwG",
			Name:    "nomblegomble",
			Health:  100,
			Head:    wire.Coord{4, 0},
			Body:    []wire.Coord{{4, 0}, {4, 1}, {4, 2}, {3, 2}, {3, 3}, {3, 4}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {0, 6}, {0, 7}, {1, 7}, {2, 7}, {2, 7}},
			Length:  15,
			Latency: "48",
			Shout:   "",
		},
	}

	treeState := NewState(&state, 1)
	move, _ := treeState.FindBestMove(true)

	if move == Left {
		t.Errorf("snake moved into space with sauce and no food, %v", move)
	}
	if move == Up {
		t.Errorf("snake moved into self, %v", move)
	}
	if move == Down {
		t.Errorf("snake moved into wall, %v", move)
	}
}

func TestH2HTieBetterThanLoss01(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "6c374e7c-2611-4c3f-a1d4-79d887f6589a",
			Ruleset: wire.Ruleset{
				Name:    "royale",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 192,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{10, 1}, {6, 0}, {10, 0}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_MKPct3wDS9gfK8qTWKCkPYXX",
					Name:    "WhitishMeteor",
					Health:  89,
					Head:    wire.Coord{6, 6},
					Body:    []wire.Coord{{6, 6}, {6, 7}, {5, 7}, {5, 6}, {4, 6}, {3, 6}, {2, 6}, {1, 6}, {0, 6}, {0, 7}, {0, 8}, {0, 9}, {1, 9}, {2, 9}, {3, 9}, {3, 8}},
					Length:  16,
					Latency: "388",
					Shout:   "",
				},
				{
					ID:      "gs_4k8x3CJVR3Y6r7j8qB76fSXb",
					Name:    "nomblegomble",
					Health:  67,
					Head:    wire.Coord{7, 5},
					Body:    []wire.Coord{{7, 5}, {7, 6}, {7, 7}, {8, 7}, {9, 7}, {10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {9, 9}, {8, 9}, {7, 9}, {7, 10}, {6, 10}, {5, 10}},
					Length:  16,
					Latency: "423",
					Shout:   "",
				},
			},
			Hazards: []wire.Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 10}, {1, 0}, {1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 10}, {2, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}, {2, 10}, {3, 0}, {3, 1}, {3, 2}, {3, 3}, {3, 4}, {3, 10}, {4, 0}, {4, 1}, {4, 2}, {4, 3}, {4, 4}, {4, 10}, {5, 0}, {5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 10}, {6, 0}, {6, 1}, {6, 2}, {6, 3}, {6, 4}, {6, 10}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {7, 10}, {8, 0}, {8, 1}, {8, 2}, {8, 3}, {8, 4}, {8, 5}, {8, 6}, {8, 7}, {8, 8}, {8, 9}, {8, 10}, {9, 0}, {9, 1}, {9, 2}, {9, 3}, {9, 4}, {9, 5}, {9, 6}, {9, 7}, {9, 8}, {9, 9}, {9, 10}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}},
		},
		You: wire.Battlesnake{
			ID:      "gs_4k8x3CJVR3Y6r7j8qB76fSXb",
			Name:    "nomblegomble",
			Health:  67,
			Head:    wire.Coord{7, 5},
			Body:    []wire.Coord{{7, 5}, {7, 6}, {7, 7}, {8, 7}, {9, 7}, {10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {9, 9}, {8, 9}, {7, 9}, {7, 10}, {6, 10}, {5, 10}},
			Length:  16,
			Latency: "423",
			Shout:   "",
		},
	}

	treeState := NewState(&state, 10)
	move, _ := treeState.FindBestMove(true)

	// if move == Up {
	// 	t.Errorf("snake moved into too small of space, %v", move)
	// }
	if move == Up {
		t.Errorf("snake moved into self, %v", move)
	}
	// if move == Left {
	// 	t.Errorf("snake moved into h2h, %v", move)
	// }
}

func TestMaybeDontMoveIntoCornerAndDie(t *testing.T) {

	t.Skip("working on 2, not on 3 deep")

	state := wire.GameState{
		Game: wire.Game{
			ID: "32774ee7-475c-4c86-892c-66422fa0328c",
			Ruleset: wire.Ruleset{
				Name:    "royale",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 99,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{1, 5}, {10, 1}, {8, 0}, {5, 2}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_Pfdf33SjMwy7hVbKkXxTkWrB",
					Name:    "Pea Eater",
					Health:  98,
					Head:    wire.Coord{3, 6},
					Body:    []wire.Coord{{3, 6}, {2, 6}, {2, 7}, {2, 8}, {2, 9}, {3, 9}, {3, 8}, {4, 8}, {4, 9}, {5, 9}, {5, 8}, {5, 7}, {6, 7}, {6, 6}, {5, 6}},
					Length:  15,
					Latency: "398",
					Shout:   "Om nom nom nom",
				},
				{
					ID:      "gs_j6D6Kp7rTHtSwwkH7CYKKtb4",
					Name:    "WhitishMeteor",
					Health:  97,
					Head:    wire.Coord{2, 3},
					Body:    []wire.Coord{{2, 3}, {2, 2}, {2, 1}, {2, 0}, {3, 0}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {7, 1}, {7, 2}},
					Length:  11,
					Latency: "380",
					Shout:   "",
				},
				{
					ID:      "gs_8trTvWDjGQdSdTqJD67wrSVX",
					Name:    "nomblegomble",
					Health:  83,
					Head:    wire.Coord{8, 9},
					Body:    []wire.Coord{{8, 9}, {7, 9}, {7, 10}, {6, 10}, {6, 9}, {6, 8}, {7, 8}, {7, 7}, {7, 6}, {8, 6}, {9, 6}},
					Length:  11,
					Latency: "473",
					Shout:   "3",
				},
			},
			Hazards: []wire.Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {0, 9}, {0, 10}, {8, 0}, {8, 1}, {8, 2}, {8, 3}, {8, 4}, {8, 5}, {8, 6}, {8, 7}, {8, 8}, {8, 9}, {8, 10}, {9, 0}, {9, 1}, {9, 2}, {9, 3}, {9, 4}, {9, 5}, {9, 6}, {9, 7}, {9, 8}, {9, 9}, {9, 10}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}},
		},
		You: wire.Battlesnake{
			ID:      "gs_8trTvWDjGQdSdTqJD67wrSVX",
			Name:    "nomblegomble",
			Health:  83,
			Head:    wire.Coord{8, 9},
			Body:    []wire.Coord{{8, 9}, {7, 9}, {7, 10}, {6, 10}, {6, 9}, {6, 8}, {7, 8}, {7, 7}, {7, 6}, {8, 6}, {9, 6}},
			Length:  11,
			Latency: "473",
			Shout:   "3",
		},
	}

	treeState := NewState(&state, 10)
	move, _ := treeState.FindBestMove(true)

	if move == Up {
		t.Errorf("snake moved into corner, just stop it, %v", move)
	}
	if move == Right {
		t.Errorf("snake moved deeper into sauce, %v", move)
	}
	if move == Left {
		t.Errorf("snake moved into self, %v", move)
	}
}

func TestEatTheFood01(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "9ed99ee5-a8d8-4270-a706-d6f1189ffd2a",
			Ruleset: wire.Ruleset{
				Name:    "royale",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 143,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 9}, {9, 8}, {1, 3}, {9, 10}, {7, 9}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_D3QVPfpvDC3txQSQfFfPF7Bf",
					Name:    "Pr├╝zze v2",
					Health:  13,
					Head:    wire.Coord{1, 2},
					Body:    []wire.Coord{{1, 2}, {2, 2}, {3, 2}, {4, 2}, {5, 2}, {6, 2}, {7, 2}, {8, 2}},
					Length:  8,
					Latency: "369",
					Shout:   ", t=347",
				},
				{
					ID:      "gs_X8SV44B8RvjHMM6b6bPYtMhW",
					Name:    "WhitishMeteor",
					Health:  85,
					Head:    wire.Coord{4, 9},
					Body:    []wire.Coord{{4, 9}, {4, 10}, {5, 10}, {6, 10}, {6, 9}, {6, 8}, {6, 7}, {7, 7}, {7, 6}, {6, 6}, {5, 6}, {5, 7}, {5, 8}, {4, 8}, {3, 8}, {2, 8}, {2, 7}},
					Length:  17,
					Latency: "342",
					Shout:   "",
				},
				{
					ID:      "gs_cG4bdcDHbKCvgmthMkJVxR3b",
					Name:    "nomblegomble",
					Health:  60,
					Head:    wire.Coord{2, 3},
					Body:    []wire.Coord{{2, 3}, {3, 3}, {4, 3}, {5, 3}, {5, 4}, {6, 4}, {7, 4}, {8, 4}, {8, 5}, {7, 5}, {6, 5}, {5, 5}, {4, 5}, {4, 4}},
					Length:  14,
					Latency: "473",
					Shout:   "5",
				},
			},
			Hazards: []wire.Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {0, 9}, {0, 10}, {1, 8}, {1, 9}, {1, 10}, {2, 8}, {2, 9}, {2, 10}, {3, 8}, {3, 9}, {3, 10}, {4, 8}, {4, 9}, {4, 10}, {5, 8}, {5, 9}, {5, 10}, {6, 8}, {6, 9}, {6, 10}, {7, 8}, {7, 9}, {7, 10}, {8, 0}, {8, 1}, {8, 2}, {8, 3}, {8, 4}, {8, 5}, {8, 6}, {8, 7}, {8, 8}, {8, 9}, {8, 10}, {9, 0}, {9, 1}, {9, 2}, {9, 3}, {9, 4}, {9, 5}, {9, 6}, {9, 7}, {9, 8}, {9, 9}, {9, 10}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}},
		},
		You: wire.Battlesnake{
			ID:      "gs_cG4bdcDHbKCvgmthMkJVxR3b",
			Name:    "nomblegomble",
			Health:  60,
			Head:    wire.Coord{2, 3},
			Body:    []wire.Coord{{2, 3}, {3, 3}, {4, 3}, {5, 3}, {5, 4}, {6, 4}, {7, 4}, {8, 4}, {8, 5}, {7, 5}, {6, 5}, {5, 5}, {4, 5}, {4, 4}},
			Length:  14,
			Latency: "473",
			Shout:   "5",
		},
	}

	treeState := NewState(&state, 3)
	move, _ := treeState.FindBestMove(true)

	if move == Up {
		t.Errorf("snake moved away from delicious food, %v", move)
	}
	if move == Right {
		t.Errorf("snake moved into self, %v", move)
	}
	if move == Down {
		t.Errorf("snake moved into other snake, %v", move)
	}
}

func TestMaybeDontMoveIntoCornerAndDie02(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "a697f868-d48d-4fd0-b807-9a89a225e147",
			Ruleset: wire.Ruleset{
				Name:    "royale",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 125,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{10, 0}, {10, 1}, {7, 0}, {10, 3}, {9, 3}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_H3cVVxcMFMMfQx3rcr6JWyK6",
					Name:    "nomblegomble",
					Health:  100,
					Head:    wire.Coord{3, 10},
					Body:    []wire.Coord{{3, 10}, {3, 9}, {3, 8}, {2, 8}, {1, 8}, {1, 7}, {0, 7}, {0, 6}, {0, 5}, {1, 5}, {1, 6}, {2, 6}, {2, 6}},
					Length:  13,
					Latency: "473",
					Shout:   "8",
				},
				{
					ID:      "gs_4bVBjpSkGxHQDr9bHm6GRPg9",
					Name:    "Nessegrev-gamma",
					Health:  79,
					Head:    wire.Coord{5, 6},
					Body:    []wire.Coord{{5, 6}, {5, 5}, {5, 4}, {4, 4}, {4, 3}, {3, 3}, {2, 3}, {1, 3}, {1, 4}, {2, 4}, {3, 4}, {3, 5}, {4, 5}},
					Length:  13,
					Latency: "402",
					Shout:   "",
				},
			},
			Hazards: []wire.Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {0, 9}, {0, 10}, {1, 0}, {1, 1}, {1, 9}, {1, 10}, {2, 0}, {2, 1}, {2, 9}, {2, 10}, {3, 0}, {3, 1}, {3, 9}, {3, 10}, {4, 0}, {4, 1}, {4, 9}, {4, 10}, {5, 0}, {5, 1}, {5, 9}, {5, 10}, {6, 0}, {6, 1}, {6, 9}, {6, 10}, {7, 0}, {7, 1}, {7, 9}, {7, 10}, {8, 0}, {8, 1}, {8, 9}, {8, 10}, {9, 0}, {9, 1}, {9, 9}, {9, 10}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}},
		},
		You: wire.Battlesnake{
			ID:      "gs_H3cVVxcMFMMfQx3rcr6JWyK6",
			Name:    "nomblegomble",
			Health:  100,
			Head:    wire.Coord{3, 10},
			Body:    []wire.Coord{{3, 10}, {3, 9}, {3, 8}, {2, 8}, {1, 8}, {1, 7}, {0, 7}, {0, 6}, {0, 5}, {1, 5}, {1, 6}, {2, 6}, {2, 6}},
			Length:  13,
			Latency: "473",
			Shout:   "8",
		},
	}

	// 1 is fine, >2 all go left
	// seem to be overly afraid of H2H (this was true when the test was added)
	treeState := NewState(&state, 3)
	move, _ := treeState.FindBestMove(true)

	if move == Left {
		t.Errorf("snake moved into the corner and will definately die, %v", move)
	}
	if move == Down {
		t.Errorf("snake moved into self, %v", move)
	}
	if move == Up {
		t.Errorf("snake moved into wall, %v", move)
	}
}

// don't need this long term, but useful for testing obviously bad pattern
// of not counting roughly equal options that other snakes have
func TestSanityCheck01(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "b6af3e62-95fe-4b6a-8807-3de99a2e776f",
			Ruleset: wire.Ruleset{
				Name:    "royale",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 18,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 9}, {1, 1}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_GhjgBb7YxVDYttTYPvb4SFkC",
					Name:    "Salazar Slitherin",
					Health:  92,
					Head:    wire.Coord{3, 3},
					Body:    []wire.Coord{{3, 3}, {4, 3}, {5, 3}, {5, 2}, {6, 2}},
					Length:  5,
					Latency: "151",
					Shout:   "5 4 17",
				},
				{
					ID:      "gs_Vq7WM3YmkWtfdSFQ3yGdGwC6",
					Name:    "Devious Devin",
					Health:  84,
					Head:    wire.Coord{1, 9},
					Body:    []wire.Coord{{1, 9}, {2, 9}, {3, 9}, {4, 9}},
					Length:  4,
					Latency: "431",
					Shout:   "",
				},
				{
					ID:      "gs_SSFrSBMjgJtCFgFS8ryDqDCY",
					Name:    "nomblegomble",
					Health:  94,
					Head:    wire.Coord{4, 10},
					Body:    []wire.Coord{{4, 10}, {5, 10}, {6, 10}, {7, 10}, {8, 10}},
					Length:  5,
					Latency: "474",
					Shout:   "3",
				},
				{
					ID:      "gs_jXrjh6XCWdfxxRQ6kDc7hwXQ",
					Name:    "≡ƒÆÇ≡ƒÆÇ≡ƒÆÇ≡ƒÆÇ≡ƒÆÇ",
					Health:  91,
					Head:    wire.Coord{6, 6},
					Body:    []wire.Coord{{6, 6}, {7, 6}, {8, 6}, {8, 5}, {8, 4}, {9, 4}},
					Length:  6,
					Latency: "84",
					Shout:   "",
				},
			},
			Hazards: []wire.Coord{},
		},
		You: wire.Battlesnake{
			ID:      "gs_SSFrSBMjgJtCFgFS8ryDqDCY",
			Name:    "nomblegomble",
			Health:  94,
			Head:    wire.Coord{4, 10},
			Body:    []wire.Coord{{4, 10}, {5, 10}, {6, 10}, {7, 10}, {8, 10}},
			Length:  5,
			Latency: "474",
			Shout:   "3",
		},
	}

	treeState := NewState(&state, 3)
	move, _ := treeState.FindBestMove(true)

	if move == Left {
		t.Errorf("snake moved into trapped corner (will die), %v", move)
	}
	if move == Right {
		t.Errorf("snake moved into self, %v", move)
	}
	if move == Up {
		t.Errorf("snake moved into wall, %v", move)
	}
}

// Don't try for a H2H when there are obviously better options that don't
// involve a a 50/50 chance of dying
// 1-5 gave left (good), then 6+ gave up (bad)...
// fixed by increasing the prune threshold from 50->100 (75 was enough, but 100
// seems like a good value to try)
func TestBadH2H01(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "77746b26-22a2-4013-a3c8-2791c6da6522",
			Ruleset: wire.Ruleset{
				Name:    "royale",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 161,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 0}, {0, 6}, {8, 8}, {10, 6}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_9vQ9WghMBmJt7fkhh4y686qF",
					Name:    "nomblegomble",
					Health:  97,
					Head:    wire.Coord{1, 8},
					Body:    []wire.Coord{{1, 8}, {1, 7}, {2, 7}, {2, 6}, {2, 5}, {2, 4}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, {5, 4}, {4, 4}, {3, 4}, {3, 5}},
					Length:  14,
					Latency: "473",
					Shout:   "7",
				},
				{
					ID:      "gs_gXhCJgkDjpFqTMB9bTfvPfgb",
					Name:    "≡ƒÆÇ≡ƒÆÇ≡ƒÆÇ≡ƒÆÇ≡ƒÆÇ",
					Health:  95,
					Head:    wire.Coord{2, 9},
					Body:    []wire.Coord{{2, 9}, {3, 9}, {3, 8}, {4, 8}, {5, 8}, {5, 7}, {5, 6}, {5, 5}, {6, 5}, {7, 5}, {7, 6}, {7, 7}, {8, 7}, {9, 7}, {9, 8}, {9, 9}, {9, 10}},
					Length:  17,
					Latency: "72",
					Shout:   "",
				},
			},
			Hazards: []wire.Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {0, 9}, {0, 10}, {1, 0}, {1, 1}, {1, 10}, {2, 0}, {2, 1}, {2, 10}, {3, 0}, {3, 1}, {3, 10}, {4, 0}, {4, 1}, {4, 10}, {5, 0}, {5, 1}, {5, 10}, {6, 0}, {6, 1}, {6, 10}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {7, 5}, {7, 6}, {7, 7}, {7, 8}, {7, 9}, {7, 10}, {8, 0}, {8, 1}, {8, 2}, {8, 3}, {8, 4}, {8, 5}, {8, 6}, {8, 7}, {8, 8}, {8, 9}, {8, 10}, {9, 0}, {9, 1}, {9, 2}, {9, 3}, {9, 4}, {9, 5}, {9, 6}, {9, 7}, {9, 8}, {9, 9}, {9, 10}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}},
		},
		You: wire.Battlesnake{
			ID:      "gs_9vQ9WghMBmJt7fkhh4y686qF",
			Name:    "nomblegomble",
			Health:  97,
			Head:    wire.Coord{1, 8},
			Body:    []wire.Coord{{1, 8}, {1, 7}, {2, 7}, {2, 6}, {2, 5}, {2, 4}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, {5, 4}, {4, 4}, {3, 4}, {3, 5}},
			Length:  14,
			Latency: "473",
			Shout:   "7",
		},
	}

	treeState := NewState(&state, 8)
	move, _ := treeState.FindBestMove(true)

	if move == Up {
		t.Errorf("snake moved into H2H needlessly, %v", move)
	}
	if move == Down {
		t.Errorf("snake moved into self, %v", move)
	}
	if move == Right {
		t.Errorf("snake moved into H2H or too small a space, %v", move)
	}
}

func Test__NAME__(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "6f5c63a3-ab49-45e5-8751-b261e1ca6b71",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 2,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{10, 4}, {5, 5}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_SFxtMhQthd39Q6fHMBSh7vk6",
					Name:    "nomblegomble-dev",
					Health:  98,
					Head:    wire.Coord{9, 7},
					Body:    []wire.Coord{{9, 7}, {9, 6}, {9, 5}},
					Length:  3,
					Latency: "52",
					Shout:   "",
				},
				{
					ID:      "gs_mb6V4ydFJSjTcKcfqqpgTWTP",
					Name:    "Stabby McStabbertooth",
					Health:  100,
					Head:    wire.Coord{2, 6},
					Body:    []wire.Coord{{2, 6}, {1, 6}, {1, 5}, {1, 5}},
					Length:  4,
					Latency: "36",
					Shout:   "",
				},
				{
					ID:      "gs_rRMt9kvpbKYy3HqxMvCkTFPb",
					Name:    "Stabby McStabbertooth",
					Health:  100,
					Head:    wire.Coord{0, 2},
					Body:    []wire.Coord{{0, 2}, {1, 2}, {1, 1}, {1, 1}},
					Length:  4,
					Latency: "52",
					Shout:   "",
				},
				{
					ID:      "gs_tdFX33gK6mwXFfdywSYxTV44",
					Name:    "Stabby McStabbertooth",
					Health:  100,
					Head:    wire.Coord{8, 2},
					Body:    []wire.Coord{{8, 2}, {9, 2}, {9, 1}, {9, 1}},
					Length:  4,
					Latency: "26",
					Shout:   "",
				},
				{
					ID:      "gs_TVfWbcf4thmQXhJ8VVVGkyJV",
					Name:    "Stabby McStabbertooth",
					Health:  100,
					Head:    wire.Coord{2, 8},
					Body:    []wire.Coord{{2, 8}, {1, 8}, {1, 9}, {1, 9}},
					Length:  4,
					Latency: "54",
					Shout:   "",
				},
				{
					ID:      "gs_XcfDWWMJdK4KPrJt7SQj73S3",
					Name:    "Stabby McStabbertooth",
					Health:  100,
					Head:    wire.Coord{10, 8},
					Body:    []wire.Coord{{10, 8}, {9, 8}, {9, 9}, {9, 9}},
					Length:  4,
					Latency: "54",
					Shout:   "",
				},
				{
					ID:      "gs_9xghdDSShDDxWygC8GbCDY4D",
					Name:    "Stabby McStabbertooth",
					Health:  100,
					Head:    wire.Coord{4, 0},
					Body:    []wire.Coord{{4, 0}, {5, 0}, {5, 1}, {5, 1}},
					Length:  4,
					Latency: "53",
					Shout:   "",
				},
				{
					ID:      "gs_VVGTx4h3gQtTSVxQHkVMrhfb",
					Name:    "Stabby McStabbertooth",
					Health:  100,
					Head:    wire.Coord{4, 8},
					Body:    []wire.Coord{{4, 8}, {5, 8}, {5, 9}, {5, 9}},
					Length:  4,
					Latency: "52",
					Shout:   "",
				},
			},
			Hazards: []wire.Coord{},
		},
		You: wire.Battlesnake{
			ID:      "gs_SFxtMhQthd39Q6fHMBSh7vk6",
			Name:    "nomblegomble-dev",
			Health:  98,
			Head:    wire.Coord{9, 7},
			Body:    []wire.Coord{{9, 7}, {9, 6}, {9, 5}},
			Length:  3,
			Latency: "52",
			Shout:   "",
		},
	}

	treeState := NewState(&state, 1)
	move, _ := treeState.FindBestMove(true)

	if move == Up {
		t.Errorf("snake moved into other snake, %v", move)
	}
	if move == Right {
		t.Errorf("snake moved into bad h2h, %v", move)
	}
	if move == Down {
		t.Errorf("snake moved into wall, %v", move)
	}
}

// -- benchmarks

var result Move
var benchState wire.GameState

func benchmark01(depth int, b *testing.B) {
	b.ReportAllocs()

	var r Move

	state := wire.GameState{
		Game: wire.Game{
			ID: "eb8f2aa1-fce4-473a-b33d-a8181573c478",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.22",
			},
			Timeout: 500,
		},
		Turn: 2,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 6}, {5, 5}, {8, 4}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_jCFX9pP8FBhckB9BTSBTvgRF",
					Name:    "haspid",
					Health:  98,
					Head:    wire.Coord{3, 5},
					Body:    []wire.Coord{{3, 5}, {2, 5}, {1, 5}},
					Length:  3,
					Latency: "264",
					Shout:   "",
				},
				{
					ID:      "gs_pkWJTmghRJfbWk6JGmv3y94V",
					Name:    "trentren-vilu",
					Health:  100,
					Head:    wire.Coord{4, 8},
					Body:    []wire.Coord{{4, 8}, {5, 8}, {5, 9}, {5, 9}},
					Length:  4,
					Latency: "44",
					Shout:   "",
				},
				{
					ID:      "gs_VChQfdbhMP8MQ8PCk8wqmrkS",
					Name:    "nomblegomble",
					Health:  100,
					Head:    wire.Coord{10, 6},
					Body:    []wire.Coord{{10, 6}, {9, 6}, {9, 5}, {9, 5}},
					Length:  4,
					Latency: "137",
					Shout:   "",
				},
				{
					ID:      "gs_wyd7SF9TfgCWqGVwvmwCtRPb",
					Name:    "caicai-vilu",
					Health:  100,
					Head:    wire.Coord{0, 0},
					Body:    []wire.Coord{{0, 0}, {1, 0}, {1, 1}, {1, 1}},
					Length:  4,
					Latency: "42",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_VChQfdbhMP8MQ8PCk8wqmrkS",
			Name:    "nomblegomble",
			Health:  100,
			Head:    wire.Coord{10, 6},
			Body:    []wire.Coord{{10, 6}, {9, 6}, {9, 5}, {9, 5}},
			Length:  4,
			Latency: "137",
			Shout:   "",
		},
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		treeState := NewState(&state, depth)
		r, _ = treeState.FindBestMove(false)
		// b.Error("failed!")
	}

	result = r
}

func Benchmark01_1(b *testing.B) { benchmark01(1, b) }
func Benchmark01_2(b *testing.B) { benchmark01(2, b) }
func Benchmark01_3(b *testing.B) { benchmark01(3, b) }
