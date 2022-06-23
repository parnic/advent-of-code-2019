package days

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

type moonData struct {
	pos u.Vec3[int]
	vel u.Vec3[int]
}

func (m *moonData) applyGravity(other *moonData) {
	applyGravityAxis := func(pos1, pos2, vel1, vel2 *int) {
		if *pos1 < *pos2 {
			*vel1++
			*vel2--
		} else if *pos1 > *pos2 {
			*vel1--
			*vel2++
		}
	}

	applyGravityAxis(&m.pos.X, &other.pos.X, &m.vel.X, &other.vel.X)
	applyGravityAxis(&m.pos.Y, &other.pos.Y, &m.vel.Y, &other.vel.Y)
	applyGravityAxis(&m.pos.Z, &other.pos.Z, &m.vel.Z, &other.vel.Z)
}

func (m *moonData) applyVelocity() {
	m.pos.Add(m.vel)
}

func (m moonData) getPotentialEnergy() int {
	return int(math.Abs(float64(m.pos.X))) +
		int(math.Abs(float64(m.pos.Y))) +
		int(math.Abs(float64(m.pos.Z)))
}

func (m moonData) getKineticEnergy() int {
	return int(math.Abs(float64(m.vel.X))) +
		int(math.Abs(float64(m.vel.Y))) +
		int(math.Abs(float64(m.vel.Z)))
}

func (m moonData) getTotalEnergy() int {
	return m.getPotentialEnergy() * m.getKineticEnergy()
}

type Day12 struct {
	moons []*moonData
}

func (d *Day12) Parse() {
	lines := u.GetStringLines("12p")
	d.moons = make([]*moonData, len(lines))
	for i, line := range lines {
		trimmed := line[1 : len(line)-1]
		vals := strings.Split(trimmed, ", ")
		x, _ := strconv.Atoi(vals[0][2:])
		y, _ := strconv.Atoi(vals[1][2:])
		z, _ := strconv.Atoi(vals[2][2:])
		d.moons[i] = &moonData{
			pos: u.Vec3[int]{X: x, Y: y, Z: z},
		}
	}
}

func (d Day12) Num() int {
	return 12
}

func (d Day12) copyMoons() []*moonData {
	moons := make([]*moonData, len(d.moons))
	for i, moon := range d.moons {
		moonCopy := *moon
		moons[i] = &moonCopy
	}

	return moons
}

func getAllEnergy(moons ...*moonData) int {
	energy := 0
	for _, moon := range moons {
		energy += moon.getTotalEnergy()
	}
	return energy
}

func (d *Day12) Part1() string {
	moons := d.copyMoons()

	numSteps := 1000

	for i := 0; i < numSteps; i++ {
		for i, moon1 := range moons {
			for _, moon2 := range moons[i+1:] {
				moon1.applyGravity(moon2)
			}

			moon1.applyVelocity()
		}
	}

	return fmt.Sprintf("Total energy after %d steps: %s%d%s", numSteps, u.TextBold, getAllEnergy(moons...), u.TextReset)
}

func (d *Day12) Part2() string {
	moons := d.copyMoons()

	orig := make([]u.Vec3[int], len(moons))
	for i, moon := range moons {
		orig[i] = moon.pos
	}
	period := u.Vec3[int]{}

	for loops := 0; period.X == 0 || period.Y == 0 || period.Z == 0; loops++ {
		for i, moon1 := range moons {
			for _, moon2 := range moons[i+1:] {
				moon1.applyGravity(moon2)
			}
			moon1.applyVelocity()
		}

		foundX := true
		foundY := true
		foundZ := true
		for i, moon := range moons {
			if moon.pos.X != orig[i].X || moon.vel.X != 0 {
				foundX = false
			}
			if moon.pos.Y != orig[i].Y || moon.vel.Y != 0 {
				foundY = false
			}
			if moon.pos.Z != orig[i].Z || moon.vel.Z != 0 {
				foundZ = false
			}
		}
		if foundX && period.X == 0 {
			period.X = loops + 1
		}
		if foundY && period.Y == 0 {
			period.Y = loops + 1
		}
		if foundZ && period.Z == 0 {
			period.Z = loops + 1
		}
	}

	stepsRequired := u.LCM(period.X, period.Y, period.Z)
	return fmt.Sprintf("Iterations to reach a previous state: %s%d%s", u.TextBold, stepsRequired, u.TextReset)
}
