package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"
)

var DX = [256]int16{'>': 1, '<': -1}
var DY = [256]int16{'^': 1, 'v': -1}

//go:embed input.txt
var data []byte

const GRID = 256
const OFF  = GRID / 2

func idx(x, y int16) int {
    return ((int(y)+OFF)<<8) | (int(x) + OFF)
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

	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	part1, part2 := solve(data)

	const iters = 10000
	benchStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(data)
	}
	benchTime := time.Since(benchStart)
	totalTime := time.Since(start)

	fmt.Println("solution of part 1:", part1)
	fmt.Println("solution of part 2:", part2)
	fmt.Println()

	fmt.Printf("Total bench time:  %d microseconds\n", benchTime.Microseconds())
	fmt.Printf("Average per iter:  %.4f microseconds\n",
		float64(benchTime.Nanoseconds())/1000.0/float64(iters))
	fmt.Printf("Total time:        %d microseconds\n", totalTime.Microseconds())
}
