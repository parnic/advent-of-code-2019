package days

import (
	"fmt"
	"math"

	"parnic.com/aoc2019/utilities"
)

type Day01 struct {
	nums []int64
}

func (d *Day01) Parse() {
	d.nums = utilities.GetIntLines("01p")
}

func (d *Day01) calcFuel(mass int64) int64 {
	return int64(math.Floor(float64(mass)/3)) - 2
}

func (d *Day01) Part1() string {
	var totalFuel int64
	for _, mass := range d.nums {
		fuel := d.calcFuel(mass)
		totalFuel += fuel
	}

	return fmt.Sprintf("Fuel required: %s%d%s", utilities.TextBold, totalFuel, utilities.TextReset)
}

func (d *Day01) Part2() string {
	var totalFuel int64
	for _, mass := range d.nums {
		for mass > 0 {
			fuel := d.calcFuel(mass)
			if fuel > 0 {
				totalFuel += fuel
			}
			mass = fuel
		}
	}

	return fmt.Sprintf("Fuel required: %s%d%s", utilities.TextBold, totalFuel, utilities.TextReset)
}
