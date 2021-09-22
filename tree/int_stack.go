package tree

type IntStack struct {
	s []int
}

func NewIntStack(depth int) *IntStack {
	x := make([]int, 0, depth)
	return &IntStack{s: x}
}

func (s *IntStack) Push(i int) {
	s.s = append(s.s, i)
}

func (s *IntStack) Pop(emptyValue int) int {
	if len(s.s) == 0 {
		return emptyValue
	}
	lastIndex := len(s.s) - 1
	lastItem := s.s[lastIndex]
	s.s = s.s[0:lastIndex]
	return lastItem
}
