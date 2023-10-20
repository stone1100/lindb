package surf

import (
	"io"
)

// node num = rank(i)
// i = select1(node num)
//
// first child(i) = select0(rank(i)) + 1
// last child(i) = select0(rank1(i) + 1) - 1
// parent(i) = select1(ranke0(i))
// children(i) = last child(i) - first child(i)
// child(i,num) = first child(i) + num
type Trie struct {
	loudsSparse *loudsSparse
}

func NewTrie() *Trie {
	return &Trie{}
}

func (trie *Trie) Create(keys [][]byte, values []uint32) {
	builder := NewBuilder()
	builder.Build(keys, values)

	// init Louds-Sparse
	trie.loudsSparse = &loudsSparse{}
	trie.loudsSparse.Init(builder)
}

func (trie *Trie) Get(key []byte) (value uint32, exist bool) {
	return trie.loudsSparse.lookupKey(key)
}

func (trie *Trie) Iterator() *Iterator {
	return NewIterator(trie)
}

func (trie *Trie) Write(w io.Writer) (err error) {
	// write Louds-Sparse
	if err = trie.loudsSparse.write(w); err != nil {
		return err
	}
	return nil
}

func (trie *Trie) Unmarshal(buf []byte) (err error) {
	// unmarshal Louds-Sparse
	loudsSparse := &loudsSparse{}
	if err = loudsSparse.unmarshal(buf); err != nil {
		return err
	}
	trie.loudsSparse = loudsSparse
	return nil
}

func (trie *Trie) String() string {
	return trie.loudsSparse.String()
}
