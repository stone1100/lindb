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

	height          int
	startLevel      int
	nodeCountDense  int
	childCountDense int
}

func (ls *loudsSparse) Init(builder *Builder) {
	ls.height = builder.treeHeight()
	ls.startLevel = builder.getSparseStartLevel()
	ls.nodeCountDense = 0
	for level := 0; level < ls.startLevel; level++ {
		ls.nodeCountDense += builder.nodeCounts[level]
	}
	// fmt.Println(ls.nodeCountDense)
	if ls.startLevel == 0 {
		ls.childCountDense = 0
	} else {
		ls.childCountDense = ls.nodeCountDense + builder.nodeCounts[ls.startLevel] - 1
	}

	//TODO: dense

	// init louds-sparse labels
	ls.labels = NewLabelVector()
	ls.labels.Init(builder.getLabels(), ls.startLevel, ls.height)

	numNodesPerLevel := make([]int, ls.height)
	for level := range numNodesPerLevel {
		numNodesPerLevel[level] = len(builder.lsLabels[level])
	}
	// init louds-sparse has-child
	ls.hasChild = &BitVectorRank{}
	ls.hasChild.Init(rankSparseBlockSize, builder.getHasChildBits(), numNodesPerLevel, ls.startLevel, ls.height)

	// init louds-sparse louds
	ls.louds = &BitVectorSelect{}
	ls.louds.Init(builder.getLoudsBits(), numNodesPerLevel, ls.startLevel, ls.height)

	// init suffix TODO:...
	// ls.suffixes = &SuffixVector{}
	// ls.suffixes.Init(builder.hasSuffix, numNodesPerLevel, builder.suffixes, ls.startLevel, ls.height)
}

func (ls *loudsSparse) lookupKey(key []byte, inNodeNum int) (nodeNum int, ok bool) {
	nodeNum = inNodeNum
	pos := ls.getFirstLabelPos(nodeNum)
	fmt.Printf("start pos=%d\n", pos)

	level := 0
	for level = ls.startLevel; level < len(key); level++ {
		fmt.Printf("----,pos=%d,nodeSize=%d\n", pos, ls.nodeSize(pos))
		// check labels if exist
		if pos, ok = ls.labels.Search(key[level], pos, ls.nodeSize(pos)); !ok {
			return -1, false
		}
		// if trie branch terminates
		if !ls.hasChild.ReadBit(pos) {
			//FIXME: need check suffix
			return nodeNum, true
		}

		fmt.Printf("before=>node=%d,pos=%d\n", nodeNum, pos)
		// move to child
		nodeNum = ls.getChildNodeNum(pos)
		pos = ls.getFirstLabelPos(nodeNum)
		fmt.Printf("after=>node=%d,pos=%d\n", nodeNum, pos)
	}
	if ls.labels.Read(pos) == terminator && !ls.hasChild.ReadBit(pos) {
		//FIXME: need check suffix
		return nodeNum, true
	}
	return -1, false
}

func (ls *loudsSparse) getChildNodeNum(pos int) int {
	return int(ls.hasChild.Rank(uint32(pos))) + ls.childCountDense
}

// S-ChildNodePos(pos) = select1(S-LOUDS, rank1(S-HasChild, pos) + 1)
// nodeNum = ls.getChildNodeNum(pos) => rank1(S-HasChild, pos)
func (ls *loudsSparse) getFirstLabelPos(nodeNum int) int {
	return int(ls.louds.Select(uint32(nodeNum + 1 - ls.nodeCountDense)))
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
	// write start level
	binary.LittleEndian.PutUint32(bs[:], uint32(ls.startLevel))
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
	ls.startLevel = int(reader.ReadUint32())

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
