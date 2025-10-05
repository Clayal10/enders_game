package assert_test

import (
	"testing"
	"time"

	"github.com/Clayal10/enders_game/lib/assert"
	"github.com/Clayal10/enders_game/lib/cross"
)

type mockAssert struct{}

func (ma *mockAssert) Helper()                   {}
func (ma *mockAssert) Errorf(_ string, _ ...any) {}
func (ma *mockAssert) Error(_ ...any)            {}

func TestAssertFunctions(t *testing.T) {
	a := assert.New(&mockAssert{})

	a.NoError(cross.ErrFrameTooSmall)
	a.NoError(nil)
	a.Error(cross.ErrFrameTooSmall)
	a.Error(nil)

	a.True(true)
	a.True(false)

	a.EqualSlice([]byte{0, 1, 2}, []byte{0, 1, 2})
	a.EqualSlice([]byte{0, 1, 3}, []byte{0, 1, 2})
	a.EqualSlice([]byte{0, 1}, []byte{0, 1, 2})

	a.NotNil(nil)
	a.NotNil(1)
	a.Nil(&mockAssert{})
	a.Nil(nil)

	count := 0
	a.Eventually(func() bool {
		count++
		return count > 10
	}, time.Millisecond*100, time.Millisecond)

	count = 0
	a.Eventually(func() bool {
		return count > 10
	}, time.Millisecond, time.Microsecond*500)
}
