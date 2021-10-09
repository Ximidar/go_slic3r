package slice

import "sort"

type LocalMinimum struct {
	Y          float64
	LeftBound  *TEdge
	RightBound *TEdge
}

type LocalMininumSort LocalMinimums

func (lm LocalMininumSort) Len() int {
	return len(lm)
}

func (lm LocalMininumSort) Swap(i, j int) {
	lm[i], lm[j] = lm[j], lm[i]
}

func (lm LocalMininumSort) Less(i, j int) bool {
	return lm[i].Y < lm[j].Y
}

func SortLocalMinimum(lms LocalMinimums) {
	sort.Sort(LocalMininumSort(lms))
}

// LocalMinimums is a collection of points
type LocalMinimums []*LocalMinimum

// NewLocalMinimums will construct LocalMinimums
func NewLocalMinimums() LocalMinimums {
	points := make(LocalMinimums, 0)
	return points
}

// GetCopy will return a copy of all points
func (lm LocalMinimums) GetCopy() LocalMinimums {
	copied := make(LocalMinimums, len(lm))
	copy(copied, lm)
	return copied
}

// Empty will determine if the points are empty
func (lm LocalMinimums) Empty() bool {
	return len(lm) == 0
}

// Clear will clear all points
func (lm LocalMinimums) Clear() {
	lm = NewLocalMinimums()
}

// First will get the first entry
func (lm LocalMinimums) First() *LocalMinimum {
	return lm[0]
}

// Last will get the last entry
func (lm LocalMinimums) Last() *LocalMinimum {
	return lm[len(lm)-1]
}

// EntryAtIndex will get a point at an index. If a negative index is supplied it will return a
// point from the back of the array
func (lm LocalMinimums) EntryAtIndex(index int) *LocalMinimum {
	if index < 0 {
		return lm[len(lm)+index]
	}
	return lm[index]
}

// PreviousEntry will get the point previous to the supplied index
func (lm LocalMinimums) PreviousEntry(index int) *LocalMinimum {
	idx := index - 1
	if idx < 0 {
		idx = len(lm) - 1
	}
	return lm.EntryAtIndex(idx)
}

// NextEntry will get the next point to the supplied index
func (lm LocalMinimums) NextEntry(index int) *LocalMinimum {
	idx := index + 1
	if idx > len(lm)-1 {
		idx = 0
	}
	return lm.EntryAtIndex(idx)
}

// PopBack will pop the last point in the stack
func (lm *LocalMinimums) PopBack() *LocalMinimum {
	popped, newLocalMinimums := (*lm)[len(*lm)-1], (*lm)[:len(*lm)-1]
	*lm = newLocalMinimums
	return popped
}

// PopFront will pop the first point in the stack
func (lm *LocalMinimums) PopFront() *LocalMinimum {
	popped, newLocalMinimums := (*lm)[0], (*lm)[1:]
	*lm = newLocalMinimums
	return popped
}

// Push will append a point
func (lm *LocalMinimums) Push(point ...*LocalMinimum) {
	*lm = append(*lm, point...)
}

// PushFront will push a point to the front of the stack
func (lm *LocalMinimums) PushFront(points ...*LocalMinimum) {
	*lm = append(points, *lm...)
}

// EraseAt will delete an item at index
func (lm *LocalMinimums) EraseAt(index int) {
	*lm = append((*lm)[:index], (*lm)[index+1:]...)
}

// Window returns a sliding window version of the points
func (lm LocalMinimums) Window(size int) []LocalMinimums {
	points := lm.GetCopy()
	if len(points) <= size {
		return []LocalMinimums{points}
	}

	window := make([]LocalMinimums, 0, len(points)-size+1)

	for i, j := 0, size; j <= len(points); i, j = i+1, j+1 {
		window = append(window, points[i:j])
	}
	return window
}
