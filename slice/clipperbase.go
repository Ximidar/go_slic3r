package slice

import (
	"errors"
	"math"
)

var ErrOutsideRange = errors.New("coordinate outside allowed range")

type ClipperBase struct {
	currentLM         int
	minimaList        LocalMinimums
	useFullRange      bool
	edges             [][]*TEdge
	preserveCollinear bool
	hasOpenPaths      bool
	polyOuts          []*OutRec
	activeEdges       *TEdge
	scanBeamList      []float64
}

func NewClipperBase() *ClipperBase {
	cb := new(ClipperBase)
	cb.currentLM = 0
	cb.useFullRange = false
	cb.minimaList = NewLocalMinimums()
	return cb
}

func (cb *ClipperBase) ProcessBound(e *TEdge, NextIsForward bool) *TEdge {
	result := e
	var horz *TEdge = nil

	if e.OutIdx == Skip {
		//if edges still remain in the current bound beyond the skip edge then
		//create another LocMin and call ProcessBound once more

		if NextIsForward {
			for e.Top.Y == e.Next.Bot.Y {
				e = e.Next
			}
			//don't include top horizontals when parsing a bound a second time,
			//they will be contained in the opposite bound ...
			for e != result && IsHorizontal(e) {
				e = e.Prev
			}
		} else {
			for e.Top.Y == e.Prev.Bot.Y {
				e = e.Prev
			}

			for e != result && IsHorizontal(e) {
				e = e.Next
			}
		}

		if e == result {
			if NextIsForward {
				result = e.Next
			} else {
				result = e.Prev
			}
		} else {
			if NextIsForward {
				e = result.Next
			} else {
				e = result.Prev
			}

			locMin := LocalMinimum{
				Y:          e.Bot.Y,
				LeftBound:  nil,
				RightBound: e,
			}
			e.WinDelta = 0
			result = cb.ProcessBound(e, NextIsForward)
			cb.minimaList = append(cb.minimaList, &locMin)
		}
		return result
	}

	var estart *TEdge = nil

	if IsHorizontal(e) {
		//We need to be careful with open paths because this may not be a
		//true local minima (ie E may be following a skip edge).
		//Also, consecutive horz. edges may start heading left before going right.
		if NextIsForward {
			estart = e.Prev
		} else {
			estart = e.Next
		}

		if IsHorizontal(estart) { //ie an adjoining horizontal skip edge
			if estart.Bot.X != e.Bot.X && estart.Top.X != e.Bot.X {
				e.ReverseHorizontal()
			} else if estart.Bot.X != e.Bot.X {
				e.ReverseHorizontal()
			}
		}
	}

	estart = e
	if NextIsForward {
		for result.Top.Y == result.Next.Bot.Y && result.Next.OutIdx != Skip {
			result = result.Next
		}
		if IsHorizontal(result) && result.Next.OutIdx != Skip {
			//nb: at the top of a bound, horizontals are added to the bound
			//only when the preceding edge attaches to the horizontal's left vertex
			//unless a Skip edge is encountered when that becomes the top divide
			horz = result
			for IsHorizontal(horz.Prev) {
				horz = horz.Prev
			}
			if horz.Prev.Top.X > result.Next.Top.X {
				result = horz.Prev
			}
		}
		for e != result {
			e.NextInLML = e.Next
			if IsHorizontal(e) && e != estart && e.Bot.X != e.Prev.Top.X {
				e.ReverseHorizontal()
			}
			e = e.Next
		}
		if IsHorizontal(e) && e != estart && e.Bot.X != e.Prev.Top.X {
			e.ReverseHorizontal()
		}
		result = result.Next
	} else {
		for result.Top.Y == result.Prev.Bot.Y && result.Prev.OutIdx != Skip {
			result = result.Prev
		}
		if IsHorizontal(result) && result.Prev.OutIdx != Skip {
			//nb: at the top of a bound, horizontals are added to the bound
			//only when the preceding edge attaches to the horizontal's left vertex
			//unless a Skip edge is encountered when that becomes the top divide
			horz = result
			for IsHorizontal(horz.Next) {
				horz = horz.Next
			}
			if horz.Next.Top.X == result.Prev.Top.X ||
				horz.Next.Top.X > result.Prev.Top.X {
				result = horz.Next
			}
		}
		for e != result {
			e.NextInLML = e.Prev
			if IsHorizontal(e) && e != estart && e.Bot.X != e.Next.Top.X {
				e.ReverseHorizontal()
			}
			e = e.Prev
		}
		if IsHorizontal(e) && e != estart && e.Bot.X != e.Next.Top.X {
			e.ReverseHorizontal()
		}
		result = result.Prev
	}

	return result
}

