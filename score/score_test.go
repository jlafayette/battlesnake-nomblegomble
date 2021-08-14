package score

import "testing"

func TestCombineWeights(t *testing.T) {
	tests := []struct {
		scores   []WeightedScore
		expected Scored
	}{
		{
			scores: []WeightedScore{
				{true, 1.0, Scored{"up": 1.0, "down": 1.0, "left": 1.0, "right": 1.0}},
			},
			expected: Scored{"up": 1.0, "down": 1.0, "left": 1.0, "right": 1.0},
		},
		{
			scores: []WeightedScore{
				{true, 1.0, Scored{"up": 1.0, "down": 1.0, "left": 1.0, "right": 1.0}},
				{true, 1.0, Scored{"up": 0.0, "down": 0.5, "left": 1.0, "right": 1.0}},
			},
			expected: Scored{"up": 0.0, "down": 0.5, "left": 1.0, "right": 1.0},
		},
		{
			scores: []WeightedScore{
				{true, 1.0, Scored{"up": 0.0, "down": 1.0, "left": 1.0, "right": 1.0}},
				{true, 1.0, Scored{"up": 0.5, "down": 0.01, "left": 0.5, "right": 0.5}},
				{false, 1.0, Scored{"up": 0.0, "down": 1.0, "left": 0.8, "right": 0.8}},
			},
			expected: Scored{"up": 0.0, "down": 0.02, "left": 0.9, "right": 0.9},
		},
		{
			scores: []WeightedScore{
				{true, 1, Scored{"down": 1, "left": 1, "right": 1, "up": 0}},
				{true, 1, Scored{"down": 0.25, "left": 0.25, "right": 0.1, "up": 0.25}},
				{false, 1, Scored{"down": 0, "left": 0, "right": 0, "up": 0}},
			},
			expected: Scored{"up": 0.0, "down": 0.25, "left": 0.25, "right": 0.1},
		},
	}

	for _, tc := range tests {

		actual := CombineMoves(tc.scores)

		for move, score := range actual {
			if actual[move] != tc.expected[move] {
				t.Errorf("%s: expected %.2f, got %.2f", move, tc.expected[move], score)
			}
		}
	}
}
