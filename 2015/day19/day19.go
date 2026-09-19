package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"time"
)

//go:embed input.txt
var data []byte

type replacement struct {
	from, to []byte
}

type input struct {
	molecule     []byte
	replacements []replacement
}

func parse() input {

	raw := bytes.TrimRight(data, "\n")

	split := bytes.LastIndex(raw, []byte("\n\n"))
	block, molecule := raw[:split], raw[split+2:]

	var reps []replacement
	for len(block) > 0 {

		line := block
		if i := bytes.IndexByte(block, '\n'); i >= 0 {
			line, block = block[:i], block[i+1:]
		} else {
			block = nil
		}

		i := bytes.Index(line, []byte(" => "))
		reps = append(reps, replacement{from: line[:i], to: line[i+4:]})
	}

	return input{molecule: molecule, replacements: reps}
}

type span struct {
	start, end int
	to         []byte
}

func part1(in input) int {
	mol := in.molecule

	groups := map[int][]replacement{}
	for _, r := range in.replacements {
		extra := len(r.to) - len(r.from)
		groups[extra] = append(groups[extra], r)
	}

	result := 0
	var modified []span

	for _, group := range groups {
		modified = modified[:0]

		for _, r := range group {
			for off := 0; ; {
				i := bytes.Index(mol[off:], r.from)
				if i < 0 {
					break
				}
				start := off + i
				modified = append(modified, span{start, start + len(r.from), r.to})
				off = start + 1
			}
		}

		sort.Slice(modified, func(a, b int) bool {
			return modified[a].start < modified[b].start
		})

	outer:
		for i, m := range modified {
			for _, m2 := range modified[i+1:] {

				if m2.start >= m.start+len(m.to) {
					break
				}

				lenFirst := len(m.to) + len(mol) - m.end
				lenSecond := (m2.start - m.start) + len(m2.to)
				n := min(lenFirst, lenSecond)

				same := true
				for k := 0; k < n; k++ {
					var a, b byte
					if k < len(m.to) {
						a = m.to[k]
					} else {
						a = mol[m.end+k-len(m.to)]
					}
					if k < m2.start-m.start {
						b = mol[m.start+k]
					} else {
						b = m2.to[k-(m2.start-m.start)]
					}
					if a != b {
						same = false
						break
					}
				}

				if same {
					continue outer
				}
			}
			result++
		}
	}

	return result
}

func part2(in input) int {
	mol := in.molecule

	elements := 0
	for _, b := range mol {
		if b >= 'A' && b <= 'Z' {
			elements++
		}
	}

	rn := bytes.Count(mol, []byte("Rn"))
	ar := bytes.Count(mol, []byte("Ar"))
	y := bytes.Count(mol, []byte("Y"))

	return elements - rn - ar - 2*y - 1
}

func main() {
	start := time.Now()
	const iters = 1000
	in := parse()
	part1Fn, part2Fn := part1, part2
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1 = part1Fn(in)
		part2 = part2Fn(in)
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
