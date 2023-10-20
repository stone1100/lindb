package surf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrie_Build(t *testing.T) {
	trie := NewTrie()
	trie.Create([][]byte{
		[]byte("f"),
		[]byte("far"),
		[]byte("fas"),
		[]byte("fast"),
		[]byte("fat"),
		[]byte("s"),
		[]byte("top"),
		[]byte("toy"),
		[]byte("trie"),
		[]byte("trip"),
		[]byte("try"),
	}, []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})

	fmt.Println(trie.String())
	num, ok := trie.Get([]byte("trip"))
	assert.Equal(t, uint32(10), num)
	assert.True(t, ok)

	it := trie.Iterator()
	it.First()
	for it.IsValid() {
		fmt.Printf("key=%s,value=%d\n", string(it.Key()), it.Value())
		it.Next()
	}
	fmt.Println("Seek")
	it.Seek([]byte("fasra"))
	for it.IsValid() {
		fmt.Printf("key=%s,value=%d\n", string(it.Key()), it.Value())
		it.Next()
	}
	// w := bytes.NewBuffer([]byte{})
	// err := trie.Write(w)
	// assert.NoError(t, err)
	// data := w.Bytes()
	// fmt.Printf("size=%d\n", len(data))
	// trie2 := NewTrie()
	// err = trie2.Unmarshal(data)
	// assert.NoError(t, err)
	// fmt.Println(trie2.String())
	// it = trie2.Iterator()
	// it.First()
	// for it.IsValid() {
	// 	fmt.Printf("key=%s,value=%d\n", string(it.Key()), it.Value())
	// 	it.Next()
	// }
}

func TestTrie_BuildSuffix(t *testing.T) {
	trie := NewTrie()
	trie.Create([][]byte{
		[]byte("hello"),
		[]byte("her"),
		[]byte("hi"),
		[]byte("how"),
		[]byte("seor"),
		[]byte("so"),
	}, []uint32{1, 2, 3, 4, 5, 6})
	fmt.Println(trie.String())
	num, ok := trie.Get([]byte("hello"))
	assert.Equal(t, uint32(1), num)
	assert.True(t, ok)
	num, ok = trie.Get([]byte("hel"))
	assert.Equal(t, uint32(0), num)
	assert.False(t, ok)
	it := trie.Iterator()
	it.First()
	for it.IsValid() {
		fmt.Printf("key=%s,value=%d\n", string(it.Key()), it.Value())
		it.Next()
	}

	fmt.Println("seek")
	it = trie.Iterator()
	it.Seek([]byte("hello1"))
	for it.IsValid() {
		fmt.Printf("key=%s,value=%d\n", string(it.Key()), it.Value())
		it.Next()
	}
}
