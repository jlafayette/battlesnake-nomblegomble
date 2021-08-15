package main

import (
	"encoding/json"
	"log"
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
				Width:  7,
				Height: 7,
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

func TestSpaceBasic(t *testing.T) {
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
				Food:   []Coord{{1, 0}, {0, 0}}, // 0,0 simulates bad luck of food spawning
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

func TestAvoidBadHead2Head2(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "294cee29-202c-4fe4-9482-db64cf19fad6",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 41,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{6, 7}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_BgWvHQ7yGrhmrWkrfqdGHQrc",
					Name:    "nomblegomble",
					Health:  71,
					Head:    Coord{0, 5},
					Body:    []Coord{{0, 5}, {1, 5}, {1, 6}, {2, 6}, {3, 6}},
					Length:  5,
					Latency: "23",
					Shout:   "",
				},
				{
					ID:      "gs_7S76jFcGSmVCrKwydvth4fJ7",
					Name:    "bsnek2",
					Health:  78,
					Head:    Coord{8, 3},
					Body:    []Coord{{8, 3}, {8, 2}, {8, 1}, {9, 1}, {9, 2}, {10, 2}},
					Length:  6,
					Latency: "190",
					Shout:   "",
				},
				{
					ID:      "gs_ch43JRySwjgX3MyW8hr9HMkC",
					Name:    "Boomslang",
					Health:  99,
					Head:    Coord{2, 3},
					Body:    []Coord{{2, 3}, {3, 3}, {3, 2}, {3, 1}, {3, 0}, {4, 0}, {5, 0}, {6, 0}},
					Length:  8,
					Latency: "253",
					Shout:   "",
				},
				{
					ID:      "gs_RRmMymt6W9bwjHh9mtrfvPt3",
					Name:    "Crimson",
					Health:  83,
					Head:    Coord{5, 4},
					Body:    []Coord{{5, 4}, {5, 5}, {6, 5}, {6, 4}, {6, 3}, {5, 3}},
					Length:  6,
					Latency: "198",
					Shout:   "",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_BgWvHQ7yGrhmrWkrfqdGHQrc",
			Name:    "nomblegomble",
			Health:  71,
			Head:    Coord{0, 5},
			Body:    []Coord{{0, 5}, {1, 5}, {1, 6}, {2, 6}, {3, 6}},
			Length:  5,
			Latency: "23",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved into head2head zone, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into wall, %s", nextMove.Move)
	}
}

func TestAvoidBadHead2Head2_diff(t *testing.T) {
	data := `{"game":{"id":"294cee29-202c-4fe4-9482-db64cf19fad6","ruleset":{"name":"standard","version":"v1.0.17"},"timeout":500},"turn":41,"board":{"height":11,"width":11,"food":[{"x":6,"y":7}],"snakes":[{"id":"gs_7S76jFcGSmVCrKwydvth4fJ7","name":"bsnek2","health":78,"body":[{"x":8,"y":3},{"x":8,"y":2},{"x":8,"y":1},{"x":9,"y":1},{"x":9,"y":2},{"x":10,"y":2}],"head":{"x":8,"y":3},"length":6,"latency":"190","shout":"","squad":""},{"id":"gs_ch43JRySwjgX3MyW8hr9HMkC","name":"Boomslang","health":99,"body":[{"x":2,"y":3},{"x":3,"y":3},{"x":3,"y":2},{"x":3,"y":1},{"x":3,"y":0},{"x":4,"y":0},{"x":5,"y":0},{"x":6,"y":0}],"head":{"x":2,"y":3},"length":8,"latency":"253","shout":"","squad":""},{"id":"gs_BgWvHQ7yGrhmrWkrfqdGHQrc","name":"nomblegomble","health":71,"body":[{"x":0,"y":5},{"x":1,"y":5},{"x":1,"y":6},{"x":2,"y":6},{"x":3,"y":6}],"head":{"x":0,"y":5},"length":5,"latency":"23","shout":"","squad":""},{"id":"gs_RRmMymt6W9bwjHh9mtrfvPt3","name":"Crimson","health":83,"body":[{"x":5,"y":4},{"x":5,"y":5},{"x":6,"y":5},{"x":6,"y":4},{"x":6,"y":3},{"x":5,"y":3}],"head":{"x":5,"y":4},"length":6,"latency":"198","shout":"","squad":""}],"hazards":[]},"you":{"id":"gs_BgWvHQ7yGrhmrWkrfqdGHQrc","name":"nomblegomble","health":71,"body":[{"x":0,"y":5},{"x":1,"y":5},{"x":1,"y":6},{"x":2,"y":6},{"x":3,"y":6}],"head":{"x":0,"y":5},"length":5,"latency":"23","shout":"","squad":""}}`
	// {"move":"down"}
	state := GameState{}
	err := json.Unmarshal([]byte(data), &state)
	if err != nil {
		log.Fatalln(err)
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved into head2head zone, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into wall, %s", nextMove.Move)
	}
}

// Other snakes are likely to go for food, so don't go for it if you don't have to.
func TestAvoidFoodInEqualHead2Head1(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "f2ce457a-7eb1-43cf-b495-831498b753e0",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 54,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{2, 7}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_qW3xGphphfBRBPjfqt7phy6X",
					Name:    "TC5001",
					Health:  98,
					Head:    Coord{8, 2},
					Body:    []Coord{{8, 2}, {9, 2}, {10, 2}, {10, 1}, {10, 0}, {9, 0}, {8, 0}, {7, 0}, {7, 1}, {6, 1}},
					Length:  10,
					Latency: "215",
					Shout:   "",
				},
				{
					ID:      "gs_97jPPpfDBqHYhxkkrdHHrQTY",
					Name:    "nomblegomble",
					Health:  98,
					Head:    Coord{3, 7},
					Body:    []Coord{{3, 7}, {3, 8}, {4, 8}, {5, 8}, {5, 7}, {5, 6}, {4, 6}, {3, 6}},
					Length:  8,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_7w9rQ6VDrSpYpgRCjygk8mmF",
					Name:    "bsnek2",
					Health:  99,
					Head:    Coord{9, 9},
					Body:    []Coord{{9, 9}, {9, 8}, {9, 7}, {9, 6}, {9, 5}, {8, 5}, {8, 4}},
					Length:  7,
					Latency: "284",
					Shout:   "",
				},
				{
					ID:      "gs_RmhhbkKtHSff73R66GMT33gJ",
					Name:    "SnakeJS",
					Health:  69,
					Head:    Coord{2, 6},
					Body:    []Coord{{2, 6}, {2, 5}, {3, 5}, {3, 4}, {3, 3}, {3, 2}, {3, 1}, {4, 1}},
					Length:  8,
					Latency: "239",
					Shout:   "[2 5] --> [2 7]",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_97jPPpfDBqHYhxkkrdHHrQTY",
			Name:    "nomblegomble",
			Health:  98,
			Head:    Coord{3, 7},
			Body:    []Coord{{3, 7}, {3, 8}, {4, 8}, {5, 8}, {5, 7}, {5, 6}, {4, 6}, {3, 6}},
			Length:  8,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "left" {
		t.Errorf("snake moved into likely H2H on food, %s", nextMove.Move)
	}
	if nextMove.Move == "up" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into too small a space, %s", nextMove.Move)
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

// seeing different results on live game from the server
// is there different data?
// need to log the actual data recieved on the server for comparison
func TestSpaceAvoidTiedHead2Head1(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "19ec09a8-0cea-40c5-a50b-f85f6cd0f400",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 6,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{5, 5}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_dF743rmHfbQ6YqkCbDSXmfj7",
					Name:    "nomblegomble",
					Health:  96,
					Head:    Coord{0, 8},
					Body:    []Coord{{0, 8}, {1, 8}, {1, 7}, {0, 7}},
					Length:  4,
					Latency: "23",
					Shout:   "",
				},
				{
					ID:      "gs_GhqqhY7mhc3vbmXtRyVGkwMR",
					Name:    "Fairy Rust",
					Health:  96,
					Head:    Coord{0, 2},
					Body:    []Coord{{0, 2}, {1, 2}, {1, 1}, {1, 0}},
					Length:  4,
					Latency: "45",
					Shout:   "",
				},
				{
					ID:      "gs_fHQbrqVXFYQcvDk74wmmfv9Q",
					Name:    "Legless Lizard",
					Health:  96,
					Head:    Coord{4, 2},
					Body:    []Coord{{4, 2}, {3, 2}, {3, 1}, {3, 0}},
					Length:  4,
					Latency: "51",
					Shout:   "",
				},
				{
					ID:      "gs_rj6ktV78gbfyqKdGB3ffXJF9",
					Name:    "The Jabberwock",
					Health:  98,
					Head:    Coord{1, 5},
					Body:    []Coord{{1, 5}, {1, 4}, {0, 4}, {0, 5}},
					Length:  4,
					Latency: "4",
					Shout:   "From hellΓÇÖs heart I stab at thee",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_dF743rmHfbQ6YqkCbDSXmfj7",
			Name:    "nomblegomble",
			Health:  96,
			Head:    Coord{0, 8},
			Body:    []Coord{{0, 8}, {1, 8}, {1, 7}, {0, 7}},
			Length:  4,
			Latency: "23",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved into potential tied head2head, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into wall, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
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

func TestSpaceCornered1(t *testing.T) {

	state := GameState{
		Game: Game{
			ID: "7560784f-f380-427d-8350-80725b25207a",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 68,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{6, 2}, {0, 2}, {6, 6}, {10, 3}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_bVHWqPM7PxRHkqCjGTgSkjhY",
					Name:    "nomblegomble",
					Health:  96,
					Head:    Coord{8, 2},
					Body:    []Coord{{8, 2}, {7, 2}, {7, 1}, {7, 0}, {6, 0}, {5, 0}, {4, 0}, {4, 1}, {4, 2}, {5, 2}},
					Length:  10,
					Latency: "22",
				},
				{
					ID:      "gs_XPtjRWfB3VT7GXFHGkvhXXbd",
					Name:    "MAsterStudentSlayer666",
					Health:  85,
					Head:    Coord{3, 3},
					Body:    []Coord{{3, 3}, {2, 3}, {2, 4}, {1, 4}, {0, 4}, {0, 5}, {0, 6}},
					Length:  7,
					Latency: "219",
				},
				{
					ID:      "gs_XyrrYqC3pg8W4DbSpF3jpPWc",
					Name:    "leshchenko-1",
					Health:  100,
					Head:    Coord{9, 1},
					Body:    []Coord{{9, 1}, {9, 2}, {9, 3}, {9, 4}, {9, 4}},
					Length:  5,
					Latency: "234",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_bVHWqPM7PxRHkqCjGTgSkjhY",
			Name:    "nomblegomble",
			Health:  96,
			Head:    Coord{8, 2},
			Body:    []Coord{{8, 2}, {7, 2}, {7, 1}, {7, 0}, {6, 0}, {5, 0}, {4, 0}, {4, 1}, {4, 2}, {5, 2}},
			Length:  10,
			Latency: "22",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into other snake, %s", nextMove.Move)
	}
}

func TestSpaceNoFoodIfLessThanBodyLen1(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "3988391b-ee86-466e-ab0c-d39c38283d38",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 126,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{6, 1}, {10, 4}, {0, 9}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_X6DRymbGqtcDBWxfTqmkBhjB",
					Name:    "snek",
					Health:  5,
					Head:    Coord{1, 5},
					Body:    []Coord{{1, 5}, {1, 6}, {2, 6}, {3, 6}, {3, 5}},
					Length:  5,
					Latency: "281",
					Shout:   "",
				},
				{
					ID:      "gs_TfjMPmMkjSWmf4dYWjX7rrjK",
					Name:    "msbs",
					Health:  69,
					Head:    Coord{3, 1},
					Body:    []Coord{{3, 1}, {3, 2}, {2, 2}, {2, 1}, {1, 1}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}},
					Length:  12,
					Latency: "53",
					Shout:   "",
				},
				{
					ID:      "gs_PXMmW3xbSxDBJCg4hycg7xyG",
					Name:    "nomblegomble",
					Health:  86,
					Head:    Coord{0, 8},
					Body:    []Coord{{0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {4, 9}, {4, 10}, {5, 10}, {6, 10}, {7, 10}, {7, 9}, {6, 9}, {5, 9}, {5, 8}, {6, 8}, {7, 8}, {7, 7}},
					Length:  17,
					Latency: "22",
					Shout:   "",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_PXMmW3xbSxDBJCg4hycg7xyG",
			Name:    "nomblegomble",
			Health:  86,
			Head:    Coord{0, 8},
			Body:    []Coord{{0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {4, 9}, {4, 10}, {5, 10}, {6, 10}, {7, 10}, {7, 9}, {6, 9}, {5, 9}, {5, 8}, {6, 8}, {7, 8}, {7, 7}},
			Length:  17,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake moved into too small of space, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into wall, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}

func TestSpaceCutoff2(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "5ff70484-ac66-4025-90b6-9af1554b74b5",
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
			Food:   []Coord{{10, 0}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_PVdSm9cYDRMk3R6Tk3Qqpw64",
					Name:    "nomblegomble",
					Health:  99,
					Head:    Coord{10, 1},
					Body:    []Coord{{10, 1}, {9, 1}, {8, 1}, {7, 1}, {7, 2}, {7, 3}, {6, 3}, {6, 4}, {6, 5}, {6, 6}, {6, 7}, {6, 8}, {6, 9}, {5, 9}},
					Length:  14,
					Latency: "21",
				},
				{
					ID:      "gs_WJQ63xb7t6mpFCXqHVdPqBr4",
					Name:    "Ifarus",
					Health:  77,
					Head:    Coord{3, 0},
					Body:    []Coord{{3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}},
					Length:  8,
					Latency: "76",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_PVdSm9cYDRMk3R6Tk3Qqpw64",
			Name:    "nomblegomble",
			Health:  99,
			Head:    Coord{10, 1},
			Body:    []Coord{{10, 1}, {9, 1}, {8, 1}, {7, 1}, {7, 2}, {7, 3}, {6, 3}, {6, 4}, {6, 5}, {6, 6}, {6, 7}, {6, 8}, {6, 9}, {5, 9}},
			Length:  14,
			Latency: "21",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into wall, %s", nextMove.Move)
	}
}

func TestSpaceCutoff3(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "3c3b7dcc-4f7d-48d2-9449-ee22bda84390",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 32,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{8, 2}, {9, 4}, {4, 1}, {3, 4}, {9, 7}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_YwVHBvWKTXXKKjXvM6y93VKH",
					Name:    "nomblegomble",
					Health:  83,
					Head:    Coord{8, 4},
					Body:    []Coord{{8, 4}, {8, 3}, {7, 3}, {7, 2}, {7, 1}},
					Length:  5,
					Latency: "22",
				},
				{
					ID:      "gs_3W9Dm9F4Hw73fXgxPPmVRwFX",
					Name:    "nates_python",
					Health:  98,
					Head:    Coord{5, 9},
					Body:    []Coord{{5, 9}, {6, 9}, {7, 9}, {7, 8}, {7, 7}},
					Length:  5,
					Latency: "220",
				},
				{
					ID:      "gs_Bgj94MjGM8c7Mqppbhvjcx3K",
					Name:    "carl",
					Health:  86,
					Head:    Coord{5, 3},
					Body:    []Coord{{5, 3}, {4, 3}, {3, 3}, {3, 2}, {3, 1}, {2, 1}},
					Length:  6,
					Latency: "250",
				},
				{
					ID:      "gs_pgQRpdYWW7cxQbhMpbwxwC84",
					Name:    "Morley",
					Health:  68,
					Head:    Coord{7, 5},
					Body:    []Coord{{7, 5}, {8, 5}, {9, 5}},
					Length:  3,
					Latency: "77",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_YwVHBvWKTXXKKjXvM6y93VKH",
			Name:    "nomblegomble",
			Health:  83,
			Head:    Coord{8, 4},
			Body:    []Coord{{8, 4}, {8, 3}, {7, 3}, {7, 2}, {7, 1}},
			Length:  5,
			Latency: "22",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "left" {
		t.Errorf("snake moved into fatal H2H, %s (can be cut off)", nextMove.Move)
	}
	if nextMove.Move == "down" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "up" {
		t.Errorf("snake moved into other snake, %s", nextMove.Move)
	}
}

func TestSpaceCutoff4H2HWeaker(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "8cb97ac0-f405-41a1-b007-a9a4b53bbbfa",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 96,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{7, 0}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_tmWR3BFhBHMPhHYM37rqdB37",
					Name:    "nomblegomble",
					Health:  100,
					Head:    Coord{8, 10},
					Body:    []Coord{{8, 10}, {8, 9}, {8, 8}, {8, 7}, {7, 7}, {7, 8}, {6, 8}, {5, 8}, {5, 9}, {4, 9}, {3, 9}, {2, 9}, {2, 9}},
					Length:  13,
					Latency: "21",
				},
				{
					ID:      "gs_XJ44wjQRyT3MPqqwTB8WKmpX",
					Name:    "Ouroboros 2",
					Health:  97,
					Head:    Coord{9, 7},
					Body:    []Coord{{9, 7}, {9, 6}, {8, 6}, {7, 6}, {6, 6}, {5, 6}, {4, 6}, {4, 5}, {3, 5}, {3, 4}, {3, 3}, {3, 2}, {3, 1}},
					Length:  13,
					Latency: "214",
				},
				{
					ID:      "gs_6HpdTvFPJJHk8KKxXGSq4vdb",
					Name:    "Canadian Bacon",
					Health:  6,
					Head:    Coord{2, 10},
					Body:    []Coord{{2, 10}, {1, 10}, {1, 9}, {1, 8}},
					Length:  4,
					Latency: "197",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_tmWR3BFhBHMPhHYM37rqdB37",
			Name:    "nomblegomble",
			Health:  100,
			Head:    Coord{8, 10},
			Body:    []Coord{{8, 10}, {8, 9}, {8, 8}, {8, 7}, {7, 7}, {7, 8}, {6, 8}, {5, 8}, {5, 9}, {4, 9}, {3, 9}, {2, 9}, {2, 9}},
			Length:  13,
			Latency: "21",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "right" {
		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
	}
	if nextMove.Move == "down" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "up" {
		t.Errorf("snake moved into wall, %s", nextMove.Move)
	}
}

func TestSpaceCutoff5(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "4ac4ba8c-4c68-4a03-9607-583264860222",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 77,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{0, 0}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_PQYck9Y4W3MQpS3HPWrJ8HMd",
					Name:    "nomblegomble",
					Health:  96,
					Head:    Coord{0, 3},
					Body:    []Coord{{0, 3}, {1, 3}, {1, 2}, {1, 1}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}},
					Length:  9,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_8TvRDTjqhcJCb6w8YhjSXHcH",
					Name:    "Cool_as_ice",
					Health:  76,
					Head:    Coord{7, 6},
					Body:    []Coord{{7, 6}, {6, 6}, {6, 5}, {6, 4}, {6, 3}, {7, 3}, {7, 4}, {7, 5}},
					Length:  8,
					Latency: "74",
					Shout:   "",
				},
				{
					ID:      "gs_x88cHrQDvKxJHkTwC7TSWk3W",
					Name:    "moon-snake-pika",
					Health:  99,
					Head:    Coord{9, 8},
					Body:    []Coord{{9, 8}, {10, 8}, {10, 7}, {10, 6}, {10, 5}, {10, 4}, {10, 3}, {10, 2}, {10, 1}, {10, 0}, {9, 0}, {9, 1}, {9, 2}, {9, 3}},
					Length:  14,
					Latency: "217",
					Shout:   "",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_PQYck9Y4W3MQpS3HPWrJ8HMd",
			Name:    "nomblegomble",
			Health:  96,
			Head:    Coord{0, 3},
			Body:    []Coord{{0, 3}, {1, 3}, {1, 2}, {1, 1}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}},
			Length:  9,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved into too small of space, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into wall, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}

func TestSpaceCutoff6(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "e74d6f1d-a38a-4135-bb89-d17f387ba9ae",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 108,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{1, 9}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_cYKtRkqxfMfyj46WH9J9VhmX",
					Name:    "nomblegomble",
					Health:  90,
					Head:    Coord{6, 4},
					Body:    []Coord{{6, 4}, {7, 4}, {7, 3}, {7, 2}, {7, 1}, {6, 1}, {5, 1}, {4, 1}, {3, 1}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {2, 2}},
					Length:  17,
					Latency: "23",
					Shout:   "",
				},
				{
					ID:      "gs_M7h8gyWqTGRFg8GX4BJf3whb",
					Name:    "Ophiophagus One",
					Health:  100,
					Head:    Coord{4, 4},
					Body:    []Coord{{4, 4}, {4, 5}, {5, 5}, {5, 6}, {5, 7}, {6, 7}, {6, 7}},
					Length:  7,
					Latency: "211",
					Shout:   "",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_cYKtRkqxfMfyj46WH9J9VhmX",
			Name:    "nomblegomble",
			Health:  90,
			Head:    Coord{6, 4},
			Body:    []Coord{{6, 4}, {7, 4}, {7, 3}, {7, 2}, {7, 1}, {6, 1}, {5, 1}, {4, 1}, {3, 1}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {2, 2}},
			Length:  17,
			Latency: "23",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "left" {
		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}

func TestSpaceOkToTailChase1(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "71116f92-59d0-4f88-a578-a75035b4c1be",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 224,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{8, 2}, {7, 9}, {0, 8}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_v8SbjcXHGrxjwjmwpmWPy67Q",
					Name:    "nomblegomble",
					Health:  98,
					Head:    Coord{3, 1},
					Body:    []Coord{{3, 1}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 4}, {1, 3}, {2, 3}, {2, 2}, {3, 2}},
					Length:  14,
					Latency: "21",
					Shout:   "",
				},
				{
					ID:      "gs_bVJmK8pX8mkM7G4wVyjKCxtT",
					Name:    "caicai-vilu",
					Health:  99,
					Head:    Coord{7, 1},
					Body:    []Coord{{7, 1}, {7, 0}, {6, 0}, {5, 0}, {5, 1}, {5, 2}, {4, 2}, {4, 3}, {3, 3}, {3, 4}, {3, 5}, {2, 5}, {1, 5}, {1, 6}, {1, 7}, {1, 8}, {1, 9}, {2, 9}, {3, 9}, {4, 9}, {5, 9}, {6, 9}, {6, 8}, {5, 8}, {4, 8}, {3, 8}, {2, 8}, {2, 7}, {3, 7}, {4, 7}, {5, 7}, {6, 7}, {7, 7}, {8, 7}, {9, 7}},
					Length:  35,
					Latency: "75",
					Shout:   "",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_v8SbjcXHGrxjwjmwpmWPy67Q",
			Name:    "nomblegomble",
			Health:  98,
			Head:    Coord{3, 1},
			Body:    []Coord{{3, 1}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 4}, {1, 3}, {2, 3}, {2, 2}, {3, 2}},
			Length:  14,
			Latency: "21",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "right" {
		t.Errorf("snake moved into too small of space, %s", nextMove.Move)
	}
	if nextMove.Move == "down" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}

func TestFoodStart0(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "4c46aa82-936c-46c6-aeb2-6e33da287a3b",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 0,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{0, 0}, {8, 4}, {10, 10}, {2, 10}, {5, 5}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_r4JCVS8Hbjq87Cg3BQM37HPf",
					Name:    "nomblegomble",
					Health:  100,
					Head:    Coord{1, 1},
					Body:    []Coord{{1, 1}, {1, 1}, {1, 1}},
					Length:  3,
					Latency: "",
				},
				{
					ID:      "gs_VG6tp7kfmSXSkQyPHKjk3vC6",
					Name:    "DDT",
					Health:  100,
					Head:    Coord{9, 5},
					Body:    []Coord{{9, 5}, {9, 5}, {9, 5}},
					Length:  3,
					Latency: "",
				},
				{
					ID:      "gs_VMrTMQrtfRrYdRPbdqbJYphd",
					Name:    "Yung Snek V0",
					Health:  100,
					Head:    Coord{9, 9},
					Body:    []Coord{{9, 9}, {9, 9}, {9, 9}},
					Length:  3,
					Latency: "",
				},
				{
					ID:      "gs_BChKFRcw7qVwTmfCQYfSgB4P",
					Name:    "Leonardo",
					Health:  100,
					Head:    Coord{1, 9},
					Body:    []Coord{{1, 9}, {1, 9}, {1, 9}},
					Length:  3,
					Latency: "",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_r4JCVS8Hbjq87Cg3BQM37HPf",
			Name:    "nomblegomble",
			Health:  100,
			Head:    Coord{1, 1},
			Body:    []Coord{{1, 1}, {1, 1}, {1, 1}},
			Length:  3,
			Latency: "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake moved away from starting food, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved away from starting food, %s", nextMove.Move)
	}
}

