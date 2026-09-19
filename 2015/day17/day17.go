package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

func parseNumber(i *int) int64 {
	var v int64

	for *i < len(data) {
		c := data[*i]

		if c < '0' || c > '9' {
			break
		}

		v = v*10 + int64(c-'0')
		*i = *i + 1
	}

	return v
}

func parse() []int64 {
	var containers []int64
	i := 0

	for i < len(data) {
		c := data[i]

		if c < '0' || c > '9' {
			i++
			continue
		}

		containers = append(containers, parseNumber(&i))
	}

	return containers
}

func part1(containers []int64, goal int64) int64 {
	ways := make([]int64, goal+1)
	ways[0] = 1

	for _, item := range containers {

		for i := goal; i >= item; i-- {
			ways[i] += ways[i-item]
		}
	}

	return ways[goal]
}

func part2(containers []int64, goal int64) int64 {
	ways := make([]int64, goal+1)
	ways[0] = 1

	const maxInt64 = int64(1) << 62
	minimum := make([]int64, goal+1)
	for i := range minimum {
		minimum[i] = maxInt64
	}
	minimum[0] = 0

	for _, item := range containers {
		for i := goal; i >= item; i-- {
			take := minimum[i-item]
			if take != maxInt64 {
				take++
			}
			notTake := minimum[i]

			switch {
			case take < notTake:

				ways[i] = ways[i-item]
				minimum[i] = take
			case take == notTake:

				ways[i] += ways[i-item]

			}
		}
	}

	return ways[goal]
}

func main() {
	start := time.Now()
	const iters = 10000
	containers := parse()
	const goal = 150
	part1Fn, part2Fn := part1, part2
	var part1, part2 int64
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1 = part1Fn(containers, goal)
		part2 = part2Fn(containers, goal)
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
