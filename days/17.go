package days

import (
	"fmt"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

type camViewCellType int
type botFacing int
type day17Grid [][]camViewCellType

const (
	cellTypeScaffold camViewCellType = iota
	cellTypeOpen
	cellTypeInvalid
)

const (
	botFacingUp botFacing = iota
	botFacingLeft
	botFacingDown
	botFacingRight

	botFacingFirst = botFacingUp
	botFacingLast  = botFacingRight
)

const (
	dirLeft                 = 1
	dirRight                = -1
	maxInstructionSetLength = 20
)

var (
	day17AdjacentOffsets = []u.Vec2i{
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: 0, Y: 1},
	}
)

type Day17 struct {
	program u.IntcodeProgram
}

func (d *Day17) Parse() {
	d.program = u.LoadIntcodeProgram("17p")
	// d.program.SetDebugASCIIPrint(true)
}

func (d Day17) Num() int {
	return 17
}

func (currentDir botFacing) getNewFacingDir(turnDir int) botFacing {
	currentDir += botFacing(turnDir)
	if currentDir < botFacingFirst {
		currentDir = botFacingLast
	} else if currentDir > botFacingLast {
		currentDir = botFacingFirst
	}

	return currentDir
}

func (grid day17Grid) Draw(botLocation u.Vec2i, botFacingDir botFacing, endLocation u.Vec2i) {
	for y := range grid {
		for x := range grid[y] {
			switch grid[y][x] {
			case cellTypeOpen:
				fmt.Print(" ")
			case cellTypeScaffold:
				char := "█"
				color := u.ColorBlack
				if botLocation.X == x && botLocation.Y == y {
					switch botFacingDir {
					case botFacingUp:
						char = "^"
					case botFacingLeft:
						char = "<"
					case botFacingDown:
						char = "v"
					case botFacingRight:
						char = ">"
					}
				} else if endLocation.X == x && endLocation.Y == y {
					char = "@"
				} else {
					color = u.ColorWhite
				}
				fmt.Printf("%s%s%s%s", u.BackgroundWhite, color, char, u.TextReset)
			}
		}
		fmt.Println()
	}
}

func (grid day17Grid) getAdjacentScaffolds(y, x int) []u.Vec2i {
	retval := make([]u.Vec2i, 0)
	for _, offset := range day17AdjacentOffsets {
		offY := y + offset.Y
		offX := x + offset.X
		if offY < 0 || offY >= len(grid) ||
			offX < 0 || offX >= len(grid[0]) {
			continue
		}

		if grid[offY][offX] == cellTypeScaffold {
			retval = append(retval, u.Vec2i{X: offX, Y: offY})
		}
	}

	return retval
}

func (grid day17Grid) forEachCellOfType(t camViewCellType, f func(y, x int)) {
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == t {
				f(y, x)
			}
		}
	}
}

func (grid *day17Grid) processGridUpdate(y int, rVal rune, currBotLocation u.Vec2i, currBotFacing botFacing) (int, u.Vec2i, botFacing) {
	grid.appendValue(rVal, y)

	switch rVal {
	case '\n':
		y++
	case '^', '<', 'v', '>':
		currBotLocation = u.Vec2i{X: len((*grid)[y]) - 1, Y: y}
		switch rVal {
		case '^':
			currBotFacing = botFacingUp
		case '<':
			currBotFacing = botFacingLeft
		case 'v':
			currBotFacing = botFacingDown
		case '>':
			currBotFacing = botFacingRight
		}
	}

	return y, currBotLocation, currBotFacing
}

func (grid day17Grid) getCellTypeInDirection(y, x int, facingDir botFacing) (camViewCellType, int, int) {
	newX := x
	newY := y
	switch facingDir {
	case botFacingUp:
		newY--
	case botFacingLeft:
		newX--
	case botFacingDown:
		newY++
	case botFacingRight:
		newX++
	}

	if newY < 0 || newY >= len(grid) || newX < 0 || newX >= len(grid[0]) {
		return cellTypeInvalid, newY, newX
	}

	return grid[newY][newX], newY, newX
}

