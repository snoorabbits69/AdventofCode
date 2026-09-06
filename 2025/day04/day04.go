package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

type Point struct{ r, c int }

func parse() [][]byte {
	var result [][]byte
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			row := make([]byte, i-start)
			copy(row, data[start:i])
			result = append(result, row)
			start = i + 1
		}
	}

	if start < len(data) {
		row := make([]byte, len(data)-start)
		copy(row, data[start:])
		result = append(result, row)
	}

	return result
}

func solve() (int, int) {

	grid := parse()
	height := len(grid)
	count := 0

	dr := [8]int{-1, -1, -1, 0, 0, 1, 1, 1}
	dc := [8]int{-1, 0, 1, -1, 1, -1, 0, 1}

	removed := 0

	for r := 0; r < height; r++ {
		rowLen := len(grid[r])
		for c := 0; c < rowLen; c++ {
			if grid[r][c] != '@' {
				continue
			}

			neighbors := 0
			for i := 0; i < 8; i++ {
				rr := r + dr[i]
				cc := c + dc[i]

				if rr >= 0 && rr < height {
					if cc >= 0 && cc < len(grid[rr]) {
						if grid[rr][cc] == '@' {
							neighbors++
						}
					}
				}
			}

			if neighbors < 4 {
				count++
			}

		}
	}

	for {
		accessible := make([]Point, 0)

		for r := 0; r < height; r++ {
			w := len(grid[r])
			for c := 0; c < w; c++ {
				if grid[r][c] != '@' {
					continue
				}

				neighbors := 0
				for i := 0; i < 8; i++ {
					rr := r + dr[i]
					cc := c + dc[i]
					if rr >= 0 && rr < height && cc >= 0 && cc < len(grid[rr]) && grid[rr][cc] == '@' {
						neighbors++
					}
				}

				if neighbors < 4 {
					accessible = append(accessible, Point{r, c})
				}
			}
		}

		if len(accessible) == 0 {
			break
		}

		for _, p := range accessible {
			if p.c >= 0 && p.c < len(grid[p.r]) && grid[p.r][p.c] == '@' {
				grid[p.r][p.c] = '.'
				removed++
			}
		}
	}

	return count, removed
}

func main() {
	start := time.Now()
	const iters = 1000
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve()
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
