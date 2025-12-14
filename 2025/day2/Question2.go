package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReaderSize(f, 1<<20)

	var sum int64

	for {
		start, ok := readInt(r)
		if !ok {
			break
		}
		if err := skipUntilDashOrEOF(r); err == io.EOF {
			break
		}
		end, ok := readInt(r)
		if !ok {
			break
		}
		if start > end {
			start, end = end, start
		}

		for x := int64(start); x <= int64(end); x++ {
			if isInvalidByRepetitionHash(x) {
				sum += x
			}
		}

		if err := skipUntilDigitOrEOF(r); err == io.EOF {
			break
		}
	}

	fmt.Println(sum)
}

func isInvalidByRepetitionHash(n int64) bool {
	if n < 0 {
		return false
	}

	s := toDigits(n)
	L := len(s)
	if L < 2 {
		return false
	}

	prefix, pow := buildRollingHash(s)

	for p := 1; p <= L/2; p++ {
		if L%p != 0 {
			continue
		}
		k := L / p
		if k < 2 {
			continue
		}

		h0 := subHash(prefix, pow, 0, p)
		ok := true
		for i := p; i < L; i += p {
			if subHash(prefix, pow, i, i+p) != h0 {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}

	return false
}

func toDigits(n int64) []byte {
	var buf [32]byte
	i := len(buf)
	for {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
		if n == 0 {
			break
		}
	}
	out := make([]byte, len(buf)-i)
	copy(out, buf[i:])
	return out
}

const base uint64 = 911382323

func buildRollingHash(s []byte) (prefix []uint64, pow []uint64) {
	L := len(s)
	prefix = make([]uint64, L+1)
	pow = make([]uint64, L+1)
	pow[0] = 1

	for i := 0; i < L; i++ {
		prefix[i+1] = prefix[i]*base + uint64(s[i])
		pow[i+1] = pow[i] * base
	}
	return prefix, pow
}

func subHash(prefix, pow []uint64, l, r int) uint64 {
	return prefix[r] - prefix[l]*pow[r-l]
}

func readInt(r *bufio.Reader) (int, bool) {
	val := 0
	found := false

	for {
		c, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return 0, false
			}
			panic(err)
		}
		if c >= '0' && c <= '9' {
			found = true
			val = int(c - '0')
			break
		}
	}

	for {
		c, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return val, true
			}
			panic(err)
		}
		if c < '0' || c > '9' {
			_ = r.UnreadByte()
			break
		}
		val = val*10 + int(c-'0')
	}

	return val, found
}

func skipUntilDashOrEOF(r *bufio.Reader) error {
	for {
		c, err := r.ReadByte()
		if err != nil {
			return err
		}
		if c == '-' {
			return nil
		}
	}
}

func skipUntilDigitOrEOF(r *bufio.Reader) error {
	for {
		c, err := r.ReadByte()
		if err != nil {
			return err
		}
		if c >= '0' && c <= '9' {
			_ = r.UnreadByte()
			return nil
		}
	}
}
