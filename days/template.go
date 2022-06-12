package days

import (
	"fmt"

	u "parnic.com/aoc2019/utilities"
)

type DayTemplate struct {
}

func (d *DayTemplate) Parse() {
	u.GetIntLines("Templatep")
}

func (d DayTemplate) Num() int {
	return -1
}

func (d *DayTemplate) Part1() string {
	return fmt.Sprintf("%s%d%s", u.TextBold, 0, u.TextReset)
}

func (d *DayTemplate) Part2() string {
	return fmt.Sprintf("%s%d%s", u.TextBold, 0, u.TextReset)
}
