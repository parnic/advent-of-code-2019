package days

import (
	"fmt"
	"strconv"
	"strings"

	"parnic.com/aoc2019/utilities"
)

type Day04 struct {
	max int
	min int
}

func (d *Day04) Parse() {
	contents := utilities.GetStringContents("04p")
	vals := strings.Split(contents, "-")
	d.min, _ = strconv.Atoi(vals[0])
	d.max, _ = strconv.Atoi(vals[1])

	d.test()
}

func (d Day04) Num() int {
	return 4
}

func (d Day04) test() {
	if !d.isValidP2(112233) {
		panic("112233")
	}
	if d.isValidP2(123444) {
		panic("123444")
	}
	if !d.isValidP2(111122) {
		panic("111122")
	}
	if !d.isValidP2(112222) {
		panic("112222")
	}
}

func (d Day04) isValidP1(num int) bool {
	// if num < d.min || num > d.max {
	// 	return false
	// }

	numStr := strconv.Itoa(num)
	lastNum := -1
	for _, ch := range numStr {
		thisNum, _ := strconv.Atoi(string(ch))
		if lastNum > -1 && thisNum < lastNum {
			return false
		}
		lastNum = thisNum
	}

	foundDouble := false
	for idx := range numStr {
		if idx == 0 {
			continue
		}
		if numStr[idx-1] == numStr[idx] {
			foundDouble = true
			break
		}
	}

	return foundDouble
}

func (d Day04) isValidP2(num int) bool {
	if !d.isValidP1(num) {
		return false
	}

	numStr := strconv.Itoa(num)
	for i := 0; i < len(numStr); i++ {
		consec := 1
		for j := i + 1; j < len(numStr) && numStr[j] == numStr[i]; j++ {
			consec++
			i = j
		}
		if consec == 2 {
			return true
		}
	}

	return false
}

// these could be sped up quite a bit by intelligently jumping to the
// next valid number when encountering an invalid one, rather than
// simply incrementing by 1
func (d *Day04) Part1() string {
	numValid := 0
	for x := d.min; x <= d.max; x++ {
		if d.isValidP1(x) {
			numValid++
		}
	}

	return fmt.Sprintf("Total valid passwords: %s%d%s", utilities.TextBold, numValid, utilities.TextReset)
}

func (d *Day04) Part2() string {
	numValid := 0
	for x := d.min; x <= d.max; x++ {
		if d.isValidP2(x) {
			numValid++
		}
	}

	return fmt.Sprintf("Total valid passwords: %s%d%s", utilities.TextBold, numValid, utilities.TextReset)
}
