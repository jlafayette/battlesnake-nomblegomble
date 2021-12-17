package tree

type ChoiceOutcome uint8

const (
	ChoiceSafe ChoiceOutcome = iota
	ChoiceBadH2H
	ChoiceDead
)

type choiceTracker struct {
	m map[Move]ChoiceOutcome
}

func newChoiceTracker() *choiceTracker {
	m := make(map[Move]ChoiceOutcome, 4)
	return &choiceTracker{
		m: m,
	}
}

func (c *choiceTracker) reset() {
	for k := range c.m {
		delete(c.m, k)
	}
}

func (c *choiceTracker) add(move Move, outcome ChoiceOutcome) {
	c.m[move] = outcome
}

func (c *choiceTracker) getSafe() int {
	safe := 0
	for _, v := range c.m {
		if v == ChoiceSafe {
			safe += 1
		}
	}
	return safe
}

func (c *choiceTracker) getBadH2h() int {
	bad := 0
	for _, v := range c.m {
		if v == ChoiceBadH2H {
			bad += 1
		}
	}
	return bad
}
