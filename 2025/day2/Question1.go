package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	content, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	input := strings.TrimSpace(string(content))
	ranges := strings.Split(input, ",")

	var sum int64

	for _, rangeStr := range ranges {
		rangeStr = strings.TrimSpace(rangeStr)
		if rangeStr == "" {
			continue
		}

		parts := strings.Split(rangeStr, "-")
		if len(parts) != 2 {
			continue
		}

		start, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
		if err != nil {
			continue
		}

		end, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			continue
		}

		for num := start; num <= end; num++ {
			if isInvalid(num) {
				sum += num
			}
		}
	}

	fmt.Println(sum)
}

func isInvalid(n int64) bool {
	s := strconv.FormatInt(n, 10)
	length := len(s)

	if length%2 != 0 {
		return false
	}

	half := length / 2
	firstHalf := s[:half]
	secondHalf := s[half:]

	return firstHalf == secondHalf
}
