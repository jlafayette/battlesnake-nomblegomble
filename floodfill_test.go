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
	if r.Area < 16 {
		t.Errorf("expected at least 16, got %d", r.Area)
	}
	if r.Food != 1 {
		t.Errorf("expected 1 food, got %d", r.Food)
	}
}

func TestFloodFill02(t *testing.T) {

	you := tt.Battlesnake{
		Head:   tt.Coord{1, 0},
		Body:   []tt.Coord{{1, 0}, {1, 1}, {0, 1}},
		Length: 3,
	}
	other := tt.Battlesnake{
		Head:   tt.Coord{2, 3},
		Body:   []tt.Coord{{2, 3}, {3, 3}, {3, 2}},
		Length: 3,
	}

	state := tt.GameState{
		Board: tt.Board{
			Snakes: []tt.Battlesnake{you, other},
			Width:  4,
			Height: 4,
			Food:   []tt.Coord{{0, 0}},
		},
		You: you,
	}

	board := NewBoard(&state)
	result := board.Fill()

	r0, ok0 := result[0]
	r1, ok1 := result[1]
	if !ok0 || !ok1 {
		t.Errorf("expected snakeIndex 0 and 1 to exist in result map")
	}
	if r0.Area < 7 {
		t.Errorf("expected at least 7, got %d", r0.Area)
	}
	if r0.Food != 1 {
		t.Errorf("expected 1 food, got %d", r0.Food)
	}
	if r1.Area < 7 {
		t.Errorf("expected at least 7, got %d", r1.Area)
	}
	if r1.Food != 0 {
		t.Errorf("expected 0 food, got %d", r1.Food)
	}
}
