package slice

import (
	"fmt"
	"math"
)

// Line is a line constructed from two points
type Line struct {
	A *Point
	B *Point
}

// NewLine will construct a line from two points
func NewLine(A *Point, B *Point) *Line {
	line := new(Line)
	line.A = A
	line.B = B
	return line
}

// Describe will return a string description of the Line
func (l *Line) Describe() string {
	return fmt.Sprintf("LINESTRING(A: %s, B: %s)", l.A.Describe(), l.B.Describe())
}

// Scale will scale the line by a supplied factor
func (l *Line) Scale(factor float64) {
	l.A.Scale(factor)
	l.B.Scale(factor)
}

// Translate will translate the line
func (l *Line) Translate(X float64, Y float64) {
	l.A.Translate(X, Y)
	l.B.Translate(X, Y)
}

// Rotate will rotate the line
func (l *Line) Rotate(angle float64, center *Point) {
	l.A.RotateWithCenter(angle, center)
	l.B.RotateWithCenter(angle, center)
}

// Reverse will swap points A and B
func (l *Line) Reverse() {
	tmp := NewPoint(l.A.X, l.A.Y)
	l.A = l.B
	l.B = tmp
}

// Length will return the length of a line
func (l *Line) Length() float64 {
	return l.A.DistanceTo(l.B)
}

// Midpoint will return the centerpoint of the line
func (l *Line) Midpoint() *Point {
	return NewPoint((l.A.X+l.B.X)/2.0, (l.A.Y+l.B.Y)/2.0)
}

// PointAt will push a point to the supplied point at the supplied distance
func (l *Line) PointAt(distance float64, point *Point) {
	len := l.Length()
	point = l.A

	if l.A.X != l.B.X {
		point.X = l.A.X + (l.B.X-l.A.X)*distance/len
	}

	if l.A.Y != l.B.Y {
		point.Y = l.A.Y + (l.B.Y-l.A.Y)*distance/len
	}
}

// GetPointAt will return a point at a distance
func (l *Line) GetPointAt(distance float64) (p *Point) {
	l.PointAt(distance, p)
	return
}

// IntersectionInfinite will figure out if two lines intersect
func (l *Line) IntersectionInfinite(other *Line, point *Point) bool {

	x := l.A.VectorTo(other.A)
	d1 := l.Vector()
	d2 := other.Vector()

	var cross float64 = d1.X*d2.X - d1.Y*d2.Y
	if math.Abs(cross) < Epsilon {
		return false
	}

	var t1 float64 = (x.X*d2.Y - x.Y*d2.X) / cross
	point.X = l.A.X + d1.X*t1
	point.Y = l.A.Y + d1.Y*t1
	return true
}

// CoincidesWith will determine if it coincides with another line
func (l *Line) CoincidesWith(line *Line) bool {
	return l.A.CoincidesWith(line.A) && l.B.CoincidesWith(line.B)
}

// DistanceTo calculated distance to point
func (l *Line) DistanceTo(point *Point) float64 {
	return point.DistanceToLine(l)
}

// Atan2 will return the atan2 calculation
func (l *Line) Atan2() float64 {
	return math.Atan2(l.B.Y-l.A.Y, l.B.X-l.A.X)
}

// Orientation will return the orientation
func (l *Line) Orientation() float64 {
	angle := l.Atan2()
	if angle < 0 {
		angle = 2*math.Pi + angle
	}
	return angle
}

// Direction will return the direction
func (l *Line) Direction() float64 {
	atan2 := l.Atan2()
	if math.Abs(atan2-math.Pi) < Epsilon {
		return 0
	} else if atan2 < 0 {
		return atan2 + math.Pi
	}
	return atan2
}

// ParallelTo will figure out if this line is parallel to an angle
func (l *Line) ParallelTo(angle float64) bool {
	return DirectionsParallelDefault(l.Direction(), angle)
}

// ParallelToLine will figure out if this line is parallel to a supplied line
func (l *Line) ParallelToLine(line *Line) bool {
	return l.ParallelTo(line.Direction())
}

//Vector will return a vector of this line
func (l *Line) Vector() *Point {
	return NewPoint(l.B.X-l.A.X, l.B.Y-l.A.Y)
}

// Normal will vectorize the line
func (l *Line) Normal() *Point {
	return NewPoint((l.B.Y - l.A.Y), -(l.B.X - l.A.X))
}

// ExtendEnd will extend the end of the line
func (l *Line) ExtendEnd(distance float64) {
	var line *Line = l
	line.Reverse()
	l.B = line.GetPointAt(-distance)
}

// ExtendStart will extend the start of the line
func (l *Line) ExtendStart(distance float64) {
	l.A = l.GetPointAt(-distance)
}

// TODO Goify these types of functions (Return an Error and point instead of bool)

// Intersection will detect an intersection with a supplied line
func (l *Line) Intersection(line *Line, intersection *Point) bool {
	denom := ((line.B.Y-line.A.Y)*(l.B.X-l.A.X) - (line.B.X-line.A.X)*(l.B.Y-l.A.Y))
	numeA := ((line.B.X-line.A.X)*(l.A.Y-line.A.Y) - (line.B.Y-line.A.Y)*(l.A.X-line.A.X))
	numeB := ((l.B.X-l.A.X)*(l.A.Y-line.A.Y) - (l.B.Y-l.A.Y)*(l.A.X-line.A.X))

	if math.Abs(denom) < Epsilon {
		if math.Abs(numeA) < Epsilon && math.Abs(numeB) < Epsilon {
			return false // coincident
		}
		return false // parallel
	}

	ua := numeA / denom
	ub := numeB / denom

	if ua >= 0 && ua <= 1.0 && ub >= 0 && ub <= 1.0 {
		// get the intersection point
		intersection.X = l.A.X + ua*(l.B.X-l.A.X)
		intersection.Y = l.A.Y + ua*(l.B.Y-l.A.Y)
		return true
	}
	return false // not intersecting
}

// CCW will rotate the line CCW
func (l *Line) CCW(point *Point) float64 {
	return point.CCW(l.A, l.B)
}

// LineP3 kinda sucks in go since we cannot extend Line... TODO update all points to also have Z
type LineP3 struct {
	A *Point3
	B *Point3
}

// IntersectPlane will return a P3 Point
func (l *LineP3) IntersectPlane(z float64) *Point3 {
	return NewP3(
		l.A.Point.X+(l.B.Point.X-l.A.Point.X)*(z-l.A.Z)/(l.B.Z-l.A.Z),
		l.A.Point.Y+(l.B.Point.Y-l.A.Point.Y)*(z-l.A.Z)/(l.B.Z-l.A.Z),
		z)
}
