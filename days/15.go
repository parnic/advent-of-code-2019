package days

//go:generate stringer -type=cellStatus,responseType,dirType -output=15_types_string.go

import (
	"fmt"
	"math"
	"sort"

	u "parnic.com/aoc2019/utilities"
)

type point u.Pair[int, int]
type cellStatus int
type responseType int
type dirType int

const (
	cellStatusUnknown cellStatus = iota
	cellStatusWall
	cellStatusOpen
	cellStatusGoal
)

const (
	responseWall responseType = iota
	responseSuccess
	responseFoundGoal
)

const maxVisited = 3

const (
	dirNorth dirType = iota + 1
	dirSouth
	dirWest
	dirEast

	dirFirst = dirNorth
	dirLast  = dirEast
)

var dirOrders = [][]dirType{
	{dirNorth, dirSouth, dirWest, dirEast},
	{dirSouth, dirWest, dirEast, dirNorth},
	{dirWest, dirEast, dirNorth, dirSouth},
	{dirEast, dirNorth, dirSouth, dirWest},
}

// turned out to be unnecessary on multiple datasets i tried. increases the iterations 6x
// var dirOrders = u.GetPermutations(dirNorth, dirSouth, dirWest, dirEast)

type visitedStatus struct {
	timesVisited      int
	distanceFromStart int
}

type Day15 struct {
	program      u.IntcodeProgram
	grid         map[point]cellStatus
	visited      map[point]*visitedStatus
	shortestPath []point
	pos          point
	goalPos      point
}

func (d *Day15) Parse() {
	d.program = u.LoadIntcodeProgram("15p")
	d.grid = map[point]cellStatus{
		{First: 0, Second: 0}: cellStatusOpen,
	}
	d.visited = map[point]*visitedStatus{
		{First: 0, Second: 0}: {timesVisited: 1},
	}
	d.shortestPath = []point{{}}
}

func (d Day15) Num() int {
	return 15
}

func (d Day15) getPointInDirection(pos point, dir dirType) point {
	target := pos
	switch dir {
	case dirNorth:
		target.First--
	case dirSouth:
		target.First++
	case dirWest:
		target.Second--
	case dirEast:
		target.Second++
	}

	return target
}

func (d Day15) getCellTypeInDirection(pos point, dir dirType) (cellStatus, point) {
	target := d.getPointInDirection(pos, dir)
	return d.grid[target], target
}

func (d Day15) getAdjacentCellsOfType(pos point, cellType cellStatus) []point {
	points := make([]point, 0, 4)
	for i := dirFirst; i <= dirLast; i++ {
		adjacentCell := d.getPointInDirection(pos, i)
		if d.grid[adjacentCell] == cellType {
			points = append(points, adjacentCell)
		}
	}
	return points
}

func (d Day15) getDirToNextCellType(pos point, t cellStatus, maxNumVisited int, dirs []dirType) (dirType, point, error) {
	for _, dir := range dirs {
		cellInDirection, targetCell := d.getCellTypeInDirection(pos, dir)
		if cellInDirection == t {
			_, visitedTargetExists := d.visited[targetCell]
			foundUnknown := t == cellStatusUnknown && !visitedTargetExists
			foundOther := t != cellStatusUnknown && visitedTargetExists && d.visited[targetCell].timesVisited <= maxNumVisited
			if foundUnknown || foundOther {
				return dir, targetCell, nil
			}
		}
	}

	return dirFirst, point{}, fmt.Errorf("no %v tiles around %v", t, pos)
}

