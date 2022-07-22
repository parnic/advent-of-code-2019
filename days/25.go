package days

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

type Day25 struct {
	program u.IntcodeProgram
}

func (d *Day25) Parse() {
	d.program = u.LoadIntcodeProgram("25p")
}

func (d Day25) Num() int {
	return 25
}

func (d *Day25) Part1() string {
	lastCmdStrings := make([]string, 0)
	sb := strings.Builder{}

	inReader := bufio.NewReader(os.Stdin)
	d.program.SetDebugASCIIPrint(true)
	d.program.RunIn(func(inputStep int) int64 {
		lastCmdStrings = lastCmdStrings[:0]

		text, _ := inReader.ReadString('\n')
		d.program.FeedInputString(text[1:])
		return int64(text[0])
	}, func(val int64, state u.IntcodeProgramState) {
		if val == '\n' {
			str := sb.String()
			if len(str) > 0 {
				lastCmdStrings = append(lastCmdStrings, sb.String())
				sb.Reset()
			}
		} else {
			sb.WriteRune(rune(val))
		}
	})

	lastString := lastCmdStrings[len(lastCmdStrings)-1]
	var answer string
	if idx := strings.Index(lastString, " by typing "); idx >= 0 {
		startIdx := idx + len(" by typing ")
		endIdx := startIdx + strings.Index(lastString[startIdx:], " ")
		answer = lastString[startIdx:endIdx]
	}

	return fmt.Sprintf("Door passcode: %s%s%s", u.TextBold, answer, u.TextReset)
}

func (d *Day25) Part2() string {
	return "There is no part 2"
}
