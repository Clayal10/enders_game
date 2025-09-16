package lurk

const (
	defenseDivisor = 200
	regenDivisor   = 500
)

// CalculateFight will take two characters, run a fight calculation on them and
// modify their stats accordingly.
//
// Calculations will be done with the following:
//
// - Attack decreases Health.
//
// - Defense will decrease the attack's effect by a % of 200.
//
// - Regen will add health as a percentage of how much damage was inflicted / 5
//
// These 3 stats can max out to 100 each. Health is limited by the stat limit in competition with gold.
func CalculateFight(c1, c2 *Character) {
	c2Damage := int16(float32(c1.Attack) - float32(c1.Attack)*(float32(c2.Defense)/defenseDivisor))
	c1Damage := int16(float32(c2.Attack) - float32(c2.Attack)*(float32(c1.Defense)/defenseDivisor))

	c2.Health -= c2Damage
	c1.Health -= c1Damage

	if c1.Health <= 0 {
		c1.Flags[Alive] = false
	} else {
		c1.Health += int16(float32(c1Damage) * float32(c1.Regen) / regenDivisor)
	}

	if c2.Health <= 0 {
		c2.Flags[Alive] = false
	} else {
		c2.Health += int16(float32(c2Damage) * float32(c2.Regen) / regenDivisor)
	}
}
