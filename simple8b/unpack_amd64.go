package simple8b

//go:generate python -m peachpy.x86_64 unpack.py -S -o unpack_amd64.s -mabi=goasm

//go:noescape
func unpack240SSE(v uint64, dst *[240]uint64)

//go:noescape
func unpack240SSEu(v uint64, dst *[240]uint64)

//go:noescape
func unpack240AVX2(v uint64, dst *[240]uint64)

//go:noescape
func unpack60AVX2(v uint64, dst *[240]uint64)

//go:noescape
func unpack30AVX2(v uint64, dst *[240]uint64)

//go:noescape
func unpack20AVX2(v uint64, dst *[240]uint64)

//go:noescape
func unpack15AVX2(v uint64, dst *[240]uint64)

var (
	support_avx2 bool
)

func init() {
	cpu_info()
	if !support_avx2 {
		return
	}

	selector[0].unpack = unpack240AVX2
	selector[2].unpack = unpack60AVX2
	selector[3].unpack = unpack30AVX2
	selector[4].unpack = unpack20AVX2
	selector[5].unpack = unpack15AVX2
}
