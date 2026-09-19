package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

const (
	hp = iota
	mana
	bossHP
	shield
	poison
	recharge
	spent
)

var costs = [5]int{53, 73, 113, 173, 229}

func parse() [2]int {
	var nums [2]int
	count := 0
	i, n := 0, len(data)

	for i < n && count < 2 {
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

func effects(state *[7]int) int {
	armor := 0
	if state[shield] > 0 {
		armor = 7
		state[shield]--
	}
	if state[poison] > 0 {
		state[bossHP] -= 3
		state[poison]--
	}
	if state[recharge] > 0 {
		state[mana] += 101
		state[recharge]--
	}
	return armor
}

func dfs(state [7]int, damage int, hard bool, best *int) {
	if hard {
		state[hp]--
		if state[hp] <= 0 {
			return
		}
	}
	effects(&state)
	if state[bossHP] <= 0 {
		*best = min(*best, state[spent])
		return
	}

	for spell := 0; spell < 5; spell++ {
		cost := costs[spell]
		if cost > state[mana] || state[spent]+cost >= *best {
			continue
		}

		n := state
		n[mana] -= cost
		n[spent] += cost

		switch spell {
		case 0:
			n[bossHP] -= 4
		case 1:
			n[bossHP] -= 2
			n[hp] += 2
		case 2:
			if n[shield] > 0 {
				continue
			}
			n[shield] = 6
		case 3:
			if n[poison] > 0 {
				continue
			}
			n[poison] = 6
		case 4:
			if n[recharge] > 0 {
				continue
			}
			n[recharge] = 5
		}

		if n[bossHP] <= 0 {
			*best = n[spent]
			continue
		}

		armor := effects(&n)
		if n[bossHP] <= 0 {
			*best = n[spent]
			continue
		}
		n[hp] -= max(1, damage-armor)
		if n[hp] <= 0 {
			continue
		}

		dfs(n, damage, hard, best)
	}
}

func solve(boss [2]int, hard bool) int {
	best := 1 << 30
	var start [7]int
	start[hp] = 50
	start[mana] = 500
	start[bossHP] = boss[0]
	dfs(start, boss[1], hard, &best)
	return best
}

func main() {
	start := time.Now()
	const iters = 100
	boss := parse()
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1 = solve(boss, false)
		part2 = solve(boss, true)
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
