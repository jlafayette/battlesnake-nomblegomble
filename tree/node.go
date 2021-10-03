package tree

import (
	"fmt"
	"strings"
)

const (
	HIGHEST = 9999999.0
	LOWEST  = -9999999.0
)

func (mn *MoveNode) EvalSiblings(snakeCount, level int) (*MoveNode, bool) {
	// this is a variation of minmax
	if mn == nil {
		panic("can't eval siblings since node is nil")
	}

	// find the lowest scored move where the snake has an alternative and drop
	// it repeat until all snakes have only one option

	// find a bad move for one of the snakes (the worst)
	// trim it
	didAPrune := true
	for didAPrune {
		didAPrune = mn.pruneLowestSibling(snakeCount, level)
	}

	nextNode, done := mn.resolve(level)
	if !done {
		nextNode = mn.nextUnscoredSibling(level)
		p, t := nextNode.place()
		fmt.Printf("Not done yet, nextNode (%d/%d): %v\n", p, t, nextNode)
	} else {
		fmt.Println("DONE!!")
	}

	return nextNode, done
	// return MoveNodeEval{NextNode: nextNode, FinalNode: finalNode}

	// move counts (how many options)
	// if all snakes have only one option then we are done

	// do we have enough info yet to do any pruning?
}

func (mn *MoveNode) pruneLowestSibling(snakeCount, level int) bool {
	// if we have enough information, prune a move for one of the snakes
	// (all the move combos that contain that move can be pruned)

	// to prune, we need all the node scored for that move
	// LLL LLR LRL LRR (0 can do L)

	// And... we need one other move that is better than the lowest of the fully scored move
	// LLL LLR LRL LRR RLL (0 can do L because if that one R is higher)
	firstNode := mn.FirstSibling()

	pruneIndex := -1
	pruneMove := NoMove
	found := false

	// for id in snakeids
	// TODO: randomize snake index order and do other snakes before mySnake
	for index := 0; index < snakeCount; index++ {
		if found {
			break
		}
		lowest := HIGHEST
		lowestMove := NoMove
		//   for move in Up,Lf,Rt,Dn
		for _, m1 := range []Move{Up, Down, Right, Left} {
			if found {
				break
			}
			lowest2 := HIGHEST
			complete := true
			atLeastOne := false
			//     for node in nodes
			for node := firstNode; node != nil; node = node.nextSibling {
				if found {
					break
				}
				if node.pruned {
					continue
				}
				//       check for complete
				//       optionally record lowest
				m2 := node.moves[index].move
				if m1 != m2 {
					continue
				}
				atLeastOne = true
				if node.scoredLevel != level {
					complete = false
					break
				}
				lowest2 = minf(lowest2, node.scores[index])
			}
			if !complete || !atLeastOne {
				continue
			}
			// ok, now we have a complete set and a lowest2 score for that set
			// if lowestMove is already found, we can compare and do a prune,
			fmt.Printf("found a complete set for %d:%v\n", index, m1)
			if lowestMove != NoMove {
				if lowest2 < lowest {
					pruneIndex = index
					pruneMove = m1
				} else {
					pruneIndex = index
					pruneMove = lowestMove
				}
				found = true
			}
			// if lowestMove is not already found, then we store this one as
			// the lowest (it's always the lowest since it's the only move)
			lowestMove = m1
			lowest = minf(lowest2, lowest)
		}
	}

	if !found {
		return false
	}
	count := 0
	for node := firstNode; node != nil; node = node.nextSibling {
		if node.moves[pruneIndex].move == pruneMove {
			count += 1
		}
	}
	fmt.Printf("found a lowest set to prune %d:%v (count: %d)\n", pruneIndex, pruneMove, count)
	actuallyPrunedSomthing := false
	for node := firstNode; node != nil; node = node.nextSibling {

		if node.pruned || node.scoredLevel != level {
			continue
		}
		if node.moves[pruneIndex].move == pruneMove {
			p1, t := node.place()
			fmt.Printf("%d pruned a move: %d/%d %v\n", level, p1, t, node)
			node.pruned = true
			actuallyPrunedSomthing = true
		}
	}
	fmt.Printf("...prunedSomething: %v\n", actuallyPrunedSomthing)
	return actuallyPrunedSomthing
}

