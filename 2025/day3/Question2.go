package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func best12(line string) string {
	const k = 12
	n := len(line)
	if n < k {
		return ""
	}

	drop := n - k

	stack := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		c := line[i]
		for drop > 0 && len(stack) > 0 && stack[len(stack)-1] < c {
			stack = stack[:len(stack)-1]
			drop--
		}
		stack = append(stack, c)
	}

	if drop > 0 {
		stack = stack[:len(stack)-drop]
	}
	return string(stack[:k])
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	total := big.NewInt(0)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 12 {
			continue
		}
		s := best12(line)

		x := new(big.Int)
		x.SetString(s, 10)
		total.Add(total, x)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println(total.String())
}
