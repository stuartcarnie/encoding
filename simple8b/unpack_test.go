package simple8b

import (
	"fmt"
	"testing"
)

func toHex(x []uint64) string {
	s := ""
	for _, v := range x {
		s += fmt.Sprintf("%0b ", v)
	}
	return s
}

func BenchmarkUnpack240(b *testing.B) {
	b.SetBytes(240 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack240(0, &dst)
	}
}

func BenchmarkUnpack60(b *testing.B) {
	b.SetBytes(60 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack60(0x6666666666666666, &dst)
	}
}

func BenchmarkUnpack30(b *testing.B) {
	b.SetBytes(30 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack30(0xe4e4e4e4e4e4e4e4, &dst)
	}
}

func BenchmarkUnpack20(b *testing.B) {
	b.SetBytes(20 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack20(0xc688fac688fac688, &dst)
	}
}

func BenchmarkUnpack15(b *testing.B) {
	b.SetBytes(15 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack15(0xc688fac688fac688, &dst)
	}
}
