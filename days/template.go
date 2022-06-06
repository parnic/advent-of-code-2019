package days

import (
	"parnic.com/aoc2019/utilities"
)

type DayTemplate struct {
}

func (d *DayTemplate) Parse() {
	utilities.GetIntLines("Templatep")
}

func (d DayTemplate) Num() int {
	return -1
}

func (d *DayTemplate) Part1() string {
	return ""
}

func (d *DayTemplate) Part2() string {
	return ""
}
