package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func countZeroHitsLeft(pos, k int) int {
	if k <= 0 {
		return 0
	}

	t0 := pos
	if t0 == 0 {
		t0 = 100
	}

	if k < t0 {
		return 0
	}
	return 1 + (k-t0)/100
}

func countZeroHitsRight(pos, k int) int {
	if k <= 0 {
		return 0
	}

	t0 := (100 - pos) % 100
	if t0 == 0 {
		t0 = 100
	}

	if k < t0 {
		return 0
	}
	return 1 + (k-t0)/100
}

func wrap100(x int) int {
	x %= 100
	if x < 0 {
		x += 100
	}
	return x
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReaderSize(f, 1<<20)

	pos := 50
	zeroCount := 0

	for {
		line, err := r.ReadBytes('\n')
		if len(line) == 0 && err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		n := len(line)
		for n > 0 && (line[n-1] == '\n' || line[n-1] == '\r') {
			n--
		}
		if n == 0 {
			if err == io.EOF {
				break
			}
			continue
		}

		dir := line[0]
		val := 0
		for i := 1; i < n; i++ {
			c := line[i]
			if c < '0' || c > '9' {
				break
			}
			val = val*10 + int(c-'0')
		}

		if dir == 'L' {
			zeroCount += countZeroHitsLeft(pos, val)
			pos = wrap100(pos - val)
		} else if dir == 'R' {
			zeroCount += countZeroHitsRight(pos, val)
			pos = wrap100(pos + val)
		}

		if err == io.EOF {
			break
		}
	}

	fmt.Println("Zero count:", zeroCount)
}
