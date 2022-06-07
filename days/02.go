package days

import (
	"fmt"

	"parnic.com/aoc2019/utilities"
)

type Day02 struct {
	program []int64
}

func (d *Day02) Parse() {
	d.program = utilities.ParseIntcodeProgram(utilities.GetStringContents("02p"))
}

func (d Day02) Num() int {
	return 2
}

func (d *Day02) getProgramWithParams(param1, param2 int64) []int64 {
	program := make([]int64, len(d.program))
	copy(program, d.program)
	program[1] = param1
	program[2] = param2
	return program
}

func (d *Day02) Part1() string {
	program := d.getProgramWithParams(12, 2)
	utilities.RunIntcodeProgram(program)

	return fmt.Sprintf("Position 0 = %s%d%s", utilities.TextBold, program[0], utilities.TextReset)
}

func (d *Day02) Part2() string {
	sentinel := int64(19690720)

	var noun int64
	var verb int64
	found := false
	for noun = 0; noun <= 99; noun++ {
		for verb = 0; verb <= 99; verb++ {
			program := d.getProgramWithParams(noun, verb)
			utilities.RunIntcodeProgram(program)

			if program[0] == sentinel {
				found = true
				break
			}
		}

		if found {
			break
		}
	}

	if !found {
		panic("!found")
	}

	return fmt.Sprintf("%d created by noun=%d, verb=%d. 100 * noun + verb = %s%d%s",
		sentinel,
		noun,
		verb,
		utilities.TextBold,
		100*noun+verb,
		utilities.TextReset,
	)
}
