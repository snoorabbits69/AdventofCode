package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

func zeroHitsLeft(pos, val int16) int16 {
	if val <= 0 {
		return 0
	}
	t0 := pos
	if t0 == 0 {
		t0 = 100
	}
	if val < t0 {
		return 0
	}
	return 1 + (val-t0)/100

}
func zeroHitsRight(pos, val int16) int16 {
	if val <= 0 {
		return 0
	}
	t0 := (100 - pos) % 100
	if t0 == 0 {
		t0 = 100
	}

	if val < t0 {
		return 0
	}
	return 1 + (val-t0)/100

}

func wrap100(x int16) int16 {
	x %= 100
	if x < 0 {
		x += 100
	}
	return x
}

func solve() (int16, int16) {
	n := len(data)
	i := 0
	var pos int16 = 50
	var occurpos int16 = 50
	var zeroCount int16
	var zeroOccurs int16
	var val int16
	for i < n {
		if data[i] == 'L' {
			i++

			val = parseNum(&i)
			pos -= val

			zeroOccurs += zeroHitsLeft(occurpos, val)
			occurpos = wrap100(occurpos - val)

			pos %= 100
			if pos < 0 {
				pos += 100
			}
			if pos == 0 {
				zeroCount++

			}

			continue
		}

		if data[i] == 'R' {
			i++

			val = parseNum(&i)
			zeroOccurs += zeroHitsRight(occurpos, val)
			occurpos = wrap100(occurpos + val)
			pos += val

			pos %= 100
			if pos < 0 {
				pos += 100
			}
			if pos == 0 {
				zeroCount++

			}

			continue
		}
		i++
	}
	return zeroCount, zeroOccurs
}

func parseNum(i *int) int16 {
	var v int16

	for *i < len(data) && data[*i] != '\n' {
		v = v*10 + int16(data[*i]-'0')
		*i++
	}

	return v
}

func main() {
	start := time.Now()
	const iters = 10000
	var part1, part2 int16
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve()
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
