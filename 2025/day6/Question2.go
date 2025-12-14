package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
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
	if len(data) == 0 {
		return
	}

	symbols := data[len(data)-1]
	lines := data[:len(data)-1]

	width := len(symbols)
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	pad := func(s string) string {
		if len(s) < width {
			return s + strings.Repeat(" ", width-len(s))
		}
		return s
	}
	symbols = pad(symbols)
	for i := range lines {
		lines[i] = pad(lines[i])
	}

	isOp := func(b byte) bool { return b == '+' || b == '-' || b == '*' || b == '/' }

	total := 0

	var nums []int
	var op byte = 0
	inProblem := false

	flush := func() {
		if !inProblem || len(nums) == 0 {
			nums = nil
			op = 0
			inProblem = false
			return
		}

		result := nums[0]
		switch op {
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
		default:
			panic("missing operator in problem block")
		}

		total += result

		nums = nil
		op = 0
		inProblem = false
	}

	for col := width - 1; col >= 0; col-- {
		allBlank := symbols[col] == ' '
		for _, line := range lines {
			if line[col] != ' ' {
				allBlank = false
				break
			}
		}
		if allBlank {
			flush()
			continue
		}

		inProblem = true

		if op == 0 && isOp(symbols[col]) {
			op = symbols[col]
		}

		var sb strings.Builder
		for _, line := range lines {
			r := rune(line[col])
			if unicode.IsDigit(r) {
				sb.WriteRune(r)
			}
		}
		if sb.Len() > 0 {
			n, err := strconv.Atoi(sb.String())
			if err != nil {
				panic(err)
			}
			nums = append(nums, n)
		}
	}

	flush()

	fmt.Println(total)
}
