package main

import (
	_ "embed"
	"fmt"
	"math"
	"sort"
	"time"
)

//go:embed input.txt
var data []byte

type Tile [2]int

type Interval struct {
	l, r int
}

func newInterval(l, r int) Interval {
	return Interval{l: l, r: r}
}

func (iv Interval) contains(x int) bool {
	return iv.l <= x && x <= iv.r
}

func (iv Interval) intersection(other Interval) Interval {
	l := iv.l
	if other.l > l {
		l = other.l
	}
	r := iv.r
	if other.r < r {
		r = other.r
	}
	return Interval{l: l, r: r}
}

type Candidate struct {
	x, y     int
	interval Interval
}

func solve(tiles []Tile) (uint64, uint64) {
	return part1(tiles), part2(tiles)
}

func part1(tiles []Tile) uint64 {
	topLeftTiles, topRightTiles := potentialCornerTiles(tiles)
	bottomLeftTiles, bottomRightTiles := potentialCornerTiles(reversedTiles(tiles))

	a := findLargestFromAllCorners(topLeftTiles, bottomRightTiles, true)
	b := findLargestFromAllCorners(bottomLeftTiles, topRightTiles, false)
	if b > a {
		return b
	}
	return a
}

func part2(tiles []Tile) uint64 {
	var largestArea uint64

	candidates := make([]Candidate, 0, 512)
	var descendingEdges []int
	var intervalsFromDescendingEdges []Interval

	for k := 0; k+1 < len(tiles); k += 2 {
		x0, y := tiles[k][0], tiles[k][1]
		x1 := tiles[k+1][0]

		for _, x := range [2]int{x0, x1} {
			toggleValueMembershipInOrderedList(&descendingEdges, x)
		}

		updateIntervalsFromDescendingEdges(descendingEdges, &intervalsFromDescendingEdges)

		for _, candidate := range candidates {
			for _, x := range [2]int{x0, x1} {
				if candidate.interval.contains(x) {
					dx := uint64(absDiff(candidate.x, x) + 1)
					dy := uint64(absDiff(candidate.y, y) + 1)
					if area := dx * dy; area > largestArea {
						largestArea = area
					}
				}
			}
		}

		retained := candidates[:0]
		for _, candidate := range candidates {
			for _, iv := range intervalsFromDescendingEdges {
				if iv.contains(candidate.x) {
					candidate.interval = iv.intersection(candidate.interval)
					retained = append(retained, candidate)
					break
				}
			}
		}
		candidates = retained

		for _, x := range [2]int{x0, x1} {
			for _, iv := range intervalsFromDescendingEdges {
				if iv.contains(x) {
					candidates = append(candidates, Candidate{
						x:        x,
						y:        y,
						interval: iv,
					})
					break
				}
			}
		}
	}

	return largestArea
}

func potentialCornerTiles(sortedTiles []Tile) ([]Tile, []Tile) {
	leftTiles := make([]Tile, 0)
	leftTilesLastX := math.MaxInt

	rightTiles := make([]Tile, 0)
	rightTilesLastX := math.MinInt

	n := len(sortedTiles)
	i := 0

	for i < n {
		firstInRow := sortedTiles[i]
		lastInRow := firstInRow

		j := i + 1
		for j < n && sortedTiles[j][1] == firstInRow[1] {
			lastInRow = sortedTiles[j]
			j++
		}

		y := firstInRow[1]

		leftX := firstInRow[0]
		if lastInRow[0] < leftX {
			leftX = lastInRow[0]
		}

		rightX := firstInRow[0]
		if lastInRow[0] > rightX {
			rightX = lastInRow[0]
		}

		if leftX < leftTilesLastX {
			leftTiles = append(leftTiles, Tile{leftX, y})
			leftTilesLastX = leftX
		}

		if rightX > rightTilesLastX {
			rightTiles = append(rightTiles, Tile{rightX, y})
			rightTilesLastX = rightX
		}

		i = j
	}

	for l, r := 0, len(rightTiles)-1; l < r; l, r = l+1, r-1 {
		rightTiles[l], rightTiles[r] = rightTiles[r], rightTiles[l]
	}

	return leftTiles, rightTiles
}

