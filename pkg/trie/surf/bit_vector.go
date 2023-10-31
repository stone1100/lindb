package surf

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/bits"
	"strings"
)

const (
	rankSparseBlockSize  = 512
	selectSampleInterval = 64
)

type BitVector struct {
	numBits int
	bits    []uint64
}

func (bv *BitVector) Init(levels []*Level, bitmapType BitmapType) {
	bv.totalNumBits(levels)
	bv.bits = make([]uint64, bv.numWords())
	bitShift := 0
	wordID := 0
	for level := range levels {
		levelObj := levels[level]
		n := levelObj.item
		if n == 0 {
			continue
		}
		bitsBlock := levelObj.GetBitmap(bitmapType)
		numCompleteWords := n / bitsSize
		for word := 0; word < numCompleteWords; word++ {
			bv.bits[wordID] |= bitsBlock[word] << bitShift
			wordID++
			if bitShift > 0 {
				bv.bits[wordID] |= bitsBlock[word] >> (bitsSize - bitShift)
			}
		}
		remain := n % bitsSize
		if remain > 0 {
			lastWord := bitsBlock[numCompleteWords]
			bv.bits[wordID] |= lastWord << bitShift
			if bitShift+remain <= bitsSize {
				bitShift = (bitShift + remain) % bitsSize
				if bitShift == 0 {
					wordID++
				}
			} else {
				wordID++
				bv.bits[wordID] |= lastWord >> (bitsSize - bitShift)
				bitShift = bitShift + remain - bitsSize
			}
		}
	}
}

func (bv *BitVector) bitsSize() int {
	return bv.numWords() * 8
}

func (bv *BitVector) write(w io.Writer) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(bv.numBits))
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	if _, err := w.Write(u64SliceToBytes(bv.bits)); err != nil {
		return err
	}
	return nil
}

func (bv *BitVector) unmarshal(buf []byte, pos int) (int, error) {
	bv.numBits = int(UnmarshalUint32(buf, pos))
	pos += 4
	words := bv.numWords()

	bv.bits = bytesToU64Slice(buf, pos, words)
	pos += words * 8
	return pos, nil
}

func (bv *BitVector) DistanceToNextSetBit(pos int) int {
	var distance = 1
	wordOff := (pos + 1) / bitsSize
	bitsOff := (pos + 1) % bitsSize

	if wordOff >= len(bv.bits) {
		return 0
	}

	testBits := bv.bits[wordOff] >> bitsOff
	if testBits > 0 {
		return distance + bits.TrailingZeros64(testBits)
	}

	numWords := bv.numWords()
	if wordOff == numWords-1 {
		return bv.numBits - pos
	}
	distance += bitsSize - bitsOff

	for wordOff < numWords-1 {
		wordOff++
		testBits = bv.bits[wordOff]
		if testBits > 0 {
			return distance + bits.TrailingZeros64(testBits)
		}
		distance += bitsSize
	}

	if wordOff == numWords-1 && bv.numBits%64 != 0 {
		distance -= bitsSize - bv.numBits%64
	}

	return distance
}

func (v *BitVector) ReadBit(pos int) bool {
	return readBit(v.bits, pos)
}

func (v *BitVector) numWords() int {
	if v.numBits%bitsSize == 0 {
		return v.numBits / bitsSize
	}
	return v.numBits/bitsSize + 1
}

func (bv *BitVector) totalNumBits(levels []*Level) {
	for level := range levels {
		bv.numBits += levels[level].item
	}
}

func (bv *BitVector) String() string {
	var s strings.Builder
	for i := 0; i < bv.numBits; i++ {
		if readBit(bv.bits, i) {
			s.WriteString("1")
		} else {
			s.WriteString("0")
		}
	}
	return s.String()
}

type BitVectorSelect struct {
	BitVector

	numOnes uint32
	// LookUp Table(LUTSs) to store a sampling of precomputed results
	selectLut []uint32

	previousLutIdx  int
	previousWordOff uint32
	previousBuf     []int
	previous        int
}

