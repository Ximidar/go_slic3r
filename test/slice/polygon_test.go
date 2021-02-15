package slice_test

import (
	"goSlicer/slice"
	"testing"
)

func TestArea(t *testing.T) {
	poly := slice.NewPolygon()
	poly.Push(slice.NewPoint(100, 100))
	poly.Push(slice.NewPoint(200, 100))
	poly.Push(slice.NewPoint(200, 200))
	poly.Push(slice.NewPoint(100, 200))

	supposedArea := 100.00 * 100.00

	if poly.Area() != supposedArea {
		t.Fail()
	}
}
