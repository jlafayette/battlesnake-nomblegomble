package tree

import (
	"fmt"

	"github.com/jlafayette/battlesnake-go/wire"
)

type snakeMove struct {
	snakeIndex int
	move       Move
}

func (s snakeMove) String() string {
	return fmt.Sprintf("%d: %v", s.snakeIndex, s.move)
}

// Root of the tree
type MoveNode struct {
	// These are moves to achieve this position from the parent.
	moves []snakeMove

	score  float64
	scored bool

	// The parent position, if nil, then we are the root of the tree
	parent *MoveNode

	// The first child position that can be reached, to iterate over all the
	// child positions, start with the child and go to the sibling
	child *MoveNode

	// The sibling positions (same parent)
	// each parent only has one child (the first reachable position). The other
	// children can be reached by using these two links. If nextSibling is nil,
	// then we are at the end of the children. If prevSibling is nil, then we
	// the first child position.

	// Links example, child2 has a link to the parent, but not the other way
	// around.
	// p
	// ^      ^
	// v
	// c1 <> c2
	nextSibling *MoveNode
	prevSibling *MoveNode
}

// State keeps track of the current position
// Supports pushing and popping positions, so backtracking through past moves
// is possible
type State struct {
	Width  int
	Height int

	MyIndex int

	// The turn at the root position
	InitialTurn int

	// How many moves deep to go
	maxDepth     int
	currentDepth int

	// The tree, for moving forward and backwards and evaluating each position
	Root *MoveNode
	// The current node in the tree
	node *MoveNode

	// The eval board evaluates a position and returns a score
	evalBoard *Board

	// Snakes
	Snakes []*Snake

	// Food (eaten, needs to be un-eaten when reversing)
	Food []Coord

	// Hazards (static for now)
	Hazards []Coord
}

func NewState(wireState *wire.GameState, depth int) *State {
	snakeCount := len(wireState.Board.Snakes)
	snakes := make([]*Snake, 0, snakeCount)
	for i, srcSnake := range wireState.Board.Snakes {
		coords := make([]Coord, 0, len(srcSnake.Body))
		for _, srcCoord := range srcSnake.Body {
			coords = append(coords, Coord{X: srcCoord.X, Y: srcCoord.Y})
		}
		snakes = append(snakes, NewSnake(i, int(srcSnake.Health), coords, wireState.Turn, depth))
	}
	food := make([]Coord, 0, len(wireState.Board.Food))
	for _, wireFood := range wireState.Board.Food {
		food = append(food, Coord{X: wireFood.X, Y: wireFood.Y})
	}
	myIndex := -1
	for i, snake := range wireState.Board.Snakes {
		if snake.ID == wireState.You.ID {
			myIndex = i
			break
		}
	}
	hazards := make([]Coord, 0)
	board := NewBoard(wireState.Board.Width, wireState.Board.Height, snakes, food, hazards)

	root := &MoveNode{}
	return &State{
		Width:        wireState.Board.Width,
		Height:       wireState.Board.Height,
		MyIndex:      myIndex,
		InitialTurn:  wireState.Turn,
		maxDepth:     depth,
		currentDepth: 0,
		Root:         root,
		node:         root,
		evalBoard:    board,
		Snakes:       snakes,
		Food:         food,
		Hazards:      hazards,
	}
}

type DebugInfo struct {
	CurrentDepth    int
	Move            int
	Moves           int
	MoveDescription string
	History         []string
}

func (s *State) DebugInfo() DebugInfo {
	node := 0
	nodes := 0
	description := ""
	history := make([]string, 0)
	if s.node != nil {
		prev := s.node.prevSibling
		numPrev := 0
		for prev != nil {
			nodes += 1
			numPrev += 1
			prev = prev.prevSibling
		}
		next := s.node.nextSibling
		numNext := 0
		for next != nil {
			nodes += 1
			numNext += 1
			next = next.nextSibling
		}
		node = numPrev + 1
		nodes = numPrev + 1 + numNext

		for _, m := range s.node.moves {
			description = description + m.String() + " "
		}
		history = append(history, "> "+description)
		parent := s.node.parent
		for parent != nil {
			s := "  "
			for _, m := range parent.moves {
				s = s + m.String() + " "
			}
			history = append(history, s)
			parent = parent.parent
		}
	}
	return DebugInfo{
		CurrentDepth:    s.currentDepth,
		Move:            node,
		Moves:           nodes,
		MoveDescription: description,
		History:         history,
	}
}

type CoordInfo struct {
	Snake *Snake
}

func (s *State) CoordInfo(c Coord) CoordInfo {
	var selSnake *Snake
	breakAll := false
	for _, snake := range s.Snakes {
		if breakAll {
			break
		}
		for _, bc := range snake.Body {
			if c.Equals(bc) {
				selSnake = snake
				breakAll = true
				break
			}
		}
	}
	return CoordInfo{Snake: selSnake}
}

