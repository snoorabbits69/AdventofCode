package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const MaxShapes = 6

type Tree struct {
	w, h     int
	presents [6]uint32
}

func hasSufficientArea(t Tree) bool {
	area := uint32(t.w * t.h)

	total := uint32(0)
	for _, p := range t.presents {
		total += p
	}

	return area >= total*9
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	trees := make([]Tree, 0, 1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		t, ok := parseTreeLine(line)
		if ok {
			trees = append(trees, t)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	areaFeasibleCount := uint32(0)
	for _, t := range trees {
		if hasSufficientArea(t) {
			areaFeasibleCount++
		}
	}

	fmt.Println(areaFeasibleCount)
}

func parseTreeLine(line string) (Tree, bool) {
	colon := strings.IndexByte(line, ':')
	if colon < 0 {
		return Tree{}, false
	}

	left := strings.TrimSpace(line[:colon])
	right := strings.TrimSpace(line[colon+1:])

	x := strings.IndexByte(left, 'x')
	if x < 0 {
		return Tree{}, false
	}

	w, err1 := strconv.Atoi(strings.TrimSpace(left[:x]))
	h, err2 := strconv.Atoi(strings.TrimSpace(left[x+1:]))
	if err1 != nil || err2 != nil || w <= 0 || h <= 0 {
		return Tree{}, false
	}

	fields := strings.Fields(right)
	if len(fields) < MaxShapes {
		return Tree{}, false
	}

	var pres [MaxShapes]uint32
	for i := 0; i < MaxShapes; i++ {
		v, err := strconv.ParseUint(fields[i], 10, 32)
		if err != nil {
			return Tree{}, false
		}
		pres[i] = uint32(v)
	}

	return Tree{w: w, h: h, presents: pres}, true
}
