package tree

import (
	"reflect"
	"testing"
)

func TestEating01(t *testing.T) {
	// Example of eating in normal gameplay
	// {"move":"left"}
	// 19
	// health: 83,  body: [{6,7},{7,7},{7,6},{7,5}],  head: {6,7},  length: 4
	// {"move":"up"}
	// 20
	// health: 82,  body: [{6,8},{6,7},{7,7},{7,6}],  head: {6,8},  length: 4
	// {"move":"up"}
	// 21
	// health: 100,  body: [{6,9},{6,8},{6,7},{7,7},{7,7}],  head: {6,9},  length: 5
	// {"move":"left"}
	// 22
	// health: 99,  body: [{5,9},{6,9},{6,8},{6,7},{7,7}],  head: {5,9},  length: 5
	// {"move":"down"}
	// 23
	// health: 98,  body: [{5,8},{5,9},{6,9},{6,8},{6,7}],  head: {5,8},  length: 5
	// {"move":"down"}

	snake := NewSnake(0, 83, []Coord{{6, 7}, {7, 7}, {7, 6}, {7, 5}}, 19, 4)

	snake.Move(Up, false, false, false, 15)
	snake.Move(Up, true, false, false, 15)
	snake.Move(Left, false, false, false, 15)
	snake.Move(Down, false, false, false, 15)
	// snake.Move(Down, false)  // panics because depth=4

	if snake.Length != 5 {
		t.Errorf("expected Length to be 5, got %d", snake.Length)
	}
	if snake.turn != 23 {
		t.Errorf("expected turn to be 23, got %d", snake.turn)
	}

	expectedBody := []Coord{{5, 8}, {5, 9}, {6, 9}, {6, 8}, {6, 7}}
	if !reflect.DeepEqual(snake.Body, expectedBody) {
		t.Errorf("expected body to be %v, got %v", expectedBody, snake.Body)
	}

	snake.UndoMove()
	snake.UndoMove()
	snake.UndoMove()
	snake.UndoMove()

	if snake.Length != 4 {
		t.Errorf("expected Length to be 4, got %d", snake.Length)
	}
	if snake.turn != 19 {
		t.Errorf("expected turn to be 19, got %d", snake.turn)
	}

	expectedBody = []Coord{{6, 7}, {7, 7}, {7, 6}, {7, 5}}
	if !reflect.DeepEqual(snake.Body, expectedBody) {
		t.Errorf("expected body to be %v, got %v", expectedBody, snake.Body)
	}
}

func TestEating02(t *testing.T) {
	// Example of eating from start position (already doubled)
	// 0
	// health:100, body: [{9,1},{9,1},{9,1}],  head: {9,1},  length: 3
	// {"move":"right"}
	// 1
	// health:99, body: [{10:1},{9,1},{9,1}],  head: {10,1},  length: 3
	// {"move":"down"}
	// 2
	// health:100, body: [{10:0},{10,1},{9,1},{9,1}],  head: {10,0},  length: 4
	// {"move":"left"}
	// 3
	// health:99, body: [{9,0},{10,0},{10,1},{9,1}],  head: {9,0},  length: 4
	// {"move":"up"}
	// 4
	// health:98, body: [{9,1},{9,0},{10,0},{10,1}],  head: {9,1},  length: 4
	// {"move":"up"}
	// 5
	// health:97, body: [{9,2},{9,1},{9,0},{10,0}],  head: {9,2},  length: 4
	// {"move":"left"}

	snake := NewSnake(0, 100, []Coord{{9, 1}, {9, 1}, {9, 1}}, 0, 5)

	snake.Move(Right, false, false, false, 15)
	snake.Move(Down, true, false, false, 15)
	snake.Move(Left, false, false, false, 15)
	snake.Move(Up, false, false, false, 15)
	snake.Move(Up, false, false, false, 15)

	if snake.Length != 4 {
		t.Errorf("expected Length to be 4, got %d", snake.Length)
	}
	if snake.turn != 5 {
		t.Errorf("expected turn to be 5, got %d", snake.turn)
	}

	expectedBody := []Coord{{9, 2}, {9, 1}, {9, 0}, {10, 0}}
	if !reflect.DeepEqual(snake.Body, expectedBody) {
		t.Errorf("expected body to be %v, got %v", expectedBody, snake.Body)
	}

	snake.UndoMove()
	snake.UndoMove()
	snake.UndoMove()
	snake.UndoMove()
	snake.UndoMove()

	if snake.Length != 3 {
		t.Errorf("expected Length to be 3, got %d", snake.Length)
	}
	if snake.turn != 0 {
		t.Errorf("expected turn to be 0, got %d", snake.turn)
	}

	expectedBody = []Coord{{9, 1}, {9, 1}, {9, 1}}
	if !reflect.DeepEqual(snake.Body, expectedBody) {
		t.Errorf("expected body to be %v, got %v", expectedBody, snake.Body)
	}
}

