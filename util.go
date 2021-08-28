package main

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func remap(old, oldMin, oldMax, newMin, newMax float64) float64 {
	oldRange := oldMax - oldMin
	if oldRange == 0 {
		return newMin
	} else {
		newRange := newMax - newMin
		return (((old - oldMin) * newRange) / oldRange) + newMin
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
