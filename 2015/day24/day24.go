package main

import (
	_ "embed"
	"fmt"
	"sort"
	"time"
)

//go:embed input.txt
var data []byte

func parse() []int {
	var nums []int
	v, in := 0, false
	for _, c := range data {
		if c >= '0' && c <= '9' {
			v = v*10 + int(c-'0')
			in = true
		} else if in {
			nums = append(nums, v)
			v, in = 0, false
		}
	}
	if in {
		nums = append(nums, v)
	}
	return nums
}

func solution(nums []int, groups int) uint64 {
	total := 0
	for _, v := range nums {
		total += v
	}
	target := total / groups
	n := len(nums)

	var canSplit func(used uint64, g int) bool
	var pick func(used, cur uint64, idx, sum, g int) bool

	pick = func(used, cur uint64, idx, sum, g int) bool {
		if sum == target {
			return canSplit(used|cur, g-1)
		}
		for i := idx; i < n; i++ {
			if used&(1<<i) != 0 || sum+nums[i] > target {
				continue
			}
			if pick(used, cur|1<<i, i+1, sum+nums[i], g) {
				return true
			}
		}
		return false
	}

	canSplit = func(used uint64, g int) bool {
		if g == 1 {
			return true
		}
		return pick(used, 0, 0, 0, g)
	}

	var best uint64
	found := false

	var choose func(start, left, sum int, prod, mask uint64)
	choose = func(start, left, sum int, prod, mask uint64) {
		if left == 0 {
			if sum != target || (found && prod >= best) {
				return
			}
			if canSplit(mask, groups-1) {
				best, found = prod, true
			}
			return
		}
		for i := start; i < n; i++ {
			if sum+nums[i] > target {
				continue
			}
			choose(i+1, left-1, sum+nums[i], prod*uint64(nums[i]), mask|1<<i)
		}
	}

	for k := 1; k <= n; k++ {
		choose(0, k, 0, 1, 0)
		if found {
			return best
		}
	}
	return 0
}

func solve(nums []int) (uint64, uint64) {
	return solution(nums, 3), solution(nums, 4)
}

func main() {
	start := time.Now()
	const iters = 10
	nums := parse()
	sort.Sort(sort.Reverse(sort.IntSlice(nums)))
	var part1, part2 uint64
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(nums)
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
