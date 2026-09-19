package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

var DX = [256]int16{'>': 1, '<': -1}
var DY = [256]int16{'^': 1, 'v': -1}

const GRID = 256
const OFF = GRID / 2

func idx(x, y int16) int {
	return ((int(y) + OFF) << 8) | (int(x) + OFF)
}

func solve(data []byte) (int, int) {
	var vis1 [GRID * GRID]uint8
	var vis2 [GRID * GRID]uint8

	var x1, y1 int16
	var sx, sy int16
	var rx, ry int16

	start := idx(0, 0)
	vis1[start] = 1
	vis2[start] = 1

	count1, count2 := 1, 1

	for i := 0; i < len(data); i++ {
		c := data[i]
		dx := DX[c]
		dy := DY[c]

		x1 += dx
		y1 += dy
		p1 := idx(x1, y1)
		if vis1[p1] == 0 {
			vis1[p1] = 1
			count1++
		}

		if i&1 == 0 {
			sx += dx
			sy += dy
			p2 := idx(sx, sy)
			if vis2[p2] == 0 {
				vis2[p2] = 1
				count2++
			}
		} else {
			rx += dx
			ry += dy
			p2 := idx(rx, ry)
			if vis2[p2] == 0 {
				vis2[p2] = 1
				count2++
			}
		}
	}

	return count1, count2
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
