package flatrtree

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

var _ Builder = &HilbertBuilder{}

type HilbertBuilder struct {
	count                  int
	refs                   []int64
	boxes                  []float64
	minX, minY, maxX, maxY float64
}

func NewHilbertBuilder() *HilbertBuilder {
	return &HilbertBuilder{
		minX: math.Inf(1),
		minY: math.Inf(1),
		maxX: math.Inf(-1),
		maxY: math.Inf(-1),
	}
}

func (b *HilbertBuilder) Add(ref int64, minX, minY, maxX, maxY float64) {
	b.count++
	b.refs = append(b.refs, ref)
	b.boxes = append(b.boxes, minX, minY, maxX, maxY)
	b.minX = math.Min(b.minX, minX)
	b.minY = math.Min(b.minY, minY)
	b.maxX = math.Max(b.maxX, maxX)
	b.maxY = math.Max(b.maxY, maxY)
}

func (b *HilbertBuilder) Finish(degree int) (*RTree, error) {
	if degree < 2 {
		return nil, fmt.Errorf("degree < 2")
	}

	if b.count == 0 {
		return &RTree{}, nil
	}

	if len(b.refs) != b.count {
		return nil, errors.New("Finish called more than once")
	}

	b.sort()
	b.pack(degree)

	return &RTree{
		count: b.count,
		refs:  b.refs,
		boxes: b.boxes,
	}, nil
}

func (b *HilbertBuilder) sort() {
	const hilbertMax float64 = (1 << 16) - 1

	var (
		xScale, yScale float64
		midX, midY     float64
		x, y           uint32
	)

	width := b.maxX - b.minX
	if width > 0 {
		xScale = hilbertMax / width
	}

	height := b.maxY - b.minY
	if height > 0 {
		yScale = hilbertMax / height
	}

	hilbertValues := make([]uint32, b.count)
	for i := 0; i < b.count; i++ {
		midX = (b.boxes[i*4] + b.boxes[i*4+2]) / 2
		midY = (b.boxes[i*4+1] + b.boxes[i*4+3]) / 2
		x = uint32(math.Round(xScale * (midX - b.minX)))
		y = uint32(math.Round(yScale * (midY - b.minY)))
		hilbertValues[i] = hilbert(x, y)
	}

	sort.Sort(sortByValues{
		boxes:  b.boxes,
		refs:   b.refs,
		values: hilbertValues,
	})
}

func (b *HilbertBuilder) pack(degree int) {
	count := b.count
	numNodes := count
	start, end := int64(0), int64(len(b.boxes))
	b.refs = append(b.refs, start)

	var nodeMinX, nodeMinY, nodeMaxX, nodeMaxY float64
	for {
		for start < end {
			nodeMinX = math.Inf(1)
			nodeMinY = math.Inf(1)
			nodeMaxX = math.Inf(-1)
			nodeMaxY = math.Inf(-1)
			for j := 0; j < degree && start < end; j++ {
				nodeMinX = math.Min(nodeMinX, b.boxes[start])
				nodeMinY = math.Min(nodeMinY, b.boxes[start+1])
				nodeMaxX = math.Max(nodeMaxX, b.boxes[start+2])
				nodeMaxY = math.Max(nodeMaxY, b.boxes[start+3])
				start += 4
			}
			b.refs = append(b.refs, start)
			b.boxes = append(b.boxes, nodeMinX, nodeMinY, nodeMaxX, nodeMaxY)
		}

		count = int(math.Ceil(float64(count) / float64(degree)))
		numNodes += count
		end = int64(numNodes * 4)
		if count == 1 {
			break
		}
	}
}

// Based on public domain code at https://github.com/rawrunprotected/hilbert_curves
func hilbert(x, y uint32) uint32 {
	a := x ^ y
	b := 0xFFFF ^ a
	c := 0xFFFF ^ (x | y)
	d := x & (y ^ 0xFFFF)

	aa := a | (b >> 1)
	bb := (a >> 1) ^ a
	cc := ((c >> 1) ^ (b & (d >> 1))) ^ c
	dd := ((a & (c >> 1)) ^ (d >> 1)) ^ d

	a = aa
	b = bb
	c = cc
	d = dd
	aa = (a & (a >> 2)) ^ (b & (b >> 2))
	bb = (a & (b >> 2)) ^ (b & ((a ^ b) >> 2))
	cc ^= (a & (c >> 2)) ^ (b & (d >> 2))
	dd ^= (b & (c >> 2)) ^ ((a ^ b) & (d >> 2))

	a = aa
	b = bb
	c = cc
	d = dd
	aa = (a & (a >> 4)) ^ (b & (b >> 4))
	bb = (a & (b >> 4)) ^ (b & ((a ^ b) >> 4))
	cc ^= (a & (c >> 4)) ^ (b & (d >> 4))
	dd ^= (b & (c >> 4)) ^ ((a ^ b) & (d >> 4))

	a = aa
	b = bb
	c = cc
	d = dd
	cc ^= (a & (c >> 8)) ^ (b & (d >> 8))
	dd ^= (b & (c >> 8)) ^ ((a ^ b) & (d >> 8))

	a = cc ^ (cc >> 1)
	b = dd ^ (dd >> 1)

	i0 := x ^ y
	i1 := b | (0xFFFF ^ (i0 | a))

	i0 = (i0 | (i0 << 8)) & 0x00FF00FF
	i0 = (i0 | (i0 << 4)) & 0x0F0F0F0F
	i0 = (i0 | (i0 << 2)) & 0x33333333
	i0 = (i0 | (i0 << 1)) & 0x55555555

	i1 = (i1 | (i1 << 8)) & 0x00FF00FF
	i1 = (i1 | (i1 << 4)) & 0x0F0F0F0F
	i1 = (i1 | (i1 << 2)) & 0x33333333
	i1 = (i1 | (i1 << 1)) & 0x55555555

	return (i1 << 1) | i0
}

type sortByValues struct {
	refs   []int64
	boxes  []float64
	values []uint32
}

func (s sortByValues) Len() int {
	return len(s.values)
}

func (s sortByValues) Less(i, j int) bool {
	return s.values[i] < s.values[j]
}

func (s sortByValues) Swap(i, j int) {
	s.refs[i], s.refs[j] = s.refs[j], s.refs[i]

	s.boxes[i*4], s.boxes[j*4] = s.boxes[j*4], s.boxes[i*4]
	s.boxes[i*4+1], s.boxes[j*4+1] = s.boxes[j*4+1], s.boxes[i*4+1]
	s.boxes[i*4+2], s.boxes[j*4+2] = s.boxes[j*4+2], s.boxes[i*4+2]
	s.boxes[i*4+3], s.boxes[j*4+3] = s.boxes[j*4+3], s.boxes[i*4+3]

	s.values[i], s.values[j] = s.values[j], s.values[i]
}
