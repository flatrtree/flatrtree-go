package flatrtree

import (
	"github.com/invisiblefunnel/flatqueue-go/v2"
)

type RTree struct {
	count int
	refs  []int64
	boxes []float64
}

// Count returns the number of items in the index
func (r *RTree) Count() int {
	return r.count
}

// Search calls the iterf function for all items intersecting the
// search box. If iterf returns false the search will terminate.
func (r *RTree) Search(
	minX, minY, maxX, maxY float64,
	iterf func(ref int64) (next bool),
) {
	if iterf == nil {
		panic("iterf nil")
	}

	if r.count == 0 {
		return
	}

	rootNodeIdx := int64(len(r.boxes) - 4)
	if r.intersects(rootNodeIdx, minX, minY, maxX, maxY) {
		r.search(rootNodeIdx, minX, minY, maxX, maxY, iterf)
	}
}

func (r *RTree) search(
	nodeIdx int64,
	minX, minY, maxX, maxY float64,
	iterf func(ref int64) (next bool),
) bool {
	var (
		refIdx       int64 = nodeIdx / 4
		childNodeIdx int64
		childRefIdx  int64
		count        int64 = int64(r.count)
	)

	for childNodeIdx = r.refs[refIdx]; childNodeIdx < r.refs[refIdx+1]; childNodeIdx += 4 {
		if r.intersects(childNodeIdx, minX, minY, maxX, maxY) {
			childRefIdx = childNodeIdx / 4
			if childRefIdx < count {
				if !iterf(r.refs[childRefIdx]) {
					return false
				}
			} else {
				if !r.search(childNodeIdx, minX, minY, maxX, maxY, iterf) {
					return false
				}
			}
		}
	}

	return true
}

func (r *RTree) intersects(nodeIdx int64, minX, minY, maxX, maxY float64) bool {
	return !(maxX < r.boxes[nodeIdx] || maxY < r.boxes[nodeIdx+1] ||
		minX > r.boxes[nodeIdx+2] || minY > r.boxes[nodeIdx+3])
}

// Neighbors calls the iterf function for all items in ascending order of distance
// to the given coordinates. If iterf returns false the search will terminate.
//
// Distances are calculated with the boxDist function. The itemDist function
// can optionally be supplied to calculate a more accurate distance to items
// in the index. Take care that itemDist and boxDist return distances in the
// same units.
func (r *RTree) Neighbors(
	x, y float64,
	iterf func(ref int64, dist float64) (next bool),
	boxDist func(pX, pY, minX, minY, maxX, maxY float64) (dist float64),
	itemDist func(pX, pY float64, ref int64) (dist float64),
) {
	if iterf == nil {
		panic("iterf nil")
	}

	if boxDist == nil {
		panic("boxDist nil")
	}

	if r.count == 0 {
		return
	}

	var (
		queue        flatqueue.FlatQueue[int64, float64]
		refIdx       int64
		childRefIdx  int64
		childNodeIdx int64
		leafRefIdx   int64
		dist         float64
		count        int64 = int64(r.count)
	)

	rootRefIdx := int64(len(r.refs) - 2)
	queue.Push(rootRefIdx, 0)

	for queue.Len() > 0 {
		refIdx = queue.Pop()
		for childNodeIdx = r.refs[refIdx]; childNodeIdx < r.refs[refIdx+1]; childNodeIdx += 4 {
			childRefIdx = childNodeIdx / 4
			if childRefIdx < count && itemDist != nil {
				dist = itemDist(x, y, r.refs[childRefIdx])
			} else {
				dist = boxDist(
					x, y,
					r.boxes[childNodeIdx], r.boxes[childNodeIdx+1],
					r.boxes[childNodeIdx+2], r.boxes[childNodeIdx+3],
				)
			}
			queue.Push(childRefIdx, dist)
		}

		for queue.Len() > 0 && queue.Peek() < count {
			dist = queue.PeekValue()
			leafRefIdx = queue.Pop()
			if !iterf(r.refs[leafRefIdx], dist) {
				return
			}
		}
	}
}
