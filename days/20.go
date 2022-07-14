package days

import (
	"container/heap"
	"fmt"
	"math"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

type day20Cell int8
type day20Graph map[day20Portal][]u.Pair[day20Portal, int]

const (
	day20CellWall day20Cell = iota
	day20CellPath
	day20CellDonutHole
)

var (
	day20AdjacentOffsets = []u.Vec2i{
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: 0, Y: 1},
	}

	entrancePortal = day20Portal{name: "AA"}
	exitPortal     = day20Portal{name: "ZZ"}
)

type day20Portal struct {
	name  string
	inner bool
	depth int
}

type Day20 struct {
	grid        [][]day20Cell
	entrance    u.Vec2i
	exit        u.Vec2i
	portals     map[day20Portal]u.Vec2i
	portalNames []string
}

func (d *Day20) Parse() {
	d.portals = make(map[day20Portal]u.Vec2i)
	d.portalNames = make([]string, 0)

	lines := u.GetStringLines("20p")
	d.grid = make([][]day20Cell, len(lines)-4)
	currPortal := strings.Builder{}

	for row, line := range lines {
		y := row - 2
		isGridRow := row >= 2 && row < len(lines)-2

		if isGridRow {
			d.grid[y] = make([]day20Cell, len(line)-4)
		}

		for col, ch := range lines[row] {
			x := col - 2
			isGridCol := col >= 2 && col < len(line)-2

			if isGridRow && isGridCol {
				if ch == '#' {
					d.grid[y][x] = day20CellWall
				} else if ch == '.' {
					d.grid[y][x] = day20CellPath
				} else {
					d.grid[y][x] = day20CellDonutHole
				}
			}

			if ch >= 'A' && ch <= 'Z' {
				portalX := 0
				portalY := 0
				if len(line) > col+1 && line[col+1] >= 'A' && line[col+1] <= 'Z' {
					currPortal.WriteRune(ch)
					currPortal.WriteRune(rune(line[col+1]))

					if len(line) > col+2 && line[col+2] == '.' {
						portalY = y
						portalX = x + 2
					} else if col-1 >= 0 && line[col-1] == '.' {
						portalY = y
						portalX = x - 1
					} else {
						panic("!")
					}
				} else if len(lines) > row+1 && lines[row+1][col] >= 'A' && lines[row+1][col] <= 'Z' {
					currPortal.WriteRune(ch)
					currPortal.WriteRune(rune(lines[row+1][col]))

					if len(lines) > row+2 && lines[row+2][col] == '.' {
						portalY = y + 2
						portalX = x
					} else if row-1 >= 0 && lines[row-1][col] == '.' {
						portalY = y - 1
						portalX = x
					} else {
						panic("!")
					}
				}

				if currPortal.Len() == 2 {
					portalName := currPortal.String()
					portalVec := u.Vec2i{X: portalX, Y: portalY}

					if portalName == entrancePortal.name {
						d.entrance = portalVec
					} else if portalName == exitPortal.name {
						d.exit = portalVec
					} else {
						portal := day20Portal{
							name:  portalName,
							inner: !d.isOuterPortal(portalVec),
						}
						d.portals[portal] = portalVec

						u.AddToArray(&d.portalNames, portalName)
					}

					currPortal.Reset()
				}
			}
		}
	}
}

func (d Day20) Num() int {
	return 20
}

func (d Day20) isPortal(vec u.Vec2i) (bool, int) {
	if d.grid[vec.Y][vec.X] != day20CellPath {
		return false, 0
	}

	for i, name := range d.portalNames {
		p, exists := d.portals[day20Portal{name: name, inner: !d.isOuterPortal(vec)}]
		if exists && vec == p {
			return true, i
		}
	}

	return false, 0
}

func (d Day20) Draw() {
	for y := range d.grid {
		for x := range d.grid[y] {
			switch d.grid[y][x] {
			case day20CellWall:
				fmt.Print("█")
			case day20CellDonutHole:
				fmt.Print(" ")
			case day20CellPath:
				posVec := u.Vec2i{X: x, Y: y}
				if posVec == d.entrance {
					fmt.Print("@")
				} else if posVec == d.exit {
					fmt.Print("!")
				} else {
					isPortal, portalIdx := d.isPortal(posVec)
					if isPortal {
						ch := 'a' + portalIdx
						if ch > 'z' {
							ch = 'A' + (portalIdx - 26)
						}
						fmt.Printf("%c", ch)
					} else {
						fmt.Print(" ")
					}
				}
			}
		}

		fmt.Println()
	}
}

func (d Day20) isOuterPortal(pos u.Vec2i) bool {
	return pos.X == 0 || pos.Y == 0 || pos.X == len(d.grid[0])-1 || pos.Y == len(d.grid)-1
}

