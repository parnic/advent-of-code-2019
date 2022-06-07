package utilities

import (
	"bufio"
	"fmt"
	"strconv"

	"parnic.com/aoc2019/inputs"
)

func getData(filename string, lineHandler func(line string)) {
	file, err := inputs.Sets.Open(fmt.Sprintf("%s.txt", filename))
	// version that doesn't use embedded files:
	// file, err := os.Open(fmt.Sprintf("inputs/%s.txt", filename))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lineHandler(scanner.Text())
	}
}

func GetStringContents(filename string) string {
	var retval string
	getData(filename, func(line string) {
		if len(retval) != 0 {
			panic("tried to parse multi-line file as a single line")
		}
		retval = line
	})
	return retval
}

func GetStringLines(filename string) []string {
	retval := make([]string, 0)
	getData(filename, func(line string) {
		retval = append(retval, line)
	})
	return retval
}

func GetIntLines(filename string) []int64 {
	retval := make([]int64, 0)
	getData(filename, func(line string) {
		val, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			panic(err)
		}
		retval = append(retval, val)
	})
	return retval
}
