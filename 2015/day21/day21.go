package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

type Item struct {
	Cost, Damage, Armor int
}

var weapons = map[int]Item{
	1: {8, 4, 0},
	2: {10, 5, 0},
	3: {25, 6, 0},
	4: {40, 7, 0},
	5: {74, 8, 0},
}

var armors = map[int]Item{
	1: {13, 0, 1},
	2: {31, 0, 2},
	3: {53, 0, 3},
	4: {75, 0, 4},
	5: {102, 0, 5},
}

var rings = map[int]Item{
	1: {25, 1, 0},
	2: {50, 2, 0},
	3: {100, 3, 0},
	4: {20, 0, 1},
	5: {40, 0, 2},
	6: {80, 0, 3},
}

func parse() [3]int {
	var nums [3]int
	count := 0
	i, n := 0, len(data)

	for i < n && count < 3 {
		if data[i] >= '0' && data[i] <= '9' {
			v := 0
			for i < n && data[i] >= '0' && data[i] <= '9' {
				v = v*10 + int(data[i]-'0')
				i++
			}
			nums[count] = v
			count++
		} else {
			i++
		}
	}

	return nums
}

func playerWins(boss [3]int, damage, armor int) bool {
	playerHP, bossHP := 100, boss[0]
	hitBoss := max(damage-boss[2], 1)
	hitPlayer := max(boss[1]-armor, 1)
	for {
		bossHP -= hitBoss
		if bossHP <= 0 {
			return true
		}
		playerHP -= hitPlayer
		if playerHP <= 0 {
			return false
		}
	}

}

func solve(boss [3]int) (int, int) {
	cheapestWin := 1 << 30
	priciestLoss := 0

	for w := 1; w <= 5; w++ {
		for a := 0; a <= 5; a++ {
			for r1 := 0; r1 <= 6; r1++ {
				for r2 := 0; r2 <= 6; r2++ {

					if r1 != 0 && r1 == r2 {
						continue
					}

					gear := []Item{weapons[w], armors[a], rings[r1], rings[r2]}

					var total Item
					for _, g := range gear {
						total.Cost += g.Cost
						total.Damage += g.Damage
						total.Armor += g.Armor
					}

					if playerWins(boss, total.Damage, total.Armor) {
						cheapestWin = min(cheapestWin, total.Cost)
					} else {
						priciestLoss = max(priciestLoss, total.Cost)
					}

				}
			}
		}
	}
	return cheapestWin, priciestLoss
}

func main() {
	start := time.Now()
	const iters = 10000
	boss := parse()
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(boss)
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
