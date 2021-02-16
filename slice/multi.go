package slice

// This is here in order to create tools for
// slices generics do not exist in go so
// we will have to incorporate the same tools
// for every slice type... Annoying.
// Even more annoying is dereferencing a pointer
// anytime you see (*pts) it's being dereferenced
// so we can mutate the slice.

// Points is a collection of points
type Points []*Point

// NewPoints will construct Points
func NewPoints() Points {
	points := make(Points, 0)
	return points
}

// GetCopy will return a copy of all points
func (pts Points) GetCopy() Points {
	copied := make(Points, len(pts))
	copy(copied, pts)
	return copied
}

// Empty will determine if the points are empty
func (pts Points) Empty() bool {
	return len(pts) == 0
}

// Clear will clear all points
func (pts Points) Clear() {
	pts = NewPoints()
}

// First will get the first entry
func (pts Points) First() *Point {
	return pts[0]
}

// Last will get the last entry
func (pts Points) Last() *Point {
	return pts[len(pts)-1]
}

// EntryAtIndex will get a point at an index. If a negative index is supplied it will return a
// point from the back of the array
func (pts Points) EntryAtIndex(index int) *Point {
	if index < 0 {
		return pts[len(pts)+index]
	}
	return pts[index]
}

// PreviousEntry will get the point previous to the supplied index
func (pts Points) PreviousEntry(index int) *Point {
	idx := index - 1
	if idx < 0 {
		idx = len(pts) - 1
	}
	return pts.EntryAtIndex(idx)
}

// NextEntry will get the next point to the supplied index
func (pts Points) NextEntry(index int) *Point {
	idx := index + 1
	if idx > len(pts)-1 {
		idx = 0
	}
	return pts.EntryAtIndex(idx)
}

// PopBack will pop the last point in the stack
func (pts *Points) PopBack() *Point {
	popped, newPoints := (*pts)[len(*pts)-1], (*pts)[:len(*pts)-1]
	*pts = newPoints
	return popped
}

// PopFront will pop the first point in the stack
func (pts *Points) PopFront() *Point {
	popped, newPoints := (*pts)[0], (*pts)[1:]
	*pts = newPoints
	return popped
}

// Push will append a point
func (pts *Points) Push(point ...*Point) {
	*pts = append(*pts, point...)
}

// PushFront will push a point to the front of the stack
func (pts *Points) PushFront(points ...*Point) {
	*pts = append(points, *pts...)
}

// EraseAt will delete an item at index
func (pts *Points) EraseAt(index int) {
	*pts = append((*pts)[:index], (*pts)[index+1:]...)
}

// Window returns a sliding window version of the points
func (pts Points) Window(size int) []Points {
	points := pts.GetCopy()
	if len(points) <= size {
		return []Points{points}
	}

	window := make([]Points, 0, len(points)-size+1)

	for i, j := 0, size; j <= len(points); i, j = i+1, j+1 {
		window = append(window, points[i:j])
	}
	return window
}

// Polygons is a collection of Polygons
type Polygons []*Polygon

// NewPolygons will construct Polygons
func NewPolygons() Polygons {
	polys := make(Polygons, 0)
	return polys
}

// GetCopy will return a copy of all Polygons
func (pts Polygons) GetCopy() Polygons {
	copied := make(Polygons, len(pts))
	copy(copied, pts)
	return copied
}

// Empty will determine if the points are empty
func (pts Polygons) Empty() bool {
	return len(pts) == 0
}

// Clear will clear all points
func (pts Polygons) Clear() {
	pts = NewPolygons()
}

// First will get the first entry
func (pts Polygons) First() *Polygon {
	return pts[0]
}

// Last will get the last entry
func (pts Polygons) Last() *Polygon {
	return pts[len(pts)-1]
}

// EntryAtIndex will get an Entry at an index. If a negative index is supplied it will return a
// Entry from the back of the array
func (pts Polygons) EntryAtIndex(index int) *Polygon {
	if index < 0 {
		return pts[len(pts)+index]
	}
	return pts[index]
}

// PreviousEntry will get the point previous to the supplied index
func (pts Polygons) PreviousEntry(index int) *Polygon {
	idx := index - 1
	if idx < 0 {
		idx = len(pts) - 1
	}
	return pts.EntryAtIndex(idx)
}

// NextEntry will get the next point to the supplied index
func (pts Polygons) NextEntry(index int) *Polygon {
	idx := index + 1
	if idx > len(pts)-1 {
		idx = 0
	}
	return pts.EntryAtIndex(idx)
}

// PopBack will pop the last point in the stack
func (pts *Polygons) PopBack() *Polygon {
	popped, newPolygons := (*pts)[len((*pts))-1], (*pts)[:len((*pts))-1]
	*pts = newPolygons
	return popped
}

// PopFront will pop the first point in the stack
func (pts *Polygons) PopFront() *Polygon {
	popped, newPolygons := (*pts)[0], (*pts)[1:]
	*pts = newPolygons
	return popped
}

// Push will append a point
func (pts *Polygons) Push(poly ...*Polygon) {
	*pts = append(*pts, poly...)
}

// PushFront will push a point to the front of the stack
func (pts *Polygons) PushFront(polys ...*Polygon) {
	*pts = append(polys, *pts...)
}

// EraseAt will delete an item at index
func (pts *Polygons) EraseAt(index int) {
	*pts = append((*pts)[:index], (*pts)[index+1:]...)
}
