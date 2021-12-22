package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

const (
	puzzleInput = "input.txt"
)

var (
	lineFormat = regexp.MustCompile(`^(on|off) x=(-?\d+)\.\.(-?\d+),y=(-?\d+)\.\.(-?\d+),z=(-?\d+)\.\.(-?\d+)$`)
)

type (
	Zone struct {
		prio int
		on   bool
		x1, x2,
		y1, y2,
		z1, z2 int
	}

	Boundary struct {
		prio int
		on   bool
		stop bool
		val  int
	}

	Vec3 struct {
		x, y, z int
	}
)

func getZones(zones []Zone, ids map[int]struct{}) []Zone {
	k := make([]Zone, 0, len(ids))
	for _, i := range zones {
		if _, ok := ids[i.prio]; ok {
			k = append(k, i)
		}
	}
	return k
}

func calcState(votes map[int]bool) bool {
	prio := -1
	state := false
	for k, v := range votes {
		if k > prio {
			prio = k
			state = v
		}
	}
	return state
}

func calculateOnZ(zones []Zone) int {
	bounds := make([]Boundary, 0, len(zones)*2)
	for _, i := range zones {
		bounds = append(bounds, Boundary{
			prio: i.prio,
			on:   i.on,
			stop: false,
			val:  i.z1,
		}, Boundary{
			prio: i.prio,
			on:   i.on,
			stop: true,
			val:  i.z2 + 1,
		})
	}
	sort.Slice(bounds, func(i, j int) bool {
		return bounds[i].val < bounds[j].val
	})

	count := 0
	prevZ := -999999
	inBounds := map[int]bool{}
	stateOn := false
	for _, i := range bounds {
		// interval in [prevZ, i.val)
		if i.stop {
			if stateOn {
				count += i.val - prevZ
			}
			prevZ = i.val
			delete(inBounds, i.prio)
			stateOn = calcState(inBounds)
		} else {
			if stateOn {
				count += i.val - prevZ
			}
			prevZ = i.val
			inBounds[i.prio] = i.on
			stateOn = calcState(inBounds)
		}
	}
	return count
}

func calculateOnY(zones []Zone) int {
	bounds := make([]Boundary, 0, len(zones)*2)
	for _, i := range zones {
		bounds = append(bounds, Boundary{
			prio: i.prio,
			on:   i.on,
			stop: false,
			val:  i.y1,
		}, Boundary{
			prio: i.prio,
			on:   i.on,
			stop: true,
			val:  i.y2 + 1,
		})
	}
	sort.Slice(bounds, func(i, j int) bool {
		return bounds[i].val < bounds[j].val
	})

	area := 0
	prevY := -999999 // always inclusive
	inBounds := map[int]struct{}{}
	for _, i := range bounds {
		// interval in [prevY, i.val)
		if i.stop {
			width := i.val - prevY
			if width > 0 {
				area += calculateOnZ(getZones(zones, inBounds)) * width
			}
			prevY = i.val
			delete(inBounds, i.prio)
		} else {
			width := i.val - prevY
			if width > 0 {
				area += calculateOnZ(getZones(zones, inBounds)) * width
			}
			prevY = i.val
			inBounds[i.prio] = struct{}{}
		}
	}
	return area
}

func calculateOnX(zones []Zone) int {
	bounds := make([]Boundary, 0, len(zones)*2)
	for _, i := range zones {
		bounds = append(bounds, Boundary{
			prio: i.prio,
			on:   i.on,
			stop: false,
			val:  i.x1,
		}, Boundary{
			prio: i.prio,
			on:   i.on,
			stop: true,
			val:  i.x2 + 1,
		})
	}
	sort.Slice(bounds, func(i, j int) bool {
		return bounds[i].val < bounds[j].val
	})

	volume := 0
	prevX := -999999 // always inclusive
	inBounds := map[int]struct{}{}
	for _, i := range bounds {
		if i.stop {
			// area in [prevX, i.val)
			width := i.val - prevX
			if width > 0 {
				volume += calculateOnY(getZones(zones, inBounds)) * width
			}
			prevX = i.val
			delete(inBounds, i.prio)
		} else {
			width := i.val - prevX
			if width > 0 {
				volume += calculateOnY(getZones(zones, inBounds)) * width
			}
			prevX = i.val
			inBounds[i.prio] = struct{}{}
		}
	}
	return volume
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

	var zones1 []Zone
	var zones2 []Zone

	scanner := bufio.NewScanner(file)
	prioCounter := 0
	for scanner.Scan() {
		m := lineFormat.FindStringSubmatch(scanner.Text())
		if len(m) == 0 {
			log.Fatalln("Invalid line")
		}
		on := m[1] == "on"
		x1, err := strconv.Atoi(m[2])
		if err != nil {
			log.Fatalln(err)
		}
		x2, err := strconv.Atoi(m[3])
		if err != nil {
			log.Fatalln(err)
		}
		y1, err := strconv.Atoi(m[4])
		if err != nil {
			log.Fatalln(err)
		}
		y2, err := strconv.Atoi(m[5])
		if err != nil {
			log.Fatalln(err)
		}
		z1, err := strconv.Atoi(m[6])
		if err != nil {
			log.Fatalln(err)
		}
		z2, err := strconv.Atoi(m[7])
		if err != nil {
			log.Fatalln(err)
		}
		if x1 > 50 || x1 < -50 {
			zones2 = append(zones2, Zone{prioCounter, on, x1, x2, y1, y2, z1, z2})
		} else {
			zones1 = append(zones1, Zone{prioCounter, on, x1, x2, y1, y2, z1, z2})
		}
		prioCounter++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	part1 := calculateOnX(zones1)
	fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", part1+calculateOnX(zones2))
}
