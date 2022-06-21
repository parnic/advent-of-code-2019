package days

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

type Day03 struct {
	line1    []u.Pair[byte, int]
	line2    []u.Pair[byte, int]
	visited  map[u.Pair[int, int]]int
	overlaps []u.Pair[u.Pair[int, int], int]
}

func (d *Day03) Parse() {
	lines := u.GetStringLines("03p")

	line1data := strings.Split(lines[0], ",")
	line2data := strings.Split(lines[1], ",")

	d.line1 = make([]u.Pair[byte, int], len(line1data))
	d.line2 = make([]u.Pair[byte, int], len(line2data))

	for idx, instr := range line1data {
		dir := instr[0]
		amt := instr[1:]
		iAmt, err := strconv.Atoi(amt)
		if err != nil {
			panic(err)
		}

		d.line1[idx] = u.Pair[byte, int]{First: dir, Second: iAmt}
	}

	for idx, instr := range line2data {
		dir := instr[0]
		amt := instr[1:]
		iAmt, err := strconv.Atoi(amt)
		if err != nil {
			panic(err)
		}

		d.line2[idx] = u.Pair[byte, int]{First: dir, Second: iAmt}
	}
}

func (d Day03) Num() int {
	return 3
}

func (d *Day03) Part1() string {
	d.visited = make(map[u.Pair[int, int]]int)
	var x int
	var y int
	var steps int
	for _, inst := range d.line1 {
		switch inst.First {
		case 'R':
			for i := 1; i <= inst.Second; i++ {
				steps++
				d.visited[u.Pair[int, int]{First: x + i, Second: y}] = steps
			}
			x += inst.Second
		case 'U':
			for i := 1; i <= inst.Second; i++ {
				steps++
				d.visited[u.Pair[int, int]{First: x, Second: y + i}] = steps
			}
			y += inst.Second
		case 'L':
			for i := 1; i <= inst.Second; i++ {
				steps++
				d.visited[u.Pair[int, int]{First: x - i, Second: y}] = steps
			}
			x -= inst.Second
		case 'D':
			for i := 1; i <= inst.Second; i++ {
				steps++
				d.visited[u.Pair[int, int]{First: x, Second: y - i}] = steps
			}
			y -= inst.Second
		}
	}

	x = 0
	y = 0
	steps = 0
	d.overlaps = make([]u.Pair[u.Pair[int, int], int], 0)
	for _, inst := range d.line2 {
		switch inst.First {
		case 'R':
			for i := 1; i <= inst.Second; i++ {
				steps++
				if _, exists := d.visited[u.Pair[int, int]{First: x + i, Second: y}]; exists {
					d.overlaps = append(d.overlaps, u.Pair[u.Pair[int, int], int]{First: u.Pair[int, int]{First: x + i, Second: y}, Second: steps})
				}
			}
			x += inst.Second
		case 'U':
			for i := 1; i <= inst.Second; i++ {
				steps++
				if _, exists := d.visited[u.Pair[int, int]{First: x, Second: y + i}]; exists {
					d.overlaps = append(d.overlaps, u.Pair[u.Pair[int, int], int]{First: u.Pair[int, int]{First: x, Second: y + i}, Second: steps})
				}
			}
			y += inst.Second
		case 'L':
			for i := 1; i <= inst.Second; i++ {
				steps++
				if _, exists := d.visited[u.Pair[int, int]{First: x - i, Second: y}]; exists {
					d.overlaps = append(d.overlaps, u.Pair[u.Pair[int, int], int]{First: u.Pair[int, int]{First: x - i, Second: y}, Second: steps})
				}
			}
			x -= inst.Second
		case 'D':
			for i := 1; i <= inst.Second; i++ {
				steps++
				if _, exists := d.visited[u.Pair[int, int]{First: x, Second: y - i}]; exists {
					d.overlaps = append(d.overlaps, u.Pair[u.Pair[int, int], int]{First: u.Pair[int, int]{First: x, Second: y - i}, Second: steps})
				}
			}
			y -= inst.Second
		}
	}

	minDist := math.MaxInt
	for _, overlap := range d.overlaps {
		dist := int(math.Abs(float64(overlap.First.First))) + int(math.Abs(float64(overlap.First.Second)))
		if dist < minDist {
			minDist = dist
		}
	}

	return fmt.Sprintf("Closest overlap manhattan distance = %s%d%s", u.TextBold, minDist, u.TextReset)
}

func (d *Day03) Part2() string {
	minOverlap := math.MaxInt
	for _, overlap := range d.overlaps {
		line1Steps := d.visited[overlap.First]
		line2Steps := overlap.Second

		totalSteps := line1Steps + line2Steps
		if totalSteps < minOverlap {
			minOverlap = totalSteps
		}
	}

	return fmt.Sprintf("Minimum steps to overlap = %s%d%s", u.TextBold, minOverlap, u.TextReset)
}
