package slice_test

import (
	"fmt"
	"goSlicer/slice"
	"testing"
)

func TestReversePolyPtLinks(t *testing.T) {
	fmt.Println("Making OutPts")

	outpt := new(slice.OutPt)
	head := outpt
	last := new(slice.OutPt)

	for i := 0; i < 5; i++ {
		outpt.Idx = i
		outpt.Pt = slice.NewPoint(1.0+float64(i), 2.0+float64(i))
		outpt.Next = new(slice.OutPt)
		outpt.Prev = last

		last = outpt
		outpt = outpt.Next
	}
	outpt = last
	last.Next = head
	head.Prev = last

	outpt = outpt.Next

	fmt.Println("original order")
	origIdx := outpt.Idx
	var count int = 0
	for {
		fmt.Println(outpt.Idx)
		outpt = outpt.Next
		if outpt.Idx == origIdx {
			break
		}
		count += 1
	}

	if count < 3 {
		t.Fail()
	}

	fmt.Println("Reversing OutPts")
	slice.ReversePolyPtLinks(outpt)

	origIdx = outpt.Idx
	count = 0
	for {
		fmt.Println(outpt.Idx)
		outpt = outpt.Next
		if outpt.Idx == origIdx {
			break
		}
		count += 1
	}
	if count < 3 {
		t.Fail()
	}
}
