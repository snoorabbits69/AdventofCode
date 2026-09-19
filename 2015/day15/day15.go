package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

type Ingredient [5]int

func parseNumber(i *int) int {
	neg := false
	if data[*i] == '-' {
		neg = true
		*i++
	}
	v := 0
	for *i < len(data) {
		c := data[*i]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int(c-'0')
		*i++
	}
	if neg {
		v = -v
	}
	return v
}

func skipNonDigit(i *int) {
	for *i < len(data) {
		c := data[*i]
		if c == '-' || (c >= '0' && c <= '9') {
			return
		}
		*i++
	}
}

func parse() []Ingredient {
	n := len(data)
	i := 0
	var recipe []Ingredient

	for i < n {

		for data[i] != ':' {
			i++
		}
		i++

		var ing Ingredient
		for field := 0; field < 5; field++ {
			skipNonDigit(&i)
			ing[field] = parseNumber(&i)
		}
		recipe = append(recipe, ing)

		for i < n && data[i] != '\n' {
			i++
		}
		if i < n {
			i++
		}
	}

	return recipe
}

func solve(recipe []Ingredient) (int, int) {
	partOne, partTwo := 0, 0

	for a := 0; a <= 100; a++ {
		var first Ingredient
		for k := 0; k < 5; k++ {
			first[k] = a * recipe[0][k]
		}

	outer:
		for b := 0; b <= 100-a; b++ {
			var second Ingredient
			for k := 0; k < 5; k++ {
				second[k] = first[k] + b*recipe[1][k]
			}

			for k := 0; k < 4; k++ {
				y, z := recipe[2][k], recipe[3][k]
				m := y
				if z > m {
					m = z
				}
				if second[k]+m*(100-a-b) <= 0 {
					continue outer
				}
			}

			for c := 0; c <= 100-a-b; c++ {
				d := 100 - a - b - c
				var third, fourth Ingredient
				for k := 0; k < 5; k++ {
					third[k] = second[k] + c*recipe[2][k]
					fourth[k] = third[k] + d*recipe[3][k]
				}

				score := 1
				for k := 0; k < 4; k++ {
					v := fourth[k]
					if v < 0 {
						v = 0
					}
					score *= v
				}
				calories := fourth[4]

				if score > partOne {
					partOne = score
				}
				if calories == 500 && score > partTwo {
					partTwo = score
				}
			}
		}
	}

	return partOne, partTwo
}

func main() {
	start := time.Now()
	const iters = 100
	recipe := parse()
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(recipe)
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
