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
	Machine struct {
		reg   [4]int
		stdin []int
		inp   int
	}

	Arg struct {
		Imm bool
		Val int
	}

	Instr struct {
		Kind int
		Arg  [2]Arg
	}
)

const (
	instrKindInp = iota
	instrKindAdd
	instrKindMul
	instrKindDiv
	instrKindMod
	instrKindEql
)

func (i Instr) String() string {
	b := strings.Builder{}
	switch i.Kind {
	case instrKindInp:
		b.WriteString("inp")
	case instrKindAdd:
		b.WriteString("add")
	case instrKindMul:
		b.WriteString("mul")
	case instrKindDiv:
		b.WriteString("div")
	case instrKindMod:
		b.WriteString("mod")
	case instrKindEql:
		b.WriteString("eql")
	}
	for _, i := range i.Arg {
		b.WriteByte(' ')
		if i.Imm {
			b.WriteString(strconv.Itoa(i.Val))
		} else {
			b.WriteByte(byte(i.Val) + 'w')
		}
	}
	return b.String()
}

func NewMachine(stdin []int) *Machine {
	return &Machine{
		reg:   [4]int{},
		stdin: stdin,
		inp:   0,
	}
}

func (m *Machine) getArg(a Arg) int {
	if a.Imm {
		return a.Val
	}
	return m.reg[a.Val]
}

func (m *Machine) getInp() int {
	k := m.stdin[m.inp]
	m.inp++
	return k
}

func (m *Machine) Exec(instr Instr) {
	switch instr.Kind {
	case instrKindInp:
		m.reg[instr.Arg[0].Val] = m.getInp()
	case instrKindAdd:
		m.reg[instr.Arg[0].Val] = m.getArg(instr.Arg[0]) + m.getArg(instr.Arg[1])
	case instrKindMul:
		m.reg[instr.Arg[0].Val] = m.getArg(instr.Arg[0]) * m.getArg(instr.Arg[1])
	case instrKindDiv:
		m.reg[instr.Arg[0].Val] = m.getArg(instr.Arg[0]) / m.getArg(instr.Arg[1])
	case instrKindMod:
		m.reg[instr.Arg[0].Val] = m.getArg(instr.Arg[0]) % m.getArg(instr.Arg[1])
	case instrKindEql:
		k := 0
		if m.getArg(instr.Arg[0]) == m.getArg(instr.Arg[1]) {
			k = 1
		}
		m.reg[instr.Arg[0].Val] = k
	}
}

func writeArg(b *strings.Builder, a Arg) {
	if a.Imm {
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(a.Val))
		b.WriteByte(')')
	} else {
		b.WriteByte(byte(a.Val) + 'w')
	}
}

func transpile(i Instr) string {
	b := strings.Builder{}
	switch i.Kind {
	case instrKindInp:
		writeArg(&b, i.Arg[0])
		b.WriteString(" = m.getInp()")
	case instrKindAdd:
		writeArg(&b, i.Arg[0])
		b.WriteString(" = ")
		writeArg(&b, i.Arg[0])
		b.WriteString(" + ")
		writeArg(&b, i.Arg[1])
	case instrKindMul:
		writeArg(&b, i.Arg[0])
		b.WriteString(" = ")
		writeArg(&b, i.Arg[0])
		b.WriteString(" * ")
		writeArg(&b, i.Arg[1])
	case instrKindDiv:
		writeArg(&b, i.Arg[0])
		b.WriteString(" = ")
		writeArg(&b, i.Arg[0])
		b.WriteString(" / ")
		writeArg(&b, i.Arg[1])
	case instrKindMod:
		writeArg(&b, i.Arg[0])
		b.WriteString(" = ")
		writeArg(&b, i.Arg[0])
		b.WriteString(" % ")
		writeArg(&b, i.Arg[1])
	case instrKindEql:
		writeArg(&b, i.Arg[0])
		b.WriteString(" = eql(")
		writeArg(&b, i.Arg[0])
		b.WriteString(", ")
		writeArg(&b, i.Arg[1])
		b.WriteString(")")
	}
	return b.String()
}

type (
	M2 struct {
		stdin []int
		inp   int
	}
)

func NewM2(stdin []int) *M2 {
	return &M2{
		stdin: stdin,
		inp:   0,
	}
}

