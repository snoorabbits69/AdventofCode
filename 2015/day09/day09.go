package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var data []byte

type Edge struct {
	To     int
	Weight int64
}

type Graph struct {
	IDs   map[string]int
	Names []string
	Adj   [][]Edge
}

func parseNumber(i *int) int64 {
	var v int64

	for *i < len(data) {
		c := data[*i]

		if c < '0' || c > '9' {
			break
		}

		v = v*10 + int64(c-'0')
		*i = *i + 1
	}

	return v
}

func (g *Graph) getID(name string) int {
	if id, ok := g.IDs[name]; ok {
		return id
	}

	id := len(g.Adj)

	g.IDs[name] = id
	g.Names = append(g.Names, name)
	g.Adj = append(g.Adj, nil)

	return id
}

func parse() Graph {
	graph := Graph{
		IDs: make(map[string]int),
	}

	i := 0
	n := len(data)

	for i < n {
		for i < n && (data[i] == ' ' || data[i] == '\n' || data[i] == '\r' || data[i] == '\t') {
			i++
		}

		if i >= n {
			break
		}

		start := i
		for i < n && data[i] != ' ' {
			i++
		}
		if i >= n {
			break
		}
		from := string(data[start:i])

		for i < n && data[i] == ' ' {
			i++
		}

		if i+1 >= n || data[i] != 't' || data[i+1] != 'o' {
			break
		}
		i += 2

		for i < n && data[i] == ' ' {
			i++
		}

		start = i
		for i < n && data[i] != ' ' && data[i] != '=' {
			i++
		}
		if i >= n {
			break
		}
		to := string(data[start:i])

		for i < n && data[i] == ' ' {
			i++
		}

		if i >= n || data[i] != '=' {
			break
		}
		i++

		for i < n && data[i] == ' ' {
			i++
		}

		weight := parseNumber(&i)

		fromID := graph.getID(from)
		toID := graph.getID(to)

		graph.Adj[fromID] = append(graph.Adj[fromID], Edge{
			To:     toID,
			Weight: weight,
		})

		graph.Adj[toID] = append(graph.Adj[toID], Edge{
			To:     fromID,
			Weight: weight,
		})

		for i < n && data[i] != '\n' {
			i++
		}
		if i < n {
			i++
		}
	}

	return graph
}

func solve(graph Graph) (minDist int64, maxDist int64) {
	numNodes := len(graph.Adj)
	if numNodes == 0 {
		return 0, 0
	}

	const noEdge = int64(-1)
	dist := make([][]int64, numNodes)
	for i := range dist {
		dist[i] = make([]int64, numNodes)
		for j := range dist[i] {
			dist[i][j] = noEdge
		}
	}
	for u, edges := range graph.Adj {
		for _, e := range edges {
			dist[u][e.To] = e.Weight
		}
	}

	fullMask := 1 << numNodes
	const inf = int64(1) << 62
	const unset = int64(-1)

	dpMin := make([][]int64, fullMask)
	dpMax := make([][]int64, fullMask)
	for m := range dpMin {
		dpMin[m] = make([]int64, numNodes)
		dpMax[m] = make([]int64, numNodes)
		for j := range dpMin[m] {
			dpMin[m][j] = inf
			dpMax[m][j] = unset
		}
	}

	for i := 0; i < numNodes; i++ {
		dpMin[1<<i][i] = 0
		dpMax[1<<i][i] = 0
	}

	for mask := 1; mask < fullMask; mask++ {
		for j := 0; j < numNodes; j++ {
			if mask&(1<<j) == 0 {
				continue
			}

			curMin := dpMin[mask][j]
			curMax := dpMax[mask][j]
			minValid := curMin != inf
			maxValid := curMax != unset

			if !minValid && !maxValid {
				continue
			}

			for k := 0; k < numNodes; k++ {
				if mask&(1<<k) != 0 {
					continue
				}
				w := dist[j][k]
				if w == noEdge {
					continue
				}

				next := mask | (1 << k)

				if minValid {
					if nd := curMin + w; nd < dpMin[next][k] {
						dpMin[next][k] = nd
					}
				}
				if maxValid {
					if nd := curMax + w; nd > dpMax[next][k] {
						dpMax[next][k] = nd
					}
				}
			}
		}
	}

	last := fullMask - 1
	bestMin := inf
	bestMax := unset
	for j := 0; j < numNodes; j++ {
		if dpMin[last][j] < bestMin {
			bestMin = dpMin[last][j]
		}
		if dpMax[last][j] > bestMax {
			bestMax = dpMax[last][j]
		}
	}

	return bestMin, bestMax
}

func main() {
	graph := parse()
	minResult, maxResult := solve(graph)
	fmt.Printf("Shortest distance to visit all nodes: %d\n", minResult)
	fmt.Printf("Longest distance to visit all nodes: %d\n", maxResult)
}
