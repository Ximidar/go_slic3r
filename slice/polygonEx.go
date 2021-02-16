package slice

type PolygonEx struct {
	Contour *Polygon
	Holes   *Polygon
}

// NewPolygonEx will create an External Polygon
func NewPolygonEx() *PolygonEx {
	pgx := new(PolygonEx)
	return pgx
}

// Points will return all points
func (pgx *PolygonEx) Points() []*Point {

}
