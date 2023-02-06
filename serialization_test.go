package flatrtree

import (
	"math"
	"testing"

	allthecities "github.com/invisiblefunnel/all-the-cities-go"
	"github.com/stretchr/testify/require"
)

func TestSerializationRoundTrip(t *testing.T) {
	for _, tc := range createTestCases(t) {
		a := tc.index

		data, err := Serialize(a, 1)
		require.Nil(t, err)

		b, err := Deserialize(data)
		require.Nil(t, err)

		require.Equal(t, a.count, b.count)

		require.Equal(t, len(a.refs), len(b.refs))
		for i, key := range a.refs {
			require.Equal(t, key, b.refs[i])
		}

		require.Equal(t, len(a.boxes), len(b.boxes))
		for i, coord := range a.boxes {
			require.Equal(t, coord, b.boxes[i])
		}
	}
}

func TestSerializationRoundTripCities(t *testing.T) {
	cities, err := allthecities.Load()
	require.Nil(t, err)

	builder := NewHilbertBuilder()
	for i, city := range cities {
		builder.Add(int64(i), city.Lon, city.Lat, city.Lon, city.Lat)
	}

	before, err := builder.Finish(DefaultDegree)
	require.Nil(t, err)

	for prec := uint32(0); prec < 6; prec++ {
		data, err := Serialize(before, prec)
		require.Nil(t, err)

		after, err := Deserialize(data)
		require.Nil(t, err)

		require.Equal(t, before.count, after.count)

		require.Equal(t, len(before.refs), len(after.refs))
		for i, key := range before.refs {
			require.Equal(t, key, after.refs[i])
		}

		require.Equal(t, len(before.boxes), len(after.boxes))
		for i, coord := range before.boxes {
			require.InDelta(t, coord, after.boxes[i], math.Pow10(-int(prec)))
		}
	}
}

func TestSerializationRoundtripEmpty(t *testing.T) {
	before, err := NewHilbertBuilder().Finish(DefaultDegree)
	require.Nil(t, err)

	data, err := Serialize(before, 7)
	require.Nil(t, err)

	require.Equal(t, 2, len(data))

	after, err := Deserialize(data)
	require.Nil(t, err)

	require.Equal(t, 0, after.count)
	require.Equal(t, 0, len(after.refs))
	require.Equal(t, 0, len(after.boxes))
}

func TestDeserializeEmpty(t *testing.T) {
	rtree, err := Deserialize([]byte{})
	require.Nil(t, err)

	require.Equal(t, 0, rtree.count)
	require.Equal(t, 0, len(rtree.refs))
	require.Equal(t, 0, len(rtree.boxes))
}
