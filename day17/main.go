package main

import "fmt"

type (
	Vec2 struct {
		x int
		y int
	}
)

func main() {
	//target area: x=25..67, y=-260..-200
	x1, x2, y1, y2 := 25, 67, -260, -200

	maxvy := 0
	{
		maxy := 0
		lower := 256
		upper := 264
		for upper >= lower {
			vy := lower + (upper-lower)/2
			k, ok := simulate(Vec2{0, 0}, Vec2{7, vy}, x1, x2, y1, y2)
			if !ok {
				upper = vy - 1
			} else {
				if k > maxy {
					maxy = k
					maxvy = vy
				}
				lower = vy + 1
			}
		}
		fmt.Println("Part 1:", maxy)
	}
	{
		count := 0
		// 6th triangular number is 21
		for x := 7; x < x2+1; x++ {
			for y := y1 - 1; y < maxvy+1; y++ {
				_, ok := simulate(Vec2{0, 0}, Vec2{x, y}, x1, x2, y1, y2)
				if ok {
					count++
				}
			}
		}
		fmt.Println("Part 2:", count)
	}
}

func simulate(p Vec2, v Vec2, x1, x2, y1, y2 int) (int, bool) {
	maxy := 0
	for {
		p.x += v.x
		p.y += v.y

		if v.x > 0 {
			v.x--
		} else if v.x < 0 {
			v.x++
		}
		v.y--

		if p.y > maxy {
			maxy = p.y
		}

		if inTarget(p, x1, x2, y1, y2) {
			return maxy, true
		}
		if p.y < y1 {
			return 0, false
		}
		if v.x == 0 && (p.x < x1 || p.x > x2) {
			return 0, false
		}
	}
}

func inTarget(p Vec2, x1, x2, y1, y2 int) bool {
	return p.x >= x1 && p.x <= x2 && p.y >= y1 && p.y <= y2
}
