package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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
			pos -= val
		} else if dir == 'R' {
			pos += val
		} else {
			if err == io.EOF {
				break
			}
			continue
		}

		pos %= 100
		if pos < 0 {
			pos += 100
		}

		if pos == 0 {
			zeroCount++
		}

		if err == io.EOF {
			break
		}
	}

	fmt.Println("Zero count:", zeroCount)
}
