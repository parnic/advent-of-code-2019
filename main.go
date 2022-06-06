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
	Part1() string
	Part2() string
}

var dayMap = map[int]day{
	1: &days.Day01{},
}

func main() {
	arg := "1"
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
	fmt.Println(utilities.ColorCyan + "Day 01" + utilities.TextReset)

	parseStart := time.Now()
	d.Parse()
	parseTime := time.Since(parseStart)

	part1Start := time.Now()
	part1 := d.Part1()
	part1Time := time.Since(part1Start)

	part2Start := time.Now()
	part2 := d.Part2()
	part2Time := time.Since(part2Start)

	fmt.Println(utilities.ColorGreen + "Part1:" + utilities.TextReset)
	fmt.Println(part1)
	fmt.Println(utilities.ColorGreen + "Part2:" + utilities.TextReset)
	fmt.Println(part2)
	fmt.Print(utilities.ColorBrightBlack)
	fmt.Printf("Parsed in %s\n", parseTime)
	fmt.Printf("Part01 in %s\n", part1Time)
	fmt.Printf("Part02 in %s\n", part2Time)
	fmt.Println(utilities.TextReset)
}
