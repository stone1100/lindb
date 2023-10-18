package surf

import (
	"bytes"
	"fmt"
)

const (
	bitsSize = 64
	// terminator($) indicates the situation where a prefix strings
	// leading to a node is also a valid key
	terminator       = 0xff
	fanout           = 256
	sparseDenseRatio = 16
)

var (
	emptyKey = []byte{}
)

// build baisc SuRF(input keys must be sorted)
// include dense = true
// sprase dense ratio = 16
// suffix type = none
// hash suffix len = 0
// real suffix len = 0
type Builder struct {
	// trie level < sparse start level  => LOUDS-Dense
	// trie level >= sparse start level => LOUDS-Sparse

	// LOUDS-Dense context: labels/hasChild/hasPrefix
	//
	// bitmap stores the branching labels for each node.
	// (0<=label<=255), 256 is terminator
	ldLabels [][]uint64
	// one bit for each byte in labels to indicate whether
	// a child branch continues(i.e. points to a sub-trie)
	// or terminals(i.e. points to a value)
	ldHasChild [][]uint64
	// one bit per node to indicates whether the perfix that leads
	// to the node is also a valid key.
	ldIsPrefixKey [][]uint64

	// LOUDS-Sparse context: labels/hasChild/louds
	//
	// store all the branching labels for each trie node
	lsLabels [][]byte
	// like LOUDS-Dense's hasChild
	lsHasChild [][]uint64
	// one bit for each byte in labels to indicate if a lable
	// is the first node in trie
	lsLouds [][]uint64

	nodeCounts           []int // first node counts??? TODO:
	isLastItemTerminator []bool

	hasSuffix    [][]uint64
	suffixes     [][][]byte
	suffixCounts []int

	sparseStartLevel int
}

func (b *Builder) Build(keys [][]byte) {
	b.buildSparse(keys)

	// b.determineCutoffLevel()
	// fmt.Println(b.sparseStartLevel)
	// b.buildDense()
}

func (b *Builder) buildDense() {
	for level := 0; level < b.sparseStartLevel; level++ {
		b.initDenseVectors(level)
		if b.numNodes(level) == 0 {
			continue
		}
		nodeNum := 0
		if b.isTerminator(level, 0) {
			setBit(b.ldIsPrefixKey[level], 0)
		} else {
			b.setLabelAndChildIndicatorBitmap(level, nodeNum, 0)
		}

		for pos := 1; pos < b.numNodes(level); pos++ {
			if b.isStartOfNode(level, pos) {
				nodeNum++
				if b.isTerminator(level, pos) {
					setBit(b.ldIsPrefixKey[level], nodeNum)
					continue
				}
			}
			b.setLabelAndChildIndicatorBitmap(level, nodeNum, pos)
		}
	}
}

func (b *Builder) setLabelAndChildIndicatorBitmap(level, nodeNum, pos int) {
	label := b.lsLabels[level][pos]
	setBit(b.ldLabels[level], nodeNum*fanout+int(label))
	if readBit(b.lsHasChild[level], pos) {
		setBit(b.ldHasChild[level], nodeNum*fanout+int(label))
	}
}

func (b *Builder) initDenseVectors(level int) {
	b.ldLabels = append(b.ldLabels, []uint64{})
	b.ldHasChild = append(b.ldHasChild, []uint64{})
	b.ldIsPrefixKey = append(b.ldIsPrefixKey, []uint64{})
	for nc := 0; nc < b.nodeCounts[level]; nc++ {
		for i := 0; i < fanout; i += bitsSize {
			b.ldLabels[level] = append(b.ldLabels[level], 0)
			b.ldHasChild[level] = append(b.ldHasChild[level], 0)
		}
		if nc%bitsSize == 0 {
			b.ldIsPrefixKey[level] = append(b.ldIsPrefixKey[level], 0)
		}
	}
}

