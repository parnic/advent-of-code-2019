package days

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

type reaction struct {
	inputs map[string]int64
	output u.Pair[string, int]
}

type Day14 struct {
	reactions []reaction
}

func (d *Day14) Parse() {
	lines := u.GetStringLines("14p")
	d.reactions = make([]reaction, len(lines))
	for i, line := range lines {
		sides := strings.Split(line, " => ")
		inputs := strings.Split(sides[0], ", ")
		output := sides[1]

		outPair := strings.Split(output, " ")
		outAmt, _ := strconv.Atoi(outPair[0])
		d.reactions[i].output = u.Pair[string, int]{First: outPair[1], Second: outAmt}
		d.reactions[i].inputs = make(map[string]int64)
		for _, input := range inputs {
			pair := strings.Split(input, " ")
			d.reactions[i].inputs[pair[1]], _ = strconv.ParseInt(pair[0], 10, 64)
		}
	}
}

func (d Day14) getReactionProducing(chem string) *reaction {
	for _, reaction := range d.reactions {
		if reaction.output.First == chem {
			return &reaction
		}
	}

	return nil
}

func (d Day14) Num() int {
	return 14
}

func (d *Day14) getOreRequiredForFuel(qty int64) int64 {
	oreRequired := int64(0)
	needs := map[string]int64{
		"FUEL": qty,
	}
	excess := make(map[string]int64)

	getFromExcess := func(qty int64, chemical string) int64 {
		available := u.Min(excess[chemical], qty)
		excess[chemical] -= available
		return available
	}

	for len(needs) > 0 {
		keys := u.MapKeys(needs)
		producing := keys[0]
		qtyRequired := needs[producing]
		delete(needs, producing)

		fromExcess := getFromExcess(qtyRequired, producing)
		if fromExcess == qtyRequired {
			continue
		}
		qtyRequired -= fromExcess

		reaction := d.getReactionProducing(producing)

		qtyProduced := int64(reaction.output.Second)
		reactionsNeeded := int64(math.Ceil(float64(qtyRequired) / float64(qtyProduced)))

		excess[producing] = (qtyProduced * reactionsNeeded) - qtyRequired

		for reagent, inputQty := range reaction.inputs {
			qtyNeeded := inputQty * reactionsNeeded
			if reagent == "ORE" {
				oreRequired += qtyNeeded
			} else {
				needs[reagent] += qtyNeeded
			}
		}
	}

	return oreRequired
}

func (d *Day14) Part1() string {
	neededOre := d.getOreRequiredForFuel(1)
	return fmt.Sprintf("Minimum ore to produce 1 FUEL: %s%d%s", u.TextBold, neededOre, u.TextReset)
}

func (d *Day14) Part2() string {
	oreAvailable := int64(1000000000000)
	estimate := oreAvailable / d.getOreRequiredForFuel(1)
	lastSuccess := u.Bisect(estimate, estimate*2, 1, func(val int64) bool {
		oreConsumed := d.getOreRequiredForFuel(val)
		return oreConsumed < oreAvailable
	})

	return fmt.Sprintf("Maximum fuel we can make from 1 trillion ore: %s%d%s", u.TextBold, lastSuccess, u.TextReset)
}
