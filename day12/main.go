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
	Node struct {
		name string
		big  bool
	}

	Graph struct {
		nodes map[string]Node
		edges map[string][]string
	}
)

func NewGraph() *Graph {
	return &Graph{
		nodes: map[string]Node{},
		edges: map[string][]string{},
	}
}

func isUpper(s string) bool {
	return s == strings.ToUpper(s)
}

func (g *Graph) AddEdge(a, b string) {
	an := Node{
		name: a,
		big:  isUpper(a),
	}
	bn := Node{
		name: b,
		big:  isUpper(b),
	}
	g.nodes[a] = an
	g.nodes[b] = bn
	g.edges[a] = append(g.edges[a], b)
	g.edges[b] = append(g.edges[b], a)
}

func (g *Graph) findPath(start, end string, path []string) int {
	if start == end {
		return 1
	}
	visitedSet := map[string]struct{}{}
	for _, i := range path {
		visitedSet[i] = struct{}{}
	}
	count := 0
	for _, i := range g.edges[start] {
		if _, ok := visitedSet[i]; ok && !g.nodes[i].big {
			continue
		}
		count += g.findPath(i, end, append(path, i))
	}
	return count
}

func (g *Graph) FindPath(start, end string) int {
	return g.findPath(start, end, []string{start})
}

func (g *Graph) findPath2(start, end string, path []string) int {
	if start == end {
		return 1
	}
	visitedSet := map[string]struct{}{}
	hasTwice := false
	for _, i := range path {
		if _, ok := visitedSet[i]; ok && !g.nodes[i].big {
			hasTwice = true
		} else {
			visitedSet[i] = struct{}{}
		}
	}
	count := 0
	for _, i := range g.edges[start] {
		if _, ok := visitedSet[i]; ok && !g.nodes[i].big {
			if hasTwice || i == "start" || i == "end" {
				continue
			}
		}
		count += g.findPath2(i, end, append(path, i))
	}
	return count
}

func (g *Graph) FindPath2(start, end string) int {
	return g.findPath2(start, end, []string{start})
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

	fmt.Println(graph.FindPath("start", "end"))
	fmt.Println(graph.FindPath2("start", "end"))
}
