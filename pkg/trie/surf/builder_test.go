package surf

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_BuildSuffix(t *testing.T) {
	b := &Builder{}
	b.Build([][]byte{
		[]byte("hello"),
		[]byte("her"),
		[]byte("hi"),
		[]byte("how"),
		[]byte("seor"),
		[]byte("so"),
	}, []uint32{1, 2, 3, 4, 5, 6})
	w := bytes.NewBuffer([]byte{})
	err := b.Write(w)
	assert.NoError(t, err)
	data := w.Bytes()
	fmt.Printf("size=%d\n", len(data))
	trie2 := NewTrie()
	err = trie2.Unmarshal(data)
	assert.NoError(t, err)
	fmt.Println(trie2.String())
	it := trie2.Iterator()
	it.First()
	for it.IsValid() {
		fmt.Printf("key=%s,value=%d\n", string(it.Key()), it.Value())
		it.Next()
	}
}

func TestBuilder_Write(t *testing.T) {
	b := NewBuilder()
	b.Build([][]byte{
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
	w := bytes.NewBuffer([]byte{})
	err := b.Write(w)
	assert.NoError(t, err)
	data := w.Bytes()
	fmt.Printf("size=%d\n", len(data))
	trie2 := NewTrie()
	err = trie2.Unmarshal(data)
	assert.NoError(t, err)
	fmt.Println(trie2.String())
	it := trie2.Iterator()
	it.First()
	for it.IsValid() {
		fmt.Printf("key=%s,value=%d\n", string(it.Key()), it.Value())
		it.Next()
	}
}
