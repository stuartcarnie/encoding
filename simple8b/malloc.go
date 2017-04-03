package simple8b

import (
	"reflect"
	"unsafe"
)

const (
	// addrAlignment is the memory address alignment of
	// input and output data. This is required by some
	// SIMD / AVX instructions.
	addrAlignment = 64
)

// from https://github.com/robskie/bp128/blob/master/utils.go#L10, MIT licensed

// MakeAlignedSlice loads an addrAlignment aligned slice
// to dst. dst must be a pointer to an integer slice.
func MakeAlignedSlice(length int, dst interface{}) {
	intSize := 0
	switch dst.(type) {
	case *[]int, *[]uint, *[]int64, *[]uint64:
		intSize = 64
	case *[]int32, *[]uint32:
		intSize = 32
	case *[]int16, *[]uint16:
		intSize = 16
	case *[]int8, *[]uint8:
		intSize = 8
	default:
		panic("bp128: dst is not a pointer to integer slice")
	}

	padding := (addrAlignment * 8) / intSize

	c := length + padding
	vdst := reflect.ValueOf(dst).Elem()
	vslice := reflect.MakeSlice(vdst.Type(), c, c)

	idx := 0
	addr := unsafe.Pointer(vslice.Pointer())
	for !isAligned(intSize, uintptr(addr), idx) {
		idx++
	}

	vdst.Set(vslice.Slice(idx, idx+length))
}

func MakeAligned240() *[240]uint64 {
	var mem []uint64
	MakeAlignedSlice(240, &mem)
	return (*[240]uint64)(unsafe.Pointer(&mem[0]))
}

func isAligned(intSize int, addr uintptr, index int) bool {
	addr += uintptr(index * (intSize / 8))
	return addr&(addrAlignment-1) == 0
}
