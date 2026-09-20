package main

import (
	"fmt"
	"time"
)

func solve(target, mult, maxVisits int32) int32 {
	need := (target-1)/mult + 1
	limit := need/16 + 1

	for {
		if limit > need {
			limit = need
		}

		houses := make([]int32, limit+1)

		for elf := int32(1); elf <= limit; elf++ {
			end := limit
			if maxVisits > 0 && elf <= end/maxVisits {
				end = elf * maxVisits
			}
			for h := elf; ; h += elf {
				houses[h] += elf
				if h > end-elf {
					break
				}
			}
		}

		for h := int32(1); h <= limit; h++ {
			if houses[h] >= need {
				return h
			}
		}

		if limit == need {
			return -1
		}
		if limit > need/2 {
			limit = need
		} else {
			limit *= 2
		}
	}
}

func main() {
	start := time.Now()
	const iters = 10
	const target int32 = 29000000
	var part1, part2 int32

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
