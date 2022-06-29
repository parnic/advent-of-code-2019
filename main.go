package main

import (
	"flag"
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

var (
	flagPart1 = flag.Bool("part1", false, "whether to run part1 or not; if no flags are present, all parts are run")
	flagPart2 = flag.Bool("part2", false, "whether to run part2 or not; if no flags are present, all parts are run")
)

var dayMap = []day{
	&days.Day01{},
	&days.Day02{},
	&days.Day03{},
	&days.Day04{},
	&days.Day05{},
	&days.Day06{},
	&days.Day07{},
	&days.Day08{},
	&days.Day09{},
	&days.Day10{},
	&days.Day11{},
	&days.Day12{},
	&days.Day13{},
	&days.Day14{},
	&days.Day15{},
}

func main() {
	flag.Parse()

	arg := strconv.Itoa(len(dayMap))
	flagArgs := flag.Args()
	if len(flagArgs) > 0 && len(flagArgs[0]) > 0 {
		arg = flagArgs[0]
	}
	if strings.ToLower(arg) == "all" {
		startTime := time.Now()
		for _, v := range dayMap {
			solve(v)
		}
		fmt.Printf("%sAll days completed in %v%s\n", utilities.ColorBrightBlack, time.Since(startTime), utilities.TextReset)
	} else {
		iArg, err := strconv.Atoi(arg)
		if err != nil {
			log.Fatalf("Invalid day " + utilities.ColorCyan + arg + utilities.TextReset)
		}

		if iArg < 0 || iArg > len(dayMap) {
			log.Fatalf("Unknown day " + utilities.ColorCyan + arg + utilities.TextReset)
		}

		solve(dayMap[iArg-1])
	}

	os.Exit(0)
}

func solve(d day) {
	fmt.Printf("%sDay %d%s\n", utilities.ColorCyan, d.Num(), utilities.TextReset)
	fmt.Printf("----%s\n", strings.Repeat("-", len(strconv.Itoa(d.Num()))))

	runPart1 := (!*flagPart1 && !*flagPart2) || *flagPart1
	runPart2 := (!*flagPart1 && !*flagPart2) || *flagPart2

	parseStart := time.Now()
	d.Parse()
	parseTime := time.Since(parseStart)

	part1Start := time.Now()
	var part1Text string
	if runPart1 {
		part1Text = d.Part1()
	}
	part1Time := time.Since(part1Start)
	if runPart1 {
		fmt.Println(part1Header)
		fmt.Println(">", part1Text)
		fmt.Println()
	}

	part2Start := time.Now()
	var part2Text string
	if runPart2 {
		part2Text = d.Part2()
	}
	part2Time := time.Since(part2Start)
	if runPart2 {
		fmt.Println(part2Header)
		fmt.Println(">", part2Text)
		fmt.Println()
	}

	fmt.Print(utilities.ColorBrightBlack)
	fmt.Println("Parsed in", parseTime)
	if runPart1 {
		fmt.Println("Part01 in", part1Time)
	}
	if runPart2 {
		fmt.Println("Part02 in", part2Time)
	}
	fmt.Println(utilities.TextReset)
}
