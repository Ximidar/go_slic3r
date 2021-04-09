package slice

import (
	"math"
)

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

func IsHorizontal(edge *TEdge) bool {
	return edge.Dx == HORIZONTAL
}

func GetDx2Pt(pt1, pt2 *Point) float64 {
	if pt1.Y == pt2.Y {
		return HORIZONTAL
	}
	return (pt2.X - pt1.X) / (pt2.Y - pt1.Y)
}

func SetDxFromTedge(edge *TEdge) {
	dy := edge.Top.Y - edge.Bot.Y
	if dy == 0 {
		edge.Dx = HORIZONTAL
		return
	}

	edge.Dx = (edge.Top.X - edge.Bot.Y) / dy
}

// Small bit of go-ifying
func (edge *TEdge) SetDx() {
	dy := edge.Top.Y - edge.Bot.Y
	if dy == 0 {
		edge.Dx = HORIZONTAL
		return
	}

	edge.Dx = (edge.Top.X - edge.Bot.Y) / dy
}

func SwapSides(edge1, edge2 *TEdge) {
	edge1.Side, edge2.Side = edge2.Side, edge1.Side
}

func SwapPolyIndexes(edge1, edge2 *TEdge) {
	edge1.OutIdx, edge2.OutIdx = edge2.OutIdx, edge1.OutIdx
}

func (edge *TEdge) TopX(currentY float64) float64 {
	if currentY == edge.Top.Y {
		return edge.Top.X
	}
	return edge.Bot.X + math.Round(edge.Dx*(currentY-edge.Bot.Y))
}

func IntersectPoint(edge1, edge2 *TEdge, point *Point) {
	var b1 float64
	var b2 float64

	if edge1.Dx == edge2.Dx {
		point.Y = edge1.Curr.Y
		point.X = edge1.TopX(point.Y)
		return
	} else if edge1.Dx == 0 {
		point.X = edge1.Bot.X
		if IsHorizontal(edge2) {
			point.Y = edge2.Bot.Y
		} else {
			b2 = edge2.Bot.Y - (edge2.Bot.X / edge2.Dx)
			point.Y = math.Round(point.X/edge2.Dx + b2)
		}
	} else if edge2.Dx == 0 {
		point.X = edge2.Bot.X
		if IsHorizontal(edge1) {
			point.Y = edge1.Bot.Y
		} else {
			b1 = edge1.Bot.Y - (edge1.Bot.X / edge1.Dx)
			point.Y = math.Round(point.X/edge1.Dx + b1)
		}
	} else {
		b1 = edge1.Bot.X - edge1.Bot.Y*edge1.Dx
		b2 = edge2.Bot.X - edge2.Bot.Y*edge2.Dx
		var q float64 = (b2 - b1) / (edge1.Dx - edge2.Dx)
		point.Y = math.Round(q)
		if math.Abs(edge1.Dx) < math.Abs(edge2.Dx) {
			point.X = math.Round(edge1.Dx*q + b1)
		} else {
			point.X = math.Round(edge2.Dx*q + b2)
		}
	}

	if point.Y < edge1.Top.Y || point.Y < edge2.Top.Y {
		if edge1.Top.Y > edge2.Top.Y {
			point.Y = edge1.Top.Y
		} else {
			point.Y = edge2.Top.Y
		}

		if math.Abs(edge1.Dx) < math.Abs(edge2.Dx) {
			point.X = edge1.TopX(point.Y)
		} else {
			point.X = edge2.TopX(point.Y)
		}
	}

	// don't allow point to be below curr.y (ie bottom of the scanbeam)
	if point.Y > edge1.Curr.Y {
		point.Y = edge1.Curr.Y
		// use the more vertical edge to derive x
		if math.Abs(edge1.Dx) > math.Abs(edge2.Dx) {
			point.X = edge2.TopX(point.Y)
		} else {
			point.X = edge1.TopX(point.Y)
		}
	}
}

func ReversePolyPtLinks(pp *OutPt) {
	if pp == nil {
		return
	}
	OrigIdx := pp.Idx
	pp1 := pp
	var pp2 *OutPt

	for {
		pp2 = pp1.Next
		pp1.Next = pp1.Prev
		pp1.Prev = pp2
		pp1 = pp2

		if pp1.Idx == OrigIdx {
			break
		}
	}
}

// continue https://github.com/slic3r/Slic3r/blob/master/xs/src/clipper.cpp
