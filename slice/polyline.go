package slice

import "errors"

// Polyline is a line made up of multiple points
type Polyline struct {
	mp    *MultiPoint
	Width float64
}

// NewPolyline will make a line
func NewPolyline() *Polyline {
	p := new(Polyline)
	p.mp = NewMultiPointFromInterface(p)
	return p
}

// ToLine will Convert a Polyline to a Line
func (pl *Polyline) ToLine() (*Line, error) {
	if len(pl.mp.Points) > 2 {
		return nil, errors.New("Cannot Convert a Polyline to Line with more than two points")
	}
	return NewLine(pl.mp.FirstPoint(), pl.mp.LastPoint()), nil
}

// LeftmostPoint will grab the leftmost point
func (pl *Polyline) LeftmostPoint() *Point {
	leftP := pl.mp.FirstPoint()
	for _, point := range pl.mp.Points {
		if point.X < leftP.X {
			leftP = point
		}
	}
	return leftP
}

// Lines will return all lines
func (pl *Polyline) Lines() []*Line {
	lines := make([]*Line, 0)
	for i := 0; i < len(pl.mp.Points); i += 2 {
		lines = append(lines, NewLine(pl.mp.PointAtIndex(i), pl.mp.PointAtIndex(i+1)))
	}
	return lines
}

// GetLines is the same as Lines()
func (pl *Polyline) GetLines() []*Line {
	return pl.Lines()
}

// ClipEnd will remove a bit from the end of the polyline
func (pl *Polyline) ClipEnd(distance float64) {
	for distance > 0 {
		lastPoint := pl.mp.PopBack()
		if pl.mp.Empty() {
			break
		}

		LastSegmentLength := lastPoint.DistanceTo(pl.mp.LastPoint())
		if LastSegmentLength <= distance {
			distance -= LastSegmentLength
			continue
		}

		segment := NewLine(lastPoint, pl.mp.LastPoint())
		pl.mp.Push(segment.GetPointAt(distance))
		distance = 0
	}
}

// ClipFront will remove a bit from the start of the polyline
func (pl *Polyline) ClipFront(distance float64) {
	pl.mp.Reverse()
	pl.ClipEnd(distance)
	if len(pl.mp.Points) >= 2 {
		pl.mp.Reverse()
	}
}

// ExtendEnd will extend the end of a polyline
func (pl *Polyline) ExtendEnd(distance float64) {
	backPoint := pl.mp.LastPoint()
	backPoint2 := pl.mp.PointAtIndex(len(pl.mp.Points) - 2)
	backline := NewLine(backPoint, backPoint2)
	pl.mp.Points[len(pl.mp.Points)-1] = backline.GetPointAt(-distance)
}

// ExtendStart will extend the front of a polyline
func (pl *Polyline) ExtendStart(distance float64) {
	frontPoint := pl.mp.FirstPoint()
	frontPoint2 := pl.mp.PointAtIndex(1)
	frontLine := NewLine(frontPoint, frontPoint2)
	pl.mp.Points[0] = frontLine.GetPointAt(-distance)
}

// EquallySpacedPoints will return a collection of points picked
// on the polygon countour that are evenly spaced
func (pl *Polyline) EquallySpacedPoints(distance float64) []*Point {
	mp := NewMultiPointNoInterface()
	mp.Push(pl.mp.FirstPoint())
	var len float64 = 0

	for i := 1; !EqualPoints(pl.mp.PointAtIndex(i), pl.mp.LastPoint()); i++ {
		currentPoint := pl.mp.PointAtIndex(i)
		previousPoint := pl.mp.PointAtIndex(i - 1)
		segmentLength := currentPoint.DistanceTo(previousPoint)

		len += segmentLength
		if len < distance {
			continue
		}

		if len == distance {
			mp.Push(currentPoint)
			len = 0
			continue
		}

		var take float64 = segmentLength - (len - distance)
		segment := NewLine(previousPoint, currentPoint)
		mp.Push(segment.GetPointAt(take))
		i--
		len = -take
	}
	return mp.Points
}

// SplitAt will split the polyline at a point
func (pl *Polyline) SplitAt(point *Point, pline1 *Polyline, pline2 *Polyline) {
	if pl.mp.Empty() {
		return
	}

	var lineIdx int = 0
	p := pl.mp.FirstPoint()
	min := point.DistanceTo(p)
	lines := pl.Lines()

	for index, line := range lines {
		tempPoint := point.ProjectionOntoLine(line)
		if point.DistanceTo(tempPoint) < min {
			p = tempPoint
			min = point.DistanceTo(p)
			lineIdx = index
		}
	}

	// Create First Half
	pline1.mp.Clear()
	for _, line := range lines[:lineIdx+1] {
		if !line.A.CoincidesWith(p) {
			pline1.mp.Push(line.A)
		}
	}

	pline1.mp.Push(point)

	// Create Second Half
	pline2.mp.Clear()
	pline2.mp.Push(point)
	for _, line := range lines[lineIdx:] {
		pline2.mp.Push(line.B)
	}
}

// IsStraight will Check that each segment's direction is equal to the line connecting
// first point and last point. (Checking each line against the previous
// one would cause the error to accumulate.)
func (pl *Polyline) IsStraight() bool {
	dir := NewLine(pl.mp.FirstPoint(), pl.mp.LastPoint()).Direction()

	for _, line := range pl.Lines() {
		if !line.ParallelTo(dir) {
			return false
		}
	}
	return true
}

// Describe will return a string representation of the Polyline
func (pl *Polyline) Describe() string {
	description := "POLYLINE(("
	for _, point := range pl.mp.Points {
		description += point.Describe()
		if !point.CoincidesWith(pl.mp.LastPoint()) {
			description += ","
		}
	}
	description += "))"
	return description
}