func TestFoodStart1(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "4c46aa82-936c-46c6-aeb2-6e33da287a3b",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 1,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{0, 0}, {8, 4}, {10, 10}, {2, 10}, {5, 5}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_r4JCVS8Hbjq87Cg3BQM37HPf",
					Name:    "nomblegomble",
					Health:  99,
					Head:    Coord{1, 0},
					Body:    []Coord{{1, 0}, {1, 1}, {1, 1}},
					Length:  3,
					Latency: "46",
				},
				{
					ID:      "gs_VG6tp7kfmSXSkQyPHKjk3vC6",
					Name:    "DDT",
					Health:  99,
					Head:    Coord{9, 4},
					Body:    []Coord{{9, 4}, {9, 5}, {9, 5}},
					Length:  3,
					Latency: "293",
				},
				{
					ID:      "gs_VMrTMQrtfRrYdRPbdqbJYphd",
					Name:    "Yung Snek V0",
					Health:  99,
					Head:    Coord{8, 9},
					Body:    []Coord{{8, 9}, {9, 9}, {9, 9}},
					Length:  3,
					Latency: "282",
				},
				{
					ID:      "gs_BChKFRcw7qVwTmfCQYfSgB4P",
					Name:    "Leonardo",
					Health:  99,
					Head:    Coord{1, 10},
					Body:    []Coord{{1, 10}, {1, 9}, {1, 9}},
					Length:  3,
					Latency: "288",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_r4JCVS8Hbjq87Cg3BQM37HPf",
			Name:    "nomblegomble",
			Health:  99,
			Head:    Coord{1, 0},
			Body:    []Coord{{1, 0}, {1, 1}, {1, 1}},
			Length:  3,
			Latency: "46",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake moved away from starting food, %s", nextMove.Move)
	}
	if nextMove.Move == "down" {
		t.Errorf("snake into wall, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved away from starting food, %s", nextMove.Move)
	}
}

