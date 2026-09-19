package main

import (
	"fmt"
	"time"
)

func solve(target, mult, maxVisits int) int {

	limit := target / mult
	houses := make([]int, limit+1)

	for elf := 1; elf <= limit; elf++ {
		visits := 0
		for h := elf; h <= limit; h += elf {
			houses[h] += elf * mult
			visits++
			if maxVisits > 0 && visits >= maxVisits {
				break
			}
		}
	}

	for h := 1; h <= limit; h++ {
		if houses[h] >= target {
			return h
		}
	}
	return -1
}

func main() {
	start := time.Now()
	const iters = 10
	const target = 29000000
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1 = solve(target, 10, 0)
		part2 = solve(target, 11, 50)
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