func (d *Day15) Draw() {
	min := point{First: math.MaxInt, Second: math.MaxInt}
	max := point{First: math.MinInt, Second: math.MinInt}
	for p := range d.grid {
		if p.First < min.First {
			min.First = p.First
		}
		if p.First > max.First {
			max.First = p.First
		}
		if p.Second < min.Second {
			min.Second = p.Second
		}
		if p.Second > max.Second {
			max.Second = p.Second
		}
	}

	for x := min.First; x <= max.First; x++ {
		for y := min.Second; y <= max.Second; y++ {
			p := point{First: x, Second: y}
			switch d.grid[p] {
			case cellStatusGoal:
				fmt.Printf("%s@%s", u.ColorBrightGreen, u.TextReset)
			case cellStatusOpen:
				if p == d.pos {
					fmt.Print(u.BackgroundBrightRed)
				} else if x == 0 && y == 0 {
					fmt.Print(u.BackgroundYellow)
				} else if u.ArrayContains(d.shortestPath, p) {
					fmt.Print(u.BackgroundGreen)
				} else if d.visited[p] != nil && d.visited[p].timesVisited > maxVisited {
					fmt.Print(u.ColorYellow)
					fmt.Print(u.BackgroundBlack)
				} else if d.visited[p] != nil && d.visited[p].timesVisited > 1 {
					fmt.Print(u.BackgroundMagenta)
				} else {
					fmt.Print(u.BackgroundBlue)
				}
				fmt.Printf(".%s", u.TextReset)
			case cellStatusWall:
				fmt.Print("█")
			case cellStatusUnknown:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (d *Day15) markShortestPath() {
	pos := d.goalPos
	checkOffsets := []point{
		{First: -1, Second: 0},
		{First: 1, Second: 0},
		{First: 0, Second: -1},
		{First: 0, Second: 1},
	}

	checkPt := func(pt point) (bool, int) {
		if v, exists := d.visited[pt]; exists && d.grid[pt] == cellStatusOpen {
			return true, v.distanceFromStart
		}
		return false, math.MaxInt
	}

	d.shortestPath = []point{d.goalPos}

	for pos.First != 0 || pos.Second != 0 {
		lowestDist := math.MaxInt
		lowestPoint := point{}
		for _, pt := range checkOffsets {
			newPt := point{First: pos.First + pt.First, Second: pos.Second + pt.Second}
			if found, dist := checkPt(newPt); found && dist < lowestDist {
				lowestDist = dist
				lowestPoint = newPt
			}
		}

		d.shortestPath = append(d.shortestPath, lowestPoint)
		pos = lowestPoint
	}
}

func (d *Day15) exploreFullMap() map[point]*visitedStatus {
	grids := make([]map[point]cellStatus, 0, len(dirOrders))
	goalVisited := d.visited

	for _, dirOrder := range dirOrders {
		d.program.Reset()

		targetPos := point{}
		nextDir := dirFirst
		distFromStart := 0

		d.pos = point{}
		d.visited = map[point]*visitedStatus{
			{First: 0, Second: 0}: {timesVisited: 1},
		}
		d.grid = map[point]cellStatus{
			{First: 0, Second: 0}: cellStatusOpen,
		}
		d.program.RunIn(func(inputStep int) int64 {
			var err error
			nextDir, targetPos, err = d.getDirToNextCellType(d.pos, cellStatusUnknown, 0, dirOrder)
			if err != nil {
				// ensure we never try to go back into the trapped spot
				d.visited[d.pos].timesVisited = maxVisited + 1
				for x := 1; x <= maxVisited && err != nil; x++ {
					nextDir, targetPos, err = d.getDirToNextCellType(d.pos, cellStatusOpen, x, dirOrder)
				}
			}
			if err != nil {
				// d.Draw()
				// panic(err)
				d.program.Stop()
			}

			return int64(nextDir)
		}, func(val int64, state u.IntcodeProgramState) {
			rVal := responseType(val)

			p := d.getPointInDirection(d.pos, nextDir)
			shouldMove := true
			switch rVal {
			case responseWall:
				d.grid[p] = cellStatusWall
				shouldMove = false
			case responseSuccess:
				d.grid[p] = cellStatusOpen
			case responseFoundGoal:
				d.grid[p] = cellStatusGoal
			}

			if shouldMove {
				d.pos = targetPos
				if d.visited[d.pos] == nil {
					d.visited[d.pos] = &visitedStatus{}
					distFromStart++
				} else {
					distFromStart--
				}
				d.visited[d.pos].timesVisited++
				d.visited[d.pos].distanceFromStart = distFromStart
			}

			if rVal == responseFoundGoal {
				// d.Draw()
				d.goalPos = targetPos
				goalVisited = d.visited
			}
		})

		grids = append(grids, d.grid)
	}

	d.grid = map[point]cellStatus{
		{First: 0, Second: 0}: cellStatusOpen,
	}
	for _, grid := range grids {
		keys := u.MapKeys(grid)
		for _, key := range keys {
			d.grid[key] = grid[key]
		}
	}

	return goalVisited
}

func (d *Day15) tagDistanceRecursive(pos, last point, dist int, distances map[point]int) {
	distances[pos] = dist
	for _, cell := range d.getAdjacentCellsOfType(pos, cellStatusOpen) {
		if cell == last {
			continue
		}
		d.tagDistanceRecursive(cell, pos, dist+1, distances)
	}
}

func (d *Day15) Part1() string {
	d.visited = d.exploreFullMap()
	d.markShortestPath()

	for _, visited := range d.visited {
		visited.timesVisited = 1
	}
	d.pos = point{}
	// d.Draw()

	return fmt.Sprintf("Moves required to reach target: %s%d%s", u.TextBold, d.visited[d.goalPos].distanceFromStart, u.TextReset)
}

func (d *Day15) Part2() string {
	startLoc := d.goalPos
	distanceMap := map[point]int{startLoc: 0}

	d.tagDistanceRecursive(startLoc, point{}, 0, distanceMap)

	cellDistances := u.MapValues(distanceMap)
	sort.Slice(cellDistances, func(i, j int) bool { return cellDistances[i] > cellDistances[j] })

	return fmt.Sprintf("Time to fill the area with oxygen: %s%d%s minutes", u.TextBold, cellDistances[0], u.TextReset)
}
