package q

import (
	"errors"
)

type Fifo struct {
	q        chan Item
	qlen     int
	capacity int
}

func NewFifo(capacity int) Fifo {
	return Fifo{
		q:        make(chan Item, capacity),
		qlen:     0,
		capacity: capacity,
	}
}

type Item struct {
	X       int
	Y       int
	Turn    int
	OpenIn  int // number of turns until this is open (for body collision squares)
	Collide bool
}

func (q *Fifo) Push(item Item) error {
	if q.qlen == q.capacity {
		return errors.New("queue is full")
	}
	q.q <- item
	// log.Printf("push %v", item)
	q.qlen += 1
	return nil
}

func (q *Fifo) Pop() (Item, bool) {
	if q.qlen <= 0 {
		// log.Printf("pop empty")
		return Item{}, true
	}
	q.qlen -= 1
	item := <-q.q
	// log.Printf("pop %v", item)
	return item, false
}
