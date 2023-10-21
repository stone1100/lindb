package surf

import (
	"encoding/binary"
	"fmt"
	"io"
)

type loudsSparse struct {
	labels   *LabelVector
	hasChild *BitVectorRank
	louds    *BitVectorSelect
	suffixes *SuffixVector
	values   *ValueVector

	height    int
	totalKeys int
}

func (ls *loudsSparse) Init(builder *Builder) {
	ls.height = len(builder.lsLabels)
	ls.totalKeys = builder.totalKeys

	// init louds-sparse labels
	ls.labels = NewLabelVector()
	ls.labels.Init(builder.lsLabels)

	numNodesPerLevel := make([]int, ls.height)
	for level := range numNodesPerLevel {
		numNodesPerLevel[level] = len(builder.lsLabels[level])
	}
	// init louds-sparse has-child
	ls.hasChild = &BitVectorRank{}
	ls.hasChild.Init(rankSparseBlockSize, builder.bitmaps, HasChildIdx, numNodesPerLevel)

	// init louds-sparse louds
	ls.louds = &BitVectorSelect{}
	ls.louds.Init(builder.bitmaps, LoudsIdx, numNodesPerLevel)

	// init suffix
	ls.suffixes = &SuffixVector{}
	ls.suffixes.Init(builder.bitmaps, HasSuffixIdx, numNodesPerLevel, builder.suffixes)

	// init values
	ls.values = &ValueVector{}
	ls.values.Init(builder.values)
}

func (ls *loudsSparse) lookupKey(key []byte) (value uint32, result bool) {
	nodeNum := 0
	pos := ls.getFirstLabelPos(nodeNum)

	ok := false
	level := 0
	for ; level < len(key); level++ {
		// check labels if exist
		if pos, ok = ls.labels.Search(key[level], pos, ls.nodeSize(pos)); !ok {
			return
		}
		// if trie branch terminates
		if !ls.hasChild.ReadBit(pos) {
			if ok = ls.suffixes.CheckSuffix(key, level+1, pos); ok {
				value = ls.values.Get(ls.valuePos(pos))
				result = true
			}
			return
		}

		// move to child
		nodeNum = ls.getChildNodeNum(pos)
		pos = ls.getFirstLabelPos(nodeNum)
	}
	if ls.labels.Read(pos) == terminator && !ls.hasChild.ReadBit(pos) {
		if ok = ls.suffixes.CheckSuffix(key, level+1, pos); ok {
			value = ls.values.Get(ls.valuePos(pos))
			result = true
		}
		return
	}
	return
}

func (ls *loudsSparse) getChildNodeNum(pos int) int {
	return ls.hasChild.Rank(pos)
}

// S-ValuePos(pos) = pos - rank1(S-HasChild,pos)
func (ls *loudsSparse) valuePos(pos int) int {
	return pos - ls.hasChild.Rank(pos)
}

// S-ChildNodePos(pos) = select1(S-LOUDS, rank1(S-HasChild, pos) + 1)
// nodeNum = ls.getChildNodeNum(pos) => rank1(S-HasChild, pos)
func (ls *loudsSparse) getFirstLabelPos(nodeNum int) int {
	return int(ls.louds.Select(nodeNum + 1))
}

func (ls *loudsSparse) nodeSize(pos int) int {
	return ls.louds.DistanceToNextSetBit(pos)
}

func (ls *loudsSparse) isEndOfNode(pos int) bool {
	return pos == ls.louds.numBits-1 || ls.louds.ReadBit(pos+1)
}

func (ls *loudsSparse) write(w io.Writer) error {
	var (
		bs [4]byte
	)
	// write total keys
	binary.LittleEndian.PutUint32(bs[:], uint32(ls.totalKeys))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	// write height
	binary.LittleEndian.PutUint32(bs[:], uint32(ls.height))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	// write labels
	if err := ls.labels.write(w); err != nil {
		return err
	}
	// write has child
	if err := ls.hasChild.write(w); err != nil {
		return err
	}
	// write louds
	if err := ls.louds.write(w); err != nil {
		return err
	}
	// write suffixes
	if err := ls.suffixes.write(w); err != nil {
		return err
	}
	// write values
	if err := ls.values.write(w); err != nil {
		return err
	}
	return nil
}

func (ls *loudsSparse) unmarshal(buf []byte) (err error) {
	pos := 0
	ls.totalKeys = int(UnmarshalUint32(buf, pos))
	pos += 4
	ls.height = int(UnmarshalUint32(buf, pos))
	pos += 4

	// read labels
	labels := &LabelVector{}
	if pos, err = labels.unmarshal(buf, pos); err != nil {
		return err
	}
	ls.labels = labels
	// read has child
	hasChild := &BitVectorRank{}
	if pos, err = hasChild.unmarshal(buf, pos); err != nil {
		return nil
	}
	ls.hasChild = hasChild
	// read louds
	louds := &BitVectorSelect{}
	if pos, err = louds.unmarshal(buf, pos); err != nil {
		return nil
	}
	ls.louds = louds
	// read suffixes
	suffixes := &SuffixVector{}
	if pos, err = suffixes.unmarshal(buf, pos); err != nil {
		return nil
	}
	ls.suffixes = suffixes
	// read values
	values := &ValueVector{}
	if _, err = values.unmarshal(ls.totalKeys, buf, pos); err != nil {
		return nil
	}
	ls.values = values
	return nil
}

func (ls *loudsSparse) String() string {
	return fmt.Sprintf(`louse sparse:
labels:  %s
hasChild:%s
louds:   %s`,
		ls.labels.String(),
		ls.hasChild.String(),
		ls.louds.String())
}
