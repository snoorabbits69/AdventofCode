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
	nums := make([]int, 0, 32)
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

type solver struct {
	nums   []int
	n      int
	target int
	groups int
	best   uint64
	found  bool
}

func (s *solver) canSplit(used uint64, g int) bool {
	if g == 1 {
		return true
	}
	return s.pick(used, 0, 0, 0, g)
}

func (s *solver) pick(used, cur uint64, idx, sum, g int) bool {
	if sum == s.target {
		return s.canSplit(used|cur, g-1)
	}
	for i := idx; i < s.n; i++ {
		v := s.nums[i]
		if used&(1<<uint(i)) != 0 || sum+v > s.target {
			continue
		}
		if s.pick(used, cur|1<<uint(i), i+1, sum+v, g) {
			return true
		}
	}
	return false
}

func (s *solver) choose(start, left, sum int, prod, mask uint64) {
	if s.found && prod >= s.best {
		return
	}
	if left == 1 {
		need := s.target - sum
		lo, hi := start, s.n
		for lo < hi {
			m := (lo + hi) >> 1
			if s.nums[m] > need {
				lo = m + 1
			} else {
				hi = m
			}
		}
		if lo < s.n && s.nums[lo] == need {
			p := prod * uint64(need)
			if (!s.found || p < s.best) && s.canSplit(mask|1<<uint(lo), s.groups-1) {
				s.best, s.found = p, true
			}
		}
		return
	}
	for i := start; i <= s.n-left; i++ {
		v := s.nums[i]
		if sum+v*left < s.target {
			break
		}
		if sum+v > s.target {
			continue
		}
		s.choose(i+1, left-1, sum+v, prod*uint64(v), mask|1<<uint(i))
	}
}

func solution(nums []int, groups int) uint64 {
	total := 0
	for _, v := range nums {
		total += v
	}
	s := solver{nums: nums, n: len(nums), target: total / groups, groups: groups}

	k, acc := 0, 0
	for acc < s.target {
		acc += nums[k]
		k++
	}
	for ; k <= s.n; k++ {
		s.choose(0, k, 0, 1, 0)
		if s.found {
			return s.best
		}
	}
	return 0
}

func solve(nums []int) (uint64, uint64) {
	return solution(nums, 3), solution(nums, 4)
}

func main() {
	start := time.Now()
	const iters = 10000
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
