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

	pos := 0
	depth := 0
	pos2 := 0
	depth2 := 0
	aim := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr := strings.SplitN(scanner.Text(), " ", 2)
		if len(arr) < 2 {
			log.Fatalln("Invalid line format")
		}
		dir := arr[0]
		num, err := strconv.Atoi(arr[1])
		if err != nil {
			log.Fatal(err)
		}
		switch dir {
		case "forward":
			pos += num
			pos2 += num
			depth2 += aim * num
		case "down":
			depth += num
			aim += num
		case "up":
			depth -= num
			aim -= num
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", pos*depth)
	fmt.Println("Part 2:", pos2*depth2)
}