func (bvs *BitVectorSelect) Init(levels []*Level, bitmapType BitmapType) {
	bvs.BitVector.Init(levels, bitmapType)

	bvs.initLut()
}

func (bvs *BitVectorSelect) initLut() {
	lut := []uint32{0}
	sampledOnes := selectSampleInterval
	onesUptoWord := 0
	for i, w := range bvs.bits {
		ones := bits.OnesCount64(w)
		for sampledOnes <= onesUptoWord+ones {
			diff := sampledOnes - onesUptoWord
			targetPos := uint32(i*bitsSize) + uint32(select64(w, int64(diff)))
			lut = append(lut, targetPos)
			sampledOnes += selectSampleInterval
		}
		onesUptoWord += ones
	}

	bvs.numOnes = uint32(onesUptoWord)
	bvs.selectLut = make([]uint32, len(lut))
	copy(bvs.selectLut, lut)
}

func (bvs *BitVectorSelect) writeLut(write io.Writer) error {
	var buf [4]byte
	sampledOnes := selectSampleInterval
	onesUptoWord := 0
	for i, w := range bvs.bits {
		ones := bits.OnesCount64(w)
		for sampledOnes <= onesUptoWord+ones {
			diff := sampledOnes - onesUptoWord
			targetPos := uint32(i*bitsSize) + uint32(select64(w, int64(diff)))
			binary.LittleEndian.PutUint32(buf[:], targetPos)
			if _, err := write.Write(buf[:]); err != nil {
				return err
			}
			sampledOnes += selectSampleInterval
		}
		onesUptoWord += ones
	}

	return nil
}

func (bvs *BitVectorSelect) write(w io.Writer) error {
	if err := bvs.BitVector.write(w); err != nil {
		return err
	}
	// if err := bvs.writeLut(w); err != nil {
	// 	return err
	// }
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], bvs.numOnes)
	_, err := w.Write(buf[:])
	if err != nil {
		return err
	}
	_, err = w.Write(u32SliceToBytes(bvs.selectLut))
	if err != nil {
		return err
	}

	return nil
}

func (bvs *BitVectorSelect) lutSize() uint32 {
	return (bvs.numOnes/selectSampleInterval + 1) * 4
}

func (bvs *BitVectorSelect) unmarshal(buf []byte, pos int) (r int, err error) {
	if r, err = bvs.BitVector.unmarshal(buf, pos); err != nil {
		return 0, err
	}
	bvs.numOnes = UnmarshalUint32(buf, r)
	r += 4
	// read lut
	lutSize := int(bvs.lutSize())
	if len(buf) < lutSize {
		return 0, fmt.Errorf("cannot read lut: %d from selectVector:%d", lutSize, len(buf))
	}
	bvs.selectLut = bytesToU32Slice(buf[r : r+lutSize])

	return r + lutSize, nil
}

