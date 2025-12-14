package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type memoVal struct {
	count  int
	cached bool
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	graph := make(map[string][]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		graph[key] = strings.Fields(parts[1])
	}

	start := "svr"
	end := "out"

	memo := make(map[string]memoVal)

	dacSeen := false
	fftSeen := false
	if start == "dac" {
		dacSeen = true
	}
	if start == "fft" {
		fftSeen = true
	}

	fmt.Println(dfsWithParams(graph, start, end, dacSeen, fftSeen, memo))
}

func dfsWithParams(
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
		total += dfsWithParams(graph, next, end, nd, nf, memo)
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
