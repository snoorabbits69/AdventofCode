package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

type Point struct {
	x, y int
}

func parseGrid(data []byte) ([][]byte, Point) {
	var grid [][]byte
	var line []byte
	var start Point

	for _, b := range data {
		switch b {
		case '\n':
			grid = append(grid, line)
			line = nil
		case '\r':
			// skip
		default:
			line = append(line, b)
		}
	}
	if len(line) > 0 {
		grid = append(grid, line)
	}

	for y, row := range grid {
		for x, c := range row {
			if c == 'S' {
				start = Point{x, y}
			}
		}
	}
	return grid, start
}

func inBounds(grid [][]byte, p Point) bool {
	return p.y >= 0 && p.y < len(grid) && p.x >= 0 && p.x < len(grid[p.y])
}

func solve(grid [][]byte, start Point) (int, int) {
	splits := 0
	queue := []Point{start}
	visited := map[Point]bool{start: true}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if grid[cur.y][cur.x] == '^' {
			splits++
			for _, next := range []Point{{cur.x - 1, cur.y}, {cur.x + 1, cur.y}} {
				if inBounds(grid, next) && !visited[next] {
					visited[next] = true
					queue = append(queue, next)
				}
			}
		} else {
			down := Point{cur.x, cur.y + 1}
			if inBounds(grid, down) && !visited[down] {
				visited[down] = true
				queue = append(queue, down)
			}
		}
	}
	total := 0
	current := map[Point]int{start: 1}

	for len(current) > 0 {
		next := map[Point]int{}
		for p, count := range current {
			if grid[p.y][p.x] == '^' {
				for _, n := range []Point{{p.x - 1, p.y}, {p.x + 1, p.y}} {
					if inBounds(grid, n) {
						next[n] += count
					} else {
						total += count
					}
				}
			} else {
				d := Point{p.x, p.y + 1}
				if inBounds(grid, d) {
					next[d] += count
				} else {
					total += count
				}
			}
		}
		current = next
	}

	return splits, total
}

func main() {
	start := time.Now()
	const iters = 100
	grid, startPt := parseGrid(data)
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(grid, startPt)
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
