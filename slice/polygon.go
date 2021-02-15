package slice

import "math"

// Polygon is a collection of points that make up a polygon
type Polygon struct {
	mp *MultiPoint
}

// NewPolygon will construct a polygon
func NewPolygon() *Polygon {
	p := new(Polygon)
	p.mp = NewMultiPointNoInterface()
	return p
}

func (pg *Polygon) Push(point *Point) {
	pg.mp.Push(point)
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
	for index, p := range pg.mp.Points {
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

// Area will calculate the area of the polygon
func (pg *Polygon) Area() float64 {
	if len(pg.mp.Points) < 3 {
		return 0
	}

	area := 0.00
	for i, p1 := range pg.mp.Points {
		p2 := pg.mp.PreviousPoint(i)
		area += (p2.X + p1.X) * (p2.Y - p1.Y)
	}

	return -area * 0.5

}

// IsValid will tell if a polygon is valid
func (pg *Polygon) IsValid() bool {
	return len(pg.mp.Points) >= 3
}

// ContainsPoint will check if the polygon contains a point
// TODO implement if this is actually used.
func (pg *Polygon) ContainsPoint(point *Point) bool {
	return false
}

// RemoveCollinearPoints will remove collinear points
// TODO Test this since you changed how it works from the original
func (pg *Polygon) RemoveCollinearPoints() {
	if len(pg.mp.Points) <= 2 {
		return
	}

	points := pg.mp.GetPoints()
	pg.mp.Clear()

	for i := 0; i < len(points); i++ {
		p1 := points[i]
		for j, p2 := range points[i+1:] {
			p3 := points[j+1]
			l := NewLine(p1, p3)

			if l.DistanceTo(p2) > ScaledEpsilon {
				pg.mp.Push(p1)
				i = j
			}
		}
	}
}

// RemoveVerticalCollinearPoints will remove vertical CLP
// TODO test this as well since you changed it heavily
func (pg *Polygon) RemoveVerticalCollinearPoints(tolerance float64) {
	points := pg.mp.GetPoints()
	erasureIndex := []int{}
	for i, p1 := range points {
		for _, p2 := range points[i+1:] {
			if p2.X == p1.X && math.Abs(p2.Y-p1.Y) <= tolerance {
				erasureIndex = append(erasureIndex, i)
			}
		}
	}

	for offset, val := range erasureIndex {
		pg.mp.EraseAt(val + offset)
	}
}

// TriangulateConvex will only work on convex polygons
func (pg *Polygon) TriangulateConvex(polygons []*Polygon) {
	for i, point := range pg.mp.Points[2:] {
		poly := NewPolygon()
		poly.mp.Push(pg.mp.FirstPoint())
		poly.mp.Push(pg.mp.Points[(i+2)-1])
		poly.mp.Push(point)

		if poly.Area() > 0 {
			polygons = append(polygons, poly)
		}
	}
}

// Centroid will calculate the center of mass
func (pg *Polygon) Centroid() *Point {
	area := pg.Area()
	tmpX := 0.00
	tmpY := 0.00

	pline := pg.SplitAtFirstPoint()
	for i, point := range pline.mp.Points {
		nextP := pg.mp.PointAtIndex(i + 1)
		tmpX += (point.X + nextP.X) * (point.X*nextP.Y - nextP.X*point.Y)
		tmpY += (point.Y + nextP.Y) * (point.X*nextP.Y - nextP.X*point.Y)
	}

	return NewPoint(tmpX/(6*area), tmpY/(6*area))
}
