package main

import (
	"testing"
)

func TestNeckAvoidance(t *testing.T) {
	tests := []struct {
		name  string
		input Battlesnake
		noGo  string
	}{
		{
			name: "neck avoidance 1",
			input: Battlesnake{
				// Length 3, facing right
				Head: Coord{X: 2, Y: 0},
				Body: []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
			}, noGo: "left",
		},
		{
			name: "neck avoidance 2",
			input: Battlesnake{
				// Length 3, facing left
				Head: Coord{X: 7, Y: 0},
				Body: []Coord{{X: 7, Y: 0}, {X: 8, Y: 0}, {X: 9, Y: 0}},
			}, noGo: "right",
		},
		{
			name: "neck avoidance 3",
			input: Battlesnake{
				// Length 3, facing up
				Head: Coord{X: 5, Y: 10},
				Body: []Coord{{X: 5, Y: 9}, {X: 5, Y: 8}, {X: 5, Y: 7}},
			}, noGo: "down",
		},
		{
			name: "neck avoidance 4",
			input: Battlesnake{
				// Length 3, facing down
				Head: Coord{X: 5, Y: 1},
				Body: []Coord{{X: 5, Y: 2}, {X: 5, Y: 3}, {X: 5, Y: 4}},
			}, noGo: "up",
		},
	}

	for _, tc := range tests {
		state := GameState{
			Board: Board{
				Snakes: []Battlesnake{tc.input},
				Width:  11,
				Height: 11,
			},
			You: tc.input,
		}

		nextMove := move(state)

		if nextMove.Move == tc.noGo {
			t.Errorf("%s: snake moved onto its own neck, %s", tc.name, nextMove.Move)
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
		name     string
		input    Battlesnake
		intoNeck string
		intoWall []string
	}{
		{
			name: "wall avoidance 1",
			input: Battlesnake{
				// Lower left corner
				Head: Coord{X: 0, Y: 0},
				Body: []Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			},
			intoNeck: "right",
			intoWall: []string{"left", "down"},
		},
		{
			name: "wall avoidance 2",
			input: Battlesnake{
				// top right corner
				Head: Coord{X: 11, Y: 11},
				Body: []Coord{{X: 11, Y: 11}, {X: 10, Y: 11}, {X: 9, Y: 11}},
			},
			intoNeck: "left",
			intoWall: []string{"up", "right"},
		},
		{
			name: "wall avoidance 3",
			input: Battlesnake{
				// bottom right corner (facing down)
				Head: Coord{X: 11, Y: 0},
				Body: []Coord{{X: 11, Y: 0}, {X: 11, Y: 1}, {X: 11, Y: 2}},
			},
			intoNeck: "up",
			intoWall: []string{"down", "right"},
		},
		{
			name: "wall avoidance 4",
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
			t.Errorf("%s: snake moved onto its own neck, %s", tc.name, nextMove.Move)
		}
		if contains(tc.intoWall, nextMove.Move) {
			t.Errorf("%s: snake moved into a wall, %s", tc.name, nextMove.Move)
		}
	}
}

func TestSelfAvoidance(t *testing.T) {
	tests := []struct {
		name     string
		input    Battlesnake
		intoSelf []string
		intoWall []string
	}{
		{
			name: "body check",
			input: Battlesnake{
				Head: Coord{X: 5, Y: 5},
				Body: []Coord{{X: 5, Y: 5}, {X: 5, Y: 4}, {X: 6, Y: 4}, {X: 6, Y: 5}, {X: 6, Y: 6}, {X: 5, Y: 6}, {X: 4, Y: 6}},
			},
			intoSelf: []string{"up", "right", "down"},
		},
		// tail is ok if not at full health
		{
			name: "tail chase ok 1",
			input: Battlesnake{
				Head:   Coord{X: 11, Y: 0},
				Body:   []Coord{{X: 11, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 1}, {X: 11, Y: 1}},
				Health: 99,
			},
			intoSelf: []string{"left"},
			intoWall: []string{"down", "right"},
		},
		{
			name: "tail chase ok 2",
			input: Battlesnake{
				Head:   Coord{X: 0, Y: 11},
				Body:   []Coord{{X: 0, Y: 11}, {X: 1, Y: 11}, {X: 1, Y: 10}, {X: 0, Y: 10}},
				Health: 99,
			},
			intoSelf: []string{"right"},
			intoWall: []string{"up", "left"},
		},
		{
			name: "tail chase not ok (just eaten)",
			input: Battlesnake{
				Head:   Coord{X: 1, Y: 1},
				Body:   []Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}},
				Health: 100,
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
			t.Errorf("%s: snake moved into self, %s", tc.name, nextMove.Move)
		}
		if contains(tc.intoWall, nextMove.Move) {
			t.Errorf("%s: snake moved into wall, %s", tc.name, nextMove.Move)
		}
	}
}

