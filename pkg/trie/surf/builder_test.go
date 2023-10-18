package surf

import "testing"

func TestBuild(t *testing.T) {
	b := &Builder{}
	b.Build([][]byte{
		[]byte("hello"),
		[]byte("her"),
		[]byte("hi"),
		[]byte("how"),
		[]byte("seor"),
		[]byte("so"),
	})
}
