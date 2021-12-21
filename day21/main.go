package main

import "fmt"

type (
	DetDie struct {
		k     int
		rolls int
	}
)

func (d *DetDie) Roll() int {
	d.rolls++
	v := d.k + 1
	d.k = v % 100
	return v
}

func turnRolls(d *DetDie) int {
	return d.Roll() + d.Roll() + d.Roll()
}

type (
	Point struct {
		x, y, z, a, b int
	}

	Vec2 struct {
		x, y int
	}
)

var (
	memo  = map[Point]Vec2{}
	steps = []Vec2{
		{3, 1},
		{4, 3},
		{5, 6},
		{6, 7},
		{7, 6},
		{8, 3},
		{9, 1},
	}
)

func getUniverses(s1, s2 int, turn int, a, b int) Vec2 {
	if s1 > 20 {
		return Vec2{1, 0}
	}
	if s2 > 20 {
		return Vec2{0, 1}
	}
	if v, ok := memo[Point{s1, s2, turn, a, b}]; ok {
		return v
	}
	count := Vec2{0, 0}
	if turn == 0 {
		for _, i := range steps {
			a1 := (a + i.x) % 10
			k := getUniverses(s1+a1+1, s2, 1, a1, b)
			count.x += k.x * i.y
			count.y += k.y * i.y
		}
	} else {
		for _, i := range steps {
			b1 := (b + i.x) % 10
			k := getUniverses(s1, s2+b1+1, 0, a, b1)
			count.x += k.x * i.y
			count.y += k.y * i.y
		}
	}
	memo[Point{s1, s2, turn, a, b}] = count
	return count
}

func main() {
	//Player 1 starting position: 7
	//Player 2 starting position: 3
	{
		s1 := 0
		s2 := 0
		p1 := 6
		p2 := 2
		d := &DetDie{}
		for {
			k := turnRolls(d)
			p1 = (p1 + k) % 10
			s1 += p1 + 1
			if s1 >= 1000 {
				break
			}
			k = turnRolls(d)
			p2 = (p2 + k) % 10
			s2 += p2 + 1
			if s2 >= 1000 {
				break
			}
		}
		l := s1
		if s2 < s1 {
			l = s2
		}
		fmt.Println("Part 1:", l*d.rolls)
	}
	{
		// 3: 1
		// 3 = 1, 1, 1

		// 4: 3
		// 4 = 1, 1, 2
		// 4 = 1, 2, 1
		// 4 = 2, 1, 1

		// 5: 6
		// 5 = 1, 1, 3
		// 5 = 1, 2, 2
		// 5 = 1, 3, 1
		// 5 = 2, 1, 2
		// 5 = 2, 2, 1
		// 5 = 3, 1, 1

		// 6: 7
		// 6 = 1, 2, 3
		// 6 = 1, 3, 2
		// 6 = 2, 1, 3
		// 6 = 2, 2, 2
		// 6 = 2, 3, 1
		// 6 = 3, 1, 2
		// 6 = 3, 2, 1

		// 7: 6
		// 7 = 1, 3, 3
		// 7 = 2, 2, 3
		// 7 = 2, 3, 2
		// 7 = 3, 1, 3
		// 7 = 3, 2, 2
		// 7 = 3, 3, 1

		// 8: 3
		// 8 = 2, 3, 3
		// 8 = 3, 2, 3
		// 8 = 3, 3, 2

		// 9: 1
		// 9 = 3, 3, 3

		k := getUniverses(0, 0, 0, 6, 2)
		max := k.x
		if k.y > max {
			max = k.y
		}
		fmt.Println("Part 2:", max)
	}
}
