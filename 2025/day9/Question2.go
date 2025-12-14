package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type P struct{ x, y int }

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var raw []struct{ x, y int }
	xMap := make([]int, 0, 1024)
	yMap := make([]int, 0, 1024)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		sp := strings.Split(line, ",")
		x, _ := strconv.Atoi(strings.TrimSpace(sp[0]))
		y, _ := strconv.Atoi(strings.TrimSpace(sp[1]))
		raw = append(raw, struct{ x, y int }{x, y})

		xMap = append(xMap, x-1, x, x+1)
		yMap = append(yMap, y-1, y, y+1)
	}

	sort.Ints(xMap)
	sort.Ints(yMap)
	xMap = uniq(xMap)
	yMap = uniq(yMap)

	xIdx := make(map[int]int, len(xMap))
	yIdx := make(map[int]int, len(yMap))
	for i, v := range xMap {
		xIdx[v] = i
	}
	for i, v := range yMap {
		yIdx[v] = i
	}

	tiles := make([]P, len(raw))
	for i := range raw {
		tiles[i] = P{x: xIdx[raw[i].x], y: yIdx[raw[i].y]}
	}

	W := len(xMap)
	H := len(yMap)

	wall := make([][]bool, H)
	out := make([][]bool, H)
	good := make([][]bool, H)
	for y := 0; y < H; y++ {
		wall[y] = make([]bool, W)
		out[y] = make([]bool, W)
		good[y] = make([]bool, W)
	}

	for i := 0; i < len(tiles); i++ {
		a := tiles[i]
		b := tiles[(i+1)%len(tiles)]

		if a.x == b.x {
			x := a.x
			y1, y2 := a.y, b.y
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			for y := y1; y <= y2; y++ {
				wall[y][x] = true
			}
		} else {
			y := a.y
			x1, x2 := a.x, b.x
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			for x := x1; x <= x2; x++ {
				wall[y][x] = true
			}
		}
	}

	qx := []int{0}
	qy := []int{0}
	out[0][0] = true
	for h := 0; h < len(qx); h++ {
		x := qx[h]
		y := qy[h]

		if x+1 < W && !out[y][x+1] && !wall[y][x+1] {
			out[y][x+1] = true
			qx = append(qx, x+1)
			qy = append(qy, y)
		}
		if x-1 >= 0 && !out[y][x-1] && !wall[y][x-1] {
			out[y][x-1] = true
			qx = append(qx, x-1)
			qy = append(qy, y)
		}
		if y+1 < H && !out[y+1][x] && !wall[y+1][x] {
			out[y+1][x] = true
			qx = append(qx, x)
			qy = append(qy, y+1)
		}
		if y-1 >= 0 && !out[y-1][x] && !wall[y-1][x] {
			out[y-1][x] = true
			qx = append(qx, x)
			qy = append(qy, y-1)
		}
	}

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			if wall[y][x] || !out[y][x] {
				good[y][x] = true
			}
		}
	}

	bestArea := 0
	var bestA, bestB struct{ x, y int }

	for i := 0; i < len(tiles); i++ {
	next:
		for j := i + 1; j < len(tiles); j++ {
			a := tiles[i]
			b := tiles[j]

			ax, ay := xMap[a.x], yMap[a.y]
			bx, by := xMap[b.x], yMap[b.y]

			minX, maxX := ax, bx
			if minX > maxX {
				minX, maxX = maxX, minX
			}
			minY, maxY := ay, by
			if minY > maxY {
				minY, maxY = maxY, minY
			}

			area := (maxX - minX + 1) * (maxY - minY + 1)
			if area <= bestArea {
				continue
			}

			cx1, cx2 := a.x, b.x
			cy1, cy2 := a.y, b.y
			if cx1 > cx2 {
				cx1, cx2 = cx2, cx1
			}
			if cy1 > cy2 {
				cy1, cy2 = cy2, cy1
			}

			for y := cy1; y <= cy2; y++ {
				for x := cx1; x <= cx2; x++ {
					if !good[y][x] {
						continue next
					}
				}
			}

			bestArea = area
			bestA = struct{ x, y int }{ax, ay}
			bestB = struct{ x, y int }{bx, by}
		}
	}

	fmt.Println("largest area:", bestArea)
	fmt.Printf("corners: (%d,%d) and (%d,%d)\n", bestA.x, bestA.y, bestB.x, bestB.y)
}

func uniq(a []int) []int {
	if len(a) == 0 {
		return a
	}
	out := a[:1]
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			out = append(out, a[i])
		}
	}
	return out
}
