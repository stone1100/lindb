package surf

type loudsDense struct {
	height int
}

func (ld *loudsDense) Init(builder *Builder) {

}

func (ld *loudsDense) lookupKey(key []byte) (int, bool) {
	nodeNum := 0
	// pos := 0
	for level := 0; level < ld.height; level++ {

	}
	// search will continue in LOUDS-Sparse
	return nodeNum, true
}
