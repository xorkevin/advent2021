package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

type (
	Item struct {
		value Point
		g, f  int
		index int
	}

	PriorityQueue struct {
		q []*Item
		s map[Point]int
	}

	Point struct {
		x int
		y int
	}

	OpenSet struct {
		q *PriorityQueue
	}

	ClosedSet map[Point]struct{}
)

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		s: map[Point]int{},
	}
}

func (q PriorityQueue) Len() int { return len(q.q) }
func (q PriorityQueue) Less(i, j int) bool {
	return q.q[i].f < q.q[j].f
}
func (q PriorityQueue) Swap(i, j int) {
	q.q[i], q.q[j] = q.q[j], q.q[i]
	q.q[i].index = i
	q.q[j].index = j
	q.s[q.q[i].value] = i
	q.s[q.q[j].value] = j
}
func (q *PriorityQueue) Push(x interface{}) {
	n := len(q.q)
	item := x.(*Item)
	item.index = n
	q.q = append(q.q, item)
	q.s[item.value] = n
}
func (q *PriorityQueue) Pop() interface{} {
	n := len(q.q)
	item := q.q[n-1]
	q.q[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	q.q = q.q[:n-1]
	delete(q.s, item.value)
	return item
}
func (q *PriorityQueue) Update(value Point, g, f int) bool {
	idx, ok := q.s[value]
	if !ok {
		return false
	}
	item := q.q[idx]
	item.g = g
	item.f = f
	heap.Fix(q, item.index)
	return true
}

func NewOpenSet() *OpenSet {
	return &OpenSet{
		q: NewPriorityQueue(),
	}
}

func (s *OpenSet) Empty() bool {
	return s.q.Len() == 0
}

func (s *OpenSet) Has(val Point) bool {
	_, ok := s.q.s[val]
	return ok
}

func (s *OpenSet) Get(val Point) (*Item, bool) {
	idx, ok := s.q.s[val]
	if !ok {
		return nil, false
	}
	return s.q.q[idx], true
}

func (s *OpenSet) Push(value Point, g, f int) {
	heap.Push(s.q, &Item{
		value: value,
		g:     g,
		f:     f,
	})
}

func (s *OpenSet) Pop() (Point, int, int) {
	item := heap.Pop(s.q).(*Item)
	return item.value, item.g, item.f
}

func (s *OpenSet) Update(value Point, g, f int) bool {
	return s.q.Update(value, g, f)
}

func NewClosedSet() ClosedSet {
	return ClosedSet{}
}

func (cs ClosedSet) Has(val Point) bool {
	_, ok := cs[val]
	return ok
}

func (cs ClosedSet) Push(val Point) {
	cs[val] = struct{}{}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func manhattanDistance(a, b Point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func main() {
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var grid [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []byte(scanner.Text())
		row := make([]int, 0, len(line))
		for _, i := range line {
			row = append(row, int(i)-'0')
		}
		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", pathfind(grid, Point{0, 0}, Point{len(grid[0]) - 1, len(grid) - 1}))
	fmt.Println("Part 2:", pathfind2(grid, Point{0, 0}, Point{len(grid[0])*5 - 1, len(grid)*5 - 1}))
}

func pathfind2(grid [][]int, start, end Point) int {
	openSet := NewOpenSet()
	openSet.Push(start, 0, manhattanDistance(start, end))
	closedSet := NewClosedSet()
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		if cur == end {
			return curg
		}
		for _, i := range neighbors2(grid, cur) {
			if closedSet.Has(i) {
				continue
			}
			g := curg + getVal(grid, i)
			f := g + manhattanDistance(i, end)
			if v, ok := openSet.Get(i); ok {
				if g < v.g {
					openSet.Update(i, g, f)
				}
				continue
			}
			openSet.Push(i, g, f)
		}
	}
	return -1
}

func getVal(grid [][]int, p Point) int {
	sx := len(grid[0])
	sy := len(grid)
	return (grid[p.y%sy][p.x%sx]+p.x/sx+p.y/sy-1)%9 + 1
}

func neighbors2(grid [][]int, p Point) []Point {
	points := make([]Point, 0, 4)
	if k := (Point{x: p.x - 1, y: p.y}); inBounds2(grid, k) {
		points = append(points, k)
	}
	if k := (Point{x: p.x, y: p.y - 1}); inBounds2(grid, k) {
		points = append(points, k)
	}
	if k := (Point{x: p.x + 1, y: p.y}); inBounds2(grid, k) {
		points = append(points, k)
	}
	if k := (Point{x: p.x, y: p.y + 1}); inBounds2(grid, k) {
		points = append(points, k)
	}
	return points
}

func inBounds2(grid [][]int, p Point) bool {
	return p.x >= 0 && p.x < len(grid[0])*5 && p.y >= 0 && p.y < len(grid)*5
}

func pathfind(grid [][]int, start, end Point) int {
	openSet := NewOpenSet()
	openSet.Push(start, 0, manhattanDistance(start, end))
	closedSet := NewClosedSet()
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		if cur == end {
			return curg
		}
		for _, i := range neighbors(grid, cur) {
			if closedSet.Has(i) {
				continue
			}
			g := curg + grid[i.y][i.x]
			f := g + manhattanDistance(i, end)
			if v, ok := openSet.Get(i); ok {
				if g < v.g {
					openSet.Update(i, g, f)
				}
				continue
			}
			openSet.Push(i, g, f)
		}
	}
	return -1
}

func neighbors(grid [][]int, p Point) []Point {
	points := make([]Point, 0, 4)
	if k := (Point{x: p.x - 1, y: p.y}); inBounds(grid, k) {
		points = append(points, k)
	}
	if k := (Point{x: p.x, y: p.y - 1}); inBounds(grid, k) {
		points = append(points, k)
	}
	if k := (Point{x: p.x + 1, y: p.y}); inBounds(grid, k) {
		points = append(points, k)
	}
	if k := (Point{x: p.x, y: p.y + 1}); inBounds(grid, k) {
		points = append(points, k)
	}
	return points
}

func inBounds(grid [][]int, p Point) bool {
	return p.x >= 0 && p.x < len(grid[0]) && p.y >= 0 && p.y < len(grid)
}
