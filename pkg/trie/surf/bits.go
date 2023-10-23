package surf

import (
	"encoding/binary"
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

func u16SliceToBytes(u []uint16) []byte {
	if len(u) == 0 {
		return nil
	}
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Len = len(u) * 2
	hdr.Cap = hdr.Len
	hdr.Data = uintptr(unsafe.Pointer(&u[0]))
	return b
}

func bytesToU16Slice(b []byte) []uint16 {
	if len(b) == 0 {
		return nil
	}
	var u16s []uint16
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&u16s))
	hdr.Len = len(b) / 2
	hdr.Cap = hdr.Len
	hdr.Data = uintptr(unsafe.Pointer(&b[0]))
	return u16s
}

func u32SliceToBytes(u []uint32) []byte {
	if len(u) == 0 {
		return nil
	}
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Len = len(u) * 4
	hdr.Cap = hdr.Len
	hdr.Data = uintptr(unsafe.Pointer(&u[0]))
	return b
}

// MarshalUint32 appends marshaled v to dst and returns the result.
func MarshalUint32(dst []byte, u uint32) []byte {
	return append(dst, byte(u>>24), byte(u>>16), byte(u>>8), byte(u))
}

// UnmarshalUint32 returns unmarshaled uint32 from src.
func UnmarshalUint32(src []byte, pos int) uint32 {
	// This is faster than the manual conversion.
	return binary.LittleEndian.Uint32(src[pos:])
}

func bytesToU32Slice(b []byte) []uint32 {
	if len(b) == 0 {
		return nil
	}
	var u32s []uint32
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&u32s))
	hdr.Len = len(b) / 4
	hdr.Cap = hdr.Len
	hdr.Data = uintptr(unsafe.Pointer(&b[0]))
	return u32s
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

func bytesToU64Slice(b []byte, pos, size int) []uint64 {
	if len(b) == 0 {
		return nil
	}
	var u32s []uint64
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&u32s))
	hdr.Len = size
	hdr.Cap = hdr.Len
	hdr.Data = uintptr(unsafe.Pointer(&b[pos]))
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
	// wordOff := pos / bitsSize
	// bitsOff := pos % bitsSize
	// return bs[wordOff]&(uint64(1)<<bitsOff) != 0
	// wordOff := pos / bitsSize
	// bitsOff := pos % bitsSize
	return bs[pos>>6]&(one<<wordsIndex(uint(pos))) != 0
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

func popcountBlock(bs []uint64, off, nbits int) int {
	if nbits == 0 {
		return 0
	}

	lastWord := (nbits - 1) / bitsSize
	lastBits := (nbits - 1) % bitsSize
	var i, p int

	for i = 0; i < lastWord; i++ {
		p += bits.OnesCount64(bs[off+i])
	}
	last := bs[off+lastWord] << (bitsSize - 1 - lastBits)
	return p + bits.OnesCount64(last)
}
