package tree

import (
	"fmt"
	"strings"
	"time"

	"github.com/jlafayette/battlesnake-go/wire"
)

type snakeMove struct {
	snakeIndex int
	move       Move
}

func (s snakeMove) String() string {
	return fmt.Sprintf("%d: %v", s.snakeIndex, s.move)
}

func (s snakeMove) ShortString() string {
	return fmt.Sprintf("%d:%s", s.snakeIndex, s.move.ShortString())
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
	maxDepth       int
	currentDepth   int
	deepeningLevel int
	timeout        int64 // timeout in ms

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
	hazards := make([]Coord, 0, len(wireState.Board.Hazards))
	for _, wireHazard := range wireState.Board.Hazards {
		hazards = append(hazards, Coord{wireHazard.X, wireHazard.Y})
	}
	board := NewBoard(wireState.Board.Width, wireState.Board.Height, snakes, food, hazards)

	// timeout should have 100ms buffer, but always be at least 50ms
	// the min of 50 is mostly for test cases where this is not specified
	timeout := max(50, int(wireState.Game.Timeout)-50)

	root := &MoveNode{}
	return &State{
		Width:        wireState.Board.Width,
		Height:       wireState.Board.Height,
		MyIndex:      myIndex,
		InitialTurn:  wireState.Turn,
		maxDepth:     depth,
		currentDepth: 0,
		timeout:      int64(timeout),
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
	snake  *Snake
	coord  Coord
	food   bool
	move   Move
	die    bool
	hazard bool
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
			err := snake.Move(Dead, false, false, false)
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
		hazard := false
		for _, h := range s.Hazards {
			if h.Equals(newHead) {
				hazard = true
				break
			}
		}
		newHeads = append(newHeads,
			&newHeadInfo{
				snake:  snake,
				coord:  newHead,
				food:   food,
				move:   move.move,
				hazard: hazard,
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
				die = head.snake.VsSelf(head.move)
				continue
			}
			die = head.snake.Vs(otherHead.snake, head.move, otherHead.move)
		}
		head.die = die
	}
	for _, head := range newHeads {
		err := head.snake.Move(head.move, head.food, head.die, head.hazard)
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
	s.node.scoredLevel = s.deepeningLevel
}

func (s *State) printNodeStack() {
	n := s.node
	for n != nil {
		fmt.Println(n.moves)
		n = n.parent
	}
}

func (s *State) FindBestMove(verbose bool) (Move, int) {
	// Not going to do iterative deepening yet, just a set depth

	start := time.Now()
	s.deepeningLevel = 1

	// The best move so far found. This is updated each time we complete a new
	// deepening level. If a level times out, then we return the best move
	// found so far (from previous level).
	bestMove := Up

	// The deepening loop
	for {
		mv, timeout, failed := s.findBestMove(start, verbose)
		if timeout || failed {
			return bestMove, s.deepeningLevel
		}
		bestMove = mv
		fmt.Printf("got best move %v at level %d\n", mv, s.deepeningLevel)
		if s.deepeningLevel >= s.maxDepth {
			fmt.Printf("got best move %v at max depth of %d\n", bestMove, s.maxDepth)
			return bestMove, s.deepeningLevel
		}
		s.deepeningLevel += 1
	}
}

func (s *State) findBestMove(start time.Time, verbose bool) (Move, bool, bool) {

	eval_count := 0

	// what happened this loop
	// !scored !atMaxDepth ---> go down
	// !scored atMaxDepth ---> score the node
	// scored nextSibling ---> go to next sibling
	// scored !nextSibling ---> calculate level score and push it up a level
	// 		(optional) prune node/nodes
	// end
	var lastScored bool
	var lastDepth int
	var lastDepth2 int
	var lastNode *MoveNode
	// var lastNode2 *MoveNode
	var lastPruneCount int
	var calc bool
	var sb strings.Builder
	loop_count := 0
	times := 0

	for {
		loop_count += 1

		// if loop_count%5001 == 0 {
		// 	fmt.Println(loop_count)
		// 	s.printNodeStack()
		// 	times += 1
		// }

		// check for timeout here
		elapsed := time.Since(start)
		if elapsed.Milliseconds() > s.timeout {
			fmt.Printf("timing out on level %d after %v (%d)\n", s.deepeningLevel, elapsed, eval_count)
			return Up, true, false
		}

		if s.node == nil {
			fmt.Printf("ERROR: s.node == nil after %d evals\n", eval_count)
			return Up, false, true
		}

		atMaxDepth := s.currentDepth >= s.deepeningLevel
		scored := s.node.scoredLevel == s.deepeningLevel

		//--
		calc = false
		lastPruneCount = 0
		lastScored = false
		lastDepth = s.currentDepth
		lastNode = s.node
		//--

		// If the root node is scored, then we are done
		if s.node.parent == nil && scored {
			panic("should have returned already")
		}

		// If not scored, either go down or score the node (if at max depth)
		if !scored { // not scored
			if !atMaxDepth {
				// if node is not scored and depth is not max

				s.DownLevel()
			} else { // at max depth (leaf nodes)
				// if node is not scored and depth is max
				// score the current node
				// fmt.Printf("%v\n", s.node.moves)
				// s.printNodeStack()

				// before scoring... can this be pruned?
				//--
				lastScored = true
				//--

				s.evalBoard.Load(s.Snakes, s.Food, s.Hazards)
				score := s.evalBoard.Eval(SnakeIndex(s.MyIndex))
				s.node.score = score
				// fmt.Printf("score: %.2f\n", score)
				s.node.scoredLevel = s.deepeningLevel
				eval_count += 1
			}
		} else { // If scored
			nextSibling := s.node.nextSibling
			if nextSibling != nil { // node is scored and next sibling
				s.NextSibling()
			} else { // node is scored and no next sibling
				// eval set (backtrack over the siblings and get the group score)
				// then go up to parent and score it
				// scan over neighbors, taking the min of each group (grouped by myMove)
				// take the max of the groups

				// debug print
				// fmt.Println("--- choosing a best move from siblings")
				// bkNode := s.node
				// s.node = s.node.parent
				// s.printNodeStack()
				// s.node = bkNode

				//--
				calc = true
				//--

				lastChild := s.node
				maxScore := LOWEST

				bestMove := Dead // this is kind of the None state for moves
				// Sometimes there is no good move if the opponent plays perfectly (and can
				// predict the future), so in this case, we need a backup 'lucky' move to
				// play instead. This comes up when there are losing H2H moves, but the
				// other snake might go a different direction, so it's better to try it
				// and hope for the best.
				luckyMove := Up // Use this if bestMove is Dead
				luckyMoveFound := false

				for _, m1 := range []Move{Left, Right, Up, Down} {
					minScore := HIGHEST
					luckyScore := LOWEST
					atLeastOne := false
					node := lastChild
					for node != nil {
						if node.moves[s.MyIndex].move == m1 && !node.pruned {
							atLeastOne = true
							minScore = minf(minScore, node.score)
							luckyScore = maxf(luckyScore, node.score)
						}
						node = node.prevSibling
					}
					if atLeastOne {
						if minScore > maxScore {
							bestMove = m1
							// fmt.Printf("  found new best move: %v\n", m1)
						}
						if luckyScore > maxScore {
							luckyMove = m1
							luckyMoveFound = true
							// fmt.Printf("  found new lucky move %v\n", m1)
						}
						maxScore = maxf(maxScore, minScore)
					} else {
						// fmt.Printf("  atLeastOne: %v for %v\n", atLeastOne, m1)
					}
				}

				doPruning := true

				// before going up to the parent, sort the moves from best to
				// worst so that next deepening level can do better some
				// pruning
				if doPruning {
					// if s.currentDepth == 1 {
					// 	fmt.Printf("before sorting (%d)\n", s.MyIndex)
					// 	s.node.PrintSiblings()
					// 	fmt.Printf("sorting with (%d) %s\n", s.MyIndex, bestMove.ShortString())
					// }
					s.node.SortSiblings(s.MyIndex, bestMove)
					// if s.currentDepth == 1 {
					// 	fmt.Printf("after sorting (%d)\n", s.MyIndex)
					// 	s.node.PrintSiblings()
					// }
				}

				// go up to parent, apply score from children
				if doPruning {
					s.node.ResetPrunedSiblings()
				}
				s.UpLevel()
				s.node.score = maxScore
				s.node.scoredLevel = s.deepeningLevel
				// fmt.Printf("  Pushed a new score (%.1f) up to the parent\n", maxScore)

				// make this a loop?
				for {
					// can we attempt to prune
					if !doPruning || s.node.parent == nil || s.node.nextSibling == nil {
						break
					}
					newNode, count := s.node.NodeAfterPrune(s.MyIndex, s.deepeningLevel)
					if newNode == s.node {
						break
					}
					//--
					lastPruneCount = count
					//--
					// fmt.Printf("successfully pruned! jumping  %v  ->  %v\n", s.node, newNode)
					if newNode != nil {
						// s.NextSibling() but jumping forward
						s.UpLevel()
						s.node = newNode
						s.currentDepth += 1
						s.ApplyMove()
						break
					}

					// all the rest were pruned, so we can go up another level
					var score float64
					bestMove, luckyMove, score, _ = s.node.BestSoFar(s.MyIndex, s.deepeningLevel)
					luckyMoveFound = luckyMove != NoMove
					s.node.ResetPrunedSiblings()
					s.UpLevel()
					s.node.score = score
					s.node.scoredLevel = s.deepeningLevel
				}

				if s.node.parent == nil {
					if verbose {
						fmt.Printf("Evaluated %d positions\n", eval_count)
					}
					if bestMove == Dead || bestMove == NoMove {
						if luckyMoveFound {
							if verbose {
								fmt.Printf("No good move, let's hope '%v' works\n", luckyMove)
							}
							return luckyMove, false, false
						} else {
							if verbose {
								fmt.Println("No good or lucky move found")
							}
							return Dead, false, true
						}
					}
					return bestMove, false, false
				}
			}
		}
		//--
		lastDepth2 = s.currentDepth
		// lastNode2 = s.node
		//--

		if eval_count >= 99999 {
			// if eval_count >= 1854 {
			sb.WriteString(fmt.Sprintf("%d->%d %d ", lastDepth, lastDepth2, s.deepeningLevel))
			sb.WriteString(lastNode.hierarchy())
			if lastDepth2 == lastDepth+1 && lastNode.child != nil {
				sb.WriteString(" child->")
				sb.WriteString(lastNode.child.hierarchy())
				sb.WriteString(" child.parent->")
				sb.WriteString(lastNode.child.parent.hierarchy())

				// if lastNode != lastNode.child.parent {
				// 	panic("this is weird")
				// }
			}
			if lastScored {
				sb.WriteString(" !EVAL! ")
			}
			// sb.WriteString(lastNode.String())
			// if lastNode != lastNode2 {
			// 	sb.WriteString(" ---> ")
			// 	sb.WriteString(lastNode2.String())
			// }
			if lastPruneCount > 0 {
				sb.WriteString(fmt.Sprintf(" !PRUNED %d! ", lastPruneCount))
			}
			if calc {
				sb.WriteString("  !CALC LVL!")
			}
			fmt.Println(sb.String())
			sb.Reset()

			// fmt.Printf("loops: %d, evals: %d\n", loop_count, eval_count)
			// atMaxDepth := s.currentDepth >= s.deepeningLevel
			// scored := s.node.scoredLevel == s.deepeningLevel
			// fmt.Printf("  dpth: %d s.node: %v  scored: %v, atMaxDepth: %v\n", s.currentDepth, s.node, scored, atMaxDepth)
			// s.printNodeStack()
			// fmt.Printf("siblings: %d\n", s.node.SiblingCount())

			times += 1
			if times > 1000 {
				return Up, true, false
			}
		}
	}
}
