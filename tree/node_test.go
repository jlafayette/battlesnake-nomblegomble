package tree

import (
	"testing"
)

func TestSingleNodeBestSoFar(t *testing.T) {

	root := &MoveNode{}
	node1 := &MoveNode{
		moves:       []snakeMove{{0, Left}, {1, Left}},
		scores:      []float64{5, 0},
		scoredLevel: 2,
		parent:      root,
	}
	root.child = node1

	// -- one scored node

	_, _, score, ok := node1.BestSoFar(0, 2)

	if score != 5 || !ok {
		t.Errorf("score should be found for a single node, expected 5.00,true, got %.2f,%v", score, ok)
	}
}

func TestNodeBestSoFar01(t *testing.T) {

	root := &MoveNode{}
	node1 := &MoveNode{
		moves:       []snakeMove{{0, Left}, {1, Left}},
		scores:      []float64{5, 0},
		scoredLevel: 2,
		parent:      root,
	}
	root.child = node1
	node2 := &MoveNode{
		moves:       []snakeMove{{0, Left}, {1, Right}},
		scores:      []float64{0, 0},
		scoredLevel: 0,
		parent:      root,
		prevSibling: node1,
	}
	node1.nextSibling = node2
	node3 := &MoveNode{
		moves:       []snakeMove{{0, Right}, {1, Left}},
		scores:      []float64{0, 0},
		scoredLevel: 0,
		parent:      root,
		prevSibling: node2,
	}
	node2.nextSibling = node3
	node4 := &MoveNode{
		moves:       []snakeMove{{0, Right}, {1, Right}},
		scores:      []float64{0, 0},
		scoredLevel: 0,
		parent:      root,
		prevSibling: node3,
	}
	node3.nextSibling = node4

	// -- one scored node is incomplete

	_, _, _, ok := node1.BestSoFar(0, 2)

	if ok {
		t.Errorf("no complete move set is scored, expected false, got %v", ok)
	}

	// -- two scored nodes (both L)

	node2.scores[0] = 4
	node2.scoredLevel = 2

	_, _, score, ok := node2.BestSoFar(0, 2)

	if score != 4 || !ok {
		t.Errorf("score should be min of the moves, expected 4.00,true, got %.2f,%v", score, ok)
	}

	// -- LL LR RL (RR)

	node3.scores[0] = 2
	node3.scoredLevel = 2

	_, _, score, ok = node3.BestSoFar(0, 2)

	if score != 4 || !ok {
		t.Errorf("score should be min of L since R incomplete, expected 4.00,true, got %.2f,%v", score, ok)
	}

	// -- LL LR RL RR

	node4.scores[0] = 2
	node4.scoredLevel = 2

	_, _, score, ok = node4.BestSoFar(0, 2)

	if score != 4 || !ok {
		t.Errorf("score should be min of L since R has lower min, expected 4.00,true, got %.2f,%v", score, ok)
	}
}

func TestNodePrune01(t *testing.T) {

	root := &MoveNode{}
	node1 := &MoveNode{
		moves:       []snakeMove{{0, Left}, {1, Left}},
		scores:      []float64{5, 0},
		scoredLevel: 2,
		parent:      root,
	}
	root.child = node1
	node2 := &MoveNode{
		moves:       []snakeMove{{0, Left}, {1, Right}},
		scores:      []float64{4, 0},
		scoredLevel: 2,
		parent:      root,
		prevSibling: node1,
	}
	node1.nextSibling = node2
	node3 := &MoveNode{
		moves:       []snakeMove{{0, Right}, {1, Left}},
		scores:      []float64{0, 0},
		scoredLevel: 1,
		parent:      root,
		prevSibling: node2,
	}
	node2.nextSibling = node3
	node4 := &MoveNode{
		moves:       []snakeMove{{0, Right}, {1, Right}},
		scores:      []float64{0, 0},
		scoredLevel: 1,
		parent:      root,
		prevSibling: node3,
	}
	node3.nextSibling = node4
	node5 := &MoveNode{
		moves:       []snakeMove{{0, Up}, {1, Left}},
		scores:      []float64{0, 0},
		scoredLevel: 1,
		parent:      root,
		prevSibling: node4,
	}
	node4.nextSibling = node5

	// RL if lower than bestSoFar (4 from LR), so the rest of R moves can be pruned
	node3.scores[0] = 2
	node3.scoredLevel = 2

	nextNode, _ := node3.NodeAfterPrune(0, 2)

	if nextNode != node5 {
		t.Errorf("other R nodes should be pruned, expected %v, got %v", node5, nextNode)
	}
	if !node4.pruned {
		t.Errorf("%v should be pruned, expected %v, got %v", node4, true, node4.pruned)
	}
}
