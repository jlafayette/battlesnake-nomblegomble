package q

import "testing"

func TestFifo(t *testing.T) {

	q := NewFifo(10)

	item1 := Item{X: 0, Y: 0}
	item2 := Item{X: 1, Y: 2}
	item3 := Item{X: 3, Y: 4}

	q.Push(item1)
	q.Push(item2)
	q.Push(item3)

	r1, e1 := q.Pop()
	r2, e2 := q.Pop()
	r3, e3 := q.Pop()
	_, e4 := q.Pop()

	equal := func(i1, i2 Item) bool {
		return i1.X == i2.X && i1.Y == i2.Y
	}

	if e1 || e2 || e3 || !e4 {
		t.Errorf("expected Fifo queue to be empty after popping all items")
	}

	if !equal(item1, r1) || !equal(item2, r2) || !equal(item3, r3) {
		t.Errorf("expected first in to be first out")
	}
}

func TestFifoFull(t *testing.T) {

	q := NewFifo(2)

	item1 := Item{X: 0, Y: 0}
	item2 := Item{X: 1, Y: 2}
	item3 := Item{X: 3, Y: 4}

	err1 := q.Push(item1)
	err2 := q.Push(item2)
	err3 := q.Push(item3)

	if err1 != nil || err2 != nil {
		t.Errorf("expected Fifo queue to have capacity for 2 items")
	}

	if err3 == nil {
		t.Errorf("expected error if queue is full")
	}
}
