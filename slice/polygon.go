package slice

import (
	"fmt"
	"math"
)

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

// Push will push a point into the polygon. Take this out at some point
func (pg *Polygon) Push(point *Point) {
	pg.mp.Push(point)
}

// Polyline will split the polygon at the first point
func (pg *Polygon) Polyline() *Polyline {
	return pg.SplitAtFirstPoint()
}

// GetPointAtIndex will retreive a point
func (pg *Polygon) GetPointAtIndex(index int) *Point {
	return pg.mp.PointAtIndex(index)
}

// GetLastPoint will retreive the last point
func (pg *Polygon) GetLastPoint() *Point {
	return pg.mp.FirstPoint() // last point == first point for polygons
}

// Lines will retrieve the polygon lines
func (pg *Polygon) Lines() []*Line {
	lines := make([]*Line, 0)
	for i := 0; i < len(pg.mp.Points); i += 2 {
		lines = append(lines, NewLine(pg.GetPointAtIndex(i), pg.GetPointAtIndex(i+1)))
	}
	lines = append(lines, NewLine(pg.mp.Points[len(pg.mp.Points)-1], pg.mp.Points[0]))
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

// EquallySpacedPoints will space out points
func (pg *Polygon) EquallySpacedPoints(distance float64) []*Point {
	return pg.SplitAtFirstPoint().EquallySpacedPoints(distance)
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

// Orientation will get the orientation of the polygon
func (pg *Polygon) Orientation() bool {
	return pg.Area() > 0
}

// IsCounterClockwise will determine if the polygon is CCW
// It seems like we could just simplify this. I had originally thought
// the orientation function from the clipperlib was going to be much more complicated.
func (pg *Polygon) IsCounterClockwise() bool {
	return pg.Orientation()
}

// IsClockwise will tell if a polygon is clockwise
func (pg *Polygon) IsClockwise() bool {
	return !pg.IsCounterClockwise()
}

// MakeCounterClockwise will reverse the polygon
func (pg *Polygon) MakeCounterClockwise() bool {
	if !pg.IsCounterClockwise() {
		pg.mp.Reverse()
		return true
	}
	return false
}

// MakeClockwise will make the polygon Clockwise
func (pg *Polygon) MakeClockwise() bool {
	if pg.IsClockwise() {
		pg.mp.Reverse()
		return true
	}
	return false
}

// IsValid will tell if a polygon is valid
func (pg *Polygon) IsValid() bool {
	return len(pg.mp.Points) >= 3
}

// ContainsPoint will check if the polygon contains a point
func (pg *Polygon) ContainsPoint(point *Point) (result bool) {
	result = false
	for index, p1 := range pg.mp.Points {
		p2 := pg.mp.PreviousPoint(index)

		// Theres a warning here:
		// FIXME this test is not numerically robust. Particularly, it does not handle horizontal segments at y == point.y well.
		// Does the ray with y == point.y intersect this line segment?
		if ((p1.Y > point.Y) != (p2.Y > point.Y)) && (point.X > (p2.X-p1.X)*(point.Y-p1.Y)/(p2.Y-p1.Y)+p1.X) {
			result = !result
		}
	}
	return
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

// Describe will return a string describing the polygon
func (pg *Polygon) Describe() string {
	describe := "POLYGON(("
	for _, point := range pg.mp.Points {
		describe += point.Describe()
		if !EqualPoints(point, pg.mp.LastPoint()) {
			describe += ","
		}
	}
	describe += "))"
	return describe
}

// ConcavePoints will find all concave point in the polygon
func (pg *Polygon) ConcavePoints(angle float64) []*Point {
	angle = 2.00*math.Pi - angle + Epsilon
	concavePoints := make([]*Point, 0)

	// check whether first point forms a concave angle
	if pg.mp.FirstPoint().CCWAngle(
		pg.mp.LastPoint(),
		pg.mp.NextPoint(0)) <= angle {
		concavePoints = append(concavePoints, pg.mp.FirstPoint())
	}

	// Check whether points [1:] form concave angles
	for index, point := range pg.mp.Points[1:] {
		if point.CCWAngle(pg.mp.PreviousPoint(index), pg.mp.NextPoint(index)) <= angle {
			concavePoints = append(concavePoints, point)
		}
	}

	// Check whether last point forms a concave angle
	if pg.mp.LastPoint().CCWAngle(pg.mp.PointAtIndex(-2), pg.mp.FirstPoint()) <= angle {
		concavePoints = append(concavePoints, pg.mp.LastPoint())
	}

	return concavePoints
}

// ConvexPoints will return all convex points
func (pg *Polygon) ConvexPoints(angle float64) []*Point {
	angle = 2.00*math.Pi - angle - Epsilon
	convexPoints := make([]*Point, 0)

	// check whether first point forms a convex angle
	if pg.mp.FirstPoint().CCWAngle(pg.mp.LastPoint(), pg.mp.NextPoint(0)) >= angle {
		convexPoints = append(convexPoints, pg.mp.FirstPoint())
	}

	// Check whether points [1:] form convex angles
	for index, point := range pg.mp.Points[1:] {
		if point.CCWAngle(pg.mp.PreviousPoint(index), pg.mp.NextPoint(index)) >= angle {
			convexPoints = append(convexPoints, point)
		}
	}

	// Check whether last point forms a convex angle
	if pg.mp.LastPoint().CCWAngle(pg.mp.PointAtIndex(-2), pg.mp.FirstPoint()) >= angle {
		convexPoints = append(convexPoints, pg.mp.LastPoint())
	}

	return convexPoints
}

// NewScale will scale the polygon
func (pg *Polygon) NewScale() {
	fmt.Println("Not implemented yet")
}
