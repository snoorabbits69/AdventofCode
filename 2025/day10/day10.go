package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"math"
	"math/bits"
	"time"
)

//go:embed input.txt
var data []byte

const (
	maxButtons  = 14
	maxJoltages = 11
)

type Column [maxJoltages]int64

type Machine struct {
	lights   uint32
	buttons  []uint32
	joltages []int64
}

type Basis struct {
	limit int64
	cost  int64
	vs    Column
}

type Subspace struct {
	rank    int
	nullity int
	lcm     int64
	rhs     Column
	basis   []Basis
}

func parseMachine(line []byte) Machine {
	i, n := 0, len(line)
	var m Machine

	for i < n && line[i] == ' ' {
		i++
	}

	start := i
	for i < n && line[i] != ' ' {
		i++
	}
	lightsTok := line[start:i]
	if len(lightsTok) > 1 {
		for idx, b := range lightsTok[1:] {
			if b == '#' {
				m.lights |= 1 << uint(idx)
			}
		}
	}

	for i < n {
		switch {
		case line[i] == ' ':
			i++

		case line[i] == '(':
			i++
			var mask uint32
			for i < n && line[i] != ')' {
				if d := line[i]; d >= '0' && d <= '9' {
					v := 0
					for i < n && line[i] >= '0' && line[i] <= '9' {
						v = v*10 + int(line[i]-'0')
						i++
					}
					mask |= 1 << uint(v)
				} else {
					i++
				}
			}
			if i < n {
				i++
			}
			m.buttons = append(m.buttons, mask)

		case line[i] == '{':
			i++
			var joltages []int64
			for i < n && line[i] != '}' {
				if d := line[i]; d == '-' || (d >= '0' && d <= '9') {
					neg := false
					if line[i] == '-' {
						neg = true
						i++
					}
					v := 0
					for i < n && line[i] >= '0' && line[i] <= '9' {
						v = v*10 + int(line[i]-'0')
						i++
					}
					if neg {
						v = -v
					}
					joltages = append(joltages, int64(v))
				} else {
					i++
				}
			}
			if i < n {
				i++
			}
			m.joltages = joltages

		default:
			i++
		}
	}

	return m
}

func parse(input []byte) []Machine {
	var machines []Machine
	for _, line := range bytes.Split(input, []byte("\n")) {
		line = bytes.TrimRight(line, "\r")
		if len(line) == 0 {
			continue
		}
		machines = append(machines, parseMachine(line))
	}
	return machines
}
func biterator(x uint32) []int {
	var out []int
	for x != 0 {
		out = append(out, bits.TrailingZeros32(x))
		x &= x - 1
	}
	return out
}

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcmOf(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	return a / gcd(a, b) * b
}

func configureLights(m Machine) uint32 {
	width := len(m.joltages)
	height := len(m.buttons)

	rank := 0
	var h [maxButtons]uint32
	var u [maxButtons]uint32
	for row := range u {
		u[row] = 1 << uint(row)
	}
	copy(h[:height], m.buttons)

	for col := 0; col < width; col++ {
		mask := uint32(1) << uint(col)
		found := -1
		for row := rank; row < height; row++ {
			if h[row]&mask != 0 {
				found = row
				break
			}
		}
		if found == -1 {
			continue
		}

		h[rank], h[found] = h[found], h[rank]
		u[rank], u[found] = u[found], u[rank]

		for row := 0; row < height; row++ {
			if row != rank && h[row]&mask != 0 {
				h[row] ^= h[rank]
				u[row] ^= u[rank]
			}
		}
		rank++
	}

	nullity := height - rank
	var particular uint32
	for row := 0; row < rank; row++ {
		mask := h[row] & (-h[row]) // isolate lowest one
		if m.lights&mask != 0 {
			particular ^= u[row]
		}
	}

	min := -1
	for i := 0; i < (1 << uint(nullity)); i++ {
		presses := particular
		for _, j := range biterator(uint32(i)) {
			presses ^= u[rank+j]
		}
		if cnt := bits.OnesCount32(presses); min == -1 || cnt < min {
			min = cnt
		}
	}
	return uint32(min)
}

