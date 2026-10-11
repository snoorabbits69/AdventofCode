package main

import (
	"fmt"
	"math"
	"time"
)

const maxK = 50

var (
	primeTab []int
	ratio    []float64
	h50      float64
	tail12   float64
	tail60   float64
)

func init() {
	const lim = 1 << 14
	comp := make([]bool, lim+1)
	for i := 2; i <= lim; i++ {
		if !comp[i] {
			primeTab = append(primeTab, i)
			ratio = append(ratio, float64(i)/float64(i-1))
			for j := i * i; j <= lim; j += i {
				comp[j] = true
			}
		}
	}
	var t3, t4, t5 float64
	for k := 1; k <= maxK; k++ {
		x := 1 / float64(k)
		h50 += x
		if k%3 != 0 {
			t3 += x
		}
		if k%4 != 0 {
			t4 += x
		}
		if k%5 != 0 {
			t5 += x
		}
	}
	tail12 = max(t3, t4)
	tail60 = max(t3, t4, t5)
}

var (
	dNeed, dBest int
	dIncomplete  bool
)

func ext(i, r int) float64 {
	e, prod := 1.0, 1
	for ; i < len(primeTab); i++ {
		if prod*primeTab[i] > r {
			return e
		}
		prod *= primeTab[i]
		e *= ratio[i]
	}
	return math.Inf(1)
}

func dfs(m, s, i int) {
	if s >= dNeed && m < dBest {
		dBest = m
	}
	lim := dBest - 1
	r := lim / m
	for j := i; j < len(primeTab); j++ {
		p := primeTab[j]
		if p > r {
			return
		}
		if float64(s)*ext(j, r)*float64(lim) < float64(dNeed)*float64(m)*(1-1e-12) {
			return
		}
		mm, pe, sp := m*p, p, 1+p
		for {
			dfs(mm, s*sp, j+1)
			lim = dBest - 1
			if mm > lim/p {
				break
			}
			mm *= p
			pe *= p
			sp += pe
		}
		r = lim / m
	}
	if r > primeTab[len(primeTab)-1] {
		dIncomplete = true
	}
}

func part1(target int) int {
	need := (target-1)/10 + 1
	if need <= 1 {
		return 1
	}
	for ub := need/8 + 1; ; ub += ub/8 + 1 {
		dNeed, dBest, dIncomplete = need, ub+1, false
		dfs(1, 1, 0)
		if dIncomplete {
			return part1Sieve(int32(target))
		}
		if dBest <= ub {
			return dBest
		}
	}
}

const (
	blk = 512
	inf = math.MaxInt
)

var gbuf [blk]int32

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func scan(lo, hi, m, need int) int {
	q0 := (lo + m - 1) / m
	var jn, ev, de, stride [maxK]int
	sm, nk := 0, 0
	for k := 1; k <= maxK; k++ {
		if m%k == 0 {
			sm += m / k
			continue
		}
		l := m / gcd(m, k) * k
		n1 := (m*q0 + l - 1) / l * l
		jn[nk], ev[nk], de[nk], stride[nk] = (n1-m*q0)/m, n1/k, l/k, l/m
		nk++
	}
	for off := 0; ; off += blk {
		base := m * (q0 + off)
		if base >= hi {
			return 0
		}
		for j := range gbuf {
			gbuf[j] = int32((q0 + off + j) * sm)
		}
		for t := 0; t < nk; t++ {
			j, e, d, s := jn[t]-off, ev[t], de[t], stride[t]
			for ; j < blk; j += s {
				gbuf[j] += int32(e)
				e += d
			}
			jn[t], ev[t] = j+off, e
		}
		for j, v := range gbuf {
			if int(v) >= need {
				if n := base + j*m; n < hi {
					return n
				}
				return 0
			}
		}
	}
}

func part2(target int) int {
	need := (target-1)/11 + 1
	lo := max(1, int(float64(need)/h50*(1-1e-9)))
	hi60 := int(float64(need) / tail60 * (1 - 1e-9))
	hi12 := int(float64(need) / tail12 * (1 - 1e-9))
	if n := scan(lo, hi60, 60, need); n != 0 {
		return n
	}
	if n := scan(max(lo, hi60), hi12, 12, need); n != 0 {
		return n
	}
	return scan(max(lo, hi12), inf, 1, need)
}

var (
	lowBuf    []int32
	restBuf   []int32
	primesBuf []int32
)

func part1Sieve(target int32) int {
	need := (target-1)/10 + 1
	if need <= 1 {
		return 1
	}
	limit := need/3 + 1
	for {
		if int(limit)+1 > cap(lowBuf) {
			lowBuf = make([]int32, limit+1)
			restBuf = make([]int32, limit+1)
			primesBuf = make([]int32, 0, limit/8+16)
		}
		low, rest, primes := lowBuf[:limit+1], restBuf[:limit+1], primesBuf[:0]
		clear(low)
		for i := int32(2); i <= limit; i++ {
			var s int32
			if low[i] == 0 {
				primes = append(primes, i)
				low[i], rest[i], s = i+1, 1, i+1
			} else {
				s = rest[i] * low[i]
			}
			if s >= need {
				return int(i)
			}
			for _, p := range primes {
				ip := i * p
				if ip > limit {
					break
				}
				if i%p == 0 {
					rest[ip] = rest[i]
					low[ip] = low[i]*p + 1
					break
				}
				rest[ip] = s
				low[ip] = p + 1
			}
		}
		limit *= 2
	}
}

func main() {
	start := time.Now()
	const iters = 10000
	const target = 29000000
	var p1, p2 int

	processStart := time.Now()
	for i := 0; i < iters; i++ {
		p1 = part1(target)
		p2 = part2(target)
	}
	processTime := time.Since(processStart)
	totalTime := time.Since(start)

	elapsedUs := float64(processTime.Nanoseconds()) / 1000.0
	fmt.Printf("Total: %.2f microseconds\n", elapsedUs)
	fmt.Printf("Average: %.4f microseconds\n", elapsedUs/float64(iters))
	fmt.Printf("Total time:        %dns\n", totalTime.Nanoseconds())
	fmt.Println("part1", p1)
	fmt.Println("part2", p2)
}
