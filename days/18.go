package days

import (
	"container/heap"
	"fmt"
	"math"
	"strings"

	u "parnic.com/aoc2019/utilities"
)

type day18Cell int
type day18Vec u.Vec2[int]
type day18Graph map[rune][]u.Pair[rune, int]

const (
	day18CellWall day18Cell = iota
	day18CellOpen
)

var (
	day18AdjacentOffsets = []day18Vec{
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: 0, Y: 1},
	}
)

type reachableKeysMemo struct {
	pos       rune
	keysFound int
}

type minStepsMemo struct {
	pos        string
	keysToFind int
	keysFound  int
}

type Day18 struct {
	entrance           day18Vec
	grid               [][]day18Cell
	doors              map[day18Vec]int
	keys               map[day18Vec]int
	knownReachableKeys map[reachableKeysMemo][]u.Pair[rune, int]
	knownMinimumSteps  map[minStepsMemo]int
}

func (d *Day18) Parse() {
	d.doors = make(map[day18Vec]int)
	d.keys = make(map[day18Vec]int)
	d.knownReachableKeys = make(map[reachableKeysMemo][]u.Pair[rune, int])
	d.knownMinimumSteps = make(map[minStepsMemo]int, 0)

	lines := u.GetStringLines("18p")
	d.grid = make([][]day18Cell, len(lines))
	for i, line := range lines {
		d.grid[i] = make([]day18Cell, len(line))
		for j, char := range line {
			if char == '#' {
				d.grid[i][j] = day18CellWall
			} else if char == '.' {
				d.grid[i][j] = day18CellOpen
			} else if char == '@' {
				d.grid[i][j] = day18CellOpen
				d.entrance = day18Vec{X: j, Y: i}
			} else if char >= 'A' && char <= 'Z' {
				d.grid[i][j] = day18CellOpen
				d.doors[day18Vec{X: j, Y: i}] = int(char - 'A')
			} else if char >= 'a' && char <= 'z' {
				d.grid[i][j] = day18CellOpen
				d.keys[day18Vec{X: j, Y: i}] = int(char - 'a')
			}
		}
	}
}

func (d Day18) Num() int {
	return 18
}