func (cb *ClipperBase) AddPath(pg *Polygon, Ptype PolyType, Closed bool) (bool, error) {
	if !Closed && Ptype == ptClip {
		return false, errors.New("add path: open path must be subject")
	}

	highI := len(pg.MP.Points) - 1
	if Closed {
		for highI > 0 && EqualPoints(pg.MP.Points.EntryAtIndex(highI), pg.MP.Points.EntryAtIndex(0)) {
			highI -= 1
		}
	}

	for highI > 0 && EqualPoints(pg.MP.Points.EntryAtIndex(highI), pg.MP.Points.EntryAtIndex(highI-1)) {
		highI -= 1
	}

	if (Closed && highI < 2) || (!Closed && highI < 1) {
		return false, nil
	}

	// create a new edge array
	var edges []*TEdge = make([]*TEdge, highI+1)
	isFlat := true

	// Basic Initialization
	edges[1].Curr = pg.MP.Points.EntryAtIndex(1)
	err1 := RangeTest(pg.MP.Points.First(), cb.useFullRange)
	err2 := RangeTest(pg.MP.Points.EntryAtIndex(highI), cb.useFullRange)
	if err1 != nil || err2 != nil {
		return false, errors.New("add path: range test failed")
	}
	edges[0].InitEdge(edges[1], edges[highI], pg.MP.Points.EntryAtIndex(0))
	edges[highI].InitEdge(edges[0], edges[highI-1], pg.MP.Points.EntryAtIndex(highI))

	for i := highI - 1; i >= 1; i-- {
		err := RangeTest(pg.MP.Points.EntryAtIndex(i), cb.useFullRange)
		if err != nil {
			return false, errors.New("add path: range test failed")
		}
		edges[i].InitEdge(edges[i+1], edges[i-1], pg.MP.Points.EntryAtIndex(i))
	}

	eStart := edges[0]

	//2. Remove duplicate vertices, and when closed collinear edges
	E, eLoopStop := eStart, eStart

	for {
		if E.Curr == E.Next.Curr && (Closed || E.Next != eStart) {
			if E == E.Next {
				break
			}
			if E == eStart {
				eStart = E.Next
			}

			E.RemoveEdge()
			eLoopStop = E
			continue
		}

		if E.Prev == E.Next {
			break //only two vertices
		} else if Closed &&
			SlopesEqual3Pt(E.Prev.Curr, E.Curr, E.Next.Curr) &&
			(!cb.preserveCollinear || !Pt2IsBetweenPt1AndPt3(E.Prev.Curr, E.Curr, E.Next.Curr)) {

			//Collinear edges are allowed for open paths but in closed paths
			//the default is to merge adjacent collinear edges into a single edge.
			//However, if the PreserveCollinear property is enabled, only overlapping
			//collinear edges (ie spikes) will be removed from closed paths.

			if E == eStart {
				eStart = E.Next
			}
			E.RemoveEdge()
			E = E.Prev
			eLoopStop = E
			continue
		}

		E = E.Next
		if E == eLoopStop || (!Closed && E.Next == eStart) {
			break
		}
	}

	if (!Closed && E == E.Next) || (Closed && E.Prev == E.Next) {
		return false, nil
	}

	if !Closed {
		cb.hasOpenPaths = true
		eStart.Prev.OutIdx = Skip
	}

	// 3. Do second stage of edge initialization
	E = eStart
	for {
		E.InitEdgeWithPolyType(Ptype)
		E = E.Next
		if isFlat && E.Curr.Y != eStart.Curr.Y {
			isFlat = false
		}

		if E == eStart {
			break
		}
	}

	//4. Finally, add edge bounds to LocalMinima list

	//Totally flat paths must be handled differently when adding them
	//to LocalMinima list to avoid endless loops etc ...
	if isFlat {
		if Closed {
			return false, nil
		}

		E.Prev.OutIdx = Skip
		locMin := &LocalMinimum{
			Y:          E.Bot.Y,
			LeftBound:  nil,
			RightBound: E,
		}
		locMin.RightBound.Side = esRight
		locMin.RightBound.WinDelta = 0

		for {
			if E.Bot.X != E.Prev.Top.X {
				E.ReverseHorizontal()
			}
			if E.Next.OutIdx == Skip {
				break
			}
			E.NextInLML = E.Next
			E = E.Next
		}

		cb.minimaList = append(cb.minimaList, locMin)
		cb.edges = append(cb.edges, edges)
		return true, nil
	}

	cb.edges = append(cb.edges, edges)
	leftBoundIsForward := false

	var EMin *TEdge = nil

	//workaround to avoid an endless loop in the while loop below when
	//open paths have matching start and end points ...
	if EqualPoints(E.Prev.Bot, E.Prev.Top) {
		E = E.Next
	}

	for {
		E.FindNextLocMin()
		if E == EMin {
			break
		} else if EMin == nil {
			EMin = E
		}
		//E and E.Prev now share a local minima (left aligned if horizontal).
		//Compare their slopes to find which starts which bound ...
		locMin := &LocalMinimum{
			Y: E.Bot.Y,
		}

		if E.Dx < E.Prev.Dx {
			locMin.LeftBound = E.Prev
			locMin.RightBound = E
			leftBoundIsForward = false //Q.nextInLML = Q.prev
		} else {
			locMin.LeftBound = E
			locMin.RightBound = E.Prev
			leftBoundIsForward = true //Q.nextInLML = Q.next
		}

		if !Closed {
			locMin.LeftBound.WinDelta = 0
		} else if locMin.LeftBound.Next == locMin.RightBound {
			locMin.LeftBound.WinDelta = -1
		} else {
			locMin.LeftBound.WinDelta = 1
		}

		locMin.RightBound.WinDelta = -locMin.LeftBound.WinDelta

		E = cb.ProcessBound(locMin.LeftBound, leftBoundIsForward)
		if E.OutIdx == Skip {
			E = cb.ProcessBound(E, leftBoundIsForward)
		}

		E2 := cb.ProcessBound(locMin.RightBound, !leftBoundIsForward)
		if E2.OutIdx == Skip {
			E2 = cb.ProcessBound(E2, !leftBoundIsForward)
		}

		if locMin.LeftBound.OutIdx == Skip {
			locMin.LeftBound = nil
		} else if locMin.RightBound.OutIdx == Skip {
			locMin.RightBound = nil
		}

		cb.minimaList = append(cb.minimaList, locMin)
		if !leftBoundIsForward {
			E = E2
		}
	}
	return true, nil
}

