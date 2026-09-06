package main

import (
	"fmt"
	"os"
	"time"
)

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

	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	var floor, firstBasementPos int
	for i := 0; i < 100; i++ {
		floor, firstBasementPos = solve(data)
	}

	const iters = 10000
	processStart := time.Now()

	for i := 0; i < iters; i++ {
		floor, firstBasementPos = solve(data)
	}

	processTime := time.Since(processStart)
	totalTime := time.Since(start)

	elapsedUs := float64(processTime.Nanoseconds()) / 1000.0

	fmt.Println("Solution of part1:", floor)
	fmt.Println("Solution of part2:", firstBasementPos)

	fmt.Printf("Total: %.2f microseconds\n", elapsedUs)
	fmt.Printf("Average: %.4f microseconds\n", elapsedUs/float64(iters))
	fmt.Printf("Total time:        %dns\n", totalTime.Nanoseconds())
}
