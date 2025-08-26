// Package assert contains assert functions to make testing less verbose.
package assert

import (
	"fmt"
	"testing"
)

type Assert struct {
	t *testing.T
}

func New(t *testing.T) *Assert {
	return &Assert{t: t}
}

func (a *Assert) NoError(err error) {
	if err != nil {
		fmt.Printf("No error expected, got %v instead.", err)
		a.t.Fail()
	}
}

func (a *Assert) Error(err error) {
	if err == nil {
		fmt.Printf("Error expected, got %v instead.", err)
		a.t.Fail()
	}
}

func (a *Assert) True(stmt bool) {
	if !stmt {
		fmt.Println("Condition not true!")
		a.t.Fail()
	}
}

// EqualSlice will check each element in each slice and compare that they are the same.
func (a *Assert) EqualSlice(one []byte, two []byte) {
	if len(one) != len(two) {
		a.t.Fail()
	}
	for i := range one {
		if one[i] != two[i] {
			fmt.Printf("%v is not equal to %v\n", one, two)
			a.t.Fail()
		}
	}
}
