//go:build ignore
// arrays1
// Make me compile!

package main

import "fmt"

func main() {
	var colors [3]string

	colors[0] = "red"
	colors[1] = "green"
	colors[2] = "blue"

	fmt.Printf("First color is %s
", colors[])
	fmt.Printf("Last color is %s
", colors[])
}
