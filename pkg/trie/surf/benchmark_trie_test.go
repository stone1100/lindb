package surf

import (
	"bytes"
	"math"
	"testing"
)

// 1321 ns/op
func BenchmarkBuilder_Write(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := NewBuilder()
		builder.Build([][]byte{
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
		_ = builder.Write(w)
	}
}

func BenchmarkTrie_Write(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
		w := bytes.NewBuffer([]byte{})
		_ = trie.Write(w)
	}
}
func BenchmarkTrie_Iterator(b *testing.B) {
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := trie.Iterator()
		it.First()
		for it.IsValid() {
			it.Next()
		}
	}
}

var (
	ips, ranks = newTestIPs(1 << 8)
	maxLen     = 0
)

func init() {
	for _, k := range ips {
		maxLen = int(math.Max(float64(maxLen), float64(len(k))))
	}
}

func BenchmarkTrie_MarshalBinary(b *testing.B) {
	b.StopTimer()
	builder := NewBuilder()
	buf := &bytes.Buffer{}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		builder.SetLevel(maxLen + 1)
		builder.Build(ips, ranks)
		_ = builder.Write(buf)
		buf.Reset()
		builder.Reset()
	}
}
