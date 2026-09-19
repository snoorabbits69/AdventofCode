package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

var delta = [256]int{
	'(': 1,
	')': -1,
}

func solve(data []byte) (int, int) {
	floor := 0
	firstBasementPos := -1

	for i, c := range data {
		floor += delta[c]

		if floor == -1 && firstBasementPos == -1 {
			firstBasementPos = i + 1
		}
	}

	return floor, firstBasementPos
}

func main() {
	start := time.Now()
	const iters = 10000
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(data)
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
