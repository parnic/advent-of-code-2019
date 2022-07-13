package days

import (
	"fmt"

	u "parnic.com/aoc2019/utilities"
)

type Day19 struct {
	program u.IntcodeProgram
}

func (d *Day19) Parse() {
	d.program = u.LoadIntcodeProgram("19p")
}

func (d Day19) Num() int {
	return 19
}

func (d *Day19) Part1() string {
	grid := make([][]bool, 50)
	for y := 0; y < len(grid); y++ {
		grid[y] = make([]bool, 50)
	}

	count := int64(0)

	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			d.program.Reset()
			d.program.RunIn(func(inputStep int) int64 {
				if inputStep == 1 {
					return int64(x)
				}

				return int64(y)
			}, func(val int64, state u.IntcodeProgramState) {
				res := val == 1
				grid[y][x] = res
				if res {
					count++
				}
			})
		}
	}

	// fmt.Println("50x50 tractor view:")
	// for y := 0; y < len(grid); y++ {
	// 	for x := 0; x < len(grid[y]); x++ {
	// 		if grid[y][x] {
	// 			fmt.Print("█")
	// 		} else {
	// 			fmt.Print(" ")
	// 		}
	// 	}
	// 	fmt.Println()
	// }

	return fmt.Sprintf("Points affected in 50x50 area: %s%d%s", u.TextBold, count, u.TextReset)
}

func (d *Day19) Part2() string {
	f := func(x, y int) bool {
		ret := false
		d.program.Reset()
		d.program.RunIn(func(inputStep int) int64 {
			if inputStep == 1 {
				return int64(x)
			}
			return int64(y)
		}, func(val int64, state u.IntcodeProgramState) {
			ret = val == 1
		})

		return ret
	}

	// find lower bound
	// this may not be necessary, but helps seed the bisect with a known-good lower bound
	startY := 0
	startX := 0
	for y := 1; startY == 0; y++ {
		for x := 0; x < 10*y; x++ {
			if f(x, y) {
				startY = y
				startX = x
				break
			}
		}
	}

	lastGoodX := 0
	threshold := 1
	// add 100 to start y since we know it has to be a 100x100 square.
	// since we multiply x by 10,000 for the final result, that tells us y cannot be 10k+
	y := u.Bisect(startY+100, 9999, threshold, func(y int) bool {
		foundX := false
		for x := startX; ; x++ {
			// check top left
			if !f(x, y) {
				if !foundX {
					continue
				} else {
					return true
				}
			}
			if !foundX {
				foundX = true
			}

			// check top right
			if !f(x+99, y) {
				return true
			}
			// check bottom left
			if !f(x, y+99) {
				continue
			}

			// we believe the corners work, so run final validations on the full borders.
			// this may not be necessary, but i've seen some rows end up shorter than a
			// previous row because of the angle of the beam and our integer fidelity.
			// plus it's really not that much more expensive to do to be certain we're correct.
			for y2 := y; y2 < y+100; y2++ {
				// right wall
				if !f(x+99, y2) {
					return true
				}
				// left wall
				if !f(x, y2) {
					return true
				}
			}

			lastGoodX = x
			return false
		}
	})

	// since we invert our bisect success returns, we need to increment y to
	// tip back over into the "success" range
	y += threshold
	result := (lastGoodX * 10000) + y
	return fmt.Sprintf("Closest 100x100 square for the ship starts at %d,%d = %s%d%s", lastGoodX, y, u.TextBold, result, u.TextReset)
}
