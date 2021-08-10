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
