package main

import (
	_ "embed"
	"fmt"
	"sort"
	"time"
)

//go:embed input.txt
var data []byte

const (
	buckets = 4
	size    = 10_000 * 10_000
)

type pair struct {
	i, j     uint16
	distance int
}

type node struct {
	parent int
	size   int
}

func distance(a, b []int) int {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	dz := a[2] - b[2]
	return dx*dx + dy*dy + dz*dz
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	start := time.Now()
	const iters = 1000
	i := 0
	var grids [][]int
	var grid []int

	for i < len(data) {
		n := parseNum(&i)
		grid = append(grid, n)
		if i < len(data) && data[i] == '\n' {
			grids = append(grids, grid)
			grid = []int{}

		}
		if i < len(data) {
			i++
		}
	}

	buckets := buildBuckets(grids)

	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(grids, buckets)
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

func solve(grids [][]int, buckets [][]pair) (int, int) {
	return part1(grids, buckets), part2(grids, buckets)
}

func parseNum(i *int) int {
	var v int
	for *i < len(data) && data[*i] != '\n' && data[*i] != ',' {
		v = v*10 + int(data[*i]-'0')
		*i++
	}
	return v
}

func buildBuckets(boxes [][]int) [][]pair {
	result := make([][]pair, buckets)

	for i := range boxes {
		v1 := boxes[i]
		for j := i + 1; j < len(boxes); j++ {
			v2 := boxes[j]

			dx := abs(v1[0] - v2[0])
			dy := abs(v1[1] - v2[1])
			dz := abs(v1[2] - v2[2])
			d := dx*dx + dy*dy + dz*dz

			index := d / size
			if index < buckets {
				result[index] = append(result[index], pair{uint16(i), uint16(j), d})
			}
		}
	}

	for b := range result {
		sort.Slice(result[b], func(x, y int) bool {
			return result[b][x].distance < result[b][y].distance
		})
	}

	return result
}

func flatten(buckets [][]pair) []pair {
	var out []pair
	for _, b := range buckets {
		out = append(out, b...)
	}
	return out
}

func find(nodes []node, x int) int {
	for nodes[x].parent != x {
		parent := nodes[x].parent
		grandparent := nodes[parent].parent
		nodes[x].parent = grandparent
		x = parent
	}
	return x
}

func union(nodes []node, x, y int) int {
	x = find(nodes, x)
	y = find(nodes, y)
	if x != y {
		if nodes[x].size < nodes[y].size {
			x, y = y, x
		}
		nodes[y].parent = x
		nodes[x].size += nodes[y].size
	}
	return nodes[x].size
}

func newNodes(n int) []node {
	nodes := make([]node, n)
	for i := range nodes {
		nodes[i] = node{parent: i, size: 1}
	}
	return nodes
}

func part1(boxes [][]int, buckets [][]pair) int {
	return part1Testable(boxes, buckets, 1000)
}

func part1Testable(boxes [][]int, buckets [][]pair, limit int) int {
	nodes := newNodes(len(boxes))

	pairs := flatten(buckets)
	if limit > len(pairs) {
		limit = len(pairs)
	}
	for _, p := range pairs[:limit] {
		union(nodes, int(p.i), int(p.j))
	}

	sizes := make([]int, len(nodes))
	for i, n := range nodes {
		sizes[i] = n.size
	}
	sort.Ints(sizes)

	result := 1
	for i := len(sizes) - 1; i >= 0 && i >= len(sizes)-3; i-- {
		result *= sizes[i]
	}
	return result
}

func part2(boxes [][]int, buckets [][]pair) int {
	nodes := newNodes(len(boxes))

	for _, p := range flatten(buckets) {
		i, j := int(p.i), int(p.j)
		if union(nodes, i, j) == len(boxes) {
			return boxes[i][0] * boxes[j][0]
		}
	}

	panic("unreachable")
}
