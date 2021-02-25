package slice

// PolygonEx is a representation of an external polygon
type PolygonEx struct {
	Contour *Polygon
	Holes   Polygons
}

// NewPolygonEx will create an External Polygon
func NewPolygonEx() *PolygonEx {
	pgx := new(PolygonEx)
	pgx.Contour = NewPolygon()
	pgx.Holes = make(Polygons, 0)
	return pgx
}

// Points will return all points
func (pgx *PolygonEx) Points() Points {
	points := make(Points, 0)

	for _, poly := range pgx.Holes {
		for _, point := range poly.MP.Points {
			points = append(points, point)
		}
	}

	return points
}

// Polygons will return all polygons
func (pgx *PolygonEx) Polygons() Polygons {
	polys := make(Polygons, 0)
	polys.Push(pgx.Contour)
	polys.Push(pgx.Holes...)
	return polys
}

// Scale will scale the polygon
func (pgx *PolygonEx) Scale(factor float64) {
	pgx.Contour.MP.Scale(factor)
	for _, poly := range pgx.Holes {
		poly.MP.Scale(factor)
	}
}

// Translate will move the polygonEx
func (pgx *PolygonEx) Translate(x, y float64) {
	p := NewPoint(x, y)
	pgx.Contour.MP.Translate(p)
	for _, poly := range pgx.Holes {
		poly.MP.Translate(p)
	}
}

// Rotate will rotate
func (pgx *PolygonEx) Rotate(angle float64) {
	pgx.Contour.MP.Rotate(angle)
	for _, poly := range pgx.Holes {
		poly.MP.Rotate(angle)
	}
}

// RotateWithCenter will rotate around a center point
func (pgx *PolygonEx) RotateWithCenter(angle float64, center *Point) {
	pgx.Contour.MP.RotateWithCenter(angle, center)
	for _, poly := range pgx.Holes {
		poly.MP.RotateWithCenter(angle, center)
	}
}

// Area will find the area
func (pgx *PolygonEx) Area() float64 {
	var area float64 = pgx.Contour.Area()
	for _, poly := range pgx.Holes {
		area -= -poly.Area() // holes have negative area
	}
	return area
}

// IsValid will find if this PolygonEx is valid
func (pgx *PolygonEx) IsValid() bool {
	if !pgx.Contour.IsValid() || !pgx.Contour.IsCounterClockwise() {
		return false
	}

	for _, poly := range pgx.Holes {
		if !poly.IsValid() || poly.IsCounterClockwise() {
			return false
		}
	}
	return true
}

// Contains a line
func (pgx *PolygonEx) Contains(line *Line) bool{
	pl := NewPolyline()
	pl.MP.Points.Push(line.A, line.B)
	return pgx.ContainsPline(pl)
}

// ContainsPline detirmine if this contains a pline
func (pgx *PolygonEx) ContainsPline(pline *Polyline) bool{
	pline.
}