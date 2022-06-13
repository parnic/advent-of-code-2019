package days

import (
	"fmt"

	"parnic.com/aoc2019/utilities"
)

type Day02 struct {
	program utilities.IntcodeProgram
}

func (d *Day02) Parse() {
	d.program = utilities.LoadIntcodeProgram("02p")
}

func (d Day02) Num() int {
	return 2
}

func (d *Day02) setParams(param1, param2 int64) {
	d.program.Reset()
	d.program.SetMemory(1, param1)
	d.program.SetMemory(2, param2)
}

func (d *Day02) Part1() string {
	d.setParams(12, 2)
	d.program.Run()

	if d.program.GetMemory(0) != 4138658 {
		panic("")
	}
	return fmt.Sprintf("Position 0 = %s%d%s", utilities.TextBold, d.program.GetMemory(0), utilities.TextReset)
}

func (d *Day02) Part2() string {
	sentinel := int64(19690720)

	var noun int64
	var verb int64
	found := false
	for noun = 0; noun <= 99; noun++ {
		for verb = 0; verb <= 99; verb++ {
			d.setParams(noun, verb)
			d.program.Run()

			if d.program.GetMemory(0) == sentinel {
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
	if noun != 72 || verb != 64 {
		panic("")
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
