package main

import (
	"fmt"
	"time"
)

const (
	firstCode  = 20151125
	multiplier = 252533
	modulus    = 33554393
)

func index(row, column int) int {
	d := row + column - 1
	return d*(d-1)/2 + column
}

func solve(row, column int) int {
	result := firstCode
	for i := 1; i < index(row, column); i++ {
		result = result * multiplier % modulus
	}
	return result
}

func main() {
	start := time.Now()
	const iters = 10
	row, column := 2947, 3029
	var part1 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1 = solve(row, column)
	}
	processTime := time.Since(processStart)
	totalTime := time.Since(start)
	elapsedUs := float64(processTime.Nanoseconds()) / 1000.0
	fmt.Printf("Total: %.2f microseconds\n", elapsedUs)
	fmt.Printf("Average: %.4f microseconds\n", elapsedUs/float64(iters))
	fmt.Printf("Total time:        %dns\n", totalTime.Nanoseconds())
	fmt.Println("part1", part1)
}
