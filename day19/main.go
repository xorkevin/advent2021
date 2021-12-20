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
	puzzleInput = "input2.txt"
)

type (
	Vec3 struct {
		x, y, z int
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

func NewScannerLog(scans []Vec3) *ScannerLog {
	dists := map[Vec3][]Edge{}
	l := len(scans)
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			k := orderedDist(scans[i], scans[j])
			dists[k] = append(dists[k], Edge{scans[i], scans[j]})
		}
	}
	return &ScannerLog{
		Scans: scans,
		Dists: dists,
	}
}

func (s *ScannerLog) IntersectDists(o *ScannerLog) ([]PossibleEdges, int) {
	var dists []PossibleEdges
	count := 0
	for k, v := range o.Dists {
		m := min(len(s.Dists[k]), len(v))
		count += m
		if m > 0 {
			a := make([]Edge, len(s.Dists[k]))
			copy(a, s.Dists[k])
			b := make([]Edge, len(v))
			copy(b, v)
			dists = append(dists, PossibleEdges{
				a: a,
				b: b,
			})
		}
	}
	return dists, count
}

func calculateTranslation(possibleEdges []PossibleEdges, pointTranslation map[Vec3]Vec3) bool {
	if len(possibleEdges) == 0 {
		return len(pointTranslation) > 2
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
				if ok := calculateTranslation(possibleEdges[1:], pointTranslation); ok {
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
				if ok := calculateTranslation(possibleEdges[1:], pointTranslation); ok {
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
			if ok := calculateTranslation(possibleEdges[1:], pointTranslation); ok {
				return true
			}
		}
	}
	return false
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

	for i := 0; i < len(scannerlogs); i++ {
		for j := i + 1; j < len(scannerlogs); j++ {
			possibleEdges, count := scannerlogs[i].IntersectDists(scannerlogs[j])
			if count < 11 {
				continue
			}
			assignment := map[Vec3]Vec3{}
			if !calculateTranslation(possibleEdges, assignment) {
				continue
			}
			if len(assignment) < 12 {
				continue
			}
			fmt.Println(i, j, assignment)
		}
	}
}
