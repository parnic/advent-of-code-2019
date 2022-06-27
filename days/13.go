package days

import (
	"fmt"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

const (
	tileEmpty   = 0
	tileWall    = 1
	tileBlock   = 2
	tileHPaddle = 3
	tileBall    = 4
)

type tile struct {
	pos u.Vec2[int]
	id  int
}

type Day13 struct {
	program   u.IntcodeProgram
	tiles     []tile
	gameBoard [24][45]int
}

func (d *Day13) Parse() {
	d.program = u.LoadIntcodeProgram("13p")
	d.tiles = make([]tile, 0, 1080)
}

func (d Day13) Num() int {
	return 13
}

func (d Day13) getNumBlocks() int {
	blockTiles := 0
	for _, tile := range d.tiles {
		if tile.id == tileBlock {
			blockTiles++
		}
	}

	return blockTiles
}

func (d Day13) Draw() {
	s := strings.Builder{}
	for x := range d.gameBoard {
		for y := range d.gameBoard[x] {
			block := d.gameBoard[x][y]
			if block == tileBlock {
				s.WriteString(u.ColorBlue)
				s.WriteRune('█')
				s.WriteString(u.TextReset)
			} else if block == tileBall {
				s.WriteString(u.ColorGreen)
				s.WriteRune('█')
				s.WriteString(u.TextReset)
			} else if block == tileWall {
				s.WriteString(u.ColorWhite)
				s.WriteRune('█')
				s.WriteString(u.TextReset)
			} else if block == tileHPaddle {
				s.WriteString(u.ColorRed)
				s.WriteRune('█')
				s.WriteString(u.TextReset)
			} else if block == tileEmpty {
				s.WriteRune(' ')
			}
		}
		s.WriteRune('\n')
	}

	fmt.Print(s.String())
}

func (d *Day13) Part1() string {
	outputStep := 0
	var newTilePos u.Vec2[int]
	d.program.RunIn(func(inputStep int) int64 {
		return 0
	}, func(val int64, state u.IntcodeProgramState) {
		if outputStep == 0 {
			newTilePos.X = int(val)
			outputStep++
		} else if outputStep == 1 {
			newTilePos.Y = int(val)
			outputStep++
		} else {
			d.tiles = append(d.tiles, tile{
				pos: newTilePos,
				id:  int(val),
			})
			outputStep = 0
		}
	})

	return fmt.Sprintf("%d total tiles, # block tiles: %s%d%s", len(d.tiles), u.TextBold, d.getNumBlocks(), u.TextReset)
}

func (d *Day13) Part2() string {
	d.program.Reset()
	d.program.SetMemory(0, 2)

	outputStep := 0
	newTilePos := u.Vec2[int]{}
	var ball u.Vec2[int]
	var paddle u.Vec2[int]
	var score int64
	d.program.RunIn(func(inputStep int) int64 {
		if ball.X < paddle.X {
			return -1
		} else if ball.X > paddle.X {
			return 1
		}

		return 0
	}, func(val int64, state u.IntcodeProgramState) {
		if outputStep == 0 {
			newTilePos.X = int(val)
			outputStep++
		} else if outputStep == 1 {
			newTilePos.Y = int(val)
			outputStep++
		} else {
			if newTilePos.Equals(u.Vec2[int]{X: -1, Y: 0}) {
				score = val
			} else {
				d.gameBoard[newTilePos.Y][newTilePos.X] = int(val)

				if val == tileBall {
					ball = newTilePos
				} else if val == tileHPaddle {
					paddle = newTilePos
					// d.Draw()
					// time.Sleep(time.Millisecond * 33)
				}
			}

			outputStep = 0
		}
	})

	return fmt.Sprintf("Game over! Score: %s%d%s", u.TextBold, score, u.TextReset)
}
