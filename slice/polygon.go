package slice

// Polygon is a collection of points that make up a polygon
type Polygon struct {
	mp *MultiPoint
}

// NewPolygon will construct a polygon
func NewPolygon() *Polygon {
	p := new(Polygon)

	return p
}

// Polyline will split the polygon at the first point
func (pg *Polygon) Polyline() *Polyline {
	return pg.SplitAtFirstPoint()
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
func (pg *Polygon) SplitAtVertex(point *Point) *Polyline {
	for index, p := range p.mp.Points {
		if p.CoincidesWith(point) {
			return pg.SplitAtIndex(index)
		}
	}
	return pg.Polyline()
}

// SplitAtIndex splits a closed polygon into and open polyline
func (pg *Polygon) SplitAtIndex(index int) *Polyline {
	pline := NewPolyline()
	for _, point := range pg.mp.Points[index:] {
		pline.mp.Push(point)
	}
	for _, point := range pg.mp.Points[:index+1] {
		pline.mp.Push(point)
	}
	return pline
}

// SplitAtFirstPoint will split the polygon at the first point
func (pg *Polygon) SplitAtFirstPoint() *Polyline {
	return pg.SplitAtIndex(0)
}

// TODO implement ClipperLib

// IsValid will tell if a polygon is valid
func (pg *Polygon) IsValid() bool {
	return len(pg.mp.Points) >= 3
}
