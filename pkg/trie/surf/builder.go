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

	height    int
	totalKeys int

	// pooling data-structures
	cachedLabel   [][]byte
	cachedUint64s [][]uint64
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Write(w io.Writer) error {
	height := len(b.lsLabels)
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
	hasSuffixes := &BitVector{}
	hasSuffixes.Init(b.hasSuffix, numNodesPerLevel)
	if err := hasSuffixes.write(w); err != nil {
		return err
	}
	suffixes := &SuffixVector{}
	suffixes.initData(numNodesPerLevel, b.suffixes)
	if err := suffixes.write(w); err != nil {
		return err
	}
	// write values
	for level := range b.values {
		values := b.values[level]
		for _, val := range values {
			binary.LittleEndian.PutUint32(bs[:], uint32(val))
			if _, err := w.Write(bs[:]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Builder) Reset() {
	b.totalKeys = 0
	// cache lsLabels
	for idx := range b.lsLabels {
		b.cachedLabel = append(b.cachedLabel, b.lsLabels[idx][:0])
	}
	b.lsLabels = b.lsLabels[:0]

	// cache lsHasChild
	for idx := range b.lsHasChild {
		b.cachedUint64s = append(b.cachedUint64s, b.lsHasChild[idx][:0])
	}
	b.lsHasChild = b.lsHasChild[:0]

	// cache lsLoudsBits
	for idx := range b.lsLouds {
		b.cachedUint64s = append(b.cachedUint64s, b.lsLouds[idx][:0])
	}
	b.lsLouds = b.lsLouds[:0]

	// reset values
	b.values = b.values[:0]

	// cache has suffix
	for idx := range b.hasSuffix {
		b.hasSuffix = append(b.hasSuffix, b.hasSuffix[idx][:0])
	}
	// reset suffixes
	b.hasSuffix = b.hasSuffix[:0]
	b.suffixes = b.suffixes[:0]

	// reset nodeCounts
	b.isLastItemTerminator = b.isLastItemTerminator[:0]

	// reset suffixes
	b.hasSuffix = b.hasSuffix[:0]
	b.suffixes = b.suffixes[:0]
	b.isLastItemTerminator = b.isLastItemTerminator[:0]
}

func (b *Builder) Build(keys [][]byte, values []uint32) {
	b.buildSparse(keys, values)
}

func (b *Builder) buildSparse(keys [][]byte, values []uint32) {
	var previousKey []byte
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		if len(key) == 0 {
			// ignore empty key
			continue
		}
		// skip common prefix
		level := 0
		if previousKey != nil {
			for level < len(key) && level < len(previousKey) && key[level] == previousKey[level] {
				setBit(b.lsHasChild[level], b.numNodes(level)-1)
				level++
			}
		}
		// level := b.skipCommonPrefix(keys[i])
		if i < len(keys)-1 {
			level = b.insertKeyBytesToTrieUntilUnique(key, keys[i+1], level)
		} else {
			// for last key, there is no successor key in the list
			level = b.insertKeyBytesToTrieUntilUnique(key, emptyKey, level)
		}
		b.ensureLevel(level)
		b.insertValue(values[i], level)
		b.totalKeys++
		if level < len(key) {
			// insert suffix if has suffix
			b.insertSuffix(key, level)
		}
		previousKey = key
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
	return level >= len(b.lsLabels) || len(b.lsLabels[level]) == 0
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
	return level < b.height &&
		// because all keys is sorted, so just check last label
		c == b.lsLabels[level][len(b.lsLabels[level])-1] &&
		!b.isLastItemTerminator[level]
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
	if level >= len(b.lsLabels) {
		b.addLevel()
	}
}

func (b *Builder) addLevel() {
	// cached
	b.lsLabels = append(b.lsLabels, b.pickLabels())
	b.lsHasChild = append(b.lsHasChild, b.pickUint64Slice())
	b.lsLouds = append(b.lsLouds, b.pickUint64Slice())
	b.hasSuffix = append(b.hasSuffix, b.pickUint64Slice())

	b.values = append(b.values, []uint32{})
	b.suffixes = append(b.suffixes, [][]byte{})
	b.isLastItemTerminator = append(b.isLastItemTerminator, false)

	level := len(b.lsLabels) - 1
	b.lsHasChild[level] = append(b.lsHasChild[level], 0)
	b.lsLouds[level] = append(b.lsLouds[level], 0)
	b.hasSuffix[level] = append(b.hasSuffix[level], 0)
}

func (b *Builder) isStartOfNode(level, pos int) bool {
	return readBit(b.lsLouds[level], pos)
}

func (b *Builder) isTerminator(level, pos int) bool {
	label := b.lsLabels[level][pos]
	return label == terminator && !readBit(b.lsHasChild[level], pos)
}

func (b *Builder) pickLabels() []byte {
	if len(b.cachedLabel) == 0 {
		return []byte{}
	}
	tailIndex := len(b.cachedLabel) - 1
	ptr := b.cachedLabel[tailIndex]
	b.cachedLabel = b.cachedLabel[:tailIndex]
	return ptr
}

func (b *Builder) pickUint64Slice() []uint64 {
	if len(b.cachedUint64s) == 0 {
		return []uint64{}
	}
	tailIndex := len(b.cachedUint64s) - 1
	ptr := b.cachedUint64s[tailIndex]
	b.cachedUint64s = b.cachedUint64s[:tailIndex]
	return ptr
}