func (b *Builder) determineCutoffLevel() {
	cutoffLevel := 0
	denseMem := b.computeDenseMem(cutoffLevel)
	sparseMem := b.computeSparseMem(cutoffLevel)
	for cutoffLevel < b.treeHeight() && denseMem*sparseDenseRatio < sparseMem {
		fmt.Printf("mm=%d,sm=%d,cutleve=%d\n", denseMem, sparseMem, cutoffLevel)
		cutoffLevel++
		denseMem = b.computeDenseMem(cutoffLevel)
		sparseMem = b.computeSparseMem(cutoffLevel)
	}
	// cutoffLevel--
	b.sparseStartLevel = cutoffLevel
}

func (b *Builder) computeDenseMem(downToLevel int) int {
	mem := 0
	for level := 0; level < downToLevel; level++ {
		mem += (2 * fanout * b.nodeCounts[level])
		if level > 0 {
			mem += (b.nodeCounts[level-1]/8 + 1)
		}
		mem += (b.suffixCounts[level] * b.getSuffixLen() / 8)
	}
	return mem
}

func (b *Builder) computeSparseMem(startLevel int) int {
	mem := 0
	for level := startLevel; level < b.treeHeight(); level++ {
		numNodes := len(b.lsLabels[level])
		mem += (numNodes + 2*numNodes/8 + 1)
		mem += (b.suffixCounts[level] * b.getSuffixLen() / 8)
	}
	return mem
}

func (b *Builder) buildSparse(keys [][]byte) {
	for i := 0; i < len(keys); i++ {
		// skip common prefix
		level := b.skipCommonPrefix(keys[i])
		if i < len(keys)-1 {
			level = b.insertKeyBytesToTrieUntilUnique(keys[i], keys[i+1], level)
		} else {
			// for last key, there is no successor key in the list
			level = b.insertKeyBytesToTrieUntilUnique(keys[i], emptyKey, level)
		}
		// insert suffix if has suffix
		if level < len(keys[i]) {
			b.insertSuffix(keys[i], level)
		}
	}
}

func (b *Builder) insertKeyBytesToTrieUntilUnique(key []byte, nextKey []byte, startLevel int) (level int) {
	level = startLevel
	isStartOfNode := false
	isTerm := false

	// if it is the start of level, the louds bit needs to be set
	if b.isLevelEmpty(level) {
		isStartOfNode = true
	}
	// after skipping the common prefix, the first following byte
	// shoud be in an the node as the previous key.
	b.insertKeyByte(key[level], level, isStartOfNode, isTerm)
	level++ // goto next
	if level > len(nextKey) || !b.isSameKey(key[0:level], nextKey[0:level]) {
		return level
	}

	// all the follwing bytes inserted must be the start of a new node,
	// becase generate new sub trie
	isStartOfNode = true
	for level < len(key) && level < len(nextKey) && key[level] == nextKey[level] {
		b.insertKeyByte(key[level], level, isStartOfNode, isTerm)
		level++
	}

	// the last byte inserted makes key unique in the tire
	if level < len(key) {
		b.insertKeyByte(key[level], level, isStartOfNode, isTerm)
		level++ // goto next
	} else {
		// insert terminator char
		isTerm = true
		b.insertKeyByte(terminator, level, isStartOfNode, isTerm)
	}

	return level
}

func (b *Builder) insertKeyByte(key byte, level int, isStartOfNode, isTerm bool) {
	// level should be at most equal to tree height
	b.ensureLevel(level)

	// store paren has child
	// sets parent node's child indicator
	if level > 0 {
		// all keys is sorted, so new key will append right
		setBit(b.lsHasChild[level-1], b.numNodes(level-1)-1)
	}
	// store label
	b.lsLabels[level] = append(b.lsLabels[level], key)

	// store louds
	if isStartOfNode {
		setBit(b.lsLouds[level], b.numNodes(level)-1)
		b.nodeCounts[level]++
	}
	b.isLastItemTerminator[level] = isTerm

	b.moveToNextNodeSlot(level)
}