func (grid *day17Grid) appendValue(rVal rune, row int) {
	ensureCapacity := func(y int) {
		for len(*grid) <= y {
			*grid = append(*grid, make([]camViewCellType, 0))
		}
	}

	switch rVal {
	case '#':
		ensureCapacity(row)
		(*grid)[row] = append((*grid)[row], cellTypeScaffold)
	case '.':
		ensureCapacity(row)
		(*grid)[row] = append((*grid)[row], cellTypeOpen)
	case '^', '<', 'v', '>':
		ensureCapacity(row)
		(*grid)[row] = append((*grid)[row], cellTypeScaffold)
	}
}

func (grid day17Grid) findEndLocation(botLocation u.Vec2i) u.Vec2i {
	var endLocation u.Vec2i
	grid.forEachCellOfType(cellTypeScaffold, func(y, x int) {
		if numSurrounding := len(grid.getAdjacentScaffolds(y, x)); numSurrounding == 1 {
			if botLocation.X != x || botLocation.Y != y {
				endLocation = u.Vec2i{X: x, Y: y}
			}
		}
	})

	return endLocation
}

func (grid day17Grid) getTurnDirectionFromCorner(pos u.Vec2i, botFacingDir botFacing) (int, string) {
	adj := grid.getAdjacentScaffolds(pos.Y, pos.X)
	turnDirection := 0
	// this is so awful. i'm sure there's a better way, but i'm tired.
	if botFacingDir == botFacingUp || botFacingDir == botFacingDown {
		if u.ArrayContains(adj, u.Vec2i{X: pos.X - 1, Y: pos.Y}) {
			if botFacingDir == botFacingUp {
				turnDirection = dirLeft
			} else if botFacingDir == botFacingDown {
				turnDirection = dirRight
			}
		} else if u.ArrayContains(adj, u.Vec2i{X: pos.X + 1, Y: pos.Y}) {
			if botFacingDir == botFacingUp {
				turnDirection = dirRight
			} else if botFacingDir == botFacingDown {
				turnDirection = dirLeft
			}
		}
	} else {
		if u.ArrayContains(adj, u.Vec2i{X: pos.X, Y: pos.Y - 1}) {
			if botFacingDir == botFacingLeft {
				turnDirection = dirRight
			} else if botFacingDir == botFacingRight {
				turnDirection = dirLeft
			}
		} else if u.ArrayContains(adj, u.Vec2i{X: pos.X, Y: pos.Y + 1}) {
			if botFacingDir == botFacingLeft {
				turnDirection = dirLeft
			} else if botFacingDir == botFacingRight {
				turnDirection = dirRight
			}
		}
	}

	dirAscii := "L"
	if turnDirection == dirRight {
		dirAscii = "R"
	}

	return turnDirection, dirAscii
}

func buildInstructionString(instructions []string) string {
	workingInstructions := make([]string, len(instructions))
	copy(workingInstructions, instructions)

	minimumRecurrence := 3
	initialInstructionSubsetLen := 4

	instructionStr := strings.Join(workingInstructions, ",")
	progs := make([][]string, 3)
	for i := range progs {
		numFound := minimumRecurrence
		subLen := initialInstructionSubsetLen
		for numFound >= minimumRecurrence {
			numFound = 1
			instructionSubset := strings.Join(workingInstructions[0:subLen], ",")
			if len(instructionSubset) > maxInstructionSetLength {
				break
			}
			for x := len(instructionSubset); x <= len(instructionStr)-len(instructionSubset); x++ {
				if instructionStr[x:x+len(instructionSubset)] == instructionSubset {
					numFound++
					x += len(instructionSubset)
				}
			}
			if numFound >= minimumRecurrence {
				subLen += 2
			}
		}
		if numFound < minimumRecurrence {
			subLen -= 2
		}
		progs[i] = make([]string, subLen)
		copy(progs[i], workingInstructions[0:subLen])

		instructionStr = strings.ReplaceAll(instructionStr, strings.Join(progs[i], ","), "")
		instructionStr = strings.TrimPrefix(strings.ReplaceAll(instructionStr, ",,", ","), ",")

		if len(instructionStr) == 0 {
			workingInstructions = nil
		} else {
			workingInstructions = strings.Split(instructionStr, ",")
		}
	}

	if workingInstructions != nil {
		panic("failed to use up all instructions")
	}

	programStr := strings.Join(instructions, ",")
	for i := range progs {
		programStr = strings.ReplaceAll(programStr, strings.Join(progs[i], ","), fmt.Sprintf("%c", 'A'+i))
	}

	sb := strings.Builder{}
	sb.WriteString(programStr)
	sb.WriteRune('\n')

	for i := range progs {
		sb.WriteString(strings.Join(progs[i], ","))
		sb.WriteRune('\n')
	}

	runDebug := 'n'
	sb.WriteRune(runDebug)
	sb.WriteRune('\n')

	return sb.String()
}

