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

	Line struct {
		P1 Pos
		P2 Pos
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

	var segs []Line

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), " -> ")
		lhs := strings.Split(arr[0], ",")
		rhs := strings.Split(arr[1], ",")
		lhs1, err := strconv.Atoi(lhs[0])
		if err != nil {
			log.Fatal(err)
		}
		lhs2, err := strconv.Atoi(lhs[1])
		if err != nil {
			log.Fatal(err)
		}
		rhs1, err := strconv.Atoi(rhs[0])
		if err != nil {
			log.Fatal(err)
		}
		rhs2, err := strconv.Atoi(rhs[1])
		if err != nil {
			log.Fatal(err)
		}
		segs = append(segs, Line{
			P1: Pos{
				X: lhs1,
				Y: lhs2,
			},
			P2: Pos{
				X: rhs1,
				Y: rhs2,
			},
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	grid := map[Pos]int{}
	grid2 := map[Pos]int{}
	for _, i := range segs {
		if i.P1.X == i.P2.X {
			start, stop := minmax(i.P1.Y, i.P2.Y)
			for j := start; j <= stop; j++ {
				k := Pos{X: i.P1.X, Y: j}
				if _, ok := grid[k]; !ok {
					grid[k] = 0
				}
				grid[k]++
				if _, ok := grid2[k]; !ok {
					grid2[k] = 0
				}
				grid2[k]++
			}
		} else if i.P1.Y == i.P2.Y {
			start, stop := minmax(i.P1.X, i.P2.X)
			for j := start; j <= stop; j++ {
				k := Pos{X: j, Y: i.P1.Y}
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
			start := i.P1
			stop := i.P2
			dirX := sign(i.P2.X - i.P1.X)
			dirY := sign(i.P2.Y - i.P1.Y)
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
