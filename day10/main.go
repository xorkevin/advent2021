package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

const (
	puzzleInput = "input.txt"
)

type (
	Stack struct {
		k []byte
	}
)

func (s *Stack) Push(c byte) {
	s.k = append(s.k, c)
}

func (s *Stack) Pop() (byte, bool) {
	l := len(s.k)
	if l == 0 {
		return 0, false
	}
	v := s.k[l-1]
	s.k = s.k[:l-1]
	return v, true
}

func (s *Stack) Peek() (byte, bool) {
	l := len(s.k)
	if l == 0 {
		return 0, false
	}
	return s.k[l-1], true
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

	score := 0
	var completes []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r, p, c := parsePair([]byte(scanner.Text()))
		if r == 2 {
			score += p
		} else if r == 1 {
			completes = append(completes, c)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Ints(completes)

	fmt.Println("Part 1:", score)
	fmt.Println("Part 2:", completes[len(completes)/2])
}

func parsePair(s []byte) (int, int, int) {
	stack := Stack{
		k: make([]byte, 0, len(s)),
	}
	for _, i := range s {
		switch i {
		case '(', '[', '{', '<':
			stack.Push(i)
		case ')', ']', '}', '>':
			c, ok := stack.Pop()
			if !ok {
				return 3, 0, 0
			}
			switch i {
			case ')':
				if c != '(' {
					return 2, 3, 0
				}
			case ']':
				if c != '[' {
					return 2, 57, 0
				}
			case '}':
				if c != '{' {
					return 2, 1197, 0
				}
			case '>':
				if c != '<' {
					return 2, 25137, 0
				}
			}
		default:
			return 4, 0, 0
		}
	}
	if len(stack.k) != 0 {
		return 1, 0, translate(stack.k)
	}
	return 0, 0, 0
}

func translate(s []byte) int {
	l := len(s)
	k := 0
	for i := l - 1; i >= 0; i-- {
		k *= 5
		switch s[i] {
		case '(':
			k += 1
		case '[':
			k += 2
		case '{':
			k += 3
		case '<':
			k += 4
		}
	}
	return k
}
