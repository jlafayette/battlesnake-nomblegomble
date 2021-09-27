package tree

import (
	"fmt"
	"strings"
)

const (
	HIGHEST = 9999999.0
	LOWEST  = -9999999.0
)

// A node in the tree
type MoveNode struct {
	// These are moves to achieve this position from the parent.
	moves []snakeMove

	// scores for iterative deepening
	score       float64 // the score for the current deepening level
	scoredLevel int     // the deepening level this has been scored at last
	pruned      bool    // if the node has been pruned (no score needed)

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

func (m1 *MoveNode) swap(m2 *MoveNode) {
	tmpM := m2.moves
	tmpS := m2.score
	tmpL := m2.scoredLevel
	tmpC := m2.child
	tmpP := m2.pruned
	m2.moves = m1.moves
	m2.score = m1.score
	m2.scoredLevel = m1.scoredLevel
	m2.child = m1.child
	m2.pruned = m1.pruned
	m1.moves = tmpM
	m1.score = tmpS
	m1.scoredLevel = tmpL
	m1.child = tmpC
	m1.pruned = tmpP
}

// Sort
// Move for my snake with max min score grouped at the start
// Subsequent moves grouped with their minscore at start of group
func (mn *MoveNode) SortSiblings(myIndex int, bestMove Move) {

	// sorting is bugged, needs to be done twice
	// fmt.Printf("sorting with (%d) %s\n", myIndex, bestMove.ShortString())

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

	// Sort by scores (within move groups)
	for i := mn.FirstSibling(); i != nil; i = i.nextSibling {
		for j := i.nextSibling; j != nil; j = j.nextSibling {
			if i.moves[myIndex].move != j.moves[myIndex].move {
				continue
			}
			if i.score > j.score {
				i.swap(j)
			}
		}
	}
}

func (mn *MoveNode) NodeAfterPrune(myIndex, level int) *MoveNode {
	// If pruning can occur, then mark all the nodes as pruned and return the
	// next node that cannot be pruned. It's ok to return nil if all the nodes
	// can be pruned or no

	// fmt.Printf("--- prune opportunity for %v (%d) level: (%d)\n", mn, myIndex, level)

	if mn == nil || mn.nextSibling == nil {
		return mn
	}

	_, _, best, ok := mn.BestSoFar(myIndex, level)
	if !ok {
		return mn
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
		lowest = minf(lowest, node.score)
		found = true
	}
	if !found {
		return mn
	}
	// this comparison is a mystery, seems like it should be opposite?
	// but lowest > best make all the test fail
	// we should not prune if lowest > best (this makes sense)
	if lowest > best {
		return mn
	}

	// pruned := 0
	for node := mn.nextSibling; node != nil; node = node.nextSibling {
		if node.moves[myIndex].move != myMove {
			// if pruned > 0 {
			// 	fmt.Printf("%d->successfully pruned! jumping  %v -> %v (%d)\n", pruned, mn, node, myIndex)
			// }
			return node
		}
		// pruned += 1
		node.pruned = true
	}
	// if pruned > 0 {
	// 	fmt.Printf("%d->successfully pruned! jumping  %v -> %v (%d)\n", pruned, mn, nil, myIndex)
	// }
	return nil
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
			minScore = minf(minScore, node.score)
			// fmt.Printf("%s setting minScore to %.2f\n", node.String(), minScore)
			atLeastOne = true
			luckyScore = maxf(luckyScore, node.score)
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
	}
	prunedStr := "-"
	if node.pruned {
		prunedStr = "X"
	}
	sb.WriteString(fmt.Sprintf("%d%s %.1f", node.scoredLevel, prunedStr, node.score))
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
