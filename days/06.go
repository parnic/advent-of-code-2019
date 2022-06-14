package days

import (
	"fmt"
	"strings"

	"parnic.com/aoc2019/utilities"
)

type body struct {
	orbits *body
	obj    string
}

type Day06 struct {
	allBodies []*body
}

func (d *Day06) Parse() {
	d.allBodies = make([]*body, 0)

	getOrAddBody := func(obj string) *body {
		target := d.findBody(obj)
		if target == nil {
			target = &body{
				obj: obj,
			}
			d.allBodies = append(d.allBodies, target)
		}

		return target
	}

	lines := utilities.GetStringLines("06p")
	for _, line := range lines {
		bodies := strings.Split(line, ")")
		newBody := getOrAddBody(bodies[1])

		target := getOrAddBody(bodies[0])
		newBody.orbits = target
	}
}

func (d *Day06) findBody(obj string) *body {
	for _, checkBody := range d.allBodies {
		if checkBody.obj == obj {
			return checkBody
		}
	}

	return nil
}

func (d Day06) Num() int {
	return 6
}

func (d *Day06) Part1() string {
	orbits := 0
	for _, obj := range d.allBodies {
		next := obj.orbits
		for next != nil {
			next = next.orbits
			orbits++
		}
	}

	return fmt.Sprintf("Total orbits: %s%d%s", utilities.TextBold, orbits, utilities.TextReset)
}

func (d *Day06) Part2() string {
	you := d.findBody("YOU")
	san := d.findBody("SAN")

	youChildren := make([]*body, 0)
	next := you.orbits
	for next != nil {
		youChildren = append(youChildren, next)
		next = next.orbits
	}

	var linkingNode *body
	next = san.orbits
	for next != nil {
		if utilities.ArrayContains(youChildren, next) {
			linkingNode = next
			break
		}
		next = next.orbits
	}

	if linkingNode == nil {
		panic("")
	}

	getDistToLinking := func(start *body) int {
		dist := 0
		next = start.orbits
		for next != nil {
			if next == linkingNode {
				break
			}
			dist++
			next = next.orbits
		}

		return dist
	}

	distYouToLinking := getDistToLinking(you)
	distSanToLinking := getDistToLinking(san)

	return fmt.Sprintf("Transfers to get to Santa: %s%d%s", utilities.TextBold, distYouToLinking+distSanToLinking, utilities.TextReset)
}
