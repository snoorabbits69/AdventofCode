package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	Start int64
	End   int64
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var ranges []Range
	var fresh int32 = 0
	readingRanges := true

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			readingRanges = false
			continue
		}

		if readingRanges {

			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				continue
			}

			start, err1 := strconv.ParseInt(parts[0], 10, 64)
			end, err2 := strconv.ParseInt(parts[1], 10, 64)
			if err1 != nil || err2 != nil {
				continue
			}

			ranges = append(ranges, Range{Start: start, End: end})
		} else {

			value, err := strconv.ParseInt(line, 10, 64)
			if err != nil {
				continue
			}

			if inRange(value, ranges) {
				fresh++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println("Fresh count:", fresh)
}

func inRange(value int64, ranges []Range) bool {
	for _, r := range ranges {
		if value >= r.Start && value <= r.End {
			return true
		}
	}
	return false
}
