package days

import (
	"fmt"

	"parnic.com/aoc2019/utilities"
)

type Day05 struct {
	program utilities.IntcodeProgram
}

func (d *Day05) Parse() {
	d.program = utilities.LoadIntcodeProgram("05p")
	d.test()
}

func (d Day05) Num() int {
	return 5
}

func (d Day05) test() {
	// Using position mode, consider whether the input is equal to 8; output 1 (if it is) or 0 (if it is not).
	program := utilities.ParseIntcodeProgram("3,9,8,9,10,9,4,9,99,-1,8")
	program.RunIn(func(int) int64 {
		return 0
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 0 {
			panic("")
		}
	})
	program.Reset()
	program.RunIn(func(int) int64 {
		return 8
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 1 {
			panic("")
		}
	})

	// Using position mode, consider whether the input is less than 8; output 1 (if it is) or 0 (if it is not).
	program = utilities.ParseIntcodeProgram("3,9,7,9,10,9,4,9,99,-1,8")
	program.RunIn(func(int) int64 {
		return 0
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 1 {
			panic("")
		}
	})
	program.Reset()
	program.RunIn(func(int) int64 {
		return 8
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 0 {
			panic("")
		}
	})

	// Using immediate mode, consider whether the input is equal to 8; output 1 (if it is) or 0 (if it is not).
	program = utilities.ParseIntcodeProgram("3,3,1108,-1,8,3,4,3,99")
	program.RunIn(func(int) int64 {
		return 0
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 0 {
			panic("")
		}
	})
	program.Reset()
	program.RunIn(func(int) int64 {
		return 8
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 1 {
			panic("")
		}
	})

	// Using immediate mode, consider whether the input is less than 8; output 1 (if it is) or 0 (if it is not).
	program = utilities.ParseIntcodeProgram("3,3,1107,-1,8,3,4,3,99")
	program.RunIn(func(int) int64 {
		return 0
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 1 {
			panic("")
		}
	})
	program.Reset()
	program.RunIn(func(int) int64 {
		return 8
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 0 {
			panic("")
		}
	})

	// jump tests that take an input, then output 0 if the input was zero or 1 if the input was non-zero
	// position mode
	program = utilities.ParseIntcodeProgram("3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9")
	program.RunIn(func(int) int64 {
		return 0
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 0 {
			panic("")
		}
	})
	program.Reset()
	program.RunIn(func(int) int64 {
		return 8
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 1 {
			panic("")
		}
	})

	// jump tests that take an input, then output 0 if the input was zero or 1 if the input was non-zero
	// immediate mode
	program = utilities.ParseIntcodeProgram("3,3,1105,-1,9,1101,0,0,12,4,12,99,1")
	program.RunIn(func(int) int64 {
		return 0
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 0 {
			panic("")
		}
	})
	program.Reset()
	program.RunIn(func(int) int64 {
		return 8
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 1 {
			panic("")
		}
	})

	// uses an input instruction to ask for a single number. The program will then output 999 if the input value is below 8, output 1000 if the input value is equal to 8, or output 1001 if the input value is greater than 8.
	program = utilities.ParseIntcodeProgram("3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99")
	program.RunIn(func(int) int64 {
		return 0
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 999 {
			panic("")
		}
	})
	program.Reset()
	program.RunIn(func(int) int64 {
		return 8
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 1000 {
			panic("")
		}
	})
	program.Reset()
	program.RunIn(func(int) int64 {
		return 9
	}, func(val int64, state utilities.IntcodeProgramState) {
		if val != 1001 {
			panic("")
		}
	})
}

func (d *Day05) Part1() string {
	diagCode := int64(-1)
	d.program.RunIn(func(int) int64 {
		return 1
	}, func(val int64, state utilities.IntcodeProgramState) {
		if state.IsHalting() {
			diagCode = val
		} else if val != 0 {
			panic("test failed")
		}
	})

	return fmt.Sprintf("Diagnostic code: %s%d%s", utilities.TextBold, diagCode, utilities.TextReset)
}

func (d *Day05) Part2() string {
	d.program.Reset()
	diagCode := int64(-1)
	d.program.RunIn(func(int) int64 {
		return 5
	}, func(val int64, state utilities.IntcodeProgramState) {
		if !state.IsHalting() {
			panic("unexpected output received")
		}
		diagCode = val
	})

	return fmt.Sprintf("Diagnostic code: %s%d%s", utilities.TextBold, diagCode, utilities.TextReset)
}
