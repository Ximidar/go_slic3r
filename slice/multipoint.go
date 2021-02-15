package slice

import (
	"errors"
	"math"
)

// Lines is an interface to gather lines
type Lines interface {
	GetLines() []*Line
}

// MultiPoint Holds multiple points
type MultiPoint struct {
	Points      []*Point
	BoundingBox float64
	Lines       Lines
}

// NewMultiPoint will construct a MultiPoint
func NewMultiPoint(lines Lines, points ...*Point) *MultiPoint {
	mp := new(MultiPoint)
	mp.Lines = lines
	mp.Points = points

	return mp
}

// NewMultiPointFromInterface will construct a Multipoint
func NewMultiPointFromInterface(iface Lines) *MultiPoint {
	mp := new(MultiPoint)
	mp.Lines = iface
	mp.Points = make([]*Point, 0)
	return mp
}

// NewMultiPointNoInterface will construct a Multipoint
func NewMultiPointNoInterface() *MultiPoint {
	mp := new(MultiPoint)
	mp.Points = make([]*Point, 0)
	return mp
}

// GetPoints will return a copy of all points
func (mp *MultiPoint) GetPoints() []*Point {
	copied := make([]*Point, len(mp.Points))
	copy(copied, mp.Points)
	return copied
}

// Empty will determine if the points are empty
func (mp *MultiPoint) Empty() bool {
	return len(mp.Points) == 0
}

// Clear will empty the multipoint
func (mp *MultiPoint) Clear() {
	mp.Points = make([]*Point, 0)
}

// Scale will scale all points
func (mp *MultiPoint) Scale(factor float64) {
	for _, point := range mp.Points {
		point.Scale(factor)
	}
}

// Translate will translate all points
func (mp *MultiPoint) Translate(vector *Point) {
	for _, point := range mp.Points {
		point.Translate(vector.X, vector.Y)
	}
}

// Rotate will rotate all points
func (mp *MultiPoint) Rotate(angle float64) {
	s := math.Sin(angle)
	c := math.Cos(angle)

	for _, point := range mp.Points {
		curX := point.X
		curY := point.Y
		point.X = math.Round(c*curX - s*curY)
		point.Y = math.Round(c*curY + s*curX)
	}
}

// Reverse will reverse the points
func (mp *MultiPoint) Reverse() {
	for i := len(mp.Points)/2 - 1; i >= 0; i-- {
		opp := len(mp.Points) - 1 - i
		mp.Points[i], mp.Points[opp] = mp.Points[opp], mp.Points[i]
	}
}

// FirstPoint will retreive the first point
func (mp *MultiPoint) FirstPoint() *Point {
	return mp.Points[0]
}

// LastPoint will retreive the last point
func (mp *MultiPoint) LastPoint() *Point {
	return mp.Points[len(mp.Points)-1]
}

// PointAtIndex will get a point at an index. If a negative index is supplied it will return a
// point from the back of the array
func (mp *MultiPoint) PointAtIndex(index int) *Point {
	if index < 0 {
		return mp.Points[len(mp.Points)+index]
	}
	return mp.Points[index]
}

// PreviousPoint will get the point previous to the supplied index
func (mp *MultiPoint) PreviousPoint(index int) *Point {
	idx := index - 1
	if idx < 0 {
		idx = len(mp.Points) - 1
	}
	return mp.PointAtIndex(idx)
}

// NextPoint will get the next point to the supplied index
func (mp *MultiPoint) NextPoint(index int) *Point {
	idx := index + 1
	if idx > len(mp.Points)-1 {
		idx = 0
	}
	return mp.PointAtIndex(idx)
}

// PopBack will pop the last point in the stack
func (mp *MultiPoint) PopBack() *Point {
	popped, newPoints := mp.Points[len(mp.Points)-1], mp.Points[:len(mp.Points)-1]
	mp.Points = newPoints
	return popped
}

// PopFront will pop the first point in the stack
func (mp *MultiPoint) PopFront() *Point {
	popped, newPoints := mp.Points[0], mp.Points[1:]
	mp.Points = newPoints
	return popped
}

// Push will append a point
func (mp *MultiPoint) Push(point ...*Point) {
	mp.Append(point...)
}

// PushFront will push a point to the front of the stack
func (mp *MultiPoint) PushFront(points ...*Point) {
	mp.Points = append(points, mp.Points...)
}

// EraseAt will delete an item at index
func (mp *MultiPoint) EraseAt(index int) {
	mp.Points = append(mp.Points[:index], mp.Points[index+1:]...)
}

// Window returns a sliding window version of the points
func (mp *MultiPoint) Window(size int) [][]*Point {
	points := mp.GetPoints()
	if len(points) <= size {
		return [][]*Point{points}
	}

	window := make([][]*Point, 0, len(points)-size+1)

	for i, j := 0, size; j <= len(points); i, j = i+1, j+1 {
		window = append(window, points[i:j])
	}
	return window
}

// Length will return the length of all lines
func (mp *MultiPoint) Length() float64 {
	lines := mp.Lines.GetLines()
	var len float64 = 0.00
	for _, line := range lines {
		len += line.Length()
	}
	return len
}

// FindPoint will attempt to find a point in points
func (mp *MultiPoint) FindPoint(p *Point) (int, error) {
	for index, point := range mp.Points {
		if point.CoincidesWith(p) {
			return index, nil
		}
	}
	return -1, errors.New("Point Not Found")
}

// HasBoundaryPoint will find boundary points
func (mp *MultiPoint) HasBoundaryPoint(point *Point) bool {
	// TODO fill this in when bounding box is made
	return false
}

// HasDuplicatePoints will return true if there are any duplicate points
func (mp *MultiPoint) HasDuplicatePoints() bool {
	for i := 1; i < len(mp.Points); i++ {
		if mp.Points[i-1].CoincidesWith(mp.Points[i]) {
			return true
		}
	}
	return false
}

// RemoveDuplicatePoints will remove duplicate points
// TODO test and or recode this
func (mp *MultiPoint) RemoveDuplicatePoints() bool {
	j := 0
	for i := 1; i < len(mp.Points); i++ {
		if !mp.Points[j].CoincidesWith(mp.Points[i]) {
			j++
			if j < i {
				mp.Points[j] = mp.Points[i]
			}
		}
	}

	if j+1 < len(mp.Points) {
		mp.Points = mp.Points[:j]
		return true
	}
	return false
}

// Append will append points to MultiPoint
func (mp *MultiPoint) Append(points ...*Point) {
	mp.Points = append(mp.Points, points...)
}

// Intersection will find intersections
func (mp *MultiPoint) Intersection(line *Line, intersection *Point) bool {
	lines := mp.Lines.GetLines()
	for _, line := range lines {
		if line.Intersection(line, intersection) {
			return true
		}
	}
	return false
}