func eql(a, b int) int {
	if a == b {
		return 1
	}
	return 0
}

func (m *M2) getInp() int {
	k := m.stdin[m.inp]
	m.inp++
	return k
}

func (m *M2) Exec() int {
	var w, x, y, z int

	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (1)
	x = x + (10)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (2)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (1)
	x = x + (15)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (16)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (1)
	x = x + (14)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (9)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (1)
	x = x + (15)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (0)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (26)
	x = x + (-8)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (1)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (1)
	x = x + (10)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (12)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (26)
	x = x + (-16)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (6)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (26)
	x = x + (-4)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (6)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (1)
	x = x + (11)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (3)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (26)
	x = x + (-3)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (5)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (1)
	x = x + (12)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (9)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (26)
	x = x + (-7)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (3)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (26)
	x = x + (-15)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (2)
	y = y * x
	z = z + y
	w = m.getInp()
	x = x * (0)
	x = x + z
	x = x % (26)
	z = z / (26)
	x = x + (-7)
	x = eql(x, w)
	x = eql(x, (0))
	y = y * (0)
	y = y + (25)
	y = y * x
	y = y + (1)
	z = z * y
	y = y * (0)
	y = y + w
	y = y + (3)
	y = y * x
	z = z + y

	return z
}

func neq(a, b int) int {
	if a != b {
		return 1
	}
	return 0
}

func (m *M2) Exec2() int {
	var w, x, z int

	w = m.getInp()

	z = (w + 2) * neq(10, w)

	w = m.getInp()

	x = neq((z%26)+15, w)

	z = z*((25*x)+1) + (w+16)*x

	w = m.getInp()

	x = neq((z%26)+14, w)

	z = z*((25*x)+1) + ((w + 9) * x)

	w = m.getInp()

	x = neq((z%26)+15, w)

	z = z*((25*x)+1) + (w * x)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-8, w)

	z = z*((25*x)+1) + ((w + 1) * x)

	w = m.getInp()

	x = neq((z%26)+10, w)

	z = z*((25*x)+1) + ((w + 12) * x)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-16, w)

	z = z*((25*x)+1) + ((w + 6) * x)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-4, w)

	z = z*((25*x)+1) + ((w + 6) * x)

	w = m.getInp()

	x = neq((z%26)+11, w)

	z = z*((25*x)+1) + ((w + 3) * x)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-3, w)

	z = z*((25*x)+1) + ((w + 5) * x)

	w = m.getInp()

	x = neq((z%26)+12, w)

	z = z*((25*x)+1) + ((w + 9) * x)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-7, w)

	z = z*((25*x)+1) + ((w + 3) * x)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-15, w)

	z = z*((25*x)+1) + ((w + 2) * x)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-7, w)

	return z*((25*x)+1) + ((w + 3) * x)
}

func (m *M2) Exec3() int {
	var w, x, z int

	w = m.getInp()

	z = (w + 2)

	w = m.getInp()

	z = z*26 + (w + 16)

	w = m.getInp()

	z = z*26 + (w + 9)

	w = m.getInp()

	z = z*26 + w

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-8, w)

	z = z*((25*x)+1) + ((w + 1) * x)

	w = m.getInp()

	z = z*26 + (w + 12)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-16, w)

	z = z*((25*x)+1) + ((w + 6) * x)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-4, w)

	z = z*((25*x)+1) + ((w + 6) * x)

	w = m.getInp()

	z = z*26 + (w + 3)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-3, w)

	z = z*((25*x)+1) + ((w + 5) * x)

	w = m.getInp()

	z = z*26 + (w + 9)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-7, w)

	z = z*((25*x)+1) + ((w + 3) * x)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-15, w)

	z = z*((25*x)+1) + ((w + 2) * x)

	w = m.getInp()

	x = z % 26

	z = z / 26

	x = neq(x-7, w)

	return z*((25*x)+1) + ((w + 3) * x)
}

