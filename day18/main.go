package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"xorkevin.dev/gnom"
)

const (
	puzzleInput = "input.txt"
)

const (
	tokenKindDefault = iota
	tokenKindEOF
	tokenKindLparen
	tokenKindRparen
	tokenKindComma
	tokenKindNum
)

var (
	ErrParse = errors.New("Parse error")
)

type (
	Pair struct {
		val int
		lhs *Pair
		rhs *Pair
	}
)

func (p Pair) IsLiteral() bool {
	return p.lhs == nil
}

func (p Pair) buildString(b *strings.Builder) {
	if p.IsLiteral() {
		b.WriteString(strconv.Itoa(p.val))
		return
	}
	b.WriteByte('[')
	p.lhs.buildString(b)
	b.WriteByte(',')
	p.rhs.buildString(b)
	b.WriteByte(']')
}

func (p Pair) String() string {
	b := strings.Builder{}
	p.buildString(&b)
	return b.String()
}

func parsePairs(tokens []gnom.Token) (*Pair, []gnom.Token, error) {
	if len(tokens) == 0 {
		return nil, nil, ErrParse
	}
	top := tokens[0]
	switch top.Kind() {
	case tokenKindNum:
		{
			num, err := strconv.Atoi(top.Val())
			if err != nil {
				return nil, nil, ErrParse
			}
			return &Pair{
				val: num,
			}, tokens[1:], nil
		}
	case tokenKindLparen:
		{
			var lhs *Pair
			var err error
			lhs, tokens, err = parsePairs(tokens[1:])
			if err != nil {
				return nil, nil, err
			}
			if len(tokens) == 0 {
				return nil, nil, ErrParse
			}
			if tokens[0].Kind() != tokenKindComma {
				return nil, nil, ErrParse
			}
			var rhs *Pair
			rhs, tokens, err = parsePairs(tokens[1:])
			if err != nil {
				return nil, nil, ErrParse
			}
			if len(tokens) == 0 {
				return nil, nil, ErrParse
			}
			if tokens[0].Kind() != tokenKindRparen {
				return nil, nil, ErrParse
			}
			return &Pair{
				lhs: lhs,
				rhs: rhs,
			}, tokens[1:], nil
		}
	default:
		return nil, nil, ErrParse
	}
}

func (p *Pair) addLeft(v int) {
	if p.IsLiteral() {
		p.val += v
		return
	}
	p.lhs.addLeft(v)
}

func (p *Pair) addRight(v int) {
	if p.IsLiteral() {
		p.val += v
		return
	}
	p.rhs.addRight(v)
}

func (p *Pair) explode(depth int) (int, int, bool) {
	if p.IsLiteral() {
		return 0, 0, false
	}
	if depth > 3 {
		if !p.lhs.IsLiteral() || !p.rhs.IsLiteral() {
			log.Fatalln("Invalid reduction state")
		}
		l := p.lhs.val
		r := p.rhs.val
		p.val = 0
		p.lhs = nil
		p.rhs = nil
		return l, r, true
	}
	if l, r, ok := p.lhs.explode(depth + 1); ok {
		if r != 0 {
			p.rhs.addLeft(r)
		}
		return l, 0, true
	}
	if l, r, ok := p.rhs.explode(depth + 1); ok {
		if l != 0 {
			p.lhs.addRight(l)
		}
		return 0, r, true
	}
	return 0, 0, false
}

func (p *Pair) split() bool {
	if p.IsLiteral() {
		if p.val > 9 {
			l := p.val / 2
			r := p.val - l
			p.val = 0
			p.lhs = &Pair{
				val: l,
			}
			p.rhs = &Pair{
				val: r,
			}
			return true
		}
		return false
	}
	if ok := p.lhs.split(); ok {
		return true
	}
	return p.rhs.split()
}

func (p *Pair) reduceStep() bool {
	if _, _, ok := p.explode(0); ok {
		return true
	}
	return p.split()
}

func (p *Pair) Reduce() {
	for p.reduceStep() {
	}
}

func (p Pair) Magnitude() int {
	if p.IsLiteral() {
		return p.val
	}
	return 3*p.lhs.Magnitude() + 2*p.rhs.Magnitude()
}

func (p Pair) Clone() *Pair {
	if p.IsLiteral() {
		return &Pair{
			val: p.val,
		}
	}
	return &Pair{
		lhs: p.lhs.Clone(),
		rhs: p.rhs.Clone(),
	}
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

	dfa := gnom.NewDfa(tokenKindDefault)
	dfaNum := gnom.NewDfa(tokenKindNum)
	dfa.AddDfa([]rune("0123456789"), dfaNum)
	dfaNum.AddDfa([]rune("0123456789"), dfaNum)
	dfa.AddPath([]rune("["), tokenKindLparen, tokenKindDefault)
	dfa.AddPath([]rune("]"), tokenKindRparen, tokenKindDefault)
	dfa.AddPath([]rune(","), tokenKindComma, tokenKindDefault)
	lexer := gnom.NewDfaLexer(dfa, tokenKindDefault, tokenKindEOF, map[int]struct{}{})

	var nums []*Pair
	var root *Pair
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens, err := lexer.Tokenize([]rune(scanner.Text()))
		if err != nil {
			log.Fatal(err)
		}
		pair, _, err := parsePairs(tokens)
		nums = append(nums, pair.Clone())
		if root == nil {
			root = pair
		} else {
			root = &Pair{
				lhs: root,
				rhs: pair,
			}
		}
		root.Reduce()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Part 1:", root.Magnitude())

	maxmag := 0
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			{
				k := &Pair{
					lhs: nums[i].Clone(),
					rhs: nums[j].Clone(),
				}
				k.Reduce()
				maxmag = max(maxmag, k.Magnitude())
			}
			{
				k := &Pair{
					lhs: nums[j].Clone(),
					rhs: nums[i].Clone(),
				}
				k.Reduce()
				maxmag = max(maxmag, k.Magnitude())
			}
		}
	}
	fmt.Println("Part 2:", maxmag)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
