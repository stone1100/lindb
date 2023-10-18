package surf

import (
	"math/bits"
	"reflect"
	"unsafe"
)

var selectInByteLut [256][8]uint8

func init() {
	for i := 0; i < 256; i++ {
		for j := 0; j < 8; j++ {
			selectInByteLut[i][j] = selectInByte(i, j)
		}
	}
}

func u64SliceToBytes(u []uint64) []byte {
	if len(u) == 0 {
		return nil
	}
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Len = len(u) * 8
	hdr.Cap = hdr.Len
	hdr.Data = uintptr(unsafe.Pointer(&u[0]))
	return b
}

func bytesToU64Slice(b []byte) []uint64 {
	if len(b) == 0 {
		return nil
	}
	var u32s []uint64
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&u32s))
	hdr.Len = len(b) / 8
	hdr.Cap = hdr.Len
	hdr.Data = uintptr(unsafe.Pointer(&b[0]))
	return u32s
}

func findFirstSet(x int) int {
	return bits.TrailingZeros64(uint64(x)) + 1
}

func selectInByte(i, j int) uint8 {
	r := 0
	for ; j != 0; j-- {
		s := findFirstSet(i)
		r += s
		i >>= s
	}
	if i == 0 {
		return 8
	}
	return uint8(r + findFirstSet(i) - 1)
}

func readBit(bs []uint64, pos int) bool {
	wordOff := pos / bitsSize
	bitsOff := pos % bitsSize
	return bs[wordOff]&(uint64(1)<<bitsOff) != 0
}

func select64Broadword(x uint64, nth int64) int64 {
	const (
		onesStep4 = uint64(0x1111111111111111)
		onesStep8 = uint64(0x0101010101010101)
		msbsStep8 = uint64(0x80) * onesStep8
	)

	k := uint64(nth - 1)
	s := x
	s -= (s & (0xa * onesStep4)) >> 1
	s = (s & (0x3 * onesStep4)) + ((s >> 2) & (0x3 * onesStep4))
	s = (s + (s >> 4)) & (0xf * onesStep8)
	byteSums := s * onesStep8

	step8 := k * onesStep8
	geqKStep8 := ((step8 | msbsStep8) - byteSums) & msbsStep8
	place := bits.OnesCount64(geqKStep8) * 8
	byteRank := k - (((byteSums << 8) >> place) & uint64(0xff))
	return int64(place + int(selectInByteLut[(x>>place)&0xff][byteRank]))
}

func popcountBlock(bs []uint64, off, nbits uint32) uint32 {
	if nbits == 0 {
		return 0
	}

	lastWord := (nbits - 1) / bitsSize
	lastBits := (nbits - 1) % bitsSize
	var i, p uint32

	for i = 0; i < lastWord; i++ {
		p += uint32(bits.OnesCount64(bs[off+i]))
	}
	last := bs[off+lastWord] << (bitsSize - 1 - lastBits)
	return p + uint32(bits.OnesCount64(last))
}
