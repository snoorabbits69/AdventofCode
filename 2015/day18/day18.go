package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

const (
	size  = 100
	width = size + 2
	steps = 100
)

type Grid [width * width]uint8

func parse() Grid {
	var g Grid
	x, y := 1, 1
	for _, b := range data {
		switch b {
		case '#':
			g[y*width+x] = 1
			x++
		case '.':
			x++
		case '\n':
			y++
			x = 1
		}
	}
	return g
}

func setCorners(g *Grid) {
	g[1*width+1] = 1
	g[1*width+size] = 1
	g[size*width+1] = 1
	g[size*width+size] = 1
}

func step(cur, next *Grid) {
	for y := 1; y <= size; y++ {
		for x := 1; x <= size; x++ {
			i := y*width + x
			n := cur[i-width-1] + cur[i-width] + cur[i-width+1] +
				cur[i-1] + cur[i+1] +
				cur[i+width-1] + cur[i+width] + cur[i+width+1]

			if n == 3 || (n == 2 && cur[i] == 1) {
				next[i] = 1
			} else {
				next[i] = 0
			}
		}
	}
}

func run(g Grid, stuck bool) int {
	cur, next := g, g
	if stuck {
		setCorners(&cur)
	}
	for s := 0; s < steps; s++ {
		step(&cur, &next)
		if stuck {
			setCorners(&next)
		}
		cur, next = next, cur
	}

	total := 0
	for _, v := range cur {
		total += int(v)
	}
	return total
}

func main() {
	start := time.Now()
	const iters = 100
	g := parse()
	var part1, part2 int
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
