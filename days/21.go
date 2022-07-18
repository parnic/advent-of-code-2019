package days

import (
	"fmt"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

type Day21 struct {
	program u.IntcodeProgram
}

func (d *Day21) Parse() {
	d.program = u.LoadIntcodeProgram("21p")
	// d.program.SetDebugASCIIPrint(true)
}

func (d Day21) Num() int {
	return 21
}

func (d *Day21) Part1() string {
	// if there's any hole up to 3 ahead of us but there's ground where we'd land if we jumped
	// (a jump takes 4 spaces), go ahead and jump
	cmds := []string{
		// check if a hole at 1 or 2 ahead
		"NOT A T",
		"NOT B J",
		// store that result in J
		"OR T J",
		// check if a hole at 3 ahead
		"NOT C T",
		// store hole in 1, 2, or 3 in T
		"OR J T",
		// set J true if hole in 1, 2, or 3
		"OR T J",
		// set J true if also no hole at 4 ahead
		"AND D J",
		"WALK",
	}
	instructionStr := strings.Join(cmds, "\n") + "\n"
	d.program.FeedInputString(instructionStr)

	res := d.program.Run()

	return fmt.Sprintf("Hull damage value when walking: %s%d%s", u.TextBold, res, u.TextReset)
}

func (d *Day21) Part2() string {
	d.program.Reset()
	//   @
	// #####.#.##.##.###
	//    ABCDEFGHI
	// using the first program, this kills us. if we jump, we land at D and H becomes our new D, so it won't jump again.
	// but if we waited to jump until we got one more ahead, we'd be ok.
	// so now we want to know essentially the same thing as part 1, but also if our multiple (immediate second jump) would be successful.
	// in problem terms, that's: if there's a hole at 1 or 2 ahead, and there's a hole at C with ground at H, and there's ground at D.
	// so now for the above example we'd wait to jump until here:
	//     @
	// #####.#.##.##.###
	//      ABCDEFGHI
	// and all will be well.
	cmds := []string{
		// check if a hole at 1 or 2 ahead
		"NOT A J",
		"NOT B T",
		// store that result in J
		"OR T J",
		// check if a hole at 3 ahead...
		"NOT C T",
		// and ground at 8 ahead (so we can immediately jump again if needed)...
		"AND H T",
		// combine those into J
		"OR T J",
		// and ensure we also still have a place to land if we jumped right away
		"AND D J",
		"RUN",
	}
	instructionStr := strings.Join(cmds, "\n") + "\n"
	d.program.FeedInputString(instructionStr)

	res := d.program.Run()

	return fmt.Sprintf("Hull damage value when running: %s%d%s", u.TextBold, res, u.TextReset)
}
