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

	Board struct {
		Width    int
		Height   int
		Unmarked map[int]Pos
		Marked   map[int]Pos
		Rows     []int
		Cols     []int
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

	var nums []int
	var boards []*Board
	board := &Board{
		Height:   0,
		Unmarked: map[int]Pos{},
		Marked:   map[int]Pos{},
	}
	first := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if first {
			first = false
			for _, i := range strings.Split(scanner.Text(), ",") {
				num, err := strconv.Atoi(i)
				if err != nil {
					log.Fatal(err)
				}
				nums = append(nums, num)
			}
			continue
		}
		line := scanner.Text()
		if line == "" {
			if board.Height > 0 {
				board.Rows = make([]int, board.Height)
				board.Cols = make([]int, board.Width)
				boards = append(boards, board)
				board = &Board{
					Height:   0,
					Unmarked: map[int]Pos{},
					Marked:   map[int]Pos{},
				}
			}
			continue
		}
		h := board.Height
		board.Height++
		row := strings.Fields(line)
		board.Width = len(row)
		for n, i := range row {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatal(err)
			}
			if _, ok := board.Unmarked[num]; ok {
				log.Fatalln("Duplicate num")
			}
			board.Unmarked[num] = Pos{
				X: n,
				Y: h,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if board.Height > 0 {
		board.Rows = make([]int, board.Height)
		board.Cols = make([]int, board.Width)
		boards = append(boards, board)
		board = nil
	}

	skipMap := map[int]struct{}{}
	first = true
	for _, i := range nums {
		for n, j := range boards {
			if _, ok := skipMap[n]; ok {
				continue
			}
			k := markBoard(j, i)
			if k > -1 {
				if first {
					first = false
					fmt.Println("Part 1:", k*i)
				}
				skipMap[n] = struct{}{}
				if len(skipMap) >= len(boards) {
					fmt.Println("Part 2:", k*i)
					return
				}
			}
		}
	}
}

func markBoard(board *Board, num int) int {
	p, ok := board.Unmarked[num]
	if !ok {
		return -1
	}
	board.Marked[num] = p
	delete(board.Unmarked, num)
	board.Rows[p.Y]++
	board.Cols[p.X]++
	if board.Rows[p.Y] >= board.Height {
		return scoreBoard(board)
	}
	if board.Cols[p.X] >= board.Width {
		return scoreBoard(board)
	}
	return -1
}

func scoreBoard(board *Board) int {
	sum := 0
	for v := range board.Unmarked {
		sum += v
	}
	return sum
}
