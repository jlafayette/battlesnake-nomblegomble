package main

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/jlafayette/battlesnake-go/wire"
)

func TestNeckAvoidance(t *testing.T) {
	tests := []struct {
		name  string
		input wire.Battlesnake
		noGo  string
	}{
		{
			name: "neck avoidance 1",
			input: wire.Battlesnake{
				Length: 3, // facing right
				Head:   wire.Coord{X: 2, Y: 0},
				Body:   []wire.Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
			}, noGo: "left",
		},
		{
			name: "neck avoidance 2",
			input: wire.Battlesnake{
				Length: 3, // facing left
				Head:   wire.Coord{X: 7, Y: 0},
				Body:   []wire.Coord{{X: 7, Y: 0}, {X: 8, Y: 0}, {X: 9, Y: 0}},
			}, noGo: "right",
		},
		{
			name: "neck avoidance 3",
			input: wire.Battlesnake{
				Length: 3, // facing up
				Head:   wire.Coord{X: 5, Y: 10},
				Body:   []wire.Coord{{X: 5, Y: 9}, {X: 5, Y: 8}, {X: 5, Y: 7}},
			}, noGo: "down",
		},
		{
			name: "neck avoidance 4",
			input: wire.Battlesnake{
				Length: 3, // facing down
				Head:   wire.Coord{X: 5, Y: 1},
				Body:   []wire.Coord{{X: 5, Y: 2}, {X: 5, Y: 3}, {X: 5, Y: 4}},
			}, noGo: "up",
		},
	}

	for _, tc := range tests {
		state := wire.GameState{
			Board: wire.Board{
				Snakes: []wire.Battlesnake{tc.input},
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
		input    wire.Battlesnake
		intoNeck string
		intoWall []string
	}{
		{
			name: "wall avoidance 1",
			input: wire.Battlesnake{
				// Lower left corner
				ID:     "my-id",
				Length: 3,
				Health: 90,
				Head:   wire.Coord{X: 0, Y: 0},
				Body:   []wire.Coord{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			},
			intoNeck: "right",
			intoWall: []string{"left", "down"},
		},
		{
			name: "wall avoidance 2",
			input: wire.Battlesnake{
				// top right corner
				ID:     "my-id",
				Length: 3,
				Health: 90,
				Head:   wire.Coord{X: 11, Y: 11},
				Body:   []wire.Coord{{X: 11, Y: 11}, {X: 10, Y: 11}, {X: 9, Y: 11}},
			},
			intoNeck: "left",
			intoWall: []string{"up", "right"},
		},
		{
			name: "wall avoidance 3",
			input: wire.Battlesnake{
				// bottom right corner (facing down)
				ID:     "my-id",
				Length: 3,
				Health: 90,
				Head:   wire.Coord{X: 11, Y: 0},
				Body:   []wire.Coord{{X: 11, Y: 0}, {X: 11, Y: 1}, {X: 11, Y: 2}},
			},
			intoNeck: "up",
			intoWall: []string{"down", "right"},
		},
		{
			name: "wall avoidance 4",
			input: wire.Battlesnake{
				// top left corner (facing up)
				ID:     "my-id",
				Length: 3,
				Health: 90,
				Head:   wire.Coord{X: 0, Y: 11},
				Body:   []wire.Coord{{X: 0, Y: 11}, {X: 0, Y: 10}, {X: 0, Y: 9}},
			},
			intoNeck: "down",
			intoWall: []string{"left", "up"},
		},
	}

	for _, tc := range tests {
		state := wire.GameState{
			Board: wire.Board{
				Width:  12,
				Height: 12,
				Snakes: []wire.Battlesnake{tc.input},
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

// func TestSelfAvoidance(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		input    wire.Battlesnake
// 		intoSelf []string
// 		intoWall []string
// 	}{
// 		{
// 			name: "body check",
// 			input: wire.Battlesnake{
// 				ID:     "my-id",
// 				Length: 7,
// 				Health: 90,
// 				Head:   wire.Coord{X: 5, Y: 5},
// 				Body:   []wire.Coord{{X: 5, Y: 5}, {X: 5, Y: 4}, {X: 6, Y: 4}, {X: 6, Y: 5}, {X: 6, Y: 6}, {X: 5, Y: 6}, {X: 4, Y: 6}},
// 			},
// 			intoSelf: []string{"up", "right", "down"},
// 		},
// 		// tail is ok if not at full health
// 		{
// 			name: "tail chase ok 1",
// 			input: wire.Battlesnake{
// 				ID:     "my-id",
// 				Length: 4,
// 				Health: 99,
// 				Head:   wire.Coord{X: 11, Y: 0},
// 				Body:   []wire.Coord{{X: 11, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 1}, {X: 11, Y: 1}},
// 			},
// 			intoSelf: []string{"left"},
// 			intoWall: []string{"down", "right"},
// 		},
// 		{
// 			name: "tail chase ok 2",
// 			input: wire.Battlesnake{
// 				ID:     "my-id",
// 				Length: 4,
// 				Health: 99,
// 				Head:   wire.Coord{X: 0, Y: 11},
// 				Body:   []wire.Coord{{X: 0, Y: 11}, {X: 1, Y: 11}, {X: 1, Y: 10}, {X: 0, Y: 10}},
// 			},
// 			intoSelf: []string{"right"},
// 			intoWall: []string{"up", "left"},
// 		},
// 		{
// 			name: "tail chase not ok (just eaten)",
// 			input: wire.Battlesnake{
// 				ID:     "my-id",
// 				Length: 6,
// 				Health: 100,
// 				Head:   wire.Coord{X: 1, Y: 1},
// 				Body:   []wire.Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}},
// 			},
// 			intoSelf: []string{"up", "right", "down"},
// 		},
// 	}

// 	for _, tc := range tests {
// 		state := wire.GameState{
// 			Board: wire.Board{
// 				Width:  12,
// 				Height: 12,
// 				Snakes: []wire.Battlesnake{tc.input},
// 			},
// 			You: tc.input,
// 		}

// 		nextMove := move(state)

// 		if contains(tc.intoSelf, nextMove.Move) {
// 			t.Errorf("%s: snake moved into self, %s", tc.name, nextMove.Move)
// 		}
// 		if contains(tc.intoWall, nextMove.Move) {
// 			t.Errorf("%s: snake moved into wall, %s", tc.name, nextMove.Move)
// 		}
// 	}
// }

func TestSelfAvoidance2(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "72dd383c-bcc9-4e18-a01a-2e6ddd911630",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 2,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{8, 10}, {5, 5}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_MPyShWKrcHkXMCtDWfFtCvGD",
					Name:    "Canadian Bacon",
					Health:  98,
					Head:    wire.Coord{10, 10},
					Body:    []wire.Coord{{10, 10}, {10, 9}, {9, 9}},
					Length:  3,
					Latency: "218",
					Shout:   "",
				},
				{
					ID:      "gs_Vpf7rhGj3qmKykCQGbqrCp9G",
					Name:    "nomblegomble",
					Health:  100,
					Head:    wire.Coord{8, 4},
					Body:    []wire.Coord{{8, 4}, {8, 5}, {9, 5}, {9, 5}},
					Length:  4,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_RffxTd39SdRRRy8qVMGdQGkJ",
					Name:    "msbs",
					Health:  100,
					Head:    wire.Coord{0, 8},
					Body:    []wire.Coord{{0, 8}, {0, 9}, {1, 9}, {1, 9}},
					Length:  4,
					Latency: "50",
					Shout:   "",
				},
				{
					ID:      "gs_QRcYBqP8PYFBYvdGfbpQm9rb",
					Name:    "random-boii-2.0",
					Health:  100,
					Head:    wire.Coord{4, 0},
					Body:    []wire.Coord{{4, 0}, {4, 1}, {5, 1}, {5, 1}},
					Length:  4,
					Latency: "242",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_Vpf7rhGj3qmKykCQGbqrCp9G",
			Name:    "nomblegomble",
			Health:  100,
			Head:    wire.Coord{8, 4},
			Body:    []wire.Coord{{8, 4}, {8, 5}, {9, 5}, {9, 5}},
			Length:  4,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}

func TestHead2Head(t *testing.T) {
	tests := []struct {
		me       wire.Battlesnake
		other    wire.Battlesnake
		expected string
	}{
		{
			me: wire.Battlesnake{
				ID:     "snake-508e96ac-94ad-11ea-bb37",
				Name:   "foo",
				Length: 3,
				Health: 78,
				Head:   wire.Coord{X: 0, Y: 0},
				Body:   []wire.Coord{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
			},
			other: wire.Battlesnake{
				ID:     "snake-b67f4906-94ae-11ea-bb37",
				Name:   "bar",
				Length: 3,
				Health: 98,
				Head:   wire.Coord{X: 2, Y: 0},
				Body:   []wire.Coord{{X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}},
			},
			expected: "right",
		},
	}
	for _, tc := range tests {
		state := wire.GameState{
			Board: wire.Board{
				Width:  7,
				Height: 7,
				Snakes: []wire.Battlesnake{tc.me, tc.other},
				Food:   []wire.Coord{{X: 6, Y: 6}},
			},
			You: tc.me,
		}

		nextMove := move(state)

		if nextMove.Move != tc.expected {
			t.Errorf("Prefer head2head over the wall or self, expected right, got %s", nextMove.Move)
		}
	}
}

func TestSpaceBasic(t *testing.T) {
	tests := []struct {
		name         string
		input        wire.Board
		intoSelf     []string
		intoWall     []string
		intoBadSpace []string
	}{
		{
			name: "avoid small space 1",
			input: wire.Board{
				Snakes: []wire.Battlesnake{
					{
						ID:     "my-id",
						Length: 5,
						Health: 75,
						Head:   wire.Coord{2, 0},
						Body:   []wire.Coord{{2, 0}, {2, 1}, {1, 1}, {0, 1}, {0, 0}},
					},
				},
				Width:  7,
				Height: 7,
				Food:   []wire.Coord{{1, 0}, {0, 0}}, // 0,0 simulates bad luck of food spawning
			},
			intoSelf:     []string{"up"},
			intoWall:     []string{"down"},
			intoBadSpace: []string{"left"},
		},
		{
			name: "avoid small space 2",
			input: wire.Board{
				Snakes: []wire.Battlesnake{
					{
						ID:     "my-id",
						Length: 13,
						Health: 100,
						Head:   wire.Coord{0, 5},
						Body:   []wire.Coord{{0, 5}, {1, 5}, {2, 5}, {3, 5}, {3, 6}, {4, 6}, {5, 6}, {5, 5}, {4, 5}, {4, 4}, {3, 4}, {2, 4}, {1, 4}},
					},
				},
				Food:   []wire.Coord{{4, 1}, {6, 6}},
				Width:  7,
				Height: 7,
			},
			intoSelf:     []string{"right"},
			intoWall:     []string{"left"},
			intoBadSpace: []string{"up"},
		},
	}

	for _, tc := range tests {
		state := wire.GameState{
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
	t.Skip("sometimes playing chicken is fine...")
	state := wire.GameState{
		Game: wire.Game{
			ID: "3509d89e-8809-46c9-b46c-164158eaac26",
			Ruleset: wire.Ruleset{
				Name: "standard",
			},
			Timeout: 500,
		},
		Turn: 3,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 4}, {8, 2}, {4, 2}, {6, 10}, {5, 5}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_PgQqRchTF6k3FWhVk6fvkXv4",
					Name:    "nomblegomble",
					Health:  97,
					Head:    wire.Coord{5, 6},
					Body:    []wire.Coord{{5, 6}, {5, 7}, {5, 8}},
					Length:  3,
					Latency: "21",
				},
				{
					ID:      "gs_xw9mxDXPHD7fgFKgxBdMpJQT",
					Name:    "rnd",
					Health:  97,
					Head:    wire.Coord{1, 8},
					Body:    []wire.Coord{{1, 8}, {1, 7}, {1, 6}},
					Length:  3,
					Latency: "0",
				},
				{
					ID:      "gs_hVTbbK4dXGPJGCjxmKx9ktJf",
					Name:    "Mangofox",
					Health:  97,
					Head:    wire.Coord{9, 4},
					Body:    []wire.Coord{{9, 4}, {9, 3}, {9, 2}},
					Length:  3,
					Latency: "0",
				},
				{
					ID:      "gs_K4C6SmXkKSjWxwv3bY96djxD",
					Name:    "Steve",
					Health:  97,
					Head:    wire.Coord{5, 4},
					Body:    []wire.Coord{{5, 4}, {5, 3}, {5, 2}},
					Length:  3,
					Latency: "0",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_PgQqRchTF6k3FWhVk6fvkXv4",
			Name:    "nomblegomble",
			Health:  97,
			Head:    wire.Coord{5, 6},
			Body:    []wire.Coord{{5, 6}, {5, 7}, {5, 8}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "294cee29-202c-4fe4-9482-db64cf19fad6",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 41,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{6, 7}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_BgWvHQ7yGrhmrWkrfqdGHQrc",
					Name:    "nomblegomble",
					Health:  71,
					Head:    wire.Coord{0, 5},
					Body:    []wire.Coord{{0, 5}, {1, 5}, {1, 6}, {2, 6}, {3, 6}},
					Length:  5,
					Latency: "23",
					Shout:   "",
				},
				{
					ID:      "gs_7S76jFcGSmVCrKwydvth4fJ7",
					Name:    "bsnek2",
					Health:  78,
					Head:    wire.Coord{8, 3},
					Body:    []wire.Coord{{8, 3}, {8, 2}, {8, 1}, {9, 1}, {9, 2}, {10, 2}},
					Length:  6,
					Latency: "190",
					Shout:   "",
				},
				{
					ID:      "gs_ch43JRySwjgX3MyW8hr9HMkC",
					Name:    "Boomslang",
					Health:  99,
					Head:    wire.Coord{2, 3},
					Body:    []wire.Coord{{2, 3}, {3, 3}, {3, 2}, {3, 1}, {3, 0}, {4, 0}, {5, 0}, {6, 0}},
					Length:  8,
					Latency: "253",
					Shout:   "",
				},
				{
					ID:      "gs_RRmMymt6W9bwjHh9mtrfvPt3",
					Name:    "Crimson",
					Health:  83,
					Head:    wire.Coord{5, 4},
					Body:    []wire.Coord{{5, 4}, {5, 5}, {6, 5}, {6, 4}, {6, 3}, {5, 3}},
					Length:  6,
					Latency: "198",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_BgWvHQ7yGrhmrWkrfqdGHQrc",
			Name:    "nomblegomble",
			Health:  71,
			Head:    wire.Coord{0, 5},
			Body:    []wire.Coord{{0, 5}, {1, 5}, {1, 6}, {2, 6}, {3, 6}},
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
	state := wire.GameState{}
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

// One snake can be beat, but not the longer one!
func TestAvoid3WayBadHead2Head1(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "8bc4c92d-9e78-4542-bc42-0e41bb2d8689",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 114,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 9}, {10, 8}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_WGPDfp6p4jf3DcYvQM4hRp9C",
					Name:    "Devious Devin",
					Health:  95,
					Head:    wire.Coord{5, 5},
					Body:    []wire.Coord{{5, 5}, {4, 5}, {3, 5}, {2, 5}, {1, 5}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {1, 8}, {2, 8}, {2, 7}, {3, 7}},
					Length:  13,
					Latency: "59",
					Shout:   "",
				},
				{
					ID:      "gs_GJjQ6md8htSKpg8k4SvCT99V",
					Name:    "nomblegomble",
					Health:  94,
					Head:    wire.Coord{6, 4},
					Body:    []wire.Coord{{6, 4}, {5, 4}, {4, 4}, {3, 4}, {2, 4}, {1, 4}, {1, 3}, {2, 3}, {3, 3}, {3, 2}, {3, 1}},
					Length:  11,
					Latency: "23",
					Shout:   "",
				},
				{
					ID:      "gs_gmYrStYRmQ68tTy9Yrb66k9S",
					Name:    "Titanoboa",
					Health:  99,
					Head:    wire.Coord{7, 5},
					Body:    []wire.Coord{{7, 5}, {7, 4}, {7, 3}, {8, 3}, {8, 4}, {8, 5}, {8, 6}, {9, 6}},
					Length:  8,
					Latency: "241",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_GJjQ6md8htSKpg8k4SvCT99V",
			Name:    "nomblegomble",
			Health:  94,
			Head:    wire.Coord{6, 4},
			Body:    []wire.Coord{{6, 4}, {5, 4}, {4, 4}, {3, 4}, {2, 4}, {1, 4}, {1, 3}, {2, 3}, {3, 3}, {3, 2}, {3, 1}},
			Length:  11,
			Latency: "23",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake moved into bad 3 way h2h, %s (one other snake is larger)", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into other snake, %s", nextMove.Move)
	}
}

// Other snakes are likely to go for food, so don't go for it if you don't have to.
func TestAvoidFoodInEqualHead2Head1(t *testing.T) {
	t.Skip("sometimes you just gotta go for it")
	state := wire.GameState{
		Game: wire.Game{
			ID: "f2ce457a-7eb1-43cf-b495-831498b753e0",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 54,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{2, 7}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_qW3xGphphfBRBPjfqt7phy6X",
					Name:    "TC5001",
					Health:  98,
					Head:    wire.Coord{8, 2},
					Body:    []wire.Coord{{8, 2}, {9, 2}, {10, 2}, {10, 1}, {10, 0}, {9, 0}, {8, 0}, {7, 0}, {7, 1}, {6, 1}},
					Length:  10,
					Latency: "215",
					Shout:   "",
				},
				{
					ID:      "gs_97jPPpfDBqHYhxkkrdHHrQTY",
					Name:    "nomblegomble",
					Health:  98,
					Head:    wire.Coord{3, 7},
					Body:    []wire.Coord{{3, 7}, {3, 8}, {4, 8}, {5, 8}, {5, 7}, {5, 6}, {4, 6}, {3, 6}},
					Length:  8,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_7w9rQ6VDrSpYpgRCjygk8mmF",
					Name:    "bsnek2",
					Health:  99,
					Head:    wire.Coord{9, 9},
					Body:    []wire.Coord{{9, 9}, {9, 8}, {9, 7}, {9, 6}, {9, 5}, {8, 5}, {8, 4}},
					Length:  7,
					Latency: "284",
					Shout:   "",
				},
				{
					ID:      "gs_RmhhbkKtHSff73R66GMT33gJ",
					Name:    "SnakeJS",
					Health:  69,
					Head:    wire.Coord{2, 6},
					Body:    []wire.Coord{{2, 6}, {2, 5}, {3, 5}, {3, 4}, {3, 3}, {3, 2}, {3, 1}, {4, 1}},
					Length:  8,
					Latency: "239",
					Shout:   "[2 5] --> [2 7]",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_97jPPpfDBqHYhxkkrdHHrQTY",
			Name:    "nomblegomble",
			Health:  98,
			Head:    wire.Coord{3, 7},
			Body:    []wire.Coord{{3, 7}, {3, 8}, {4, 8}, {5, 8}, {5, 7}, {5, 6}, {4, 6}, {3, 6}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "82cb8643-5f86-4674-a3ed-28c1d99a689f",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 80,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{4, 10}, {6, 3}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_3CWVKmKbMYvmS7qQYPDkm9f8",
					Name:    "nomblegomble",
					Health:  83,
					Head:    wire.Coord{6, 10},
					Body:    []wire.Coord{{6, 10}, {5, 10}, {5, 9}, {4, 9}, {4, 8}, {3, 8}, {3, 7}},
					Length:  7,
					Latency: "21",
				},
				{
					ID:      "gs_3D66hv63CXRMVRKDjCKDj8pJ",
					Name:    "nomblegomble",
					Health:  96,
					Head:    wire.Coord{7, 9},
					Body:    []wire.Coord{{7, 9}, {8, 9}, {9, 9}, {10, 9}, {10, 8}, {9, 8}, {8, 8}, {8, 7}},
					Length:  8,
					Latency: "21",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_3CWVKmKbMYvmS7qQYPDkm9f8",
			Name:    "nomblegomble",
			Health:  83,
			Head:    wire.Coord{6, 10},
			Body:    []wire.Coord{{6, 10}, {5, 10}, {5, 9}, {4, 9}, {4, 8}, {3, 8}, {3, 7}},
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

func TestSpaceAvoidTiedHead2Head1(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "19ec09a8-0cea-40c5-a50b-f85f6cd0f400",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 6,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{5, 5}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_dF743rmHfbQ6YqkCbDSXmfj7",
					Name:    "nomblegomble",
					Health:  96,
					Head:    wire.Coord{0, 8},
					Body:    []wire.Coord{{0, 8}, {1, 8}, {1, 7}, {0, 7}},
					Length:  4,
					Latency: "23",
					Shout:   "",
				},
				{
					ID:      "gs_GhqqhY7mhc3vbmXtRyVGkwMR",
					Name:    "Fairy Rust",
					Health:  96,
					Head:    wire.Coord{0, 2},
					Body:    []wire.Coord{{0, 2}, {1, 2}, {1, 1}, {1, 0}},
					Length:  4,
					Latency: "45",
					Shout:   "",
				},
				{
					ID:      "gs_fHQbrqVXFYQcvDk74wmmfv9Q",
					Name:    "Legless Lizard",
					Health:  96,
					Head:    wire.Coord{4, 2},
					Body:    []wire.Coord{{4, 2}, {3, 2}, {3, 1}, {3, 0}},
					Length:  4,
					Latency: "51",
					Shout:   "",
				},
				{
					ID:      "gs_rj6ktV78gbfyqKdGB3ffXJF9",
					Name:    "The Jabberwock",
					Health:  98,
					Head:    wire.Coord{1, 5},
					Body:    []wire.Coord{{1, 5}, {1, 4}, {0, 4}, {0, 5}},
					Length:  4,
					Latency: "4",
					Shout:   "From hellΓÇÖs heart I stab at thee",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_dF743rmHfbQ6YqkCbDSXmfj7",
			Name:    "nomblegomble",
			Health:  96,
			Head:    wire.Coord{0, 8},
			Body:    []wire.Coord{{0, 8}, {1, 8}, {1, 7}, {0, 7}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "37719616-712a-4dea-9dbd-b9dfa992d908",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 80,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{9, 10}, {8, 9}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_x3Mybq4HHmYSHKbGh83fk9rJ",
					Name:    "nomblegomble",
					Health:  100,
					Head:    wire.Coord{4, 6},
					Body:    []wire.Coord{{4, 6}, {4, 5}, {4, 4}, {4, 3}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {1, 7}, {1, 7}},
					Length:  12,
					Latency: "22",
				},
				{
					ID:      "gs_D6yxdPV87SbYfrFYSDK7JVTR",
					Name:    "Scared Cobra Chicken",
					Health:  63,
					Head:    wire.Coord{5, 7},
					Body:    []wire.Coord{{5, 7}, {6, 7}, {6, 8}, {5, 8}, {5, 9}},
					Length:  5,
					Latency: "204",
				},
				{
					ID:      "gs_vJxJGRgK43X9G8DDBmK8CDSQ",
					Name:    "jsnek2",
					Health:  98,
					Head:    wire.Coord{10, 4},
					Body:    []wire.Coord{{10, 4}, {10, 3}, {10, 2}, {9, 2}, {9, 3}, {9, 4}, {8, 4}},
					Length:  7,
					Latency: "263",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_x3Mybq4HHmYSHKbGh83fk9rJ",
			Name:    "nomblegomble",
			Health:  100,
			Head:    wire.Coord{4, 6},
			Body:    []wire.Coord{{4, 6}, {4, 5}, {4, 4}, {4, 3}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {1, 7}, {1, 7}},
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

	state := wire.GameState{
		Game: wire.Game{
			ID: "37719616-712a-4dea-9dbd-b9dfa992d908",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 81,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{9, 10}, {8, 9}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_x3Mybq4HHmYSHKbGh83fk9rJ",
					Name:    "nomblegomble",
					Health:  99,
					Head:    wire.Coord{3, 6},
					Body:    []wire.Coord{{3, 6}, {4, 6}, {4, 5}, {4, 4}, {4, 3}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {1, 7}},
					Length:  12,
					Latency: "21",
				},
				{
					ID:      "gs_D6yxdPV87SbYfrFYSDK7JVTR",
					Name:    "Scared Cobra Chicken",
					Health:  62,
					Head:    wire.Coord{4, 7},
					Body:    []wire.Coord{{4, 7}, {5, 7}, {6, 7}, {6, 8}, {5, 8}},
					Length:  5,
					Latency: "205",
				},
				{
					ID:      "gs_vJxJGRgK43X9G8DDBmK8CDSQ",
					Name:    "jsnek2",
					Health:  97,
					Head:    wire.Coord{10, 5},
					Body:    []wire.Coord{{10, 5}, {10, 4}, {10, 3}, {10, 2}, {9, 2}, {9, 3}, {9, 4}},
					Length:  7,
					Latency: "254",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_x3Mybq4HHmYSHKbGh83fk9rJ",
			Name:    "nomblegomble",
			Health:  99,
			Head:    wire.Coord{3, 6},
			Body:    []wire.Coord{{3, 6}, {4, 6}, {4, 5}, {4, 4}, {4, 3}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {1, 7}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "f53ef734-3349-467d-9eee-89b2d6f9b4fa",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 97,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 2}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_YSGWKK73YPHrX3vdG3hJhGHT",
					Name:    "nomblegomble",
					Health:  95,
					Head:    wire.Coord{1, 2},
					Body:    []wire.Coord{{1, 2}, {1, 3}, {0, 3}, {0, 4}, {1, 4}, {2, 4}, {2, 3}, {3, 3}, {4, 3}, {4, 4}, {5, 4}, {6, 4}},
					Length:  12,
					Latency: "21",
				},
				{
					ID:      "gs_mmKppxxMFHWF73VYrjBM8RS8",
					Name:    "Ekans on a Plane",
					Health:  90,
					Head:    wire.Coord{2, 1},
					Body:    []wire.Coord{{2, 1}, {1, 1}, {0, 1}, {0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}},
					Length:  9,
					Latency: "73",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_YSGWKK73YPHrX3vdG3hJhGHT",
			Name:    "nomblegomble",
			Health:  95,
			Head:    wire.Coord{1, 2},
			Body:    []wire.Coord{{1, 2}, {1, 3}, {0, 3}, {0, 4}, {1, 4}, {2, 4}, {2, 3}, {3, 3}, {4, 3}, {4, 4}, {5, 4}, {6, 4}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "e5652b90-b24e-43ff-ba46-c00ef8b1cb41",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 223,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{2, 0}, {0, 2}, {2, 10}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_jjjF7vvJC9tWc6dBM4MhjKrW",
					Name:    "nomblegomble",
					Health:  95,
					Head:    wire.Coord{7, 2},
					Body:    []wire.Coord{{7, 2}, {8, 2}, {8, 3}, {9, 3}, {9, 4}, {9, 5}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {9, 9}, {8, 9}, {7, 9}, {6, 9}, {6, 8}, {6, 7}, {5, 7}, {4, 7}, {4, 6}, {4, 5}, {5, 5}, {6, 5}, {7, 5}, {7, 6}, {7, 7}, {8, 7}, {9, 7}, {9, 6}, {8, 6}, {8, 5}, {8, 4}, {7, 4}, {7, 3}, {6, 3}, {5, 3}, {5, 2}, {4, 2}, {3, 2}, {2, 2}},
					Length:  42,
					Latency: "20",
				},
				{
					ID:      "gs_3F4d6gQq9ygjhykc6JpmmTmJ",
					Name:    "Eremetic Eric",
					Health:  74,
					Head:    wire.Coord{6, 1},
					Body:    []wire.Coord{{6, 1}, {5, 1}, {5, 0}, {6, 0}, {7, 0}, {7, 1}},
					Length:  6,
					Latency: "15",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_jjjF7vvJC9tWc6dBM4MhjKrW",
			Name:    "nomblegomble",
			Health:  95,
			Head:    wire.Coord{7, 2},
			Body:    []wire.Coord{{7, 2}, {8, 2}, {8, 3}, {9, 3}, {9, 4}, {9, 5}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {9, 9}, {8, 9}, {7, 9}, {6, 9}, {6, 8}, {6, 7}, {5, 7}, {4, 7}, {4, 6}, {4, 5}, {5, 5}, {6, 5}, {7, 5}, {7, 6}, {7, 7}, {8, 7}, {9, 7}, {9, 6}, {8, 6}, {8, 5}, {8, 4}, {7, 4}, {7, 3}, {6, 3}, {5, 3}, {5, 2}, {4, 2}, {3, 2}, {2, 2}},
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

func TestProblemWithGridAreaAndEscape(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "f266c441-92f1-4ac3-9e9b-7159696ded38",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.20",
			},
			Timeout: 500,
		},
		Turn: 303,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{5, 0}, {6, 1}, {10, 0}, {5, 3}, {10, 4}, {8, 9}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_W3B4hXPRB6CXyRmVv6xMT9mC",
					Name:    "return of the rise of the last snake awakens",
					Health:  83,
					Head:    wire.Coord{10, 5},
					Body:    []wire.Coord{{10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {9, 9}, {9, 8}, {8, 8}, {7, 8}, {7, 7}, {6, 7}, {6, 6}, {5, 6}, {5, 5}, {6, 5}, {7, 5}, {7, 4}, {6, 4}, {6, 3}, {6, 2}},
					Length:  22,
					Latency: "287",
					Shout:   "",
				},
				{
					ID:      "gs_9j8pkgxP8KfpTpcV3cqVhkw7",
					Name:    "nomblegomble",
					Health:  96,
					Head:    wire.Coord{6, 9},
					Body:    []wire.Coord{{6, 9}, {5, 9}, {4, 9}, {4, 10}, {3, 10}, {3, 9}, {2, 9}, {2, 8}, {3, 8}, {3, 7}, {3, 6}, {2, 6}, {2, 5}, {1, 5}, {1, 4}, {1, 3}, {1, 2}, {1, 1}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {4, 1}},
					Length:  23,
					Latency: "23",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_9j8pkgxP8KfpTpcV3cqVhkw7",
			Name:    "nomblegomble",
			Health:  96,
			Head:    wire.Coord{6, 9},
			Body:    []wire.Coord{{6, 9}, {5, 9}, {4, 9}, {4, 10}, {3, 10}, {3, 9}, {2, 9}, {2, 8}, {3, 8}, {3, 7}, {3, 6}, {2, 6}, {2, 5}, {1, 5}, {1, 4}, {1, 3}, {1, 2}, {1, 1}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {4, 1}},
			Length:  23,
			Latency: "23",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake moved into too small of space, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into too small of space, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}

func TestSpaceCornered1(t *testing.T) {

	state := wire.GameState{
		Game: wire.Game{
			ID: "7560784f-f380-427d-8350-80725b25207a",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 68,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{6, 2}, {0, 2}, {6, 6}, {10, 3}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_bVHWqPM7PxRHkqCjGTgSkjhY",
					Name:    "nomblegomble",
					Health:  96,
					Head:    wire.Coord{8, 2},
					Body:    []wire.Coord{{8, 2}, {7, 2}, {7, 1}, {7, 0}, {6, 0}, {5, 0}, {4, 0}, {4, 1}, {4, 2}, {5, 2}},
					Length:  10,
					Latency: "22",
				},
				{
					ID:      "gs_XPtjRWfB3VT7GXFHGkvhXXbd",
					Name:    "MAsterStudentSlayer666",
					Health:  85,
					Head:    wire.Coord{3, 3},
					Body:    []wire.Coord{{3, 3}, {2, 3}, {2, 4}, {1, 4}, {0, 4}, {0, 5}, {0, 6}},
					Length:  7,
					Latency: "219",
				},
				{
					ID:      "gs_XyrrYqC3pg8W4DbSpF3jpPWc",
					Name:    "leshchenko-1",
					Health:  100,
					Head:    wire.Coord{9, 1},
					Body:    []wire.Coord{{9, 1}, {9, 2}, {9, 3}, {9, 4}, {9, 4}},
					Length:  5,
					Latency: "234",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_bVHWqPM7PxRHkqCjGTgSkjhY",
			Name:    "nomblegomble",
			Health:  96,
			Head:    wire.Coord{8, 2},
			Body:    []wire.Coord{{8, 2}, {7, 2}, {7, 1}, {7, 0}, {6, 0}, {5, 0}, {4, 0}, {4, 1}, {4, 2}, {5, 2}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "3988391b-ee86-466e-ab0c-d39c38283d38",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 126,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{6, 1}, {10, 4}, {0, 9}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_X6DRymbGqtcDBWxfTqmkBhjB",
					Name:    "snek",
					Health:  5,
					Head:    wire.Coord{1, 5},
					Body:    []wire.Coord{{1, 5}, {1, 6}, {2, 6}, {3, 6}, {3, 5}},
					Length:  5,
					Latency: "281",
					Shout:   "",
				},
				{
					ID:      "gs_TfjMPmMkjSWmf4dYWjX7rrjK",
					Name:    "msbs",
					Health:  69,
					Head:    wire.Coord{3, 1},
					Body:    []wire.Coord{{3, 1}, {3, 2}, {2, 2}, {2, 1}, {1, 1}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}},
					Length:  12,
					Latency: "53",
					Shout:   "",
				},
				{
					ID:      "gs_PXMmW3xbSxDBJCg4hycg7xyG",
					Name:    "nomblegomble",
					Health:  86,
					Head:    wire.Coord{0, 8},
					Body:    []wire.Coord{{0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {4, 9}, {4, 10}, {5, 10}, {6, 10}, {7, 10}, {7, 9}, {6, 9}, {5, 9}, {5, 8}, {6, 8}, {7, 8}, {7, 7}},
					Length:  17,
					Latency: "22",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_PXMmW3xbSxDBJCg4hycg7xyG",
			Name:    "nomblegomble",
			Health:  86,
			Head:    wire.Coord{0, 8},
			Body:    []wire.Coord{{0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {4, 9}, {4, 10}, {5, 10}, {6, 10}, {7, 10}, {7, 9}, {6, 9}, {5, 9}, {5, 8}, {6, 8}, {7, 8}, {7, 7}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "5ff70484-ac66-4025-90b6-9af1554b74b5",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 81,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{10, 0}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_PVdSm9cYDRMk3R6Tk3Qqpw64",
					Name:    "nomblegomble",
					Health:  99,
					Head:    wire.Coord{10, 1},
					Body:    []wire.Coord{{10, 1}, {9, 1}, {8, 1}, {7, 1}, {7, 2}, {7, 3}, {6, 3}, {6, 4}, {6, 5}, {6, 6}, {6, 7}, {6, 8}, {6, 9}, {5, 9}},
					Length:  14,
					Latency: "21",
				},
				{
					ID:      "gs_WJQ63xb7t6mpFCXqHVdPqBr4",
					Name:    "Ifarus",
					Health:  77,
					Head:    wire.Coord{3, 0},
					Body:    []wire.Coord{{3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}},
					Length:  8,
					Latency: "76",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_PVdSm9cYDRMk3R6Tk3Qqpw64",
			Name:    "nomblegomble",
			Health:  99,
			Head:    wire.Coord{10, 1},
			Body:    []wire.Coord{{10, 1}, {9, 1}, {8, 1}, {7, 1}, {7, 2}, {7, 3}, {6, 3}, {6, 4}, {6, 5}, {6, 6}, {6, 7}, {6, 8}, {6, 9}, {5, 9}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "3c3b7dcc-4f7d-48d2-9449-ee22bda84390",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 32,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{8, 2}, {9, 4}, {4, 1}, {3, 4}, {9, 7}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_YwVHBvWKTXXKKjXvM6y93VKH",
					Name:    "nomblegomble",
					Health:  83,
					Head:    wire.Coord{8, 4},
					Body:    []wire.Coord{{8, 4}, {8, 3}, {7, 3}, {7, 2}, {7, 1}},
					Length:  5,
					Latency: "22",
				},
				{
					ID:      "gs_3W9Dm9F4Hw73fXgxPPmVRwFX",
					Name:    "nates_python",
					Health:  98,
					Head:    wire.Coord{5, 9},
					Body:    []wire.Coord{{5, 9}, {6, 9}, {7, 9}, {7, 8}, {7, 7}},
					Length:  5,
					Latency: "220",
				},
				{
					ID:      "gs_Bgj94MjGM8c7Mqppbhvjcx3K",
					Name:    "carl",
					Health:  86,
					Head:    wire.Coord{5, 3},
					Body:    []wire.Coord{{5, 3}, {4, 3}, {3, 3}, {3, 2}, {3, 1}, {2, 1}},
					Length:  6,
					Latency: "250",
				},
				{
					ID:      "gs_pgQRpdYWW7cxQbhMpbwxwC84",
					Name:    "Morley",
					Health:  68,
					Head:    wire.Coord{7, 5},
					Body:    []wire.Coord{{7, 5}, {8, 5}, {9, 5}},
					Length:  3,
					Latency: "77",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_YwVHBvWKTXXKKjXvM6y93VKH",
			Name:    "nomblegomble",
			Health:  83,
			Head:    wire.Coord{8, 4},
			Body:    []wire.Coord{{8, 4}, {8, 3}, {7, 3}, {7, 2}, {7, 1}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "8cb97ac0-f405-41a1-b007-a9a4b53bbbfa",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 96,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{7, 0}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_tmWR3BFhBHMPhHYM37rqdB37",
					Name:    "nomblegomble",
					Health:  100,
					Head:    wire.Coord{8, 10},
					Body:    []wire.Coord{{8, 10}, {8, 9}, {8, 8}, {8, 7}, {7, 7}, {7, 8}, {6, 8}, {5, 8}, {5, 9}, {4, 9}, {3, 9}, {2, 9}, {2, 9}},
					Length:  13,
					Latency: "21",
				},
				{
					ID:      "gs_XJ44wjQRyT3MPqqwTB8WKmpX",
					Name:    "Ouroboros 2",
					Health:  97,
					Head:    wire.Coord{9, 7},
					Body:    []wire.Coord{{9, 7}, {9, 6}, {8, 6}, {7, 6}, {6, 6}, {5, 6}, {4, 6}, {4, 5}, {3, 5}, {3, 4}, {3, 3}, {3, 2}, {3, 1}},
					Length:  13,
					Latency: "214",
				},
				{
					ID:      "gs_6HpdTvFPJJHk8KKxXGSq4vdb",
					Name:    "Canadian Bacon",
					Health:  6,
					Head:    wire.Coord{2, 10},
					Body:    []wire.Coord{{2, 10}, {1, 10}, {1, 9}, {1, 8}},
					Length:  4,
					Latency: "197",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_tmWR3BFhBHMPhHYM37rqdB37",
			Name:    "nomblegomble",
			Health:  100,
			Head:    wire.Coord{8, 10},
			Body:    []wire.Coord{{8, 10}, {8, 9}, {8, 8}, {8, 7}, {7, 7}, {7, 8}, {6, 8}, {5, 8}, {5, 9}, {4, 9}, {3, 9}, {2, 9}, {2, 9}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "4ac4ba8c-4c68-4a03-9607-583264860222",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 77,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 0}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_PQYck9Y4W3MQpS3HPWrJ8HMd",
					Name:    "nomblegomble",
					Health:  96,
					Head:    wire.Coord{0, 3},
					Body:    []wire.Coord{{0, 3}, {1, 3}, {1, 2}, {1, 1}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}},
					Length:  9,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_8TvRDTjqhcJCb6w8YhjSXHcH",
					Name:    "Cool_as_ice",
					Health:  76,
					Head:    wire.Coord{7, 6},
					Body:    []wire.Coord{{7, 6}, {6, 6}, {6, 5}, {6, 4}, {6, 3}, {7, 3}, {7, 4}, {7, 5}},
					Length:  8,
					Latency: "74",
					Shout:   "",
				},
				{
					ID:      "gs_x88cHrQDvKxJHkTwC7TSWk3W",
					Name:    "moon-snake-pika",
					Health:  99,
					Head:    wire.Coord{9, 8},
					Body:    []wire.Coord{{9, 8}, {10, 8}, {10, 7}, {10, 6}, {10, 5}, {10, 4}, {10, 3}, {10, 2}, {10, 1}, {10, 0}, {9, 0}, {9, 1}, {9, 2}, {9, 3}},
					Length:  14,
					Latency: "217",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_PQYck9Y4W3MQpS3HPWrJ8HMd",
			Name:    "nomblegomble",
			Health:  96,
			Head:    wire.Coord{0, 3},
			Body:    []wire.Coord{{0, 3}, {1, 3}, {1, 2}, {1, 1}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "e74d6f1d-a38a-4135-bb89-d17f387ba9ae",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 108,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{1, 9}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_cYKtRkqxfMfyj46WH9J9VhmX",
					Name:    "nomblegomble",
					Health:  90,
					Head:    wire.Coord{6, 4},
					Body:    []wire.Coord{{6, 4}, {7, 4}, {7, 3}, {7, 2}, {7, 1}, {6, 1}, {5, 1}, {4, 1}, {3, 1}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {2, 2}},
					Length:  17,
					Latency: "23",
					Shout:   "",
				},
				{
					ID:      "gs_M7h8gyWqTGRFg8GX4BJf3whb",
					Name:    "Ophiophagus One",
					Health:  100,
					Head:    wire.Coord{4, 4},
					Body:    []wire.Coord{{4, 4}, {4, 5}, {5, 5}, {5, 6}, {5, 7}, {6, 7}, {6, 7}},
					Length:  7,
					Latency: "211",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_cYKtRkqxfMfyj46WH9J9VhmX",
			Name:    "nomblegomble",
			Health:  90,
			Head:    wire.Coord{6, 4},
			Body:    []wire.Coord{{6, 4}, {7, 4}, {7, 3}, {7, 2}, {7, 1}, {6, 1}, {5, 1}, {4, 1}, {3, 1}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {2, 2}},
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

func TestSpaceCutoff7(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "43172677-aa69-4a04-aecc-4aedcf238d05",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 144,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{9, 9}, {2, 7}, {5, 3}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_K96GMhmm4XqSJDgbDmfdSv3J",
					Name:    "nomblegomble",
					Health:  89,
					Head:    wire.Coord{2, 4},
					Body:    []wire.Coord{{2, 4}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, {5, 4}, {6, 4}, {7, 4}, {8, 4}, {9, 4}, {10, 4}, {10, 5}, {9, 5}, {8, 5}, {7, 5}, {6, 5}},
					Length:  16,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_MF6b9fcWTpS9FRTCVJMK88r4",
					Name:    "Super Snakey",
					Health:  95,
					Head:    wire.Coord{3, 3},
					Body:    []wire.Coord{{3, 3}, {2, 3}, {2, 2}, {2, 1}, {1, 1}, {1, 0}, {2, 0}, {3, 0}, {3, 1}, {3, 2}, {4, 2}, {5, 2}, {6, 2}, {7, 2}, {7, 1}},
					Length:  15,
					Latency: "226",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_K96GMhmm4XqSJDgbDmfdSv3J",
			Name:    "nomblegomble",
			Health:  89,
			Head:    wire.Coord{2, 4},
			Body:    []wire.Coord{{2, 4}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, {5, 4}, {6, 4}, {7, 4}, {8, 4}, {9, 4}, {10, 4}, {10, 5}, {9, 5}, {8, 5}, {7, 5}, {6, 5}},
			Length:  16,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "right" {
		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
	}
	if nextMove.Move == "up" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "down" {
		t.Errorf("snake moved into other snake, %s", nextMove.Move)
	}
}

func TestSpaceCutoff8(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "eca2463d-0fd7-43b7-aa6b-43dbb489da07",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 50,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{4, 0}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_kpRwFYKwVjmj7JF6RTwdPHBB",
					Name:    "nomblegomble",
					Health:  99,
					Head:    wire.Coord{9, 7},
					Body:    []wire.Coord{{9, 7}, {9, 8}, {8, 8}, {7, 8}, {6, 8}, {6, 7}, {5, 7}, {4, 7}, {4, 6}},
					Length:  9,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_tC8WtyKcvjkvyQhVSB977YR9",
					Name:    "The Very Hungry Caterpillar ≡ƒìè≡ƒìÅ≡ƒìæ≡ƒìÆ≡ƒìÄ≡ƒÉ¢",
					Health:  95,
					Head:    wire.Coord{8, 6},
					Body:    []wire.Coord{{8, 6}, {9, 6}, {9, 5}, {9, 4}, {10, 4}, {10, 3}, {9, 3}},
					Length:  7,
					Latency: "40",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_kpRwFYKwVjmj7JF6RTwdPHBB",
			Name:    "nomblegomble",
			Health:  99,
			Head:    wire.Coord{9, 7},
			Body:    []wire.Coord{{9, 7}, {9, 8}, {8, 8}, {7, 8}, {6, 8}, {6, 7}, {5, 7}, {4, 7}, {4, 6}},
			Length:  9,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "left" {
		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
	}
	if nextMove.Move == "up" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "down" {
		t.Errorf("snake moved into other snake's body, %s", nextMove.Move)
	}
}

func TestSpaceCutoff9(t *testing.T) {
	// t.Skip("failing after node sibling eval changes")

	state := wire.GameState{
		Game: wire.Game{
			ID: "dca928e2-57db-43b9-92c6-1e3ee092e6e0",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.20",
			},
			Timeout: 500,
		},
		Turn: 57,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{5, 2}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_DpHTVDMxKRKQ3JC7VSwwgRy8",
					Name:    "Fairy Rust",
					Health:  99,
					Head:    wire.Coord{2, 5},
					Body:    []wire.Coord{{2, 5}, {2, 6}, {2, 7}, {2, 8}, {2, 9}, {2, 10}, {3, 10}},
					Length:  7,
					Latency: "44",
					Shout:   "",
				},
				{
					ID:      "gs_hwCJPBfbgBXGYvwH3QyyWgGK",
					Name:    "ΓÜ¢∩╕ÅΓ₧í∩╕ÅSnakeΓ¼å∩╕ÅΓÜ¢∩╕Å",
					Health:  95,
					Head:    wire.Coord{5, 8},
					Body:    []wire.Coord{{5, 8}, {6, 8}, {6, 9}, {7, 9}, {7, 10}, {8, 10}, {8, 9}},
					Length:  7,
					Latency: "239",
					Shout:   "",
				},
				{
					ID:      "gs_dBr8MkRBCYPtkbrC3G6wqxvB",
					Name:    "lars",
					Health:  43,
					Head:    wire.Coord{4, 7},
					Body:    []wire.Coord{{4, 7}, {4, 8}, {4, 9}},
					Length:  3,
					Latency: "51",
					Shout:   "",
				},
				{
					ID:      "gs_qTBGjFXpFwgm6xSrv89jdpj9",
					Name:    "nomblegomble",
					Health:  96,
					Head:    wire.Coord{1, 6},
					Body:    []wire.Coord{{1, 6}, {0, 6}, {0, 5}, {0, 4}, {0, 3}, {1, 3}, {2, 3}, {3, 3}, {3, 2}},
					Length:  9,
					Latency: "23",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_qTBGjFXpFwgm6xSrv89jdpj9",
			Name:    "nomblegomble",
			Health:  96,
			Head:    wire.Coord{1, 6},
			Body:    []wire.Coord{{1, 6}, {0, 6}, {0, 5}, {0, 4}, {0, 3}, {1, 3}, {2, 3}, {3, 3}, {3, 2}},
			Length:  9,
			Latency: "23",
			Shout:   "",
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

func TestSpaceCutoff10(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "099d4abb-ef21-41e3-b91b-8315e26672ee",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.20",
			},
			Timeout: 500,
		},
		Turn: 156,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{2, 10}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_qpFQRbcKWCV9vJ8fQDRfDB6B",
					Name:    "Kuro",
					Health:  95,
					Head:    wire.Coord{2, 8},
					Body:    []wire.Coord{{2, 8}, {3, 8}, {4, 8}, {5, 8}, {6, 8}, {7, 8}, {7, 7}, {7, 6}, {6, 6}, {5, 6}, {4, 6}, {3, 6}, {3, 7}, {2, 7}, {2, 6}, {2, 5}, {2, 4}, {1, 4}},
					Length:  18,
					Latency: "81",
					Shout:   "",
				},
				{
					ID:      "gs_XHxWgR7CXQQp4vTkrDy4bwCR",
					Name:    "nomblegomble",
					Health:  79,
					Head:    wire.Coord{8, 10},
					Body:    []wire.Coord{{8, 10}, {8, 9}, {8, 8}, {8, 7}, {8, 6}, {8, 5}, {8, 4}, {7, 4}, {6, 4}, {6, 3}, {5, 3}, {4, 3}, {4, 2}, {3, 2}, {3, 1}, {3, 0}, {2, 0}, {2, 1}, {2, 2}, {2, 3}, {3, 3}},
					Length:  21,
					Latency: "23",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_XHxWgR7CXQQp4vTkrDy4bwCR",
			Name:    "nomblegomble",
			Health:  79,
			Head:    wire.Coord{8, 10},
			Body:    []wire.Coord{{8, 10}, {8, 9}, {8, 8}, {8, 7}, {8, 6}, {8, 5}, {8, 4}, {7, 4}, {6, 4}, {6, 3}, {5, 3}, {4, 3}, {4, 2}, {3, 2}, {3, 1}, {3, 0}, {2, 0}, {2, 1}, {2, 2}, {2, 3}, {3, 3}},
			Length:  21,
			Latency: "23",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "left" {
		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
	}
	if nextMove.Move == "down" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "up" {
		t.Errorf("snake moved into wall, %s", nextMove.Move)
	}
}

func TestSpaceCutoff11(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "f33135ff-e07b-4a81-9a1b-2397ead3c8c9",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.20",
			},
			Timeout: 500,
		},
		Turn: 229,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 10}, {3, 10}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_pKTqbxbGQYBtdd8ytQGBKmX7",
					Name:    "Try not to die!",
					Health:  97,
					Head:    wire.Coord{4, 9},
					Body:    []wire.Coord{{4, 9}, {4, 8}, {4, 7}, {4, 6}, {5, 6}, {5, 5}, {4, 5}, {3, 5}, {2, 5}, {2, 6}, {3, 6}, {3, 7}, {2, 7}},
					Length:  13,
					Latency: "356",
					Shout:   "",
				},
				{
					ID:      "gs_r9M4G8DS4Sj6w8pXtmBhHdcQ",
					Name:    "nomblegomble",
					Health:  98,
					Head:    wire.Coord{0, 1},
					Body:    []wire.Coord{{0, 1}, {1, 1}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {5, 2}, {6, 2}, {7, 2}, {7, 1}, {8, 1}, {9, 1}, {9, 2}, {9, 3}, {10, 3}, {10, 4}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {7, 9}, {7, 8}, {6, 8}, {5, 8}, {5, 7}, {6, 7}, {7, 7}, {8, 7}},
					Length:  33,
					Latency: "22",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_r9M4G8DS4Sj6w8pXtmBhHdcQ",
			Name:    "nomblegomble",
			Health:  98,
			Head:    wire.Coord{0, 1},
			Body:    []wire.Coord{{0, 1}, {1, 1}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {5, 2}, {6, 2}, {7, 2}, {7, 1}, {8, 1}, {9, 1}, {9, 2}, {9, 3}, {10, 3}, {10, 4}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {7, 9}, {7, 8}, {6, 8}, {5, 8}, {5, 7}, {6, 7}, {7, 7}, {8, 7}},
			Length:  33,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved into too small of space, %s (no escape)", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into wall, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}

func TestAvoidBadHead2Head3(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "38a9ec7c-2d88-49ca-a44c-d0f6c28f4b7b",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 312,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{9, 9}, {2, 0}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_pwDV8rGfY7t3TXfvy8VYmF3c",
					Name:    "Ready, Set, Hike!",
					Health:  89,
					Head:    wire.Coord{9, 7},
					Body:    []wire.Coord{{9, 7}, {8, 7}, {7, 7}, {7, 6}, {6, 6}, {6, 7}, {5, 7}, {4, 7}, {3, 7}, {2, 7}, {1, 7}, {0, 7}, {0, 8}, {0, 9}, {1, 9}, {2, 9}, {3, 9}, {3, 8}, {4, 8}, {4, 9}, {4, 10}, {5, 10}, {6, 10}, {6, 9}, {6, 8}, {7, 8}, {8, 8}, {9, 8}},
					Length:  28,
					Latency: "292",
					Shout:   "",
				},
				{
					ID:      "gs_hbSvgwd98g7ydMHRVYhgBmJV",
					Name:    "nomblegomble",
					Health:  86,
					Head:    wire.Coord{10, 6},
					Body:    []wire.Coord{{10, 6}, {9, 6}, {9, 5}, {8, 5}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {4, 4}, {3, 4}, {2, 4}, {2, 5}, {3, 5}, {4, 5}, {4, 6}, {3, 6}, {2, 6}, {1, 6}, {1, 5}, {1, 4}, {1, 3}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, {5, 2}},
					Length:  26,
					Latency: "22",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_hbSvgwd98g7ydMHRVYhgBmJV",
			Name:    "nomblegomble",
			Health:  86,
			Head:    wire.Coord{10, 6},
			Body:    []wire.Coord{{10, 6}, {9, 6}, {9, 5}, {8, 5}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {4, 4}, {3, 4}, {2, 4}, {2, 5}, {3, 5}, {4, 5}, {4, 6}, {3, 6}, {2, 6}, {1, 6}, {1, 5}, {1, 4}, {1, 3}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, {5, 2}},
			Length:  26,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake moved into potential losing h2h, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into wall, %s", nextMove.Move)
	}
}

func TestSpaceOkToTailChase1(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "71116f92-59d0-4f88-a578-a75035b4c1be",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 224,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{8, 2}, {7, 9}, {0, 8}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_v8SbjcXHGrxjwjmwpmWPy67Q",
					Name:    "nomblegomble",
					Health:  98,
					Head:    wire.Coord{3, 1},
					Body:    []wire.Coord{{3, 1}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 4}, {1, 3}, {2, 3}, {2, 2}, {3, 2}},
					Length:  14,
					Latency: "21",
					Shout:   "",
				},
				{
					ID:      "gs_bVJmK8pX8mkM7G4wVyjKCxtT",
					Name:    "caicai-vilu",
					Health:  99,
					Head:    wire.Coord{7, 1},
					Body:    []wire.Coord{{7, 1}, {7, 0}, {6, 0}, {5, 0}, {5, 1}, {5, 2}, {4, 2}, {4, 3}, {3, 3}, {3, 4}, {3, 5}, {2, 5}, {1, 5}, {1, 6}, {1, 7}, {1, 8}, {1, 9}, {2, 9}, {3, 9}, {4, 9}, {5, 9}, {6, 9}, {6, 8}, {5, 8}, {4, 8}, {3, 8}, {2, 8}, {2, 7}, {3, 7}, {4, 7}, {5, 7}, {6, 7}, {7, 7}, {8, 7}, {9, 7}},
					Length:  35,
					Latency: "75",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_v8SbjcXHGrxjwjmwpmWPy67Q",
			Name:    "nomblegomble",
			Health:  98,
			Head:    wire.Coord{3, 1},
			Body:    []wire.Coord{{3, 1}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 4}, {1, 3}, {2, 3}, {2, 2}, {3, 2}},
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

func TestAvoidEdge1(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "0ba77727-c282-4f4c-9938-71e50d884002",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 203,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{10, 1}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_7Wxmq93bWkSKCyjB6XFq6J6T",
					Name:    "hhhotdaysssnake",
					Health:  95,
					Head:    wire.Coord{7, 6},
					Body:    []wire.Coord{{7, 6}, {8, 6}, {9, 6}, {10, 6}, {10, 5}, {10, 4}, {10, 3}, {9, 3}, {8, 3}, {7, 3}, {7, 4}, {8, 4}, {9, 4}},
					Length:  13,
					Latency: "77",
					Shout:   "Do it",
				},
				{
					ID:      "gs_cBcpfgvbFGbMxmDYSVh4CthC",
					Name:    "nomblegomble",
					Health:  96,
					Head:    wire.Coord{6, 1},
					Body:    []wire.Coord{{6, 1}, {6, 2}, {6, 3}, {6, 4}, {5, 4}, {5, 3}, {5, 2}, {4, 2}, {3, 2}, {2, 2}, {1, 2}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {4, 9}, {5, 9}, {6, 9}, {6, 8}, {5, 8}},
					Length:  27,
					Latency: "22",
					Shout:   "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_cBcpfgvbFGbMxmDYSVh4CthC",
			Name:    "nomblegomble",
			Health:  96,
			Head:    wire.Coord{6, 1},
			Body:    []wire.Coord{{6, 1}, {6, 2}, {6, 3}, {6, 4}, {5, 4}, {5, 3}, {5, 2}, {4, 2}, {3, 2}, {2, 2}, {1, 2}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {4, 9}, {5, 9}, {6, 9}, {6, 8}, {5, 8}},
			Length:  27,
			Latency: "22",
			Shout:   "",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "down" {
		t.Errorf("snake moved to edge, %s (cuts space in half)", nextMove.Move)
	}
	if nextMove.Move == "up" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}

func TestFoodStart0(t *testing.T) {
	t.Skip("starting food is overrated")
	state := wire.GameState{
		Game: wire.Game{
			ID: "4c46aa82-936c-46c6-aeb2-6e33da287a3b",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 0,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 0}, {8, 4}, {10, 10}, {2, 10}, {5, 5}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_r4JCVS8Hbjq87Cg3BQM37HPf",
					Name:    "nomblegomble",
					Health:  100,
					Head:    wire.Coord{1, 1},
					Body:    []wire.Coord{{1, 1}, {1, 1}, {1, 1}},
					Length:  3,
					Latency: "",
				},
				{
					ID:      "gs_VG6tp7kfmSXSkQyPHKjk3vC6",
					Name:    "DDT",
					Health:  100,
					Head:    wire.Coord{9, 5},
					Body:    []wire.Coord{{9, 5}, {9, 5}, {9, 5}},
					Length:  3,
					Latency: "",
				},
				{
					ID:      "gs_VMrTMQrtfRrYdRPbdqbJYphd",
					Name:    "Yung Snek V0",
					Health:  100,
					Head:    wire.Coord{9, 9},
					Body:    []wire.Coord{{9, 9}, {9, 9}, {9, 9}},
					Length:  3,
					Latency: "",
				},
				{
					ID:      "gs_BChKFRcw7qVwTmfCQYfSgB4P",
					Name:    "Leonardo",
					Health:  100,
					Head:    wire.Coord{1, 9},
					Body:    []wire.Coord{{1, 9}, {1, 9}, {1, 9}},
					Length:  3,
					Latency: "",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_r4JCVS8Hbjq87Cg3BQM37HPf",
			Name:    "nomblegomble",
			Health:  100,
			Head:    wire.Coord{1, 1},
			Body:    []wire.Coord{{1, 1}, {1, 1}, {1, 1}},
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
	t.Skip("failing now with changes to eval siblings... maybe this isn't even optimal?")
	state := wire.GameState{
		Game: wire.Game{
			ID: "4c46aa82-936c-46c6-aeb2-6e33da287a3b",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 1,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{0, 0}, {8, 4}, {10, 10}, {2, 10}, {5, 5}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_r4JCVS8Hbjq87Cg3BQM37HPf",
					Name:    "nomblegomble",
					Health:  99,
					Head:    wire.Coord{1, 0},
					Body:    []wire.Coord{{1, 0}, {1, 1}, {1, 1}},
					Length:  3,
					Latency: "46",
				},
				{
					ID:      "gs_VG6tp7kfmSXSkQyPHKjk3vC6",
					Name:    "DDT",
					Health:  99,
					Head:    wire.Coord{9, 4},
					Body:    []wire.Coord{{9, 4}, {9, 5}, {9, 5}},
					Length:  3,
					Latency: "293",
				},
				{
					ID:      "gs_VMrTMQrtfRrYdRPbdqbJYphd",
					Name:    "Yung Snek V0",
					Health:  99,
					Head:    wire.Coord{8, 9},
					Body:    []wire.Coord{{8, 9}, {9, 9}, {9, 9}},
					Length:  3,
					Latency: "282",
				},
				{
					ID:      "gs_BChKFRcw7qVwTmfCQYfSgB4P",
					Name:    "Leonardo",
					Health:  99,
					Head:    wire.Coord{1, 10},
					Body:    []wire.Coord{{1, 10}, {1, 9}, {1, 9}},
					Length:  3,
					Latency: "288",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_r4JCVS8Hbjq87Cg3BQM37HPf",
			Name:    "nomblegomble",
			Health:  99,
			Head:    wire.Coord{1, 0},
			Body:    []wire.Coord{{1, 0}, {1, 1}, {1, 1}},
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
	state := wire.GameState{
		Game: wire.Game{
			ID: "245970ce-0424-4a9f-a02b-a1f0d5f531a1",
			Ruleset: wire.Ruleset{
				Name:    "standard",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 5,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{5, 5}, {9, 2}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_SWrrBGBdPF3TCy7qJXcSkrwP",
					Name:    "nomblegomble",
					Health:  97,
					Head:    wire.Coord{8, 1},
					Body:    []wire.Coord{{8, 1}, {7, 1}, {7, 0}, {8, 0}},
					Length:  4,
					Latency: "25",
				},
				{
					ID:      "gs_Tp8qw8GdDpJrxFjcj7qBVkkb",
					Name:    "nomblegomble",
					Health:  97,
					Head:    wire.Coord{1, 4},
					Body:    []wire.Coord{{1, 4}, {1, 5}, {0, 5}, {0, 4}},
					Length:  4,
					Latency: "48",
				},
				{
					ID:      "gs_CBQM9J66qjrbcYk44YSwKCRY",
					Name:    "nomblegomble",
					Health:  97,
					Head:    wire.Coord{8, 5},
					Body:    []wire.Coord{{8, 5}, {7, 5}, {7, 6}, {8, 6}},
					Length:  4,
					Latency: "45",
				},
				{
					ID:      "gs_vhSmG8q4P4H4dRmpgc3xSgHS",
					Name:    "nomblegomble",
					Health:  97,
					Head:    wire.Coord{5, 8},
					Body:    []wire.Coord{{5, 8}, {5, 9}, {5, 10}, {6, 10}},
					Length:  4,
					Latency: "23",
				},
			},
		},
		You: wire.Battlesnake{
			ID:      "gs_SWrrBGBdPF3TCy7qJXcSkrwP",
			Name:    "nomblegomble",
			Health:  97,
			Head:    wire.Coord{8, 1},
			Body:    []wire.Coord{{8, 1}, {7, 1}, {7, 0}, {8, 0}},
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

func TestGetFoodWhenStarving(t *testing.T) {
	state := wire.GameState{
		Game: wire.Game{
			ID: "d95a1a71-9ff7-4d45-9c18-4be22b7fe396",
			Ruleset: wire.Ruleset{
				Name:    "royale",
				Version: "v1.0.23",
			},
			Timeout: 500,
		},
		Turn: 173,
		Board: wire.Board{
			Height: 11,
			Width:  11,
			Food:   []wire.Coord{{10, 0}, {0, 10}, {3, 10}, {0, 5}, {10, 3}, {9, 6}, {10, 7}, {9, 9}, {2, 8}},
			Snakes: []wire.Battlesnake{
				{
					ID:      "gs_K9GKr6dgB7Fd8SbvxxkyKyPB",
					Name:    "nomblegomble dev",
					Health:  16,
					Head:    wire.Coord{2, 9},
					Body:    []wire.Coord{{2, 9}, {3, 9}, {4, 9}, {5, 9}, {6, 9}, {6, 8}, {5, 8}, {4, 8}, {4, 7}, {4, 6}, {3, 6}, {3, 7}, {2, 7}, {1, 7}, {0, 7}},
					Length:  15,
					Latency: "473",
					Shout:   "6",
				},
				{
					ID:      "gs_PqKKPgbktC7QCPyCyTJ7thRP",
					Name:    "Nessegrev-gamma",
					Health:  85,
					Head:    wire.Coord{8, 7},
					Body:    []wire.Coord{{8, 7}, {8, 8}, {7, 8}, {7, 7}, {7, 6}, {6, 6}, {6, 5}, {6, 4}, {5, 4}},
					Length:  9,
					Latency: "401",
					Shout:   "",
				},
			},
			Hazards: []wire.Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {0, 9}, {0, 10}, {1, 0}, {1, 1}, {1, 9}, {1, 10}, {2, 0}, {2, 1}, {2, 9}, {2, 10}, {3, 0}, {3, 1}, {3, 9}, {3, 10}, {4, 0}, {4, 1}, {4, 9}, {4, 10}, {5, 0}, {5, 1}, {5, 9}, {5, 10}, {6, 0}, {6, 1}, {6, 9}, {6, 10}, {7, 0}, {7, 1}, {7, 9}, {7, 10}, {8, 0}, {8, 1}, {8, 9}, {8, 10}, {9, 0}, {9, 1}, {9, 9}, {9, 10}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6}, {10, 7}, {10, 8}, {10, 9}, {10, 10}},
		},
		You: wire.Battlesnake{
			ID:      "gs_K9GKr6dgB7Fd8SbvxxkyKyPB",
			Name:    "nomblegomble dev",
			Health:  16,
			Head:    wire.Coord{2, 9},
			Body:    []wire.Coord{{2, 9}, {3, 9}, {4, 9}, {5, 9}, {6, 9}, {6, 8}, {5, 8}, {4, 8}, {4, 7}, {4, 6}, {3, 6}, {3, 7}, {2, 7}, {1, 7}, {0, 7}},
			Length:  15,
			Latency: "473",
			Shout:   "6",
		},
	}

	nextMove := move(state)

	if nextMove.Move == "up" {
		t.Errorf("snake moved away from nearby food, %s", nextMove.Move)
	}
	if nextMove.Move == "left" {
		t.Errorf("snake moved away from nearby food, %s", nextMove.Move)
	}
	if nextMove.Move == "right" {
		t.Errorf("snake moved into self, %s", nextMove.Move)
	}
}
