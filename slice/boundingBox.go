package slice

import (
	"fmt"
	"math"
)

// BoundingBox defines the bounds of a box
type BoundingBox struct {
	Min *Point
	Max *Point
}

// NewBoundingBox will construct a bounding box
// TODO alter this constructor to also accept Z unit
func NewBoundingBox(points ...*Point) *BoundingBox {
	bb := new(BoundingBox)
	if len(points) == 0 {
		fmt.Println("Empty Point Set Supplied to BoundingBox")
		return bb
	}

	bb.Min = NewPoint(points[0].X, points[0].Y)
	bb.Max = NewPoint(points[0].X, points[0].Y)

	for _, point := range points {
		bb.Min.X = math.Min(bb.Min.X, point.X)
		bb.Min.Y = math.Min(bb.Min.Y, point.Y)

		bb.Max.X = math.Max(bb.Max.X, point.X)
		bb.Max.Y = math.Max(bb.Max.Y, point.Y)
	}
	return bb
}

// NewBoundingBoxLines will construct a bounding box from lines
func NewBoundingBoxLines(lines ...*Line) *BoundingBox {
	var points []*Point = make([]*Point, 0)
	for _, line := range lines {
		points = append(points, line.A)
		points = append(points, line.B)
	}
	return NewBoundingBox(points...)
}

// Polygon will create a bounding box polygon
func (bb *BoundingBox) Polygon() {
	print("Bounding Box Polygon is not implemented yet")
}

// Rotate will rotate the bounding box
func (bb *BoundingBox) Rotate(angle float64) {
	bb.Min.Rotate(angle)
	bb.Max.Rotate(angle)
}

// Todo incorporate the rest of Bounding Box
