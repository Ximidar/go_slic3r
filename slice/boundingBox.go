package slice

import (
	"fmt"
	"math"
)

// BoundingBox defines the bounds of a box
type BoundingBox struct {
	Min     *Point
	Max     *Point
	defined bool
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
	bb.defined = true
	return bb
}

// NewBoundingBoxLines will construct a bounding box from lines
func NewBoundingBoxLines(lines ...*Line) *BoundingBox {
	var points Points = make(Points, 0)
	for _, line := range lines {
		points.Push(line.A, line.B)
	}
	return NewBoundingBox(points...)
}

// Polygon will alter a polygon to be a bounding box
func (bb *BoundingBox) Polygon(poly *Polygon) {
	poly.MP.Points.Clear()
	poly.Push(NewPoint(bb.Min.X, bb.Min.Y))
	poly.Push(NewPoint(bb.Max.X, bb.Min.Y))
	poly.Push(NewPoint(bb.Max.X, bb.Max.Y))
	poly.Push(NewPoint(bb.Min.X, bb.Max.Y))
}

// GetPolygon will get the polygon of the bounding box
func (bb *BoundingBox) GetPolygon() *Polygon {
	poly := NewPolygon()
	bb.Polygon(poly)
	return poly
}

// MergeBox will merge a point with the bounding box
func (bb *BoundingBox) MergeBox(mergeBox *BoundingBox) {
	if !bb.defined {
		bb.Min, bb.Max = mergeBox.Min, mergeBox.Max
		bb.defined = true
	}

	bb.Min.X = math.Min(mergeBox.Min.X, bb.Min.X)
	bb.Min.Y = math.Min(mergeBox.Min.Y, bb.Min.Y)

	bb.Max.X = math.Max(mergeBox.Max.X, bb.Max.X)
	bb.Max.Y = math.Max(mergeBox.Max.Y, bb.Max.Y)

}

// MergePoint will merge a point into the bounding box
func (bb *BoundingBox) MergePoint(point *Point) {
	if !bb.defined {
		bb.Min, bb.Max = point, point
		bb.defined = true
	}

	bb.Min.X = math.Min(point.X, bb.Min.X)
	bb.Min.Y = math.Min(point.Y, bb.Min.Y)

	bb.Max.X = math.Max(point.X, bb.Max.X)
	bb.Max.Y = math.Max(point.X, bb.Max.Y)
}

// Rotate will Rotate this bounding box
func (bb *BoundingBox) Rotate(angle float64) {
	bb = bb.Rotated(angle)
}

// RotateWithCenter will Rotate this bounding box
func (bb *BoundingBox) RotateWithCenter(angle float64, center *Point) {
	bb = bb.RotatedWithCenter(angle, center)
}

// Rotated will return a rotated boundingbox
func (bb *BoundingBox) Rotated(angle float64) (box *BoundingBox) {
	min := bb.Min.Rotated(angle)
	max := bb.Max.Rotated(angle)

	box.MergePoint(min)
	box.MergePoint(max)

	p1 := NewPoint(bb.Min.X, bb.Max.Y).Rotated(angle)
	p2 := NewPoint(bb.Max.X, bb.Min.Y).Rotated(angle)

	box.MergePoint(p1)
	box.MergePoint(p2)

	return
}

// RotatedWithCenter will return a Rotated around center boundingbox
func (bb *BoundingBox) RotatedWithCenter(angle float64, center *Point) (box *BoundingBox) {
	min := bb.Min.RotatedWithCenter(angle, center)
	max := bb.Max.RotatedWithCenter(angle, center)
	box.MergePoint(min)
	box.MergePoint(max)

	p1 := NewPoint(bb.Min.X, bb.Max.Y).RotatedWithCenter(angle, center)
	p2 := NewPoint(bb.Max.X, bb.Min.Y).RotatedWithCenter(angle, center)
	box.MergePoint(p1)
	box.MergePoint(p2)

	return
}

// Scale will scale the boundingbox
func (bb *BoundingBox) Scale(factor float64) {
	bb.Min.Scale(factor)
	bb.Max.Scale(factor)
}

// Size will give a point back for some reason
// Todo implement for Z
func (bb *BoundingBox) Size() *Point {
	return NewPoint(bb.Max.X-bb.Min.X, bb.Max.Y-bb.Min.Y)
}

// Radius will return the Bounding Box Radius
// TODO implement for Z
func (bb *BoundingBox) Radius() float64 {
	return 0
}

// Translate will translate the Bounding Box
func (bb *BoundingBox) Translate(x float64, y float64) {
	bb.Min.Translate(x, y)
	bb.Max.Translate(x, y)
}

// Offset will make an offset to the Bounding Box
func (bb *BoundingBox) Offset(delta float64) {
	bb.Min.Translate(-delta, -delta)
	bb.Max.Translate(delta, delta)
}

// Center will get the center
func (bb *BoundingBox) Center() *Point {
	return NewPoint(
		(bb.Max.X+bb.Min.X)/2,
		(bb.Max.Y+bb.Max.Y)/2)
}

// ContainsPoint will figure out if the bounding box contains a point
func (bb *BoundingBox) ContainsPoint(point *Point) bool {
	return point.X >= bb.Min.X && point.X <= bb.Max.X && point.Y >= bb.Min.Y && point.Y <= bb.Max.Y
}

// EqualBBoxi will Equate bounding boxes
func EqualBBoxi(bbox1 *BoundingBox, bbox2 *BoundingBox) bool {
	return EqualPoints(bbox1.Min, bbox2.Min) && EqualPoints(bbox1.Max, bbox2.Max)
}

// NotEqualBBoxi will Equate Bounding Boxes
func NotEqualBBoxi(bbox1 *BoundingBox, bbox2 *BoundingBox) bool {
	return !EqualBBoxi(bbox1, bbox2)
}

//Todo incorporate BoundingBox functions for Z
