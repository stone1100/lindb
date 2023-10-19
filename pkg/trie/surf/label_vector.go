package surf

import (
	"bytes"
	"encoding/binary"
	"io"
	"sort"
	"strings"

	"github.com/lindb/lindb/pkg/stream"
)

type LabelVector struct {
	labels []byte
}

func NewLabelVector() *LabelVector {
	return &LabelVector{}
}

func (lv *LabelVector) Init(labels [][]byte, startLevel, trieHeight int) {
	numBytes := labelsSize(labels, startLevel, trieHeight)
	lv.labels = make([]byte, numBytes)

	pos := 0
	for l := startLevel; l < trieHeight; l++ {
		copy(lv.labels[pos:], labels[l])
		pos += len(labels[l])
	}
}

func labelsSize(labels [][]byte, startLevel, trieHeight int) int {
	numBytes := 0 // TODO: 1=> root?
	for l := startLevel; l < trieHeight; l++ {
		numBytes += len(labels[l])
	}
	return numBytes
}

func (ls *LabelVector) Read(pos int) byte {
	return ls.labels[pos]
}

func (lv *LabelVector) Search(k byte, off, size int) (int, bool) {
	start := off
	if size > 1 && lv.labels[start] == terminator {
		start++
		size--
	}

	end := start + size
	if end > len(lv.labels) {
		end = len(lv.labels)
	}
	result := bytes.IndexByte(lv.labels[start:end], k)
	if result < 0 {
		return off, false
	}
	return start + result, true
}

func (lv *LabelVector) SearchGreaterThan(label byte, pos, size int) (int, bool) {
	if size > 1 && lv.labels[pos] == terminator {
		pos++
		size--
	}

	result := sort.Search(size, func(i int) bool { return lv.labels[pos+i] > label })
	if result == size {
		return pos + result - 1, false
	}
	return pos + result, true
}

func (lv *LabelVector) write(w io.Writer) error {
	var bs [4]byte
	binary.LittleEndian.PutUint32(bs[:], uint32(len(lv.labels)))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	if _, err := w.Write(lv.labels); err != nil {
		return err
	}
	// FIXME:??
	// padding := v.MarshalSize() - v.rawMarshalSize()
	// var zeros [8]byte
	// _, err := w.Write(zeros[:padding])
	return nil
}

func (lv *LabelVector) unmarshal(reader *stream.Reader) error {
	size := reader.ReadUint32()
	lv.labels = reader.ReadSlice(int(size))
	return nil
}

func (lv *LabelVector) String() string {
	sb := strings.Builder{}
	for _, c := range lv.labels {
		sb.WriteRune(rune(c))
	}
	return sb.String()
}

type compressPathVector struct {
	hasPathVector BitVectorRank
	offsets       []int
	data          []byte
}

func (cpv *compressPathVector) Init(hasPathBits [][]uint64, numNodesPerLevel []int, data [][][]byte, startLevel, trieHeight int) {
	cpv.hasPathVector.Init(rankSparseBlockSize, hasPathBits, numNodesPerLevel, startLevel, trieHeight)
	offset := 0
	for level := startLevel; level < trieHeight; level++ {
		levelData := data[level]
		for idx := range levelData {
			d := levelData[idx]
			cpv.offsets = append(cpv.offsets, offset)
			offset += len(d)
			cpv.data = append(cpv.data, d...)
		}
	}
}

type SuffixVector struct {
	compressPathVector
}
