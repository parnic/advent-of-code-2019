package days

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"parnic.com/aoc2019/utilities"
)

type Pair[T, U any] struct {
	a T
	b U
}

type Day03 struct {
	line1    []Pair[byte, int]
	line2    []Pair[byte, int]
	visited  map[Pair[int, int]]int
	overlaps []Pair[Pair[int, int], int]
}

func (d *Day03) Parse() {
	lines := utilities.GetStringLines("03p")

	line1data := strings.Split(lines[0], ",")
	line2data := strings.Split(lines[1], ",")

	d.line1 = make([]Pair[byte, int], len(line1data))
	d.line2 = make([]Pair[byte, int], len(line2data))

	for idx, instr := range line1data {
		dir := instr[0]
		amt := instr[1:]
		iAmt, err := strconv.Atoi(amt)
		if err != nil {
			panic(err)
		}

		d.line1[idx] = Pair[byte, int]{a: dir, b: iAmt}
	}

	for idx, instr := range line2data {
		dir := instr[0]
		amt := instr[1:]
		iAmt, err := strconv.Atoi(amt)
		if err != nil {
			panic(err)
		}

		d.line2[idx] = Pair[byte, int]{a: dir, b: iAmt}
	}
}

func (d Day03) Num() int {
	return 3
}

func (d *Day03) Part1() string {
	d.visited = make(map[Pair[int, int]]int)
	var x int
	var y int
	var steps int
	for _, inst := range d.line1 {
		switch inst.a {
		case 'R':
			for i := 1; i <= inst.b; i++ {
				steps++
				d.visited[Pair[int, int]{a: x + i, b: y}] = steps
			}
			x += inst.b
		case 'U':
			for i := 1; i <= inst.b; i++ {
				steps++
				d.visited[Pair[int, int]{a: x, b: y + i}] = steps
			}
			y += inst.b
		case 'L':
			for i := 1; i <= inst.b; i++ {
				steps++
				d.visited[Pair[int, int]{a: x - i, b: y}] = steps
			}
			x -= inst.b
		case 'D':
			for i := 1; i <= inst.b; i++ {
				steps++
				d.visited[Pair[int, int]{a: x, b: y - i}] = steps
			}
			y -= inst.b
		}
	}

	x = 0
	y = 0
	steps = 0
	d.overlaps = make([]Pair[Pair[int, int], int], 0)
	for _, inst := range d.line2 {
		switch inst.a {
		case 'R':
			for i := 1; i <= inst.b; i++ {
				steps++
				if _, exists := d.visited[Pair[int, int]{x + i, y}]; exists {
					d.overlaps = append(d.overlaps, Pair[Pair[int, int], int]{a: Pair[int, int]{x + i, y}, b: steps})
				}
			}
			x += inst.b
		case 'U':
			for i := 1; i <= inst.b; i++ {
				steps++
				if _, exists := d.visited[Pair[int, int]{x, y + i}]; exists {
					d.overlaps = append(d.overlaps, Pair[Pair[int, int], int]{a: Pair[int, int]{x, y + i}, b: steps})
				}
			}
			y += inst.b
		case 'L':
			for i := 1; i <= inst.b; i++ {
				steps++
				if _, exists := d.visited[Pair[int, int]{x - i, y}]; exists {
					d.overlaps = append(d.overlaps, Pair[Pair[int, int], int]{a: Pair[int, int]{x - i, y}, b: steps})
				}
			}
			x -= inst.b
		case 'D':
			for i := 1; i <= inst.b; i++ {
				steps++
				if _, exists := d.visited[Pair[int, int]{x, y - i}]; exists {
					d.overlaps = append(d.overlaps, Pair[Pair[int, int], int]{a: Pair[int, int]{x, y - i}, b: steps})
				}
			}
			y -= inst.b
		}
	}

	minDist := math.MaxInt
	for _, overlap := range d.overlaps {
		dist := int(math.Abs(float64(overlap.a.a))) + int(math.Abs(float64(overlap.a.b)))
		if dist < minDist {
			minDist = dist
		}
	}

	return fmt.Sprintf("Closest overlap manhattan distance = %s%d%s", utilities.TextBold, minDist, utilities.TextReset)
}

func (d *Day03) Part2() string {
	minOverlap := math.MaxInt
	for _, overlap := range d.overlaps {
		line1Steps := d.visited[overlap.a]
		line2Steps := overlap.b

		totalSteps := line1Steps + line2Steps
		if totalSteps < minOverlap {
			minOverlap = totalSteps
		}
	}

	return fmt.Sprintf("Minimum steps to overlap = %s%d%s", utilities.TextBold, minOverlap, utilities.TextReset)
}
