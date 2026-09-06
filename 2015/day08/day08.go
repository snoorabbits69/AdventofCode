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
	line_started :=0

	for i < n {
		c := data[i]

		if c == '\n' {
			total_code--

			if line_started==1 {
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

	if line_started==1 {
		extra_code += 2
	}

	return total_code - total_char, extra_code
}

func main() {
	start := time.Now()
	var part1,part2 int

	for i := 0; i < 100; i++ {
		part1,part2= solve()
	}

	const iters = 10000
	benchStart := time.Now()
	for i := 0; i < iters; i++ {
	part1,part2= solve()
	}
	benchTime := time.Since(benchStart)
	totalTime := time.Since(start)

	fmt.Println("solution of part1:", part1)
	fmt.Println("solution of part2:", part2)
	fmt.Println()

	fmt.Printf("Total bench time:  %d microseconds\n", benchTime.Microseconds())
	avgUs := float64(benchTime.Nanoseconds()) / 1_000.0/float64(iters)
	fmt.Printf("Average per iter:  %.4f microseconds\n", avgUs)
	fmt.Printf("Total time:        %d microseconds\n", totalTime.Microseconds())
}
