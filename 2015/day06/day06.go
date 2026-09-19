package main

import (
	_ "embed"
	"fmt"
	"time"
	"unsafe"
)

//go:embed input.txt
var data []byte

const (
	W = 1000
	N = W * W
)

var (
	grid1 [N]uint8
	grid2 [N]uint8
)

//go:noescape
//go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)

func sumBytes(data []uint8) uint64 {
	var sum uint64
	i := 0
	for ; i+7 < len(data); i += 8 {
		sum += uint64(data[i]) + uint64(data[i+1]) + uint64(data[i+2]) + uint64(data[i+3]) +
			uint64(data[i+4]) + uint64(data[i+5]) + uint64(data[i+6]) + uint64(data[i+7])
	}
	for ; i < len(data); i++ {
		sum += uint64(data[i])
	}
	return sum
}

func solve(input []byte) (int, int) {
	memclrNoHeapPointers(unsafe.Pointer(&grid1[0]), N)
	memclrNoHeapPointers(unsafe.Pointer(&grid2[0]), N)

	i := 0
	for i < len(input) && input[i] == 't' {
		if i+6 >= len(input) {
			break
		}

		op := uint8(2)
		if input[i+6] == 'n' {
			op = 0
			i += 8
		} else if input[i+6] == 'f' {
			op = 1
			i += 9
		} else {
			i += 7
		}

		if i >= len(input) {
			break
		}

		x1 := 0
		for i < len(input) && input[i] != ',' {
			x1 = x1*10 + int(input[i]-'0')
			i++
		}
		i++
		if i >= len(input) {
			break
		}

		y1 := 0
		for i < len(input) && input[i] != ' ' {
			y1 = y1*10 + int(input[i]-'0')
			i++
		}
		i += 9
		if i >= len(input) {
			break
		}

		x2 := 0
		for i < len(input) && input[i] != ',' {
			x2 = x2*10 + int(input[i]-'0')
			i++
		}
		i++
		if i > len(input) {
			break
		}

		y2 := 0
		for i < len(input) && input[i] >= '0' && input[i] <= '9' {
			y2 = y2*10 + int(input[i]-'0')
			i++
		}

		for i < len(input) && (input[i] == '\n' || input[i] == '\r') {
			i++
		}

		width := x2 - x1 + 1
		height := y2 - y1 + 1

		for y := 0; y < height; y++ {
			base := (y1+y)*W + x1
			x := 0
			for ; x+15 < width; x += 16 {
				idx := base + x
				if op == 0 {
					grid1[idx] = 1
					grid1[idx+1] = 1
					grid1[idx+2] = 1
					grid1[idx+3] = 1
					grid1[idx+4] = 1
					grid1[idx+5] = 1
					grid1[idx+6] = 1
					grid1[idx+7] = 1
					grid1[idx+8] = 1
					grid1[idx+9] = 1
					grid1[idx+10] = 1
					grid1[idx+11] = 1
					grid1[idx+12] = 1
					grid1[idx+13] = 1
					grid1[idx+14] = 1
					grid1[idx+15] = 1
					grid2[idx]++
					grid2[idx+1]++
					grid2[idx+2]++
					grid2[idx+3]++
					grid2[idx+4]++
					grid2[idx+5]++
					grid2[idx+6]++
					grid2[idx+7]++
					grid2[idx+8]++
					grid2[idx+9]++
					grid2[idx+10]++
					grid2[idx+11]++
					grid2[idx+12]++
					grid2[idx+13]++
					grid2[idx+14]++
					grid2[idx+15]++
				} else if op == 1 {
					grid1[idx] = 0
					grid1[idx+1] = 0
					grid1[idx+2] = 0
					grid1[idx+3] = 0
					grid1[idx+4] = 0
					grid1[idx+5] = 0
					grid1[idx+6] = 0
					grid1[idx+7] = 0
					grid1[idx+8] = 0
					grid1[idx+9] = 0
					grid1[idx+10] = 0
					grid1[idx+11] = 0
					grid1[idx+12] = 0
					grid1[idx+13] = 0
					grid1[idx+14] = 0
					grid1[idx+15] = 0
					if grid2[idx] > 0 {
						grid2[idx]--
					}
					if grid2[idx+1] > 0 {
						grid2[idx+1]--
					}
					if grid2[idx+2] > 0 {
						grid2[idx+2]--
					}
					if grid2[idx+3] > 0 {
						grid2[idx+3]--
					}
					if grid2[idx+4] > 0 {
						grid2[idx+4]--
					}
					if grid2[idx+5] > 0 {
						grid2[idx+5]--
					}
					if grid2[idx+6] > 0 {
						grid2[idx+6]--
					}
					if grid2[idx+7] > 0 {
						grid2[idx+7]--
					}
					if grid2[idx+8] > 0 {
						grid2[idx+8]--
					}
					if grid2[idx+9] > 0 {
						grid2[idx+9]--
					}
					if grid2[idx+10] > 0 {
						grid2[idx+10]--
					}
					if grid2[idx+11] > 0 {
						grid2[idx+11]--
					}
					if grid2[idx+12] > 0 {
						grid2[idx+12]--
					}
					if grid2[idx+13] > 0 {
						grid2[idx+13]--
					}
					if grid2[idx+14] > 0 {
						grid2[idx+14]--
					}
					if grid2[idx+15] > 0 {
						grid2[idx+15]--
					}
				} else {
					grid1[idx] ^= 1
					grid1[idx+1] ^= 1
					grid1[idx+2] ^= 1
					grid1[idx+3] ^= 1
					grid1[idx+4] ^= 1
					grid1[idx+5] ^= 1
					grid1[idx+6] ^= 1
					grid1[idx+7] ^= 1
					grid1[idx+8] ^= 1
					grid1[idx+9] ^= 1
					grid1[idx+10] ^= 1
					grid1[idx+11] ^= 1
					grid1[idx+12] ^= 1
					grid1[idx+13] ^= 1
					grid1[idx+14] ^= 1
					grid1[idx+15] ^= 1
					grid2[idx] += 2
					grid2[idx+1] += 2
					grid2[idx+2] += 2
					grid2[idx+3] += 2
					grid2[idx+4] += 2
					grid2[idx+5] += 2
					grid2[idx+6] += 2
					grid2[idx+7] += 2
					grid2[idx+8] += 2
					grid2[idx+9] += 2
					grid2[idx+10] += 2
					grid2[idx+11] += 2
					grid2[idx+12] += 2
					grid2[idx+13] += 2
					grid2[idx+14] += 2
					grid2[idx+15] += 2
				}
			}
			for ; x < width; x++ {
				idx := base + x
				if op == 0 {
					grid1[idx] = 1
					grid2[idx]++
				} else if op == 1 {
					grid1[idx] = 0
					if grid2[idx] > 0 {
						grid2[idx]--
					}
				} else {
					grid1[idx] ^= 1
					grid2[idx] += 2
				}
			}
		}
	}

	p1 := int(sumBytes(grid1[:]))
	p2 := int(sumBytes(grid2[:]))
	return p1, p2
}

func main() {
	_ = unsafe.Sizeof(0)
	start := time.Now()
	const iters = 100
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
