package tree

import (
	"fmt"
	"testing"
)

func TestRemap1(t *testing.T) {

	rawArea := -159
	myLen := 7
	// -1614.3 := remap(float64(-159), 0, float64(21), -100, 100)
	result := remap(float64(rawArea), 0, float64(myLen*3), -100, 100)
	fmt.Println(result)
}
