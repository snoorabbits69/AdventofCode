package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	f, err := os.Open("input.txt")

	if err != nil {
		panic(err)
	}

	graph := make(map[string][]string)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		key := strings.TrimSpace(parts[0])
		var values []string = strings.Fields(parts[1])
		graph[key] = values
	}

	var start string = "you"
	var end string = "out"

	memo := make(map[string]int)
	var path int = dfs(graph, start, end, memo)

	fmt.Println(path)

}

func dfs(
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
		total += dfs(graph, next, end, memo)
	}

	memo[node] = total
	return total
}
