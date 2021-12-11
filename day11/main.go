package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

type (
	Pos struct {
		x, y int
	}

	Grid struct {
		rows [][]int
		w    int
		h    int
	}
)

func NewGrid(rows [][]int) *Grid {
	w := 0
	h := len(rows)
	if h > 0 {
		w = len(rows[0])
	}
	return &Grid{
		rows: rows,
		w:    w,
		h:    h,
	}
}

func (g *Grid) Step() int {
	flashes := map[Pos]struct{}{}
	var flashQueue []Pos
	for r, i := range g.rows {
		for c := range i {
			g.rows[r][c]++
			if g.rows[r][c] > 9 {
				k := Pos{
					x: c,
					y: r,
				}
				flashes[k] = struct{}{}
				flashQueue = append(flashQueue, k)
			}
		}
	}
	for len(flashQueue) > 0 {
		var next []Pos
		for _, i := range flashQueue {
			if g.incr(i.x-1, i.y, flashes) {
				k := Pos{
					x: i.x - 1,
					y: i.y,
				}
				flashes[k] = struct{}{}
				next = append(next, k)
			}
			if g.incr(i.x-1, i.y-1, flashes) {
				k := Pos{
					x: i.x - 1,
					y: i.y - 1,
				}
				flashes[k] = struct{}{}
				next = append(next, k)
			}
			if g.incr(i.x, i.y-1, flashes) {
				k := Pos{
					x: i.x,
					y: i.y - 1,
				}
				flashes[k] = struct{}{}
				next = append(next, k)
			}
			if g.incr(i.x+1, i.y-1, flashes) {
				k := Pos{
					x: i.x + 1,
					y: i.y - 1,
				}
				flashes[k] = struct{}{}
				next = append(next, k)
			}
			if g.incr(i.x+1, i.y, flashes) {
				k := Pos{
					x: i.x + 1,
					y: i.y,
				}
				flashes[k] = struct{}{}
				next = append(next, k)
			}
			if g.incr(i.x+1, i.y+1, flashes) {
				k := Pos{
					x: i.x + 1,
					y: i.y + 1,
				}
				flashes[k] = struct{}{}
				next = append(next, k)
			}
			if g.incr(i.x, i.y+1, flashes) {
				k := Pos{
					x: i.x,
					y: i.y + 1,
				}
				flashes[k] = struct{}{}
				next = append(next, k)
			}
			if g.incr(i.x-1, i.y+1, flashes) {
				k := Pos{
					x: i.x - 1,
					y: i.y + 1,
				}
				flashes[k] = struct{}{}
				next = append(next, k)
			}
		}
		flashQueue = next
	}
	for i := range flashes {
		g.rows[i.y][i.x] = 0
	}
	return len(flashes)
}

func (g *Grid) incr(x, y int, flashes map[Pos]struct{}) bool {
	if g.outBounds(x, y) {
		return false
	}
	g.rows[y][x]++
	if g.rows[y][x] > 9 {
		_, ok := flashes[Pos{x: x, y: y}]
		return !ok
	}
	return false
}

func (g *Grid) outBounds(x, y int) bool {
	return x < 0 || y < 0 || x >= g.w || y >= g.h
}

func (g *Grid) inBounds(x, y int) bool {
	return !g.outBounds(x, y)
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

	var rows [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []byte(scanner.Text())
		row := make([]int, 0, len(line))
		for _, i := range line {
			row = append(row, int(i)-'0')
		}
		rows = append(rows, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	grid := NewGrid(rows)
	count := 0
	for i := 0; i < 100; i++ {
		count += grid.Step()
	}
	fmt.Println("Part 1:", count)
	total := grid.w * grid.h
	step := 100
	for {
		step++
		if grid.Step() == total {
			fmt.Println("Part 2:", step)
			return
		}
	}
}
