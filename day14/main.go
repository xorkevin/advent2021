package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	var first, last byte
	pairs := map[string]int{}
	rules := map[string]byte{}
	scanner := bufio.NewScanner(file)
	batch1 := true
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			batch1 = false
			continue
		}
		if batch1 {
			for i := 0; i < len(line)-1; i++ {
				p := line[i : i+2]
				if _, ok := pairs[p]; !ok {
					pairs[p] = 0
				}
				pairs[p]++
			}
			first = line[0]
			last = line[len(line)-1]
			continue
		}
		arr := strings.SplitN(line, " -> ", 2)
		if len(arr) != 2 {
			log.Fatalln("Invalid line")
		}
		if len(arr[1]) != 1 {
			log.Fatalln("Invalid line")
		}
		rules[arr[0]] = arr[1][0]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		pairs = processStep(rules, pairs)
	}
	max, min := maxminCount(pairs, first, last)
	fmt.Println("Part 1:", max-min)
	for i := 0; i < 30; i++ {
		pairs = processStep(rules, pairs)
	}
	max, min = maxminCount(pairs, first, last)
	fmt.Println("Part 2:", max-min)
}

func processStep(rules map[string]byte, pairs map[string]int) map[string]int {
	next := map[string]int{}
	for k, v := range pairs {
		p := rules[k]
		p1 := string([]byte{k[0], p})
		p2 := string([]byte{p, k[1]})
		if _, ok := next[p1]; !ok {
			next[p1] = 0
		}
		next[p1] += v
		if _, ok := next[p2]; !ok {
			next[p2] = 0
		}
		next[p2] += v
	}
	return next
}

func maxminCount(pairs map[string]int, first, last byte) (int, int) {
	max := 0
	min := 0
	counts := map[byte]int{
		first: 1,
		last:  1,
	}
	for k, v := range pairs {
		c1 := k[0]
		if _, ok := counts[c1]; !ok {
			counts[c1] = 0
		}
		counts[c1] += v
		c2 := k[1]
		if _, ok := counts[c2]; !ok {
			counts[c2] = 0
		}
		counts[c2] += v
		max = counts[c2] / 2
		min = counts[c2] / 2
	}
	for _, v := range counts {
		if v/2 > max {
			max = v / 2
		}
		if v/2 < min {
			min = v / 2
		}
	}
	return max, min
}
