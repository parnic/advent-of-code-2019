package utilities

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	opAdd          = 1
	opMultiply     = 2
	opInput        = 3
	opOutput       = 4
	opJumpIfTrue   = 5
	opJumpIfFalse  = 6
	opLessThan     = 7
	opEquals       = 8
	opRelativeBase = 9
	opHalt         = 99

	modePosition  = 0
	modeImmediate = 1
	modeRelative  = 2
)

type IntcodeProgram struct {
	memory        []int64
	program       []int64
	relativeBase  int
	haltRequested bool
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
		copy(p.memory, p.program)
	}
}

func (p *IntcodeProgram) getParamValue(param, mode int) int64 {
	switch mode {
	case modePosition:
		return p.GetMemory(param)

	case modeImmediate:
		return int64(param)

	case modeRelative:
		return p.GetMemory(param + p.relativeBase)
	}

	panic("unhandled param mode")
}

func (p *IntcodeProgram) GetMemory(idx int) int64 {
	p.ensureMemoryCapacity(idx)
	return p.memory[idx]
}

func (p *IntcodeProgram) SetMemory(idx int, val int64) {
	p.init()
	p.ensureMemoryCapacity(idx)
	p.memory[idx] = val
}

func (p *IntcodeProgram) setMemory(idx int, val int64, mode int) {
	if mode == modeImmediate {
		panic("exception executing program - write parameter must never be in immediate mode")
	}
	if mode == modeRelative {
		idx = idx + p.relativeBase
	}

	p.SetMemory(idx, val)
}

func (p *IntcodeProgram) ensureMemoryCapacity(address int) {
	if len(p.memory) > address {
		return
	}

	p.memory = append(p.memory, make([]int64, address+1-len(p.memory))...)
}

func (p *IntcodeProgram) Reset() {
	wiped := false
	if len(p.memory) != len(p.program) {
		wiped = true
		p.memory = nil
	}
	p.init()
	if !wiped {
		copy(p.memory, p.program)
	}
	p.relativeBase = 0
}

func (p *IntcodeProgram) Run() {
	p.RunIn(func(int) int64 { return 0 }, func(int64, IntcodeProgramState) {})
}

func (p *IntcodeProgram) RunIn(inputFunc ProvideInputFunc, outputFunc ReceiveOutputFunc) {
	p.init()

	inputsRequested := 0
	for instructionPointer := 0; instructionPointer < len(p.program) && !p.haltRequested; {
		instruction := p.GetMemory(instructionPointer)
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
			param1 := p.GetMemory(instructionPointer)
			param2 := p.GetMemory(instructionPointer + 1)
			param3 := p.GetMemory(instructionPointer + 2)
			p.setMemory(int(param3), p.getParamValue(int(param1), paramModes[0])+p.getParamValue(int(param2), paramModes[1]), paramModes[2])

			instructionPointer += 3

		case opMultiply:
			param1 := p.GetMemory(instructionPointer)
			param2 := p.GetMemory(instructionPointer + 1)
			param3 := p.GetMemory(instructionPointer + 2)
			p.setMemory(int(param3), p.getParamValue(int(param1), paramModes[0])*p.getParamValue(int(param2), paramModes[1]), paramModes[2])

			instructionPointer += 3

		case opInput:
			inputsRequested++
			param1 := p.GetMemory(instructionPointer)
			p.setMemory(int(param1), inputFunc(inputsRequested), paramModes[0])

			instructionPointer += 1

		case opOutput:
			param1 := p.GetMemory(instructionPointer)
			outputFunc(p.getParamValue(int(param1), paramModes[0]), p.makeState(instructionPointer))

			instructionPointer += 1

		case opJumpIfTrue:
			param1 := p.GetMemory(instructionPointer)
			param2 := p.GetMemory(instructionPointer + 1)

			if p.getParamValue(int(param1), paramModes[0]) != 0 {
				instructionPointer = int(p.getParamValue(int(param2), paramModes[1]))
			} else {
				instructionPointer += 2
			}

		case opJumpIfFalse:
			param1 := p.GetMemory(instructionPointer)
			param2 := p.GetMemory(instructionPointer + 1)

			if p.getParamValue(int(param1), paramModes[0]) == 0 {
				instructionPointer = int(p.getParamValue(int(param2), paramModes[1]))
			} else {
				instructionPointer += 2
			}

		case opLessThan:
			param1 := p.GetMemory(instructionPointer)
			param2 := p.GetMemory(instructionPointer + 1)
			param3 := p.GetMemory(instructionPointer + 2)

			if p.getParamValue(int(param1), paramModes[0]) < p.getParamValue(int(param2), paramModes[1]) {
				p.setMemory(int(param3), 1, paramModes[2])
			} else {
				p.setMemory(int(param3), 0, paramModes[2])
			}

			instructionPointer += 3

		case opEquals:
			param1 := p.GetMemory(instructionPointer)
			param2 := p.GetMemory(instructionPointer + 1)
			param3 := p.GetMemory(instructionPointer + 2)

			if p.getParamValue(int(param1), paramModes[0]) == p.getParamValue(int(param2), paramModes[1]) {
				p.setMemory(int(param3), 1, paramModes[2])
			} else {
				p.setMemory(int(param3), 0, paramModes[2])
			}

			instructionPointer += 3

		case opRelativeBase:
			param1 := p.GetMemory(instructionPointer)

			p.relativeBase += int(p.getParamValue(int(param1), paramModes[0]))

			instructionPointer += 1

		case opHalt:
			instructionPointer = len(p.program)

		default:
			panic(fmt.Sprintf("exception executing program - unhandled opcode %d", opcode))
		}
	}

	p.haltRequested = false
}

func (p *IntcodeProgram) Stop() {
	p.haltRequested = true
}
