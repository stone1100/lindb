package surf

import (
	"encoding/binary"
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

func (bv *BitVector) Init(bitsPerLevel [][][]uint64, idx int, numNodesPerLevel []int) {
	bv.totalNumBits(numNodesPerLevel)
	bv.bits = make([]uint64, bv.numWords())

	bitShift := 0
	wordID := 0
	for level, bitsBlock := range bitsPerLevel {
		n := numNodesPerLevel[level]
		if n == 0 {
			continue
		}
		numCompleteWords := n / bitsSize
		for word := 0; word < numCompleteWords; word++ {
			bv.bits[wordID] |= bitsBlock[word][idx] << bitShift
			wordID++
			if bitShift > 0 {
				bv.bits[wordID] |= bitsBlock[word][idx] >> (bitsSize - bitShift)
			}
		}
		remain := n % bitsSize
		if remain > 0 {
			lastWord := bitsBlock[numCompleteWords][idx]
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

func (bv *BitVector) totalNumBits(numNodesPerLevel []int) {
	for level := range numNodesPerLevel {
		bv.numBits += numNodesPerLevel[level]
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

	numOnes int
	// LookUp Table(LUTSs) to store a sampling of precomputed results
	selectLut []int
}

func (bvs *BitVectorSelect) Init(bitsPerLevel [][][]uint64, idx int, numNodesPerLevel []int) {
	bvs.BitVector.Init(bitsPerLevel, idx, numNodesPerLevel)

	bvs.initLut()
}

func (bvs *BitVectorSelect) initLut() {
	lut := []int{0}
	sampledOnes := selectSampleInterval
	onesUptoWord := 0
	for i, w := range bvs.bits {
		ones := bits.OnesCount64(w)
		for sampledOnes <= onesUptoWord+ones {
			diff := sampledOnes - onesUptoWord
			targetPos := i*bitsSize + int(select64(w, int64(diff)))
			lut = append(lut, targetPos)
			sampledOnes += selectSampleInterval
		}
		onesUptoWord += ones
	}

	bvs.numOnes = onesUptoWord
	bvs.selectLut = make([]int, len(lut))
	copy(bvs.selectLut, lut)
}

func (bvs *BitVectorSelect) unmarshal(buf []byte, pos int) (r int, err error) {
	if r, err = bvs.BitVector.unmarshal(buf, pos); err != nil {
		return 0, err
	}
	bvs.initLut()
	return r, nil
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
		return pos
	}

	wordOff := pos / bitsSize
	bitsOff := pos % bitsSize
	if bitsOff == bitsSize-1 {
		wordOff++
		bitsOff = 0
	} else {
		bitsOff++
	}

	// clear low level bits
	w := bvs.bits[wordOff] >> bitsOff << bitsOff
	ones := bits.OnesCount64(w)
	for ones < rankLeft {
		wordOff++
		w = bvs.bits[wordOff]
		rankLeft -= ones
		ones = bits.OnesCount64(w)
	}

	return wordOff*bitsSize + int(select64(w, int64(rankLeft)))
}

type BitVectorRank struct {
	BitVector

	blockSize int
	rankLut   []int
}

func (bvr *BitVectorRank) Init(blockSize int, bitsPerLevel [][][]uint64, idx int, numNodesPerLevel []int) {
	bvr.BitVector.Init(bitsPerLevel, idx, numNodesPerLevel)
	bvr.blockSize = blockSize
	bvr.initLut()
}

func (bvr *BitVectorRank) initLut() {
	wordPerBlk := bvr.blockSize / bitsSize
	nblks := 0
	if bvr.numBits%bvr.blockSize == 0 {
		nblks = bvr.numBits / bvr.blockSize
	} else {
		nblks = bvr.numBits/bvr.blockSize + 1
	}
	bvr.rankLut = make([]int, nblks)

	var totalRank int
	for i := 0; i < nblks-1; i++ {
		bvr.rankLut[i] = totalRank
		totalRank += popcountBlock(bvr.bits, i*wordPerBlk, bvr.blockSize)
	}
	bvr.rankLut[nblks-1] = totalRank
}

// one count [0 pos]
func (bvr *BitVectorRank) Rank(pos int) int {
	wordPreBlk := rankSparseBlockSize / bitsSize
	blockOff := pos / rankSparseBlockSize
	bitsOff := pos % rankSparseBlockSize

	return bvr.rankLut[blockOff] + popcountBlock(bvr.bits, blockOff*wordPreBlk, bitsOff+1)
}

func (bvr *BitVectorRank) unmarshal(buf []byte, pos int) (r int, err error) {
	if r, err = bvr.BitVector.unmarshal(buf, pos); err != nil {
		return 0, nil
	}
	bvr.blockSize = rankSparseBlockSize
	bvr.initLut()
	return r, nil
}
