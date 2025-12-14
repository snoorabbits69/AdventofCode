package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 50*1024*1024)

	var grid []string
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	h := len(grid)
	w := len(grid[0])

	start := Point{}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if grid[y][x] == 'S' {
				start = Point{x, y}
			}
		}
	}

	inBounds := func(p Point) bool {
		return p.x >= 0 && p.x < w && p.y >= 0 && p.y < h
	}

	current := map[Point]int{}
	current[start] = 1

	totalTimelines := 0

	for len(current) > 0 {
		next := map[Point]int{}

		for p, count := range current {
			cell := grid[p.y][p.x]

			if cell == '^' {
				left := Point{p.x - 1, p.y}
				right := Point{p.x + 1, p.y}

				if inBounds(left) {
					next[left] += count
				} else {
					totalTimelines += count
				}

				if inBounds(right) {
					next[right] += count
				} else {
					totalTimelines += count
				}

			} else {
				down := Point{p.x, p.y + 1}
				if inBounds(down) {
					next[down] += count
				} else {
					totalTimelines += count
				}
			}
		}

		current = next
	}

	fmt.Println(totalTimelines)
}
