package flatrtree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidDegree(t *testing.T) {
	for builderName, newBuilder := range testBuilders {
		for _, invalidDegree := range []int{-1, 0, 1} {
			t.Run(builderName, func(t *testing.T) {
				builder := newBuilder()
				builder.Add(0, 1, 1, 2, 2)
				index, err := builder.Finish(invalidDegree)
				require.Nil(t, index)
				require.NotNil(t, err)
			})
		}
	}
}

func TestFinishAgain(t *testing.T) {
	for builderName, newBuilder := range testBuilders {
		t.Run(builderName, func(t *testing.T) {
			builder := newBuilder()
			builder.Add(0, 1, 1, 2, 2)

			index, err := builder.Finish(DefaultDegree)
			require.NotNil(t, index)
			require.Nil(t, err)

			index, err = builder.Finish(DefaultDegree)
			require.Nil(t, index)
			require.NotNil(t, err)
			require.Contains(t, err.Error(), "called more than once")
		})
	}
}

func TestHilbertFunction(t *testing.T) {
	// min and max pairs
	require.Equal(t, uint32(0), hilbert(0, 0))
	require.Equal(t, uint32(1431655765), hilbert(0, 65535))
	require.Equal(t, uint32(4294967295), hilbert(65535, 0))
	require.Equal(t, uint32(2863311530), hilbert(65535, 65535))

	// the first few
	require.Equal(t, uint32(0), hilbert(0, 0))
	require.Equal(t, uint32(1), hilbert(1, 0))
	require.Equal(t, uint32(2), hilbert(1, 1))
	require.Equal(t, uint32(3), hilbert(0, 1))
	require.Equal(t, uint32(4), hilbert(0, 2))
	require.Equal(t, uint32(5), hilbert(0, 3))
	require.Equal(t, uint32(6), hilbert(1, 3))
	require.Equal(t, uint32(7), hilbert(1, 2))

	// random sample
	require.Equal(t, uint32(980776996), hilbert(2971, 17497))
	require.Equal(t, uint32(3277697163), hilbert(62026, 27915))
	require.Equal(t, uint32(1534664434), hilbert(13890, 60206))
	require.Equal(t, uint32(3525267956), hilbert(43827, 27885))
	require.Equal(t, uint32(1058908279), hilbert(2794, 32229))
	require.Equal(t, uint32(1141222447), hilbert(8443, 33752))
	require.Equal(t, uint32(876709497), hilbert(13931, 24390))
	require.Equal(t, uint32(4219223461), hilbert(63456, 10643))
	require.Equal(t, uint32(534213004), hilbert(18084, 14710))
	require.Equal(t, uint32(1194905159), hilbert(11890, 39641))
}