type work struct {
	pLo, pHi, qLo, qHi int
}

func addRange(w *[]work, pLo, pHi, qLo, qHi int) {
	if pLo <= pHi && qLo <= qHi {
		*w = append(*w, work{
			pLo: pLo,
			pHi: pHi,
			qLo: qLo,
			qHi: qHi,
		})
	}
}

func findLargestFromAllCorners(corner, oppositeCorner []Tile, topLeft bool) uint64 {
	if len(corner) == 0 || len(oppositeCorner) == 0 {
		return 0
	}

	var largest uint64

	workList := []work{{
		pLo: 0,
		pHi: len(corner) - 1,
		qLo: 0,
		qHi: len(oppositeCorner) - 1,
	}}

	for len(workList) > 0 {
		job := workList[len(workList)-1]
		workList = workList[:len(workList)-1]

		pMid := job.pLo + (job.pHi-job.pLo)/2
		p := corner[pMid]

		bestI := -1
		maxSize := uint64(0)
		qLim := job.qLo

		for qI := job.qLo; qI <= job.qHi; qI++ {
			q := oppositeCorner[qI]

			if p[0] > q[0] {
				qLim = qI
			} else if (p[1] < q[1]) == topLeft {
				dx := uint64(absDiff(p[0], q[0]) + 1)
				dy := uint64(absDiff(p[1], q[1]) + 1)

				if size := dx * dy; size > maxSize {
					maxSize = size
					bestI = qI
				}
			}
		}

		if bestI != -1 {
			if maxSize > largest {
				largest = maxSize
			}

			if pMid > 0 {
				addRange(
					&workList,
					job.pLo,
					pMid-1,
					job.qLo,
					bestI,
				)
			}

			addRange(
				&workList,
				pMid+1,
				job.pHi,
				bestI,
				job.qHi,
			)
		} else {
			if pMid > 0 && qLim > 0 {
				addRange(
					&workList,
					job.pLo,
					pMid-1,
					job.qLo,
					qLim-1,
				)
			}

			addRange(
				&workList,
				pMid+1,
				job.pHi,
				qLim,
				job.qHi,
			)
		}
	}

	return largest
}

func reversedTiles(tiles []Tile) []Tile {
	n := len(tiles)
	out := make([]Tile, n)

	for i, t := range tiles {
		out[n-1-i] = t
	}

	return out
}

func toggleValueMembershipInOrderedList(orderedList *[]int, value int) {
	l := *orderedList
	i := sort.SearchInts(l, value)

	if i < len(l) && l[i] == value {
		*orderedList = append(l[:i], l[i+1:]...)
		return
	}

	l = append(l, 0)
	copy(l[i+1:], l[i:])
	l[i] = value
	*orderedList = l
}

func updateIntervalsFromDescendingEdges(
	descendingEdges []int,
	toUpdate *[]Interval,
) {
	result := (*toUpdate)[:0]

	for i := 0; i+1 < len(descendingEdges); i += 2 {
		result = append(
			result,
			newInterval(
				descendingEdges[i],
				descendingEdges[i+1],
			),
		)
	}

	*toUpdate = result
}

func absDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

func parseNum(i *int) int {
	var v int

	for *i < len(data) && data[*i] != '\n' && data[*i] != ',' {
		v = v*10 + int(data[*i]-'0')
		*i++
	}

	return v
}

func main() {
	start := time.Now()

	const iters = 100

	i := 0

	var grids [][]int
	var grid []int

	for i < len(data) {
		n := parseNum(&i)
		grid = append(grid, n)

		if data[i] == '\n' {
			grids = append(grids, grid)
			grid = []int{}
		}

		i++
	}

	tiles := make([]Tile, len(grids))

	for idx, g := range grids {
		tiles[idx] = Tile{g[0], g[1]}
	}

	sort.Slice(tiles, func(a, b int) bool {
		if tiles[a][1] != tiles[b][1] {
			return tiles[a][1] < tiles[b][1]
		}

		return tiles[a][0] < tiles[b][0]
	})

	var part1, part2 uint64

	processStart := time.Now()

	for i := 0; i < iters; i++ {
		part1, part2 = solve(tiles)
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
