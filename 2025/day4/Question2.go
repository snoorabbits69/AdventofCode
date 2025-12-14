package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Point struct{ r, c int }

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var grid [][]byte
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if len(grid) == 0 {
			line = strings.TrimPrefix(line, "\uFEFF")
		}
		grid = append(grid, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if len(grid) == 0 {
		fmt.Println(0)
		return
	}

	dr := [8]int{-1, -1, -1, 0, 0, 1, 1, 1}
	dc := [8]int{-1, 0, 1, -1, 1, -1, 0, 1}

	removed := 0
	h := len(grid)

	for {
		accessible := make([]Point, 0)

		for r := 0; r < h; r++ {
			w := len(grid[r])
			for c := 0; c < w; c++ {
				if grid[r][c] != '@' {
					continue
				}

				neighbors := 0
				for i := 0; i < 8; i++ {
					rr := r + dr[i]
					cc := c + dc[i]
					if rr >= 0 && rr < h && cc >= 0 && cc < len(grid[rr]) && grid[rr][cc] == '@' {
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

	fmt.Println(removed)
}
