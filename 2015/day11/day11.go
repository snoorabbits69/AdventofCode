package main

import (
	"fmt"
	"time"
)

func charIncrement(char rune) rune {
	if char == 'z' {
		return 'a'
	}
	return char + 1
}
func incrementPassword(pw []rune) {
	for i := len(pw) - 1; i >= 0; i-- {
		if pw[i] == 'z' {
			pw[i] = 'a'
		} else {
			pw[i] = charIncrement(pw[i])
			return
		}
	}
}

func hasIncreasingStraight(pw []rune) bool {
	for i := 0; i < len(pw)-2; i++ {
		if pw[i]+1 == pw[i+1] && pw[i+1]+1 == pw[i+2] {
			return true
		}
	}
	return false
}

func hasInvalidLetters(pw []rune) bool {
	for _, char := range pw {
		if char == 'i' || char == 'o' || char == 'l' {
			return true
		}
	}
	return false
}

func hasTwoPairs(pw []rune) bool {
	pairs := make(map[rune]bool)

	for i := 0; i < len(pw)-1; i++ {
		if pw[i] == pw[i+1] {
			pairs[pw[i]] = true
			i++
		}
	}

	return len(pairs) >= 2
}

func isValid(pw []rune) bool {
	return hasIncreasingStraight(pw) &&
		!hasInvalidLetters(pw) &&
		hasTwoPairs(pw)
}

func nextValid(pw string) string {
	password := []rune(pw)

	for {
		incrementPassword(password)

		if isValid(password) {
			return string(password)
		}
	}
}

func solve(input string) (string, string) {
	part1 := nextValid(input)
	part2 := nextValid(part1)
	return part1, part2
}

func main() {
	start := time.Now()
	const iters = 100
	input := "cqjxjnds"
	var part1, part2 string
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
