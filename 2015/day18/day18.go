package main

import (
	_ "embed"
	"fmt"
	"math/bits"
	"time"
)

//go:embed input.txt
var data []byte

const (
	size  int32 = 100
	steps int32 = 100

	mask        uint64 = (1 << 50) - 1
	leftCorner  uint64 = 1 << 49
	rightCorner uint64 = 1 << 0
)

type Row struct {
	even uint64
	odd  uint64
}

type Grid [101]Row

func parse() Grid {
	var g Grid
	var row Row
	var x, y int32

	for _, b := range data {
		switch b {
		case '#', '.':
			bit := uint64(b & 1)
			if x&1 == 0 {
				row.even = row.even<<1 | bit
			} else {
				row.odd = row.odd<<1 | bit
			}
			x++
		case '\n':
			if x > 0 && y < size {
				g[y] = row
				y++
			}
			row = Row{}
			x = 0
		}
	}

	if x > 0 && y < size {
		g[y] = row
	}
	return g
}

func setCorners(g *Grid) {
	g[0].even |= leftCorner
	g[0].odd |= rightCorner
	g[99].even |= leftCorner
	g[99].odd |= rightCorner
}

func halfAdder(a, b uint64) (uint64, uint64) {
	return a & b, a ^ b
}

func fullAdder(a, b, c uint64) (uint64, uint64) {
	return a&b | c&(a^b), a ^ b ^ c
}

func step(grid *Grid) {
	var above Row
	var row Row
	below := grid[0]

	for y := int32(0); y < size; y++ {

		above, row, below = row, below, grid[y+1]

		j, k := fullAdder(above.odd, row.odd, below.odd)
		l, m := halfAdder(above.even, below.even)
		n, o := j>>1, k>>1

		p, q := fullAdder(j, l, n)
		r, s := fullAdder(k, m, o)
		grid[y].even = (grid[y].even | s) & (q ^ r) &^ p & mask

		j, k = fullAdder(above.even, row.even, below.even)
		l, m = halfAdder(above.odd, below.odd)
		n, o = j<<1, k<<1

		p, q = fullAdder(j, l, n)
		r, s = fullAdder(k, m, o)
		grid[y].odd = (grid[y].odd | s) & (q ^ r) &^ p & mask
	}
}

func run(g Grid, stuck bool) int32 {

	cur := g
	if stuck {
		setCorners(&cur)
	}

	for s := int32(0); s < steps; s++ {
		step(&cur)
		if stuck {
			setCorners(&cur)
		}
	}

	var total int32
	for i := int32(0); i < size; i++ {
		total += int32(bits.OnesCount64(cur[i].even) + bits.OnesCount64(cur[i].odd))
	}
	return total
}

func main() {
	start := time.Now()
	const iters = 100
	g := parse()
	var part1, part2 int32
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1 = run(g, false)
		part2 = run(g, true)
	}
	processTime := time.Since(processStart)
	totalTime := time.Since(start)
	elapsedUs := float64(processTime.Nanoseconds()) / 1000.0
	fmt.Printf("Total: %.2f microseconds\n", elapsedUs)
	fmt.Printf("Average: %.4f microseconds\n", elapsedUs/float64(iters))
	fmt.Printf("Total time:        %dns\n", totalTime.Nanoseconds())
	fmt.Println("part1", part1)
	fmt.Println("part2", part2)
}