func gaussianElimination(m Machine) Subspace {
	width := len(m.buttons)
	height := len(m.joltages)
	if width >= maxButtons || height >= maxJoltages {
		panic("machine exceeds fixed bounds")
	}

	var equations [maxJoltages][maxButtons]int64

	for row := 0; row < height; row++ {
		equations[row][width] = m.joltages[row]
	}

	for col := 0; col < width; col++ {
		limit := int64(math.MaxInt64)
		for _, row := range biterator(m.buttons[col]) {
			equations[row][col] = 1
			if m.joltages[row] < limit {
				limit = m.joltages[row]
			}
		}
		equations[height][col] = limit
	}

	rank, last := 0, width

	for rank < height && rank < last {
		found, bestAbs := -1, int64(-1)
		for row := rank; row < height; row++ {
			if equations[row][rank] != 0 {
				a := equations[row][rank]
				if a < 0 {
					a = -a
				}
				if found == -1 || a < bestAbs {
					found, bestAbs = row, a
				}
			}
		}

		if found != -1 {
			equations[rank], equations[found] = equations[found], equations[rank]
			pivot := equations[rank][rank]
			if pivot < 0 {
				pivot = -pivot
				for c := rank; c <= width; c++ {
					equations[rank][c] *= -1
				}
			}
			for row := 0; row < height; row++ {
				coeff := equations[row][rank]
				if row != rank && coeff != 0 {
					for col := range equations[row] {
						equations[row][col] = pivot*equations[row][col] - coeff*equations[rank][col]
					}
				}
			}
			rank++
		} else {
			last--
			for row := 0; row <= height; row++ {
				equations[row][rank], equations[row][last] = equations[row][last], equations[row][rank]
			}
		}
	}

	var lcmVal int64 = 1
	for pivot := 0; pivot < rank; pivot++ {
		lcmVal = lcmOf(lcmVal, equations[pivot][pivot])
	}
	for pivot := 0; pivot < rank; pivot++ {
		q := lcmVal / equations[pivot][pivot]
		for c := rank; c <= width; c++ {
			equations[pivot][c] *= q
		}
	}

	nullity := width - rank
	var rhs Column
	for row := 0; row < maxJoltages; row++ {
		rhs[row] = equations[row][width]
	}

	basis := make([]Basis, nullity)
	for col := 0; col < nullity; col++ {
		limit := equations[height][col+rank]
		var vs Column
		for row := 0; row < maxJoltages; row++ {
			vs[row] = equations[row][rank+col]
		}
		var sum int64
		for row := 0; row < rank; row++ {
			sum += vs[row]
		}
		basis[col] = Basis{limit: limit, cost: lcmVal - sum, vs: vs}
	}

	return Subspace{rank: rank, nullity: nullity, lcm: lcmVal, rhs: rhs, basis: basis}
}

func recurse(sub *Subspace, rhs Column, remaining uint32, presses int64) (int64, bool) {
	rank := sub.rank
	temp := rhs

	for _, i := range biterator(remaining) {
		free := sub.basis[i]
		for row := 0; row < rank; row++ {
			if v := free.vs[row]; v < 0 {
				temp[row] -= v * free.limit
			}
		}
	}

	minValue := int64(math.MaxInt64)
	minIndex := -1
	var globalLower, globalUpper int64

	for _, i := range biterator(remaining) {
		free := sub.basis[i]
		var lower int64
		upper := free.limit

		for row := 0; row < rank; row++ {
			v, r := free.vs[row], temp[row]
			if v > 0 {
				if q := r / v; q < upper {
					upper = q
				}
			}
			if v < 0 {
				rr := r + v*free.limit
				if q := (rr + v + 1) / v; q > lower {
					lower = q
				}
			}
		}

		if size := upper - lower + 1; size > 0 && size < minValue {
			minValue, minIndex, globalLower, globalUpper = size, i, lower, upper
		}
	}

	if minIndex == -1 {
		return 0, false
	}

	remaining2 := remaining ^ (1 << uint(minIndex))
	lower, upper := globalLower, globalUpper
	b := sub.basis[minIndex]
	cost, vs, lcmVal := b.cost, b.vs, sub.lcm

	if remaining2 != 0 {
		for row := 0; row < rank; row++ {
			rhs[row] -= (lower - 1) * vs[row]
		}
		best, found := int64(math.MaxInt64), false
		for n := lower; n <= upper; n++ {
			for row := 0; row < rank; row++ {
				rhs[row] -= vs[row]
			}
			if v, ok := recurse(sub, rhs, remaining2, presses+n*cost); ok && (!found || v < best) {
				best, found = v, true
			}
		}
		return best, found
	}

	check := func(n int64) (int64, bool) {
		total := (presses + n*cost) / lcmVal
		for row := 0; row < rank; row++ {
			if (rhs[row]-n*vs[row])%lcmVal != 0 {
				return 0, false
			}
		}
		return total, true
	}

	if cost >= 0 {
		for n := lower; n <= upper; n++ {
			if v, ok := check(n); ok {
				return v, true
			}
		}
	} else {
		for n := upper; n >= lower; n-- {
			if v, ok := check(n); ok {
				return v, true
			}
		}
	}
	return 0, false
}

func configureJoltages(m Machine) int64 {
	sub := gaussianElimination(m)

	var particular int64
	for row := 0; row < sub.rank; row++ {
		particular += sub.rhs[row]
	}

	if sub.nullity == 0 {
		return particular / sub.lcm
	}

	remaining := uint32(1)<<uint(len(sub.basis)) - 1
	result, ok := recurse(&sub, sub.rhs, remaining, particular)
	if !ok {
		panic("no integer solution found")
	}
	return result
}

func part1(machines []Machine) uint32 {
	var sum uint32
	for _, m := range machines {
		sum += configureLights(m)
	}
	return sum
}

func part2(machines []Machine) int64 {
	var sum int64
	for _, m := range machines {
		sum += configureJoltages(m)
	}
	return sum
}

func solve(machines []Machine) (uint32, int64) {
	return part1(machines), part2(machines)
}

func main() {
	start := time.Now()
	const iters = 1000
	var part1 uint32
	var part2 int64
	var machines []Machine
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		machines = parse(data)
		part1, part2 = solve(machines)

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
