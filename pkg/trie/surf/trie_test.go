package surf

import (
	"bytes"
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
	})

	fmt.Println(trie.String())
	// num, ok := trie.Get([]byte("trip"))
	// fmt.Printf("num=%d,ok=%t\n", num, ok)

	it := trie.Iterator()
	it.First()
	for it.IsValid() {
		fmt.Printf("key=%s\n", string(it.Key()))
		it.Next()
	}
	fmt.Println("Seek")
	it.Seek([]byte("fasra"))
	for it.IsValid() {
		fmt.Printf("key=%s\n", string(it.Key()))
		it.Next()
	}
	w := bytes.NewBuffer([]byte{})
	err := trie.Write(w)
	assert.NoError(t, err)
	data := w.Bytes()
	fmt.Printf("size=%d\n", len(data))
	trie2 := NewTrie()
	err = trie2.Unmarshal(data)
	assert.NoError(t, err)
	fmt.Println(trie2.String())
	it = trie2.Iterator()
	it.First()
	for it.IsValid() {
		fmt.Printf("key=%s\n", string(it.Key()))
		it.Next()
	}
}
