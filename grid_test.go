package main

import "testing"

func TestGrid(t *testing.T) {

	you := Battlesnake{
		Head: Coord{2, 0},
		Body: []Coord{{2, 0}, {2, 1}, {1, 1}, {0, 1}, {0, 0}},
	}

	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{you},
			Width:  4,
			Height: 4,
		},
		You: you,
	}

	grid := NewGrid(4, 4)
	leftArea := grid.Area(&state, "left")
	rightArea := grid.Area(&state, "right")

	if leftArea != 1 {
		t.Errorf("expected 1, got %d", leftArea)
	}
	if rightArea != 10 {
		t.Errorf("expected 10, got %d", rightArea)
	}
}
