package main

import (
	"fmt"
	"time"
)

const n = 8

type searcher struct {
	start [n]byte
	buf   [n]byte
}

func minNeeded(pairs int, straight, consec bool) int {
	if straight {
		return [3]int{3, 1, 0}[pairs]
	}
	if consec {
		return [3]int{4, 2, 1}[pairs]
	}
	return [3]int{4, 3, 2}[pairs]
}

func (s *searcher) dfs(pos int, tight bool, pairs int, pairLetter byte, straight bool) bool {
	if pos == n {
		return pairs == 2 && straight
	}
	lo := byte('a')
	if tight {
		lo = s.start[pos]
	}
	for c := lo; c <= 'z'; c++ {
		if c == 'i' || c == 'o' || c == 'l' {
			continue
		}

		np, npl, ns := pairs, pairLetter, straight
		consec := false
		if pos >= 1 {
			p := s.buf[pos-1]
			consec = p+1 == c

			if c == p && np < 2 && (np == 0 || c != npl) {
				if np == 0 {
					npl = c
				}
				np++
			}
			if !ns && pos >= 2 && consec && s.buf[pos-2]+1 == p {
				ns = true
			}
		}

		if n-pos-1 < minNeeded(np, ns, consec) {
			continue
		}

		s.buf[pos] = c
		if s.dfs(pos+1, tight && c == lo, np, npl, ns) {
			return true
		}
	}
	return false
}

func next(pw [n]byte) [n]byte {

	for i := n - 1; i >= 0; i-- {
		if pw[i] == 'z' {
			pw[i] = 'a'
			continue
		}
		pw[i]++
		break
	}
	s := searcher{start: pw}
	s.dfs(0, true, 0, 0, false)
	return s.buf
}

func solve() (string, string) {
	var in [n]byte
	copy(in[:], "cqjxjnds")
	p1 := next(in)
	p2 := next(p1)
	return string(p1[:]), string(p2[:])
}

func main() {
	start := time.Now()
	const iters = 1000

	var part1, part2 string
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve()
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
