package days

import (
	"fmt"
	"math"

	u "parnic.com/aoc2019/utilities"
)

var (
	day24AdjacentOffsets = []u.Vec2i{
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: 0, Y: 1},
	}
)

type Day24 struct {
	grid [][]bool
}

func (d *Day24) Parse() {
	lines := u.GetStringLines("24p")
	d.grid = make([][]bool, len(lines))
	for i, line := range lines {
		d.grid[i] = make([]bool, len(line))
		for j, ch := range line {
			d.grid[i][j] = ch == '#'
		}
	}
}

func (d Day24) Num() int {
	return 24
}

func (d Day24) calcActivatedNeighbors(grid [][]bool, i, j int) int {
	activatedNeighbors := 0
	for _, o := range day24AdjacentOffsets {
		newI := i + o.X
		newJ := j + o.Y
		if newI < 0 || newI >= len(grid) || newJ < 0 || newJ >= len(grid[i]) {
			continue
		}
		if grid[newI][newJ] {
			activatedNeighbors++
		}
	}

	return activatedNeighbors
}

func (d Day24) recursiveCalcActivatedNeighbors(gridMap map[int][][]bool, mapIdx, i, j int) int {
	activatedNeighbors := 0
	numNeighbors := 0
	thisGrid := gridMap[mapIdx]
	for _, o := range day24AdjacentOffsets {
		newI := i + o.X
		newJ := j + o.Y
		if newI < 0 || newI >= len(thisGrid) || newJ < 0 || newJ >= len(thisGrid[i]) {
			continue
		}
		if newI == 2 && newJ == 2 {
			continue
		}
		numNeighbors++
		if thisGrid[newI][newJ] {
			activatedNeighbors++
		}
	}

	checkLower := (i == 1 && j == 2) ||
		(i == 2 && (j == 1 || j == 3)) ||
		(i == 3 && j == 2)
	if checkLower {
		if lowerGrid, exists := gridMap[mapIdx+1]; exists {
			if i == 1 {
				for _, b := range lowerGrid[0] {
					numNeighbors++
					if b {
						activatedNeighbors++
					}
				}
			} else if i == 2 {
				if j == 1 {
					for _, r := range lowerGrid {
						numNeighbors++
						if r[0] {
							activatedNeighbors++
						}
					}
				} else if j == 3 {
					for _, r := range lowerGrid {
						numNeighbors++
						if r[len(lowerGrid[0])-1] {
							activatedNeighbors++
						}
					}
				}
			} else if i == 3 {
				for _, b := range lowerGrid[len(lowerGrid)-1] {
					numNeighbors++
					if b {
						activatedNeighbors++
					}
				}
			}
		}
	}

	checkUpper := (i == 0) || (i == len(thisGrid)-1) ||
		((i != 0 && i != len(thisGrid)) && (j == 0 || j == len(thisGrid[0])-1))
	if checkUpper {
		if upperGrid, exists := gridMap[mapIdx-1]; exists {
			if i == 0 {
				numNeighbors++
				if upperGrid[1][2] {
					activatedNeighbors++
				}
			} else if i == len(thisGrid)-1 {
				numNeighbors++
				if upperGrid[3][2] {
					activatedNeighbors++
				}
			}
			if j == 0 {
				numNeighbors++
				if upperGrid[2][1] {
					activatedNeighbors++
				}
			} else if j == len(thisGrid[0])-1 {
				numNeighbors++
				if upperGrid[2][3] {
					activatedNeighbors++
				}
			}
		}
	}

	return activatedNeighbors
}

func (d Day24) calcRating(grid [][]bool) int {
	rating := 0
	for i, r := range grid {
		for j := range r {
			pow := (i * len(r)) + j
			if grid[i][j] {
				result := int(math.Pow(2, float64(pow)))
				rating += result
			}
		}
	}
	return rating
}

func (d Day24) getNumBugs(gridMap map[int][][]bool) int {
	ret := 0
	for _, v := range gridMap {
		for _, r := range v {
			for _, b := range r {
				if b {
					ret++
				}
			}
		}
	}
	return ret
}

func copy2d[T comparable](dest [][]T, src [][]T) {
	for i, r := range src {
		copy(dest[i], r)
	}
}

func (d Day24) Draw(grid [][]bool) {
	for _, r := range grid {
		for _, c := range r {
			if c {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (d *Day24) Part1() string {
	grid := make([][]bool, len(d.grid))
	scratch := make([][]bool, len(grid))
	for i, g := range d.grid {
		grid[i] = make([]bool, len(g))
		scratch[i] = make([]bool, len(g))
		copy(grid[i], d.grid[i])
	}

	found := false
	answer := 0
	seenRatings := make([]int, 0)
	for i := 1; !found; i++ {
		// d.Draw(grid)
		for i, r := range grid {
			for j := range r {
				numActivated := d.calcActivatedNeighbors(grid, i, j)
				if grid[i][j] {
					scratch[i][j] = numActivated == 1
				} else {
					scratch[i][j] = numActivated == 1 || numActivated == 2
				}
			}
		}

		rating := d.calcRating(scratch)
		if u.ArrayContains(seenRatings, rating) {
			found = true
			// d.Draw(scratch)
			answer = rating
		}
		seenRatings = append(seenRatings, rating)
		copy2d(grid, scratch)
	}

	return fmt.Sprintf("First repeated biodiversity rating is %s%d%s", u.TextBold, answer, u.TextReset)
}

func (d *Day24) Part2() string {
	makeGrid := func(initialGrid [][]bool) ([][]bool, [][]bool) {
		grid := make([][]bool, len(d.grid))
		scratch := make([][]bool, len(grid))
		for i, g := range d.grid {
			grid[i] = make([]bool, len(g))
			scratch[i] = make([]bool, len(g))
			if initialGrid != nil {
				copy(grid[i], initialGrid[i])
			}
		}

		return grid, scratch
	}

	gridMap := make(map[int][][]bool)
	scratchMap := make(map[int][][]bool)
	gridMap[0], scratchMap[0] = makeGrid(d.grid)

	min := 0
	max := 0

	for i := 0; i < 200; i++ {
		gridMap[min-1], scratchMap[min-1] = makeGrid(nil)
		gridMap[max+1], scratchMap[max+1] = makeGrid(nil)
		min, max = min-1, max+1

		// if i == 10 {
		// 	keys := u.MapKeys(gridMap)
		// 	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
		// 	for _, k := range keys {
		// 		fmt.Println("Depth", k)
		// 		d.Draw(gridMap[k])
		// 	}
		// 	fmt.Println("# bugs:", d.numBugs(gridMap))
		// }

		for depth, grid := range gridMap {
			for i, r := range grid {
				for j := range r {
					if i == 2 && j == 2 {
						continue
					}
					numActivated := d.recursiveCalcActivatedNeighbors(gridMap, depth, i, j)
					if grid[i][j] {
						scratchMap[depth][i][j] = numActivated == 1
					} else {
						scratchMap[depth][i][j] = numActivated == 1 || numActivated == 2
					}
				}
			}
		}

		for d := range gridMap {
			copy2d(gridMap[d], scratchMap[d])
		}
	}

	return fmt.Sprintf("Bugs present after 200 minutes: %s%d%s", u.TextBold, d.getNumBugs(gridMap), u.TextReset)
}
