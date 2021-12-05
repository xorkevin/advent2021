package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Pos struct {
		X int
		Y int
	}
)

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

	grid := map[Pos]int{}
	grid2 := map[Pos]int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), " -> ")
		lhs := strings.Split(arr[0], ",")
		rhs := strings.Split(arr[1], ",")
		x1, err := strconv.Atoi(lhs[0])
		if err != nil {
			log.Fatal(err)
		}
		y1, err := strconv.Atoi(lhs[1])
		if err != nil {
			log.Fatal(err)
		}
		x2, err := strconv.Atoi(rhs[0])
		if err != nil {
			log.Fatal(err)
		}
		y2, err := strconv.Atoi(rhs[1])
		if err != nil {
			log.Fatal(err)
		}
		if x1 == x2 {
			start, stop := minmax(y1, y2)
			for j := start; j <= stop; j++ {
				k := Pos{X: x1, Y: j}
				if _, ok := grid[k]; !ok {
					grid[k] = 0
				}
				grid[k]++
				if _, ok := grid2[k]; !ok {
					grid2[k] = 0
				}
				grid2[k]++
			}
		} else if y1 == y2 {
			start, stop := minmax(x1, x2)
			for j := start; j <= stop; j++ {
				k := Pos{X: j, Y: y1}
				if _, ok := grid[k]; !ok {
					grid[k] = 0
				}
				grid[k]++
				if _, ok := grid2[k]; !ok {
					grid2[k] = 0
				}
				grid2[k]++
			}
		} else {
			start := Pos{X: x1, Y: y1}
			stop := Pos{X: x2, Y: y2}
			dirX := sign(x2 - x1)
			dirY := sign(y2 - y1)
			for start.X != stop.X {
				if _, ok := grid2[start]; !ok {
					grid2[start] = 0
				}
				grid2[start]++
				start.X += dirX
				start.Y += dirY
			}
			if _, ok := grid2[stop]; !ok {
				grid2[stop] = 0
			}
			grid2[stop]++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, v := range grid {
		if v > 1 {
			count++
		}
	}
	fmt.Println("Part 1:", count)
	count2 := 0
	for _, v := range grid2 {
		if v > 1 {
			count2++
		}
	}
	fmt.Println("Part 2:", count2)
}

func minmax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func sign(a int) int {
	if a > 0 {
		return 1
	}
	return -1
}
