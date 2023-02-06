package flatrtree

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	inside    = [2]float64{-73.649597, 45.51982}
	north     = [2]float64{-73.627625, 45.815401}
	northEast = [2]float64{-72.951965, 45.823057}
	east      = [2]float64{-72.927246, 45.512121}
	southEast = [2]float64{-72.946472, 45.154927}
	south     = [2]float64{-73.624878, 45.13168}
	southWest = [2]float64{-74.382935, 45.182037}
	west      = [2]float64{-74.374695, 45.494796}
	northWest = [2]float64{-74.344482, 45.811572}

	bboxMin = [2]float64{-74.19342, 45.265222}
	bboxMax = [2]float64{-73.157959, 45.704261}

	epsilon = 1e-5
)

func haversine(a, b [2]float64) float64 {
	return earthRadiusMeters * distRad(
		a[1]*math.Pi/180, a[0]*math.Pi/180,
		b[1]*math.Pi/180, b[0]*math.Pi/180,
	)
}

func TestGeodeticBoxDistInside(t *testing.T) {
	expected := 0.0
	actual := GeodeticBoxDist(
		inside[0], inside[1],
		bboxMin[0], bboxMin[1],
		bboxMax[0], bboxMax[1],
	)
	require.Equal(t, expected, actual)
}

func TestGeodeticBoxDistNorth(t *testing.T) {
	expected := haversine(north, [2]float64{north[0], bboxMax[1]})
	actual := GeodeticBoxDist(
		north[0], north[1],
		bboxMin[0], bboxMin[1],
		bboxMax[0], bboxMax[1],
	)
	require.InEpsilon(t, expected, actual, epsilon)
}

func TestGeodeticBoxDistNorthEast(t *testing.T) {
	expected := haversine(northEast, bboxMax)
	actual := GeodeticBoxDist(
		northEast[0], northEast[1],
		bboxMin[0], bboxMin[1],
		bboxMax[0], bboxMax[1],
	)
	require.InEpsilon(t, expected, actual, epsilon)
}

func TestGeodeticBoxDistEast(t *testing.T) {
	expected := haversine(east, [2]float64{bboxMax[0], east[1]})
	actual := GeodeticBoxDist(
		east[0], east[1],
		bboxMin[0], bboxMin[1],
		bboxMax[0], bboxMax[1],
	)
	require.InEpsilon(t, expected, actual, epsilon)
}

func TestGeodeticBoxDistSouthEast(t *testing.T) {
	expected := haversine(southEast, [2]float64{bboxMax[0], bboxMin[1]})
	actual := GeodeticBoxDist(
		southEast[0], southEast[1],
		bboxMin[0], bboxMin[1],
		bboxMax[0], bboxMax[1],
	)
	require.InEpsilon(t, expected, actual, epsilon)
}

func TestGeodeticBoxDistSouth(t *testing.T) {
	expected := haversine(south, [2]float64{south[0], bboxMin[1]})
	actual := GeodeticBoxDist(
		south[0], south[1],
		bboxMin[0], bboxMin[1],
		bboxMax[0], bboxMax[1],
	)
	require.InEpsilon(t, expected, actual, epsilon)
}

func TestGeodeticBoxDistSouthWest(t *testing.T) {
	expected := haversine(southWest, bboxMin)
	actual := GeodeticBoxDist(
		southWest[0], southWest[1],
		bboxMin[0], bboxMin[1],
		bboxMax[0], bboxMax[1],
	)
	require.InEpsilon(t, expected, actual, epsilon)
}

func TestGeodeticBoxDistWest(t *testing.T) {
	expected := haversine(west, [2]float64{bboxMin[0], west[1]})
	actual := GeodeticBoxDist(
		west[0], west[1],
		bboxMin[0], bboxMin[1],
		bboxMax[0], bboxMax[1],
	)
	require.InEpsilon(t, expected, actual, epsilon)
}

func TestGeodeticBoxDistNorthWest(t *testing.T) {
	expected := haversine(northWest, [2]float64{bboxMin[0], bboxMax[1]})
	actual := GeodeticBoxDist(
		northWest[0], northWest[1],
		bboxMin[0], bboxMin[1],
		bboxMax[0], bboxMax[1],
	)
	require.InEpsilon(t, expected, actual, epsilon)
}

func TestPlanarBoxDist(t *testing.T) {
	var (
		minX = 0.0
		minY = 0.0
		maxX = 2.0
		maxY = 2.0
		midY = (maxX - minX) / 2
		midy = (maxY - minY) / 2
	)

	// inside
	require.Equal(t, 0.0, PlanarBoxDist(midY, midy, minX, minY, maxX, maxY))
	// top
	require.Equal(t, 1.0, PlanarBoxDist(midY, maxY+1, minX, minY, maxX, maxY))
	// top right
	require.Equal(t, 2.0, PlanarBoxDist(maxX+1, maxY+1, minX, minY, maxX, maxY))
	// right
	require.Equal(t, 1.0, PlanarBoxDist(maxX+1, midy, minX, minY, maxX, maxY))
	// bottom right
	require.Equal(t, 2.0, PlanarBoxDist(maxX+1, minY-1, minX, minY, maxX, maxY))
	// bottom
	require.Equal(t, 1.0, PlanarBoxDist(midY, minY-1, minX, minY, maxX, maxY))
	// bottom left
	require.Equal(t, 2.0, PlanarBoxDist(minX-1, minY-1, minX, minY, maxX, maxY))
	// left
	require.Equal(t, 1.0, PlanarBoxDist(minX-1, midy, minX, minY, maxX, maxY))
	// top right
	require.Equal(t, 2.0, PlanarBoxDist(minX-1, maxY+1, minX, minY, maxX, maxY))
}