func (s *State) moves(index int) []snakeMove {
	// return Dead move if applicable (also for solo, anything beyond 0)
	if index >= len(s.Snakes) {
		m := []snakeMove{{snakeIndex: index, move: Dead}}
		// fmt.Printf("moves for snake %d: %v\n", index, m)
		return m
	}
	if s.Snakes[index].Dead {
		m := []snakeMove{{snakeIndex: index, move: Dead}}
		// fmt.Printf("moves for snake %d: %v\n", index, m)
		return m
	}

	// Future enhancements
	// Don't bother with moves that kill the snake instantly (if there are any
	// other options)
	// if all moves for a snake lead to death, we don't need to explore all of
	// those, just one move leading to death (and H2H if both snakes die)
	snake := s.Snakes[index]
	head := snake.Body[0]

	m := make([]snakeMove, 0)
	leftHead := head.Move(Left)
	rightHead := head.Move(Right)
	upHead := head.Move(Up)
	downHead := head.Move(Down)
	leftCollide := false
	rightCollide := false
	upCollide := false
	downCollide := false
	for _, snake := range s.Snakes {
		for ci, c := range snake.Body {
			// Skipping tail
			if ci == len(snake.Body)-1 && !snake.ateLastTurn {
				continue
			}
			if leftHead.Equals(c) {
				leftCollide = true
			} else if rightHead.Equals(c) {
				rightCollide = true
			} else if upHead.Equals(c) {
				upCollide = true
			} else if downHead.Equals(c) {
				downCollide = true
			}
		}
	}
	if leftHead.InBounds(s.Width, s.Height) && !leftCollide {
		m = append(m, snakeMove{snakeIndex: index, move: Left})
	}
	if rightHead.InBounds(s.Width, s.Height) && !rightCollide {
		m = append(m, snakeMove{snakeIndex: index, move: Right})
	}
	if upHead.InBounds(s.Width, s.Height) && !upCollide {
		m = append(m, snakeMove{snakeIndex: index, move: Up})
	}
	if downHead.InBounds(s.Width, s.Height) && !downCollide {
		m = append(m, snakeMove{snakeIndex: index, move: Down})
	}
	if len(m) == 0 {
		// fmt.Printf("All moves end in death for %d\n", index)
		m = append(m, snakeMove{snakeIndex: index, move: Left})
	}
	// fmt.Printf("moves for snake %d: %v\n", index, m)
	return m
}

func (s *State) createPossibleMoves() {
	// This directly creates child of the current node, then siblings until
	// all possible moves are found
	// fmt.Println("creating possible moves")

	var prev *MoveNode
	m0moves := s.moves(0)
	m1moves := s.moves(1)
	m2moves := s.moves(2)
	m3moves := s.moves(3)
	for _, m0 := range m0moves {
		for _, m1 := range m1moves {
			for _, m2 := range m2moves {
				for _, m3 := range m3moves {
					node := &MoveNode{moves: []snakeMove{m0, m1, m2, m3}, parent: s.node}
					if prev == nil {
						s.node.child = node
					} else {
						prev.nextSibling = node
						node.prevSibling = prev
					}
					prev = node
				}
			}
		}
	}
}

func (s *State) DownLevel() {
	// just for testing, move down the tree
	// fmt.Printf("max depth %d, current depth: %d\n", s.Depth, s.currentDepth)
	// fmt.Printf("DownLevel: currentDepth %d\n", s.currentDepth)
	if s.currentDepth >= s.maxDepth {
		return
	}

	// First find a move
	// if no children, populate next level of children nodes
	if s.node.child == nil {
		s.createPossibleMoves()
	}
	s.node = s.node.child
	s.currentDepth += 1

	// Then, apply the move
	s.ApplyMove()
}

type newHeadInfo struct {
	snake *Snake
	coord Coord
	food  bool
	move  Move
	die   bool
}

