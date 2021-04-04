package slice

// When attempting to convert clipper.cpp/.hpp there were a bunch of constants
// that needed to be defined. Here they are!

type ClipType int

const (
	ctIntersection ClipType = iota
	ctUnion
	ctDifference
	ctXor
)

type PolyType int

const (
	ptSubject PolyType = iota
	ptClip
)

type PolyFillType int

const (
	pftEvenOdd PolyFillType = iota
	pftNonZero
	pftPositive
	pftNegative
)

type InitOptions int

const (
	ioReverseSolution InitOptions = iota
	ioStrictlySimple
	ioPreserveCollinear
)

type JoinType int

const (
	jtSquare JoinType = iota
	jtRound
	jtMiter
)

type EndType int

const (
	etClosedPolygon EndType = iota
	etClosedLine
	etOpenButt
	etOpenSquare
	etOpenRound
)

type EdgeSide int

const (
	esLeft EdgeSide = iota
	esRight
)

const pi float64 = 3.141592653589793238
const twoPi float64 = pi * 2
const defArcTolerance = 0.25

type Direction int

const (
	dRightToLeft Direction = iota
	dLeftToRight
)

const Unassigned = -1
const Skip = -2

const HORIZONTAL = -1.0e+40
const TOLERANCE = 1.0e-20
