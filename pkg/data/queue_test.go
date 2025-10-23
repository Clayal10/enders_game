package data_test

import (
	"testing"

	"github.com/Clayal10/enders_game/lib/assert"
	"github.com/Clayal10/enders_game/lib/cross"
	"github.com/Clayal10/enders_game/lib/data"
)

func TestBasicQueuing(t *testing.T) {
	a := assert.New(t)

	type myObj struct {
		num int
	}

	q := data.NewQueue[*myObj](10)
	q.Enqueue(&myObj{num: 10})
	a.False(q.IsEmpty())
	obj, err := q.Dequeue()
	a.NoError(err)
	a.True(obj.num == 10)
	_, err = q.Dequeue()
	a.ErrorIs(err, cross.ErrQueueEmpty)
}