func TestHead2Head(t *testing.T) {
	tests := []struct {
		me       Battlesnake
		other    Battlesnake
		expected string
	}{
		{
			me: Battlesnake{
				ID:   "snake-508e96ac-94ad-11ea-bb37",
				Head: Coord{X: 0, Y: 0},
				Body: []Coord{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
			},
			other: Battlesnake{
				ID:   "snake-b67f4906-94ae-11ea-bb37",
				Head: Coord{X: 2, Y: 0},
				Body: []Coord{{X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}},
			},
			expected: "right",
		},
	}
	for _, tc := range tests {
		state := GameState{
			Board: Board{
				Width:  12,
				Height: 12,
				Snakes: []Battlesnake{tc.me, tc.other},
			},
			You: tc.me,
		}

		nextMove := move(state)

		if nextMove.Move != tc.expected {
			t.Errorf("Prefer head2head over the wall or self, expected right, got %s", nextMove.Move)
		}
	}
}

func TestHead2HeadJson1(t *testing.T) {
	tests := []struct {
		state GameState
		ok    []string
		notOk []string
	}{
		{
			state: GameState{
				Board: Board{
					Width:  11,
					Height: 11,
					Food:   []Coord{{5, 5}, {6, 10}, {4, 2}, {9, 2}, {0, 5}},
					Snakes: []Battlesnake{
						{
							ID:     "snake-0-id",
							Name:   "snake0",
							Health: 97,
							Head:   Coord{5, 6},
							Body:   []Coord{{5, 6}, {5, 7}, {5, 8}},
							Length: 3,
						},
						{
							ID:     "snake-1-id",
							Name:   "snake1",
							Health: 97,
							Head:   Coord{5, 4},
							Body:   []Coord{{5, 4}, {5, 3}, {5, 2}},
							Length: 3,
						},
					},
				},
				You: Battlesnake{
					ID:     "snake-0-id",
					Name:   "snake0",
					Health: 97,
					Head:   Coord{5, 6},
					Body:   []Coord{{5, 6}, {5, 7}, {5, 8}},
					Length: 3,
				},
			},
			ok:    []string{"left", "right"},
			notOk: []string{"down", "up"},
		},
	}
	for _, tc := range tests {

		nextMove := move(tc.state)

		for _, badMove := range tc.notOk {
			if nextMove.Move == badMove {
				t.Errorf("expected one of %v, got %s", tc.ok, nextMove.Move)
			}
		}
		found := false
		for _, okMove := range tc.ok {
			if nextMove.Move == okMove {
				found = true
			}
		}
		if !found {
			t.Errorf("expected one of %v, got %s", tc.ok, nextMove.Move)
		}
	}
}

func TestSpace(t *testing.T) {
	tests := []struct {
		name         string
		input        Board
		intoSelf     []string
		intoWall     []string
		intoBadSpace []string
	}{
		{
			name: "avoid small space 1",
			input: Board{
				Snakes: []Battlesnake{
					{
						Head: Coord{2, 0},
						Body: []Coord{{2, 0}, {2, 1}, {1, 1}, {0, 1}, {0, 0}},
					},
				},
				Width:  7,
				Height: 7,
			},
			intoSelf:     []string{"up"},
			intoWall:     []string{"down"},
			intoBadSpace: []string{"left"},
		},
		{
			name: "avoid small space 2",
			input: Board{
				Snakes: []Battlesnake{
					{
						Health: 100,
						Head:   Coord{0, 5},
						Body:   []Coord{{0, 5}, {1, 5}, {2, 5}, {3, 5}, {3, 6}, {4, 6}, {5, 6}, {5, 5}, {4, 5}, {4, 4}, {3, 4}, {2, 4}, {1, 4}},
					},
				},
				Food:   []Coord{{4, 1}, {6, 6}},
				Width:  7,
				Height: 7,
			},
			intoSelf:     []string{"right"},
			intoWall:     []string{"left"},
			intoBadSpace: []string{"up"},
		},
	}

	for _, tc := range tests {
		state := GameState{
			Board: tc.input,
			You:   tc.input.Snakes[0],
		}

		nextMove := move(state)

		if contains(tc.intoSelf, nextMove.Move) {
			t.Errorf("%s: snake moved into self, %s", tc.name, nextMove.Move)
		}
		if contains(tc.intoWall, nextMove.Move) {
			t.Errorf("%s: snake moved into wall, %s", tc.name, nextMove.Move)
		}
		if contains(tc.intoBadSpace, nextMove.Move) {
			t.Errorf("%s: snake moved into small space with, %s", tc.name, nextMove.Move)
		}
	}
}

func TestAvoidBadHead2Head(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "3509d89e-8809-46c9-b46c-164158eaac26",
			Ruleset: Ruleset{
				Name: "standard",
			},
			Timeout: 500,
		},
		Turn: 3,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{0, 4}, {8, 2}, {4, 2}, {6, 10}, {5, 5}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_PgQqRchTF6k3FWhVk6fvkXv4",
					Name:    "nomblegomble",
					Health:  97,
					Head:    Coord{5, 6},
					Body:    []Coord{{5, 6}, {5, 7}, {5, 8}},
					Length:  3,
					Latency: "21",
				},
				{
					ID:      "gs_xw9mxDXPHD7fgFKgxBdMpJQT",
					Name:    "rnd",
					Health:  97,
					Head:    Coord{1, 8},
					Body:    []Coord{{1, 8}, {1, 7}, {1, 6}},
					Length:  3,
					Latency: "0",
				},
				{
					ID:      "gs_hVTbbK4dXGPJGCjxmKx9ktJf",
					Name:    "Mangofox",
					Health:  97,
					Head:    Coord{9, 4},
					Body:    []Coord{{9, 4}, {9, 3}, {9, 2}},
					Length:  3,
					Latency: "0",
				},
				{
					ID:      "gs_K4C6SmXkKSjWxwv3bY96djxD",
					Name:    "Steve",
					Health:  97,
					Head:    Coord{5, 4},
					Body:    []Coord{{5, 4}, {5, 3}, {5, 2}},
					Length:  3,
					Latency: "0",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_PgQqRchTF6k3FWhVk6fvkXv4",
			Name:    "nomblegomble",
			Health:  97,
			Head:    Coord{5, 6},
			Body:    []Coord{{5, 6}, {5, 7}, {5, 8}},
			Length:  3,
			Latency: "21",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "down" {
		t.Errorf("snake moved into a head2head that will kill both snakes, %s", nextMove.Move)
	}
}

