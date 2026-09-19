package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

func min(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

func minPerimeter(a, b, c int) int {
	if a > b {
		a, b = b, a
	}
	if b > c {
		b, c = c, b
	}
	if a > b {
		a, b = b, a
	}
	return 2 * (a + b)

}

func solve(data []byte) (int, int) {
	surfaceSum := 0
	slackSum := 0
	totalRibbon := 0

	i := 0
	for i < len(data) {
		l := 0
		for i < len(data) && data[i] >= '0' && data[i] <= '9' {
			l = l*10 + int(data[i]-'0')
			i++
		}
		i++

		w := 0
		for i < len(data) && data[i] >= '0' && data[i] <= '9' {
			w = w*10 + int(data[i]-'0')
			i++
		}
		i++

		h := 0
		for i < len(data) && data[i] >= '0' && data[i] <= '9' {
			h = h*10 + int(data[i]-'0')
			i++
		}

		for i < len(data) && (data[i] == '\n' || data[i] == '\r') {
			i++
		}

		lw := l * w
		wh := w * h
		hl := h * l

		slack := min(lw, wh, hl)
		surfaceSum += lw + wh + hl
		slackSum += slack

		perimeter := minPerimeter(l, w, h)
		volume := lw * h
		totalRibbon += perimeter + volume
	}

	totalPaper := 2*surfaceSum + slackSum
	return totalPaper, totalRibbon
}

func main() {
	start := time.Now()
	const iters = 10000
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(data)
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
