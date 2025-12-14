package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}

		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			continue
		}

		start, err1 := strconv.ParseInt(parts[0], 10, 64)
		end, err2 := strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		if start > end {
			start, end = end, start
		}

		ranges = append(ranges, Range{Start: start, End: end})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fresh := countFresh(ranges)
	fmt.Println("Fresh count:", fresh)
}

func countFresh(ranges []Range) int64 {
	if len(ranges) == 0 {
		return 0
	}

	sort.Slice(ranges, func(i, j int) bool {
		if ranges[i].Start == ranges[j].Start {
			return ranges[i].End < ranges[j].End
		}
		return ranges[i].Start < ranges[j].Start
	})

	cur := ranges[0]
	var total int64 = 0

	for i := 1; i < len(ranges); i++ {
		r := ranges[i]
		if r.Start <= cur.End+1 {
			if r.End > cur.End {
				cur.End = r.End
			}
		} else {
			total += cur.End - cur.Start + 1
			cur = r
		}
	}

	total += cur.End - cur.Start + 1
	return total
}