func TestHead2HeadBetterThanWall(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "82cb8643-5f86-4674-a3ed-28c1d99a689f",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 80,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{4, 10}, {6, 3}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_3CWVKmKbMYvmS7qQYPDkm9f8",
					Name:    "nomblegomble",
					Health:  83,
					Head:    Coord{6, 10},
					Body:    []Coord{{6, 10}, {5, 10}, {5, 9}, {4, 9}, {4, 8}, {3, 8}, {3, 7}},
					Length:  7,
					Latency: "21",
				},
				{
					ID:      "gs_3D66hv63CXRMVRKDjCKDj8pJ",
					Name:    "nomblegomble",
					Health:  96,
					Head:    Coord{7, 9},
					Body:    []Coord{{7, 9}, {8, 9}, {9, 9}, {10, 9}, {10, 8}, {9, 8}, {8, 8}, {8, 7}},
					Length:  8,
					Latency: "21",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_3CWVKmKbMYvmS7qQYPDkm9f8",
			Name:    "nomblegomble",
			Health:  83,
			Head:    Coord{6, 10},
			Body:    []Coord{{6, 10}, {5, 10}, {5, 9}, {4, 9}, {4, 8}, {3, 8}, {3, 7}},
			Length:  7,
			Latency: "21",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake moved into wall, %s  (should prefer h2h chance)", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into self, %s  (should prefer h2h chance)", nextMove.Move)
	}
}

