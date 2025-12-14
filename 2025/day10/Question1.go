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

	var indicators []int

	var buttons [][]int

	for scanner.Scan() {
		line := scanner.Text()

		{
			start := strings.Index(line, "[")
			end := strings.Index(line, "]")
			indicators = append(indicators, parseIndicator(line[start+1:end]))
		}

		{
			start := strings.Index(line, "(")
			end := strings.Index(line, "{")

			content := line[start:end]
			groups := strings.Fields(content)
			var button []int

			for _, group := range groups {
				cleanedbtn := parseIntsfromContent(group)

				button = append(button, int(buttonsToMask(cleanedbtn)))

			}
			buttons = append(buttons, button)
		}

	}

	for i, _ := range indicators {
		fmt.Println(indicators[i], buttons[i])
	}

	total := 0

	for i, indicator := range indicators {

		buttonMasks := buttons[i]

		numLights := 0
		tmp := indicator

		for tmp > 0 {
			numLights++
			tmp >>= 1
		}

		for _, bm := range buttonMasks {
			t := bm

			bit := 0

			for t > 0 {
				if t&1 == 1 && bit+1 > numLights {
					numLights = bit + 1
				}
				t >>= 1
				bit++
			}
		}
		presses := minBFS(numLights, indicator, buttonMasks)

		if presses < 0 {
			panic("no solution")
		}
		total += presses
	}

	fmt.Println("Total presses:", total)
}
