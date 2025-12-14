package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	f, err := os.Open("input.txt")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	var buttons [][][]int

	var joltings [][]int

	for scanner.Scan() {
		line := scanner.Text()

		{
			start := strings.Index(line, "(")
			end := strings.Index(line, "{")

			content := line[start:end]
			groups := strings.Fields(content)
			var button [][]int

			for _, group := range groups {
				cleanedbtn := parseIntsfromContent(group)

				button = append(button, cleanedbtn)

			}
			buttons = append(buttons, button)
		}

		{
			start := strings.Index(line, "{")
			end := strings.Index(line, "}")

			joltings = append(joltings, parseIntsfromContent(line[start:end+1]))

		}

	}

	var total int = 0

	for i, jolting := range joltings {

		min := minLinearAlgebra(buttons[i], jolting)
		total += min
	}

	fmt.Println(total)

}
