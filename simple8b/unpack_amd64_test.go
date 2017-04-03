package simple8b

import (
	"testing"
)

func compare(t *testing.T, a, e []uint64) {
	for i, v := range a {
		if v != e[i] {
			t.Fatalf("mismatch v[%d]; %d != %d ", i, v, e[i])
		}
	}
}

func TestUnpack240SSE(t *testing.T) {
	dst := MakeAligned240()
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack240SSE(0, dst)
	var exp [240]uint64
	unpack240(0, &exp)
	compare(t, dst[:], exp[:])
}

func TestUnpack240SSEUnaligned(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack240SSEu(0, &dst)
	var exp [240]uint64
	unpack240(0, &exp)
	compare(t, dst[:], exp[:])
}

func TestUnpack240AVX2(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack240AVX2(0, &dst)
	var exp [240]uint64
	unpack240(0, &exp)
	compare(t, dst[:], exp[:])
}

func TestUnpack60AVX2(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack60AVX2(0x6666666666666666, &dst)
	var exp [240]uint64
	unpack60(0x6666666666666666, &exp)
	compare(t, dst[:60], exp[:60])
}

func TestUnpack30AVX2(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack30AVX2(0xe4e4e4e4e4e4e4e4, &dst)
	var exp [240]uint64
	unpack30(0xe4e4e4e4e4e4e4e4, &exp)
	compare(t, dst[:30], exp[:30])
}

func TestUnpack20AVX2(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack20AVX2(0xc688fac688fac688, &dst)
	var exp [240]uint64
	unpack20(0xc688fac688fac688, &exp)
	compare(t, dst[:20], exp[:20])
}

func TestUnpack15AVX2(t *testing.T) {
	var dst [240]uint64
	for i := range dst {
		dst[i] = uint64(i)
	}
	unpack15AVX2(0xc688fac688fac689, &dst)
	var exp [240]uint64
	unpack15(0xc688fac688fac689, &exp)
	compare(t, dst[:15], exp[:15])
}

func BenchmarkUnpack240SSE(b *testing.B) {
	b.SetBytes(240 * 8)
	dst := MakeAligned240()
	for i := 0; i < b.N; i++ {
		unpack240SSE(0, dst)
	}
}

func BenchmarkUnpack240SSEuUsingUnaligned(b *testing.B) {
	b.SetBytes(240 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack240SSEu(0, &dst)
	}
}

func BenchmarkUnpack240SSEuUsingAligned(b *testing.B) {
	b.SetBytes(240 * 8)
	dst := MakeAligned240()
	for i := 0; i < b.N; i++ {
		unpack240SSEu(0, dst)
	}
}

func BenchmarkUnpack240AVX2(b *testing.B) {
	b.SetBytes(240 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack240AVX2(0, &dst)
	}
}

func BenchmarkUnpack240AVX2UsingAligned(b *testing.B) {
	dst := MakeAligned240()
	b.SetBytes(240 * 8)
	for i := 0; i < b.N; i++ {
		unpack240AVX2(0, dst)
	}
}

func BenchmarkUnpack60AVX2(b *testing.B) {
	b.SetBytes(60 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack60AVX2(0x6666666666666666, &dst)
	}
}

func BenchmarkUnpack60AVX2UsingAligned(b *testing.B) {
	b.SetBytes(60 * 8)
	dst := MakeAligned240()
	for i := 0; i < b.N; i++ {
		unpack60AVX2(0x6666666666666666, dst)
	}
}

func BenchmarkUnpack30AVX2(b *testing.B) {
	b.SetBytes(30 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack30AVX2(0xe4e4e4e4e4e4e4e4, &dst)
	}
}

func BenchmarkUnpack30AVX2UsingAligned(b *testing.B) {
	b.SetBytes(30 * 8)
	dst := MakeAligned240()
	for i := 0; i < b.N; i++ {
		unpack30AVX2(0xe4e4e4e4e4e4e4e4, dst)
	}
}

func BenchmarkUnpack20AVX2(b *testing.B) {
	b.SetBytes(20 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack20AVX2(0xc688fac688fac688, &dst)
	}
}

func BenchmarkUnpack20AVX2UsingAligned(b *testing.B) {
	b.SetBytes(20 * 8)
	dst := MakeAligned240()
	for i := 0; i < b.N; i++ {
		unpack20AVX2(0xc688fac688fac688, dst)
	}
}

func BenchmarkUnpack15AVX2(b *testing.B) {
	b.SetBytes(15 * 8)
	var dst [240]uint64
	for i := 0; i < b.N; i++ {
		unpack15AVX2(0xc688fac688fac688, &dst)
	}
}

func BenchmarkUnpack15AVX2UsingAligned(b *testing.B) {
	b.SetBytes(15 * 8)
	dst := MakeAligned240()
	for i := 0; i < b.N; i++ {
		unpack15AVX2(0xc688fac688fac688, dst)
	}
}