func (s *State) ApplyMove() {
	if s.node == nil {
		fmt.Println("ERROR: s.node == nil")
		return
	}
	newHeads := make([]*newHeadInfo, 0, 4)
	for _, move := range s.node.moves {
		// shouldn't happen...
		if move.snakeIndex >= len(s.Snakes) {
			continue
		}
		snake := s.Snakes[move.snakeIndex]
		if snake.Dead {
			// We need to call move even though the snake is dead
			// so it's internal turn count stay accurate.
			err := snake.Move(Dead, false, false)
			if err != nil {
				fmt.Println(err.Error())
			}
			continue
		}
		head := snake.Head()
		newHead := head.Move(move.move)
		food := false
		removeIndex := -1
		for i, f := range s.Food {
			if f.Equals(newHead) {
				food = true
				removeIndex = i
				break
			}
		}
		if food {
			s.Food = remove(s.Food, removeIndex)
		}
		newHeads = append(newHeads,
			&newHeadInfo{
				snake: snake,
				coord: newHead,
				food:  food,
				move:  move.move,
			},
		)
	}
	for _, head := range newHeads {
		if head.snake.Dead {
			continue
		}
		die := false

		// check wall, other snakes
		if !head.coord.InBounds(s.Width, s.Height) {
			die = true
		}
		for _, otherHead := range newHeads {
			if die {
				break
			}
			if head.snake.Index == otherHead.snake.Index {
				continue
			}
			die = head.snake.Vs(otherHead.snake, head.move, otherHead.move)
		}
		head.die = die
	}
	for _, head := range newHeads {
		err := head.snake.Move(head.move, head.food, head.die)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (s *State) UpLevel() {
	// for testing, move up the tree
	// fmt.Printf("UpLevel: currentDepth %d\n", s.currentDepth)
	if s.currentDepth <= 0 {
		return
	}
	s.node = s.node.parent
	s.currentDepth -= 1

	// unapply the move
	for _, snake := range s.Snakes {
		food, err := snake.UndoMove()
		if err != nil {
			fmt.Println(err.Error())
			n := s.node
			for n != nil {
				fmt.Printf("  %v, %d\n", n.moves, s.currentDepth)
				n = n.parent
			}
		}
		if food != nil {
			s.Food = append(s.Food, *food)
		}
	}
}

func (s *State) NextSibling() {
	// check for next - if nil, return
	nextNode := s.node.nextSibling
	if nextNode == nil {
		return
	}
	// go up to parent (uplevel)
	s.UpLevel()
	// go to new child
	s.node = nextNode
	s.currentDepth += 1
	// Then, apply the move
	s.ApplyMove()
}

func (s *State) PrevSibling() {
	// check for prev - if nil, return
	prevNode := s.node.prevSibling
	if prevNode == nil {
		return
	}
	// go up to parent (uplevel)
	s.UpLevel()
	// go to new child
	s.node = prevNode
	s.currentDepth += 1
	// Then, apply the move
	s.ApplyMove()
}

// Score sets node score from the current position loaded into eval.Board
func (s *State) Score() {
	score := s.evalBoard.Eval(SnakeIndex(s.MyIndex))
	s.node.score = score
	s.node.scored = true
}

func (s *State) FindBestMove(verbose bool) Move {
	// Not going to do iterative deepening yet, just a set depth

	loop_count := 0

	for {
		loop_count += 1
		if s.node == nil {
			fmt.Printf("ERROR: s.node == nil after %d loops\n", loop_count)
			return Up
		}
		// If the root node is scored, then we are done
		if s.node.parent == nil && s.node.scored {
			panic("should have returned already")
		}

		// If not scored, either go down or score the node (if at max depth)
		if !s.node.scored { // not scored
			atMaxDepth := s.currentDepth >= s.maxDepth
			if !atMaxDepth {
				// if node is not scored and depth is not max
				s.DownLevel()
			} else { // at max depth
				// if node is not scored and depth is max
				// score the current node
				if !s.node.scored && (s.currentDepth >= s.maxDepth) {
					s.evalBoard.Load(s.Snakes, s.Food, s.Hazards)
					score := s.evalBoard.Eval(SnakeIndex(s.MyIndex))
					s.node.score = score
					// fmt.Printf("score: %.2f\n", score)
					s.node.scored = true
				}
			}
		} else { // If scored
			nextSibling := s.node.nextSibling
			if nextSibling != nil { // node is scored and next sibling
				s.NextSibling()
			} else { // node is scored and no next sibling
				// eval set (backtrack over the siblings and get the group score)
				// then go up to parent and score it
				// scan over neighbors, taking the min of each group (grouped by myMove)
				// take the max of the group
				lastChild := s.node
				maxScore := -99999.0
				bestMove := Dead // this is kind of the None state for moves
				for _, m1 := range []Move{Left, Right, Up, Down} {
					minScore := 99999.0
					atLeastOne := false
					node := lastChild
					for node != nil {
						if node.moves[s.MyIndex].move == m1 {
							atLeastOne = true
							minScore = minf(minScore, node.score)
						}
						node = node.prevSibling
					}
					if atLeastOne {
						if minScore > maxScore {
							bestMove = m1
						}
						maxScore = maxf(maxScore, minScore)
					}
				}

				// go up to parent, apply score from children
				s.UpLevel()
				s.node.score = maxScore
				s.node.scored = true
				if s.node.parent == nil {
					if verbose {
						fmt.Printf("Ran %d loops\n", loop_count)
					}
					return bestMove
				}
			}
		}
	}
}