func TestFood3(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "245970ce-0424-4a9f-a02b-a1f0d5f531a1",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 5,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{5, 5}, {9, 2}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_SWrrBGBdPF3TCy7qJXcSkrwP",
					Name:    "nomblegomble",
					Health:  97,
					Head:    Coord{8, 1},
					Body:    []Coord{{8, 1}, {7, 1}, {7, 0}, {8, 0}},
					Length:  4,
					Latency: "25",
				},
				{
					ID:      "gs_Tp8qw8GdDpJrxFjcj7qBVkkb",
					Name:    "nomblegomble",
					Health:  97,
					Head:    Coord{1, 4},
					Body:    []Coord{{1, 4}, {1, 5}, {0, 5}, {0, 4}},
					Length:  4,
					Latency: "48",
				},
				{
					ID:      "gs_CBQM9J66qjrbcYk44YSwKCRY",
					Name:    "nomblegomble",
					Health:  97,
					Head:    Coord{8, 5},
					Body:    []Coord{{8, 5}, {7, 5}, {7, 6}, {8, 6}},
					Length:  4,
					Latency: "45",
				},
				{
					ID:      "gs_vhSmG8q4P4H4dRmpgc3xSgHS",
					Name:    "nomblegomble",
					Health:  97,
					Head:    Coord{5, 8},
					Body:    []Coord{{5, 8}, {5, 9}, {5, 10}, {6, 10}},
					Length:  4,
					Latency: "23",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_SWrrBGBdPF3TCy7qJXcSkrwP",
			Name:    "nomblegomble",
			Health:  97,
			Head:    Coord{8, 1},
			Body:    []Coord{{8, 1}, {7, 1}, {7, 0}, {8, 0}},
			Length:  4,
			Latency: "25",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved away from nearby food, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
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
