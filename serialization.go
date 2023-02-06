package flatrtree

import (
	"math"

	"github.com/flatrtree/flatrtree-go/internal"
	"google.golang.org/protobuf/proto"
)

func Serialize(index *RTree, precision uint32) ([]byte, error) {
	count := uint32(index.count)

	scale := math.Pow10(int(precision))

	boxes := make([]int64, len(index.boxes))
	for i := 0; i < len(index.boxes); i++ {
		boxes[i] = int64(math.Round(index.boxes[i] * scale))
	}

	// Note: I did not see a performance improvement using
	// vtprotobuf for serialization. Any ideas?
	return proto.Marshal(&internal.RTree{
		Count:     count,
		Refs:      index.refs,
		Boxes:     boxes,
		Precision: precision,
	})
}

func Deserialize(b []byte) (*RTree, error) {
	msg := &internal.RTree{}

	// Note: vtprotobuf is much faster than stock
	// protobuf for deserialization
	if err := msg.UnmarshalVT(b); err != nil {
		return nil, err
	}

	count := int(msg.GetCount())

	scale := math.Pow10(int(msg.GetPrecision()))

	msgBoxes := msg.GetBoxes()
	boxes := make([]float64, len(msgBoxes))
	for i := 0; i < len(msgBoxes); i++ {
		boxes[i] = float64(msgBoxes[i]) / scale
	}

	return &RTree{
		count: count,
		refs:  msg.GetRefs(),
		boxes: boxes,
	}, nil
}
