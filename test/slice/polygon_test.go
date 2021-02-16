package slice_test

import (
	"fmt"
	"goSlicer/slice"
	"testing"
)

func TestArea(t *testing.T) {
	poly := slice.NewPolygon()

	poly.MP.Points.Push(
		slice.NewPoint(100, 100),
		slice.NewPoint(200, 100),
		slice.NewPoint(200, 200),
		slice.NewPoint(100, 200))

	supposedArea := 100.00 * 100.00

	if poly.Area() != supposedArea {
		fmt.Printf("Area %f != %f", poly.Area(), supposedArea)
		t.Fail()
	}
}
