package main

import (
	_ "embed"
	"fmt"
	"math/bits"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

//go:embed input.txt
var data []byte

const (
	W  = 1000
	WB = 128
	WT = 16

	lo01 = 0x0101010101010101
	lo02 = 0x0202020202020202
	lo7F = 0x7F7F7F7F7F7F7F7F
	hi80 = 0x8080808080808080
	m16  = 0x00FF00FF00FF00FF
	mul4 = 0x0001000100010001
)

type instr struct {
	op       uint8
	y1, y2   int32
	bl, bh   int32
	bml, bmh uint64
	tl, th   int32
	tml, tmh uint64
}

var instrs [1024]instr

type seg struct{ y, h int32 }

type result struct {
	p1, p2 uint64
	_      [48]byte
}

func makeInstr(op uint8, x1, x2, y1, y2 int) instr {

	in := instr{op: op, y1: int32(y1), y2: int32(y2)}

	last := x2 - 1
	in.bl, in.bh = int32(x1>>3), int32(last>>3)
	in.bml = ^uint64(0) << (8 * uint(x1&7))
	in.bmh = ^uint64(0) >> (8 * uint(7-(last&7)))
	if in.bl == in.bh {
		in.bml &= in.bmh
		in.bmh = in.bml
	}

	in.tl, in.th = int32(x1>>6), int32(last>>6)
	in.tml = ^uint64(0) << uint(x1&63)
	in.tmh = ^uint64(0) >> uint(63-(last&63))
	if in.tl == in.th {
		in.tml &= in.tmh
		in.tmh = in.tml
	}
	return in
}

func applyBits(row *[WT]uint64, in *instr) {
	l, h := in.tl, in.th
	if l == h {
		m := in.tml
		switch in.op {
		case 0:
			row[l] |= m
		case 1:
			row[l] &^= m
		default:
			row[l] ^= m
		}
		return
	}
	switch in.op {
	case 0:
		row[l] |= in.tml
		for w := l + 1; w < h; w++ {
			row[w] = ^uint64(0)
		}
		row[h] |= in.tmh
	case 1:
		row[l] &^= in.tml
		for w := l + 1; w < h; w++ {
			row[w] = 0
		}
		row[h] &^= in.tmh
	default:
		row[l] ^= in.tml
		for w := l + 1; w < h; w++ {
			row[w] = ^row[w]
		}
		row[h] ^= in.tmh
	}
}

func applyBytes(row *[WB]uint64, in *instr) {
	l, h := in.bl, in.bh
	switch in.op {
	case 0:
		if l == h {
			row[l] += lo01 & in.bml
			return
		}
		row[l] += lo01 & in.bml
		for w := l + 1; w < h; w++ {
			row[w] += lo01
		}
		row[h] += lo01 & in.bmh
	case 2:
		if l == h {
			row[l] += lo02 & in.bml
			return
		}
		row[l] += lo02 & in.bml
		for w := l + 1; w < h; w++ {
			row[w] += lo02
		}
		row[h] += lo02 & in.bmh
	default:

		if l == h {
			v := row[l]
			nz := ((((v & lo7F) + lo7F) | v) & hi80) >> 7
			row[l] = v - (nz & in.bml)
			return
		}
		v := row[l]
		nz := ((((v & lo7F) + lo7F) | v) & hi80) >> 7
		row[l] = v - (nz & in.bml)
		for w := l + 1; w < h; w++ {
			v = row[w]
			nz = ((((v & lo7F) + lo7F) | v) & hi80) >> 7
			row[w] = v - nz
		}
		v = row[h]
		nz = ((((v & lo7F) + lo7F) | v) & hi80) >> 7
		row[h] = v - (nz & in.bmh)
	}
}

func worker(list []instr, segs []seg, next *atomic.Int32, out *result) {
	var r1 [WT]uint64
	var r2 [WB]uint64
	var p1, p2 uint64

	for {
		s := int(next.Add(1)) - 1
		if s >= len(segs) {
			break
		}
		y, h := segs[s].y, segs[s].h

		r1 = [WT]uint64{}
		r2 = [WB]uint64{}
		for k := range list {
			in := &list[k]
			if y < in.y1 || y >= in.y2 {
				continue
			}
			applyBits(&r1, in)
			applyBytes(&r2, in)
		}

		var c1, c2 uint64
		for _, w := range r1 {
			c1 += uint64(bits.OnesCount64(w))
		}
		for _, v := range r2 {
			a := (v & m16) + ((v >> 8) & m16)
			c2 += (a * mul4) >> 48
		}
		p1 += c1 * uint64(h)
		p2 += c2 * uint64(h)
	}

	out.p1 = p1
	out.p2 = p2
}

func solve(input []byte) (int, int) {
	n := 0
	i := 0
	for i < len(input) && input[i] == 't' {
		if i+6 >= len(input) {
			break
		}

		op := uint8(2)
		if input[i+6] == 'n' {
			op = 0
			i += 8
		} else if input[i+6] == 'f' {
			op = 1
			i += 9
		} else {
			i += 7
		}

		var v [4]int
		for k := 0; k < 4; k++ {
			for i < len(input) && (input[i] < '0' || input[i] > '9') {
				i++
			}
			for i < len(input) && input[i] >= '0' && input[i] <= '9' {
				v[k] = v[k]*10 + int(input[i]-'0')
				i++
			}
		}
		for i < len(input) && (input[i] == '\n' || input[i] == '\r') {
			i++
		}

		instrs[n] = makeInstr(op, v[0], v[2]+1, v[1], v[3]+1)
		n++
	}
	list := instrs[:n]

	var edge [W + 1]bool
	edge[0], edge[W] = true, true
	for k := range list {
		edge[list[k].y1] = true
		edge[list[k].y2] = true
	}
	segs := make([]seg, 0, 2*n+2)
	prev := 0
	for y := 1; y <= W; y++ {
		if edge[y] {
			segs = append(segs, seg{int32(prev), int32(y - prev)})
			prev = y
		}
	}

	threads := runtime.GOMAXPROCS(0)
	if threads > len(segs) {
		threads = len(segs)
	}
	results := make([]result, threads)
	var next atomic.Int32

	var wg sync.WaitGroup
	for t := 1; t < threads; t++ {
		wg.Add(1)
		go func(out *result) {
			defer wg.Done()
			worker(list, segs, &next, out)
		}(&results[t])
	}
	worker(list, segs, &next, &results[0])
	wg.Wait()

	var p1, p2 uint64
	for t := range results {
		p1 += results[t].p1
		p2 += results[t].p2
	}
	return int(p1), int(p2)
}

func main() {
	start := time.Now()
	const iters = 100
	var part1, part2 int
	processStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(data)
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
