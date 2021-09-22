package tree

import (
	"fmt"
)

type eatEvent struct {
	turn       int
	coord      Coord
	prevHealth int
}

func removeEatEvent(s []eatEvent, i int) []eatEvent {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type Snake struct {
	Index int

	ba    []Coord // body "array" has extra capacity for rewinding moves
	start int     // start index of Body in body array
	end   int     // end index of Body in body array
	Body  []Coord // active body (window into part of the body array)

	eatEvents   []eatEvent
	turn        int
	ateLastTurn bool
	toUndo      int

	Length      int
	Health      int
	HealthStack *IntStack

	Dead   bool
	DiedOn int
}

func NewSnake(index, health int, coords []Coord, turn, depth int) *Snake {
	// TODO: Convert from double tail as last index into proper internal snake representation
	// 		 (if it matters)
	start := depth
	end := len(coords) + depth
	b := make([]Coord, 0, len(coords)+depth) // +1 extra?
	for i := 0; i < depth; i++ {
		b = append(b, Coord{0, 0})
	}
	b = append(b, coords...)
	// b = append(b, Coord{0, 0}) // extra?
	eatEvents := make([]eatEvent, 0)
	return &Snake{
		Index:       index,
		ba:          b,
		start:       start,
		end:         end,
		Body:        b[start:end],
		eatEvents:   eatEvents,
		turn:        turn,
		Length:      len(coords),
		Health:      health,
		HealthStack: NewIntStack(depth),
		Dead:        false,
		DiedOn:      -1,
	}
}

func (s *Snake) Move(m Move, food, die, hazard bool) error {
	// fmt.Printf("%d %d move: %v, food %v, die %v\n", s.Index, s.turn+1, m, food, die)

	// Eating normally mutates the previous tail into a duplicate of the new tail,
	// forming a "double tail"
	// But since this destroys some of the history, we intead record the new length
	// and add an entry into the ate field.
	s.turn += 1

	if s.Dead {
		// fmt.Printf("Snake %d is dead, ignoring move %v food: %v die %v\n", s.Index, m, food, die)
		return nil
	}

	s.toUndo += 1

	// don't extend body unless ate last turn and not this turn
	// ate last turn and ate this turn?
	growTail := false
	if s.ateLastTurn {
		growTail = true
	}

	if !growTail {
		s.end -= 1
	}
	x := s.ba[s.start]
	s.start -= 1

	s.ba[s.start] = x.Move(m)
	// This line here is always where the crashes happen
	var err error
	if s.end > len(s.ba) {
		err = fmt.Errorf("error making move %v for snake %d, %d is out of bounds (len %d)", m, s.Index, s.end, len(s.ba))
		s.end = len(s.ba)
	}
	s.Body = s.ba[s.start:s.end]

	if food {
		s.Length += 1
		s.eatEvents = append(s.eatEvents, eatEvent{turn: s.turn, coord: s.ba[s.start], prevHealth: s.Health})
		s.ateLastTurn = true
		s.HealthStack.Push(s.Health)
		s.Health = 100
	} else {
		s.ateLastTurn = false
		s.HealthStack.Push(s.Health)
		if hazard {
			s.Health -= 16
		} else {
			s.Health -= 1
		}
	}
	if s.Health <= 0 {
		die = true
	}
	if die {
		// fmt.Printf("snake %d died on turn %d\n", s.Index, s.turn)
		s.Dead = true
		s.DiedOn = s.turn
	}
	// fmt.Println(s.Body)
	return err
}

func (s *Snake) UndoMove() (*Coord, error) {
	if s.toUndo <= 0 {
		return nil, fmt.Errorf("error undoing move for snake %d, no more moves to undo", s.Index)
	}
	if s.Dead && s.DiedOn < s.turn {
		s.turn -= 1
		return nil, nil
	}

	s.toUndo -= 1

	var food *Coord
	growTail := true
	removeEventIndex := -1
	for i, event := range s.eatEvents {
		if event.turn == s.turn {
			growTail = false
			s.Length -= 1
			food = &event.coord
			removeEventIndex = i
			s.Health = event.prevHealth
			break
		}
	}
	if removeEventIndex >= 0 {
		s.eatEvents = removeEatEvent(s.eatEvents, removeEventIndex)
	}
	// should this be after or before we decrement the turn?
	if s.DiedOn == s.turn {
		// fmt.Printf("snake un-died on turn %d\n", s.turn)
		s.Dead = false
		s.DiedOn = -1
	}
	s.turn -= 1 // decrement turn after checking eat events
	s.start += 1
	if growTail {
		s.end += 1
	}
	var err error
	if s.end > len(s.ba) {
		err = fmt.Errorf("error undoing move for snake %d, %d is out of bounds (len %d)", s.Index, s.end, len(s.ba))
		s.end = len(s.ba)
	}
	s.Body = s.ba[s.start:s.end]
	s.Health = s.HealthStack.Pop(100)
	return food, err
}

func (s *Snake) Head() Coord {
	if len(s.Body) == 0 {
		panic("snake has body length of zero")
	}
	return s.Body[0]
}

func (s1 *Snake) Vs(s2 *Snake, m1, m2 Move) bool {
	// Check if m1 or m2 would kill one or both of the snakes
	head1 := s1.Head().Move(m1)
	head2 := s2.Head().Move(m2)
	die1 := false

	// H2H
	if head1.Equals(head2) {
		die1 = s1.Length <= s2.Length
		// if die1 {
		// 	fmt.Printf("snake %d will die in h2h with %d\n", s1.Index, s2.Index)
		// }
	}

	// Body collisions
	for i, b2 := range s2.Body {
		if !s2.ateLastTurn && i == len(s2.Body)-1 {
			continue
		}
		if head1.Equals(b2) {
			die1 = true
			// fmt.Printf("snake %d will die from body collision with %d on %v\n", s1.Index, s2.Index, b2)
			break
		}
	}

	return die1
}

func (s *Snake) VsSelf(m Move) bool {
	die := false
	newHead := s.Head().Move(m)
	for i, b := range s.Body {
		if !s.ateLastTurn && i == len(s.Body)-1 {
			continue
		}
		if newHead.Equals(b) {
			die = true
			break
		}
	}
	return die
}
