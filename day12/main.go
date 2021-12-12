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

type (
	Stack struct {
		k []string
	}
)

func (s *Stack) Push(c string) {
	s.k = append(s.k, c)
}

func (s *Stack) Pop() (string, bool) {
	l := len(s.k)
	if l == 0 {
		return "", false
	}
	v := s.k[l-1]
	s.k = s.k[:l-1]
	return v, true
}

func (s *Stack) Peek() (string, bool) {
	l := len(s.k)
	if l == 0 {
		return "", false
	}
	return s.k[l-1], true
}

type (
	Graph struct {
		nodes map[string]bool
		edges map[string][]string
	}
)

func NewGraph() *Graph {
	return &Graph{
		nodes: map[string]bool{},
		edges: map[string][]string{},
	}
}

func isUpper(s string) bool {
	return s == strings.ToUpper(s)
}

func (g *Graph) AddEdge(a, b string) {
	g.nodes[a] = isUpper(a)
	g.nodes[b] = isUpper(b)
	g.edges[a] = append(g.edges[a], b)
	g.edges[b] = append(g.edges[b], a)
}

func (g *Graph) findPath(start, end string, path *Stack) int {
	if start == end {
		return 1
	}
	visitedSet := map[string]struct{}{}
	for _, i := range path.k {
		visitedSet[i] = struct{}{}
	}
	count := 0
	for _, i := range g.edges[start] {
		if _, ok := visitedSet[i]; ok && !g.nodes[i] {
			continue
		}
		path.Push(i)
		count += g.findPath(i, end, path)
		path.Pop()
	}
	return count
}

func (g *Graph) FindPath(start, end string) int {
	path := Stack{}
	path.Push(start)
	return g.findPath(start, end, &path)
}

func (g *Graph) findPath2(start, end string, path *Stack) int {
	if start == end {
		return 1
	}
	visitedSet := map[string]struct{}{}
	hasTwice := false
	for _, i := range path.k {
		if !hasTwice {
			if _, ok := visitedSet[i]; ok && !g.nodes[i] {
				hasTwice = true
			}
		}
		visitedSet[i] = struct{}{}
	}
	count := 0
	for _, i := range g.edges[start] {
		if _, ok := visitedSet[i]; ok && !g.nodes[i] {
			if hasTwice || i == "start" || i == "end" {
				continue
			}
		}
		path.Push(i)
		count += g.findPath2(i, end, path)
		path.Pop()
	}
	return count
}

func (g *Graph) FindPath2(start, end string) int {
	path := Stack{}
	path.Push(start)
	return g.findPath2(start, end, &path)
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

	graph := NewGraph()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr := strings.SplitN(scanner.Text(), "-", 2)
		if len(arr) != 2 {
			log.Fatalln("Invalid line")
		}
		graph.AddEdge(arr[0], arr[1])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", graph.FindPath("start", "end"))
	fmt.Println("Part 2:", graph.FindPath2("start", "end"))
}
