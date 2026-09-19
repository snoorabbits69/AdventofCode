package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

func parseNumber(i *int) int64 {
	var v int64
	for *i < len(data) {
		c := data[*i]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int64(c-'0')
		*i++
	}
	return v
}

type frame struct {
	sum      int64
	hasRed   bool
	isObject bool
}

func solve() (int64, int64) {
	n := len(data)

	part1 := int64(0)
	excluded := int64(0)

	inStr := false
	afterColon := false
	var stack []frame

	for i := 0; i < n; i++ {
		c := data[i]

		if inStr {
			if c == '\\' {
				i++
				continue
			}
			if c == '"' {
				inStr = false
			}
			continue
		}

		switch {
		case c == '"':
			if afterColon && len(stack) > 0 && stack[len(stack)-1].isObject &&
				i+4 < n && data[i+1] == 'r' && data[i+2] == 'e' && data[i+3] == 'd' {
				stack[len(stack)-1].hasRed = true
			}
			inStr = true
			afterColon = false

		case c == ':':
			afterColon = true

		case c == ',':
			afterColon = false

		case c == '{':
			stack = append(stack, frame{isObject: true})

		case c == '[':
			stack = append(stack, frame{isObject: false})

		case c == '}', c == ']':
			f := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if f.hasRed {
				excluded += f.sum
			} else if len(stack) > 0 {
				stack[len(stack)-1].sum += f.sum
			}

		case c == '-' || (c >= '0' && c <= '9'):
			neg := c == '-'
			if neg {
				i++
			}
			num := parseNumber(&i)
			i--
			if neg {
				num = -num
			}
			part1 += num
			if len(stack) > 0 {
				stack[len(stack)-1].sum += num
			}
			afterColon = false
		}
	}

	return part1, part1 - excluded
}

func main() {
	start := time.Now()
	const iters = 10000
	var part1, part2 int64
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
