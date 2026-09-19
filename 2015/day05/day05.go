package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

var badNext = [256]byte{'a': 'b', 'c': 'd', 'p': 'q', 'x': 'y'}

func solve(data []byte) (int, int) {
	return part1(data), part2(data)
}

func part1(data []byte) int {
	vowels := 0
	hasDup := false
	nice := 0
	n := len(data)

	for i := 0; i < n; i++ {
		d := data[i]

		if d == '\n' || d == '\r' {
			if vowels >= 3 && hasDup {
				nice++
			}
			vowels = 0
			hasDup = false

			for i+1 < n && (data[i+1] == '\n' || data[i+1] == '\r') {
				i++
			}
			continue
		}

		if isVowel(d) {
			vowels++
		}

		if i+1 < n && data[i+1] != '\n' && data[i+1] != '\r' {
			next := data[i+1]

			if bn := badNext[d]; bn != 0 && next == bn {
				for i < n && data[i] != '\n' {
					i++
				}
				vowels = 0
				hasDup = false

				for i+1 < n && data[i+1] == '\n' {
					i++
				}
				continue
			}

			if d == next {
				hasDup = true
			}
		}
	}

	if n > 0 && data[n-1] != '\n' && data[n-1] != '\r' {
		if vowels >= 3 && hasDup {
			nice++
		}
	}

	return nice
}

func isVowel(b byte) bool {
	return b == 'a' || b == 'e' || b == 'i' || b == 'o' || b == 'u' ||
		b == 'A' || b == 'E' || b == 'I' || b == 'O' || b == 'U'
}

func part2(data []byte) int {
	n := len(data)

	hasSandwich := false
	hasPair := false
	nice := 0

	firstPos := [256 * 256]int{}
	seenLine := [256 * 256]uint32{}
	var lineID uint32 = 1

	lineStart := 0
	linePos := 0

	for i := 0; i < n; i++ {
		d := data[i]

		if d == '\n' || d == '\r' {
			if i > lineStart {
				if hasSandwich && hasPair {
					nice++
				}
			}

			hasSandwich = false
			hasPair = false

			for i+1 < n && (data[i+1] == '\n' || data[i+1] == '\r') {
				i++
			}
			lineStart = i + 1
			linePos = 0
			lineID++
			continue
		}

		if !hasSandwich &&
			i+2 < n &&
			data[i+1] != '\n' &&
			data[i+2] != '\n' &&
			d == data[i+2] {
			hasSandwich = true
		}

		if !hasPair &&
			i+1 < n &&
			data[i+1] != '\n' {
			key := int(d)<<8 | int(data[i+1])
			if seenLine[key] != lineID {
				seenLine[key] = lineID
				firstPos[key] = linePos
			} else if linePos-firstPos[key] >= 2 {
				hasPair = true
			}
		}

		if hasSandwich && hasPair {
			for i+1 < n && data[i+1] != '\n' && data[i+1] != '\r' {
				i++
			}
		}

		linePos++
	}

	if lineStart < n {
		if hasSandwich && hasPair {
			nice++
		}
	}

	return nice
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