func (grid day17Grid) solvePath(botLocation u.Vec2i, botFacingDir botFacing) string {
	instructions := make([]string, 0)

	pos := botLocation
	endLocation := grid.findEndLocation(botLocation)
	for {
		if pos == endLocation {
			break
		}

		turnDirection, dirAscii := grid.getTurnDirectionFromCorner(pos, botFacingDir)
		if turnDirection == 0 {
			panic("at an invalid location somehow")
		}

		instructions = append(instructions, dirAscii)

		botFacingDir = botFacingDir.getNewFacingDir(turnDirection)
		numMoved := 0
		for {
			cell, newY, newX := grid.getCellTypeInDirection(pos.Y, pos.X, botFacingDir)
			if cell != cellTypeScaffold {
				break
			}
			pos.X = newX
			pos.Y = newY
			numMoved++
		}
		instructions = append(instructions, fmt.Sprintf("%d", numMoved))
	}

	return buildInstructionString(instructions)
}

func (d *Day17) Part1() string {
	grid := day17Grid{}
	y := 0
	var botLocation u.Vec2i
	var botFacingDir botFacing

	d.program.RunIn(func(inputStep int) int64 {
		return 0
	}, func(val int64, state u.IntcodeProgramState) {
		rVal := rune(val)
		y, botLocation, botFacingDir = grid.processGridUpdate(y, rVal, botLocation, botFacingDir)
	})

	alignmentParameterTotal := 0
	grid.forEachCellOfType(cellTypeScaffold, func(y, x int) {
		if numSurrounding := len(grid.getAdjacentScaffolds(y, x)); numSurrounding == 4 {
			alignmentParameterTotal += y * x
		}
	})

	// endLocation := grid.findEndLocation(botLocation)
	// grid.Draw(botLocation, botFacingDir, endLocation)

	return fmt.Sprintf("Alignment parameter sum: %s%d%s", u.TextBold, alignmentParameterTotal, u.TextReset)
}

func (d *Day17) Part2() string {
	beforeGrid := day17Grid{}
	var beforeBotLocation u.Vec2i
	var beforeBotFacing botFacing

	afterGrid := day17Grid{}
	var afterBotLocation u.Vec2i
	var afterBotFacing botFacing

	d.program.Reset()
	d.program.SetMemory(0, 2)

	row := 0
	var outputState int
	var lastOutput int64
	d.program.RunIn(func(inputStep int) int64 {
		panic("unexpected read")
	}, func(val int64, state u.IntcodeProgramState) {
		rVal := rune(val)
		if outputState == 0 {
			row, beforeBotLocation, beforeBotFacing = beforeGrid.processGridUpdate(row, rVal, beforeBotLocation, beforeBotFacing)
		} else if outputState == 2 {
			row, afterBotLocation, afterBotFacing = afterGrid.processGridUpdate(row, rVal, afterBotLocation, afterBotFacing)
		}

		if rVal == '\n' && lastOutput == '\n' {
			if outputState == 0 {
				d.program.FeedInputString(beforeGrid.solvePath(beforeBotLocation, beforeBotFacing))
			}
			outputState++
			row = 0
		}

		lastOutput = val
	})

	// fmt.Println("initial grid:")
	// beforeEndLocation := beforeGrid.findEndLocation(beforeBotLocation)
	// beforeGrid.Draw(beforeBotLocation, beforeBotFacing, beforeEndLocation)

	// fmt.Println("completed grid:")
	// afterEndLocation := afterGrid.findEndLocation(afterBotLocation)
	// afterGrid.Draw(afterBotLocation, afterBotFacing, afterEndLocation)

	return fmt.Sprintf("Dust collected after traveling all paths: %s%d%s", u.TextBold, lastOutput, u.TextReset)
}
