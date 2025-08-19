// Package assert contains assert functions to make testing less verbose.
package assert

import "testing"

type Assert struct {
	t *testing.T
}

func New(t *testing.T) *Assert {
	return &Assert{t: t}
}

func (a *Assert) NoError(err error) {
	if err != nil {
		a.t.Fail()
	}
}

func (a *Assert) Error(err error) {
	if err == nil {
		a.t.Fail()
	}
}

func (a *Assert) True(stmt bool) {
	if !stmt {
		a.t.Fail()
	}
}
