package slice

import "errors"

var ErrOutsideRange = errors.New("coordinate outside allowed range")

type ClipperBase struct {
	currentLM         int
	minimaList        []*LocalMinimum
	useFullRange      bool
	edges             []*TEdge
	preserveCollinear bool
	hasOpenPaths      bool
	polyOuts          []*OutRec
	activeEdges       *TEdge
	scanBeamList      Points
}

func NewClipperBase() *ClipperBase {
	cb := new(ClipperBase)
	cb.currentLM = 0
	cb.useFullRange = false
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
