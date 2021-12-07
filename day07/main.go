package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
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

	var nums []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, i := range strings.Split(scanner.Text(), ",") {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatal(err)
			}
			nums = append(nums, num)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Ints(nums)
	median := nums[len(nums)/2]
	sum := 0
	for _, i := range nums {
		sum += i
	}
	a1 := int(math.Floor(float64(sum) / float64(len(nums))))
	a2 := int(math.Ceil(float64(sum) / float64(len(nums))))
	diffs := 0
	diffs1 := 0
	diffs2 := 0
	for _, i := range nums {
		diffs += abs(i, median)
		k1 := abs(i, a1)
		diffs1 += k1 * (k1 + 1) / 2
		k2 := abs(i, a2)
		diffs2 += k2 * (k2 + 1) / 2
	}
	fmt.Println("Part 1:", diffs)
	fmt.Println("Part 2:", min(diffs1, diffs2))
}

func abs(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
