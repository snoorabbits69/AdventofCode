package main

import (
	"crypto/md5"
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const chunk = 1 << 14

func appendInt(buf []byte, n int64) []byte {
	if n == 0 {
		return append(buf, '0')
	}
	var tmp [20]byte
	i := len(tmp)
	for n > 0 {
		i--
		tmp[i] = byte('0' + n%10)
		n /= 10
	}
	return append(buf, tmp[i:]...)
}

func incr(buf []byte, p int) []byte {
	for i := len(buf) - 1; i >= p; i-- {
		if buf[i] == '9' {
			buf[i] = '0'
		} else {
			buf[i]++
			return buf
		}
	}
	buf = append(buf, '0')
	buf[p] = '1'
	return buf
}

func solution(input string, mask byte, start int) int {
	workers := runtime.NumCPU()
	var next atomic.Int64
	var best atomic.Int64
	next.Store(int64(start))
	best.Store(math.MaxInt64)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := make([]byte, 0, len(input)+20)
			buf = append(buf, input...)
			p := len(buf)
			for {
				s := next.Add(chunk) - chunk
				if s >= best.Load() {
					return
				}
				buf = appendInt(buf[:p], s)
				for n := s; n < s+chunk; n++ {
					h := md5.Sum(buf)
					if h[0]|h[1]|(h[2]&mask) == 0 {
						for {
							b := best.Load()
							if n >= b || best.CompareAndSwap(b, n) {
								break
							}
						}
						break
					}
					buf = incr(buf, p)
				}
			}
		}()
	}
	wg.Wait()
	return int(best.Load())
}

func solve(input string) (int, int) {
	p1 := solution(input, 0xF0, 0)
	p2 := solution(input, 0xFF, p1)
	return p1, p2
}

func main() {
	const iters = 10
	input := "bgvyzdsv"
	var part1, part2 int

	start := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(input)
	}
	elapsed := time.Since(start)

	us := float64(elapsed.Nanoseconds()) / 1000.0
	fmt.Printf("Total: %.2f microseconds\n", us)
	fmt.Printf("Average: %.4f microseconds\n", us/float64(iters))
	fmt.Println("part1", part1)
	fmt.Println("part2", part2)
}
