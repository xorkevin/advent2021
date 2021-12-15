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

func manhattan(a, b Point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func neighbors(grid [][]int, p Point) []Point {
	return []Point{
		{x: p.x - 1, y: p.y},
		{x: p.x, y: p.y - 1},
		{x: p.x + 1, y: p.y},
		{x: p.x, y: p.y + 1},
	}
}

func inBounds(p Point, w, h int) bool {
	return p.x >= 0 && p.x < w && p.y >= 0 && p.y < h
}

func getVal(grid [][]int, p Point, w, h int) int {
	return (grid[p.y%h][p.x%w]+p.x/w+p.y/h-1)%9 + 1
}

func pathfind(grid [][]int, vw, vh int, start, end Point) int {
	w := len(grid[0])
	h := len(grid)
	openSet := NewOpenSet()
	openSet.Push(start, 0, manhattan(start, end))
	closedSet := NewClosedSet()
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		if cur == end {
			return curg
		}
		for _, i := range neighbors(grid, cur) {
			if !inBounds(i, vw, vh) || closedSet.Has(i) {
				continue
			}
			g := curg + getVal(grid, i, w, h)
			f := g + manhattan(i, end)
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

	fmt.Println("Part 1:", pathfind(grid, len(grid[0]), len(grid), Point{0, 0}, Point{len(grid[0]) - 1, len(grid) - 1}))
	fmt.Println("Part 2:", pathfind(grid, len(grid[0])*5, len(grid)*5, Point{0, 0}, Point{len(grid[0])*5 - 1, len(grid)*5 - 1}))
}
