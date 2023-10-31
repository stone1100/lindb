package surf

import (
	"bytes"
	"encoding/binary"
	"io"
	"sort"
	"strings"
)

type LabelVector struct {
	labels []byte
}

func NewLabelVector() *LabelVector {
	return &LabelVector{}
}

func (lv *LabelVector) Init(levels []*Level) {
	numBytes := labelsSize(levels)
	lv.labels = make([]byte, numBytes)

	pos := 0
	for l := range levels {
		levelObj := levels[l]
		copy(lv.labels[pos:], levelObj.lsLabels[:levelObj.item])
		pos += levelObj.item
	}
}

func labelsSize(labels []*Level) int {
	numBytes := 0
	for _, l := range labels {
		numBytes += l.item
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
	if size < 4 {
		for i := start; i < end; i++ {
			if lv.labels[i] == k {
				return i, true
			}
		}
		return off, false
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

func (lv *LabelVector) unmarshal(buf []byte, pos int) (int, error) {
	size := int(UnmarshalUint32(buf, pos))
	pos += 4
	lv.labels = buf[pos : pos+size]
	pos += size
	return pos, nil
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
	offsets       []uint32
	data          []byte
}

func (cpv *compressPathVector) Init(levels []*Level, bitmapType BitmapType) {
	cpv.hasPathVector.Init(rankSparseBlockSize, levels, bitmapType)

	cpv.initData(levels)
}

func (cpv *compressPathVector) initData(levels []*Level) {
	offset := 0
	for level := range levels {
		levelData := levels[level].suffixes
		for idx := range levelData {
			d := levelData[idx]
			cpv.offsets = append(cpv.offsets, uint32(offset))
			offset += len(d)
			cpv.data = append(cpv.data, d...)
		}
	}
}

func (cpv *compressPathVector) GetPath(pos int) []byte {
	if !cpv.hasPathVector.ReadBit(pos) {
		return nil
	}
	endPos := cpv.hasPathVector.Rank(pos) - 1
	start := cpv.offsets[endPos]
	end := len(cpv.data)
	if int(endPos+1) < len(cpv.offsets) {
		end = int(cpv.offsets[endPos+1])
	}
	return cpv.data[start:end]
}

func (cpv *compressPathVector) write(w io.Writer) error {
	var length [8]byte
	binary.LittleEndian.PutUint32(length[:4], uint32(len(cpv.offsets)*4))
	binary.LittleEndian.PutUint32(length[4:], uint32(len(cpv.data)))

	if _, err := w.Write(length[:]); err != nil {
		return err
	}
	if _, err := w.Write(u32SliceToBytes(cpv.offsets)); err != nil {
		return err
	}
	if _, err := w.Write(cpv.data); err != nil {
		return err
	}

	return nil
}

func (cpv *compressPathVector) unmarshal(buf []byte, pos int) (r int, err error) {
	if r, err = cpv.hasPathVector.unmarshal(buf, pos); err != nil {
		return 0, err
	}
	pos = r
	offsetsLen := int(UnmarshalUint32(buf, pos))
	pos += 4
	dataLen := int(UnmarshalUint32(buf, pos))
	pos += 4
	end := pos + offsetsLen
	// read offsets
	cpv.offsets = bytesToU32Slice(buf[pos:end])
	pos = end
	// read data
	cpv.data = buf[pos : pos+dataLen]
	return pos + dataLen, nil
}

type SuffixVector struct {
	compressPathVector
}

func (v *SuffixVector) GetSuffix(pos int) []byte {
	return v.GetPath(pos)
}

func (v *SuffixVector) CheckSuffix(key []byte, level, pos int) bool {
	suffix := v.GetSuffix(pos)
	if level >= len(key) {
		return len(suffix) == 0
	}
	return bytes.Equal(suffix, key[level:])
}

type ValueVector struct {
	values []uint32
}

func (v *ValueVector) Init(levels []*Level) {
	size := 0
	for _, level := range levels {
		size += level.item
	}
	v.values = make([]uint32, size)

	pos := 0
	for level := range levels {
		values := levels[level].values
		for _, val := range values {
			v.values[pos] = val
			pos++
		}
	}
}

func (v *ValueVector) Get(pos int) uint32 {
	return v.values[pos]
}

func (v *ValueVector) write(w io.Writer) error {
	if _, err := w.Write(u32SliceToBytes(v.values)); err != nil {
		return err
	}

	return nil
}

func (v *ValueVector) unmarshal(totalKeys int, buf []byte, pos int) (int, error) {
	dataLen := totalKeys * 4
	end := pos + dataLen
	v.values = bytesToU32Slice(buf[pos:end])
	return end, nil
}
