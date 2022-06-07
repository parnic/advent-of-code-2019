package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"parnic.com/aoc2019/days"
	"parnic.com/aoc2019/utilities"
)

type day interface {
	Parse()
	Num() int
	Part1() string
	Part2() string
}

const (
	part1Header = utilities.ColorGreen + "Part1:" + utilities.TextReset
	part2Header = utilities.ColorGreen + "Part2:" + utilities.TextReset
)

var dayMap = map[int]day{
	1: &days.Day01{},
}

func main() {
	arg := strconv.Itoa(len(dayMap))
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	if strings.ToLower(arg) == "all" {
		for _, v := range dayMap {
			solve(v)
		}
	} else {
		iArg, err := strconv.Atoi(arg)
		if err != nil {
			log.Fatalf("Invalid day " + utilities.ColorCyan + arg + utilities.TextReset)
		}

		p, ok := dayMap[iArg]
		if !ok {
			log.Fatalf("Unknown day " + utilities.ColorCyan + arg + utilities.TextReset)
		}

		solve(p)
	}

	os.Exit(0)
}

func solve(d day) {
	fmt.Println(fmt.Sprintf("%sDay %d%s", utilities.ColorCyan, d.Num(), utilities.TextReset))

	parseStart := time.Now()
	d.Parse()
	parseTime := time.Since(parseStart)

	part1Start := time.Now()
	part1Text := d.Part1()
	part1Time := time.Since(part1Start)

	part2Start := time.Now()
	part2Text := d.Part2()
	part2Time := time.Since(part2Start)

	fmt.Println(part1Header)
	fmt.Println(part1Text)
	fmt.Println(part2Header)
	fmt.Println(part2Text)
	fmt.Print(utilities.ColorBrightBlack)
	fmt.Println("Parsed in", parseTime)
	fmt.Println("Part01 in", part1Time)
	fmt.Println("Part02 in", part2Time)
	fmt.Println(utilities.TextReset)
}