func TestKillerInstinct1(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "37719616-712a-4dea-9dbd-b9dfa992d908",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 80,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{9, 10}, {8, 9}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_x3Mybq4HHmYSHKbGh83fk9rJ",
					Name:    "nomblegomble",
					Health:  100,
					Head:    Coord{4, 6},
					Body:    []Coord{{4, 6}, {4, 5}, {4, 4}, {4, 3}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {1, 7}, {1, 7}},
					Length:  12,
					Latency: "22",
				},
				{
					ID:      "gs_D6yxdPV87SbYfrFYSDK7JVTR",
					Name:    "Scared Cobra Chicken",
					Health:  63,
					Head:    Coord{5, 7},
					Body:    []Coord{{5, 7}, {6, 7}, {6, 8}, {5, 8}, {5, 9}},
					Length:  5,
					Latency: "204",
				},
				{
					ID:      "gs_vJxJGRgK43X9G8DDBmK8CDSQ",
					Name:    "jsnek2",
					Health:  98,
					Head:    Coord{10, 4},
					Body:    []Coord{{10, 4}, {10, 3}, {10, 2}, {9, 2}, {9, 3}, {9, 4}, {8, 4}},
					Length:  7,
					Latency: "263",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_x3Mybq4HHmYSHKbGh83fk9rJ",
			Name:    "nomblegomble",
			Health:  100,
			Head:    Coord{4, 6},
			Body:    []Coord{{4, 6}, {4, 5}, {4, 4}, {4, 3}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {1, 7}, {1, 7}},
			Length:  12,
			Latency: "22",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake backed down from a kill, %s", nextMove.Move)
	}
}

