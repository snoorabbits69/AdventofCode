package main

import (
	"fmt"
	"time"
)

func lookAndSay(input string) string {
	output := make([]byte, 0, len(input)*2)

	for i := 0; i < len(input); {
		count := byte(1)

		for i+int(count) < len(input) && input[i] == input[i+int(count)] {
			count++
		}

		output = append(output, '0'+count, input[i])
		i += int(count)
	}

	return string(output)
}

func solve(input string) (int, int) {
	part1 := 0
	part2 := 0

	for i := 1; i <= 50; i++ {
		input = lookAndSay(input)

		if i == 40 {
			part1 = len(input)
		}
		if i == 50 {
			part2 = len(input)
		}
	}

	return part1, part2
}

func main() {
	start := time.Now()
	const iters = 10
	input := "3113322113"
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(input)
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