func (d Day20) findAdjacentCells(inPos u.Vec2i) []u.Pair[day20Portal, int] {
	found := make([]u.Pair[day20Portal, int], 0)

	getAdjacent := func(pos u.Vec2i) []u.Vec2i {
		retAdjacent := make([]u.Vec2i, 0, len(day20AdjacentOffsets))
		for _, off := range day20AdjacentOffsets {
			offVec := u.Vec2i{X: pos.X + off.X, Y: pos.Y + off.Y}
			if offVec.Y < 0 || offVec.Y >= len(d.grid) || offVec.X < 0 || offVec.X >= len(d.grid[0]) {
				continue
			}
			if d.grid[offVec.Y][offVec.X] != day20CellPath {
				continue
			}
			retAdjacent = append(retAdjacent, offVec)
		}

		return retAdjacent
	}

	queue := make([]u.Pair[int, u.Vec2i], 0)
	visited := map[u.Vec2i]bool{
		inPos: true,
	}
	for _, adjacent := range getAdjacent(inPos) {
		queue = append(queue, u.Pair[int, u.Vec2i]{First: 1, Second: adjacent})
	}

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]

		if _, exists := visited[next.Second]; !exists {
			visited[next.Second] = true

			adjacentIsPortal, portalIdx := d.isPortal(next.Second)
			if adjacentIsPortal || next.Second == d.entrance || next.Second == d.exit {
				var portalName string
				if next.Second == d.entrance {
					portalName = entrancePortal.name
				} else if next.Second == d.exit {
					portalName = exitPortal.name
				} else {
					portalName = d.portalNames[portalIdx]
				}

				alreadyFound := false
				for _, p := range found {
					if p.First.name == portalName {
						alreadyFound = true
						break
					}
				}
				if !alreadyFound {
					found = append(found, u.Pair[day20Portal, int]{First: day20Portal{
						name:  portalName,
						inner: !d.isOuterPortal(next.Second),
					}, Second: next.First})
					continue
				}
			}

			for _, neighbor := range getAdjacent(next.Second) {
				if _, exists := visited[neighbor]; !exists {
					queue = append(queue, u.Pair[int, u.Vec2i]{First: next.First + 1, Second: neighbor})
				}
			}
		}
	}

	return found
}

type day20PriorityQueue struct {
	distance int
	neighbor day20Portal
}
type day20PriorityQueueHeap []day20PriorityQueue

func (h day20PriorityQueueHeap) Len() int           { return len(h) }
func (h day20PriorityQueueHeap) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h day20PriorityQueueHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *day20PriorityQueueHeap) Push(x any) {
	*h = append(*h, x.(day20PriorityQueue))
}

func (h *day20PriorityQueueHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (d Day20) dijkstra(graph day20Graph, start, end day20Portal, neighborFunc func(inPortal day20Portal) []u.Pair[day20Portal, int]) int {
	distance := make(map[day20Portal]int)

	ih := day20PriorityQueueHeap{
		day20PriorityQueue{
			distance: 0,
			neighbor: start,
		},
	}

	heap.Init(&ih)
	visited := make(map[day20Portal]bool)

	for ih.Len() > 0 {
		node := heap.Pop(&ih).(day20PriorityQueue)

		if node.neighbor == end {
			return node.distance
		}

		if _, exists := visited[node.neighbor]; exists {
			continue
		}

		visited[node.neighbor] = true

		for _, p := range neighborFunc(node.neighbor) {
			if _, exists := visited[p.First]; exists {
				continue
			}

			newDistance := node.distance + p.Second
			if dist, exists := distance[p.First]; !exists || newDistance < dist {
				distance[p.First] = newDistance
				heap.Push(&ih, day20PriorityQueue{
					distance: newDistance,
					neighbor: p.First,
				})
			}
		}
	}

	return math.MaxInt
}

func (d Day20) buildGraph() day20Graph {
	graph := make(day20Graph, len(d.portals)+1)

	adjacent := d.findAdjacentCells(d.entrance)
	graph[entrancePortal] = adjacent

	for portal, portalVec := range d.portals {
		adjacent = d.findAdjacentCells(portalVec)
		graph[portal] = adjacent
	}

	return graph
}

func (d Day20) getDepthNeighbors(graph day20Graph, portal day20Portal) []u.Pair[day20Portal, int] {
	basePortal := portal
	basePortal.depth = 0
	baseNeighbors := graph[basePortal]

	neighbors := make([]u.Pair[day20Portal, int], 0)

	if portal.inner {
		n := day20Portal{name: portal.name, inner: false, depth: portal.depth + 1}
		neighbors = append(neighbors, u.Pair[day20Portal, int]{First: n, Second: 1})
	}

	if portal.depth == 0 {
		for _, i := range baseNeighbors {
			if i.First.inner || i.First.name == entrancePortal.name || i.First.name == exitPortal.name {
				neighbors = append(neighbors, u.Pair[day20Portal, int]{First: i.First, Second: i.Second})
			}
		}
	} else {
		if !portal.inner {
			n := day20Portal{name: portal.name, inner: true, depth: portal.depth - 1}
			neighbors = append(neighbors, u.Pair[day20Portal, int]{First: n, Second: 1})
		}

		for _, i := range baseNeighbors {
			if i.First.name != entrancePortal.name && i.First.name != exitPortal.name {
				n := day20Portal{name: i.First.name, inner: i.First.inner, depth: portal.depth}
				neighbors = append(neighbors, u.Pair[day20Portal, int]{First: n, Second: i.Second})
			}
		}
	}

	return neighbors
}

func (d *Day20) Part1() string {
	// d.Draw()

	graph := d.buildGraph()

	for portal, adjacent := range graph {
		if portal.name == entrancePortal.name {
			continue
		}

		graph[portal] = append(adjacent, u.Pair[day20Portal, int]{First: day20Portal{name: portal.name, inner: !portal.inner}, Second: 1})
	}

	distance := d.dijkstra(graph, entrancePortal, exitPortal, func(inPortal day20Portal) []u.Pair[day20Portal, int] { return graph[inPortal] })
	return fmt.Sprintf("Steps to traverse maze: %s%d%s", u.TextBold, distance, u.TextReset)
}

func (d *Day20) Part2() string {
	graph := d.buildGraph()
	distance := d.dijkstra(graph, entrancePortal, exitPortal, func(inPortal day20Portal) []u.Pair[day20Portal, int] { return d.getDepthNeighbors(graph, inPortal) })
	return fmt.Sprintf("Steps to traverse recursive maze: %s%d%s", u.TextBold, distance, u.TextReset)
}
