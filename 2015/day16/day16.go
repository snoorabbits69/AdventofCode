package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

type sue struct {
	num   int
	facts map[string]int
}

var target = map[string]int{
	"children":    3,
	"cats":        7,
	"samoyeds":    2,
	"pomeranians": 3,
	"akitas":      0,
	"vizslas":     0,
	"goldfish":    5,
	"trees":       3,
	"cars":        2,
	"perfumes":    1,
}

func parse() []sue {
	var sues []sue
	i, n := 0, len(data)

	for i < n {
		for i < n && data[i] != ' ' {
			i++
		}
		i++
		num := 0
		for i < n && data[i] >= '0' && data[i] <= '9' {
			num = num*10 + int(data[i]-'0')
			i++
		}

		for i < n && data[i] != ' ' {
			i++
		}
		i++

		facts := make(map[string]int, 3)

		for i < n && data[i] != '\n' {
			start := i
			for i < n && data[i] != ':' {
				i++
			}
			key := string(data[start:i])
			i++
			if i < n && data[i] == ' ' {
				i++
			}

			val := 0
			for i < n && data[i] >= '0' && data[i] <= '9' {
				val = val*10 + int(data[i]-'0')
				i++
			}
			facts[key] = val

			if i < n && data[i] == ',' {
				i++
			}
			if i < n && data[i] == ' ' {
				i++
			}
		}

		sues = append(sues, sue{num: num, facts: facts})

		if i < n && data[i] == '\n' {
			i++
		}
	}

	return sues
}

func solve(sues []sue, predicate func(key string, value int) bool) int {
	for _, s := range sues {
		matches := true
		for key, val := range s.facts {
			if !predicate(key, val) {
				matches = false
				break
			}
		}
		if matches {
			return s.num
		}
	}
	return -1
}

func part1(sues []sue) int {
	return solve(sues, func(key string, value int) bool {
		return value == target[key]
	})
}

func part2(sues []sue) int {
	return solve(sues, func(key string, value int) bool {
		switch key {
		case "cats", "trees":
			return value > target[key]
		case "pomeranians", "goldfish":
			return value < target[key]
		default:
			return value == target[key]
		}
	})
}

func main() {
	start := time.Now()
	const iters = 10000
	sues := parse()
	part1Fn, part2Fn := part1, part2
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1 = part1Fn(sues)
		part2 = part2Fn(sues)
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
