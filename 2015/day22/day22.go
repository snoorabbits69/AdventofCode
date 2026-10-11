package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

type State struct {
	bossHP   int16
	hp       int16
	mana     int16
	spent    int16
	shield   uint8
	poison   uint8
	recharge uint8
}

type solver struct {
	best   int16
	floor  int16
	damage int16
	hard   bool
}

func lowerBound(bossHP int16, poison uint8) int16 {
	rem := int32(bossHP) - 3*int32(poison)
	if rem <= 0 {
		return 0
	}
	return int16((rem*173 + 17) / 18)
}

func (s *State) tick() {
	if s.shield > 0 {
		s.shield--
	}
	if s.poison > 0 {
		s.bossHP -= 3
		s.poison--
	}
	if s.recharge > 0 {
		s.mana += 101
		s.recharge--
	}
}

func (s *State) tickArmor() int16 {
	var armor int16
	if s.shield > 0 {
		armor = 7
	}
	s.tick()
	return armor
}

func (s *solver) after(n *State) {
	if n.bossHP <= 0 {
		s.best = n.spent
		return
	}
	armor := n.tickArmor()
	if n.bossHP <= 0 {
		s.best = n.spent
		return
	}
	atk := s.damage - armor
	if atk < 1 {
		atk = 1
	}
	n.hp -= atk
	if n.hp <= 0 {
		return
	}
	if n.spent+lowerBound(n.bossHP, n.poison) >= s.best {
		return
	}
	s.dfs(n)
}

func (s *solver) dfs(st *State) {
	if s.hard {
		st.hp--
		if st.hp <= 0 {
			return
		}
	}
	st.tick()
	if st.bossHP <= 0 {
		if st.spent < s.best {
			s.best = st.spent
		}
		return
	}

	if st.poison == 0 && st.mana >= 173 && st.spent+173 < s.best {
		n := *st
		n.mana -= 173
		n.spent += 173
		n.poison = 6
		s.after(&n)
		if s.best <= s.floor {
			return
		}
	}

	if st.mana >= 53 && st.spent+53 < s.best {
		n := *st
		n.mana -= 53
		n.spent += 53
		n.bossHP -= 4
		s.after(&n)
		if s.best <= s.floor {
			return
		}
	}

	if st.recharge == 0 && st.mana >= 229 && st.spent+229 < s.best {
		n := *st
		n.mana -= 229
		n.spent += 229
		n.recharge = 5
		s.after(&n)
		if s.best <= s.floor {
			return
		}
	}

	if st.shield == 0 && st.mana >= 113 && st.spent+113 < s.best {
		n := *st
		n.mana -= 113
		n.spent += 113
		n.shield = 6
		s.after(&n)
		if s.best <= s.floor {
			return
		}
	}

	if st.mana >= 73 && st.spent+73 < s.best {
		n := *st
		n.mana -= 73
		n.spent += 73
		n.bossHP -= 2
		n.hp += 2
		s.after(&n)
	}
}

func solve(boss [2]int16, hard bool, floor int16) int16 {
	s := solver{best: 1<<15 - 1, floor: floor, damage: boss[1], hard: hard}
	s.dfs(&State{bossHP: boss[0], hp: 50, mana: 500})
	return s.best
}

func parse() [2]int16 {
	var nums [2]int16
	var count int16
	i, n := 0, len(data)

	for i < n && count < 2 {
		if data[i] >= '0' && data[i] <= '9' {
			v := 0
			for i < n && data[i] >= '0' && data[i] <= '9' {
				v = v*10 + int(data[i]-'0')
				i++
			}
			nums[count] = int16(v)
			count++
		} else {
			i++
		}
	}
	return nums
}

func main() {
	start := time.Now()
	const iters = 1
	boss := parse()
	var part1, part2 int16
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1 = solve(boss, false, 0)
		part2 = solve(boss, true, part1)
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
