package surf

import "bytes"

type Iterator struct {
	loudsSparseIt *loudsSparseIterator
}

func NewIterator(trie *Trie) *Iterator {
	it := &Iterator{
		loudsSparseIt: newLoudsSparseIterator(trie.loudsSparse),
	}
	return it
}

func (it *Iterator) First() {
	it.loudsSparseIt.moveToLeftMostKey()
}

func (it *Iterator) Next() {
	it.loudsSparseIt.next()
}

func (it *Iterator) IsValid() bool {
	return it.loudsSparseIt.IsValid()
}

func (it *Iterator) Key() []byte {
	return it.loudsSparseIt.getKey()
}

func (it *Iterator) Value() uint32 {
	return it.loudsSparseIt.getValue()
}

func (it *Iterator) Seek(prefix []byte) {
	_ = it.loudsSparseIt.seek(prefix)
}

type loudsSparseIterator struct {
	trie *loudsSparse
	// true means the iter currently points to a valid key
	isValid bool
	// start couting from start level; does NOT include suffix
	level          int // level
	key            []byte
	fullKey        []byte
	isAtTerminator bool
	posInTrie      []int
}

func newLoudsSparseIterator(trie *loudsSparse) *loudsSparseIterator {
	it := &loudsSparseIterator{
		trie:           trie,
		isValid:        false,
		level:          0,
		isAtTerminator: false,
	}
	it.key = make([]byte, trie.height)
	it.posInTrie = make([]int, trie.height)
	return it
}

func (it *loudsSparseIterator) reset() {
	it.isValid = false
	it.level = 0
	it.isAtTerminator = false

	for level := 0; level < len(it.key); level++ {
		it.key[level] = 0
		it.posInTrie[level] = 0
	}
}

func (it *loudsSparseIterator) moveToLeftMostKey() {
	if it.level == 0 {
		pos := it.trie.getFirstLabelPos(0)
		label := it.trie.labels.labels[pos]
		it.append(label, pos)
	}
	level := it.level - 1
	pos := it.posInTrie[level]
	if !it.trie.hasChild.ReadBit(pos) {
		label := it.trie.labels.labels[pos]
		if label == terminator && !it.trie.isEndOfNode(pos) {
			it.isAtTerminator = true
		}
		it.isValid = true
		return
	}

	for level < it.trie.height {
		// process child
		nodeNum := it.trie.getChildNodeNum(pos)
		pos = it.trie.getFirstLabelPos(nodeNum)
		label := it.trie.labels.labels[pos]

		// if trie branch terminates
		if !it.trie.hasChild.ReadBit(int(pos)) {
			it.append(label, pos)
			if label == terminator && !it.trie.isEndOfNode(int(pos)) {
				it.isAtTerminator = true
			}
			it.isValid = true
			return
		}

		it.append(label, pos)
		level++
	}
}

func (it *loudsSparseIterator) next() {
	if !it.isValid {
		return
	}
	it.doNext()
}

func (it *loudsSparseIterator) doNext() {
	pos, ok := it.nextPos()
	if !ok {
		return
	}
	// read next label
	label := it.trie.labels.labels[pos]
	it.append(label, pos)
	// read more lable
	it.moveToLeftMostKey()
}

func (it *loudsSparseIterator) nextPos() (int, bool) {
	it.isAtTerminator = false
	level := it.level - 1 // because append label add level, so need -1
	pos := it.posInTrie[level] + 1

	for pos >= it.trie.louds.numBits || it.trie.louds.ReadBit(pos) {
		// if not child, try find next node of parent
		if level == 0 {
			it.level = 0
			it.isValid = false
			return -1, false
		}
		level--
		pos = it.posInTrie[level] + 1 // brother node
	}
	it.level = level
	return pos, true
}

func (it *loudsSparseIterator) seek(prefix []byte) bool {
	it.reset()

	nodeNum := 0
	var ok bool
	pos := it.trie.getFirstLabelPos(nodeNum)

	level := 0
	for ; level < len(prefix); level++ {
		nodeSize := it.trie.nodeSize(pos)
		if pos, ok = it.trie.labels.Search(prefix[level], pos, nodeSize); !ok {
			// if no exact match
			it.moveToLeftInNextSubTrie(pos, nodeSize, prefix[level])
			return false
		}
		it.append(prefix[level], pos)

		// if trie brach terminates
		if !it.trie.hasChild.ReadBit(pos) {
			// check suffix
			return it.compareSuffixGreaterThan(prefix, pos, level+1)
		}

		// move to child
		nodeNum = it.trie.getChildNodeNum(pos)
		pos = it.trie.getFirstLabelPos(nodeNum)
	}
	if it.trie.labels.labels[pos] == terminator &&
		!it.trie.hasChild.ReadBit(pos) &&
		!it.trie.isEndOfNode(pos) {
		// prefix is key
		it.append(terminator, pos)
		it.isAtTerminator = true
		it.isValid = true
		return false
	}

	if len(prefix) <= level {
		// try read more label
		it.moveToLeftMostKey()
		return false
	}

	it.isValid = true
	return true
}

func (it *loudsSparseIterator) compareSuffixGreaterThan(prefix []byte, pos, level int) bool {
	if level < len(prefix) {
		// if prefix is remaining, need check suffix
		suffix := it.trie.suffixes.GetSuffix(pos)
		if !bytes.HasSuffix(suffix, prefix[level:]) {
			// suffix not match, do next
			it.doNext()
			return false
		}
	}
	it.isValid = true
	return true
}

func (it *loudsSparseIterator) moveToLeftInNextSubTrie(pos, nodeSize int, label byte) {
	// if no label is greater than key[level] in this node
	var ok bool
	if pos, ok = it.trie.labels.SearchGreaterThan(label, pos, nodeSize); !ok {
		it.appendByPos(pos + nodeSize - 1)

		if _, ok = it.nextPos(); !ok {
			return
		}

		it.moveToLeftMostKey()
	} else {
		it.appendByPos(pos)
		it.moveToLeftMostKey()
	}
}

func (it *loudsSparseIterator) IsValid() bool {
	return it.isValid
}

func (it *loudsSparseIterator) getKey() []byte {
	kLen := it.level
	if it.isAtTerminator {
		kLen--
	}
	uniqueKey := it.key[:kLen]

	pos := it.posInTrie[it.level-1]
	// check whether has suffix key
	suffix := it.trie.suffixes.GetSuffix(pos)
	if len(suffix) == 0 {
		return uniqueKey
	}
	// result key = unique key + suffix key
	expectLen := kLen + len(suffix)
	if cap(it.fullKey) < expectLen {
		it.fullKey = make([]byte, expectLen)
	}
	it.fullKey = it.fullKey[0:expectLen]
	copy(it.fullKey[:kLen], uniqueKey)
	copy(it.fullKey[kLen:], suffix)
	return it.fullKey
}

func (it *loudsSparseIterator) getValue() uint32 {
	valPos := it.trie.valuePos(it.posInTrie[it.level-1])
	return it.trie.values.Get(valPos)
}

// append appends label to key buffer, and goto next level.
func (it *loudsSparseIterator) append(label byte, pos int) {
	it.key[it.level] = label
	it.posInTrie[it.level] = pos
	it.level++
}

func (it *loudsSparseIterator) appendByPos(pos int) {
	it.append(it.trie.labels.labels[pos], pos)
}