func (b *Builder) moveToNextNodeSlot(level int) {
	numNodes := b.numNodes(level)
	if numNodes%bitsSize == 0 {
		// put next slot for bit vector context
		b.lsHasChild[level] = append(b.lsHasChild[level], 0)
		b.lsLouds[level] = append(b.lsLouds[level], 0)

		b.hasSuffix[level] = append(b.hasSuffix[level], 0)
	}
}

// isLevelEmpty returns whether level is empty.
func (b *Builder) isLevelEmpty(level int) bool {
	return level >= b.treeHeight() || len(b.lsLabels[level]) == 0
}

func (b *Builder) insertSuffix(key []byte, level int) {
	b.ensureLevel(level)

	suffixWord := key[level:]
	//FIXME: impl
	b.suffixes[level] = append(b.suffixes[level], suffixWord)
	b.suffixCounts[level]++
}

func (b *Builder) isSameKey(a, c []byte) bool {
	return bytes.Equal(a, c)
}

// skipCommonPrefix skips common prefix, returns level that different char.
func (b *Builder) skipCommonPrefix(key []byte) (level int) {
	for level < len(key) && b.isCommonPrefix(key[level], level) {
		setBit(b.lsHasChild[level], b.numNodes(level)-1)
		level++
	}
	return level
}

// isCommonPrefix returns whether char is common prefix.
func (b *Builder) isCommonPrefix(c byte, level int) bool {
	return level < b.treeHeight() &&
		!b.isLastItemTerminator[level] &&
		// because all keys is sorted, so just check last label
		c == b.lsLabels[level][len(b.lsLabels[level])-1]
}

func setBit(bs []uint64, pos int) {
	wordOff := pos / bitsSize
	bitsOff := pos % bitsSize
	bs[wordOff] |= uint64(1) << bitsOff
}

func (b *Builder) numNodes(level int) int {
	return len(b.lsLabels[level])
}

func (b *Builder) ensureLevel(level int) {
	if level >= b.treeHeight() {
		b.addLevel()
	}
}

func (b *Builder) treeHeight() int {
	return len(b.nodeCounts)
}

func (b *Builder) addLevel() {
	// cached
	b.lsLabels = append(b.lsLabels, []byte{})
	b.lsHasChild = append(b.lsHasChild, []uint64{})
	b.lsLouds = append(b.lsLouds, []uint64{})
	// b.hasPrefix = append(b.hasPrefix, b.pickUint64Slice())
	b.hasSuffix = append(b.hasSuffix, []uint64{})

	// b.values = append(b.values, []byte{})
	// b.valueCounts = append(b.valueCounts, 0)
	// b.prefixes = append(b.prefixes, [][]byte{})
	b.suffixes = append(b.suffixes, [][]byte{})
	b.suffixCounts = append(b.suffixCounts, 0)

	b.nodeCounts = append(b.nodeCounts, 0)
	b.isLastItemTerminator = append(b.isLastItemTerminator, false)

	level := b.treeHeight() - 1
	b.lsHasChild[level] = append(b.lsHasChild[level], 0)
	b.lsLouds[level] = append(b.lsLouds[level], 0)

	// b.hasPrefix[level] = append(b.hasPrefix[level], 0)
	// b.hasSuffix[level] = append(b.hasSuffix[level], 0)
}

func (b *Builder) getSparseStartLevel() int {
	return b.sparseStartLevel
}

func (b *Builder) getLabels() [][]byte {
	return b.lsLabels
}

func (b *Builder) getLoudsBits() [][]uint64 {
	return b.lsLouds
}

func (b *Builder) getHasChildBits() [][]uint64 {
	return b.lsHasChild
}

func (b *Builder) getSuffixLen() int {
	//FIXME: impl
	return 0
}

func (b *Builder) isStartOfNode(level, pos int) bool {
	return readBit(b.lsLouds[level], pos)
}

func (b *Builder) isTerminator(level, pos int) bool {
	label := b.lsLabels[level][pos]
	return label == terminator && !readBit(b.lsHasChild[level], pos)
}
