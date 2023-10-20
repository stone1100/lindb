package surf

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type kvPair struct {
	keys   [][]byte
	values []uint32
}

func (p kvPair) Len() int {
	return len(p.keys)
}

func (p kvPair) Less(i, j int) bool {
	return bytes.Compare(p.keys[i], p.keys[j]) < 0
}

func (p kvPair) Swap(i, j int) {
	p.keys[i], p.keys[j] = p.keys[j], p.keys[i]
	p.values[i], p.values[j] = p.values[j], p.values[i]
}

func (p kvPair) Sort() {
	sort.Sort(p)
}

func newTestIPs(batchSize int) (ips [][]byte, ids []uint32) {
	var count int
	for x := 10; x > 0; x-- {
		for y := 1; y < batchSize; y++ {
			for z := batchSize - 1; z > 0; z-- {
				ips = append(ips, []byte(fmt.Sprintf("%d.%d.%d.%d", x, y, y, z)))
				count++
				ids = append(ids, uint32(count))
			}
		}
	}
	kvPair{keys: ips, values: ids}.Sort()
	return
}

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
func TestTrie_Words2(t *testing.T) {
	var keys [][]byte
	var values []uint32
	keysString := []string{
		"A",
		"Aani",
		"a",
		"aa",
	}
	for idx, key := range keysString {
		keys = append(keys, []byte(key))
		values = append(values, uint32(idx))
	}
	kvPair{keys: keys, values: values}.Sort()
	trie := NewTrie()
	trie.Create(keys, values)
	examples := []struct {
		input string
		ok    bool
	}{
		{"A", true},
		{"a", true},
		{"aa", true},
		{"Aani", true},
		{"ab", false},
	}

	for _, example := range examples {
		_, ok := trie.Get([]byte(example.input))
		assert.Equalf(t, example.ok, ok, example.input)
	}
}

func TestTrie_Words(t *testing.T) {
	var keys [][]byte
	var values []uint32
	keysString := []string{
		"a",
		"ab",
		"abc",
		"abcdefgh",
		"abcdefghijklmnopqrstuvwxyz",
		"abcdefghijkl",
		"b",
		"ice",
		"zzzzzz",
	}
	for idx, key := range keysString {
		keys = append(keys, []byte(key))
		values = append(values, uint32(idx))
	}
	kvPair{keys: keys, values: values}.Sort()
	trie := NewTrie()
	trie.Create(keys, values)
	examples := []struct {
		input string
		ok    bool
	}{
		{"a", true},
		{"ab", true},
		{"abc", true},
		{"abcd", false},
		{"abcdefghijklmnopqrstuvwxyz", true},
		{"abcdefghijkl", true},
		{"abcdefghijklm", false},
		{"b", true},
		{"bb", false},
		{"i", false},
		{"ic", false},
		{"ice", true},
		{"ices", false},
		{"zzzzzz", true},
		{"zzzzz", false},
		{"zzzzzzz", false},
	}

	for _, example := range examples {
		_, ok := trie.Get([]byte(example.input))
		assert.Equalf(t, example.ok, ok, example.input)
	}
}

func assertTestData(t *testing.T, path string) {
	var keys [][]byte
	var values []uint32
	f, err := os.Open(path)
	assert.Nil(t, err)
	r, err := gzip.NewReader(f)
	assert.Nil(t, err)

	data, err := io.ReadAll(r)
	assert.Nil(t, err)
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if len([]byte(line)) == 0 {
			continue
		}
		keys = append(keys, []byte(line))
		values = append(values, uint32(i))
	}
	kvPair{keys: keys, values: values}.Sort()
	builder := NewBuilder()
	builder.Build(keys, values)
	trie := NewTrie()
	trie.Init(builder)

	if len(keys) == 0 || len(values) == 0 {
		panic("length is zero")
	}
	for idx := range keys {
		if len(keys[idx]) == 0 {
			fmt.Println(values[idx])
			continue
		}
		value, ok := trie.Get(keys[idx])
		assert.True(t, ok)
		assert.Equal(t, values[idx], value)
	}

	w := bytes.NewBuffer([]byte{})
	err = builder.Write(w)
	assert.NoError(t, err)
	data = w.Bytes()
	fmt.Printf("size=%d\n", len(data))
	trie2 := NewTrie()
	assert.NoError(t, trie2.Unmarshal(data))

	it := trie2.Iterator()
	it.First()
	var idx = 0
	for it.IsValid() {
		assert.Equal(t, values[idx], it.Value())
		assert.Equal(t, keys[idx], it.Key())
		it.Next()
		idx++
	}
}

func TestTrie_TestData_Words(t *testing.T) {
	assertTestData(t, "../testdata/words.txt.gz")
}

func TestTrie_TestData_UUID(t *testing.T) {
	assertTestData(t, "../testdata/uuid.txt.gz")
}

func TestTrie_TestData_Hsk_words(t *testing.T) {
	assertTestData(t, "../testdata/hsk_words.txt.gz")
}
