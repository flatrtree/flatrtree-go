package flatrtree

import (
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStructure(t *testing.T) {
	for _, tc := range createTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			for nodeIndex := 0; nodeIndex < len(tc.index.boxes); nodeIndex += 4 {
				refIndex := nodeIndex / 4

				minX := tc.index.boxes[nodeIndex]
				minY := tc.index.boxes[nodeIndex+1]
				maxX := tc.index.boxes[nodeIndex+2]
				maxY := tc.index.boxes[nodeIndex+3]

				// correctly formatted bounds
				require.LessOrEqual(t, minX, maxX)
				require.LessOrEqual(t, minY, maxY)

				if refIndex < tc.index.count {
					// leaf node: test item id is correct
					id := tc.index.refs[refIndex]
					require.Equal(t, minX, tc.items[id*4])
					require.Equal(t, minY, tc.items[id*4+1])
					require.Equal(t, maxX, tc.items[id*4+2])
					require.Equal(t, maxY, tc.items[id*4+3])
				} else {
					// loop through the children [start, end)
					start := tc.index.refs[refIndex]
					end := tc.index.refs[refIndex+1]
					for childNodeIdx := start; childNodeIdx < end; childNodeIdx += 4 {
						// test for containment, not just intersection
						require.LessOrEqual(t, minX, tc.index.boxes[childNodeIdx])
						require.LessOrEqual(t, minY, tc.index.boxes[childNodeIdx+1])
						require.GreaterOrEqual(t, maxX, tc.index.boxes[childNodeIdx+2])
						require.GreaterOrEqual(t, maxY, tc.index.boxes[childNodeIdx+3])
					}
				}
			}
		})
	}
}

func TestSearch(t *testing.T) {
	for _, tc := range createTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.count; i++ {
				minX := tc.items[i*4]
				minY := tc.items[i*4+1]
				maxX := tc.items[i*4+2]
				maxY := tc.items[i*4+3]

				var actual, expected []int64

				tc.index.Search(minX, minY, maxX, maxY, func(ref int64) bool {
					actual = append(actual, ref)
					return true
				})
				require.Contains(t, actual, int64(i))

				// Do a full scan to ensure correctness
				for j := 0; j < tc.count; j++ {
					if !(maxX < tc.items[j*4] ||
						maxY < tc.items[j*4+1] ||
						minX > tc.items[j*4+2] ||
						minY > tc.items[j*4+3]) {
						expected = append(expected, int64(j))
					}
				}
				require.ElementsMatch(t, expected, actual)
			}
		})
	}
}

func TestSearchEverything(t *testing.T) {
	for _, tc := range createTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			var refs []int64
			tc.index.Search(math.Inf(-1), math.Inf(-1), math.Inf(1), math.Inf(1), func(ref int64) bool {
				refs = append(refs, ref)
				return true
			})

			uniq := make(map[int64]bool)
			for _, ref := range refs {
				uniq[ref] = true
			}

			require.Equal(t, tc.index.count, len(refs))
			require.Equal(t, len(refs), len(uniq))
		})
	}
}

func TestSearchEarlyTermination(t *testing.T) {
	for _, tc := range createTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			// cutoff the test search at about 1/4 of the items
			cutoff := int(math.Ceil(float64(tc.index.count) / float64(4)))

			count := 0
			tc.index.Search(math.Inf(-1), math.Inf(-1), math.Inf(1), math.Inf(1), func(ref int64) bool {
				count++
				return count < cutoff
			})

			require.Equal(t, cutoff, count)
		})
	}
}

func TestSearchNilIterfPanics(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()

	index, _ := createIndex(t, testBuilders["Hilbert"], 100, DefaultDegree)

	index.Search(0, 0, 0, 0, nil)
}

