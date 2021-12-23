package main

import (
	"container/heap"
	"fmt"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Item struct {
		value State
		g, f  int
		index int
	}

	PriorityQueue struct {
		q []*Item
		s map[State]int
	}

	OpenSet struct {
		q *PriorityQueue
	}

	ClosedSet map[State]struct{}
)

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		s: map[State]int{},
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
func (q *PriorityQueue) Update(value State, g, f int) bool {
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

func (s *OpenSet) Has(val State) bool {
	_, ok := s.q.s[val]
	return ok
}

func (s *OpenSet) Get(val State) (*Item, bool) {
	idx, ok := s.q.s[val]
	if !ok {
		return nil, false
	}
	return s.q.q[idx], true
}

func (s *OpenSet) Push(value State, g, f int) {
	heap.Push(s.q, &Item{
		value: value,
		g:     g,
		f:     f,
	})
}

func (s *OpenSet) Pop() (State, int, int) {
	item := heap.Pop(s.q).(*Item)
	return item.value, item.g, item.f
}

func (s *OpenSet) Update(value State, g, f int) bool {
	return s.q.Update(value, g, f)
}

func NewClosedSet() ClosedSet {
	return ClosedSet{}
}

func (cs ClosedSet) Has(val State) bool {
	_, ok := cs[val]
	return ok
}

func (cs ClosedSet) Push(val State) {
	cs[val] = struct{}{}
}

type (
	Vec2 struct {
		x, y int
	}

	State [4][4]Vec2
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (v Vec2) Manhattan(o Vec2) int {
	return abs(v.x-o.x) + abs(v.y-o.y)
}

func (s State) String() string {
	grid := [][]byte{
		[]byte("#############"),
		[]byte("#...........#"),
		[]byte("###.#.#.#.###"),
		[]byte("  #.#.#.#.#"),
		[]byte("  #.#.#.#.#"),
		[]byte("  #.#.#.#.#"),
		[]byte("  #########"),
	}
	for n, i := range s {
		for _, j := range i {
			if j.x == 0 {
				continue
			}
			grid[j.y][j.x] = 'A' + byte(n)
		}
	}
	b := strings.Builder{}
	for _, i := range grid {
		b.WriteString(string(i))
		b.WriteByte('\n')
	}
	return b.String()
}

func (s State) Heuristic(num int) int {
	k := 0
	for n, i := range s {
		for m, j := range i {
			if m >= num {
				continue
			}
			k += abs(j.x-winLocs[n]) * costs[n]
		}
	}
	return k
}

type (
	StateOpt struct {
		value State
		cost  int
	}
)

var (
	winLocs = [4]int{3, 5, 7, 9}
	teamPos = [12]int{-1, -1, -1, 0, -1, 1, -1, 2, -1, 3, -1, -1}
	costs   = [4]int{1, 10, 100, 1000}
)

func isWin(s State, num int) bool {
	for n, i := range s {
		for m, j := range i {
			if m >= num {
				continue
			}
			if j.x != winLocs[n] {
				return false
			}
		}
	}
	return true
}

func isHallway(v Vec2) bool {
	return v.y < 2
}

func genHallway(state State, num, base int) ([12]bool, [12]int, [4]bool) {
	hall := [12]bool{}
	depth := [12]int{}
	clear := [4]bool{}
	for i := range depth {
		depth[i] = base
	}
	for i := range clear {
		clear[i] = true
	}
	for n, i := range state {
		for m, j := range i {
			if m >= num {
				continue
			}
			if isHallway(j) {
				hall[j.x] = true
			} else {
				if j.y <= depth[j.x] {
					depth[j.x] = j.y - 1
				}
				if j.x != winLocs[n] {
					clear[teamPos[j.x]] = false
				}
			}
		}
	}
	return hall, depth, clear
}

func hallwayClear(a, b int, hall [12]bool) bool {
	origA := a
	if a > b {
		a, b = b, a
	}
	for i := a; i <= b; i++ {
		if i == origA {
			continue
		}
		if hall[i] {
			return false
		}
	}
	return true
}

func getNeighbors(state State, num, base int) []StateOpt {
	var opts []StateOpt
	hallway, depth, clear := genHallway(state, num, base)
	for n, i := range state {
		for m, j := range i {
			if m >= num {
				continue
			}
			if isHallway(j) {
				// once in hallway, must move into correct room, unless other kind has occupied
				if !clear[n] {
					continue
				}
				tx := winLocs[n]
				if depth[tx] < 2 {
					continue
				}
				if !hallwayClear(j.x, tx, hallway) {
					continue
				}
				target := Vec2{
					x: tx,
					y: depth[tx],
				}
				k := state
				k[n][m] = target
				opts = append(opts, StateOpt{
					value: k,
					cost:  j.Manhattan(target) * costs[n],
				})
			} else {
				// if not top, cannot move
				if j.y != depth[j.x]+1 {
					continue
				}
				// if clear and in winning position, no need to move
				if clear[n] && j.x == winLocs[n] {
					continue
				}
				for _, x := range []int{1, 2, 4, 6, 8, 10, 11} {
					if !hallwayClear(j.x, x, hallway) {
						continue
					}
					target := Vec2{
						x: x,
						y: 1,
					}
					k := state
					k[n][m] = target
					opts = append(opts, StateOpt{
						value: k,
						cost:  j.Manhattan(target) * costs[n],
					})
				}
			}
		}
	}
	return opts
}

func pathfind(start State, num, base int) int {
	openSet := NewOpenSet()
	openSet.Push(start, 0, start.Heuristic(num))
	closedSet := NewClosedSet()
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		if isWin(cur, num) {
			return curg
		}
		for _, o := range getNeighbors(cur, num, base) {
			if closedSet.Has(o.value) {
				continue
			}
			g := curg + o.cost
			f := g + o.value.Heuristic(num)
			if v, ok := openSet.Get(o.value); ok {
				if g < v.g {
					openSet.Update(o.value, g, f)
				}
				continue
			}
			openSet.Push(o.value, g, f)
		}
	}
	return -1
}

func calcStart(inp string) State {
	s := State{}
	for y, i := range strings.Split(strings.TrimSpace(inp), "\n") {
		for x, j := range []byte(i) {
			switch j {
			case 'A', 'B', 'C', 'D':
			default:
				continue
			}
			d := -1
			for n, i := range s[j-'A'] {
				if i.x == 0 {
					d = n
					break
				}
			}
			s[j-'A'][d] = Vec2{x, y}
		}
	}
	return s
}

func main() {
	startState := calcStart(`
#############
#...........#
###C#B#A#D###
  #C#D#A#B#
  #########
`)

	fmt.Println(startState)
	fmt.Println(pathfind(startState, 2, 3))

	startState = calcStart(`
#############
#...........#
###C#B#A#D###
  #D#C#B#A#
  #D#B#A#C#
  #C#D#A#B#
  #########
`)

	fmt.Println(startState)
	fmt.Println(pathfind(startState, 4, 5))
}
