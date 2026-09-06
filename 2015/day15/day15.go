package main

import "fmt"

//go:embed input.txt
var data []byte

func parse() {
	n := len(data)
	i := 0
	for i < n {

		for data[i] != ':' {
			fmt.Print(data[i])
		}

		i++
	}

}

func main() {

}
