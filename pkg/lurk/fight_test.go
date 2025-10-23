package lurk_test

import (
	"testing"

	"github.com/Clayal10/enders_game/pkg/assert"
	"github.com/Clayal10/enders_game/pkg/lurk"
)

func TestFightCalculation(t *testing.T) {
	a := assert.New(t)
	t.Run("TestOneBeatingTwo", func(_ *testing.T) {
		// This should take less than 10 attacks with their health so high.
		c1 := &lurk.Character{
			Flags: map[string]bool{
				lurk.Alive: true,
			},
			Attack:  100,
			Defense: 100,
			Regen:   100,
			Health:  500,
		}
		c2 := &lurk.Character{
			Flags: map[string]bool{
				lurk.Alive: true,
			},
			Attack:  100,
			Defense: 50,
			Regen:   50,
			Health:  500,
		}
		count := 0
		for c2.Flags[lurk.Alive] {
			lurk.CalculateFight(c1, c2)
			count++
		}
		a.True(c1.Flags[lurk.Alive] == true)
		a.True(count < 10)
	})
	t.Run("TestTwoBeatingOne", func(_ *testing.T) {
		c1 := &lurk.Character{
			Flags: map[string]bool{
				lurk.Alive: true,
			},
			Attack:  100,
			Defense: 10,
			Regen:   50,
			Health:  100,
		}
		c2 := &lurk.Character{
			Flags: map[string]bool{
				lurk.Alive: true,
			},
			Attack:  100,
			Defense: 50,
			Regen:   50,
			Health:  500,
		}
		count := 0
		for c1.Flags[lurk.Alive] {
			lurk.CalculateFight(c1, c2)
			count++
		}
		a.True(c2.Flags[lurk.Alive] == true)
		a.True(count < 3)
	})
}
