package utilities

import (
	"strconv"
	"strings"
)

const (
	opAdd = 1
	opMul = 2
	opEnd = 99
)

type IntcodeProgram []int64

func ParseIntcodeProgram(programStr string) IntcodeProgram {
	nums := strings.Split(programStr, ",")
	program := make(IntcodeProgram, len(nums))
	for idx, num := range nums {
		iNum, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			panic(err)
		}
		program[idx] = iNum
	}

	return program
}

func RunIntcodeProgram(program IntcodeProgram) {
	for instructionPointer := 0; instructionPointer < len(program); {
		opcode := program[instructionPointer]
		switch opcode {
		case opAdd:
			param1 := program[instructionPointer+1]
			param2 := program[instructionPointer+2]
			param3 := program[instructionPointer+3]
			program[param3] = program[param1] + program[param2]

			instructionPointer += 4
			break
		case opMul:
			param1 := program[instructionPointer+1]
			param2 := program[instructionPointer+2]
			param3 := program[instructionPointer+3]
			program[param3] = program[param1] * program[param2]

			instructionPointer += 4
			break
		case opEnd:
			instructionPointer = len(program)
			break
		}
	}
}
