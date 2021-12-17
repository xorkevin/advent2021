package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

type (
	BitReader struct {
		bit    int
		offset int
		buffer []byte
	}
)

func NewBitReader(buffer []byte) *BitReader {
	return &BitReader{
		bit:    0,
		offset: 0,
		buffer: buffer,
	}
}

func (r *BitReader) readBit() (byte, bool) {
	if r.offset >= len(r.buffer) {
		return 0, false
	}
	b := r.buffer[r.offset]
	k := (b >> (7 - r.bit)) & 1
	r.bit = (r.bit + 1) % 8
	if r.bit == 0 {
		r.offset++
	}
	return k, true
}

func (r *BitReader) ReadBits(a int, b []byte) (n int) {
	for i := range b {
		if n >= a {
			return
		}
		k, ok := r.readBit()
		if !ok {
			return
		}
		b[i] = k
		n++
	}
	return
}

func bitsToByte(b []byte) byte {
	var k byte = 0
	for _, i := range b {
		k = (k << 1) + i
	}
	return k
}

func bitsToInt(b []byte) int {
	k := 0
	for _, i := range b {
		k = (k << 1) + int(i)
	}
	return k
}

func nibblesToInt(b []byte) int {
	k := 0
	for _, i := range b {
		k = (k << 4) + int(i)
	}
	return k
}

func parsePacketHeader(tokens []int) (int, int, []int, bool) {
	if len(tokens) < 2 {
		return 0, 0, tokens, false
	}
	return tokens[0], tokens[1], tokens[2:], true
}

func parseSubpackets0(l, offset int, tokens []int) ([]int, int, []int, bool) {
	origOffset := offset
	origTokens := tokens
	var vals []int
	for {
		if offset-origOffset >= l {
			break
		}
		var val int
		var ok bool
		val, offset, tokens, ok = evalPacket(offset, tokens)
		if !ok {
			return nil, origOffset, origTokens, false
		}
		vals = append(vals, val)
	}
	return vals, offset, tokens, true
}

func parseSubpackets1(l, offset int, tokens []int) ([]int, int, []int, bool) {
	origOffset := offset
	origTokens := tokens
	vals := make([]int, 0, l)
	for i := 0; i < l; i++ {
		var val int
		var ok bool
		val, offset, tokens, ok = evalPacket(offset, tokens)
		if !ok {
			return nil, origOffset, origTokens, false
		}
		vals = append(vals, val)
	}
	if len(vals) < l {
		return nil, origOffset, origTokens, false
	}
	return vals, offset, tokens, true
}

func evalPacket(offset int, tokens []int) (int, int, []int, bool) {
	origOffset := offset
	origTokens := tokens
	var id int
	var ok bool
	_, id, tokens, ok = parsePacketHeader(tokens)
	if !ok {
		return 0, origOffset, origTokens, false
	}
	if id == 4 {
		if len(tokens) < 2 {
			return 0, origOffset, origTokens, false
		}
		return tokens[0], tokens[1], tokens[2:], true
	}
	if len(tokens) < 3 {
		return 0, origOffset, origTokens, false
	}
	mode := tokens[0]
	l := tokens[1]
	offset = tokens[2]
	var vals []int
	if mode == 0 {
		var ok bool
		vals, offset, tokens, ok = parseSubpackets0(l, offset, tokens[3:])
		if !ok {
			return 0, origOffset, origTokens, false
		}
	} else {
		var ok bool
		vals, offset, tokens, ok = parseSubpackets1(l, offset, tokens[3:])
		if !ok {
			return 0, origOffset, origTokens, false
		}
	}
	switch id {
	case 0:
		{
			k := 0
			for _, i := range vals {
				k += i
			}
			return k, offset, tokens, true
		}
	case 1:
		{
			k := 1
			for _, i := range vals {
				k *= i
			}
			return k, offset, tokens, true
		}
	case 2:
		{
			k := vals[0]
			for _, i := range vals {
				if i < k {
					k = i
				}
			}
			return k, offset, tokens, true
		}
	case 3:
		{
			k := vals[0]
			for _, i := range vals {
				if i > k {
					k = i
				}
			}
			return k, offset, tokens, true
		}
	case 5:
		{
			k := 0
			if vals[0] > vals[1] {
				k = 1
			}
			return k, offset, tokens, true
		}
	case 6:
		{
			k := 0
			if vals[0] < vals[1] {
				k = 1
			}
			return k, offset, tokens, true
		}
	case 7:
		{
			k := 0
			if vals[0] == vals[1] {
				k = 1
			}
			return k, offset, tokens, true
		}
	}
	return 0, origOffset, origTokens, false
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

	var bitstream []byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var err error
		bitstream, err = hex.DecodeString(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var tokens []int
	bits := NewBitReader(bitstream)
	buf := make([]byte, 15)
	bitOffset := 0
	versionSum := 0
	for {
		n := bits.ReadBits(3, buf)
		if n != 3 {
			break
		}
		bitOffset += 3
		version := bitsToByte(buf[:n])
		n = bits.ReadBits(3, buf)
		if n != 3 {
			break
		}
		bitOffset += 3
		id := bitsToByte(buf[:n])
		tokens = append(tokens, int(version))
		tokens = append(tokens, int(id))
		versionSum += int(version)
		if id == 4 {
			var nibbles []byte
			for {
				n := bits.ReadBits(5, buf)
				if n != 5 {
					break
				}
				bitOffset += 5
				nibbles = append(nibbles, bitsToByte(buf[1:n]))
				if buf[0] == 0 {
					break
				}
			}
			tokens = append(tokens, nibblesToInt(nibbles))
		} else {
			n := bits.ReadBits(1, buf)
			if n != 1 {
				break
			}
			bitOffset += 1
			mode := buf[0]
			tokens = append(tokens, int(mode))
			if mode == 0 {
				n := bits.ReadBits(15, buf)
				if n != 15 {
					break
				}
				bitOffset += 15
				tokens = append(tokens, bitsToInt(buf[:n]))
			} else {
				n := bits.ReadBits(11, buf)
				if n != 11 {
					break
				}
				bitOffset += 11
				tokens = append(tokens, bitsToInt(buf[:n]))
			}
		}
		tokens = append(tokens, bitOffset)
	}
	fmt.Println("Part 1:", versionSum)
	val, _, _, ok := evalPacket(0, tokens)
	if !ok {
		log.Fatalln("Failed eval")
	}
	fmt.Println("Part 2:", val)
}
