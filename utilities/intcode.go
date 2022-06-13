package utilities

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	opAdd         = 1
	opMultiply    = 2
	opInput       = 3
	opOutput      = 4
	opJumpIfTrue  = 5
	opJumpIfFalse = 6
	opLessThan    = 7
	opEquals      = 8
	opHalt        = 99

	modePosition  = 0
	modeImmediate = 1
)

type IntcodeProgram struct {
	memory  []int64
	program []int64
}

type IntcodeProgramState struct {
	program            *IntcodeProgram
	CurrentInstruction int
	NextInstruction    int
}

func (s IntcodeProgramState) IsHalting() bool {
	return s.program.GetMemory(s.NextInstruction) == opHalt
}

type ProvideInputFunc func(inputStep int) int64
type ReceiveOutputFunc func(val int64, state IntcodeProgramState)

func ParseIntcodeProgram(programStr string) IntcodeProgram {
	nums := strings.Split(programStr, ",")
	program := IntcodeProgram{
		program: make([]int64, len(nums)),
	}
	for idx, num := range nums {
		iNum, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			panic(err)
		}
		program.program[idx] = iNum
	}

	return program
}

func (p *IntcodeProgram) makeState(instructionPointer int) IntcodeProgramState {
	return IntcodeProgramState{
		program:            p,
		CurrentInstruction: instructionPointer,
		NextInstruction:    instructionPointer + 1,
	}
}

func (p *IntcodeProgram) Copy() IntcodeProgram {
	ret := IntcodeProgram{
		program: make([]int64, len(p.program)),
	}
	copy(ret.program, p.program)
	return ret
}

func (p *IntcodeProgram) init() {
	if p.memory == nil {
		p.memory = make([]int64, len(p.program))
		p.Reset()
	}
}

func (p *IntcodeProgram) getParamValue(param, mode int) int64 {
	switch mode {
	case modePosition:
		return p.memory[param]

	case modeImmediate:
		return int64(param)
	}

	panic("unhandled param mode")
}

func (p *IntcodeProgram) GetMemory(idx int) int64 {
	return p.memory[idx]
}

func (p *IntcodeProgram) SetMemory(idx int, val int64) {
	p.init()
	p.memory[idx] = val
}

func (p *IntcodeProgram) Reset() {
	p.init()
	copy(p.memory, p.program)
}

func (p *IntcodeProgram) Run() {
	p.RunIn(func(int) int64 { return 0 }, func(int64, IntcodeProgramState) {})
}

func (p *IntcodeProgram) RunIn(inputFunc ProvideInputFunc, outputFunc ReceiveOutputFunc) {
	p.init()

	inputsRequested := 0
	for instructionPointer := 0; instructionPointer < len(p.program); {
		instruction := p.memory[instructionPointer]
		instructionPointer++

		paramModes := [3]int{
			modePosition,
			modePosition,
			modePosition,
		}
		modes := instruction / 100
		for i := 0; modes > 0; i++ {
			paramModes[i] = int(modes % 10)
			modes = modes / 10
		}

		opcode := instruction % 100
		switch opcode {
		case opAdd:
			param1 := p.memory[instructionPointer]
			param2 := p.memory[instructionPointer+1]
			param3 := p.memory[instructionPointer+2]
			p.memory[param3] = p.getParamValue(int(param1), paramModes[0]) + p.getParamValue(int(param2), paramModes[1])

			instructionPointer += 3

		case opMultiply:
			param1 := p.memory[instructionPointer]
			param2 := p.memory[instructionPointer+1]
			param3 := p.memory[instructionPointer+2]
			p.memory[param3] = p.getParamValue(int(param1), paramModes[0]) * p.getParamValue(int(param2), paramModes[1])

			instructionPointer += 3

		case opInput:
			inputsRequested++
			param1 := p.memory[instructionPointer]
			p.memory[param1] = inputFunc(inputsRequested)

			instructionPointer += 1

		case opOutput:
			param1 := p.memory[instructionPointer]
			outputFunc(p.getParamValue(int(param1), paramModes[0]), p.makeState(instructionPointer))

			instructionPointer += 1

		case opJumpIfTrue:
			param1 := p.memory[instructionPointer]
			param2 := p.memory[instructionPointer+1]

			if p.getParamValue(int(param1), paramModes[0]) != 0 {
				instructionPointer = int(p.getParamValue(int(param2), paramModes[1]))
			} else {
				instructionPointer += 2
			}

		case opJumpIfFalse:
			param1 := p.memory[instructionPointer]
			param2 := p.memory[instructionPointer+1]

			if p.getParamValue(int(param1), paramModes[0]) == 0 {
				instructionPointer = int(p.getParamValue(int(param2), paramModes[1]))
			} else {
				instructionPointer += 2
			}

		case opLessThan:
			param1 := p.memory[instructionPointer]
			param2 := p.memory[instructionPointer+1]
			param3 := p.memory[instructionPointer+2]

			if p.getParamValue(int(param1), paramModes[0]) < p.getParamValue(int(param2), paramModes[1]) {
				p.memory[param3] = 1
			} else {
				p.memory[param3] = 0
			}

			instructionPointer += 3

		case opEquals:
			param1 := p.memory[instructionPointer]
			param2 := p.memory[instructionPointer+1]
			param3 := p.memory[instructionPointer+2]

			if p.getParamValue(int(param1), paramModes[0]) == p.getParamValue(int(param2), paramModes[1]) {
				p.memory[param3] = 1
			} else {
				p.memory[param3] = 0
			}

			instructionPointer += 3

		case opHalt:
			instructionPointer = len(p.program)

		default:
			panic(fmt.Sprintf("exception executing program - unhandled opcode %d", opcode))
		}
	}
}
