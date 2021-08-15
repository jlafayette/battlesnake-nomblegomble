package main

import (
	"fmt"
	"testing"
)

func TestGrid(t *testing.T) {
	fmt.Printf("starting")

	you := Battlesnake{
		Head: Coord{2, 0},
		Body: []Coord{{2, 0}, {2, 1}, {1, 1}, {0, 1}, {0, 0}},
	}

	state := GameState{
		Board: Board{
			Snakes: []Battlesnake{you},
			Width:  4,
			Height: 4,
			Food:   []Coord{{1, 0}, {0, 0}}, // 0,0 simulates bad luck of food spawning
		},
		You: you,
	}

	grid := NewGrid(&state)
	leftArea := grid.Area(&state, "left")
	grid = NewGrid(&state)
	rightArea := grid.Area(&state, "right")

	if leftArea != 2 {
		t.Errorf("expected 2, got %d", leftArea)
	}
	if rightArea != 16 {
		t.Errorf("expected 16, got %d", rightArea)
	}
}

func TestGrid2(t *testing.T) {
	fmt.Println("starting TestGrid2")

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

	// because tail is here, it counts as infinate space
	if upArea != 49 {
		t.Errorf("expected 49, got %d", upArea)
	}
	if downArea != 15 {
		t.Errorf("expected 15, got %d", downArea)
	}

}

func TestGrid3(t *testing.T) {

	state := GameState{
		Game: Game{
			ID: "43172677-aa69-4a04-aecc-4aedcf238d05",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 144,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{9, 9}, {2, 7}, {5, 3}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_K96GMhmm4XqSJDgbDmfdSv3J",
					Name:    "nomblegomble",
					Health:  89,
					Head:    Coord{2, 4},
					Body:    []Coord{{2, 4}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, {5, 4}, {6, 4}, {7, 4}, {8, 4}, {9, 4}, {10, 4}, {10, 5}, {9, 5}, {8, 5}, {7, 5}, {6, 5}},
					Length:  16,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_MF6b9fcWTpS9FRTCVJMK88r4",
					Name:    "Super Snakey",
					Health:  95,
					Head:    Coord{3, 3},
					Body:    []Coord{{3, 3}, {2, 3}, {2, 2}, {2, 1}, {1, 1}, {1, 0}, {2, 0}, {3, 0}, {3, 1}, {3, 2}, {4, 2}, {5, 2}, {6, 2}, {7, 2}, {7, 1}},
					Length:  15,
					Latency: "226",
					Shout:   "",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_K96GMhmm4XqSJDgbDmfdSv3J",
			Name:    "nomblegomble",
			Health:  89,
			Head:    Coord{2, 4},
			Body:    []Coord{{2, 4}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, {5, 4}, {6, 4}, {7, 4}, {8, 4}, {9, 4}, {10, 4}, {10, 5}, {9, 5}, {8, 5}, {7, 5}, {6, 5}},
			Length:  16,
			Latency: "22",
			Shout:   "",
		},
	}

	grid := NewGrid(&state)
	rightArea := grid.Area(&state, "right")

	// Because the other snake can cut this off, should only count as 2
	if rightArea != 2 {
		t.Errorf("expected 2, got %d", rightArea)
	}

}

func TestGrid4(t *testing.T) {

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

	downArea := GetArea(&state, "down")

	// Because we can beath the other snake in a H2H, this area should be fairly large
	if downArea != 40 {
		t.Errorf("expected 40, got %d", downArea)
	}

}
