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

	grid := NewGrid(&state)
	leftArea := grid.Area(&state, "left")
	grid = NewGrid(&state)
	rightArea := grid.Area(&state, "right")

	if leftArea != 1 {
		t.Errorf("expected 1, got %d", leftArea)
	}
	if rightArea != 10 {
		t.Errorf("expected 10, got %d", rightArea)
	}
}

func TestGrid2(t *testing.T) {

	state := GameState{
		Game: Game{
			ID: "7e4b5f59-2e60-473b-8646-a3ce36371189",
			Ruleset: Ruleset{
				Name:    "solo",
				Version: "",
			},
			Timeout: 500,
		},
		Turn: 128,
		Board: Board{
			Height: 7,
			Width:  7,
			Food:   []Coord{{1, 0}, {0, 0}, {1, 3}, {0, 3}, {6, 3}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_wTThGxrtjvPjS7YQ77hFkdFc",
					Name:    "nomblegomble",
					Health:  98,
					Head:    Coord{0, 2},
					Body:    []Coord{{0, 2}, {1, 2}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 5}, {4, 5}, {4, 6}, {3, 6}, {2, 6}, {1, 6}, {0, 6}, {0, 5}, {0, 4}, {1, 4}, {1, 5}, {2, 5}, {3, 5}, {3, 4}, {4, 4}, {4, 3}, {3, 3}},
					Length:  27,
					Latency: "20",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_wTThGxrtjvPjS7YQ77hFkdFc",
			Name:    "nomblegomble",
			Health:  98,
			Head:    Coord{0, 2},
			Body:    []Coord{{0, 2}, {1, 2}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {5, 2}, {5, 3}, {5, 4}, {5, 5}, {4, 5}, {4, 6}, {3, 6}, {2, 6}, {1, 6}, {0, 6}, {0, 5}, {0, 4}, {1, 4}, {1, 5}, {2, 5}, {3, 5}, {3, 4}, {4, 4}, {4, 3}, {3, 3}},
			Length:  27,
			Latency: "20",
		},
	}

	grid := NewGrid(&state)
	upArea := grid.Area(&state, "up")
	grid2 := NewGrid(&state)
	downArea := grid2.Area(&state, "down")

	if upArea != 7 {
		t.Errorf("expected 7, got %d", upArea)
	}
	if downArea != 15 {
		t.Errorf("expected 15, got %d", downArea)
	}
}
