package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	total := 0

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}

		maxVal := -1
		for i := 0; i < len(line); i++ {
			di := int(line[i] - '0')
			for j := i + 1; j < len(line); j++ {
				dj := int(line[j] - '0')
				val := di*10 + dj
				if val > maxVal {
					maxVal = val
				}
			}
		}

		total += maxVal
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println(total)
}
