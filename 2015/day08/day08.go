package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

func solve() (int, int) {
	n := len(data)
	i := 0

	total_code := n
	total_char := 0
	extra_code := 0
	line_started := 0

	for i < n {
		c := data[i]

		if c == '\n' {
			total_code--

			if line_started == 1 {
				extra_code += 2
				line_started = 0
			}

			i++
			continue
		}

		line_started = 1
		if c == '"' {
			extra_code++

			i++
			continue
		}

		if c == '\\' && i+1 < n {
			extra_code++
			switch data[i+1] {
			case '\\', '"':
				extra_code++
				total_char++
				i += 2
				continue
			case 'x':
				if i+3 < n {
					total_char++
					i += 4
					continue
				}
			}
		}

		total_char++
		i++
	}

	if line_started == 1 {
		extra_code += 2
	}

	return total_code - total_char, extra_code
}

func main() {
	start := time.Now()
	const iters = 10000
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
