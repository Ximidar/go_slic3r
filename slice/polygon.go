package slice

import "fmt"

// Polygon is a collection of points that make up a polygon
type Polygon struct {
	mp *MultiPoint
}

// NewPolygon will construct a polygon
func NewPolygon() *Polygon {
	p := new(Polygon)

	return p
}

// GetPointAtIndex will retreive a point
func (p *Polygon) GetPointAtIndex(index int) *Point {
	return p.mp.Points[index]
}

// GetLastPoint will retreive the last point
func (p *Polygon) GetLastPoint() *Point {
	return p.mp.FirstPoint() // last point == first point for polygons
}

// Lines will retrieve the polygon lines
func (p *Polygon) Lines() []*Line {
	lines := make([]*Line, 0)
	for i := 0; i < len(p.mp.Points); i += 2 {
		lines = append(lines, NewLine(p.GetPointAtIndex(i), p.GetPointAtIndex(i+1)))
	}
	lines = append(lines, NewLine(p.mp.Points[len(p.mp.Points)-1], p.mp.Points[0]))
	return lines
}

// SplitAtVertex splits the polygon at a vertex and returns a polyline
func (p *Polygon) SplitAtVertex() {
	fmt.Println("Not implemented")
}
