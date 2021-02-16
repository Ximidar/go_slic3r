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
	Points      Points
	BoundingBox *BoundingBox
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
	mp.Points = make(Points, 0)
	return mp
}

// NewMultiPointNoInterface will construct a Multipoint
func NewMultiPointNoInterface() *MultiPoint {
	mp := new(MultiPoint)
	mp.Points = make(Points, 0)
	return mp
}

// GetPoints will return a copy of all points
func (mp *MultiPoint) GetPoints() Points {
	return mp.Points.GetCopy()
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
