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

type (
	Vec3 struct {
		x, y, z int
	}

	Mat3 struct {
		x, y, z Vec3
	}

	Edge struct {
		a Vec3
		b Vec3
	}

	ScannerLog struct {
		Scans []Vec3
		Dists map[Vec3][]Edge
	}

	PossibleEdges struct {
		a []Edge
		b []Edge
	}
)

func abs(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func orderedDist(a, b Vec3) Vec3 {
	k := Vec3{
		x: abs(a.x, b.x),
		y: abs(a.y, b.y),
		z: abs(a.z, b.z),
	}
	if k.x > k.y {
		k.x, k.y = k.y, k.x
	}
	if k.y > k.z {
		k.y, k.z = k.z, k.y
	}
	if k.x > k.y {
		k.x, k.y = k.y, k.x
	}
	return k
}

func vecDot(a, b Vec3) int {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

func matDot(a Mat3, b Vec3) Vec3 {
	return Vec3{
		x: vecDot(a.x, b),
		y: vecDot(a.y, b),
		z: vecDot(a.z, b),
	}
}

func vecNeg(v Vec3) Vec3 {
	return Vec3{
		x: -v.x,
		y: -v.y,
		z: -v.z,
	}
}

func vecSum(a, b Vec3) Vec3 {
	return Vec3{
		x: a.x + b.x,
		y: a.y + b.y,
		z: a.z + b.z,
	}
}

func computeDists(scans []Vec3) map[Vec3][]Edge {
	dists := map[Vec3][]Edge{}
	l := len(scans)
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			k := orderedDist(scans[i], scans[j])
			dists[k] = append(dists[k], Edge{scans[i], scans[j]})
		}
	}
	return dists
}

func NewScannerLog(scans []Vec3) *ScannerLog {
	return &ScannerLog{
		Scans: scans,
		Dists: computeDists(scans),
	}
}

func intersectDists(s, o map[Vec3][]Edge) []PossibleEdges {
	var dists []PossibleEdges
	for k, v := range o {
		if min(len(s[k]), len(v)) > 0 {
			a := make([]Edge, len(s[k]))
			copy(a, s[k])
			b := make([]Edge, len(v))
			copy(b, v)
			dists = append(dists, PossibleEdges{
				a: a,
				b: b,
			})
		}
	}
	return dists
}

func calculateTranslation(possibleEdges []PossibleEdges, pointTranslation map[Vec3]Vec3, target int) bool {
	if len(possibleEdges) == 0 {
		return len(pointTranslation) >= target
	}
	first := possibleEdges[0]
	for _, i := range first.a {
		for _, j := range first.b {
			{
				// suppose i.a is j.a and i.b is j.b
				origA, hasA := pointTranslation[i.a]
				origB, hasB := pointTranslation[i.b]
				if hasA && origA != j.a || hasB && origB != j.b {
					goto second
				}
				if !hasA {
					pointTranslation[i.a] = j.a
				}
				if !hasB {
					pointTranslation[i.b] = j.b
				}
				if ok := calculateTranslation(possibleEdges[1:], pointTranslation, target); ok {
					return true
				}
				if !hasA {
					delete(pointTranslation, i.a)
				}
				if !hasB {
					delete(pointTranslation, i.b)
				}
			}
		second:
			{
				// suppose i.a is j.b and i.b is j.a
				origA, hasA := pointTranslation[i.a]
				origB, hasB := pointTranslation[i.b]
				if hasA && origA != j.b || hasB && origB != j.a {
					goto third
				}
				if !hasA {
					pointTranslation[i.a] = j.b
				}
				if !hasB {
					pointTranslation[i.b] = j.a
				}
				if ok := calculateTranslation(possibleEdges[1:], pointTranslation, target); ok {
					return true
				}
				if !hasA {
					delete(pointTranslation, i.a)
				}
				if !hasB {
					delete(pointTranslation, i.b)
				}
			}
		third:
			if ok := calculateTranslation(possibleEdges[1:], pointTranslation, target); ok {
				return true
			}
		}
	}
	return false
}

func get3Vec(m map[Vec3]Vec3) ([]Vec3, []Vec3) {
	var a, b []Vec3
	for k, v := range m {
		if len(a) >= 3 {
			break
		}
		a = append(a, k)
		b = append(b, v)
	}
	return a, b
}

