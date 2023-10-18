package surf

import "fmt"

type Iterator struct {
	loudsDenseIt  *loudsDenseIterator
	loudsSparseIt *loudsSparseIterator
}

func NewIterator(trie *Trie) *Iterator {
	it := &Iterator{
		loudsSparseIt: newLoudsSparseIterator(trie.loudsSparse),
	}
	loudsDense := trie.loudsDense
	if loudsDense != nil {
		it.loudsDenseIt = &loudsDenseIterator{
			loudsDense: loudsDense,
		}
	}
	return it
}

func (it *Iterator) First() {
	//TODO: dense

	// it.loudsSparseIt.setToFirstLabelInRoot()
	it.loudsSparseIt.moveToLeftMostKey()
}

func (it *Iterator) Next() {
	it.loudsSparseIt.next()
}

func (it *Iterator) IsValid() bool {
	return it.loudsSparseIt.IsValid()
}

func (it *Iterator) Key() []byte {
	key := it.loudsSparseIt.getKey()
	return key
}

func (it *Iterator) passToSparse() {
	//TODO:
}

func (it *Iterator) Seek(prefix []byte) {
	ok := it.loudsSparseIt.seek(prefix)
	fmt.Println(ok)
}

type loudsDenseIterator struct {
	loudsDense *loudsDense
}

type loudsSparseIterator struct {
	trie *loudsSparse
	// true means the iter currently points to a valid key
	isValid    bool
	startLevel int
	// passed in by the dense iterator, default 0
	startNodeNum int
	// start couting from start level; does NOT include suffix
	level          int // level
	key            []byte
	isAtTerminator bool
	posInTrie      []int
}

func newLoudsSparseIterator(trie *loudsSparse) *loudsSparseIterator {
	it := &loudsSparseIterator{
		trie:           trie,
		isValid:        false,
		startNodeNum:   0,
		level:          0,
		isAtTerminator: false,
		startLevel:     trie.startLevel,
	}
	size := trie.height - trie.startLevel
	it.key = make([]byte, size)
	it.posInTrie = make([]int, size)
	return it
}

func (it *loudsSparseIterator) reset() {
	it.isValid = false
	it.startNodeNum = 0
	it.level = 0
	it.isAtTerminator = false

	for level := 0; level < len(it.key); level++ {
		it.key[level] = 0
		it.posInTrie[level] = 0
	}
}

func (it *loudsSparseIterator) setToFirstLabelInRoot() {
	// it.posInTrie[0] = 0
	// it.key = append(it.key, it.trie.labels.Read(0))
	// it.key[0] = it.trie.labels.Read(0)
}

func (it *loudsSparseIterator) moveToLeftMostKey() {
	if it.level == 0 {
		pos := it.trie.getFirstLabelPos(int(it.startNodeNum))
		label := it.trie.labels.Read(pos)
		it.append(label, pos)
	}
	level := it.level - 1
	pos := it.posInTrie[level]
	if !it.trie.hasChild.ReadBit(int(pos)) {
		label := it.trie.labels.Read(int(pos))
		if label == terminator && !it.trie.isEndOfNode(int(pos)) {
			it.isAtTerminator = true
		}
		it.isValid = true
		return
	}

	for level < it.trie.height {
		// process child
		nodeNum := it.trie.getChildNodeNum(int(pos))
		pos = it.trie.getFirstLabelPos(nodeNum)
		label := it.trie.labels.Read(int(pos))

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
	pos, ok := it.nextPos()
	if !ok {
		return
	}
	// it.isAtTerminator = false
	// pos := it.posInTrie[it.level-1]
	// nodeNum := it.trie.getChildNodeNum(pos)
	// pos = it.trie.getFirstLabelPos(nodeNum)
	//
	// for pos >= it.trie.louds.numBits || it.trie.louds.ReadBit(pos) {
	// 	// if not child, try find next node of parent
	// 	if it.level == 0 {
	// 		it.isValid = false
	// 		return
	// 	}
	// 	it.level--                       // goto parent
	// 	pos = it.posInTrie[it.level] + 1 // brother node
	// }
	//
	// read next label
	label := it.trie.labels.Read(pos)
	it.append(label, pos)
	// read more lable
	it.moveToLeftMostKey()
}

func (it *loudsSparseIterator) nextPos() (int, bool) {
	it.isAtTerminator = false
	pos := it.posInTrie[it.level-1]
	nodeNum := it.trie.getChildNodeNum(pos)
	pos = it.trie.getFirstLabelPos(nodeNum)

	for pos >= it.trie.louds.numBits || it.trie.louds.ReadBit(pos) {
		// if not child, try find next node of parent
		if it.level == 0 {
			it.isValid = false
			return -1, false
		}
		it.level--                       // goto parent
		pos = it.posInTrie[it.level] + 1 // brother node
	}
	return pos, true
}

func (it *loudsSparseIterator) seek(prefix []byte) bool {
	it.reset()

	nodeNum := it.startNodeNum
	var ok bool
	pos := it.trie.getFirstLabelPos(nodeNum)

	level := 0
	for level = it.startLevel; level < len(prefix); level++ {
		nodeSize := it.trie.nodeSize(pos)
		if pos, ok = it.trie.labels.Search(prefix[level], pos, nodeSize); !ok {
			// if no exact match
			it.moveToLeftInNextSubTrie(pos, nodeSize, prefix[level])
			return false
		}
		it.append(prefix[level], pos)

		// if trie brach terminates
		if !it.trie.hasChild.ReadBit(pos) {
			return it.compareSuffixGreaterThan()
		}

		// move to child
		nodeNum = it.trie.getChildNodeNum(pos)
		pos = it.trie.getFirstLabelPos(nodeNum)
	}
	if it.trie.labels.Read(pos) == terminator &&
		!it.trie.hasChild.ReadBit(pos) &&
		!it.trie.isEndOfNode(pos) {
		it.append(terminator, pos)
		it.isAtTerminator = true
		it.isValid = true
		return false
	}

	if len(prefix) <= level {
		it.moveToLeftMostKey()
		return false
	}

	it.isValid = true
	return true
}

func (it *loudsSparseIterator) compareSuffixGreaterThan() bool {
	//TODO: fix
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
	return it.key[:kLen]
}

// append appends label to key buffer, and goto next level.
func (it *loudsSparseIterator) append(label byte, pos int) {
	it.key[it.level] = label
	it.posInTrie[it.level] = pos
	it.level++
}

func (it *loudsSparseIterator) appendByPos(pos int) {
	it.key[it.level] = it.trie.labels.Read(pos)
	it.posInTrie[it.level] = pos
	it.level++
}
