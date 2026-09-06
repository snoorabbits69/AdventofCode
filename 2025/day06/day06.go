package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//go:embed input.txt
var data []byte

func parseNumber(data []byte, i *int, n int) int64 {
	var v int64

	for *i < n {
		c := data[*i]

		if c < '0' || c > '9' {
			break
		}

		v = v*10 + int64(c-'0')
		*i++
	}

	return v
}

func parseData() ([]int64, []byte, []int) {
	var numbers []int64
	var symbols []byte
	var rowStarts []int

	var numbercount int = 0

	n := len(data)
	i := 0

	rowStarts = append(rowStarts, 0)

	for i < n {
		c := data[i]

		if c == '\n' {
			rowStarts = append(rowStarts, numbercount)
			i++
			continue
		}

		if c >= '0' && c <= '9' {
			v := parseNumber(data, &i, n)
			numbers = append(numbers, v)
			numbercount++
			continue
		}

		if c != ' ' {
			symbols = append(symbols, c)
		}

		i++
	}
	return numbers, symbols, rowStarts

}

func numLen(a int64) int64 {
	var l int64 = 0
	for a > 0 {
		a /= 10
		l++
	}
	return l
}

func getDigit(a int64, pos int64, totalLen int64) int64 {
	var shift int64 = totalLen - pos - 1
	var i int64 = 0
	for i < shift {
		a /= 10
		i++
	}
	return a % 10
}

func digitsStart(a, b, c, d int64) []int64 {
	maxLen := numLen(a)
	if l := numLen(b); l > maxLen {
		maxLen = l
	}
	if l := numLen(c); l > maxLen {
		maxLen = l
	}
	if l := numLen(d); l > maxLen {
		maxLen = l
	}

	var results []int64

	var i int64 = 0
	var result int64
	var factor int64

	for i < maxLen {
		dA := getDigit(a, i, numLen(a))
		dB := getDigit(b, i, numLen(b))
		dC := getDigit(c, i, numLen(c))
		dD := getDigit(d, i, numLen(d))

		result = 0
		factor = 1
		if i < numLen(d) {
			result += dD * factor
			factor *= 10

		}
		if i < numLen(c) {
			result += dC * factor
			factor *= 10
		}
		if i < numLen(b) {
			result += dB * factor
			factor *= 10
		}
		if i < numLen(a) {
			result += dA * factor
		}
		results = append(results, result)
		i++
	}
	return results
}

func applyOp(symbol byte, a int64, b int64, c int64, d int64) int64 {

	switch symbol {
	case '+':
		return a + b + c + d
	case '*':
		return a * b * c * d
	default:
		return 0
	}

}

func applyOpArr(symbol byte, result []int64) int64 {
	switch symbol {
	case '+':
		var total int64
		for _, r := range result {
			total += r
		}
		return total

	case '*':
		var total int64
		for _, r := range result {
			total *= r
		}
		return total
	default:
		return 0
	}
}

func solvePart1(numbers []int64, symbols []byte, rowStarts []int) int64 {
	var row1 = numbers[rowStarts[0]:rowStarts[1]]
	var row2 = numbers[rowStarts[1]:rowStarts[2]]
	var row3 = numbers[rowStarts[2]:rowStarts[3]]
	var row4 = numbers[rowStarts[3]:rowStarts[4]]

	var part1 int64
	for i, symbol := range symbols {
		part1 += applyOp(symbol, row1[i], row2[i], row3[i], row4[i])
	}

	return part1
}

func solvePart2(data []byte) int64 {
	text := strings.TrimRight(string(data), "\n")
	if text == "" {
		return 0
	}
	lines := strings.Split(text, "\n")
	if len(lines) < 2 {
		return 0
	}

	symbols := lines[len(lines)-1]
	rows := lines[:len(lines)-1]

	width := len(symbols)
	for _, line := range rows {
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
	for i := range rows {
		rows[i] = pad(rows[i])
	}

	isOp := func(b byte) bool { return b == '+' || b == '-' || b == '*' || b == '/' }

	var total int64

	var nums []int64
	var op byte = 0
	inProblem := false

	flush := func() {
		if !inProblem || len(nums) == 0 {
			nums = nil
			op = 0
			inProblem = false
			return
		}

		var result int64
		switch op {
		case '+':
			for _, v := range nums {
				result += v
			}
		case '*':
			result = 1
			for _, v := range nums {
				result *= v
			}
		case '-':
			result = nums[0]
			for _, v := range nums[1:] {
				result -= v
			}
		case '/':
			result = nums[0]
			for _, v := range nums[1:] {
				result /= v
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
		if allBlank {
			for _, line := range rows {
				if line[col] != ' ' {
					allBlank = false
					break
				}
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
		for _, line := range rows {
			r := rune(line[col])
			if unicode.IsDigit(r) {
				sb.WriteRune(r)
			}
		}
		if sb.Len() > 0 {
			n, err := strconv.ParseInt(sb.String(), 10, 64)
			if err != nil {
				panic(err)
			}
			nums = append(nums, n)
		}
	}

	flush()

	return total
}

func solve(numbers []int64, symbols []byte, rowStarts []int) (int64, int64) {
	part1 := solvePart1(numbers, symbols, rowStarts)
	part2 := solvePart2(data)
	return part1, part2
}

func main() {
	start := time.Now()
	const iters = 10000
	numbers, symbols, rowStarts := parseData()
	var part1, part2 int64
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(numbers, symbols, rowStarts)
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
