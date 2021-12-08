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

	count := 0
	count2 := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), " | ")
		if len(arr) < 2 {
			log.Fatalln("Invalid line")
		}
		assigned := map[byte]int{}
		possible := allPossibilities()
		ex6 := fullWires()
		ex5 := fullWires()
		for _, i := range strings.Fields(arr[0]) {
			wires := []byte(i)
			l := len(wires)
			switch l {
			case 2: // 1
				reducePossibilities(assigned, possible, wires, []int{0, 1, 3, 4, 6})
			case 4: // 4
				reducePossibilities(assigned, possible, wires, []int{0, 4, 6})
			case 3: // 7
				reducePossibilities(assigned, possible, wires, []int{1, 3, 4, 6})
			case 7: // 8
			case 6: // 0, 6, 9, common 0, 1, 5, 6
				ex6 = reduceCommon(ex6, wires)
			case 5: // 2, 3, 5, common 0, 3, 6
				ex5 = reduceCommon(ex5, wires)
			}
		}
		if len(ex6) == 4 {
			reducePossibilities(assigned, possible, values(ex6), []int{2, 3, 4})
		}
		if len(ex5) == 3 {
			reducePossibilities(assigned, possible, values(ex5), []int{1, 2, 4, 5})
		}
		if len(assigned) != 7 {
			log.Fatalln("Failed to assign")
		}
		num := 0
		for _, i := range strings.Fields(arr[1]) {
			num *= 10
			wires := []byte(i)
			l := len(wires)
			switch l {
			case 2: // 1
				num += 1
				count++
			case 4: // 4
				num += 4
				count++
			case 3: // 7
				num += 7
				count++
			case 7: // 8
				num += 8
				count++
			case 6: // 0, 6, 9, common 0, 1, 5, 6
				num += translateSeg6(translateWires(wires, assigned))
			case 5: // 2, 3, 5, common 0, 3, 6
				num += translateSeg5(translateWires(wires, assigned))
			}
		}
		count2 += num
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", count)
	fmt.Println("Part 2:", count2)
}

func translateWires(wires []byte, assigned map[byte]int) []int {
	k := make([]int, 0, len(wires))
	for _, i := range wires {
		k = append(k, assigned[i])
	}
	return k
}

func translateSeg6(segs []int) int {
	missing2 := true
	missing3 := true
	missing4 := true
	for _, i := range segs {
		switch i {
		case 2:
			missing2 = false
		case 3:
			missing3 = false
		case 4:
			missing4 = false
		}
	}
	if missing2 {
		return 6
	}
	if missing3 {
		return 0
	}
	if missing4 {
		return 9
	}
	return -1
}

func translateSeg5(segs []int) int {
	has1 := false
	has2 := false
	has4 := false
	has5 := false
	for _, i := range segs {
		switch i {
		case 1:
			has1 = true
		case 2:
			has2 = true
		case 4:
			has4 = true
		case 5:
			has5 = true
		}
	}
	if has2 && has4 {
		return 2
	}
	if has2 && has5 {
		return 3
	}
	if has1 && has5 {
		return 5
	}
	return -1
}

func fullPossibilities() map[int]struct{} {
	return map[int]struct{}{
		0: {},
		1: {},
		2: {},
		3: {},
		4: {},
		5: {},
		6: {},
	}
}

func allPossibilities() map[byte]map[int]struct{} {
	return map[byte]map[int]struct{}{
		'a': fullPossibilities(),
		'b': fullPossibilities(),
		'c': fullPossibilities(),
		'd': fullPossibilities(),
		'e': fullPossibilities(),
		'f': fullPossibilities(),
		'g': fullPossibilities(),
	}
}

func reducePossibilities(assigned map[byte]int, possible map[byte]map[int]struct{}, wires []byte, segs []int) {
	changed := false
	for _, i := range wires {
		for _, j := range segs {
			if _, ok := possible[i][j]; ok {
				changed = true
				delete(possible[i], j)
			}
		}
	}
	if !changed {
		return
	}
	for {
		changed := false
		for k, v := range possible {
			if _, ok := assigned[k]; ok {
				continue
			}
			if len(v) == 1 {
				changed = true
				assigned[k] = firstval(v)
				continue
			}
			for _, b := range assigned {
				if _, ok := v[b]; ok {
					changed = true
					delete(v, b)
				}
			}
		}
		if !changed {
			return
		}
	}
}

func firstval(possible map[int]struct{}) int {
	for k := range possible {
		return k
	}
	return -1
}

func values(s map[byte]struct{}) []byte {
	res := make([]byte, 0, len(s))
	for k := range s {
		res = append(res, k)
	}
	return res
}

func fullWires() map[byte]struct{} {
	return map[byte]struct{}{
		'a': {},
		'b': {},
		'c': {},
		'd': {},
		'e': {},
		'f': {},
		'g': {},
	}
}

func reduceCommon(s map[byte]struct{}, wires []byte) map[byte]struct{} {
	k := map[byte]struct{}{}
	for _, i := range wires {
		if _, ok := s[i]; ok {
			k[i] = struct{}{}
		}
	}
	return k
}
