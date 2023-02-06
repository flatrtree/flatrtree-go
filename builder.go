package flatrtree

const DefaultDegree int = 10

type Builder interface {
	Add(ref int64, minX, minY, maxX, maxY float64)
	Finish(degree int) (*RTree, error)
}
