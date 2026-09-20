package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

var (
	cells []uint8
	stack []int32
	rolls []int32
)

func solve() (int, int) {
	input := bytes.TrimRight(data, "\r\n")
	w := bytes.IndexByte(input, '\n')
	if w < 0 {
		w = len(input)
	}
	if w > 0 && input[w-1] == '\r' {
		w--
	}
	h := bytes.Count(input, []byte{'\n'}) + 1

	stride := w + 2
	size := (h + 2) * stride
	if cap(cells) < size {
		cells = make([]uint8, size)
	}
	cells = cells[:size]
	for i := range cells {
		cells[i] = 128
	}

	offs := [8]int{-stride - 1, -stride, -stride + 1, -1, 1, stride - 1, stride, stride + 1}

	rolls = rolls[:0]
	pos := 0
	for y := 0; y < h; y++ {
		line := input[pos:]
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line = line[:i]
			pos += i + 1
		}
		if n := len(line); n > w {
			line = line[:w]
		}
		base := (y+1)*stride + 1
		for x, ch := range line {
			if ch == '@' {
				cells[base+x] = 0
				rolls = append(rolls, int32(base+x))
			}
		}
	}

	for _, p := range rolls {
		for _, o := range offs {
			cells[int(p)+o]++
		}
	}

	stack = stack[:0]
	for _, p := range rolls {
		if cells[p] < 4 {
			stack = append(stack, p)
		}
	}
	part1 := len(stack)

	removed := 0
	for len(stack) > 0 {
		p := int(stack[len(stack)-1])
		stack = stack[:len(stack)-1]
		removed++
		for _, o := range offs {
			n := p + o
			if cells[n] == 4 {
				stack = append(stack, int32(n))
			}
			cells[n]--
		}
	}
	return part1, removed
}

func main() {
	const iters = 1000
	var part1, part2 int
	start := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve()
	}
	elapsed := time.Since(start)
	us := float64(elapsed.Nanoseconds()) / 1000.0
	fmt.Printf("Total: %.2f microseconds\n", us)
	fmt.Printf("Average: %.4f microseconds\n", us/iters)
	fmt.Println("part1", part1)
	fmt.Println("part2", part2)
}
