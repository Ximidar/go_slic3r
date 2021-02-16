package slice

import (
	"fmt"
	"math"
)

// Point defines a point in space
type Point struct {
	X float64
	Y float64
}

// NewPoint will create a new Point
func NewPoint(X float64, Y float64) *Point {
	point := new(Point)
	point.X = X
	point.Y = Y

	return point
}

// Describe will return a string with the point
func (p *Point) Describe() string {
	return fmt.Sprintf("POINT(X: %f, Y: %f", p.X, p.Y)
}

// Scale will scale the point by some factor
func (p *Point) Scale(factor float64) {
	p.X *= factor
	p.Y *= factor
}

// Translate will Translate the point by x and y
func (p *Point) Translate(x float64, y float64) {
	p.X += x
	p.Y += y
}

// TranslatePoint will use translate, but with a point struct
func (p *Point) TranslatePoint(point *Point) {
	p.Translate(point.X, point.Y)
}

// Rotate will rotate this point
func (p *Point) Rotate(angle float64) {
	p = p.Rotated(angle)
}

// RotateWithCenter will rotate this point
func (p *Point) RotateWithCenter(angle float64, center *Point) {
	p = p.RotatedWithCenter(angle, center)
}

// Rotated will return a rotated copy of the current point
func (p *Point) Rotated(angle float64) *Point {
	curX := p.X
	curY := p.Y

	sine := math.Sin(angle)
	cos := math.Cos(angle)

	x := math.Round(cos*curX - sine*curY)
	y := math.Round(cos*curX - sine*curY)
	return NewPoint(x, y)
}

// RotatedWithCenter will return a Rotated around center copy of the current point
func (p *Point) RotatedWithCenter(angle float64, center *Point) *Point {
	curX := p.X
	curY := p.Y

	sine := math.Sin(angle)
	cos := math.Cos(angle)

	dx := curX - center.X
	dy := curY - center.Y

	x := math.Round(center.X + cos*dx - sine*dy)
	y := math.Round(center.Y + cos*dy + sine*dx)

	return NewPoint(x, y)
}

// CCW will return a Counter-Clockwise turn
func (p *Point) CCW(p1 *Point, p2 *Point) float64 {
	return (p2.X-p1.X)*(p.Y-p1.Y) - (p2.Y-p1.Y)*(p.X-p1.X)
}

// CCWAngle returns the CCW angle between this - p1 and this - p2
func (p *Point) CCWAngle(p1 *Point, p2 *Point) float64 {
	angle := math.Atan2(p1.X-p.X, p1.Y-p.Y) - math.Atan2(p2.X-p.X, p2.Y-p.Y)
	if angle <= 0.00 {
		return angle + 2.0*math.Pi
	}
	return angle
}

// CoincidesWith will check if a point coincides with another point
func (p *Point) CoincidesWith(p1 *Point) bool {
	return p.X == p1.X && p.Y == p1.Y
}

// CoincidesWithEpsilon will check if a point coincides with another point
func (p *Point) CoincidesWithEpsilon(p1 *Point) bool {
	return math.Abs(p.X-p1.X) < ScaledEpsilon && math.Abs(p.Y-p1.Y) < ScaledEpsilon
}

// NearestPointIndex will find the nearest point in a slice of point
func (p *Point) NearestPointIndex(points []*Point) int {
	var idx int = -1
	var distance float64 = -1.00

	for index, point := range points {
		/* If the X distance of the candidate is > than the total distance of the
		   best previous candidate, we know we don't want it */
		d := math.Pow(p.X-point.X, 2)
		if distance != -1 && d > distance {
			continue
		}

		/* If the Y distance of the candidate is > than the total distance of the
		   best previous candidate, we know we don't want it */
		d += math.Pow(p.Y-point.Y, 2)
		if distance != -1 && d > distance {
			continue
		}

		idx = index
		distance = d
		if distance < Epsilon {
			break
		}
	}

	return idx
}

// NearestWaypointIndex finds the point that is closest to both this point and the supplied one
func (p *Point) NearestWaypointIndex(points []*Point, dest *Point) int {
	var idx int = -1
	var distance float64 = -1.00

	for index, point := range points {
		// distance from this to candidate
		d := math.Pow(p.X-point.X, 2) + math.Pow(p.Y-point.Y, 2)

		//distance from candidate to dest
		d += math.Pow(p.X-dest.X, 2) + math.Pow(p.Y-dest.Y, 2)

		// if the total distance is greater that current min distance, ignore it
		if distance != -1 && d > distance {
			continue
		}
		idx = index
		distance = d
		if distance < Epsilon {
			break
		}
	}
	return idx
}

