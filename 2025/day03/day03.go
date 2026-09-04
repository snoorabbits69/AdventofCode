package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

func parse() [][]byte {
	var result [][]byte
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {

			result = append(result, data[start:i])

			start = i + 1
		}
	}

	if start < len(data) {
		end := len(data)
		result = append(result, data[start:end])
	}

	return result
}

func solve(input [][]byte, N int) int64 {
	var sum int64

	for _, bank := range input {
		batteries := make([]byte, N)

		end := len(bank) - N
		copy(batteries, bank[end:])

		for i := end - 1; i >= 0; i-- {
			next := bank[i]

			for j := 0; j < N; j++ {
				if next < batteries[j] {
					break
				}
				next, batteries[j] = batteries[j], next
			}
		}

		var value int64
		for _, b := range batteries {
			value = value*10 + int64(b-'0')
		}

		sum += value
	}

	return sum
}

func main() {
	start := time.Now()
	const iters = 10000
	lines := parse()
	var part1, part2 int64
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(lines, 2), solve(lines, 12)
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
