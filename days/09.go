package days

import (
	"fmt"

	"parnic.com/aoc2019/utilities"
)

type Day09 struct {
	program utilities.IntcodeProgram
}

func (d *Day09) Parse() {
	d.program = utilities.LoadIntcodeProgram("09p")
}

func (d Day09) Num() int {
	return 9
}

func (d *Day09) Part1() string {
	var code int64
	d.program.RunIn(func(inputStep int) int64 {
		return 1
	}, func(val int64, state utilities.IntcodeProgramState) {
		code = val
	})

	return fmt.Sprintf("BOOST keycode: %s%d%s", utilities.TextBold, code, utilities.TextReset)
}

func (d *Day09) Part2() string {
	var coordinates int64
	d.program.Reset()
	d.program.RunIn(func(inputStep int) int64 {
		return 2
	}, func(val int64, state utilities.IntcodeProgramState) {
		coordinates = val
	})

	return fmt.Sprintf("Coordinates: %s%d%s", utilities.TextBold, coordinates, utilities.TextReset)
}
