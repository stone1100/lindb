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
	// LOUDS-Sparse tire level
	levels []*Level

	height    int
	totalKeys int

	// pooling data-structures
	cachedLevel []*Level
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Write(w io.Writer) error {
	var (
		bs [4]byte
	)
	// write total keys
	binary.LittleEndian.PutUint32(bs[:], uint32(b.totalKeys))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	// write height
	binary.LittleEndian.PutUint32(bs[:], uint32(b.height))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	// write labels
	numBytes := labelsSize(b.levels)
	binary.LittleEndian.PutUint32(bs[:], uint32(numBytes))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	for level := range b.levels {
		levelObj := b.levels[level]
		if _, err := w.Write(levelObj.lsLabels); err != nil {
			return err
		}
	}
	// write has child
	child := &BitVectorRank{}
	child.Init2(rankSparseBlockSize, b.levels, HasChild)
	if err := child.write(w); err != nil {
		return err
	}
	// write louds
	louds := &BitVectorSelect{}
	louds.Init(b.levels, Louds)
	if err := louds.write(w); err != nil {
		return err
	}
	// write suffixes
	hasSuffix := &BitVectorRank{}
	hasSuffix.Init2(rankSparseBlockSize, b.levels, HasSuffix)
	if err := hasSuffix.write(w); err != nil {
		return err
	}
	suffixes := &SuffixVector{}
	suffixes.initData(b.levels)
	if err := suffixes.write(w); err != nil {
		return err
	}

	// write values
	for level := range b.levels {
		values := b.levels[level].values
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

	// cache level
	for idx := range b.levels {
		level := b.levels[idx]
		level.Reset()
		b.cachedLevel = append(b.cachedLevel, level)
	}
	b.levels = b.levels[:0]
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
		// level should be at most equal to tree height
		for b.height <= len(key)+1 { // + terminate node
			b.addLevel()
		}
		level := 0
		kenLen := len(key)
		// skip common prefix
		for level < kenLen && level < len(previousKey) && key[level] == previousKey[level] {
			level++
		}
		if i < len(keys)-1 {
			level = b.insertKeyBytesToTrieUntilUnique(key, keys[i+1], level, kenLen)
		} else {
			// for last key, there is no successor key in the list
			level = b.insertKeyBytesToTrieUntilUnique(key, emptyKey, level, kenLen)
		}
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
	// if needCheckLevel {
	// 	b.ensureLevel(level)
	// }

	// if it is the start of level, the louds bit needs to be set
	if b.levels[level].item == 0 {
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
		setBit(b.levels[level-1].lsHasChild, pos)
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
			setBit(b.levels[level].lsHasChild, pos)
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
	levelObj := b.levels[level]
	levelObj.lsLabels = append(levelObj.lsLabels, key)
	levelObj.item++

	pos := levelObj.item - 1
	// if level < keyLen {
	// 	setBit(b.lsHasChild[level], b.nodeItems[level]-1)
	// }

	// store louds
	if isStartOfNode {
		setBit(levelObj.lsLouds, pos)
	}
	// b.isLastItemTerminator[level] = isTerm

	b.moveToNextNodeSlot(level, pos)
	return pos
}

func (b *Builder) moveToNextNodeSlot(level, pos int) {
	if wordsIndex(uint(pos)) == 0 {
		levelObj := b.levels[level]
		levelObj.lsHasChild = append(levelObj.lsHasChild, 0)
		levelObj.lsLouds = append(levelObj.lsLouds, 0)
		levelObj.hasSuffixes = append(levelObj.hasSuffixes, 0)
		// put next slot for bit vector context
		// b.lsHasChild[level] = append(b.lsHasChild[level], 0)
		// b.lsLouds[level] = append(b.lsLouds[level], 0)
		// b.hasSuffix[level] = append(b.hasSuffix[level], 0)
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
	return level >= b.height || b.levels[level].item == 0
}

func (b *Builder) insertSuffix(key []byte, level int) {
	suffixLevel := level - 1 // need -1, because after insert label, level will move next
	levelObj := b.levels[suffixLevel]
	setBit(levelObj.hasSuffixes, levelObj.item-1)
	levelObj.suffixes = append(levelObj.suffixes, key[level:])
}

func (b *Builder) insertValue(value uint32, level int) {
	levelObj := b.levels[level]
	levelObj.values = append(levelObj.values, value)
}

func (b *Builder) isSameKey(a, c []byte) bool {
	return bytes.Equal(a, c)
}

// skipCommonPrefix skips common prefix, returns level that different char.
// func (b *Builder) skipCommonPrefix(key []byte) (level int) {
// 	for level < len(key) && b.isCommonPrefix(key[level], level) {
// 		setBit(b.lsHasChild[level], b.nodeItems[level]-1)
// 		level++
// 	}
// 	return level
// }

// isCommonPrefix returns whether char is common prefix.
// func (b *Builder) isCommonPrefix(c byte, level int) bool {
// 	return level < b.height &&
// 		// because all keys is sorted, so just check last label
// 		c == b.lsLabels[level][len(b.lsLabels[level])-1]
// 	// !b.isLastItemTerminator[level]
// }

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
	b.levels = append(b.levels, b.pickUint64Slice())
	// cached
	// b.lsLabels = append(b.lsLabels, b.pickLabels())
	// b.lsHasChild = append(b.lsHasChild, b.pickUint64Slice())
	// b.lsLouds = append(b.lsLouds, b.pickUint64Slice())
	// b.hasSuffix = append(b.hasSuffix, b.pickUint64Slice())
	//
	// b.values = append(b.values, b.pickValues())
	// b.suffixes = append(b.suffixes, [][]byte{})
	// // b.isLastItemTerminator = append(b.isLastItemTerminator, false)
	// b.nodeItems = append(b.nodeItems, 0)

	level := b.height - 1
	b.moveToNextNodeSlot(level, 0)
	// b.lsHasChild[level] = append(b.lsHasChild[level], 0)
	// b.lsLouds[level] = append(b.lsLouds[level], 0)
	// b.hasSuffix[level] = append(b.hasSuffix[level], 0)
}

// func (b *Builder) isStartOfNode(level, pos int) bool {
// 	return readBit(b.lsLouds[level], pos)
// }

// func (b *Builder) isTerminator(level, pos int) bool {
// 	label := b.levels[level].lsLabels[pos]
// 	return label == terminator && !readBit(b.levels[level].lsHasChild, pos)
// }

//	func (b *Builder) pickLabels() []byte {
//		if len(b.cachedLabel) == 0 {
//			return []byte{}
//		}
//		tailIndex := len(b.cachedLabel) - 1
//		ptr := b.cachedLabel[tailIndex]
//		b.cachedLabel = b.cachedLabel[:tailIndex]
//		return ptr
//	}
//
//	func (b *Builder) pickValues() []uint32 {
//		if len(b.cachedLabel) == 0 {
//			return []uint32{}
//		}
//		tailIndex := len(b.cachedValue) - 1
//		ptr := b.cachedValue[tailIndex]
//		b.cachedValue = b.cachedValue[:tailIndex]
//		return ptr
//	}
func (b *Builder) pickUint64Slice() *Level {
	if len(b.cachedLevel) == 0 {
		return NewLevel()
	}
	tailIndex := len(b.cachedLevel) - 1
	ptr := b.cachedLevel[tailIndex]
	b.cachedLevel = b.cachedLevel[:tailIndex]
	return ptr
}
