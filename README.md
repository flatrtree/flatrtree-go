# Flatrtree - Go

Flatrtree is a serialization format and set of libraries for reading and writing R-trees. It's directly inspired by [Flatbush](https://github.com/mourner/flatbush) and [FlatGeobuf](https://github.com/flatgeobuf/flatgeobuf), and aims to make tiny, portable R-trees accessible in new contexts.

- Store R-trees to disk or transport them over a network.
- Build R-trees in one language and query them in another.
- Query R-trees before reading the data they index.

## Installation

Requires Go 1.18 or later.

```console
$ go get github.com/flatrtree/flatrtree-go
```

## Usage

Flatrtree separates building and querying behavior. The builder doesn’t know how to query an index and the index doesn’t know how it was built. This is inspired by [FlatBuffers](https://google.github.io/flatbuffers/).

### Search

```golang
package main

import (
	"fmt"

	"github.com/flatrtree/flatrtree-go"
)

type myItem struct {
	name                   string
	minX, minY, maxX, maxY float64
}

func main() {
	items := []myItem{
		{"Manhattan", -74.04773, 40.682917, -73.906651, 40.879038},
		{"Bronx", -73.933606, 40.785357, -73.765332, 40.915533},
		{"Brooklyn", -74.041896, 40.56953, -73.833559, 40.739128},
		{"Staten Island", -74.255591, 40.496134, -74.049236, 40.648926},
		{"Queens", -73.96262, 40.541834, -73.700009, 40.801011},
	}

	builder := flatrtree.NewHilbertBuilder() // or OMTBuilder
	for i, item := range items {
		// The first argument is any integer reference to the item being indexed
		builder.Add(int64(i), item.minX, item.minY, item.maxX, item.maxY)
	}

	// Create an RTree from the builder
	index, err := builder.Finish(flatrtree.DefaultDegree)
	if err != nil {
		panic(err)
	}

	fmt.Println(index.Count() == len(items)) //=> true

	// Search area
	minX := items[1].minX
	minY := items[1].minY
	maxX := items[1].maxX
	maxY := items[1].maxY

	// Supply a function to be called for each search result
	iter := func(ref int64) bool {
		fmt.Println(items[ref].name)
		return true // return false to halt search
	}

	index.Search(minX, minY, maxX, maxY, iter)
	//=> Manhattan
	//=> Bronx
	//=> Queens
}
```

### Neighbors

```golang
package main

import (
	"fmt"

	"github.com/flatrtree/flatrtree-go"
)

type myItem struct {
	name                   string
	minX, minY, maxX, maxY float64
}

func main() {
	items := []myItem{
		{"Manhattan", -74.04773, 40.682917, -73.906651, 40.879038},
		{"Bronx", -73.933606, 40.785357, -73.765332, 40.915533},
		{"Brooklyn", -74.041896, 40.56953, -73.833559, 40.739128},
		{"Staten Island", -74.255591, 40.496134, -74.049236, 40.648926},
		{"Queens", -73.96262, 40.541834, -73.700009, 40.801011},
	}

	builder := flatrtree.NewHilbertBuilder() // or OMTBuilder
	for i, item := range items {
		builder.Add(int64(i), item.minX, item.minY, item.maxX, item.maxY)
	}

	index, err := builder.Finish(flatrtree.DefaultDegree)
	if err != nil {
		panic(err)
	}

	// Search point
	x := items[3].maxX
	y := items[3].maxY

	// Supply a function to be called for each result.
	// Distance units are determined by the given box distance function.
	iter := func(ref int64, dist float64) bool {
		fmt.Println(items[ref].name, dist)
		return true // customize your halt condition (usually based on result count or distance)
	}

	// Planar and geodetic box distance functions are included.
	boxDist := flatrtree.GeodeticBoxDist

	index.Neighbors(x, y, iter, boxDist, nil)
	//=> Staten Island 0
	//=> Brooklyn 619.2421164576169
	//=> Manhattan 3781.7657804872874
	//=> Queens 7307.393263218162
	//=> Bronx 18030.845120691552
}
```

## Serialization

Flatrtree uses [protocol buffers](https://protobuf.dev/) for serialization, taking advantage of varint encoding to reduce the output size in bytes. There are many tradeoffs to explore for serialization and this seems like a good place to start. It wouldn’t be hard to roll your own format with something like FlatBuffers if that better fit your needs.

```golang
package main

import "github.com/flatrtree/flatrtree-go"

type Item struct {
	name                   string
	minX, minY, maxX, maxY float64
}

func main() {
	items := []Item{
		// ...
	}

	builder := flatrtree.NewHilbertBuilder()
	for i, item := range items {
		builder.Add(int64(i), item.minX, item.minY, item.maxX, item.maxY)
	}

	index, err := builder.Finish(flatrtree.DefaultDegree)
	if err != nil {
		panic(err)
	}

	// Specify the decimal precision you need
	var precision uint32 = 7

	data, err := flatrtree.Serialize(index, precision)
	if err != nil {
		panic(err)
	}

	// Store or transport for later use

	index, err = flatrtree.Deserialize(data)
	if err != nil {
		panic(err)
	}
}
```
