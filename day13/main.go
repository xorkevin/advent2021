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
		x int
		y int
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

	points := map[Pos]struct{}{}
	trackPoints := true
	first := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			trackPoints = false
			continue
		}
		if trackPoints {
			arr := strings.SplitN(line, ",", 2)
			if len(arr) != 2 {
				log.Fatalln("Invalid line")
			}
			x, err := strconv.Atoi(arr[0])
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(arr[1])
			if err != nil {
				log.Fatal(err)
			}
			points[Pos{
				x: x,
				y: y,
			}] = struct{}{}
			continue
		}
		words := strings.Fields(line)
		if len(words) != 3 {
			log.Fatalln("Invalid line")
		}
		arr := strings.SplitN(words[2], "=", 2)
		if len(arr) != 2 {
			log.Fatalln("Invalid line")
		}
		axis := arr[0]
		axisval, err := strconv.Atoi(arr[1])
		if err != nil {
			log.Fatal(err)
		}
		points = fold(points, axis == "y", axisval)
		if first {
			first = false
			fmt.Println(len(points))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	maxx, maxy := findMax(points)
	grid := make([][]byte, maxy)
	for i := range grid {
		row := make([]byte, maxx)
		for j := range row {
			row[j] = ' '
		}
		grid[i] = row
	}
	for k := range points {
		grid[k.y][k.x] = '#'
	}

	for _, i := range grid {
		fmt.Println(string(i))
	}
}

func fold(points map[Pos]struct{}, yaxis bool, axisval int) map[Pos]struct{} {
	next := map[Pos]struct{}{}
	for k := range points {
		if yaxis {
			if k.y <= axisval {
				next[k] = struct{}{}
				continue
			}
			next[Pos{
				x: k.x,
				y: axisval - (k.y - axisval),
			}] = struct{}{}
		} else {
			if k.x <= axisval {
				next[k] = struct{}{}
				continue
			}
			next[Pos{
				x: axisval - (k.x - axisval),
				y: k.y,
			}] = struct{}{}
		}
	}
	return next
}

func findMax(points map[Pos]struct{}) (int, int) {
	x := 0
	y := 0
	for k := range points {
		if k.x > x {
			x = k.x
		}
		if k.y > y {
			y = k.y
		}
	}
	return x + 1, y + 1
}
