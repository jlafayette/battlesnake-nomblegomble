package main

import (
	"testing"
)

func TestGrid1(t *testing.T) {

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

	leftArea := GetArea(&state, "left")
	rightArea := GetArea(&state, "right")

	if leftArea.Space != 1 {
		t.Errorf("expected 1, got %d", leftArea.Space)
	}
	if rightArea.Space != 16 {
		t.Errorf("expected 16, got %d", rightArea.Space)
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

	upArea := GetArea(&state, "up")
	downArea := GetArea(&state, "down")

	// because tail is here, it counts as infinite space (49)
	// Or.. at least that's how it should be, with the current calculation it
	// doesn't work this way. That's because multiple neighbors are selected
	// each turn, so the tail doesn't have as much time to shrink as it would
	// in the real game.
	if upArea.Space != 18 {
		t.Errorf("expected 18, got %d", upArea.Space)
	}
	if downArea.Space != 15 {
		t.Errorf("expected 15, got %d", downArea.Space)
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

	rightArea := GetArea(&state, "right")

	// Because the other snake can cut this off, should only count as 2
	if rightArea.Space != 2 {
		t.Errorf("expected 2, got %d", rightArea.Space)
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
	upArea := GetArea(&state, "up")

	// Because we can beat the other snake in a H2H, this area should be fairly large
	if downArea.Space < 29 {
		t.Errorf("expected 29 or more, got %d", downArea.Space)
	}

	if upArea.Space != 8 {
		t.Errorf("expected 8, got %d", upArea.Space)
	}

}

func TestGrid5(t *testing.T) {
	state := GameState{
		Game: Game{
			ID: "eca2463d-0fd7-43b7-aa6b-43dbb489da07",
			Ruleset: Ruleset{
				Name:    "standard",
				Version: "v1.0.17",
			},
			Timeout: 500,
		},
		Turn: 50,
		Board: Board{
			Height: 11,
			Width:  11,
			Food:   []Coord{{4, 0}},
			Snakes: []Battlesnake{
				{
					ID:      "gs_kpRwFYKwVjmj7JF6RTwdPHBB",
					Name:    "nomblegomble",
					Health:  99,
					Head:    Coord{9, 7},
					Body:    []Coord{{9, 7}, {9, 8}, {8, 8}, {7, 8}, {6, 8}, {6, 7}, {5, 7}, {4, 7}, {4, 6}},
					Length:  9,
					Latency: "22",
					Shout:   "",
				},
				{
					ID:      "gs_tC8WtyKcvjkvyQhVSB977YR9",
					Name:    "The Very Hungry Caterpillar ≡ƒìè≡ƒìÅ≡ƒìæ≡ƒìÆ≡ƒìÄ≡ƒÉ¢",
					Health:  95,
					Head:    Coord{8, 6},
					Body:    []Coord{{8, 6}, {9, 6}, {9, 5}, {9, 4}, {10, 4}, {10, 3}, {9, 3}},
					Length:  7,
					Latency: "40",
					Shout:   "",
				},
			},
		},
		You: Battlesnake{
			ID:      "gs_kpRwFYKwVjmj7JF6RTwdPHBB",
			Name:    "nomblegomble",
			Health:  99,
			Head:    Coord{9, 7},
			Body:    []Coord{{9, 7}, {9, 8}, {8, 8}, {7, 8}, {6, 8}, {6, 7}, {5, 7}, {4, 7}, {4, 6}},
			Length:  9,
			Latency: "22",
			Shout:   "",
		},
	}

	area := GetArea(&state, "left")

	// Because the other snake can cut this off, should only count as 2
	if area.Space != 2 {
		t.Errorf("expected 2, got %d", area.Space)
	}
}
