package main

import "fmt"

func lookAndSay(input string) string {
	output := make([]byte, 0, len(input)*2)

	for i := 0; i < len(input); {
		count := byte(1)

		for i+int(count) < len(input) && input[i] == input[i+int(count)] {
			count++
		}

		output = append(output, '0'+count, input[i])
		i += int(count)
	}

	return string(output)
}

func solve(input string) (int, int) {
	part1 := 0
	part2 := 0

	for i := 1; i <= 50; i++ {
		input = lookAndSay(input)

		if i == 40 {
			part1 = len(input)
		}
		if i == 50 {
			part2 = len(input)
		}
	}

	return part1, part2
}

func main() {
	input := "3113322113"

	part1, part2 := solve(input)
	fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", part2)
}
