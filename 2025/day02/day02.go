package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

type Pair [2]int64
type Range [2]int64

var FIRST = []Range{{2, 1}, {4, 2}, {6, 3}, {8, 4}, {10, 5}}

var SECOND = []Range{{3, 1}, {5, 1}, {6, 2}, {7, 1}, {9, 3}, {10, 2}}

var THIRD = []Range{{6, 1}, {10, 1}}

var powerOf10 = []int64{
	1, 10, 100, 1000, 10000, 100000, 1000000,
	10000000, 100000000, 1000000000, 10000000000,
	100000000000, 1000000000000, 10000000000000,
	100000000000000, 1000000000000000,
	10000000000000000, 100000000000000000,
	1000000000000000000,
}

func parse() []Pair {
	var pairs []Pair
	var nums []int64
	i, n := 0, len(data)
	for i < n {
		for i < n && (data[i] < '0' || data[i] > '9') {
			i++
		}
		if i >= n {
			break
		}
		var v int64
		for i < n && data[i] >= '0' && data[i] <= '9' {
			v = v*10 + int64(data[i]-'0')
			i++
		}
		nums = append(nums, v)
	}
	for j := 0; j+1 < len(nums); j += 2 {
		pairs = append(pairs, Pair{nums[j], nums[j+1]})
	}
	return pairs
}

func nextMultipleOf(value, step int64) int64 {
	rem := value % step
	if rem == 0 {
		return value
	}
	return value + (step - rem)
}

func sumRanges(ranges []Range, input []Pair) int64 {
	var result int64
	for _, r := range ranges {
		digits, size := r[0], r[1]

		digitsPower := powerOf10[digits]
		sizePower := powerOf10[size]

		step := (digitsPower - 1) / (sizePower - 1)
		start := step * (sizePower / 10)
		end := step * (sizePower - 1)

		for _, pair := range input {
			from, to := pair[0], pair[1]

			lower := nextMultipleOf(from, step)
			if lower < start {
				lower = start
			}
			upper := to
			if upper > end {
				upper = end
			}

			if lower <= upper {
				n := (upper - lower) / step
				triangular := n * (n + 1) / 2
				result += lower*(n+1) + step*triangular
			}
		}
	}
	return result
}

func solve(input []Pair) (int64, int64) {
	part1 := sumRanges(FIRST, input)
	part2 := sumRanges(FIRST, input) + sumRanges(SECOND, input) - sumRanges(THIRD, input)
	return part1, part2
}

func main() {
	start := time.Now()
	const iters = 10000
	input := parse()
	var part1, part2 int64
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