func TestKillerInstinct2(t *testing.T) {

	state := GameState{
		Game: Game{
			ID: "37719616-712a-4dea-9dbd-b9dfa992d908",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 81,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{9, 10}, {8, 9}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_x3Mybq4HHmYSHKbGh83fk9rJ",
					Name:    "nomblegomble",
					Health:  99,
					Head:    Coord{3, 6},
					Body:    []Coord{{3, 6}, {4, 6}, {4, 5}, {4, 4}, {4, 3}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {1, 7}},
					Length:  12,
					Latency: "21",
				},
				{
					ID:      "gs_D6yxdPV87SbYfrFYSDK7JVTR",
					Name:    "Scared Cobra Chicken",
					Health:  62,
					Head:    Coord{4, 7},
					Body:    []Coord{{4, 7}, {5, 7}, {6, 7}, {6, 8}, {5, 8}},
					Length:  5,
					Latency: "205",
				},
				{
					ID:      "gs_vJxJGRgK43X9G8DDBmK8CDSQ",
					Name:    "jsnek2",
					Health:  97,
					Head:    Coord{10, 5},
					Body:    []Coord{{10, 5}, {10, 4}, {10, 3}, {10, 2}, {9, 2}, {9, 3}, {9, 4}},
					Length:  7,
					Latency: "254",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_x3Mybq4HHmYSHKbGh83fk9rJ",
			Name:    "nomblegomble",
			Health:  99,
			Head:    Coord{3, 6},
			Body:    []Coord{{3, 6}, {4, 6}, {4, 5}, {4, 4}, {4, 3}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {1, 7}},
			Length:  12,
			Latency: "21",
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

func TestKillerInstinct3(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "f53ef734-3349-467d-9eee-89b2d6f9b4fa",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 97,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{0, 2}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_YSGWKK73YPHrX3vdG3hJhGHT",
					Name:    "nomblegomble",
					Health:  95,
					Head:    Coord{1, 2},
					Body:    []Coord{{1, 2}, {1, 3}, {0, 3}, {0, 4}, {1, 4}, {2, 4}, {2, 3}, {3, 3}, {4, 3}, {4, 4}, {5, 4}, {6, 4}},
					Length:  12,
					Latency: "21",
				},
				{
					ID:      "gs_mmKppxxMFHWF73VYrjBM8RS8",
					Name:    "Ekans on a Plane",
					Health:  90,
					Head:    Coord{2, 1},
					Body:    []Coord{{2, 1}, {1, 1}, {0, 1}, {0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}},
					Length:  9,
					Latency: "73",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_YSGWKK73YPHrX3vdG3hJhGHT",
			Name:    "nomblegomble",
			Health:  95,
			Head:    Coord{1, 2},
			Body:    []Coord{{1, 2}, {1, 3}, {0, 3}, {0, 4}, {1, 4}, {2, 4}, {2, 3}, {3, 3}, {4, 3}, {4, 4}, {5, 4}, {6, 4}},
			Length:  12,
			Latency: "21",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "left" {
		t.Errorf("snake moved into too small of space, %s", nextMove.Move)
	}
	if nextMove.Move == "up" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "down" {
		t.Errorf("snake moved into other snake, %s", nextMove.Move)
	}
}

// This one is tricky because the Head2Head happens on the tail of the other snake
func TestKillerInstinctOtherTail4(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "e5652b90-b24e-43ff-ba46-c00ef8b1cb41",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 223,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{2, 0}, {0, 2}, {2, 10}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_jjjF7vvJC9tWc6dBM4MhjKrW",
					Name:    "nomblegomble",
					Health:  95,
					Head:    Coord{7, 2},
					Body:    []Coord{{7, 2}, {8, 2}, {8, 3}, {9, 3}, {9, 4}, {9, 5}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {9, 9}, {8, 9}, {7, 9}, {6, 9}, {6, 8}, {6, 7}, {5, 7}, {4, 7}, {4, 6}, {4, 5}, {5, 5}, {6, 5}, {7, 5}, {7, 6}, {7, 7}, {8, 7}, {9, 7}, {9, 6}, {8, 6}, {8, 5}, {8, 4}, {7, 4}, {7, 3}, {6, 3}, {5, 3}, {5, 2}, {4, 2}, {3, 2}, {2, 2}},
					Length:  42,
					Latency: "20",
				},
				{
					ID:      "gs_3F4d6gQq9ygjhykc6JpmmTmJ",
					Name:    "Eremetic Eric",
					Health:  74,
					Head:    Coord{6, 1},
					Body:    []Coord{{6, 1}, {5, 1}, {5, 0}, {6, 0}, {7, 0}, {7, 1}},
					Length:  6,
					Latency: "15",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_jjjF7vvJC9tWc6dBM4MhjKrW",
			Name:    "nomblegomble",
			Health:  95,
			Head:    Coord{7, 2},
			Body:    []Coord{{7, 2}, {8, 2}, {8, 3}, {9, 3}, {9, 4}, {9, 5}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {9, 9}, {8, 9}, {7, 9}, {6, 9}, {6, 8}, {6, 7}, {5, 7}, {4, 7}, {4, 6}, {4, 5}, {5, 5}, {6, 5}, {7, 5}, {7, 6}, {7, 7}, {8, 7}, {9, 7}, {9, 6}, {8, 6}, {8, 5}, {8, 4}, {7, 4}, {7, 3}, {6, 3}, {5, 3}, {5, 2}, {4, 2}, {3, 2}, {2, 2}},
			Length:  42,
			Latency: "20",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "left" {
		t.Errorf("snake moved into too small of space, %s", nextMove.Move)
	}
	if nextMove.Move == "up" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}

func TestFood(t *testing.T) {
	tests := []struct {
		name     string
		input    Battlesnake
		food     []Coord
		expected string
	}{
		{
			name: "eat when starving",
			input: Battlesnake{
				Head:   Coord{X: 5, Y: 5},
				Body:   []Coord{{X: 5, Y: 5}, {X: 5, Y: 4}, {X: 6, Y: 4}},
				Health: 1,
			},
			food:     []Coord{{X: 6, Y: 5}},
			expected: "right",
		},
		{
			name: "go towards food when hungry",
			input: Battlesnake{
				Head:   Coord{X: 5, Y: 5},
				Body:   []Coord{{X: 5, Y: 5}, {X: 5, Y: 4}, {X: 6, Y: 4}},
				Health: 20,
			},
			food:     []Coord{{X: 0, Y: 5}},
			expected: "left",
		},
	}

	for _, tc := range tests {
		state := GameState{
			Board: Board{
				Width:  12,
				Height: 12,
				Snakes: []Battlesnake{tc.input},
				Food:   tc.food,
			},
			You: tc.input,
		}

		nextMove := move(state)

		if nextMove.Move != tc.expected {
			t.Errorf("%s: expected %s, got %s", tc.name, tc.expected, nextMove.Move)
		}
	}
}

func TestMath(t *testing.T) {
	tests := []struct {
		head     Coord
		food     Coord
		expected int
	}{
		{
			head:     Coord{X: 0, Y: 0},
			food:     Coord{X: 2, Y: 2},
			expected: 4,
		},
		{
			head:     Coord{X: 5, Y: 5},
			food:     Coord{X: 7, Y: 3},
			expected: 4,
		},
	}

	for _, tc := range tests {

		actual := distance(tc.head, tc.food)

		if actual != tc.expected {
			t.Errorf("expected %d, got %d", tc.expected, actual)
		}

	}
}

func TestCombineWeights(t *testing.T) {
	tests := []struct {
		scores   []WeightedScore
		expected Scored
	}{
		{
			scores: []WeightedScore{
				{true, 1.0, Scored{"up": 1.0, "down": 1.0, "left": 1.0, "right": 1.0}},
			},
			expected: Scored{"up": 1.0, "down": 1.0, "left": 1.0, "right": 1.0},
		},
		{
			scores: []WeightedScore{
				{true, 1.0, Scored{"up": 1.0, "down": 1.0, "left": 1.0, "right": 1.0}},
				{true, 1.0, Scored{"up": 0.0, "down": 0.5, "left": 1.0, "right": 1.0}},
			},
			expected: Scored{"up": 0.0, "down": 0.5, "left": 1.0, "right": 1.0},
		},
		{
			scores: []WeightedScore{
				{true, 1.0, Scored{"up": 0.0, "down": 1.0, "left": 1.0, "right": 1.0}},
				{true, 1.0, Scored{"up": 0.5, "down": 0.01, "left": 0.5, "right": 0.5}},
				{false, 1.0, Scored{"up": 0.0, "down": 1.0, "left": 0.8, "right": 0.8}},
			},
			expected: Scored{"up": 0.0, "down": 0.02, "left": 0.9, "right": 0.9},
		},
	}

	for _, tc := range tests {

		actual := combineMoves(tc.scores)

		for move, score := range actual {
			if actual[move] != tc.expected[move] {
				t.Errorf("%s: expected %.2f, got %.2f", move, tc.expected[move], score)
			}
		}

	}
}

func TestRemap(t *testing.T) {
	tests := []struct {
		oldValue float64
		oldMin   float64
		oldMax   float64
		newMin   float64
		newMax   float64
		expected float64
	}{
		{
			oldValue: 1,
			oldMin:   0,
			oldMax:   2,
			newMin:   0,
			newMax:   1,
			expected: 0.5,
		},
		{
			oldValue: 0,
			oldMin:   0,
			oldMax:   0,
			newMin:   0,
			newMax:   1,
			expected: 0,
		},
	}

	for _, tc := range tests {

		actual := remap(tc.oldValue, tc.oldMin, tc.oldMax, tc.newMin, tc.newMax)

		if actual != tc.expected {
			t.Errorf("expected %.4f, got %.4f", tc.expected, actual)
		}

	}
}
