package slice

func NearZero(val float64) bool {
	return (val > -TOLERANCE && val < TOLERANCE)
}

type TEdge struct {
	Bot       *Point
	Curr      *Point
	Top       *Point
	Dx        float64
	PolyType  PolyType
	Side      EdgeSide
	WinDelta  int
	WinCnt    int
	WinCnt2   int
	OutIdx    int
	Next      *TEdge
	Prev      *TEdge
	NextInLML *TEdge
	NextInAEL *TEdge
	PrevInAEL *TEdge
	NextInSEL *TEdge
	PrevInSEL *TEdge
}

type IntersectNode struct {
	Edge1 *TEdge
	Edge2 *TEdge
	Pt    *Point
}

type LocalMinimum struct {
	Y          float64
	LeftBound  *TEdge
	RightBound *TEdge
}

func LocalMinSort(a, b LocalMinimum) bool {
	return b.Y < a.Y
}

type OutRec struct {
	Idx       int
	IsHole    bool
	IsOpen    bool
	FirstLect *OutRec
	PolyNd    *Polygon
	Pts       *OutPt
	BottomPt  *OutPt
}

type OutPt struct {
	Idx  int
	Pt   *Point
	Next *OutPt
	Prev *OutPt
}

func (op *OutPt) Area() float64 {
	StartIdx := op.Idx
	if op == nil {
		return 0
	}

	var area float64 = 0.00
	for op.Idx != StartIdx {
		area += (op.Prev.Pt.X + op.Pt.X) * (op.Prev.Pt.Y - op.Pt.Y)
		op = op.Next
	}
	return area * 0.5
}

type Join struct {
	OutPt1 *Point
	OutPt2 *Point
	OffPt  *Point
}

func PointInPolygon(pt *Point, poly *Polygon) int {
	var result int = 0
	cnt := len(poly.MP.Points)

	if cnt < 3 {
		return 0
	}

	for idx, point := range poly.MP.Points {
		nextPoint := poly.MP.Points.NextEntry(idx)

		if nextPoint.Y == point.Y {
			if (nextPoint.X == point.X) || (point.Y == pt.Y &&
				((nextPoint.X > pt.X) == (point.X < pt.X))) {
				return -1
			}
		}

		if (point.Y < pt.Y) != (nextPoint.Y < pt.Y) {
			if point.X >= pt.X {
				if nextPoint.X > pt.X {
					result = 1 - result
				} else {
					d := (point.X-pt.X)*(nextPoint.Y-pt.Y) - (nextPoint.X-pt.X)*(point.Y-pt.Y)
					if d == 0.00 {
						return -1
					}
					if d > 0.00 == (nextPoint.Y > point.Y) {
						result = 1 - result
					}
				}
			} else {
				if nextPoint.X > pt.X {
					d := (point.X-pt.X)*(nextPoint.Y-pt.Y) - (nextPoint.X-pt.X)*(point.Y-pt.Y)
					if d == 0.00 {
						return -1
					}
					if d > 0.00 == (nextPoint.Y > point.Y) {
						result = 1 - result
					}
				}
			}
		}
	}

	return result
}

func PointInOutPt(pt *Point, op *OutPt) int {
	poly := NewPolygon()
	start := op.Idx

	poly.Push(op.Pt)
	for {
		op = op.Next
		if op.Idx == start {
			break
		}
		poly.Push(op.Pt)
	}

	return PointInPolygon(pt, poly)
}

func Poly2ContainsPoly1(poly1, poly2 *OutPt) bool {
	opIdx := poly1.Idx
	for {
		res := PointInOutPt(poly1.Pt, poly2)
		if res >= 0 {
			return res > 0
		}
		poly1 = poly1.Next
		if poly1.Idx == opIdx {
			break
		}
	}
	return true
}

func SlopesEqual(edge1, edge2 *TEdge) bool {
	return (edge1.Top.Y-edge1.Bot.Y)*(edge2.Top.X-edge2.Bot.X) ==
		(edge1.Top.X-edge1.Bot.X)*(edge2.Top.Y-edge2.Bot.Y)
}

func SlopesEqual3Pt(pt1, pt2, pt3 *Point) bool {
	return (pt1.Y-pt2.Y)*(pt2.X-pt3.X) == (pt1.X-pt2.X)*(pt2.Y-pt3.Y)
}

func SlopesEqual4Pt(pt1, pt2, pt3, pt4 *Point) bool {
	return (pt1.Y-pt2.Y)*(pt3.X-pt4.X) == (pt1.X-pt2.X)*(pt3.Y-pt4.Y)
}

func IsHorizontal(edge TEdge) bool {
	return edge.Dx == HORIZONTAL
}

func GetDx2Pt(pt1, pt2 *Point) float64 {
	if pt1.Y == pt2.Y {
		return HORIZONTAL
	}
	return (pt2.X - pt1.X) / (pt2.Y - pt1.Y)
}

// continue https://github.com/slic3r/Slic3r/blob/master/xs/src/clipper.cpp
