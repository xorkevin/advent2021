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

	var nums [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		nums = append(nums, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(nums) == 0 {
		return
	}

	numbits := len(nums[0])

	count := make([]int, numbits)
	for _, i := range nums {
		for n, j := range i {
			if j == '1' {
				count[n]++
			} else {
				count[n]--
			}
		}
	}
	max := make([]byte, numbits)
	min := make([]byte, numbits)
	for n, i := range count {
		if i > 0 {
			max[n] = '1'
			min[n] = '0'
		} else {
			max[n] = '0'
			min[n] = '1'
		}
	}
	fmt.Println("Part 1:", btoi(max)*btoi(min))

	most := nums
	least := nums

	for i := 0; i < numbits; i++ {
		if len(most) < 2 {
			break
		}
		most = findCommon(true, i, most)
	}
	for i := 0; i < numbits; i++ {
		if len(least) < 2 {
			break
		}
		least = findCommon(false, i, least)
	}
	if len(most) == 1 && len(least) == 1 {
		fmt.Println("Part 2:", btoi(most[0])*btoi(least[0]))
	}
}

func findCommon(most bool, pos int, nums [][]byte) [][]byte {
	onecounts := 0
	zerocounts := 0

	for _, i := range nums {
		if i[pos] == '1' {
			onecounts++
		} else {
			zerocounts++
		}
	}

	zeros := zerocounts <= onecounts
	if most {
		zeros = zerocounts > onecounts
	}

	var res [][]byte

	for _, i := range nums {
		if zeros {
			if i[pos] == '0' {
				res = append(res, i)
			}
		} else {
			if i[pos] == '1' {
				res = append(res, i)
			}
		}
	}

	return res
}

func btoi(b []byte) int {
	c := 0
	for _, i := range b {
		k := 0
		if i == '1' {
			k = 1
		}
		c = c<<1 + k
	}
	return c
}
