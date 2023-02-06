package flatrtree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	allthecities "github.com/invisiblefunnel/all-the-cities-go"
	"github.com/stretchr/testify/require"
)

var cities []allthecities.City

func init() {
	var err error
	if cities, err = allthecities.Load(); err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixMicro())

	rand.Shuffle(len(cities), func(i, j int) {
		cities[i], cities[j] = cities[j], cities[i]
	})
}

func Benchmark_Build(b *testing.B) {
	for builderName, newBuilder := range testBuilders {
		for _, degree := range testDegrees {
			b.Run(fmt.Sprintf("%v/deg=%d", builderName, degree), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					builder := newBuilder()
					for i, city := range cities {
						builder.Add(int64(i), city.Lon, city.Lat, city.Lon, city.Lat)
					}
					builder.Finish(DefaultDegree)
				}
			})
		}
	}
}

func Benchmark_Serialize(b *testing.B) {
	for builderName, newBuilder := range testBuilders {
		for _, degree := range testDegrees {
			builder := newBuilder()
			for i, city := range cities {
				builder.Add(int64(i), city.Lon, city.Lat, city.Lon, city.Lat)
			}
			rtree, err := builder.Finish(degree)
			require.Nil(b, err)

			// sanity check
			_, err = Serialize(rtree, 5)
			require.Nil(b, err)

			b.Run(fmt.Sprintf("%v/deg=%d", builderName, degree), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, _ = Serialize(rtree, 5)
				}
			})
		}
	}
}

func Benchmark_Deserialize(b *testing.B) {
	for builderName, newBuilder := range testBuilders {
		for _, degree := range testDegrees {
			builder := newBuilder()
			for ref, city := range cities {
				builder.Add(int64(ref), city.Lon, city.Lat, city.Lon, city.Lat)
			}

			rtree, err := builder.Finish(degree)
			require.Nil(b, err)

			data, err := Serialize(rtree, 5)
			require.Nil(b, err)

			// sanity check
			_, err = Deserialize(data)
			require.Nil(b, err)

			b.Run(fmt.Sprintf("%v/deg=%d", builderName, degree), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, _ = Deserialize(data)
				}
			})
		}
	}
}

func Benchmark_Search(b *testing.B) {
	// California-ish
	minX := -124.628906
	minY := 32.509762
	maxX := -113.818359
	maxY := 42.261049
	iterf := func(int64) bool { return true }

	for builderName, newBuilder := range testBuilders {
		for _, degree := range testDegrees {
			builder := newBuilder()
			for ref, city := range cities {
				builder.Add(int64(ref), city.Lon, city.Lat, city.Lon, city.Lat)
			}
			rtree, err := builder.Finish(degree)
			require.Nil(b, err)

			b.Run(fmt.Sprintf("%v/deg=%d", builderName, degree), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					rtree.Search(minX, minY, maxX, maxY, iterf)
				}
			})
		}
	}
}

func Benchmark_NeighborsPlanar(b *testing.B) {
	// A point in Central Park, NYC
	pX := -73.97197723388672
	pY := 40.774041868909734
	iterf := func(int64, float64) bool { return false }

	for builderName, newBuilder := range testBuilders {
		for _, degree := range testDegrees {
			builder := newBuilder()
			for ref, city := range cities {
				builder.Add(int64(ref), city.Lon, city.Lat, city.Lon, city.Lat)
			}
			rtree, err := builder.Finish(degree)
			require.Nil(b, err)

			b.Run(fmt.Sprintf("%v/deg=%d", builderName, degree), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					// Find the nearest neighbor and halt search
					rtree.Neighbors(pX, pY, iterf, PlanarBoxDist, nil)
				}
			})
		}
	}
}

func Benchmark_NeighborsGeodetic(b *testing.B) {
	// A point in Central Park, NYC
	pX := -73.97197723388672
	pY := 40.774041868909734
	iterf := func(int64, float64) bool { return false }

	for builderName, newBuilder := range testBuilders {
		for _, degree := range testDegrees {
			builder := newBuilder()
			for ref, city := range cities {
				builder.Add(int64(ref), city.Lon, city.Lat, city.Lon, city.Lat)
			}
			rtree, err := builder.Finish(degree)
			require.Nil(b, err)

			b.Run(fmt.Sprintf("%v/deg=%d", builderName, degree), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					// Find the nearest neighbor and halt search
					rtree.Neighbors(pX, pY, iterf, GeodeticBoxDist, nil)
				}
			})
		}
	}
}
