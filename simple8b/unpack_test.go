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

func TestUnpack240SSE(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack240SSE(0, dst[:])
	t.Log(dst)
}

func TestUnpack240AVX2(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack240AVX2(0, dst[:])
	t.Log(dst)
}

func TestUnpack60AVX2(t *testing.T) {
	var dst [60]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack60AVX2(0x6666666666666666, dst[:])
	t.Log(toHex(dst[:]))
}

func TestUnpack60(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack60(0x6666666666666666, &dst)
	t.Log(toHex(dst[:60]))
}

func TestUnpack30AVX2(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack30AVX2(0xe4e4e4e4e4e4e4e4, dst[:30])
	t.Log(toHex(dst[:30]))
}

func TestUnpack30(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack30(0xe4e4e4e4e4e4e4e4, &dst)
	t.Log(toHex(dst[:30]))
}

func TestUnpack20AVX2(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack20AVX2(0xc688fac688fac688, dst[:20])
	t.Log(toHex(dst[:20]))
}

func TestUnpack20(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack20(0xc688fac688fac688, &dst)
	t.Log(toHex(dst[:20]))
}

func BenchmarkUnpack240SSE(b *testing.B) {
	b.SetBytes(240 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack240SSE(0, dst[:])
	}
}

func BenchmarkUnpack240AVX2(b *testing.B) {
	b.SetBytes(240 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack240AVX2(0, dst[:])
	}
}

func BenchmarkUnpack240(b *testing.B) {
	b.SetBytes(240 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack240(0, &dst)
	}
}

func BenchmarkUnpack60AVX2(b *testing.B) {
	b.SetBytes(60 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack60AVX2(0x6666666666666666, dst[:60])
	}
}

func BenchmarkUnpack60(b *testing.B) {
	b.SetBytes(60 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack60(0x6666666666666666, &dst)
	}
}

func BenchmarkUnpack30AVX2(b *testing.B) {
	b.SetBytes(30 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack30AVX2(0xe4e4e4e4e4e4e4e4, dst[:30])
	}
}

func BenchmarkUnpack30(b *testing.B) {
	b.SetBytes(30 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack30(0xe4e4e4e4e4e4e4e4, &dst)
	}
}

func BenchmarkUnpack20AVX2(b *testing.B) {
	b.SetBytes(20 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack30AVX2(0xc688fac688fac688, dst[:30])
	}
}

func BenchmarkUnpack20(b *testing.B) {
	b.SetBytes(20 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack20(0xc688fac688fac688, &dst)
	}
}
