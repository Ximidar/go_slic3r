package main

import (
	"fmt"
	slice "goSlicer/slice"
)

func main() {
	poly := slice.NewPolygon()
	poly.Push(slice.NewPoint(100, 100))
	poly.Push(slice.NewPoint(200, 100))
	poly.Push(slice.NewPoint(200, 200))
	poly.Push(slice.NewPoint(100, 200))

	supposedArea := 100.00 * 100.00

	fmt.Printf("Area: %f\n", poly.Area())
	if poly.Area() == supposedArea {
		fmt.Println("Area Matches")
	} else {
		fmt.Println("You Fool")
	}
}