// Select returns the position of the rank-th 1 bit.
// position is zero-based; rank is one-based.
// E.g., for bitvector: 100101000, select(3) = 5
func (bvs *BitVectorSelect) Select(rank int) int {
	lutIdx := rank / selectSampleInterval
	rankLeft := rank % selectSampleInterval
	if lutIdx == 0 {
		rankLeft--
	}

	pos := bvs.selectLut[lutIdx]
	if rankLeft == 0 {
		return int(pos)
	}

	wordOff := pos / bitsSize
	bitsOff := pos % bitsSize
	if bitsOff == bitsSize-1 {
		wordOff++
		bitsOff = 0
	} else {
		bitsOff++
	}

	if bvs.previous > 0 && lutIdx == bvs.previousLutIdx {
		// idx := 0

		// ones := bits.OnesCount64(w)
		// for ones < rankLeft && idx < len(bvs.previousBuf) {
		// 	rankLeft -= ones
		// 	idx++
		// 	wordOff++
		// 	w = bvs.bits[wordOff]
		// 	ones = bvs.previousBuf[idx]
		// }
		// var w uint64
		// if wordOff == bvs.previousWordOff {
		// } else {
		// 	w = bvs.bits[wordOff]
		// }
		var w uint64
		idx := 0
		ones := bvs.previousBuf[idx]
		// bits.OnesCount64(w)
		// idx++
		for ones < rankLeft {
			rankLeft -= ones
			idx++
			wordOff++
			if idx < bvs.previous {
				ones = bvs.previousBuf[idx]
			} else {
				w = bvs.bits[wordOff]
				ones = bits.OnesCount64(w)
			}
		}
		if idx > 0 {
			w = bvs.bits[wordOff]
		} else {
			w = bvs.bits[wordOff] >> bitsOff << bitsOff
		}
		// fmt.Println(bvs.previousBuf)
		// fmt.Printf("aw=%d,wo=%d,pwo=%d,l=%d,idx=%d\n", w, wordOff, bvs.previousWordOff, rankLeft, lutIdx)

		return int(wordOff*bitsSize) + int(select64(w, int64(rankLeft)))
	}

	flag := false
	if lutIdx != bvs.previousLutIdx || bvs.previous == 0 {
		bvs.previous = 0
		bvs.previousBuf = bvs.previousBuf[:0]
		bvs.previousWordOff = wordOff
		flag = true
	}

	// clear low level bits
	w := bvs.bits[wordOff] >> bitsOff << bitsOff
	ones := bits.OnesCount64(w)
	if flag {
		bvs.previousBuf = append(bvs.previousBuf, ones)
	}

	for ones < rankLeft {
		rankLeft -= ones
		wordOff++
		w = bvs.bits[wordOff]
		ones = bits.OnesCount64(w)
		if flag {
			bvs.previousBuf = append(bvs.previousBuf, ones)
		}
	}
	bvs.previousLutIdx = lutIdx
	bvs.previous = len(bvs.previousBuf)

	// fmt.Println(bvs.previousBuf)
	// fmt.Printf("w=%d,wo=%d,pwo=%d,l=%d,idx=%d\n", w, wordOff, bvs.previousWordOff, rankLeft, lutIdx)
	return int(wordOff*bitsSize) + int(select64(w, int64(rankLeft)))
}

type BitVectorRank struct {
	BitVector

	blockSize int
	rankLut   []uint32
}

func (bvr *BitVectorRank) Init(blockSize int, levels []*Level, bitmapType BitmapType) {
	bvr.BitVector.Init(levels, bitmapType)
	bvr.blockSize = blockSize
	bvr.initLut()
}

func (bvr *BitVectorRank) Init2(blockSize int, levels []*Level, bitmapType BitmapType) {
	bvr.BitVector.Init(levels, bitmapType)
	bvr.blockSize = blockSize
	// bvr.initLut()
}

func (bvr *BitVectorRank) initLut() {
	wordPerBlk := bvr.blockSize / bitsSize
	nblks := 0
	if bvr.numBits%bvr.blockSize == 0 {
		nblks = bvr.numBits / bvr.blockSize
	} else {
		nblks = bvr.numBits/bvr.blockSize + 1
	}
	bvr.rankLut = make([]uint32, nblks)

	var totalRank uint32
	for i := 0; i < nblks-1; i++ {
		bvr.rankLut[i] = totalRank
		totalRank += uint32(popcountBlock(bvr.bits, i*wordPerBlk, bvr.blockSize))
	}
	bvr.rankLut[nblks-1] = totalRank
}

