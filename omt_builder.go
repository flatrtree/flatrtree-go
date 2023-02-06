package flatrtree

import (
	"errors"
	"fmt"
	"math"

	"github.com/furstenheim/nth_element/FloydRivest"
)

var _ Builder = &OMTBuilder{}

type OMTBuilder struct {
	count int
	refs  []int64
	boxes []float64

	nodeSizes [][]int
	nodeBoxes [][]float64
}

func NewOMTBuilder() *OMTBuilder {
	return &OMTBuilder{}
}

func (b *OMTBuilder) Add(ref int64, minX, minY, maxX, maxY float64) {
	b.count++
	b.refs = append(b.refs, ref)
	b.boxes = append(b.boxes, minX, minY, maxX, maxY)
}

func (b *OMTBuilder) Finish(degree int) (*RTree, error) {
	if degree < 2 {
		return nil, fmt.Errorf("degree < 2")
	}

	if b.count == 0 {
		return &RTree{}, nil
	}

	if len(b.refs) != b.count {
		return nil, errors.New("Finish called more than once")
	}

	targetHeight := int(math.Ceil(math.Log(float64(b.count)) / math.Log(float64(degree))))
	if targetHeight < 1 {
		targetHeight = 1
	}

	b.nodeSizes = make([][]int, targetHeight)
	b.nodeBoxes = make([][]float64, targetHeight)

	b.build(degree, 0, b.count, targetHeight-1)
	b.pack()

	return &RTree{
		count: b.count,
		refs:  b.refs,
		boxes: b.boxes,
	}, nil
}

func (b *OMTBuilder) build(degree, start, end, level int) (float64, float64, float64, float64) {
	N := end - start
	childCount := int(math.Ceil(float64(N) / math.Pow(float64(degree), float64(level))))

	minX := math.Inf(1)
	minY := math.Inf(1)
	maxY := math.Inf(-1)
	maxX := math.Inf(-1)

	if N <= childCount {
		for i := start; i < end; i++ {
			minX = math.Min(minX, b.boxes[i*4])
			minY = math.Min(minY, b.boxes[i*4+1])
			maxX = math.Max(maxX, b.boxes[i*4+2])
			maxY = math.Max(maxY, b.boxes[i*4+3])
		}

		b.addNode(level, N, minX, minY, maxX, maxY)
		return minX, minY, maxX, maxY
	}

	nodeSize := 0
	nodeCapacity := int(math.Ceil(float64(N) / float64(childCount)))
	sliceCapacity := nodeCapacity * int(math.Ceil(math.Sqrt(float64(childCount))))

	sortByX(b.refs[start:end], b.boxes[start*4:end*4], sliceCapacity)

	for sliceStart := start; sliceStart < end; sliceStart += sliceCapacity {
		sliceEnd := sliceStart + sliceCapacity
		if sliceEnd > end {
			sliceEnd = end
		}

		sortByY(b.refs[sliceStart:sliceEnd], b.boxes[sliceStart*4:sliceEnd*4], nodeCapacity)

		for childStart := sliceStart; childStart < sliceEnd; childStart += nodeCapacity {
			childEnd := childStart + nodeCapacity
			if childEnd > end {
				childEnd = end
			}

			childMinX, childMinY, childMaxX, childMaxY := b.build(degree, childStart, childEnd, level-1)

			minX = math.Min(minX, childMinX)
			minY = math.Min(minY, childMinY)
			maxX = math.Max(maxX, childMaxX)
			maxY = math.Max(maxY, childMaxY)
			nodeSize++
		}
	}

	b.addNode(level, nodeSize, minX, minY, maxX, maxY)
	return minX, minY, maxX, maxY
}

func (b *OMTBuilder) pack() {
	var ref int64
	b.refs = append(b.refs, ref)

	for level := 0; level < len(b.nodeSizes); level++ {
		for i := 0; i < len(b.nodeSizes[level]); i++ {
			ref += int64(4 * b.nodeSizes[level][i])
			b.refs = append(b.refs, ref)
			b.boxes = append(b.boxes,
				b.nodeBoxes[level][i*4],
				b.nodeBoxes[level][i*4+1],
				b.nodeBoxes[level][i*4+2],
				b.nodeBoxes[level][i*4+3],
			)
		}
	}
}

func (b *OMTBuilder) addNode(level, size int, minX, minY, maxX, maxY float64) {
	b.nodeSizes[level] = append(b.nodeSizes[level], size)
	b.nodeBoxes[level] = append(b.nodeBoxes[level], minX, minY, maxX, maxY)
}

//
// Sorting
//

func sortByX(refs []int64, boxes []float64, n int) {
	FloydRivest.Buckets(sortByCoord{
		sortBy: 0, // minX
		refs:   refs,
		boxes:  boxes,
	}, n)
}

func sortByY(refs []int64, boxes []float64, n int) {
	FloydRivest.Buckets(sortByCoord{
		sortBy: 1, // minY
		refs:   refs,
		boxes:  boxes,
	}, n)
}

type sortByCoord struct {
	sortBy int
	refs   []int64
	boxes  []float64
}

func (s sortByCoord) Len() int {
	return len(s.refs)
}

func (s sortByCoord) Less(i, j int) bool {
	return s.boxes[i*4+s.sortBy] < s.boxes[j*4+s.sortBy]
}

func (s sortByCoord) Swap(i, j int) {
	s.refs[i], s.refs[j] = s.refs[j], s.refs[i]

	s.boxes[i*4], s.boxes[j*4] = s.boxes[j*4], s.boxes[i*4]
	s.boxes[i*4+1], s.boxes[j*4+1] = s.boxes[j*4+1], s.boxes[i*4+1]
	s.boxes[i*4+2], s.boxes[j*4+2] = s.boxes[j*4+2], s.boxes[i*4+2]
	s.boxes[i*4+3], s.boxes[j*4+3] = s.boxes[j*4+3], s.boxes[i*4+3]
}
