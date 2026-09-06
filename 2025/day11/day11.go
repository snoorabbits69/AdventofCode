package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"
)

//go:embed input.txt
var data []byte

type memoVal struct {
	count  int
	cached bool
}

func parseGraph(data []byte) map[string][]string {
	graph := make(map[string][]string)

	i := 0
	n := len(data)

	for i < n {
		var key strings.Builder

		for i < n && data[i] != ':' {
			if data[i] != '\n' && data[i] != '\r' {
				key.WriteByte(data[i])
			}
			i++
		}

		if i >= n {
			break
		}
		i++

		for i < n && data[i] == ' ' {
			i++
		}

		var valBuf strings.Builder
		var values []string

		for i < n && data[i] != '\n' {
			if data[i] == ' ' {
				if valBuf.Len() > 0 {
					values = append(values, valBuf.String())
					valBuf.Reset()
				}
			} else if data[i] != '\r' {
				valBuf.WriteByte(data[i])
			}
			i++
		}
		if valBuf.Len() > 0 {
			values = append(values, valBuf.String())
		}

		graph[key.String()] = values

		i++
	}

	return graph
}

func part1(
	graph map[string][]string,
	node, end string,
	memo map[string]int,
) int {
	if node == end {
		return 1
	}

	if v, ok := memo[node]; ok {
		return v
	}

	total := 0
	for _, next := range graph[node] {
		total += part1(graph, next, end, memo)
	}

	memo[node] = total
	return total
}

func part2(
	graph map[string][]string,
	node, end string,
	dacSeen, fftSeen bool,
	memo map[string]memoVal,
) int {

	if node == end {
		if dacSeen && fftSeen {
			return 1
		}
		return 0
	}

	key := node + "|" + b2s(dacSeen) + "|" + b2s(fftSeen)
	if v, ok := memo[key]; ok && v.cached {
		return v.count
	}

	total := 0
	for _, next := range graph[node] {
		nd := dacSeen || next == "dac"
		nf := fftSeen || next == "fft"
		total += part2(graph, next, end, nd, nf, memo)
	}

	memo[key] = memoVal{count: total, cached: true}
	return total
}

func b2s(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func main() {
	start := time.Now()
	const iters = 10000
	graph := parseGraph(data)
	s := "svr"
	e := "out"
	var part1Ans, part2Ans int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		m1 := make(map[string]int)
		m2 := make(map[string]memoVal)
		part1Ans = part1(graph, s, e, m1)
		part2Ans = part2(graph, s, e, false, false, m2)
	}
	processTime := time.Since(processStart)
	totalTime := time.Since(start)
	elapsedUs := float64(processTime.Nanoseconds()) / 1000.0
	fmt.Printf("Total: %.2f microseconds\n", elapsedUs)
	fmt.Printf("Average: %.4f microseconds\n", elapsedUs/float64(iters))
	fmt.Printf("Total time:        %dns\n", totalTime.Nanoseconds())
	fmt.Println("part1", part1Ans)
	fmt.Println("part2", part2Ans)
}
