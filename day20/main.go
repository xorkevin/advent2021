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
	Vec2 struct {
		x, y int
	}
)

func getPixelIndex(grid map[Vec2]struct{}, p Vec2) int {
	k := 0
	for a := -1; a <= 1; a++ {
		y := p.y + a
		for b := -1; b <= 1; b++ {
			x := p.x + b
			k = k << 1
			if _, ok := grid[Vec2{x, y}]; ok {
				k += 1
			}
		}
	}
	return k
}

const (
	modeBase = iota
	modeInv
	modeUnInv
)

func enhanceAlg(grid map[Vec2]struct{}, alg []byte, mode int) map[Vec2]struct{} {
	next := map[Vec2]struct{}{}
	for p := range grid {
		for a := -1; a <= 1; a++ {
			y := p.y + a
			for b := -1; b <= 1; b++ {
				x := p.x + b
				c := getPixelIndex(grid, Vec2{x, y})
				switch mode {
				case modeInv:
					// in this mode, we are reading '#', and are storing '.'
					if alg[c] == '.' {
						next[Vec2{x, y}] = struct{}{}
					}
				case modeUnInv:
					// in this mode, we are reading '.', and are storing '#'
					if alg[c^0b111111111] == '#' {
						next[Vec2{x, y}] = struct{}{}
					}
				default:
					// in this mode, we are reading '#', and are storing '#'
					if alg[c] == '#' {
						next[Vec2{x, y}] = struct{}{}
					}
				}
			}
		}
	}
	return next
}

func vecNeg(p Vec2) Vec2 {
	return Vec2{
		x: -p.x,
		y: -p.y,
	}
}

func vecAdd(a, b Vec2) Vec2 {
	return Vec2{
		x: a.x + b.x,
		y: a.y + b.y,
	}
}

func gridMaxMin(grid map[Vec2]struct{}) (Vec2, Vec2) {
	maxx := -999999
	maxy := -999999
	minx := 999999
	miny := 999999
	for i := range grid {
		if i.x > maxx {
			maxx = i.x
		}
		if i.x < minx {
			minx = i.x
		}
		if i.y > maxy {
			maxy = i.y
		}
		if i.y < miny {
			miny = i.y
		}
	}
	return Vec2{maxx, maxy}, Vec2{minx, miny}
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

	var alg []byte
	var grid [][]byte
	first := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []byte(scanner.Text())
		if first {
			first = false
			alg = line
			continue
		}
		if len(line) == 0 {
			continue
		}
		grid = append(grid, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	points := map[Vec2]struct{}{}
	for i, r := range grid {
		for j, v := range r {
			if v == '#' {
				points[Vec2{j, i}] = struct{}{}
			}
		}
	}

	points = enhanceAlg(points, alg, modeInv)
	points = enhanceAlg(points, alg, modeUnInv)

	fmt.Println("Part 1:", len(points))

	for i := 0; i < 24; i++ {
		points = enhanceAlg(points, alg, modeInv)
		points = enhanceAlg(points, alg, modeUnInv)
	}

	fmt.Println("Part 2:", len(points))
}
