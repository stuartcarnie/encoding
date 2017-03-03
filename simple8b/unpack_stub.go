package simple8b

//go:generate python -m peachpy.x86_64 unpack.py -S -o unpack_amd64.s -mabi=goasm

//go:noescape
func unpack240SSE(v uint64, dst []uint64)

//go:noescape
func unpack240AVX2(v uint64, dst []uint64)

//go:noescape
func unpack60AVX2(v uint64, dst []uint64)

//go:noescape
func unpack30AVX2(v uint64, dst []uint64)

//go:noescape
func unpack20AVX2(v uint64, dst []uint64)
