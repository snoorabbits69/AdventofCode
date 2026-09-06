// package main

// import (
// 	_ "embed"
// 	"fmt"
// 	"time"
// )

// //go:embed input.txt
// var data []byte

// // Movement arrays for the int32 coordinate version
// var dx = [256]int32{'>': 1, '<': -1}
// var dy = [256]int32{'^': 1, 'v': -1}

// // Function to pack two int32 coordinates into a single int64 key
// func packCoords(x, y int32) int64 {
// 	// Packs x into the high 32 bits and y into the low 32 bits
// 	return int64(x)<<32 | int64(uint32(y))
// }

// func part1(data []byte) int {
// 	n := len(data)
// 	var x, y int32 // int32 coordinates

// 	// Map key is int64
// 	visited := make(map[int64]struct{}, n+1)
// 	visited[packCoords(0, 0)] = struct{}{}

// 	for i := 0; i < n; i++ {
// 		c := data[i]
// 		x += dx[c]
// 		y += dy[c]

// 		key := packCoords(x, y)
// 		visited[key] = struct{}{}
// 	}
// 	return len(visited)
// }

// func part2(data []byte) int {
// 	n := len(data)

// 	// Santa's and Robo-Santa's coordinates are int32
// 	var sx, sy int32 = 0, 0
// 	var rx, ry int32 = 0, 0

// 	// Map key is int64
// 	visited := make(map[int64]struct{}, n+1)
// 	visited[packCoords(0, 0)] = struct{}{}

// 	// Pointer variables to swap between Santa's and Robo-Santa's coordinates
// 	var currentX, currentY *int32

// 	for i := 0; i < n; i++ {
// 		c := data[i]
// 		moveX := dx[c]
// 		moveY := dy[c]

// 		// Pointer setup based on parity (i&1 == 0 is Santa, i&1 == 1 is Robo)
// 		if i&1 == 0 {
// 			currentX = &sx
// 			currentY = &sy
// 		} else {
// 			currentX = &rx
// 			currentY = &ry
// 		}

// 		// Accumulate movement using the dereferenced pointers
// 		*currentX += moveX
// 		*currentY += moveY

// 		// Insert the current coordinate into the map
// 		visited[packCoords(*currentX, *currentY)] = struct{}{}
// 	}
// 	return len(visited)
// }

// func solve(data []byte) (int, int) {
// 	return part1(data), part2(data)
// }

// func main() {
// 	start := time.Now()

// 	var part1, part2 int
// 	// Run once to warm up/ensure correctness
// 	for i := 0; i < 100; i++ {
// 		part1, part2 = solve(data)
// 	}

// 	const iters = 10000
// 	benchStart := time.Now()
// 	for i := 0; i < iters; i++ {
// 		part1, part2 = solve(data)
// 	}

// 	benchTime := time.Since(benchStart)
// 	totalTime := time.Since(start)

// 	fmt.Println("solution of part1:", part1)
// 	fmt.Println("solution of part2:", part2)
// 	fmt.Println()

// 	fmt.Printf("Total bench time:  %d microseconds\n", benchTime.Microseconds())
// 	avgUs := float64(benchTime.Nanoseconds()) / 1_000.0 / float64(iters)
// 	fmt.Printf("Average per iter:  %.4f microseconds\n", avgUs)
// 	fmt.Printf("Total time:        %d microseconds\n", totalTime.Microseconds())
// }
