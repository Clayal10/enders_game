package data_test

import (
	"testing"

	"github.com/Clayal10/enders_game/pkg/assert"
	"github.com/Clayal10/enders_game/pkg/data"
)

func TestKeyValueList(t *testing.T) {
	a := assert.New(t)

	type object struct {
		thing string
	}
	m := data.NewMap[string, *object]()

	keys := []string{"one", "two", "three", "four"}

	m.Insert(keys[0], &object{thing: "ONE"})
	m.Insert(keys[1], &object{thing: "TWO"})
	m.Insert(keys[2], &object{thing: "THREE"})
	m.Insert(keys[3], &object{thing: "FOUR"})

	item, ok := m.Get(keys[1])
	a.True(ok)
	a.True(item.thing == "TWO")
}