var (
	rotationMatricies = []Mat3{
		{Vec3{1, 0, 0}, Vec3{0, 1, 0}, Vec3{0, 0, 1}},
		{Vec3{-1, 0, 0}, Vec3{0, -1, 0}, Vec3{0, 0, 1}},
		{Vec3{-1, 0, 0}, Vec3{0, 1, 0}, Vec3{0, 0, -1}},
		{Vec3{1, 0, 0}, Vec3{0, -1, 0}, Vec3{0, 0, -1}},
		{Vec3{-1, 0, 0}, Vec3{0, 0, 1}, Vec3{0, 1, 0}},
		{Vec3{1, 0, 0}, Vec3{0, 0, -1}, Vec3{0, 1, 0}},
		{Vec3{1, 0, 0}, Vec3{0, 0, 1}, Vec3{0, -1, 0}},
		{Vec3{-1, 0, 0}, Vec3{0, 0, -1}, Vec3{0, -1, 0}},
		{Vec3{0, -1, 0}, Vec3{1, 0, 0}, Vec3{0, 0, 1}},
		{Vec3{0, 1, 0}, Vec3{-1, 0, 0}, Vec3{0, 0, 1}},
		{Vec3{0, 1, 0}, Vec3{1, 0, 0}, Vec3{0, 0, -1}},
		{Vec3{0, -1, 0}, Vec3{-1, 0, 0}, Vec3{0, 0, -1}},
		{Vec3{0, 1, 0}, Vec3{0, 0, 1}, Vec3{1, 0, 0}},
		{Vec3{0, -1, 0}, Vec3{0, 0, -1}, Vec3{1, 0, 0}},
		{Vec3{0, -1, 0}, Vec3{0, 0, 1}, Vec3{-1, 0, 0}},
		{Vec3{0, 1, 0}, Vec3{0, 0, -1}, Vec3{-1, 0, 0}},
		{Vec3{0, 0, 1}, Vec3{1, 0, 0}, Vec3{0, 1, 0}},
		{Vec3{0, 0, -1}, Vec3{-1, 0, 0}, Vec3{0, 1, 0}},
		{Vec3{0, 0, -1}, Vec3{1, 0, 0}, Vec3{0, -1, 0}},
		{Vec3{0, 0, 1}, Vec3{-1, 0, 0}, Vec3{0, -1, 0}},
		{Vec3{0, 0, -1}, Vec3{0, 1, 0}, Vec3{1, 0, 0}},
		{Vec3{0, 0, 1}, Vec3{0, -1, 0}, Vec3{1, 0, 0}},
		{Vec3{0, 0, 1}, Vec3{0, 1, 0}, Vec3{-1, 0, 0}},
		{Vec3{0, 0, -1}, Vec3{0, -1, 0}, Vec3{-1, 0, 0}},
	}
)

func findTransform(a, b []Vec3) (Vec3, Mat3, Vec3, bool) {
	t1 := vecNeg(b[0])
	t3 := a[0]
	for _, i := range rotationMatricies {
		if a[1] == vecSum(matDot(i, vecSum(b[1], t1)), t3) && a[2] == vecSum(matDot(i, vecSum(b[2], t1)), t3) {
			return t1, i, t3, true
		}
	}
	return Vec3{}, Mat3{}, Vec3{}, false
}

func markGrid(grid map[Vec3]struct{}, points []Vec3, t1 Vec3, t2 Mat3, t3 Vec3) {
	for _, i := range points {
		grid[vecSum(matDot(t2, vecSum(i, t1)), t3)] = struct{}{}
	}
}

func gridToList(g map[Vec3]struct{}) []Vec3 {
	k := make([]Vec3, 0, len(g))
	for i := range g {
		k = append(k, i)
	}
	return k
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

	var scannerlogs []*ScannerLog
	var scans []Vec3

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			scannerlogs = append(scannerlogs, NewScannerLog(scans))
			scans = nil
			continue
		}
		if strings.HasPrefix(line, "---") {
			continue
		}
		arr := strings.SplitN(line, ",", 3)
		if len(arr) != 3 {
			log.Fatalln("Invalid line")
		}
		x, err := strconv.Atoi(arr[0])
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(arr[1])
		if err != nil {
			log.Fatal(err)
		}
		z, err := strconv.Atoi(arr[2])
		if err != nil {
			log.Fatal(err)
		}
		scans = append(scans, Vec3{x, y, z})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(scans) > 0 {
		scannerlogs = append(scannerlogs, NewScannerLog(scans))
		scans = nil
	}

	grid := map[Vec3]struct{}{}
	markGrid(grid, scannerlogs[0].Scans, Vec3{0, 0, 0}, rotationMatricies[0], Vec3{0, 0, 0})
	dists := scannerlogs[0].Dists
	scannerlogs = scannerlogs[1:]

	for len(scannerlogs) != 0 {
		for i := 0; i < len(scannerlogs); i++ {
			possibleEdges := intersectDists(dists, scannerlogs[i].Dists)
			if len(possibleEdges) < 11 {
				continue
			}
			assignment := map[Vec3]Vec3{}
			if !calculateTranslation(possibleEdges, assignment, 12) {
				continue
			}
			aa, ab := get3Vec(assignment)
			t1, t2, t3, ok := findTransform(aa, ab)
			if !ok {
				continue
			}
			markGrid(grid, scannerlogs[i].Scans, t1, t2, t3)
			dists = computeDists(gridToList(grid))
			copy(scannerlogs[i:], scannerlogs[i+1:])
			scannerlogs = scannerlogs[:len(scannerlogs)-1]
			break
		}
	}

	fmt.Println(len(grid))
}