func TestEating03(t *testing.T) {
	// Example of eating two in a row
	// 93
	// health: 86  body: [{4,3},{4,2},{4,1},{4,0},{3,0},{2,0},{2,1},{3,1},{3,2},{2,2},{1,2},{1,1},{1,0},{0,0},{0,1},{0,2},{0,3}]  head: {4,3}  length: 17
	// {"move":"up"}
	// 94
	// health: 85  body: [{4,4},{4,3},{4,2},{4,1},{4,0},{3,0},{2,0},{2,1},{3,1},{3,2},{2,2},{1,2},{1,1},{1,0},{0,0},{0,1},{0,2}]  head: {4,4}  length: 17
	// {"move":"left"}
	// 95
	// health: 100  body: [{3,4},{4,4},{4,3},{4,2},{4,1},{4,0},{3,0},{2,0},{2,1},{3,1},{3,2},{2,2},{1,2},{1,1},{1,0},{0,0},{0,1},{0,1}]  head: {3,4}  length: 18
	// {"move":"left"}
	// 96
	// health: 100  body: [{2,4},{3,4},{4,4},{4,3},{4,2},{4,1},{4,0},{3,0},{2,0},{2,1},{3,1},{3,2},{2,2},{1,2},{1,1},{1,0},{0,0},{0,1},{0,1}]  head: {2,4}  length: 19
	// {"move":"left"}
	// 97
	// health: 99  body: [{1,4},{2,4},{3,4},{4,4},{4,3},{4,2},{4,1},{4,0},{3,0},{2,0},{2,1},{3,1},{3,2},{2,2},{1,2},{1,1},{1,0},{0,0},{0,1}]  head: {1,4}  length: 19
	// {"move":"left"}
	// 98
	// health: 100  body: [{0,4},{1,4},{2,4},{3,4},{4,4},{4,3},{4,2},{4,1},{4,0},{3,0},{2,0},{2,1},{3,1},{3,2},{2,2},{1,2},{1,1},{1,0},{0,0},{0,0}]  head: {0,4}  length: 20
	// {"move":"up"}
	// 99
	// health: 99  body: [{0,5},{0,4},{1,4},{2,4},{3,4},{4,4},{4,3},{4,2},{4,1},{4,0},{3,0},{2,0},{2,1},{3,1},{3,2},{2,2},{1,2},{1,1},{1,0},{0,0}]  head: {0,5}  length: 20
	// {"move":"right"}
	// 100
	// health: 98  body: [{1,5},{0,5},{0,4},{1,4},{2,4},{3,4},{4,4},{4,3},{4,2},{4,1},{4,0},{3,0},{2,0},{2,1},{3,1},{3,2},{2,2},{1,2},{1,1},{1,0}]  head: {1,5}  length: 20

	snake := NewSnake(0, 83, []Coord{{4, 3}, {4, 2}, {4, 1}, {4, 0}, {3, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 2}, {2, 2}, {1, 2}, {1, 1}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}}, 93, 7)

	snake.Move(Up, false, false, false, 15)
	snake.Move(Left, true, false, false, 15)
	snake.Move(Left, true, false, false, 15)
	snake.Move(Left, false, false, false, 15)
	snake.Move(Left, true, false, false, 15)
	snake.Move(Up, false, false, false, 15)
	snake.Move(Right, false, false, false, 15)

	if snake.Length != 20 {
		t.Errorf("expected Length to be 20, got %d", snake.Length)
	}
	if snake.turn != 100 {
		t.Errorf("expected turn to be 100, got %d", snake.turn)
	}

	// [{1 5} {0 5} {0 4} {1 4} {2 4} {3 4} {4 4} {4 3} {4 2} {4 1} {4 0} {3 0} {2 0} {2 1} {3 1} {3 2} {2 2} {1 2} {1 1} {1 0}]
	// [{1 5} {0 5} {0 4} {1 4} {2 4} {3 4} {4 4} {4 3} {4 2} {4 1} {4 0} {3 0} {2 0} {2 1} {3 1} {3 2} {2 2} {1 2} {1 1}]

	expectedBody := []Coord{{1, 5}, {0, 5}, {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 3}, {4, 2}, {4, 1}, {4, 0}, {3, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 2}, {2, 2}, {1, 2}, {1, 1}, {1, 0}}
	if !reflect.DeepEqual(snake.Body, expectedBody) {
		t.Errorf("expected body to be %v, got %v", expectedBody, snake.Body)
	}

	snake.UndoMove()
	snake.UndoMove()
	snake.UndoMove()
	snake.UndoMove()
	snake.UndoMove()
	snake.UndoMove()
	snake.UndoMove()

	if snake.Length != 17 {
		t.Errorf("expected Length to be 17, got %d", snake.Length)
	}
	if snake.turn != 93 {
		t.Errorf("expected turn to be 93, got %d", snake.turn)
	}

	// [{4 3} {4 2} {4 1} {4 0} {3 0} {2 0} {2 1} {3 1} {3 2} {2 2} {1 2} {1 1} {1 0} {0 0} {0 1} {0 2} {0 3}]
	// [{4 3} {4 2} {4 1} {4 0} {3 0} {2 0} {2 1} {3 1} {3 2} {2 2} {1 2} {1 1} {1 0} {0 0} {0 1} {0 2}]
	expectedBody = []Coord{{4, 3}, {4, 2}, {4, 1}, {4, 0}, {3, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 2}, {2, 2}, {1, 2}, {1, 1}, {1, 0}, {0, 0}, {0, 1}, {0, 2}, {0, 3}}
	if !reflect.DeepEqual(snake.Body, expectedBody) {
		t.Errorf("expected body to be %v, got %v", expectedBody, snake.Body)
	}
}

