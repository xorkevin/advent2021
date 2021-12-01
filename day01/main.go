package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

	prev := 0
	prev1 := 0
	prevsum := 0
	count := -1
	count2 := -3
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		if num > prev {
			count++
		}
		k := num + prev + prev1
		if k > prevsum {
			count2++
		}
		prevsum = k
		prev1 = prev
		prev = num
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", count)
	fmt.Println("Part 2:", count2)
}
