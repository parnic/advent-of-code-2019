package days

import (
	"fmt"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

type Day11 struct {
	program u.IntcodeProgram
	painted map[u.Pair[int, int]]int
}

func (d *Day11) Parse() {
	d.program = u.LoadIntcodeProgram("11p")
}

func (d Day11) Num() int {
	return 11
}

func (d *Day11) paintHull() (int, u.Pair[int, int], u.Pair[int, int]) {
	pos := u.Pair[int, int]{First: 0, Second: 0}
	facing := 0

	min := pos
	max := pos

	outputState := 0
	numPainted := 1
	d.program.RunIn(func(inputStep int) int64 {
		return int64(d.painted[pos])
	}, func(val int64, state u.IntcodeProgramState) {
		if outputState == 0 {
			outputState++
			color := int(val)
			if _, exists := d.painted[pos]; !exists {
				numPainted++
			}
			d.painted[pos] = color
		} else {
			outputState = 0
			direction := val

			if direction == 0 {
				facing--
				if facing == -1 {
					facing = 3
				}
			} else {
				facing++
				if facing == 4 {
					facing = 0
				}
			}

			switch facing {
			case 0:
				pos.First--
				if pos.First < min.First {
					min.First = pos.First
				}
			case 1:
				pos.Second++
				if pos.Second > max.Second {
					max.Second = pos.Second
				}
			case 2:
				pos.First++
				if pos.First > max.First {
					max.First = pos.First
				}
			case 3:
				pos.Second--
				if pos.Second < min.Second {
					min.Second = pos.Second
				}
			}
		}
	})

	return numPainted, min, max
}

func (d *Day11) Part1() string {
	d.painted = map[u.Pair[int, int]]int{
		{First: 0, Second: 0}: 0,
	}
	numPainted, _, _ := d.paintHull()

	return fmt.Sprintf("Unique panels painted: %s%d%s", u.TextBold, numPainted, u.TextReset)
}

func (d *Day11) Part2() string {
	d.painted = map[u.Pair[int, int]]int{
		{First: 0, Second: 0}: 1,
	}
	_, min, max := d.paintHull()

	outStr := strings.Builder{}
	outStr.WriteString("Registration identifier:\n")
	outStr.WriteString(u.TextBold)
	for x := min.First; x <= max.First; x++ {
		for y := min.Second; y <= max.Second; y++ {
			val, exists := d.painted[u.Pair[int, int]{First: x, Second: y}]
			if exists && val == 1 {
				outStr.WriteRune('█')
			} else {
				outStr.WriteRune(' ')
			}
		}
		outStr.WriteRune('\n')
	}
	outStr.WriteString(u.TextReset)

	return outStr.String()
}