func (mn *MoveNode) resolve(level int) (*MoveNode, bool) {
	// here we return the final node if we can narrow it down that far

	// for this to work we need only one node left after skipping pruned nodes
	// (and the one node left can't be unscored)
	// return nil if that isn't the case

	firstNode := mn.FirstSibling()

	var candidate *MoveNode
	for node := firstNode; node != nil; node = node.nextSibling {
		if node.pruned {
			continue
		}
		if node.scoredLevel != level {
			return node, false
		}
		if candidate == nil {
			candidate = node
		} else {
			// two candidates... so we can't resolve
			fmt.Printf("%d two candidates... not done yet\n", level)
			// p, t := candidate.place()
			// fmt.Printf("  candidate1 (%d/%d): %v\n", p, t, candidate)
			// p, t = node.place()
			// fmt.Printf("  candidate2 (%d/%d): %v\n", p, t, node)
			return nil, false
		}
	}
	if candidate != nil {
		return candidate, true
	}
	panic("uh oh")
}

func (mn *MoveNode) nextUnscoredSibling(level int) *MoveNode {
	if mn.nextSibling == nil && mn.scoredLevel != level {
		panic("resolve didn't do it's job")
	}
	for node := mn; node != nil; node = node.nextSibling {
		if node.pruned {
			continue
		}
		if node.scoredLevel != level {
			return node
		}
	}
	fmt.Println("before no unscored siblings...")
	p, t := mn.place()
	fmt.Printf("%d/%d, %v\n", p, t, mn)
	for node := mn.FirstSibling(); node != nil; node = node.nextSibling {

		if node.pruned {
			continue
		}
		if node.scoredLevel != level {
			panic("bad")
		}
		p, t = node.place()
		fmt.Printf("scored %d/%d: %v\n", p, t, node)
		// p, t = node.place()
		// fmt.Printf("scored %d/%d: %v\n", p, t, node)
	}
	panic("no unscored siblings")
}