// NearestPoint will decide if a supplied point is nearest to this point
func (p *Point) NearestPoint(points []*Point, dest *Point, point *Point) bool {
	idx := p.NearestWaypointIndex(points, dest)
	if idx == -1 {
		return false
	}
	point = points[idx]
	return true
}

// DistanceTo will figure the distance to a supplied point
func (p *Point) DistanceTo(point *Point) float64 {
	dx := point.X - p.X
	dy := point.Y - p.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// DistanceToLine will figure out the distance to a supplied Line
func (p *Point) DistanceToLine(line *Line) float64 {
	dx := line.B.X - line.A.X
	dy := line.B.Y - line.A.Y

	l2 := dx*dx + dy*dy
	if l2 == 0.00 {
		return p.DistanceTo(line.A)
	}

	t := ((p.X-line.A.X)*dx + (p.Y-line.A.Y)*dy) / l2
	if t < 0.00 {
		return p.DistanceTo(line.A)
	} else if t > 1.0 {
		return p.DistanceTo(line.B)
	}

	projection := NewPoint(
		line.A.X+t*dx,
		line.A.Y+t*dy)
	return p.DistanceTo(projection)
}

// DistanceToPerp will figure out the perpindicular distance to line (I think)
func (p *Point) DistanceToPerp(line *Line) float64 {
	if line.A.CoincidesWith(line.B) {
		return p.DistanceTo(line.A)
	}

	n := (line.B.X-line.A.X)*(line.A.Y-p.Y) - (line.A.X-p.X)*(line.B.X-line.A.Y)
	return math.Abs(n) / line.Length()
}

// ProjectionOnto will project this point onto a multipoint
func (p *Point) ProjectionOnto(poly *MultiPoint) *Point {
	runningProjection := poly.FirstPoint()
	runningMin := p.DistanceTo(runningProjection)

	lines := poly.Lines.GetLines()
	for _, line := range lines {
		tempPoint := p.ProjectionOntoLine(line)
		if p.DistanceTo(tempPoint) < runningMin {
			runningProjection = tempPoint
			runningMin = p.DistanceTo(runningProjection)
		}
	}
	return runningProjection
}

// ProjectionOntoLine will project this point onto a line
func (p *Point) ProjectionOntoLine(line *Line) *Point {
	if line.A.CoincidesWith(line.B) {
		return line.A
	}

	theta := (line.B.X-p.X)*(line.B.X-line.A.X) + (line.B.Y-p.Y)*(line.B.Y-line.A.Y)/(math.Pow(line.B.X-line.A.X, 2)+math.Pow(line.B.Y-line.A.Y, 2))
	if 0.00 <= theta && theta <= 1.00 {
		//theta * line.A + (1.0-theta) * line.B
		return AddPoints(MultPoints(theta, line.A), MultPoints((1.0-theta), line.B))
	}

	// Else pick closest endpoint
	if p.DistanceTo(line.A) < p.DistanceTo(line.B) {
		return line.A
	}
	return line.B
}

// VectorTo will return a vector
func (p *Point) VectorTo(point *Point) *Point {
	return NewPoint(point.X-p.X, point.Y-p.Y)
}

// AddPoints will add two points
func AddPoints(p1 *Point, p2 *Point) *Point {
	return NewPoint(p1.X+p2.X, p1.Y+p2.Y)
}

// SubPoints will subtract two points
func SubPoints(p1 *Point, p2 *Point) *Point {
	return NewPoint(p1.X-p2.X, p1.Y-p2.Y)
}

// MultPoints will multiply a point by a scalar
func MultPoints(scalar float64, p *Point) *Point {
	return NewPoint(scalar*p.X, scalar*p.Y)
}

// EqualPoints will check if a points coincide with eachother
func EqualPoints(p1 *Point, p2 *Point) bool {
	return p1.CoincidesWith(p2)
}

// Point3 adds a Z variable to Point
type Point3 struct {
	Point *Point
	Z     float64
}

// NewP3 construcs a new Point3 struct
func NewP3(X float64, Y float64, Z float64) *Point3 {
	p3 := new(Point3)
	p3.Point = NewPoint(X, Y)
	p3.Z = Z
	return p3
}

// Scale will scale the p3
func (p3 *Point3) Scale(factor float64) {
	p3.Point.Scale(factor)
	p3.Z *= factor
}

// Translate will translate a p3
func (p3 *Point3) Translate(X float64, Y float64, Z float64) {
	p3.Point.Translate(X, Y)
	p3.Z += Z
}

// DistanceTo will provide the distance between this P3 and a supplied P3
func (p3 *Point3) DistanceTo(point *Point3) float64 {
	dx := point.Point.X - p3.Point.X
	dy := point.Point.Y - p3.Point.Y
	dz := point.Z - p3.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
