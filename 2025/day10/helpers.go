package main

import (
	"math"
	"strconv"
	"strings"
)

func containsInt(s []int, target int) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

func parseIntsfromContent(content string) []int {
	replacer := strings.NewReplacer(
		"(", " ", ")", " ",
		"{", " ", "}", " ",
		",", " ",
	)
	cleaned := replacer.Replace(content)

	fields := strings.Fields(cleaned)

	nums := make([]int, 0, len(fields))

	for _, f := range fields {
		if n, err := strconv.Atoi(f); err == nil {
			nums = append(nums, n)
		}
	}

	return nums
}

func parseIndicator(content string) int {
	var mask int
	for i := 0; i < len(content); i++ {
		if content[i] == '#' {
			mask |= 1 << i
		}
	}
	return mask
}

func buttonsToMask(idxs []int) uint64 {
	var m uint64
	for _, x := range idxs {
		m |= 1 << uint(x)
	}
	return m
}

func minBFS(numLights int, target int, buttonMasks []int) int {
	size := 1 << uint(numLights)

	dist := make([]int, size)
	for i := range dist {
		dist[i] = -1
	}

	queue := make([]int, 0, size)
	dist[0] = 0
	queue = append(queue, 0)

	head := 0
	for head < len(queue) {
		state := queue[head]
		head++

		if state == target {
			return dist[state]
		}

		for _, bm := range buttonMasks {
			next := state ^ bm

			if dist[next] == -1 {
				dist[next] = dist[state] + 1
				queue = append(queue, next)
			}
		}
	}

	return -1
}

func minLinearAlgebra(buttons [][]int, jolting []int) int {
	rows := len(jolting)
	cols := len(buttons)
	if rows == 0 {
		return 0
	}

	A := make([][]int, rows)
	for r := 0; r < rows; r++ {
		A[r] = make([]int, cols)
	}
	for c, button := range buttons {
		for _, counterIdx := range button {
			A[counterIdx][c] = 1
		}
	}

	M := make([][]float64, rows)
	for r := 0; r < rows; r++ {
		M[r] = make([]float64, cols+1)
		for c := 0; c < cols; c++ {
			M[r][c] = float64(A[r][c])
		}
		M[r][cols] = float64(jolting[r])
	}

	eps := 1e-9
	pivotRowForCol := make([]int, cols)
	for i := range pivotRowForCol {
		pivotRowForCol[i] = -1
	}

	r := 0
	for c := 0; c < cols && r < rows; c++ {
		p := -1
		for i := r; i < rows; i++ {
			if math.Abs(M[i][c]) > eps {
				p = i
				break
			}
		}
		if p == -1 {
			continue
		}

		M[r], M[p] = M[p], M[r]

		pv := M[r][c]
		for j := c; j <= cols; j++ {
			M[r][j] /= pv
		}

		for i := 0; i < rows; i++ {
			if i == r {
				continue
			}
			f := M[i][c]
			if math.Abs(f) <= eps {
				continue
			}
			for j := c; j <= cols; j++ {
				M[i][j] -= f * M[r][j]
			}
		}

		pivotRowForCol[c] = r
		r++
	}

	for i := 0; i < rows; i++ {
		allZero := true
		for j := 0; j < cols; j++ {
			if math.Abs(M[i][j]) > eps {
				allZero = false
				break
			}
		}
		if allZero && math.Abs(M[i][cols]) > eps {
			return -1
		}
	}

	isPivot := make([]bool, cols)
	for c := 0; c < cols; c++ {
		if pivotRowForCol[c] != -1 {
			isPivot[c] = true
		}
	}
	freeCols := make([]int, 0)
	for c := 0; c < cols; c++ {
		if !isPivot[c] {
			freeCols = append(freeCols, c)
		}
	}

	maxT := 0
	for _, v := range jolting {
		if v > maxT {
			maxT = v
		}
	}

	best := math.MaxInt

	eval := func(x []int) bool {
		for pc := 0; pc < cols; pc++ {
			pr := pivotRowForCol[pc]
			if pr == -1 {
				continue
			}
			val := M[pr][cols]
			for j := 0; j < cols; j++ {
				if j == pc {
					continue
				}
				val -= M[pr][j] * float64(x[j])
			}
			if val < -eps {
				return false
			}
			iv := math.Round(val)
			if math.Abs(val-iv) > 1e-7 {
				return false
			}
			x[pc] = int(iv)
			if x[pc] < 0 {
				return false
			}
		}
		return true
	}

	var dfs func(k int, x []int)
	dfs = func(k int, x []int) {
		if k == len(freeCols) {
			x2 := append([]int(nil), x...)
			if !eval(x2) {
				return
			}
			sum := 0
			for _, v := range x2 {
				sum += v
				if sum >= best {
					return
				}
			}
			if sum < best {
				best = sum
			}
			return
		}

		col := freeCols[k]
		for v := 0; v <= maxT; v++ {
			x[col] = v
			dfs(k+1, x)
		}
	}

	dfs(0, make([]int, cols))

	if best == math.MaxInt {
		return -1
	}
	return best
}
