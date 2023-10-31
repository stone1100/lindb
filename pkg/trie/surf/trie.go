package surf

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
	return &Trie{
		loudsSparse: NewLoudsSparse(),
	}
}

func (trie *Trie) Create(keys [][]byte, values []uint32) {
	builder := NewBuilder()
	builder.Build(keys, values)

	trie.Init(builder)
}

func (trie *Trie) Init(builder *Builder) {
	// init Louds-Sparse
	trie.loudsSparse.Init(builder)
}

func (trie *Trie) Get(key []byte) (value uint32, exist bool) {
	return trie.loudsSparse.lookupKey(key)
}

func (trie *Trie) Iterator() *Iterator {
	return NewIterator(trie)
}

func (trie *Trie) Unmarshal(buf []byte) (err error) {
	// unmarshal Louds-Sparse
	if err = trie.loudsSparse.unmarshal(buf); err != nil {
		return err
	}
	return nil
}

func (trie *Trie) String() string {
	return trie.loudsSparse.String()
}
