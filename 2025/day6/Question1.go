package main

import (
	"bufio"
	"fmt"
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

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 50*1024*1024)

	var data []string
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	symbols := data[len(data)-1]
	lines := data[:len(data)-1]

	total := 0

	for i, s := range symbols {
		if s != '+' && s != '-' && s != '*' && s != '/' {
			continue
		}

		var nums []int

		for _, line := range lines {
			if i >= len(line) {
				continue
			}

			end := len(line)
			for j := i + 1; j < len(symbols); j++ {
				if symbols[j] == '+' || symbols[j] == '-' ||
					symbols[j] == '*' || symbols[j] == '/' {
					end = j
					break
				}
			}

			part := strings.TrimSpace(line[i:end])
			if part == "" {
				continue
			}

			n, err := strconv.Atoi(part)
			if err != nil {
				panic(err)
			}
			nums = append(nums, n)
		}

		if len(nums) == 0 {
			continue
		}

		result := nums[0]
		switch s {
		case '+':
			result = 0
			for _, n := range nums {
				result += n
			}
		case '*':
			result = 1
			for _, n := range nums {
				result *= n
			}
		case '-':
			for _, n := range nums[1:] {
				result -= n
			}
		case '/':
			for _, n := range nums[1:] {
				result /= n
			}
		}

		total += result
	}

	fmt.Println(total)
}
