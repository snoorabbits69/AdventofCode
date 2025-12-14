package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int64
}

type Edge struct {
	i, j int
	d2   int64
}

type EdgeMaxHeap []Edge

func (h EdgeMaxHeap) Len() int           { return len(h) }
func (h EdgeMaxHeap) Less(a, b int) bool { return h[a].d2 > h[b].d2 }
func (h EdgeMaxHeap) Swap(a, b int)      { h[a], h[b] = h[b], h[a] }

func (h *EdgeMaxHeap) Push(x interface{}) {
	*h = append(*h, x.(Edge))
}

func (h *EdgeMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	sz := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{parent: p, size: sz}
}

func (d *DSU) Find(a int) int {
	for a != d.parent[a] {
		d.parent[a] = d.parent[d.parent[a]]
		a = d.parent[a]
	}
	return a
}

func (d *DSU) Union(a, b int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func dist2(a, b Point) int64 {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return dx*dx + dy*dy + dz*dz
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	var pts []Point
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			panic("bad line (expected X,Y,Z): " + line)
		}
		x, e1 := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
		y, e2 := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		z, e3 := strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 64)
		if e1 != nil || e2 != nil || e3 != nil {
			panic("bad numbers: " + line)
		}
		pts = append(pts, Point{x, y, z})
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	n := len(pts)
	if n == 0 {
		fmt.Println(0)
		return
	}

	const K = 1000

	h := &EdgeMaxHeap{}
	heap.Init(h)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d2 := dist2(pts[i], pts[j])
			e := Edge{i: i, j: j, d2: d2}

			if h.Len() < K {
				heap.Push(h, e)
			} else if (*h)[0].d2 > d2 {
				heap.Pop(h)
				heap.Push(h, e)
			}
		}
	}

	edges := make([]Edge, h.Len())
	for i := len(edges) - 1; i >= 0; i-- {
		edges[i] = heap.Pop(h).(Edge)
	}
	sort.Slice(edges, func(a, b int) bool { return edges[a].d2 < edges[b].d2 })

	dsu := NewDSU(n)
	for _, e := range edges {
		dsu.Union(e.i, e.j)
	}

	seen := make(map[int]bool)
	var sizes []int
	for i := 0; i < n; i++ {
		r := dsu.Find(i)
		if !seen[r] {
			seen[r] = true
			sizes = append(sizes, dsu.size[r])
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	if len(sizes) == 1 {
		fmt.Println(sizes[0])
		return
	}
	if len(sizes) == 2 {
		fmt.Println(sizes[0] * sizes[1])
		return
	}
	fmt.Println(sizes[0] * sizes[1] * sizes[2])
}