func (cb *ClipperBase) AddPaths(paths []*Polygon, Ptype PolyType, Closed bool) (result bool, err error) {
	for _, path := range paths {
		result, err = cb.AddPath(path, Ptype, Closed)
		if err != nil {
			break
		}
	}
	return result, err
}

func (cb *ClipperBase) Clear() {
	cb.DisposeLocalMinimaList()
	cb.edges = make([][]*TEdge, 0)
	cb.useFullRange = false
	cb.hasOpenPaths = false
}

func (cb *ClipperBase) Reset() {
	cb.currentLM = 0
	if cb.minimaList.EntryAtIndex(cb.currentLM) == cb.minimaList.Last() {
		return // Nothing to process
	}

	SortLocalMinimum(cb.minimaList)

	// clear / reset priority queue
	cb.scanBeamList = make([]float64, 0)

	for _, lm := range cb.minimaList {
		cb.InsertScanbeam(lm.Y)
		e := lm.LeftBound
		if e != nil {
			e.Curr = e.Bot
			e.Side = esLeft
			e.OutIdx = Unassigned
		}

		e = nil
		e = lm.RightBound
		if e != nil {
			e.Curr = e.Bot
			e.Side = esRight
			e.OutIdx = Unassigned
		}
	}

	cb.activeEdges = nil
	cb.currentLM = 0

}

