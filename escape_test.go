package main

// func TestEscape1_287(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 287,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 5}, {9, 1}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  98,
// 					Head:    Coord{9, 0},
// 					Body:    []Coord{{9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {4, 2}, {4, 3}, {4, 4}, {4, 5}},
// 					Length:  33,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  90,
// 					Head:    Coord{5, 8},
// 					Body:    []Coord{{5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}, {3, 9}, {3, 8}, {3, 7}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}},
// 					Length:  24,
// 					Latency: "87",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  98,
// 			Head:    Coord{9, 0},
// 			Body:    []Coord{{9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {4, 2}, {4, 3}, {4, 4}, {4, 5}},
// 			Length:  33,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into wall, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake is wasting space, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_288(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 288,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 5}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  100,
// 					Head:    Coord{9, 1},
// 					Body:    []Coord{{9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {4, 2}, {4, 3}, {4, 4}, {4, 4}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  89,
// 					Head:    Coord{5, 7},
// 					Body:    []Coord{{5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}, {3, 9}, {3, 8}, {3, 7}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}},
// 					Length:  24,
// 					Latency: "89",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  100,
// 			Head:    Coord{9, 1},
// 			Body:    []Coord{{9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {4, 2}, {4, 3}, {4, 4}, {4, 4}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake is wasting space, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_289(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 289,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 5}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  99,
// 					Head:    Coord{9, 2},
// 					Body:    []Coord{{9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {4, 2}, {4, 3}, {4, 4}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  88,
// 					Head:    Coord{6, 7},
// 					Body:    []Coord{{6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}, {3, 9}, {3, 8}, {3, 7}, {3, 6}, {4, 6}, {5, 6}, {6, 6}},
// 					Length:  24,
// 					Latency: "85",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  99,
// 			Head:    Coord{9, 2},
// 			Body:    []Coord{{9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {4, 2}, {4, 3}, {4, 4}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake is wasting space, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_290(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 290,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 5}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  98,
// 					Head:    Coord{9, 3},
// 					Body:    []Coord{{9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {4, 2}, {4, 3}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  87,
// 					Head:    Coord{6, 6},
// 					Body:    []Coord{{6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}, {3, 9}, {3, 8}, {3, 7}, {3, 6}, {4, 6}, {5, 6}},
// 					Length:  24,
// 					Latency: "89",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  98,
// 			Head:    Coord{9, 3},
// 			Body:    []Coord{{9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {4, 2}, {4, 3}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "up" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_291(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 291,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 5}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  97,
// 					Head:    Coord{8, 3},
// 					Body:    []Coord{{8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {4, 2}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  86,
// 					Head:    Coord{5, 6},
// 					Body:    []Coord{{5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}, {3, 9}, {3, 8}, {3, 7}, {3, 6}, {4, 6}},
// 					Length:  24,
// 					Latency: "89",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  97,
// 			Head:    Coord{8, 3},
// 			Body:    []Coord{{8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {4, 2}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "left" {
// 		t.Errorf("snake is wasting space, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_292(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 292,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 5}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  96,
// 					Head:    Coord{8, 4},
// 					Body:    []Coord{{8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  85,
// 					Head:    Coord{4, 6},
// 					Body:    []Coord{{4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}, {3, 9}, {3, 8}, {3, 7}, {3, 6}},
// 					Length:  24,
// 					Latency: "83",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  96,
// 			Head:    Coord{8, 4},
// 			Body:    []Coord{{8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "up" {
// 		t.Errorf("snake moved into too small of space, %s (no escape)", nextMove.Move)
// 	}
// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_293(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 293,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 5}, {1, 3}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  95,
// 					Head:    Coord{7, 4},
// 					Body:    []Coord{{7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}},
// 					Length:  34,
// 					Latency: "22",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  84,
// 					Head:    Coord{3, 6},
// 					Body:    []Coord{{3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}, {3, 9}, {3, 8}, {3, 7}},
// 					Length:  24,
// 					Latency: "99",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  95,
// 			Head:    Coord{7, 4},
// 			Body:    []Coord{{7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}},
// 			Length:  34,
// 			Latency: "22",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake is wasting space, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "up" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// }

// 'correct' moves are not set after this point

