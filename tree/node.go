package tree

import (
	"fmt"
	"strings"
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
	count := 1
	node := mn.FirstSibling()
	for {
		if node.nextSibling == nil {
			return count
		}
		node = node.nextSibling
		count += 1
	}
}

func (m1 *MoveNode) swap(m2 *MoveNode) {
	tmpM := m2.moves
	tmpS := m2.score
	tmpL := m2.scoredLevel
	tmpC := m2.child
	m2.moves = m1.moves
	m2.score = m1.score
	m2.scoredLevel = m1.scoredLevel
	m2.child = m1.child
	m1.moves = tmpM
	m1.score = tmpS
	m1.scoredLevel = tmpL
	m1.child = tmpC

}

// Sort
// Move for my snake with max min score grouped at the start
// Subsequent moves grouped with their minscore at start of group
func (mn *MoveNode) SortSiblings(myIndex int, bestMove Move) {
	for i := mn.FirstSibling(); i != nil; i = i.nextSibling {
		iSort := i.moves[myIndex].move
		if iSort == bestMove {
			iSort = 0
		}
		for j := i.nextSibling; j != nil; j = j.nextSibling {
			jSort := j.moves[myIndex].move
			if jSort == bestMove {
				jSort = 0
			}
			if iSort > jSort {
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

func (mn *MoveNode) NodeAfterPrune(myIndex int) *MoveNode {
	return nil
}

// Return score and bool (true if a best is found, otherwise false)
func (mn *MoveNode) BestSoFar(myIndex int) (float64, bool) {
	// need a complete set scored
	// take the min of the complete set

	// currentMove := mn.moves[myIndex].move

	// for n := mn.FirstSibling(); n != nil && n != mn; n = n.nextSibling {
	// }

	// minScore := 0
	// maxScore := 0
	// for _, m1 := range []Move{Left, Right, Up, Down} {
	// 	minScore := 99999.0
	// 	atLeastOne := false
	// 	for node != nil {
	// 		if node.moves[myIndex].move == m1 {
	// 			minScore = minf(minScore, node.score)
	// 		}
	// 		node = node.prevSibling
	// 	}
	// 	if atLeastOne {
	// 		if minScore > maxScore {
	// 			bestMove = m1
	// 		}
	// 		maxScore = maxf(maxScore, minScore)
	// 	}
	// }

	// take max of all the mins

	return 0, false
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
	sb.WriteString(fmt.Sprintf("%d %.1f", node.scoredLevel, node.score))
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
