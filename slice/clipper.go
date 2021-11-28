package slice

import (
	"errors"
	"fmt"
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

type OutRec struct {
	Idx       int
	IsHole    bool
	IsOpen    bool
	FirstLeft *OutRec
	PolyNd    *Polygon
	Pts       *OutPt
	BottomPt  *OutPt
}

type FloatRect struct {
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
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

func DisposeOutPts(pp *OutPt) {
	if pp == nil {
		return
	}

	pp.Prev.Next = nil
	for pp != nil {
		temp := pp
		pp = pp.Next
		temp.Next = nil
		temp.Prev = nil
		temp.Pt = nil
		temp.Idx = 0
	}
}

func (edge *TEdge) InitEdge(next, prev *TEdge, pt *Point) {
	edge.Next = next
	edge.Prev = prev
	edge.Curr = pt
	edge.OutIdx = 0
}

func (edge *TEdge) InitEdgeWithPolyType(polyType PolyType) {
	if edge.Curr.Y >= edge.Next.Curr.Y {
		edge.Bot = edge.Curr
		edge.Top = edge.Next.Curr
	} else {
		edge.Top = edge.Curr
		edge.Bot = edge.Next.Curr
	}
	edge.SetDx()
	edge.PolyType = polyType
}

// RemoveEdge will remove the current TEdge from the linked list
func (edge *TEdge) RemoveEdge() {
	edge.Prev.Next = edge.Next
	edge.Next.Prev = edge.Prev
	result := edge.Next
	edge.Prev = nil // flag as removed (see ClipperBase.Clear) (Maybe... who knows)
	edge = result
}

func (edge *TEdge) ReverseHorizontal() {
	//swap horizontal edges' Top and Bottom x's so they follow the natural
	//progression of the bounds - ie so their xbots will align with the
	//adjoining lower edge. [Helpful in the ProcessHorizontal() method.]
	edge.Top.X, edge.Bot.X = edge.Bot.X, edge.Top.X
}

func SwapPoints(pt1, pt2 *Point) (*Point, *Point) {
	return pt2, pt1
}

func GetOverlapSegment(pt1a, pt1b, pt2a, pt2b *Point) (result bool, pt1 *Point, pt2 *Point) {

	// precondition: segments are Collinear
	if math.Abs(pt1a.X-pt1b.X) > math.Abs(pt1a.Y-pt1b.Y) {
		if pt1a.X > pt1b.X {
			pt1a, pt1b = SwapPoints(pt1a, pt1b)
		}
		if pt2a.X > pt2b.X {
			pt2a, pt2b = SwapPoints(pt2a, pt2b)
		}
		if pt1a.X > pt2a.X {
			pt1 = pt1a
		} else {
			pt1 = pt2a
		}
		if pt1b.X > pt2b.X {
			pt2 = pt1b
		} else {
			pt2 = pt2b
		}
		result = pt1.X < pt2.X
		return
	}
	if pt1a.Y > pt1b.Y {
		pt1a, pt1b = SwapPoints(pt1a, pt1b)
	}
	if pt2a.Y > pt2b.Y {
		pt2a, pt2b = SwapPoints(pt2a, pt2b)
	}
	if pt1a.Y > pt2a.Y {
		pt1 = pt1a
	} else {
		pt1 = pt2a
	}
	if pt1b.Y > pt2b.Y {
		pt2 = pt1b
	} else {
		pt2 = pt2b
	}
	result = pt1.Y < pt2.Y
	return
}

func FirstIsBottomPt(btmPt1 *OutPt, btmPt2 *OutPt) bool {
	p := btmPt1.Prev
	for EqualPoints(p.Pt, btmPt1.Pt) && p.Idx != btmPt1.Idx {
		p = p.Prev
	}
	dx1p := math.Abs(GetDx2Pt(btmPt1.Pt, p.Pt))
	p = btmPt1.Next
	for EqualPoints(p.Pt, btmPt1.Pt) && p.Idx != btmPt1.Idx {
		p = p.Next
	}
	dx1n := math.Abs(GetDx2Pt(btmPt1.Pt, p.Pt))

	p = btmPt2.Prev
	for EqualPoints(p.Pt, btmPt2.Pt) && p.Idx != btmPt2.Idx {
		p = p.Prev
	}
	dx2p := math.Abs(GetDx2Pt(btmPt2.Pt, p.Pt))
	p = btmPt2.Next
	for EqualPoints(p.Pt, btmPt2.Pt) && p.Idx != btmPt2.Idx {
		p = p.Next
	}
	dx2n := math.Abs(GetDx2Pt(btmPt2.Pt, p.Pt))

	if math.Max(dx1p, dx1n) == math.Max(dx2p, dx2n) &&
		math.Min(dx1p, dx1n) == math.Min(dx2p, dx1n) {
		return btmPt1.Area() > 0
	}
	return (dx1p >= dx2p && dx1p >= dx2n) || (dx1n >= dx2p && dx1n >= dx2n)

}

func (pp *OutPt) GetBottomPt() *OutPt {
	var dups *OutPt = nil
	p := pp.Next

	for p.Idx != pp.Idx {
		if p.Pt.Y > pp.Pt.Y {
			pp = p
			dups = nil
		} else if p.Pt.Y == pp.Pt.Y && p.Pt.X <= pp.Pt.X {
			if p.Pt.X < pp.Pt.X {
				dups = nil
				pp = p
			} else {
				if p.Next.Idx != pp.Idx && p.Prev.Idx != pp.Idx {
					dups = p
				}
			}
		}
		p = p.Next
	}

	if dups != nil {
		for dups.Idx != p.Idx {
			if !FirstIsBottomPt(p, dups) {
				pp = dups
			}
			dups = dups.Next
			for dups.Pt != pp.Pt {
				dups = dups.Next
			}
		}
	}
	return pp
}

func Pt2IsBetweenPt1AndPt3(pt1, pt2, pt3 *Point) bool {
	if EqualPoints(pt1, pt3) || EqualPoints(pt1, pt2) || EqualPoints(pt3, pt2) {
		return false
	} else if pt1.X != pt3.X {
		return (pt2.X > pt1.X) == (pt2.X < pt3.X)
	}
	return (pt2.Y > pt1.Y) == (pt2.Y < pt3.Y)
}

func HorzSegmentsOverlap(seg1a, seg1b, seg2a, seg2b float64) (bool, float64, float64, float64, float64) {

	if seg1a > seg1b {
		seg1a, seg1b = seg1b, seg1a
	}
	if seg2a > seg2b {
		seg2a, seg2b = seg2b, seg2a
	}
	return (seg1a < seg2b) && (seg2a < seg1b), seg1a, seg1b, seg2a, seg2b

}

func RangeTest(pt *Point, useFullRange bool) error {
	if useFullRange {
		if pt.X > hiRange || pt.Y > hiRange || -pt.X > hiRange || -pt.Y > hiRange {
			return ErrOutsideRange
		}
	}
	if pt.X > loRange || pt.Y > loRange || -pt.X > loRange || -pt.Y > loRange {
		useFullRange = true
		return RangeTest(pt, useFullRange)
	}
	return nil
}

func (edge *TEdge) FindNextLocMin() *TEdge {
	e := edge
	for {
		for !EqualPoints(e.Bot, e.Prev.Bot) || EqualPoints(e.Curr, e.Top) {
			e = e.Next
		}

		if IsHorizontal(e) && !IsHorizontal(e.Prev) {
			break
		}

		e2 := e
		for IsHorizontal(e) {
			e = e.Next
		}

		if e.Top.Y == e.Prev.Bot.Y {
			continue // just an intermediate horizon
		}
		if e2.Prev.Bot.X < e.Bot.X {
			e = e2
		}
		break
	}
	return e
}

// Clipper Options
type ClipperOptions struct {
	ExecuteLocked     bool
	UseFullRange      bool
	ReverseOutput     bool
	StrictSimple      bool
	PreserveCollinear bool
	HasOpenPaths      bool
	ZFill             int64
}

// clipper
type Clipper struct {
	opt  ClipperOptions
	base *ClipperBase
}

// New Clipper Object
func NewClipper(options ClipperOptions) *Clipper {
	clip := new(Clipper)
	clip.opt = options
	clip.base = new(ClipperBase)

	return clip
}

// ZFill function has a callback
func (clip *Clipper) ZFillFunc() {
	// TODO introduce a way to set a callback for ZFill
	fmt.Println("ZFill Func not implemented yet.")
}

func (clip *Clipper) ExecutePolygon(clipType ClipType, solution *Polygon, fillType PolyFillType) error {
	if clip.opt.ExecuteLocked {
		return errors.New("Execution Locked")
	}

	clip.opt.ExecuteLocked = true
	solution.MP.Points.Clear() // Empty Solution

	err := clip.ExecuteInternal()

	if err != nil {
		fmt.Println("Clip Failed:", err)
		return err
	}
	clip.BuildResult(solution)

	clip.base.DisposeAllOutRecs()
	clip.opt.ExecuteLocked = false
	return nil
}

func (clip *Clipper) FixHoleLinkage(outrec *OutRec) {
	//skip OutRecs that (a) contain outermost polygons or
	//(b) already have the correct owner/child linkage ...
	if outrec.FirstLeft == nil ||
		(outrec.IsHole != outrec.FirstLeft.IsHole &&
			outrec.FirstLeft.Pts != nil) {
		return

	}

	orfl := outrec.FirstLeft
	for orfl != nil && ((orfl.IsHole == outrec.IsHole) || orfl.Pts == nil) {
		orfl = orfl.FirstLeft
	}
	outrec.FirstLeft = orfl
}

func (clip *Clipper) ExecuteInternal() error{
	clip.base.Reset()
	Maxima := clip.MaximaList()
	SortedEdges := 0

	botY, topY := 0, 0 


	if ! clip.base.PopScanbeam()
}

// LINE 1560
// continue https://github.com/slic3r/Slic3r/blob/master/xs/src/clipper.cpp
