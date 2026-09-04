package main

import (
	_ "embed"
	"fmt"
	"time"
)

type Interval struct {
	start int64
	end   int64
}

var (
	//go:embed input.txt
	data []byte
)
var n = len(data)

func parseNumber(i *int) int64 {
	var v int64
	for *i < n {
		c := data[*i]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int64(c-'0')
		*i++
	}
	return v
}

func lowerBound(a []Interval, x int64) int {
	l, r := 0, len(a)
	for l < r {
		m := (l + r) / 2
		if a[m].start <= x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func insertTree(tree *[]Interval, x Interval) {
	t := *tree

	pos := lowerBound(t, x.start)

	t = append(t, Interval{})
	copy(t[pos+1:], t[pos:])
	t[pos] = x

	res := make([]Interval, 0, len(t))

	for _, cur := range t {
		if len(res) == 0 {
			res = append(res, cur)
			continue
		}

		last := &res[len(res)-1]

		if last.end < cur.start {
			res = append(res, cur)
		} else {
			if cur.end > last.end {
				last.end = cur.end
			}
		}
	}

	*tree = res
}

func contains(tree []Interval, v int64) bool {
	l, r := 0, len(tree)

	for l < r {
		m := (l + r) / 2
		if tree[m].start <= v {
			l = m + 1
		} else {
			r = m
		}
	}

	i := l - 1
	if i < 0 {
		return false
	}

	return tree[i].end >= v
}

func parseData() ([]Interval, []int64) {
	i := 0
	n := n

	tree := make([]Interval, 0)
	values := make([]int64, 0)

	for i < n {
		for i < n {
			c := data[i]
			if c >= '0' && c <= '9' {
				break
			}
			i++
		}
		if i >= n {
			break
		}

		start := parseNumber(&i)

		if i < n && data[i] == '-' {
			i++
			end := parseNumber(&i)
			insertTree(&tree, Interval{start, end})
		} else {
			values = append(values, start)
		}

		for i < n && data[i] != '\n' {
			i++
		}
		i++
	}

	return tree, values

}

func solvePart1(tree []Interval, values []int64) int64 {
	var count int64
	for _, v := range values {
		if contains(tree, v) {
			count++
		}
	}
	return count

}

func solvePart2(tree []Interval) int64 {
	var total int64
	for _, t := range tree {
		total += (t.end - t.start + 1)
	}
	return total
}

func solve(tree []Interval, values []int64) (int64, int64) {
	part1 := solvePart1(tree, values)
	part2 := solvePart2(tree)
	return part1, part2
}

func main() {
	start := time.Now()
	const iters = 10000
	tree, values := parseData()
	var part1, part2 int64
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(tree, values)
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
