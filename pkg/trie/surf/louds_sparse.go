package surf

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/lindb/lindb/pkg/stream"
)

type loudsSparse struct {
	labels   *LabelVector
	hasChild *BitVectorRank
	louds    *BitVectorSelect
	suffixes *SuffixVector
	values   *ValueVector

	height int
}

func (ls *loudsSparse) Init(builder *Builder) {
	ls.height = builder.treeHeight()

	// init louds-sparse labels
	ls.labels = NewLabelVector()
	ls.labels.Init(builder.getLabels())

	numNodesPerLevel := make([]int, ls.height)
	for level := range numNodesPerLevel {
		numNodesPerLevel[level] = len(builder.lsLabels[level])
	}
	// init louds-sparse has-child
	ls.hasChild = &BitVectorRank{}
	ls.hasChild.Init(rankSparseBlockSize, builder.getHasChildBits(), numNodesPerLevel)

	// init louds-sparse louds
	ls.louds = &BitVectorSelect{}
	ls.louds.Init(builder.getLoudsBits(), numNodesPerLevel)

	// init suffix
	ls.suffixes = &SuffixVector{}
	ls.suffixes.Init(builder.hasSuffix, numNodesPerLevel, builder.suffixes)

	// init values
	ls.values = &ValueVector{}
	ls.values.Init(builder.values)
}

func (ls *loudsSparse) lookupKey(key []byte) (value uint32, ok bool) {
	nodeNum := 0
	pos := ls.getFirstLabelPos(nodeNum)

	level := 0
	for ; level < len(key); level++ {
		// check labels if exist
		if pos, ok = ls.labels.Search(key[level], pos, ls.nodeSize(pos)); !ok {
			return
		}
		// if trie branch terminates
		if !ls.hasChild.ReadBit(pos) {
			if ok = ls.suffixes.CheckSuffix(key, level, pos); ok {
				value = ls.values.Get(ls.valuePos(pos))
				ok = true
			}
			return
		}

		// move to child
		nodeNum = ls.getChildNodeNum(pos)
		pos = ls.getFirstLabelPos(nodeNum)
	}
	if ls.labels.Read(pos) == terminator && !ls.hasChild.ReadBit(pos) {
		if ok = ls.suffixes.CheckSuffix(key, level, pos); ok {
			value = ls.values.Get(ls.valuePos(pos))
			ok = true
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
	return int(ls.louds.Select(uint32(nodeNum + 1)))
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
	return nil
}

func (ls *loudsSparse) unmarshal(reader *stream.Reader) (err error) {
	ls.height = int(reader.ReadUint32())

	labels := &LabelVector{}
	if err := labels.unmarshal(reader); err != nil {
		return err
	}
	ls.labels = labels

	hasChild := &BitVectorRank{}
	if err := hasChild.unmarshal(reader); err != nil {
		return nil
	}
	ls.hasChild = hasChild

	louds := &BitVectorSelect{}
	if err := louds.unmarshal(reader); err != nil {
		return nil
	}
	ls.louds = louds
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
