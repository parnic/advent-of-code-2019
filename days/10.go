package days

import (
	"fmt"
	"sort"

	u "parnic.com/aoc2019/utilities"
)

type Day10 struct {
	asteroids     [][]bool
	idealLocation u.Vec2[int]
}

func (d *Day10) Parse() {
	lines := u.GetStringLines("10p")
	d.asteroids = make([][]bool, len(lines))
	for i, line := range lines {
		d.asteroids[i] = make([]bool, len(line))
		for j, ch := range line {
			d.asteroids[i][j] = ch == '#'
		}
	}
}

func (d Day10) Num() int {
	return 10
}

// func (d Day10) draw() {
// 	for i := range d.asteroids {
// 		for j := range d.asteroids[i] {
// 			if !d.asteroids[i][j].First {
// 				fmt.Print(".")
// 			} else {
// 				num := d.asteroids[i][j].Second
// 				ch := rune('0') + rune(num)
// 				if num >= 10 {
// 					ch = '+'
// 				}
// 				fmt.Printf("%c", ch)
// 			}
// 		}
// 		fmt.Println()
// 	}
// }

func (d Day10) getVisibleAsteroids(i1, j1 int) map[u.Vec2[int]]bool {
	visited := make(map[u.Vec2[int]]bool, 0)
	foundAsteroids := make(map[u.Vec2[int]]bool, 0)

	findNext := func(startX, startY, incX, incY int) *u.Vec2[int] {
		var found *u.Vec2[int]
		if incX == 0 && incY == 0 {
			return found
		}

		x := startX + incX
		y := startY + incY
		for x < len(d.asteroids) && x >= 0 && y < len(d.asteroids[x]) && y >= 0 {
			currPair := u.Vec2[int]{X: x, Y: y}
			if _, exists := visited[currPair]; !exists {
				visited[currPair] = true

				if d.asteroids[x][y] {
					if found == nil {
						found = &currPair
					}
				}
			}

			x += incX
			y += incY
		}

		return found
	}

	for incX := 0; ; {
		plusXValid := i1+incX < len(d.asteroids)
		minusXValid := i1-incX >= 0
		if !plusXValid && !minusXValid {
			break
		}

		for incY := 0; ; {
			plusYValid := j1+incY < len(d.asteroids[0])
			minusYValid := j1-incY >= 0
			if !plusYValid && !minusYValid {
				break
			}

			if found := findNext(i1, j1, incX, incY); found != nil {
				foundAsteroids[*found] = true
			}
			if found := findNext(i1, j1, incX, -incY); found != nil {
				foundAsteroids[*found] = true
			}
			if found := findNext(i1, j1, -incX, incY); found != nil {
				foundAsteroids[*found] = true
			}
			if found := findNext(i1, j1, -incX, -incY); found != nil {
				foundAsteroids[*found] = true
			}

			incY++
		}

		incX++
	}

	return foundAsteroids
}

func (d Day10) numVisibleAsteroids(i1, j1 int) int {
	return len(d.getVisibleAsteroids(i1, j1))
}

func (d *Day10) removeAsteroids(locs map[u.Vec2[int]]bool) {
	for loc := range locs {
		if !d.asteroids[loc.X][loc.Y] {
			panic("tried to remove non-asteroid")
		}

		d.asteroids[loc.X][loc.Y] = false
	}
}

func (d *Day10) Part1() string {
	mostAsteroids := 0
	for i := range d.asteroids {
		for j := range d.asteroids[i] {
			if d.asteroids[i][j] {
				numVisible := d.numVisibleAsteroids(i, j)
				if numVisible > mostAsteroids {
					mostAsteroids = numVisible
					d.idealLocation = u.Vec2[int]{X: i, Y: j}
				}
			}
		}
	}

	return fmt.Sprintf("Most visible asteroids: %s%d%s at (%d,%d)", u.TextBold, mostAsteroids, u.TextReset, d.idealLocation.Y, d.idealLocation.X)
}

func (d *Day10) Part2() string {
	findNumVaporized := 200
	var targetLocation u.Vec2[int]

	vaporized := 0
	for vaporized < findNumVaporized {
		visibleAsteroids := d.getVisibleAsteroids(d.idealLocation.X, d.idealLocation.Y)
		if len(visibleAsteroids) == 0 {
			panic("no more asteroids to vaporize")
		}

		if vaporized+len(visibleAsteroids) < findNumVaporized {
			vaporized += len(visibleAsteroids)
			d.removeAsteroids(visibleAsteroids)
			continue
		}

		vecs := u.MapKeys(visibleAsteroids)
		sort.Slice(vecs, func(i, j int) bool {
			return d.idealLocation.AngleBetween(vecs[i]) > d.idealLocation.AngleBetween(vecs[j])
		})
		targetLocation = vecs[findNumVaporized-1-vaporized]
		break
	}

	return fmt.Sprintf("#%d asteroid to be vaporized is at (%d,%d), transformed: %s%d%s", findNumVaporized, targetLocation.Y, targetLocation.X, u.TextBold, (targetLocation.Y*100)+targetLocation.X, u.TextReset)
}
