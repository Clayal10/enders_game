// Package assert contains assert functions to make testing less verbose.
package assert

import (
	"errors"
	"time"
)

type TestAsserter interface {
	Helper()
	Errorf(format string, args ...any)
	Error(args ...any)
}

type Assert struct {
	t TestAsserter
}

func New(t TestAsserter) *Assert {
	return &Assert{t: t}
}

func (a *Assert) NoError(err error) {
	a.t.Helper()
	if err != nil {
		a.t.Errorf("No error expected, got %v instead.\n", err)
	}
}

func (a *Assert) Error(err error) {
	a.t.Helper()
	if err == nil {
		a.t.Errorf("Error expected, got %v instead.\n", err)
	}
}

func (a *Assert) ErrorIs(expected, actual error) {
	a.t.Helper()
	if !errors.Is(expected, actual) {
		a.t.Errorf("Error %v expected, got error %v instead", expected, actual)
	}
}

func (a *Assert) True(stmt bool) {
	a.t.Helper()
	if !stmt {
		a.t.Error("Condition not True!")
	}
}

func (a *Assert) False(stmt bool) {
	a.t.Helper()
	if stmt {
		a.t.Error("Condition not True!")
	}
}

// EqualSlice will check each element in each slice and compare that they are the same.
func (a *Assert) EqualSlice(one []byte, two []byte) {
	a.t.Helper()

	if len(one) != len(two) {
		a.t.Error("Slices have unequal length!")
	}
	for i := range one {
		if one[i] != two[i] {
			a.t.Errorf("%v is not equal to %v\n", one, two)
		}
	}
}

func (a *Assert) NotNil(obj any) {
	a.t.Helper()

	if obj == nil {
		a.t.Error("Object is nil!")
	}
}

func (a *Assert) Nil(obj any) {
	a.t.Helper()

	if obj != nil {
		a.t.Error("Object is not nil!")
	}
}

func (a *Assert) Eventually(f func() bool, duration time.Duration, tick time.Duration) {
	a.t.Helper()
	start := time.Now()
	ticker := time.NewTicker(tick)

	passed := false
	for range ticker.C {
		if time.Since(start) > duration {
			ticker.Stop()
			break
		}

		if f() {
			passed = true
			break
		}
	}

	if !passed {
		a.t.Error("Did not meet condition!")
	}
}
