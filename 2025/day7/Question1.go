package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 50*1024*1024)

	var grid []string
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if len(grid) == 0 {
		fmt.Println(0)
		return
	}

	h := len(grid)
	w := len(grid[0])

	start := Point{0, 0}
	found := false
	for y := 0; y < h && !found; y++ {
		for x := 0; x < w; x++ {
			if grid[y][x] == 'S' {
				start = Point{x, y}
				found = true
				break
			}
		}
	}

	inBounds := func(p Point) bool {
		return p.x >= 0 && p.x < w && p.y >= 0 && p.y < h
	}

	queue := []Point{start}
	visited := map[Point]bool{start: true}
	splits := 0

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		cell := grid[cur.y][cur.x]
		if cell == '^' {
			splits++

			left := Point{cur.x - 1, cur.y}
			if inBounds(left) && !visited[left] {
				visited[left] = true
				queue = append(queue, left)
			}

			right := Point{cur.x + 1, cur.y}
			if inBounds(right) && !visited[right] {
				visited[right] = true
				queue = append(queue, right)
			}
		} else {
			down := Point{cur.x, cur.y + 1}
			if inBounds(down) && !visited[down] {
				visited[down] = true
				queue = append(queue, down)
			}
		}
	}

	fmt.Println(splits)
}
