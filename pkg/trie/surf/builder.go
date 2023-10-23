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
	one        = uint64(1)
)

type Level struct {
	labels      []byte
	hasChild    []uint64
	louds       []uint64
	hasSuffixes []uint64
	suffixes    [][]byte
	values      []uint32
	item        int
}

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

	// isLastItemTerminator []bool

	hasSuffix [][]uint64
	suffixes  [][][]byte

	values    [][]uint32
	nodeItems []int

	height    int
	totalKeys int

	// pooling data-structures
	cachedValue   [][]uint32
	cachedLabel   [][]byte
	cachedUint64s [][]uint64
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) SetLevel(maxLevel int) {
	// for level := 0; level < maxLevel; level++ {
	// 	b.addLevel()
	// }
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
	numBytes := labelsSize(b.nodeItems)
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
		if len(values) > 0 {
			if _, err := w.Write(u32SliceToBytes(values)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Builder) Reset() {
	b.height = 0
	b.totalKeys = 0
	// cache lsLabels
	for idx := range b.lsLabels {
		b.cachedLabel = append(b.cachedLabel, b.lsLabels[idx][:0])
	}
	b.lsLabels = b.lsLabels[:0]

	// cache lsHasChild/lsLouds/suffixes
	for idx := range b.lsHasChild {
		b.cachedUint64s = append(b.cachedUint64s, b.lsHasChild[idx][:0])
		b.cachedUint64s = append(b.cachedUint64s, b.lsLouds[idx][:0])
		b.cachedUint64s = append(b.cachedUint64s, b.hasSuffix[idx][:0])
	}
	b.lsHasChild = b.lsHasChild[:0]
	b.lsLouds = b.lsLouds[:0]
	b.hasSuffix = b.hasSuffix[:0]
	b.suffixes = b.suffixes[:0]

	for idx := range b.values {
		b.cachedValue = append(b.cachedValue, b.values[idx][:0])
	}
	// reset values
	b.values = b.values[:0]

	// b.isLastItemTerminator = b.isLastItemTerminator[:0]
	b.nodeItems = b.nodeItems[:0]
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
		for b.height <= len(key)+1 {
			b.addLevel()
		}
		// needCheckLevel := len(key)+1 >= b.height
		level := 0
		kenLen := len(key)
		// skip common prefix
		for level < kenLen && level < len(previousKey) && key[level] == previousKey[level] {
			// setBit(b.lsHasChild[level], b.nodeCounts[level]-1)
			level++
		}
		// level := b.skipCommonPrefix(keys[i])
		if i < len(keys)-1 {
			level = b.insertKeyBytesToTrieUntilUnique(key, keys[i+1], level, kenLen)
		} else {
			// for last key, there is no successor key in the list
			level = b.insertKeyBytesToTrieUntilUnique(key, emptyKey, level, kenLen)
		}
		// if needCheckLevel {
		// 	b.ensureLevel(level)
		// }
		b.insertValue(values[i], level)
		b.totalKeys++
		if level < kenLen {
			// insert suffix if has suffix
			b.insertSuffix(key, level)
		}
		previousKey = key
	}
}

func (b *Builder) insertKeyBytesToTrieUntilUnique(key []byte, nextKey []byte, startLevel, keyLen int) (level int) {
	level = startLevel
	isStartOfNode := false
	isTerm := false
	// level should be at most equal to tree height
	// if needCheckLevel {
	// 	b.ensureLevel(level)
	// }

	// if it is the start of level, the louds bit needs to be set
	if b.nodeItems[level] == 0 {
		isStartOfNode = true
	}
	// after skipping the common prefix, the first following byte
	// shoud be in an the node as the previous key.
	pos := b.insertKeyByte(key[level], level, keyLen, isStartOfNode, isTerm)
	level++ // goto next
	if level > len(nextKey) || key[level-1] != nextKey[level-1] || !b.isSameKey(key[0:level], nextKey[0:level]) {
		return level
	}
	if level-1 < keyLen {
		setBit(b.lsHasChild[level-1], pos)
	}

	// all the follwing bytes inserted must be the start of a new node,
	// becase generate new sub trie
	isStartOfNode = true
	for level < len(key) && level < len(nextKey) && key[level] == nextKey[level] {
		// level should be at most equal to tree height
		// if needCheckLevel {
		// 	b.ensureLevel(level)
		// }

		pos := b.insertKeyByte(key[level], level, keyLen, isStartOfNode, isTerm)
		if level < keyLen {
			setBit(b.lsHasChild[level], pos)
		}
		level++
	}

	// the last byte inserted makes key unique in the tire
	if level < len(key) {
		// level should be at most equal to tree height
		// if needCheckLevel {
		// 	b.ensureLevel(level)
		// }
		//
		b.insertKeyByte(key[level], level, keyLen, isStartOfNode, isTerm)
	} else {
		// insert terminator char
		isTerm = true
		// level should be at most equal to tree height
		// if needCheckLevel {
		// 	b.ensureLevel(level)
		// }

		b.insertKeyByte(terminator, level, keyLen, isStartOfNode, isTerm)
	}
	level++ // goto next, for storing value

	return level
}

func (b *Builder) insertKeyByte(key byte, level, keyLen int, isStartOfNode, isTerm bool) int {
	// store parent has child
	// sets parent node's child indicator
	// if level > 0 {
	// 	parent := level - 1
	// 	pos := b.nodeItems[parent] - 1
	// 	if !readBit(b.lsHasChild[parent], pos) {
	// 		// all keys is sorted, so new key will append right
	// 		setBit(b.lsHasChild[parent], pos)
	// 	}
	// 	// ok := b.child[parent][pos]
	// 	// if !ok {
	// 	// 	b.child[parent][pos] = true
	// 	// }
	// }
	// store label
	b.lsLabels[level] = append(b.lsLabels[level], key)
	b.nodeItems[level]++

	pos := b.nodeItems[level] - 1
	// if level < keyLen {
	// 	setBit(b.lsHasChild[level], b.nodeItems[level]-1)
	// }

	// store louds
	if isStartOfNode {
		setBit(b.lsLouds[level], pos)
	}
	// b.isLastItemTerminator[level] = isTerm

	b.moveToNextNodeSlot(level, pos)
	return pos
}

func (b *Builder) moveToNextNodeSlot(level, pos int) {
	if wordsIndex(uint(pos)) == 0 {
		// put next slot for bit vector context
		b.lsHasChild[level] = append(b.lsHasChild[level], 0)
		b.lsLouds[level] = append(b.lsLouds[level], 0)
		b.hasSuffix[level] = append(b.hasSuffix[level], 0)
	}
}

func wordsNeeded(d []uint64, i int) int {
	if i > (cap(d) - bitsSize + 1) {
		return int(cap(d) >> 6)
	}
	return int((i + (bitsSize - 1)) >> 6)
}

// isLevelEmpty returns whether level is empty.
func (b *Builder) isLevelEmpty(level int) bool {
	return level >= b.height || b.nodeItems[level] == 0
}

func (b *Builder) insertSuffix(key []byte, level int) {
	suffixLevel := level - 1 // need -1, because after insert label, level will move next
	setBit(b.hasSuffix[suffixLevel], b.nodeItems[suffixLevel]-1)
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
		setBit(b.lsHasChild[level], b.nodeItems[level]-1)
		level++
	}
	return level
}

