package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func solution(input string, prefix string, start int) int {
	n := start
	for {
		newinput := input + strconv.Itoa(n)
		hash := md5.Sum([]byte(newinput))
		hashHex := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashHex, prefix) {
			break
		}
		n++
	}
	return n
}

func solve(input string) (int, int) {
	solution1 := solution(input, "00000", 0)
	solution2 := solution(input, "000000", solution1+1)
	return solution1, solution2
}

func main() {
	input := "bgvyzdsv"

	start := time.Now()

	var part1, part2 int
	for i := 0; i < 10; i++ {
		part1, part2 = solve(input)
	}


	const iters = 10000
	processStart := time.Now()

	for i := 0; i < iters; i++ {
		part1, part2 = solve(input)
	}

	processTime := time.Since(processStart)
	totalTime := time.Since(start)

	fmt.Println("Answer of part 1:", part1)
	fmt.Println("Answer of part 2:", part2)
	fmt.Println()

	elapsedUs := float64(processTime.Nanoseconds()) / 1_000.0
	fmt.Printf("Total:             %.2f microseconds\n", elapsedUs)
	fmt.Printf("Average:           %.4f microseconds\n", elapsedUs/float64(iters))

	fmt.Printf("Total time:        %d µs\n", totalTime.Microseconds())
}
