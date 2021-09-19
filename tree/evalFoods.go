package tree

import (
	"math"
)

func foodCalc(turn Turn) float64 {
	// Messed around with a curve that kind of fits the shape we want
	e := 2.71828
	a := 0.1
	b := 6.0
	c := -0.16
	return maxf(0, (a + (b * math.Pow(e, (c*float64(turn))))))
}

type foodTracker struct {
	score   float64
	visited map[Coord]bool
}

func newFoodTracker() *foodTracker {
	vmap := make(map[Coord]bool)
	return &foodTracker{
		score:   0.0,
		visited: vmap,
	}
}

func (f *foodTracker) add(c Coord, t Turn) {
	visited := f.visited[c]
	if visited {
		return
	}
	f.visited[c] = true

	// Messed around with a curve that kind of fits the shape we want
	e := 2.71828
	a := 0.1
	b := 6.0
	d := -0.16
	score := maxf(0, (a + (b * math.Pow(e, (d*float64(t))))))
	f.score += score

	// fmt.Printf("got food on {%d, %d} for a score of %.1f\n", c.X, c.Y, score)
}

func (f *foodTracker) reset() {
	for c := range f.visited {
		f.visited[c] = false
	}
	f.score = 0.0
}
