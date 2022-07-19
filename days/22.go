package days

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

type day22Instruction int

const (
	day22InstructionNewStack day22Instruction = iota
	day22InstructionCut
	day22InstructionDealIncrement
)

type day22Shuffle struct {
	instruction day22Instruction
	arg         int
}

type Day22 struct {
	shuffles []day22Shuffle
}

func (d *Day22) Parse() {
	lines := u.GetStringLines("22p")
	d.shuffles = make([]day22Shuffle, len(lines))

	for idx, line := range lines {
		split := strings.Split(line, " ")
		if split[0] == "deal" {
			if split[1] == "into" {
				d.shuffles[idx] = day22Shuffle{instruction: day22InstructionNewStack}
			} else if split[1] == "with" {
				arg, err := strconv.Atoi(split[3])
				if err != nil {
					panic(err)
				}
				d.shuffles[idx] = day22Shuffle{
					instruction: day22InstructionDealIncrement,
					arg:         arg,
				}
			}
		} else if split[0] == "cut" {
			arg, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
			}
			d.shuffles[idx] = day22Shuffle{
				instruction: day22InstructionCut,
				arg:         arg,
			}
		}
	}
}

func (d Day22) Num() int {
	return 22
}

func (d Day22) applyShuffle(s day22Shuffle, stack, scratch []int) {
	switch s.instruction {
	case day22InstructionNewStack:
		for i := 0; i < len(stack)/2; i++ {
			stack[i], stack[len(stack)-1-i] = stack[len(stack)-1-i], stack[i]
		}

		// there's probably a way to do these two in place...
	case day22InstructionCut:
		absArg := int(math.Abs(float64(s.arg)))
		for i, v := range stack {
			if s.arg > 0 {
				if i < absArg {
					scratch[len(scratch)-absArg+i] = v
				} else {
					scratch[i-absArg] = v
				}
			} else {
				if i < absArg {
					scratch[i] = stack[len(stack)-absArg+i]
				} else {
					scratch[i] = stack[i-absArg]
				}
			}
		}
		copy(stack, scratch)

	case day22InstructionDealIncrement:
		for i, v := range stack {
			scratch[(i*s.arg)%len(stack)] = v
		}
		copy(stack, scratch)
	}
}

func (d *Day22) Part1() string {
	deckSize := 10007
	// deckSize := 10

	stack := make([]int, deckSize)
	for i := range stack {
		stack[i] = i
	}

	scratch := make([]int, len(stack))

	for _, s := range d.shuffles {
		d.applyShuffle(s, stack, scratch)
	}

	pos := -1
	for i, v := range stack {
		if v == 2019 {
			pos = i
			break
		}
	}

	return fmt.Sprintf("Card 2019 is at position %s%d%s", u.TextBold, pos, u.TextReset)
}

func (d *Day22) Part2() string {
	n, iter := big.NewInt(119315717514047), big.NewInt(101741582076661)
	offset, increment := big.NewInt(0), big.NewInt(1)
	for _, s := range d.shuffles {
		switch s.instruction {
		case day22InstructionNewStack:
			increment.Mul(increment, big.NewInt(-1))
			offset.Add(offset, increment)
		case day22InstructionCut:
			offset.Add(offset, big.NewInt(0).Mul(big.NewInt(int64(s.arg)), increment))
		case day22InstructionDealIncrement:
			increment.Mul(increment, big.NewInt(0).Exp(big.NewInt(int64(s.arg)), big.NewInt(0).Sub(n, big.NewInt(2)), n))
		}
	}

	finalIncr := big.NewInt(0).Exp(increment, iter, n)

	finalOffs := big.NewInt(0).Exp(increment, iter, n)
	finalOffs.Sub(big.NewInt(1), finalOffs)
	invmod := big.NewInt(0).Exp(big.NewInt(0).Sub(big.NewInt(1), increment), big.NewInt(0).Sub(n, big.NewInt(2)), n)
	finalOffs.Mul(finalOffs, invmod)
	finalOffs.Mul(finalOffs, offset)

	answer := big.NewInt(0).Mul(big.NewInt(2020), finalIncr)
	answer.Add(answer, finalOffs)
	answer.Mod(answer, n)

	return fmt.Sprintf("Card at position 2020: %s%d%s", u.TextBold, answer, u.TextReset)
}
