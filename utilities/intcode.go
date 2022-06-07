package utilities

import (
	"strconv"
	"strings"
)

func ParseIntcodeProgram(programStr string) []int64 {
	nums := strings.Split(programStr, ",")
	program := make([]int64, len(nums))
	for idx, num := range nums {
		iNum, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			panic(err)
		}
		program[idx] = iNum
	}

	return program
}

func RunIntcodeProgram(program []int64) {
	for instructionPointer := 0; instructionPointer < len(program); {
		opcode := program[instructionPointer]
		switch opcode {
		case 1:
			param1 := program[instructionPointer+1]
			param2 := program[instructionPointer+2]
			param3 := program[instructionPointer+3]
			program[param3] = program[param1] + program[param2]

			instructionPointer += 4
			break
		case 2:
			param1 := program[instructionPointer+1]
			param2 := program[instructionPointer+2]
			param3 := program[instructionPointer+3]
			program[param3] = program[param1] * program[param2]

			instructionPointer += 4
			break
		case 99:
			instructionPointer = len(program)
			break
		}
	}
}
