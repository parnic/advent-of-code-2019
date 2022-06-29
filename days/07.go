package days

import (
	"fmt"
	"sync"

	"parnic.com/aoc2019/utilities"
)

type Day07 struct {
	program utilities.IntcodeProgram
	amps    []utilities.IntcodeProgram
}

func (d *Day07) Parse() {
	d.program = utilities.LoadIntcodeProgram("07p")
	d.amps = make([]utilities.IntcodeProgram, 5)
	for i := range d.amps {
		d.amps[i] = d.program.Copy()
	}
}

func (d Day07) Num() int {
	return 7
}

func (d *Day07) Part1() string {
	var highestVal int64
	var highestSequence []int64

	allSequences := utilities.GetPermutations([]int64{0, 1, 2, 3, 4}...)
	for _, sequence := range allSequences {
		if len(sequence) != len(d.amps) {
			panic("input sequence does not match up to number of amplifiers")
		}

		input := int64(0)
		var output int64
		for i, amp := range d.amps {
			amp.RunIn(func(step int) int64 {
				if step == 1 {
					return sequence[i]
				} else if step == 2 {
					return input
				}

				panic("hit more input instructions than expected")
			}, func(val int64, state utilities.IntcodeProgramState) {
				output = val
			})

			input = output
		}

		if output > highestVal {
			highestVal = output
			if highestSequence == nil {
				highestSequence = make([]int64, len(sequence))
			}
			copy(highestSequence, sequence)
		}
	}

	return fmt.Sprintf("Max thruster signal: %s%d%s (produced by %v)", utilities.TextBold, highestVal, utilities.TextReset, highestSequence)
}

func (d *Day07) Part2() string {
	var highestVal int64
	var highestSequence []int64

	allSequences := utilities.GetPermutations([]int64{5, 6, 7, 8, 9}...)
	for _, sequence := range allSequences {
		if len(sequence) != len(d.amps) {
			panic("input sequence does not match up to number of amplifiers")
		}

		inputs := make([]chan int64, len(d.amps))
		for i := range d.amps {
			d.amps[i].Reset()
			inputs[i] = make(chan int64, 1)
			inputs[i] <- sequence[i]
		}

		var finalOutput int64
		var wg sync.WaitGroup
		for i := range d.amps {
			wg.Add(1)
			go func(idx int) {
				d.amps[idx].RunIn(func(step int) int64 {
					input := <-inputs[idx]
					return input
				}, func(val int64, state utilities.IntcodeProgramState) {
					finalOutput = val
					inputIdx := idx + 1
					if inputIdx == len(inputs) {
						inputIdx = 0
					}
					inputs[inputIdx] <- val
				})
				wg.Done()
			}(i)
		}
		inputs[0] <- 0
		wg.Wait()

		if finalOutput > highestVal {
			highestVal = finalOutput
			if highestSequence == nil {
				highestSequence = make([]int64, len(sequence))
			}
			copy(highestSequence, sequence)
		}
	}

	return fmt.Sprintf("Max thruster signal: %s%d%s (produced by %v)", utilities.TextBold, highestVal, utilities.TextReset, highestSequence)
}