func (bvr *BitVectorRank) writeLUT(w io.Writer) error {
	wordPerBlk := bvr.blockSize / bitsSize
	nblks := 0
	if bvr.numBits%bvr.blockSize == 0 {
		nblks = bvr.numBits / bvr.blockSize
	} else {
		nblks = bvr.numBits/bvr.blockSize + 1
	}
	// bvr.rankLut = make([]uint32, nblks)
	var buf [4]byte

	var totalRank uint32
	for i := 0; i < nblks-1; i++ {
		binary.LittleEndian.PutUint32(buf[:], uint32(totalRank))
		if _, err := w.Write(buf[:]); err != nil {
			return err
		}
		// bvr.rankLut[i] = totalRank
		totalRank += uint32(popcountBlock(bvr.bits, i*wordPerBlk, bvr.blockSize))
	}
	// bvr.rankLut[nblks-1] = totalRank
	binary.LittleEndian.PutUint32(buf[:], uint32(totalRank))
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	return nil
}

// one count [0 pos]
func (bvr *BitVectorRank) Rank(pos int) int {
	wordPreBlk := bvr.blockSize / bitsSize
	blockOff := pos / bvr.blockSize
	bitsOff := pos % bvr.blockSize

	return int(bvr.rankLut[blockOff]) + popcountBlock(bvr.bits, blockOff*wordPreBlk, bitsOff+1)
}

func (bvr *BitVectorRank) Select(rank int) int {
	wordPreBlk := bvr.blockSize / bitsSize
	blockOff := rank / bvr.blockSize
	// bitsOff := pos % bvr.blockSize

	rankLeft := uint32(rank) - bvr.rankLut[blockOff]
	wordOff := blockOff * wordPreBlk
	// ones := uint32(0)
	w := bvr.bits[wordOff]
	ones := uint32(bits.OnesCount64(w))
	for ones < rankLeft {
		rankLeft -= ones
		wordOff++
		// ones = bvr.rankLut[wordOff]
		w = bvr.bits[wordOff]
		ones = uint32(bits.OnesCount64(w))
	}

	return int(bvr.rankLut[blockOff]) + int(wordOff*bitsSize) + int(select64(w, int64(rankLeft)))
}

func (bvr *BitVectorRank) write(w io.Writer) error {
	if err := bvr.BitVector.write(w); err != nil {
		return err
	}
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(bvr.blockSize))
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	// if _, err := w.Write(u32SliceToBytes(bvr.rankLut)); err != nil {
	// 	return err
	// }
	return bvr.writeLUT(w)
}

func (bvr *BitVectorRank) unmarshal(buf []byte, pos int) (r int, err error) {
	if r, err = bvr.BitVector.unmarshal(buf, pos); err != nil {
		return 0, nil
	}

	bvr.blockSize = int(UnmarshalUint32(buf, r))
	r += 4
	// reading lut
	lutSize := int(bvr.lutSize())
	if len(buf) < lutSize {
		return 0, fmt.Errorf("cannot read lut: %d from rankVector: %d", lutSize, len(buf))
	}
	bvr.rankLut = bytesToU32Slice(buf[r : r+lutSize])
	return r + lutSize, nil
}

func (bvr *BitVectorRank) lutSize() int {
	return (bvr.numBits/bvr.blockSize + 1) * 4
}

type BitVectorRank2 struct {
	BitVector

	blockSize int
	rankLut   []uint32
}

func (bvr *BitVectorRank2) Init(blockSize int, levels []*Level, bitmapType BitmapType) {
	bvr.BitVector.Init(levels, bitmapType)
	bvr.blockSize = blockSize
	bvr.initLut()
}

func (bvr *BitVectorRank2) Init2(blockSize int, levels []*Level, bitmapType BitmapType) {
	bvr.BitVector.Init(levels, bitmapType)
	bvr.blockSize = blockSize
	// bvr.initLut()
}

func (bvr *BitVectorRank2) initLut() {
	wordPerBlk := bvr.blockSize / bitsSize
	nblks := 0
	if bvr.numBits%bvr.blockSize == 0 {
		nblks = bvr.numBits / bvr.blockSize
	} else {
		nblks = bvr.numBits/bvr.blockSize + 1
	}
	bvr.rankLut = make([]uint32, nblks)

	// var totalRank uint32
	for i := 0; i < nblks; i++ {
		bvr.rankLut[i] = uint32(popcountBlock(bvr.bits, i*wordPerBlk, bvr.blockSize))
		// bvr.rankLut[i] = totalRank
	}
	// bvr.rankLut[nblks-1] = totalRank
}

