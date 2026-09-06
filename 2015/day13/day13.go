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
	n := len(data)
	graph := Graph{
		IDs: make(map[string]int),
	}
	var start int
	var end int
	var num int64

	for i := 0; i < n; i++ {
		start = i

		for i < n && data[i] != ' ' {
			i++
		}
		from := string(data[start:i])
		i = i + 7

		if i < n && data[i] == 'g' {

			i += 5
			num = parseNumber(&i)
		}
		if i < n && data[i] == 'l' {

			i += 5
			num = -parseNumber(&i)
		}

		for i < n && data[i] != 'o' {
			i++
		}
		i++

		for i < n && data[i] == ' ' {
			i++
		}

		end = i
		for i < n && data[i] != '.' {
			i++
		}
		to := string(data[end:i])

		fromID := graph.getID(from)
		toID := graph.getID(to)
		graph.Adj[fromID] = append(graph.Adj[fromID], Edge{
			To:     toID,
			Weight: num,
		})
		i++

		for i < n && (data[i] == '\n' || data[i] == ' ') {
			i++
		}
		i--
	}
	return graph
}

func buildCostMatrix(graph Graph) [][]int64 {
	numPeople := len(graph.Adj)

	cost := make([][]int64, numPeople)
	for i := range cost {
		cost[i] = make([]int64, numPeople)
	}
	for u, edges := range graph.Adj {
		for _, e := range edges {
			cost[u][e.To] += e.Weight
		}
	}
	for u := 0; u < numPeople; u++ {
		for v := u + 1; v < numPeople; v++ {
			combined := cost[u][v] + cost[v][u]
			cost[u][v] = combined
			cost[v][u] = combined
		}
	}
	return cost
}

func extendCostMatrix(cost [][]int64) [][]int64 {
	n := len(cost)
	newCost := make([][]int64, n+1)
	for i := 0; i < n; i++ {
		newCost[i] = append(append([]int64{}, cost[i]...), 0)
	}
	newCost[n] = make([]int64, n+1)
	return newCost
}

func solve(graph Graph) (int64, int64) {
	cost := buildCostMatrix(graph)
	part1 := heldKarp(cost, len(graph.Adj))

	fromID := graph.getID("me")

	for _, value := range graph.IDs {
		graph.Adj[fromID] = append(graph.Adj[fromID], Edge{
			To:     value,
			Weight: 0,
		})
	}

	cost = extendCostMatrix(cost)
	part2 := heldKarp(cost, len(graph.Adj))

	return part1, part2
}

func heldKarp(cost [][]int64, n int) int64 {
	const negInf = int64(-1) << 62
	full := 1 << n

	dp := make([][]int64, full)
	for m := range dp {
		dp[m] = make([]int64, n)
		for j := range dp[m] {
			dp[m][j] = negInf
		}
	}
	dp[1][0] = 0

	for mask := 1; mask < full; mask++ {
		if mask&1 == 0 {
			continue
		}
		for last := 0; last < n; last++ {
			if mask&(1<<last) == 0 {
				continue
			}
			cur := dp[mask][last]
			if cur == negInf {
				continue
			}
			for v := 0; v < n; v++ {
				if mask&(1<<v) != 0 {
					continue
				}
				nmask := mask | (1 << v)
				val := cur + cost[last][v]
				if val > dp[nmask][v] {
					dp[nmask][v] = val
				}
			}
		}
	}

	if n == 1 {
		return 0
	}

	fullMask := full - 1
	best := negInf
	for last := 1; last < n; last++ {
		if dp[fullMask][last] == negInf {
			continue
		}
		total := dp[fullMask][last] + cost[last][0]
		if total > best {
			best = total
		}
	}
	return best
}

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

func main() {
	graph := parse()
	part1, part2 := solve(graph)
	fmt.Println("Part 1 ", part1)
	fmt.Println("Part 2 ", part2)

}
