package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open input.txt:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var grid []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if len(grid) == 0 {
			line = strings.TrimPrefix(line, "\uFEFF")
		}
		grid = append(grid, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(grid) == 0 {
		fmt.Fprintln(os.Stderr, "empty grid")
		os.Exit(1)
	}

	dr := [8]int{-1, -1, -1, 0, 0, 1, 1, 1}
	dc := [8]int{-1, 0, 1, -1, 1, -1, 0, 1}

	height := len(grid)
	count := 0

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

	fmt.Println(count)
}
