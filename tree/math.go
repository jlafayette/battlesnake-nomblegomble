package tree

func remap(old, oldMin, oldMax, newMin, newMax float64) float64 {
	oldRange := oldMax - oldMin
	if oldRange == 0 {
		return newMin
	} else {
		newRange := newMax - newMin
		return (((old - oldMin) * newRange) / oldRange) + newMin
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minf(a, b float64) float64 {
	if a <= b {
		return a
	}
	return b
}

func clamp(v, min_, max_ int) int {
	return min(max(v, min_), max_)
}