func (d Day18) Draw(grid [][]day18Cell, keys, doors map[day18Vec]int, entrances ...day18Vec) {
	for y := range grid {
		for x := range grid[y] {
			switch grid[y][x] {
			case day18CellWall:
				fmt.Print("█")
			case day18CellOpen:
				posVec := day18Vec{X: x, Y: y}
				if _, exists := doors[posVec]; exists {
					fmt.Printf("%c", rune(doors[posVec]+'A'))
				} else if _, exists := keys[posVec]; exists {
					fmt.Printf("%c", rune(keys[posVec]+'a'))
				} else if u.ArrayContains(entrances, posVec) {
					fmt.Print("@")
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
}

func (d Day18) findAdjacentCells(inPos day18Vec, keys, doors map[day18Vec]int, grid [][]day18Cell) []u.Pair[rune, int] {
	found := make([]u.Pair[rune, int], 0)

	getAdjacent := func(pos day18Vec) []day18Vec {
		retAdjacent := make([]day18Vec, 0, len(day18AdjacentOffsets))
		for _, off := range day18AdjacentOffsets {
			offVec := day18Vec{X: pos.X + off.X, Y: pos.Y + off.Y}
			if grid[offVec.Y][offVec.X] == day18CellWall {
				continue
			}
			retAdjacent = append(retAdjacent, offVec)
		}

		return retAdjacent
	}

	queue := make([]u.Pair[int, day18Vec], 0)
	visited := make(map[day18Vec]bool)
	for _, adjacent := range getAdjacent(inPos) {
		queue = append(queue, u.Pair[int, day18Vec]{First: 1, Second: adjacent})
	}

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]

		if _, exists := visited[next.Second]; !exists {
			visited[next.Second] = true

			key, adjacentIsKey := keys[next.Second]
			door, adjacentIsDoor := doors[next.Second]
			if adjacentIsKey || adjacentIsDoor {
				var rVal rune
				if adjacentIsKey {
					rVal = rune('a' + key)
				} else if adjacentIsDoor {
					rVal = rune('A' + door)
				}

				alreadyFound := false
				for _, p := range found {
					if p.First == rVal {
						alreadyFound = true
						break
					}
				}
				if !alreadyFound {
					found = append(found, u.Pair[rune, int]{First: rVal, Second: next.First})
					continue
				}
			}

			for _, neighbor := range getAdjacent(next.Second) {
				if _, exists := visited[neighbor]; !exists {
					queue = append(queue, u.Pair[int, day18Vec]{First: next.First + 1, Second: neighbor})
				}
			}
		}
	}

	return found
}

type day18PriorityQueue struct {
	distance int
	neighbor rune
}
type day18PriorityQueueHeap []day18PriorityQueue

func (h day18PriorityQueueHeap) Len() int           { return len(h) }
func (h day18PriorityQueueHeap) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h day18PriorityQueueHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *day18PriorityQueueHeap) Push(x any) {
	*h = append(*h, x.(day18PriorityQueue))
}

func (h *day18PriorityQueueHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (d Day18) reachableKeys(inPos rune, keysFound int, graph day18Graph) []u.Pair[rune, int] {
	memo := reachableKeysMemo{
		pos:       inPos,
		keysFound: keysFound,
	}
	if v, exists := d.knownReachableKeys[memo]; exists {
		return v
	}

	ret := make([]u.Pair[rune, int], 0)
	distance := make(map[rune]int)

	ih := make(day18PriorityQueueHeap, 0)

	for _, p := range graph[inPos] {
		ih = append(ih, day18PriorityQueue{
			distance: p.Second,
			neighbor: p.First,
		})
	}

	heap.Init(&ih)

	for ih.Len() > 0 {
		node := heap.Pop(&ih).(day18PriorityQueue)

		// it's a key and we haven't picked it up yet...
		if node.neighbor >= 'a' && node.neighbor <= 'z' && (1<<int(node.neighbor-'a')&keysFound) == 0 {
			ret = append(ret, u.Pair[rune, int]{First: node.neighbor, Second: node.distance})
			continue
		}

		// it's a door but we don't have the key yet...
		if node.neighbor >= 'A' && node.neighbor <= 'Z' && ((1<<int(node.neighbor-'A'))&keysFound) == 0 {
			continue
		}

		for _, p := range graph[node.neighbor] {
			newDistance := node.distance + p.Second
			if dist, exists := distance[p.First]; !exists || newDistance < dist {
				distance[p.First] = newDistance
				heap.Push(&ih, day18PriorityQueue{
					distance: newDistance,
					neighbor: p.First,
				})
			}
		}
	}

	d.knownReachableKeys[memo] = ret
	return ret
}

func (d Day18) minimumSteps(inPos string, keysToFind int, keysFound int, graph day18Graph) int {
	memo := minStepsMemo{
		pos:        inPos,
		keysToFind: keysToFind,
		keysFound:  keysFound,
	}
	if v, exists := d.knownMinimumSteps[memo]; exists {
		return v
	}

	if keysToFind == 0 {
		return 0
	}

	best := math.Inf(1)
	for _, item := range inPos {
		for _, p := range d.reachableKeys(item, keysFound, graph) {
			sb := strings.Builder{}
			oldIdx := strings.IndexRune(inPos, item)
			for i := range inPos {
				if i == oldIdx {
					sb.WriteRune(p.First)
				} else {
					sb.WriteByte(inPos[i])
				}
			}
			newKeys := keysFound + (1 << (p.First - 'a'))
			dist := p.Second

			dist += d.minimumSteps(sb.String(), keysToFind-1, newKeys, graph)

			if float64(dist) < best {
				best = float64(dist)
			}
		}
	}

	d.knownMinimumSteps[memo] = int(best)
	return int(best)
}

func (d Day18) buildGraph(pos []day18Vec, keys map[day18Vec]int, doors map[day18Vec]int, grid [][]day18Cell) day18Graph {
	graph := make(day18Graph)
	for i, p := range pos {
		adjacent := d.findAdjacentCells(p, keys, doors, grid)
		graph[rune('1'+i)] = adjacent
	}
	for keyPos, keyType := range keys {
		graph[rune('a'+keyType)] = d.findAdjacentCells(keyPos, keys, doors, grid)
	}
	for doorPos, doorType := range doors {
		graph[rune('A'+doorType)] = d.findAdjacentCells(doorPos, keys, doors, grid)
	}

	return graph
}

func (d Day18) part2PatchMap(grid [][]day18Cell, entrance day18Vec) []day18Vec {
	grid[entrance.Y-1][entrance.X] = day18CellWall
	grid[entrance.Y][entrance.X-1] = day18CellWall
	grid[entrance.Y][entrance.X] = day18CellWall
	grid[entrance.Y][entrance.X+1] = day18CellWall
	grid[entrance.Y+1][entrance.X] = day18CellWall

	return []day18Vec{
		{X: entrance.X - 1, Y: entrance.Y - 1},
		{X: entrance.X + 1, Y: entrance.Y - 1},
		{X: entrance.X - 1, Y: entrance.Y + 1},
		{X: entrance.X + 1, Y: entrance.Y + 1},
	}
}

func (d *Day18) Part1() string {
	// fmt.Println("initial state:")
	// d.Draw(d.grid, d.keys, d.doors, d.entrance)

	graph := d.buildGraph([]day18Vec{d.entrance}, d.keys, d.doors, d.grid)
	minSteps := d.minimumSteps("1", len(d.keys), 0, graph)

	return fmt.Sprintf("Total distance traveled: %s%d%s", u.TextBold, minSteps, u.TextReset)
}

func (d *Day18) Part2() string {
	// fmt.Println("initial state:")
	grid := make([][]day18Cell, len(d.grid))
	for i := range d.grid {
		grid[i] = make([]day18Cell, len(d.grid[i]))
		copy(grid[i], d.grid[i])
	}

	entrances := d.part2PatchMap(grid, d.entrance)
	// d.Draw(grid, d.keys, d.doors, entrances...)

	// clear memoized maps that (might have) came from part1
	d.knownMinimumSteps = make(map[minStepsMemo]int)
	d.knownReachableKeys = make(map[reachableKeysMemo][]u.Pair[rune, int])

	graph := d.buildGraph(entrances, d.keys, d.doors, grid)
	minSteps := d.minimumSteps("1234", len(d.keys), 0, graph)

	return fmt.Sprintf("Total distance traveled: %s%d%s", u.TextBold, minSteps, u.TextReset)
}
