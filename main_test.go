package main

import "fmt"

func ExampleScanf() {
	var n int
	var ext string
	fmt.Sscanf("702.mp3", "%d%s", &n, &ext)
	fmt.Println(n, ext)
	// Output:
	// 702 .mp3
}
