package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type P struct{ x, y int }

func main() {
	f, err := os.Open("input.txt")

	if err != nil {
		panic(err)
	}

	sc := bufio.NewScanner(f)
	var pts []P

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		p := strings.Split(line, ",")
		x, _ := strconv.Atoi(strings.TrimSpace(p[0]))
		y, _ := strconv.Atoi(strings.TrimSpace(p[1]))
		pts = append(pts, P{x, y})
	}

	bestArea := 0
	var a, b P

	for i := 0; i < len(pts); i++ {
		for j := i + 1; j < len(pts); j++ {
			dx := pts[i].x - pts[j].x
			if dx < 0 {
				dx = -dx
			}
			dy := pts[i].y - pts[j].y
			if dy < 0 {
				dy = -dy
			}
			area := (dx + 1) * (dy + 1)
			if area > bestArea {
				bestArea = area
				a, b = pts[i], pts[j]
			}
		}
	}

	fmt.Println("largest area:", bestArea)
	fmt.Printf("corners: (%d,%d) and (%d,%d)\n", a.x, a.y, b.x, b.y)
}
