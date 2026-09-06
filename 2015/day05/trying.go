// package main

// import (
// 	_ "embed"
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"
// )

// //go:embed input.txt
// var data []byte

// var badNext = [256]byte{'a': 'b', 'c': 'd', 'p': 'q', 'x': 'y'}

// func main() {
// 	start := time.Now()
// 	data, err := os.ReadFile("input.txt")
// 	if err != nil {
// 		panic(err)
// 	}
// 	readTime := time.Since(start)

// 	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

// 	var part1, part2 int
// 	for i := 0; i < 100; i++ {
// 		part1, part2 = solve(lines)
// 	}

// 	const iters = 10000
// 	benchStart := time.Now()
// 	for i := 0; i < iters; i++ {
// 		part1, part2 = solve(lines)
// 	}

// 	benchTime := time.Since(benchStart)
// 	totalTime := time.Since(start)

// 	fmt.Println("solution of part1:", part1)
// 	fmt.Println("solution of part2:", part2)
// 	fmt.Println()

// 	fmt.Printf("File read time:    %d microseconds\n", readTime.Microseconds())
// 	fmt.Printf("Total bench time:  %d microseconds\n", benchTime.Microseconds())

// 	avgUs := float64(benchTime.Nanoseconds()) / 1_000.0 / float64(iters)
// 	fmt.Printf("Average per iter:  %.4f microseconds\n", avgUs)

// 	fmt.Printf("Total time:        %d microseconds\n", totalTime.Microseconds())
// }

// func solve(lines []string) (int, int) {
// 	nice1 := 0
// 	nice2 := 0

// 	for _, line := range lines {
// 		if isNicePart1(line) {
// 			nice1++
// 		}
// 		if isNicePart2(line) {
// 			nice2++
// 		}
// 	}

// 	return nice1, nice2
// }

// func isNicePart1(line string) bool {
// 	vowelCount := 0
// 	hasDuplicate := false

// 	for i := 0; i < len(line); i++ {
// 		ch := line[i]

// 		if ch == 'a' || ch == 'e' || ch == 'i' || ch == 'o' || ch == 'u' {
// 			vowelCount++
// 		}

// 		if i > 0 {
// 			prevCh := line[i-1]

// 			if ch == prevCh {
// 				hasDuplicate = true
// 			}

// 			if bn := badNext[prevCh]; bn != 0 && ch == bn {
// 				return false
// 			}
// 		}
// 	}

// 	return vowelCount >= 3 && hasDuplicate
// }

// func isNicePart2(line string) bool {
// 	hasNonOverlappingPair := false
// 	hasSandwich := false
// 	n := len(line)

// 	for i := 0; i < n; i++ {
// 		if !hasSandwich && i < n-2 {
// 			if line[i] == line[i+2] {
// 				hasSandwich = true
// 			}
// 		}

// 		if !hasNonOverlappingPair && i < n-1 {
// 			pair := line[i : i+2]
// 			if strings.Contains(line[i+2:], pair) {
// 				hasNonOverlappingPair = true
// 			}
// 		}

// 		if hasSandwich && hasNonOverlappingPair {
// 			return true
// 		}
// 	}

// 	return hasSandwich && hasNonOverlappingPair
// }
