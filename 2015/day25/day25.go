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

func powMod(base, exp int) int {
	result := 1
	base %= modulus
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % modulus
		}
		base = base * base % modulus
		exp >>= 1
	}
	return result
}

func solve(row, column int) int {
	return firstCode * powMod(multiplier, index(row, column)-1) % modulus
}

func main() {
	start := time.Now()
	const iters = 10000
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
