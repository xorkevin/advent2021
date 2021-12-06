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

	nums := make([]int, 9)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, i := range strings.Split(scanner.Text(), ",") {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatal(err)
			}
			nums[num]++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for k := 0; k < 256; k++ {
		next := make([]int, 9)
		for n, i := range nums {
			if n == 0 {
				next[6] += i
				next[8] += i
			} else {
				next[n-1] += i
			}
		}
		nums = next
		if k == 79 {
			count := 0
			for _, i := range nums {
				count += i
			}
			fmt.Println("Part 1:", count)
		}
	}

	count := 0
	for _, i := range nums {
		count += i
	}
	fmt.Println("Part 2:", count)
}
