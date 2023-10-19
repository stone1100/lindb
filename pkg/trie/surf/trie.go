package surf

import (
	"fmt"
	"io"

	"github.com/lindb/lindb/pkg/stream"
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
	loudsDense  *loudsDense
	loudsSparse *loudsSparse
}

func NewTrie() *Trie {
	return &Trie{}
}

func (trie *Trie) Create(keys [][]byte, values []uint32) {
	builder := NewBuilder()
	builder.Build(keys, values)

	// init Louds-Dense
	if builder.getSparseStartLevel() > 0 {
		trie.loudsDense = &loudsDense{}
		trie.loudsDense.Init(builder)
	}

	// init Louds-Sparse
	trie.loudsSparse = &loudsSparse{}
	trie.loudsSparse.Init(builder)
}

func (trie *Trie) Get(key []byte) (nodeNum int, exist bool) {
	if trie.loudsDense == nil {
		return trie.loudsSparse.lookupKey(key, nodeNum)
	}

	if nodeNum, exist = trie.loudsDense.lookupKey(key); !exist {
		return -1, false
	} else if nodeNum != 0 {
		return trie.loudsSparse.lookupKey(key, nodeNum)
	}
	return
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
	reader := stream.NewReader(buf)
	// unmarshal Louds-Sparse
	loudsSparse := &loudsSparse{}
	if err = loudsSparse.unmarshal(reader); err != nil {
		return err
	}
	trie.loudsSparse = loudsSparse
	return nil
}

func (trie *Trie) String() string {
	return fmt.Sprintf("%s\n\n", trie.loudsSparse.String())
}