// A node in the tree
type MoveNode struct {
	// These are moves to achieve this position from the parent.
	moves []snakeMove

	// scores for iterative deepening
	scores      []float64 // the score for the current deepening level
	scoredLevel int       // the deepening level this has been scored at last
	pruned      bool      // if the node has been pruned (no score needed)

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

func NewMoveNode(snakeCount int, moves []snakeMove) *MoveNode {
	scores := make([]float64, 0, snakeCount)
	for i := 0; i < snakeCount; i++ {
		scores = append(scores, 0.0)
	}
	return &MoveNode{
		moves:  moves,
		scores: scores,
	}
}

func (mn *MoveNode) FirstSibling() *MoveNode {
	node := mn
	for {
		if node.prevSibling == nil {
			return node
		}
		node = node.prevSibling
	}
}

func (mn *MoveNode) LastSibling() *MoveNode {
	node := mn
	for {
		if node.nextSibling == nil {
			return node
		}
		node = node.nextSibling
	}
}

func (mn *MoveNode) SiblingCount() int {
	count := 0
	for node := mn.FirstSibling(); node != nil; node = node.nextSibling {
		count += 1
	}
	return count
}

func (mn *MoveNode) place() (int, int) {
	before := 0
	for node := mn; node != nil; node = node.prevSibling {
		before += 1
	}
	return before, mn.SiblingCount()
}

func (mn *MoveNode) hierarchy() string {
	// .2/23 ..1/123 ...4/4
	strs := make([]string, 0)
	for node := mn; node != nil; node = node.parent {
		p, t := node.place()
		strs = append(strs, fmt.Sprintf("%d/%d", p, t))
	}
	var sb strings.Builder
	for i := len(strs) - 2; i >= 0; i -= 1 {
		for j := 0; j < len(strs)-i-1; j++ {
			sb.WriteByte('.')
		}
		sb.WriteString(strs[i])
		sb.WriteByte(' ')
	}
	return sb.String()
}

func (m1 *MoveNode) swap(m2 *MoveNode) {
	// p1, _ := m1.place()
	// p2, _ := m2.place()
	// fmt.Printf("swap %d<>%d\n", p1, p2)
	tmpM := m2.moves
	tmpS := m2.scores
	tmpL := m2.scoredLevel
	tmpC := m2.child
	tmpP := m2.pruned
	m2.moves = m1.moves
	m2.scores = m1.scores
	m2.scoredLevel = m1.scoredLevel
	m2.child = m1.child
	if m2.child != nil {
		m2.child.parent = m2
	}
	m2.pruned = m1.pruned

	m1.moves = tmpM
	m1.scores = tmpS
	m1.scoredLevel = tmpL
	m1.child = tmpC
	m1.pruned = tmpP
	if m1.child != nil {
		m1.child.parent = m1
	}

	// m1.assert()
	// m2.assert()
}

func (mn *MoveNode) assert() {
	if mn.child == nil {
		return
	}
	if mn.child.parent != mn {
		panic("got all screwed up here")
	}
}

// Sort
// Move for my snake with max min score grouped at the start
// Subsequent moves grouped with their minscore at start of group
func (mn *MoveNode) SortSiblings(myIndex int, bestMove Move) {

	// fmt.Printf("sorting with (%d) %s\n", myIndex, bestMove.ShortString())

	// This seems to break everything?

	// Sort by scores (within move groups)
	// for i := mn.FirstSibling(); i != nil; i = i.nextSibling {
	// 	for j := i.nextSibling; j != nil; j = j.nextSibling {
	// 		if i.moves[myIndex].move != j.moves[myIndex].move {
	// 			continue
	// 		}
	// 		if i.score > j.score {
	// 			i.swap(j)
	// 		}
	// 	}
	// }

	for i := mn.FirstSibling(); i != nil; i = i.nextSibling {
		iSort := i.moves[myIndex].move
		if iSort == bestMove {
			iSort = NoMove
		}
		for j := i.nextSibling; j != nil; j = j.nextSibling {
			jSort := j.moves[myIndex].move
			if jSort == bestMove {
				jSort = NoMove
			}
			if iSort > jSort {
				// fmt.Printf("%d(%s) > %d(%s)  swap\n", iSort, iSort.ShortString(), jSort, jSort.ShortString())
				i.swap(j)
			}
		}
	}
	// for i := mn.FirstSibling(); i != nil; i = i.nextSibling {
	// 	iSort := i.moves[myIndex].move
	// 	if iSort == bestMove {
	// 		iSort = NoMove
	// 	}
	// 	for j := i.nextSibling; j != nil; j = j.nextSibling {
	// 		jSort := j.moves[myIndex].move
	// 		if jSort == bestMove {
	// 			jSort = NoMove
	// 		}
	// 		if iSort > jSort {
	// 			// fmt.Printf("%d(%s) > %d(%s)  swap\n", iSort, iSort.ShortString(), jSort, jSort.ShortString())
	// 			i.swap(j)
	// 		}
	// 	}
	// }

}

func (mn *MoveNode) NodeAfterPrune(myIndex, level int) (*MoveNode, int) {
	// If pruning can occur, then mark all the nodes as pruned and return the
	// next node that cannot be pruned. It's ok to return nil if all the nodes
	// can be pruned or no

	// fmt.Printf("--- prune opportunity for %v (%d) level: (%d)\n", mn, myIndex, level)

	if mn == nil || mn.nextSibling == nil {
		return mn, 0
	}

	_, _, best, ok := mn.BestSoFar(myIndex, level)
	if !ok {
		return mn, 0
	}
	myMove := mn.moves[myIndex].move
	lowest := HIGHEST
	found := false
	firstNode := mn.FirstSibling()
	for node := firstNode; node != nil; node = node.nextSibling {
		if node.moves[myIndex].move != myMove {
			continue
		}
		if node.scoredLevel != level || node.pruned {
			continue
		}
		lowest = minf(lowest, node.scores[myIndex])
		found = true
	}
	if !found {
		return mn, 0
	}
	// we should not prune if lowest > best
	if lowest > best {
		return mn, 0
	}

	pruned := 0
	for node := mn.nextSibling; node != nil; node = node.nextSibling {
		if node.moves[myIndex].move != myMove {
			// if pruned > 0 {
			// 	fmt.Printf("%d->successfully pruned! jumping  %v -> %v (%d)\n", pruned, mn, node, myIndex)
			// }
			return node, pruned
		}
		pruned += 1
		node.pruned = true
	}
	// if pruned > 0 {
	// 	fmt.Printf("%d->successfully pruned! jumping  %v -> %v (%d)\n", pruned, mn, nil, myIndex)
	// }
	return nil, pruned
}

// Return score and bool (true if a best is found, otherwise false)
func (mn *MoveNode) BestSoFar(myIndex, level int) (Move, Move, float64, bool) {
	// need a complete set scored
	// take the min of the complete set

	found := false
	maxScore := LOWEST
	bestMove := NoMove
	luckyMove := NoMove
	for _, m1 := range []Move{Left, Right, Up, Down} {
		allScored := true
		atLeastOne := false
		minScore := HIGHEST
		luckyScore := LOWEST

		for node := mn.FirstSibling(); node != nil; node = node.nextSibling {
			m2 := node.moves[myIndex].move
			if m2 != m1 || node.pruned {
				// fmt.Printf("%v != %v || %v\n", m2, m1, node.pruned)
				continue
			}
			if node.scoredLevel != level {
				// fmt.Printf("%d != %d\n", node.scoredLevel, level)
				allScored = false
				break
			}
			minScore = minf(minScore, node.scores[myIndex])
			// fmt.Printf("%s setting minScore to %.2f\n", node.String(), minScore)
			atLeastOne = true
			luckyScore = maxf(luckyScore, node.scores[myIndex])
		}
		if allScored && atLeastOne {
			if minScore > maxScore {
				bestMove = m1
			}
			if luckyScore > maxScore {
				luckyMove = m1
			}
			maxScore = maxf(maxScore, minScore)
			// fmt.Printf("%s setting maxScore to %.2f\n", m1.ShortString(), maxScore)
			found = true
		}
	}
	return bestMove, luckyMove, maxScore, found
}

func (mn *MoveNode) ResetPrunedSiblings() {
	for node := mn.FirstSibling(); node != nil; node = node.nextSibling {
		node.pruned = false
	}
}

// Convert to a short display string, for example: "0:L 1:L 2:R 3:X 1 212.0"
func (node *MoveNode) String() string {
	var sb strings.Builder
	for _, m := range node.moves {
		sb.WriteString(m.ShortString())
		sb.WriteByte(' ')
		sb.WriteString(fmt.Sprintf("%.1f", node.scores[m.snakeIndex]))
		sb.WriteByte(' ')
	}
	prunedStr := "-"
	if node.pruned {
		prunedStr = "X"
	}
	sb.WriteString(fmt.Sprintf("%d%s", node.scoredLevel, prunedStr))
	return sb.String()
}

func (mn *MoveNode) PrintSiblings() {
	node := mn.FirstSibling()
	for {
		fmt.Println(node)
		if node.nextSibling == nil {
			break
		}
		node = node.nextSibling
	}
	fmt.Printf("connected: %d\n", node.SiblingCount())
}

func (mn *MoveNode) MyMovesString(myIndex int) string {
	strs := make([]string, 0)
	for node := mn; node.parent != nil; node = node.parent {
		strs = append(strs, node.moves[myIndex].move.ShortString())
	}
	var sb strings.Builder
	for i := len(strs) - 1; i >= 0; i -= 1 {
		sb.WriteString(strs[i])
	}
	return sb.String()
}
