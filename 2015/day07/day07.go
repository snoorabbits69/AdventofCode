package main

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed input.txt
var data []byte

type Op uint8

const (
	OP_ASSIGN Op = iota
	OP_AND
	OP_OR
	OP_LSHIFT
	OP_RSHIFT
	OP_NOT
)

type Expr struct {
	op        Op
	a, b, out uint16
}

const MaxWires = 26 + 26*26

const IMM uint16 = 1 << 15
const UNKNOWN uint16 = 0xFFFF

func isImm(x uint16) bool    { return x&IMM != 0 }
func immVal(x uint16) uint16 { return x &^ IMM }

func wireId(tok []byte) uint16 {
	if len(tok) == 1 {
		return uint16(tok[0] - 'a')
	}
	return uint16(tok[0]-'a')*26 + uint16(tok[1]-'a') + 26
}

func parseNum(tok []byte) uint16 {
	var v uint16
	for _, c := range tok {
		v = v*10 + uint16(c-'0')
	}
	return v
}

func atom(tok []byte) uint16 {
	if tok[0] >= '0' && tok[0] <= '9' {
		return IMM | parseNum(tok)
	}
	return wireId(tok)
}

func trimTok(b []byte) []byte {
	i, j := 0, len(b)
	for i < j && (b[i] == ' ' || b[i] == '\t' || b[i] == '\r' || b[i] == '\n') {
		i++
	}
	for j > i && (b[j-1] == ' ' || b[j-1] == '\t' || b[j-1] == '\r' || b[j-1] == '\n') {
		j--
	}
	return b[i:j]
}

func main() {
	start := time.Now()

	exprs, defined := parse(data)

	var part1, part2 uint16
	for i := 0; i < 100; i++ {
		part1, part2 = solve(exprs, defined)
	}

	const iters = 10000
	benchStart := time.Now()
	for i := 0; i < iters; i++ {
		part1, part2 = solve(exprs, defined)
	}
	benchTime := time.Since(benchStart)
	totalTime := time.Since(start)

	fmt.Println("solution of part1:", part1)
	fmt.Println("solution of part2:", part2)
	fmt.Println()

	fmt.Printf("Total bench time:  %d microseconds\n", benchTime.Microseconds())
	avgUs := float64(benchTime.Nanoseconds()) / 1_000.0 / float64(iters)
	fmt.Printf("Average per iter:  %.4f microseconds\n", avgUs)
	fmt.Printf("Total time:        %d microseconds\n", totalTime.Microseconds())
}

func solve(exprs [MaxWires]Expr, defined [MaxWires]bool) (uint16, uint16) {
	part1 := eval(exprs, defined, nil)

	bID := wireId([]byte("b"))
	exprs[bID] = Expr{op: OP_ASSIGN, a: IMM | part1, out: bID}
	defined[bID] = true

	part2 := eval(exprs, defined, nil)
	return part1, part2
}

func eval(exprs [MaxWires]Expr, defined [MaxWires]bool, memoIn *[MaxWires]uint16) uint16 {
	var memo [MaxWires]uint16
	var seen [MaxWires]bool
	for i := 0; i < MaxWires; i++ {
		memo[i] = UNKNOWN
	}

	var evalWire func(id uint16) uint16
	valueOf := func(x uint16) uint16 {
		if isImm(x) {
			return immVal(x)
		}
		return evalWire(x)
	}

	evalWire = func(id uint16) uint16 {
		if seen[id] {
			return memo[id]
		}
		seen[id] = true

		e := exprs[id]
		var v uint16
		switch e.op {
		case OP_ASSIGN:
			v = valueOf(e.a)
		case OP_NOT:
			v = ^valueOf(e.a)
		case OP_AND:
			v = valueOf(e.a) & valueOf(e.b)
		case OP_OR:
			v = valueOf(e.a) | valueOf(e.b)
		case OP_LSHIFT:
			v = valueOf(e.a) << valueOf(e.b)
		case OP_RSHIFT:
			v = valueOf(e.a) >> valueOf(e.b)
		}

		memo[id] = v
		return v
	}

	aID := wireId([]byte("a"))
	return evalWire(aID)
}

func parse(data []byte) ([MaxWires]Expr, [MaxWires]bool) {
	n := len(data)
	i := 0

	var exprs [MaxWires]Expr
	var defined [MaxWires]bool

	buf := make([]byte, 0, 8)

	flushAtom := func() (uint16, bool) {
		if len(buf) == 0 {
			return 0, false
		}
		tok := trimTok(buf)
		buf = buf[:0]
		if len(tok) == 0 {
			return 0, false
		}
		return atom(tok), true
	}

	readUntilDash := func() []byte {
		start := i
		for i+1 < n && data[i+1] != '-' {
			i++
		}
		return data[start : i+1]
	}

	var pendOp Op
	var pendA, pendB uint16
	var havePending bool

	for i < n {
		switch data[i] {
		case 'A':
			if a, ok := flushAtom(); ok {
				pendOp, pendA, havePending = OP_AND, a, true
			}
			i += 3
			pendB = atom(trimTok(readUntilDash()))
		case 'O':
			if a, ok := flushAtom(); ok {
				pendOp, pendA, havePending = OP_OR, a, true
			}
			i += 2
			pendB = atom(trimTok(readUntilDash()))
		case 'L':
			if a, ok := flushAtom(); ok {
				pendOp, pendA, havePending = OP_LSHIFT, a, true
			}
			i += 6
			pendB = atom(trimTok(readUntilDash()))
		case 'R':
			if a, ok := flushAtom(); ok {
				pendOp, pendA, havePending = OP_RSHIFT, a, true
			}
			i += 6
			pendB = atom(trimTok(readUntilDash()))
		case 'N':
			pendOp, havePending = OP_NOT, true
			buf = buf[:0]
			i += 4
			pendA = atom(trimTok(readUntilDash()))
		case '-':
			i += 2
			for data[i] == ' ' {
				i++
			}
			start := i
			for i < n && data[i] >= 'a' && data[i] <= 'z' {
				i++
			}
			outID := wireId(data[start:i])

			if havePending {
				exprs[outID] = Expr{op: pendOp, a: pendA, b: pendB, out: outID}
			} else {
				src, _ := flushAtom()
				exprs[outID] = Expr{op: OP_ASSIGN, a: src, out: outID}
			}
			defined[outID] = true
			havePending = false
		default:
			if data[i] != ' ' && data[i] != '\n' && data[i] != '\r' {
				buf = append(buf, data[i])
			}
			i++
		}
	}

	return exprs, defined
}
