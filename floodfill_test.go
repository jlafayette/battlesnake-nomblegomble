package main

import (
	"testing"

	tt "github.com/jlafayette/battlesnake-go/t"
)

var result map[SnakeIndex]*FloodFillResult

func BenchmarkMap1(b *testing.B) {
	you := tt.Battlesnake{
		Head:   tt.Coord{2, 0},
		Body:   []tt.Coord{{2, 0}, {2, 1}, {1, 1}, {0, 1}, {0, 0}},
		Length: 5,
	}

	state := tt.GameState{
		Board: tt.Board{
			Snakes: []tt.Battlesnake{you},
			Width:  4,
			Height: 4,
			Food:   []tt.Coord{{1, 0}},
		},
		You: you,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		board := NewBoard(&state)
		result = board.Fill()
	}
}

func TestFloodFill01(t *testing.T) {

	you := tt.Battlesnake{
		Head:   tt.Coord{2, 0},
		Body:   []tt.Coord{{2, 0}, {2, 1}, {1, 1}, {0, 1}, {0, 0}},
		Length: 5,
	}

	state := tt.GameState{
		Board: tt.Board{
			Snakes: []tt.Battlesnake{you},
			Width:  4,
			Height: 4,
			Food:   []tt.Coord{{1, 0}},
		},
		You: you,
	}

	board := NewBoard(&state)
	result := board.Fill()

	r, ok := result[0]
	if !ok {
		t.Errorf("expected snakeIndex 0 to exist in result map")
	}
	if r.Area != 16 {
		t.Errorf("expected 16, got %d", r.Area)
	}
	if r.Food != 1 {
		t.Errorf("expected 1 food, got %d", r.Food)
	}
}
