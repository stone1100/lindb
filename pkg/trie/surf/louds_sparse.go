package surf

import (
	"fmt"
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

func NewLoudsSparse() *loudsSparse {
	return &loudsSparse{
		labels:   NewLabelVector(),
		hasChild: &BitVectorRank{},
		louds:    &BitVectorSelect{},
		suffixes: &SuffixVector{},
		values:   &ValueVector{},
	}
}

func (ls *loudsSparse) Init(builder *Builder) {
	ls.height = builder.height
	ls.totalKeys = builder.totalKeys

	// init louds-sparse labels
	ls.labels.Init(builder.levels)

	// init louds-sparse has-child
	ls.hasChild.Init(rankSparseBlockSize, builder.levels, HasChild)

	// init louds-sparse louds
	ls.louds.Init(builder.levels, Louds)

	// init suffix
	ls.suffixes.Init(builder.levels, HasSuffix)

	// init values
	ls.values.Init(builder.levels)
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

func (ls *loudsSparse) unmarshal(buf []byte) (err error) {
	pos := 0
	ls.totalKeys = int(UnmarshalUint32(buf, pos))
	pos += 4
	ls.height = int(UnmarshalUint32(buf, pos))
	pos += 4

	// read labels
	if pos, err = ls.labels.unmarshal(buf, pos); err != nil {
		return err
	}
	// read has child
	if pos, err = ls.hasChild.unmarshal(buf, pos); err != nil {
		return nil
	}
	// read louds
	if pos, err = ls.louds.unmarshal(buf, pos); err != nil {
		return nil
	}
	// read suffixes
	if pos, err = ls.suffixes.unmarshal(buf, pos); err != nil {
		return nil
	}
	// read values
	if _, err = ls.values.unmarshal(ls.totalKeys, buf, pos); err != nil {
		return nil
	}
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
