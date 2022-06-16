package days

import (
	"fmt"
	"math"
	"strings"

	"parnic.com/aoc2019/utilities"
)

type Day08 struct {
	img [][]int
}

const (
	imgWidth  = 25
	imgHeight = 6
)

func (d *Day08) Parse() {
	contents := utilities.GetStringContents("08p")
	imgSize := imgWidth * imgHeight
	layers := len(contents) / imgSize
	d.img = make([][]int, layers)
	for layer := 0; layer < layers; layer++ {
		d.img[layer] = make([]int, imgSize)
		for i := 0; i < imgSize; i++ {
			d.img[layer][i] = int(contents[(layer*imgSize)+i] - '0')
		}
	}
}

func (d Day08) Num() int {
	return 8
}

func (d *Day08) Part1() string {
	fewestZeroes := math.MaxInt
	var layerFewestZeroes int
	for layer := range d.img {
		zeroes := 0
		for i := range d.img[layer] {
			if d.img[layer][i] == 0 {
				zeroes++
			}
		}

		if zeroes < fewestZeroes {
			fewestZeroes = zeroes
			layerFewestZeroes = layer
		}
	}

	numOne := 0
	numTwo := 0
	for i := range d.img[layerFewestZeroes] {
		if d.img[layerFewestZeroes][i] == 1 {
			numOne++
		} else if d.img[layerFewestZeroes][i] == 2 {
			numTwo++
		}
	}

	return fmt.Sprintf("Fewest zeroes on layer %d, #1s * #2s = %d * %d = %s%d%s",
		layerFewestZeroes,
		numOne,
		numTwo,
		utilities.TextBold,
		numOne*numTwo,
		utilities.TextReset,
	)
}

func (d *Day08) Part2() string {
	imgSize := imgWidth * imgHeight
	finalImg := make([]int, imgSize)
	for i := 0; i < imgSize; i++ {
		for layer := 0; layer < len(d.img); layer++ {
			if d.img[layer][i] != 2 {
				finalImg[i] = d.img[layer][i]
				break
			}
		}
	}

	outStr := strings.Builder{}
	outStr.WriteString("Message received:\n")
	outStr.WriteString(utilities.TextBold)
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			if finalImg[(y*imgWidth)+x] == 0 {
				outStr.WriteRune(' ')
			} else {
				outStr.WriteRune('█')
			}
		}
		outStr.WriteRune('\n')
	}
	outStr.WriteString(utilities.TextReset)

	return outStr.String()
}
