package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

const (
	puzzleInput = "input.txt"
)

type (
	Grid struct {
		grid  [][]byte
		h     int
		w     int
		basin [][]byte
	}
)

func NewGrid(grid [][]byte) *Grid {
	h := len(grid)
	w := len(grid[0])
	basin := make([][]byte, h)
	for i := range basin {
		basin[i] = make([]byte, w)
	}
	return &Grid{
		grid:  grid,
		h:     h,
		w:     w,
		basin: basin,
	}
}

func (g *Grid) outBounds(x, y int) bool {
	return x < 0 || y < 0 || x >= g.w || y >= g.h
}

func (g *Grid) inBounds(x, y int) bool {
	return !g.outBounds(x, y)
}

func (g *Grid) isLow(x, y int) bool {
	v := g.grid[y][x]
	if g.inBounds(x-1, y) && g.grid[y][x-1] <= v {
		return false
	}
	if g.inBounds(x, y-1) && g.grid[y-1][x] <= v {
		return false
	}
	if g.inBounds(x+1, y) && g.grid[y][x+1] <= v {
		return false
	}
	if g.inBounds(x, y+1) && g.grid[y+1][x] <= v {
		return false
	}
	return true
}

func (g *Grid) markBasin(k byte, x, y int) int {
	if !g.inBounds(x, y) {
		return 0
	}
	if g.grid[y][x] == '9' {
		return 0
	}
	if g.basin[y][x] != 0 {
		return 0
	}
	g.basin[y][x] = k
	return 1 + g.markBasin(k, x-1, y) + g.markBasin(k, x, y-1) + g.markBasin(k, x+1, y) + g.markBasin(k, x, y+1)
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

	var rows [][]byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rows = append(rows, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	grid := NewGrid(rows)

	count := 0
	var basinnum byte = 1
	var sizes []int
	for r, i := range grid.grid {
		for c := range i {
			if grid.isLow(c, r) {
				count += int(grid.grid[r][c]-'0') + 1
				size := grid.markBasin(basinnum, c, r)
				sizes = append(sizes, size)
				basinnum++
			}
		}
	}
	fmt.Println("Part 1:", count)
	if len(sizes) >= 3 {
		sort.Ints(sizes)
		l := len(sizes)
		fmt.Println("Part 2:", sizes[l-1]*sizes[l-2]*sizes[l-3])
	}
}