// func TestEscape1_294(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 294,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 5}, {1, 3}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  94,
// 					Head:    Coord{7, 3},
// 					Body:    []Coord{{7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  83,
// 					Head:    Coord{2, 6},
// 					Body:    []Coord{{2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}, {3, 9}, {3, 8}},
// 					Length:  24,
// 					Latency: "114",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  94,
// 			Head:    Coord{7, 3},
// 			Body:    []Coord{{7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}, {2, 1}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into wall, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_295(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 295,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 5}, {1, 3}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  93,
// 					Head:    Coord{7, 2},
// 					Body:    []Coord{{7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  82,
// 					Head:    Coord{2, 5},
// 					Body:    []Coord{{2, 5}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}, {3, 9}},
// 					Length:  24,
// 					Latency: "98",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  93,
// 			Head:    Coord{7, 2},
// 			Body:    []Coord{{7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}, {1, 1}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into wall, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_296(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 296,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 3}, {10, 6}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  92,
// 					Head:    Coord{7, 1},
// 					Body:    []Coord{{7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  100,
// 					Head:    Coord{1, 5},
// 					Body:    []Coord{{1, 5}, {2, 5}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}, {4, 9}},
// 					Length:  25,
// 					Latency: "106",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  92,
// 			Head:    Coord{7, 1},
// 			Body:    []Coord{{7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}, {0, 1}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into wall, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_297(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 297,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 3}, {10, 6}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  91,
// 					Head:    Coord{7, 0},
// 					Body:    []Coord{{7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  99,
// 					Head:    Coord{1, 4},
// 					Body:    []Coord{{1, 4}, {1, 5}, {2, 5}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}, {4, 9}},
// 					Length:  25,
// 					Latency: "115",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  91,
// 			Head:    Coord{7, 0},
// 			Body:    []Coord{{7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}, {0, 0}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into wall, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_298(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 298,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 3}, {10, 6}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  90,
// 					Head:    Coord{6, 0},
// 					Body:    []Coord{{6, 0}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  98,
// 					Head:    Coord{0, 4},
// 					Body:    []Coord{{0, 4}, {1, 4}, {1, 5}, {2, 5}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}, {5, 9}},
// 					Length:  25,
// 					Latency: "95",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  90,
// 			Head:    Coord{6, 0},
// 			Body:    []Coord{{6, 0}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}, {1, 0}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into wall, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_299(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 299,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 3}, {10, 6}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  89,
// 					Head:    Coord{6, 1},
// 					Body:    []Coord{{6, 1}, {6, 0}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  97,
// 					Head:    Coord{0, 5},
// 					Body:    []Coord{{0, 5}, {0, 4}, {1, 4}, {1, 5}, {2, 5}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}, {6, 9}},
// 					Length:  25,
// 					Latency: "98",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  89,
// 			Head:    Coord{6, 1},
// 			Body:    []Coord{{6, 1}, {6, 0}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}, {2, 0}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into wall, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_300(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 300,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 3}, {10, 6}, {3, 8}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  88,
// 					Head:    Coord{6, 2},
// 					Body:    []Coord{{6, 2}, {6, 1}, {6, 0}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  96,
// 					Head:    Coord{0, 6},
// 					Body:    []Coord{{0, 6}, {0, 5}, {0, 4}, {1, 4}, {1, 5}, {2, 5}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}, {6, 10}},
// 					Length:  25,
// 					Latency: "84",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  88,
// 			Head:    Coord{6, 2},
// 			Body:    []Coord{{6, 2}, {6, 1}, {6, 0}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}, {3, 0}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into wall, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_301(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 301,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 3}, {10, 6}, {3, 8}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  87,
// 					Head:    Coord{6, 3},
// 					Body:    []Coord{{6, 3}, {6, 2}, {6, 1}, {6, 0}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  95,
// 					Head:    Coord{1, 6},
// 					Body:    []Coord{{1, 6}, {0, 6}, {0, 5}, {0, 4}, {1, 4}, {1, 5}, {2, 5}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10}},
// 					Length:  25,
// 					Latency: "302",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  87,
// 			Head:    Coord{6, 3},
// 			Body:    []Coord{{6, 3}, {6, 2}, {6, 1}, {6, 0}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}, {4, 0}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into wall, %s", nextMove.Move)
// 	}
// }

// func TestEscape1_302(t *testing.T) {
// 	state := GameState{
// 		Game: Game{
// 			ID: "9092a80c-02da-40c9-8531-25b7b4fe11ac",
// 			Ruleset: Ruleset{
// 				Name:    "standard",
// 				Version: "v1.0.20",
// 			},
// 			Timeout: 500,
// 		},
// 		Turn: 302,
// 		Board: Board{
// 			Height: 11,
// 			Width:  11,
// 			Food:   []Coord{{1, 3}, {10, 6}, {3, 8}},
// 			Snakes: []Battlesnake{
// 				{
// 					ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 					Name:    "nomblegomble",
// 					Health:  86,
// 					Head:    Coord{6, 4},
// 					Body:    []Coord{{6, 4}, {6, 3}, {6, 2}, {6, 1}, {6, 0}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}},
// 					Length:  34,
// 					Latency: "21",
// 					Shout:   "",
// 				},
// 				{
// 					ID:      "gs_C7gqcRgkjJk47ywFhRTbmH34",
// 					Name:    "Untimely Neglected Wearable",
// 					Health:  94,
// 					Head:    Coord{1, 7},
// 					Body:    []Coord{{1, 7}, {1, 6}, {0, 6}, {0, 5}, {0, 4}, {1, 4}, {1, 5}, {2, 5}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {6, 7}, {5, 7}, {5, 8}, {6, 8}, {7, 8}, {8, 8}, {8, 9}, {9, 9}, {10, 9}, {10, 10}, {9, 10}, {8, 10}},
// 					Length:  25,
// 					Latency: "78",
// 					Shout:   "",
// 				},
// 			},
// 		},
// 		You: Battlesnake{
// 			ID:      "gs_8hbpfDWwhCCVYThYMjVgbPSR",
// 			Name:    "nomblegomble",
// 			Health:  86,
// 			Head:    Coord{6, 4},
// 			Body:    []Coord{{6, 4}, {6, 3}, {6, 2}, {6, 1}, {6, 0}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {8, 4}, {8, 3}, {9, 3}, {9, 2}, {9, 1}, {9, 0}, {10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {9, 4}, {9, 5}, {9, 6}, {8, 6}, {7, 6}, {7, 5}, {6, 5}, {5, 5}, {5, 4}, {5, 3}, {5, 2}, {5, 1}, {5, 0}},
// 			Length:  34,
// 			Latency: "21",
// 			Shout:   "",
// 		},
// 	}

// 	nextMove := move(state)

// 	if nextMove.Move == "down" {
// 		t.Errorf("snake moved into too small of space, %s (can be cut off)", nextMove.Move)
// 	}
// 	if nextMove.Move == "left" {
// 		t.Errorf("snake moved into self, %s", nextMove.Move)
// 	}
// 	if nextMove.Move == "right" {
// 		t.Errorf("snake moved into wall, %s", nextMove.Move)
// 	}
// }