func (bvr *BitVectorRank2) writeLUT(w io.Writer) error {
	wordPerBlk := bvr.blockSize / bitsSize
	nblks := 0
	if bvr.numBits%bvr.blockSize == 0 {
		nblks = bvr.numBits / bvr.blockSize
	} else {
		nblks = bvr.numBits/bvr.blockSize + 1
	}
	// bvr.rankLut = make([]uint32, nblks)
	var buf [4]byte

	var totalRank uint32
	for i := 0; i < nblks; i++ {
		totalRank = uint32(popcountBlock(bvr.bits, i*wordPerBlk, bvr.blockSize))
		binary.LittleEndian.PutUint32(buf[:], uint32(totalRank))
		if _, err := w.Write(buf[:]); err != nil {
			return err
		}
		// bvr.rankLut[i] = totalRank
	}
	// bvr.rankLut[nblks-1] = totalRank
	// binary.LittleEndian.PutUint32(buf[:], uint32(totalRank))
	// if _, err := w.Write(buf[:]); err != nil {
	// 	return err
	// }
	return nil
}

// one count [0 pos]
func (bvr *BitVectorRank2) Rank(pos int) int {
	wordPreBlk := bvr.blockSize / bitsSize
	blockOff := pos / bvr.blockSize
	bitsOff := pos % bvr.blockSize

	return int(bvr.rankLut[blockOff]) + popcountBlock(bvr.bits, blockOff*wordPreBlk, bitsOff+1)
}

func (bvr *BitVectorRank2) Select(rank int) int {
	wordPreBlk := bvr.blockSize / bitsSize
	// blockOff := rank / bvr.blockSize
	// bitsOff := pos % bvr.blockSize

	rankLeft := uint32(rank)
	idx := 0
	// blockOff * wordPreBlk
	// ones := uint32(0)
	ones := bvr.rankLut[idx]
	for ones < rankLeft {
		rankLeft -= ones
		idx++
		// ones = bvr.rankLut[wordOff]
		// ones = uint32(bits.OnesCount64(w))
		ones = bvr.rankLut[idx]
	}
	// - bvr.rankLut[blockOff]
	wordOff := idx * wordPreBlk
	// blockOff * wordPreBlk
	// ones := uint32(0)
	w := bvr.bits[wordOff]
	ones = uint32(bits.OnesCount64(w))
	for ones < rankLeft {
		rankLeft -= ones
		wordOff++
		// ones = bvr.rankLut[wordOff]
		w = bvr.bits[wordOff]
		ones = uint32(bits.OnesCount64(w))
		// ones = bvr.rankLut[wordOff]
	}

	return int(wordOff*bitsSize) + int(select64(w, int64(rankLeft)))
}

func (bvr *BitVectorRank2) write(w io.Writer) error {
	if err := bvr.BitVector.write(w); err != nil {
		return err
	}
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(bvr.blockSize))
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	// if _, err := w.Write(u32SliceToBytes(bvr.rankLut)); err != nil {
	// 	return err
	// }
	return bvr.writeLUT(w)
}

func (bvr *BitVectorRank2) unmarshal(buf []byte, pos int) (r int, err error) {
	if r, err = bvr.BitVector.unmarshal(buf, pos); err != nil {
		return 0, nil
	}

	bvr.blockSize = int(UnmarshalUint32(buf, r))
	r += 4
	// reading lut
	lutSize := int(bvr.lutSize())
	if len(buf) < lutSize {
		return 0, fmt.Errorf("cannot read lut: %d from rankVector: %d", lutSize, len(buf))
	}
	bvr.rankLut = bytesToU32Slice(buf[r : r+lutSize])
	return r + lutSize, nil
}

func (bvr *BitVectorRank2) lutSize() int {
	return (bvr.numBits/bvr.blockSize + 1) * 4
}
