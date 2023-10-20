package surf

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	bitsSize = 64
	// terminator($) indicates the situation where a prefix strings
	// leading to a node is also a valid key
	terminator = 0xff
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
	// LOUDS-Sparse context: labels/hasChild/louds
	//
	// store all the branching labels for each trie node
	lsLabels [][]byte
	// one bit for each byte in labels to indicate whether
	// a child branch continues(i.e. points to a sub-trie)
	// or terminals(i.e. points to a value)
	lsHasChild [][]uint64
	// one bit for each byte in labels to indicate if a lable
	// is the first node in trie
	lsLouds [][]uint64

	isLastItemTerminator []bool

	hasSuffix [][]uint64
	suffixes  [][][]byte

	values [][]uint32

	totalKeys int
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Write(w io.Writer) error {
	height := b.treeHeight()
	var (
		bs [4]byte
	)
	// write total keys
	binary.LittleEndian.PutUint32(bs[:], uint32(b.totalKeys))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	// write height
	binary.LittleEndian.PutUint32(bs[:], uint32(height))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	// write labels
	numBytes := labelsSize(b.lsLabels)
	binary.LittleEndian.PutUint32(bs[:], uint32(numBytes))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	for level := range b.lsLabels {
		if _, err := w.Write(b.lsLabels[level]); err != nil {
			return err
		}
	}
	numNodesPerLevel := make([]int, height)
	for level := range numNodesPerLevel {
		numNodesPerLevel[level] = len(b.lsLabels[level])
	}
	// write has child
	hasChild := &BitVector{}
	hasChild.Init(b.lsHasChild, numNodesPerLevel)
	if err := hasChild.write(w); err != nil {
		return err
	}
	// write louds
	louds := &BitVector{}
	louds.Init(b.lsLouds, numNodesPerLevel)
	if err := louds.write(w); err != nil {
		return err
	}
	// write suffixes
	suffixes := &SuffixVector{}
	suffixes.Init(b.hasSuffix, numNodesPerLevel, b.suffixes)
	if err := suffixes.write(w); err != nil {
		return err
	}
	// write suffixes
	values := &ValueVector{}
	values.Init(b.values)
	if err := values.write(w); err != nil {
		return err
	}
	return nil
}

func (b *Builder) Build(keys [][]byte, values []uint32) {
	b.totalKeys = len(keys)
	b.buildSparse(keys, values)
}

func (b *Builder) buildSparse(keys [][]byte, values []uint32) {
	for i := 0; i < len(keys); i++ {
		// skip common prefix
		level := b.skipCommonPrefix(keys[i])
		if i < len(keys)-1 {
			level = b.insertKeyBytesToTrieUntilUnique(keys[i], keys[i+1], level)
		} else {
			// for last key, there is no successor key in the list
			level = b.insertKeyBytesToTrieUntilUnique(keys[i], emptyKey, level)
		}
		b.ensureLevel(level)
		b.insertValue(values[i], level)
		if level < len(keys[i]) {
			// insert suffix if has suffix
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
	} else {
		// insert terminator char
		isTerm = true
		b.insertKeyByte(terminator, level, isStartOfNode, isTerm)
	}
	level++ // goto next, for storing value

	return level
}

func (b *Builder) insertKeyByte(key byte, level int, isStartOfNode, isTerm bool) {
	// level should be at most equal to tree height
	b.ensureLevel(level)

	// store parent has child
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
	suffixLevel := level - 1 // need -1, because after insert label, level will move next
	setBit(b.hasSuffix[suffixLevel], b.numNodes(suffixLevel)-1)
	b.suffixes[suffixLevel] = append(b.suffixes[suffixLevel], key[level:])
}

func (b *Builder) insertValue(value uint32, level int) {
	b.values[level] = append(b.values[level], value)
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
	// level should be at most equal to trie height
	if level >= b.treeHeight() {
		b.addLevel()
	}
}

func (b *Builder) treeHeight() int {
	return len(b.lsLabels)
}

func (b *Builder) addLevel() {
	// cached
	b.lsLabels = append(b.lsLabels, []byte{})
	b.lsHasChild = append(b.lsHasChild, []uint64{})
	b.lsLouds = append(b.lsLouds, []uint64{})
	// b.hasPrefix = append(b.hasPrefix, b.pickUint64Slice())
	b.hasSuffix = append(b.hasSuffix, []uint64{})
	b.values = append(b.values, []uint32{})

	// b.values = append(b.values, []byte{})
	// b.valueCounts = append(b.valueCounts, 0)
	// b.prefixes = append(b.prefixes, [][]byte{})
	b.suffixes = append(b.suffixes, [][]byte{})

	b.isLastItemTerminator = append(b.isLastItemTerminator, false)

	level := b.treeHeight() - 1
	b.lsHasChild[level] = append(b.lsHasChild[level], 0)
	b.lsLouds[level] = append(b.lsLouds[level], 0)
	// b.hasPrefix[level] = append(b.hasPrefix[level], 0)
	b.hasSuffix[level] = append(b.hasSuffix[level], 0)
}

func (b *Builder) isStartOfNode(level, pos int) bool {
	return readBit(b.lsLouds[level], pos)
}

func (b *Builder) isTerminator(level, pos int) bool {
	label := b.lsLabels[level][pos]
	return label == terminator && !readBit(b.lsHasChild[level], pos)
}
