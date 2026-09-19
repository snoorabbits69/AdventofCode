package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

var (
	ops  []byte
	regs []int
	offs []int
)

func parse() {
	n := len(data)

	lines := 0
	for _, c := range data {
		if c == '\n' {
			lines++
		}
	}
	if n > 0 && data[n-1] != '\n' {
		lines++
	}

	ops = make([]byte, 0, lines)
	regs = make([]int, 0, lines)
	offs = make([]int, 0, lines)

	i := 0
	for i < n {
		var op byte
		reg := 0
		off := 0

		switch data[i] {
		case 'h', 't', 'i':
			op = data[i]
		case 'j':
			if data[i+1] == 'm' {
				op = 'm'
			} else {
				op = data[i+2]
			}
		}
		i += 4

		if op != 'm' {
			reg = int(data[i] - 'a')
			i++
		}

		if op == 'e' || op == 'o' {
			i += 2
		}

		if op == 'm' || op == 'e' || op == 'o' {
			sign := 1
			if data[i] == '-' {
				sign = -1
			}
			i++

			num := 0
			for i < n && data[i] >= '0' && data[i] <= '9' {
				num = num*10 + int(data[i]-'0')
				i++
			}
			off = sign * num
		}

		for i < n && data[i] != '\n' {
			i++
		}
		i++

		ops = append(ops, op)
		regs = append(regs, reg)
		offs = append(offs, off)
	}
}

func run(a, b int) int {
	r := [2]int{a, b}
	pc := 0

	for pc >= 0 && pc < len(ops) {
		switch ops[pc] {
		case 'h':
			r[regs[pc]] /= 2
			pc++
		case 't':
			r[regs[pc]] *= 3
			pc++
		case 'i':
			r[regs[pc]]++
			pc++
		case 'm':
			pc += offs[pc]
		case 'e':
			if r[regs[pc]]%2 == 0 {
				pc += offs[pc]
			} else {
				pc++
			}
		case 'o':
			if r[regs[pc]] == 1 {
				pc += offs[pc]
			} else {
				pc++
			}
		}
	}

	return r[1]
}

func main() {
	start := time.Now()
	const iters = 10000
	parse()
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1 = run(0, 0)
		part2 = run(1, 0)
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
