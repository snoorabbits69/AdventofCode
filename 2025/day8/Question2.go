package main

import (
	"bufio"
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

type DSU struct {
	parent []int
	size   []int
	sets   int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	sz := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{parent: p, size: sz, sets: n}
}

func (d *DSU) Find(a int) int {
	for a != d.parent[a] {
		d.parent[a] = d.parent[d.parent[a]]
		a = d.parent[a]
	}
	return a
}

func (d *DSU) Union(a, b int) bool {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return false
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	d.sets--
	return true
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
	if n < 2 {
		fmt.Println(0)
		return
	}

	edges := make([]Edge, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			edges = append(edges, Edge{i: i, j: j, d2: dist2(pts[i], pts[j])})
		}
	}

	sort.Slice(edges, func(a, b int) bool {
		if edges[a].d2 != edges[b].d2 {
			return edges[a].d2 < edges[b].d2
		}
		if edges[a].i != edges[b].i {
			return edges[a].i < edges[b].i
		}
		return edges[a].j < edges[b].j
	})

	dsu := NewDSU(n)

	lastMerge := Edge{-1, -1, 0}

	for _, e := range edges {
		if dsu.Union(e.i, e.j) {
			lastMerge = e
			if dsu.sets == 1 {
				break
			}
		}
	}

	if lastMerge.i == -1 {
		fmt.Println(0)
		return
	}

	ans := pts[lastMerge.i].x * pts[lastMerge.j].x
	fmt.Println(ans)
}
