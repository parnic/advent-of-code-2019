package days

import (
	"fmt"
	"math"

	u "parnic.com/aoc2019/utilities"
)

type Day16 struct {
	numberSet []int8
}

func (d *Day16) Parse() {
	numberSequence := u.GetStringContents("16p")
	d.numberSet = make([]int8, len(numberSequence))
	for i, numRune := range numberSequence {
		d.numberSet[i] = int8(numRune - '0')
	}
}

func (d Day16) Num() int {
	return 16
}

func (d *Day16) Part1() string {
	transformed := make([]int8, len(d.numberSet))
	copy(transformed, d.numberSet)

	transformPattern := []int8{0, 1, 0, -1}

	phases := 100
	workingSet := make([]int8, len(transformed))
	for i := 0; i < phases; i++ {
		copy(workingSet, transformed)

		// fmt.Printf("Phase %d. Input signal: %v\n", (i + 1), transformed)
		for destIdx := range transformed {
			repeated := 0
			patternIdx := 0
			workingVal := int64(0)
			for idx := range transformed {
				if repeated >= destIdx {
					repeated = 0
					patternIdx++
					if patternIdx == len(transformPattern) {
						patternIdx = 0
					}
				} else {
					repeated++
				}

				// fmt.Printf("%d*%d", transformed[idx], transformPattern[patternIdx])
				// if idx < len(transformed)-1 {
				// 	fmt.Print(" + ")
				// }
				workingVal += int64(transformed[idx] * transformPattern[patternIdx])
			}

			workingSet[destIdx] = int8(int64(math.Abs(float64(workingVal))) % 10)
			// fmt.Printf(" = %d\n", workingSet[destIdx])
		}

		copy(transformed, workingSet)
	}

	finalVal := 0
	for i := range transformed[0:8] {
		finalVal += int(transformed[i]) * int(math.Pow10(8-1-i))
	}
	return fmt.Sprintf("First 8 digits of the final output list: %s%d%s", u.TextBold, finalVal, u.TextReset)
}

func (d *Day16) Part2() string {
	transformed := make([]int8, len(d.numberSet)*10000)
	for i := 0; i < 10000; i++ {
		copy(transformed[i*len(d.numberSet):(i*len(d.numberSet))+len(d.numberSet)], d.numberSet)
	}

	finalMsgOffset := 0
	for i := 0; i < 7; i++ {
		finalMsgOffset += int(d.numberSet[i]) * int(math.Pow10(7-1-i))
	}

	if finalMsgOffset < len(transformed)/2 {
		panic("offset must be in the back half of the message for this solution to work")
	}

	phases := 100
	for p := 0; p < phases; p++ {
		rollingTotal := int8(0)
		for i := len(transformed) - 1; i >= finalMsgOffset; i-- {
			rollingTotal += transformed[i]
			rollingTotal = rollingTotal % 10
			transformed[i] = rollingTotal
		}
	}

	finalVal := 0
	for i := range transformed[finalMsgOffset : finalMsgOffset+8] {
		finalVal += int(transformed[finalMsgOffset+i]) * int(math.Pow10(8-1-i))
	}
	return fmt.Sprintf("Embedded message in the final output list: %s%d%s", u.TextBold, finalVal, u.TextReset)
}