func (m *M2) Exec4() int {
	var x, z int

	{
		w1 := m.getInp()
		z = (w1 + 2)
	}
	{
		w2 := m.getInp()
		z = z*26 + (w2 + 16)
	}
	{
		w3 := m.getInp()
		z = z*26 + (w3 + 9)
	}
	{
		w4 := m.getInp()
		z = z*26 + w4
	}

	x = z % 26
	z = z / 26

	{
		w5 := m.getInp()
		x = neq(x-8, w5)
		z = z*((25*x)+1) + ((w5 + 1) * x)
	}

	{
		w6 := m.getInp()
		z = z*26 + (w6 + 12)
	}

	x = z % 26

	z = z / 26

	{
		w7 := m.getInp()

		x = neq(x-16, w7)

		z = z*((25*x)+1) + ((w7 + 6) * x)
	}

	x = z % 26

	z = z / 26

	{
		w8 := m.getInp()

		x = neq(x-4, w8)

		z = z*((25*x)+1) + ((w8 + 6) * x)
	}

	{
		w9 := m.getInp()

		z = z*26 + (w9 + 3)
	}

	x = z % 26

	z = z / 26

	{
		w10 := m.getInp()

		x = neq(x-3, w10)

		z = z*((25*x)+1) + ((w10 + 5) * x)
	}

	{
		w11 := m.getInp()

		z = z*26 + (w11 + 9)
	}

	x = z % 26

	z = z / 26

	{
		w12 := m.getInp()

		x = neq(x-7, w12)

		z = z*((25*x)+1) + ((w12 + 3) * x)
	}

	x = z % 26

	z = z / 26

	{
		w13 := m.getInp()

		x = neq(x-15, w13)

		z = z*((25*x)+1) + ((w13 + 2) * x)
	}

	x = z % 26

	z = z / 26

	w14 := m.getInp()

	x = neq(x-7, w14)

	return z*((25*x)+1) + ((w14 + 3) * x)
}

func (m *M2) Exec5() bool {
	w1 := m.getInp()
	w2 := m.getInp()
	w3 := m.getInp()

	{
		w4 := m.getInp()
		w5 := m.getInp()
		if w4-8 != w5 {
			return false
		}
	}

	{
		w6 := m.getInp()
		w7 := m.getInp()
		if w6-4 != w7 {
			return false
		}
	}

	{
		w8 := m.getInp()
		if w3+5 != w8 {
			return false
		}
	}

	{
		w9 := m.getInp()
		w10 := m.getInp()
		if w9 != w10 {
			return false
		}
	}

	{
		w11 := m.getInp()
		w12 := m.getInp()

		if w11+2 != w12 {
			return false
		}
	}

	{
		w13 := m.getInp()
		if w2+1 != w13 {
			return false
		}
	}

	{
		w14 := m.getInp()
		if w1-5 != w14 {
			return false
		}
	}

	return true
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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr := strings.Fields(scanner.Text())
		kind := instrKindInp
		switch arr[0] {
		case "inp":
			kind = instrKindInp
		case "add":
			kind = instrKindAdd
		case "mul":
			kind = instrKindMul
		case "div":
			kind = instrKindDiv
		case "mod":
			kind = instrKindMod
		case "eql":
			kind = instrKindEql
		default:
			log.Fatalln("Invalid line")
		}
		instr := Instr{
			Kind: kind,
		}
		for n, i := range arr[1:] {
			switch i {
			case "w", "x", "y", "z":
				instr.Arg[n] = Arg{
					Imm: false,
					Val: int(i[0]) - 'w',
				}
			default:
				num, err := strconv.Atoi(i)
				if err != nil {
					log.Fatal(err)
				}
				instr.Arg[n] = Arg{
					Imm: true,
					Val: num,
				}
			}
		}
		//fmt.Println(transpile(instr))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	{
		stdin := [14]int{9, 8, 4, 9, 1, 9, 5, 9, 9, 9, 7, 9, 9, 4}
		m := NewM2(stdin[:])
		if m.Exec() == 0 {
			k := 0
			for _, i := range stdin {
				k = 10*k + i
			}
			fmt.Println("Part 1:", k)
		} else {
			log.Fatal("Wrong input")
		}
	}
	{
		stdin := [14]int{6, 1, 1, 9, 1, 5, 1, 6, 1, 1, 1, 3, 2, 1}
		m := NewM2(stdin[:])
		if m.Exec() == 0 {
			k := 0
			for _, i := range stdin {
				k = 10*k + i
			}
			fmt.Println("Part 2:", k)
		} else {
			log.Fatal("Wrong input")
		}
	}
}
