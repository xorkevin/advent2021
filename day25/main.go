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

type (
	Vec2 struct {
		x, y int
	}
)

func right(p Vec2, w int) Vec2 {
	return Vec2{(p.x + 1) % w, p.y}
}

func down(p Vec2, h int) Vec2 {
	return Vec2{p.x, (p.y + 1) % h}
}

func occupied(p Vec2, a, b map[Vec2]struct{}) bool {
	if _, ok := a[p]; ok {
		return true
	}
	_, ok := b[p]
	return ok
}

func next(east, south map[Vec2]struct{}, w, h int) (map[Vec2]struct{}, map[Vec2]struct{}, bool) {
	ne := map[Vec2]struct{}{}
	ns := map[Vec2]struct{}{}
	changed := false
	for i := range east {
		k := right(i, w)
		if occupied(k, east, south) {
			ne[i] = struct{}{}
		} else {
			ne[k] = struct{}{}
			changed = true
		}
	}
	for i := range south {
		k := down(i, h)
		if occupied(k, ne, south) {
			ns[i] = struct{}{}
		} else {
			ns[k] = struct{}{}
			changed = true
		}
	}
	return ne, ns, changed
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

	east := map[Vec2]struct{}{}
	south := map[Vec2]struct{}{}
	scanner := bufio.NewScanner(file)
	linenum := 0
	width := 0
	for scanner.Scan() {
		line := scanner.Text()
		width = len(line)
		for n, i := range []byte(line) {
			switch i {
			case '>':
				east[Vec2{n, linenum}] = struct{}{}
			case 'v':
				south[Vec2{n, linenum}] = struct{}{}
			}
		}
		linenum++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	iter := 0
	for {
		iter++
		var changed bool
		east, south, changed = next(east, south, width, linenum)
		if !changed {
			break
		}
	}
	fmt.Println(iter)
	//grid := make([][]byte, linenum)
	//for i := range grid {
	//	k := make([]byte, width)
	//	for j := range k {
	//		k[j] = ' '
	//	}
	//	grid[i] = k
	//}
	//for k := range east {
	//	grid[k.y][k.x] = '>'
	//}
	//for k := range south {
	//	grid[k.y][k.x] = 'v'
	//}
	//for _, i := range grid {
	//	fmt.Println(string(i))
	//}
}