func TestDead01(t *testing.T) {

	snake := NewSnake(0, 83, []Coord{{4, 3}, {4, 2}, {4, 1}}, 1, 5)

	// ddd {1, 4}, {2, 4}, {3, 4}, {4, 4}, {4, 3}, {4, 2}, {4, 1}
	snake.Move(Up, false, false, false, 15)
	// fmt.Printf("Body %v\n", snake.Body)
	snake.Move(Left, false, false, false, 15)
	// fmt.Printf("Body %v\n", snake.Body)
	snake.Move(Left, false, false, false, 15)
	// fmt.Printf("Body %v\n", snake.Body)
	snake.Move(Left, false, true, false, 15)
	// fmt.Printf("Body %v\n", snake.Body)
	snake.Move(Dead, false, false, false, 15)
	// fmt.Printf("Body %v\n", snake.Body)
	snake.Move(Dead, false, false, false, 15)
	// fmt.Printf("Body %v\n", snake.Body)
	snake.Move(Dead, false, false, false, 15)
	// fmt.Printf("Body %v\n", snake.Body)

	if snake.Length != 3 {
		t.Errorf("expected Length to be 3, got %d", snake.Length)
	}
	if snake.turn != 8 {
		t.Errorf("expected turn to be 8, got %d", snake.turn)
	}
	expectedBody := []Coord{{1, 4}, {2, 4}, {3, 4}}
	if !reflect.DeepEqual(snake.Body, expectedBody) {
		t.Errorf("expected body to be %v, got %v", expectedBody, snake.Body)
	}

	_, err := snake.UndoMove()
	if err != nil {
		t.Errorf("after undo got error %s", err)
	}
	_, err = snake.UndoMove()
	if err != nil {
		t.Errorf("after undo got error %s", err)
	}
	_, err = snake.UndoMove()
	if err != nil {
		t.Errorf("after undo got error %s", err)
	}
	// should be at dead position here
	expectedBody = []Coord{{1, 4}, {2, 4}, {3, 4}}
	if !reflect.DeepEqual(snake.Body, expectedBody) {
		t.Errorf("expected body to be %v, got %v", expectedBody, snake.Body)
	}
	_, err = snake.UndoMove()
	if err != nil {
		t.Errorf("after undo got error %s", err)
	}
	// should be at 1 before dead position here
	expectedBody = []Coord{{2, 4}, {3, 4}, {4, 4}}
	if !reflect.DeepEqual(snake.Body, expectedBody) {
		t.Errorf("expected body to be %v, got %v", expectedBody, snake.Body)
	}

	_, err = snake.UndoMove()
	if err != nil {
		t.Errorf("after undo got error %s", err)
	}
	_, err = snake.UndoMove()
	if err != nil {
		t.Errorf("after undo got error %s", err)
	}
	_, err = snake.UndoMove()
	if err != nil {
		t.Errorf("after undo got error %s", err)
	}

	if snake.Length != 3 {
		t.Errorf("after undos expected Length to be 3, got %d", snake.Length)
	}
	if snake.turn != 1 {
		t.Errorf("after undos expected turn to be 1, got %d", snake.turn)
	}

	expectedBody = []Coord{{4, 3}, {4, 2}, {4, 1}}
	if !reflect.DeepEqual(snake.Body, expectedBody) {
		t.Errorf("after undos expected body to be %v, got %v", expectedBody, snake.Body)
	}
}
