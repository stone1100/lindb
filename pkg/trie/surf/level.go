package surf

type BitmapType int

const (
	HasChild BitmapType = iota + 1
	Louds
	HasSuffix
)

type Level struct {
	// LOUDS-Sparse context: labels/hasChild/louds
	//
	// store all the branching labels for each trie node
	// lsLabels [][]byte
	lsLabels []byte
	// // one bit for each byte in labels to indicate whether
	// // a child branch continues(i.e. points to a sub-trie)
	// // or terminals(i.e. points to a value)
	lsHasChild []uint64
	// // one bit for each byte in labels to indicate if a lable
	// // is the first node in trie
	lsLouds []uint64

	// suffix
	hasSuffixes []uint64
	suffixes    [][]byte
	// value
	values []uint32
	// level node count
	item int
}

func NewLevel() *Level {
	return &Level{
		lsLabels:    []byte{},
		lsHasChild:  []uint64{},
		lsLouds:     []uint64{},
		hasSuffixes: []uint64{},
		suffixes:    [][]byte{},
		values:      []uint32{},
	}
}

func (l *Level) GetBitmap(t BitmapType) []uint64 {
	switch t {
	case HasChild:
		return l.lsHasChild
	case Louds:
		return l.lsLouds
	case HasSuffix:
		return l.hasSuffixes
	default:
		return []uint64{}
	}
}

func (l *Level) Reset() {
	l.lsLabels = l.lsLabels[:0]
	l.lsHasChild = l.lsHasChild[:0]
	l.lsLouds = l.lsLouds[:0]
	l.hasSuffixes = l.hasSuffixes[:0]
	l.suffixes = l.suffixes[:0]
	l.values = l.values[:0]
	l.item = 0
}
