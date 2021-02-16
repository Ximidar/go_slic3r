package slice

// PolygonEx is a representation of an external polygon
type PolygonEx struct {
	Contour *Polygon
	Holes   []*Polygon
}

// NewPolygonEx will create an External Polygon
func NewPolygonEx() *PolygonEx {
	pgx := new(PolygonEx)
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
func (pgx *PolygonEx) Polygons() []*Polygon {
	polys := make([]*Polygon, 0)
	polys = append(polys, pgx.Contour)
	polys = append(polys, pgx.Holes...)
	return polys

}