// isCommonPrefix returns whether char is common prefix.
func (b *Builder) isCommonPrefix(c byte, level int) bool {
	return level < b.height &&
		// because all keys is sorted, so just check last label
		c == b.lsLabels[level][len(b.lsLabels[level])-1]
	// !b.isLastItemTerminator[level]
}

func setBit(bs []uint64, pos int) {
	// wordOff := pos / bitsSize
	// bitsOff := pos % bitsSize
	bs[pos>>6] |= one << wordsIndex(uint(pos))
}

func wordsIndex(i uint) uint {
	return i & (bitsSize - 1)
}

// func (b *Builder) numNodes(level int) int {
// 	return len(b.lsLabels[level])
// }

// func (b *Builder) ensureLevel(level int) {
// 	// level should be at most equal to trie height
// 	// if level >= b.height {
// 	// 	b.addLevel()
// 	// }
// }

func (b *Builder) addLevel() {
	b.height++
	// cached
	b.lsLabels = append(b.lsLabels, b.pickLabels())
	b.lsHasChild = append(b.lsHasChild, b.pickUint64Slice())
	b.lsLouds = append(b.lsLouds, b.pickUint64Slice())
	b.hasSuffix = append(b.hasSuffix, b.pickUint64Slice())

	b.values = append(b.values, b.pickValues())
	b.suffixes = append(b.suffixes, [][]byte{})
	// b.isLastItemTerminator = append(b.isLastItemTerminator, false)
	b.nodeItems = append(b.nodeItems, 0)

	level := b.height - 1
	b.moveToNextNodeSlot(level, 0)
	// b.lsHasChild[level] = append(b.lsHasChild[level], 0)
	// b.lsLouds[level] = append(b.lsLouds[level], 0)
	// b.hasSuffix[level] = append(b.hasSuffix[level], 0)
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

func (b *Builder) pickValues() []uint32 {
	if len(b.cachedLabel) == 0 {
		return []uint32{}
	}
	tailIndex := len(b.cachedValue) - 1
	ptr := b.cachedValue[tailIndex]
	b.cachedValue = b.cachedValue[:tailIndex]
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