func (cb *ClipperBase) DisposeLocalMinimaList() {
	cb.minimaList.Clear()
	cb.currentLM = 0
}

func (cb *ClipperBase) PopLocalMinima(y float64) (*LocalMinimum, bool) {
	if cb.minimaList.EntryAtIndex(cb.currentLM) == cb.minimaList.Last() ||
		(cb.minimaList.EntryAtIndex(cb.currentLM).Y != y) {
		return nil, false
	}
	locmin := cb.minimaList.EntryAtIndex(cb.currentLM)
	cb.currentLM += 1
	return locmin, true
}

func (cb *ClipperBase) GetBounds() FloatRect {
	result := FloatRect{}
	lmIdx := 0
	lm := cb.minimaList.EntryAtIndex(lmIdx)
	if lm == cb.minimaList.Last() {
		return FloatRect{
			Top:    0,
			Bottom: 0,
			Left:   0,
			Right:  0,
		}
	}

	result.Left = lm.LeftBound.Bot.X
	result.Top = lm.LeftBound.Bot.Y
	result.Right = lm.LeftBound.Bot.X
	result.Bottom = lm.LeftBound.Bot.Y

	for lm != cb.minimaList.Last() {
		result.Bottom = math.Max(result.Bottom, lm.LeftBound.Bot.Y)

		e := lm.LeftBound
		for {
			bottomE := e
			for e.NextInLML != nil {
				if e.Bot.X < result.Left {
					result.Left = e.Bot.X
				}
				if e.Bot.X > result.Right {
					result.Right = e.Bot.X
				}
				e = e.NextInLML
			}
			result.Left = math.Min(result.Left, e.Bot.X)
			result.Right = math.Max(result.Right, e.Bot.X)

			result.Left = math.Min(result.Left, e.Top.X)
			result.Right = math.Max(result.Right, e.Top.X)

			result.Top = math.Min(result.Top, e.Top.Y)

			if bottomE == lm.LeftBound {
				e = lm.RightBound
			} else {
				break
			}
			lmIdx++
			lm = cb.minimaList.EntryAtIndex(lmIdx)
		}
	}
	return result
}

func (cb *ClipperBase) InsertScanbeam(y float64) {
	cb.scanBeamList = append(cb.scanBeamList, y)
}

func (cb *ClipperBase) PopScanbeam() (float64, bool) {
	if len(cb.scanBeamList) == 0 {
		return 0, false
	}
	// pop first value and shift
	var y float64
	y, cb.scanBeamList = cb.scanBeamList[0], cb.scanBeamList[1:]
	for len(cb.scanBeamList) != 0 && y == cb.scanBeamList[0] {
		// remove duplicate points
		cb.scanBeamList = cb.scanBeamList[1:]
	}
	return y, true
}

func (cb *ClipperBase) DisposeAllOutRecs() {
	for idx := range cb.polyOuts {
		cb.DisposeOutRec(idx)
	}
	cb.polyOuts = make([]*OutRec, 0)
}

func (cb *ClipperBase) DisposeOutRec(idx int) {
	if len(cb.polyOuts) < idx {
		return
	}
	OutRec := cb.polyOuts[idx]
	if OutRec.Pts != nil {
		DisposeOutPts(OutRec.Pts)
	}

	// Delete index
	if idx < len(cb.polyOuts)-1 {
		copy(cb.polyOuts[idx:], cb.polyOuts[idx+1:])
	}
	cb.polyOuts[len(cb.polyOuts)-1] = nil // or the zero value of T
	cb.polyOuts = cb.polyOuts[:len(cb.polyOuts)-1]
}

// Continue at https://github.com/slic3r/Slic3r/blob/master/xs/src/clipper.cpp