func TestNeighbors(t *testing.T) {
	for _, tc := range createTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.count; i++ {
				midX := (tc.items[i*4] + tc.items[i*4+2]) / 2
				midY := (tc.items[i*4+1] + tc.items[i*4+3]) / 2

				var expected []float64
				expectedRefsByDist := make(map[float64][]int64)
				for j := 0; j < tc.count; j++ {
					d := PlanarBoxDist(
						midX, midY,
						tc.items[j*4], tc.items[j*4+1], tc.items[j*4+2], tc.items[j*4+3],
					)
					expected = append(expected, d)
					expectedRefsByDist[d] = append(expectedRefsByDist[d], int64(j))
				}
				sort.Float64s(expected)

				var actual []float64
				actualRefsByDist := make(map[float64][]int64)
				tc.index.Neighbors(midX, midY, func(ref int64, dist float64) (next bool) {
					actual = append(actual, dist)
					actualRefsByDist[dist] = append(actualRefsByDist[dist], ref)
					return true
				}, PlanarBoxDist, nil)

				require.Equal(t, expected, actual)
				require.Equal(t, len(expectedRefsByDist), len(actualRefsByDist))

				for dist := range expectedRefsByDist {
					require.ElementsMatch(t, expectedRefsByDist[dist], actualRefsByDist[dist])
				}
			}
		})
	}
}

func TestNeighborsDist(t *testing.T) {
	for _, tc := range createTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			for _, testDist := range []float64{0, 10, math.Inf(1)} {
				for i := 0; i < tc.count; i++ {
					midX := (tc.items[i*4] + tc.items[i*4+2]) / 2
					midY := (tc.items[i*4+1] + tc.items[i*4+3]) / 2

					var expected []float64
					expectedRefsByDist := make(map[float64][]int64)
					for j := 0; j < len(tc.items)/4; j++ {
						d := PlanarBoxDist(
							midX, midY,
							tc.items[j*4], tc.items[j*4+1], tc.items[j*4+2], tc.items[j*4+3],
						)
						if d <= testDist {
							expected = append(expected, d)
							expectedRefsByDist[d] = append(expectedRefsByDist[d], int64(j))
						}
					}
					sort.Float64s(expected)

					var actual []float64
					actualRefsByDist := make(map[float64][]int64)
					tc.index.Neighbors(midX, midY, func(ref int64, dist float64) bool {
						if dist <= testDist {
							actual = append(actual, dist)
							actualRefsByDist[dist] = append(actualRefsByDist[dist], ref)
							return true
						}
						return false
					}, PlanarBoxDist, nil)

					require.Equal(t, expected, actual)
					require.Equal(t, len(expectedRefsByDist), len(actualRefsByDist))

					for dist := range expectedRefsByDist {
						require.ElementsMatch(t, expectedRefsByDist[dist], actualRefsByDist[dist])
					}
				}
			}
		})
	}
}

func TestNeighborsNilBoxDistPanics(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()

	index, _ := createIndex(t, testBuilders["Hilbert"], 100, DefaultDegree)
	iterf := func(int64, float64) bool { return true }

	index.Neighbors(0, 0, iterf, nil, nil)
}

func TestNeighborsNilIterfPanics(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()

	index, _ := createIndex(t, testBuilders["Hilbert"], 100, DefaultDegree)

	index.Neighbors(0, 0, nil, PlanarBoxDist, nil)
}

func TestNeighborsByItemDist(t *testing.T) {
	builder := NewHilbertBuilder()

	pX := -112.084665
	pY := 33.470112

	boxes := [][4]float64{
		// closest bounding box, farthest actual distance
		{-112.108612, 33.451423, -112.082519, 33.473262},
		// farthest bounding box, closest actual distance
		{-112.080888, 33.472976, -112.073764, 33.473048},
	}

	distances := []float64{
		1.204e+07,
		1.203e+07,
	}

	itemDist := func(lon, lat float64, ref int64) float64 {
		return distances[ref]
	}

	for i, box := range boxes {
		builder.Add(int64(i), box[0], box[1], box[2], box[3])
	}

	index, err := builder.Finish(DefaultDegree)
	require.Nil(t, err)

	t.Run("With", func(t *testing.T) {
		var actual []int64
		index.Neighbors(pX, pY, func(ref int64, dist float64) bool {
			actual = append(actual, ref)
			return true
		}, GeodeticBoxDist, itemDist)
		require.Equal(t, []int64{1, 0}, actual)
	})

	t.Run("Without", func(t *testing.T) {
		var actual []int64
		index.Neighbors(pX, pY, func(ref int64, dist float64) bool {
			actual = append(actual, ref)
			return true
		}, GeodeticBoxDist, nil)
		require.Equal(t, []int64{0, 1}, actual)
	})
}
