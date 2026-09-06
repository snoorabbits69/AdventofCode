package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var data []byte

const totalTime = 2503

func parseNumber(i *int) int16 {
	var v int16
	for *i < len(data) {
		c := data[*i]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int16(c-'0')
		*i++
	}
	return v
}

func distanceAt(speed, travelTime, restTime, t int16) int16 {
	cycle := travelTime + restTime
	fullCycles := t / cycle
	remainder := t % cycle
	if remainder > travelTime {
		remainder = travelTime
	}
	return speed * (travelTime*fullCycles + remainder)
}

func solve() (int16, int16) {
	n := len(data)
	i := 0

	var speeds, travelTimes, restTimes []int16
	var maxDistance, maxPoints int16

	for i < n {
		for i < n && data[i] != ' ' {
			i++
		}
		if i >= n {
			break
		}
		i += 9
		speed := parseNumber(&i)
		speeds = append(speeds, speed)
		i += 10
		travelTime := parseNumber(&i)
		travelTimes = append(travelTimes, travelTime)
		i += 33
		restTime := parseNumber(&i)
		restTimes = append(restTimes, restTime)

		d := distanceAt(speed, travelTime, restTime, totalTime)
		if d > maxDistance {
			maxDistance = d
		}

		i += 10
	}

	count := len(speeds)
	points := make([]int16, count)
	dists := make([]int16, count)

	for t := int16(1); t <= totalTime; t++ {
		var maxDist int16 = -1
		for idx := 0; idx < count; idx++ {
			dists[idx] = distanceAt(speeds[idx], travelTimes[idx], restTimes[idx], t)
			if dists[idx] > maxDist {
				maxDist = dists[idx]
			}
		}
		for idx := 0; idx < count; idx++ {
			if dists[idx] == maxDist {
				points[idx]++
			}
		}
	}

	for idx := 0; idx < count; idx++ {
		if points[idx] > maxPoints {
			maxPoints = points[idx]
		}
	}

	return maxDistance, maxPoints
}

func main() {

	fmt.Println(solve())
}
