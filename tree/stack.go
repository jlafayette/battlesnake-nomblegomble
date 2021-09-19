package tree

import "fmt"

type Stack struct {
	s [][]snakeMove
}

func NewStack(depth int) *Stack {
	x := make([][]snakeMove, 0, depth)
	return &Stack{s: x}
}

func (s *Stack) Push(moves []snakeMove) {
	s.s = append(s.s, moves)
}

func (s *Stack) Pop() []snakeMove {
	if len(s.s) == 0 {
		return nil
	}
	lastIndex := len(s.s) - 1
	lastItem := s.s[lastIndex]
	s.s = s.s[0:lastIndex]
	return lastItem
}

// -- for debugging snake

type snakeMoveArgs struct {
	move    Move
	food    bool
	die     bool
	undo    bool
	toUndo1 int
	toUndo2 int
}

func (s snakeMoveArgs) String() string {
	if s.undo {
		return fmt.Sprintf("%d->%d undo", s.toUndo1, s.toUndo2)
	}
	return fmt.Sprintf("%d->%d move: %v, food: %v, die %v", s.toUndo1, s.toUndo2, s.move, s.food, s.die)
}

type StackSnakeMoveArgs struct {
	s            []snakeMoveArgs
	pushPopCount int
}

func NewStackSnakeMoveArgs(depth int) *StackSnakeMoveArgs {
	x := make([]snakeMoveArgs, 0, depth)
	return &StackSnakeMoveArgs{s: x}
}

func (s *StackSnakeMoveArgs) Push(moves snakeMoveArgs) {
	s.pushPopCount += 1
	s.s = append(s.s, moves)
}

func (s *StackSnakeMoveArgs) Pop() (snakeMoveArgs, error) {
	s.pushPopCount += 1
	if len(s.s) == 0 {
		return snakeMoveArgs{}, fmt.Errorf("stack is empty")
	}
	lastIndex := len(s.s) - 1
	lastItem := s.s[lastIndex]
	s.s = s.s[0:lastIndex]
	return lastItem, nil
}
