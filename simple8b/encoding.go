// Package simple8b implements the 64bit integer encoding algoritm as published
// by Ann and Moffat in "Index compression using 64-bit words", Softw. Pract. Exper. 2010; 40:131–147
//
// It is capable of encoding multiple integers with values betweeen 0 and to 1^60 -1, in a single word.
package simple8b

// Simple8b is 64bit word-sized encoder that packs multiple integers into a single word using
// a 4 bit selector values and up to 60 bits for the remaining values.  Integers are encoded using
// the following table:
//
// ┌──────────────┬─────────────────────────────────────────────────────────────┐
// │   Selector   │       0    1   2   3   4   5   6   7  8  9  0 11 12 13 14 15│
// ├──────────────┼─────────────────────────────────────────────────────────────┤
// │     Bits     │       0    0   1   2   3   4   5   6  7  8 10 12 15 20 30 60│
// ├──────────────┼─────────────────────────────────────────────────────────────┤
// │      N       │     240  120  60  30  20  15  12  10  8  7  6  5  4  3  2  1│
// ├──────────────┼─────────────────────────────────────────────────────────────┤
// │   Wasted Bits│      60   60   0   0   0   0  12   0  4  4  0  0  0  0  0  0│
// └──────────────┴─────────────────────────────────────────────────────────────┘
//
// For example, when the number of values can be encoded using 4 bits, selected 5 is encoded in the
// 4 most significant bits followed by 15 values encoded used 4 bits each in the remaing 60 bits.
import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

const MaxValue = (1 << 60) - 1

// Encoder converts a stream of unsigned 64bit integers to a compressed byte slice.
type Encoder struct {
	// most recently written integers that have not been flushed
	buf []uint64

	// index in buf of the head of the buf
	h int

	// index in buf of the tail of the buf
	t int

	// index into bytes of written bytes
	bp int

	// current bytes written and flushed
	bytes []byte
	b     []byte
}

// NewEncoder returns an Encoder able to convert uint64s to compressed byte slices
func NewEncoder() *Encoder {
	return &Encoder{
		buf:   make([]uint64, 240),
		b:     make([]byte, 8),
		bytes: make([]byte, 128),
	}
}

func (e *Encoder) SetValues(v []uint64) {
	e.buf = v
	e.t = len(v)
	e.h = 0
	e.bytes = e.bytes[:0]
}

func (e *Encoder) Reset() {
	e.t = 0
	e.h = 0
	e.bp = 0

	e.buf = e.buf[:240]
	e.b = e.b[:8]
	e.bytes = e.bytes[:128]
}

func (e *Encoder) Write(v uint64) error {
	if e.t >= len(e.buf) {
		if err := e.flush(); err != nil {
			return err
		}
	}

	// The buf is full but there is space at the front, just shift
	// the values down for now. TODO: use ring buffer
	if e.t >= len(e.buf) {
		copy(e.buf, e.buf[e.h:])
		e.t -= e.h
		e.h = 0
	}
	e.buf[e.t] = v
	e.t += 1
	return nil
}

func (e *Encoder) flush() error {
	if e.t == 0 {
		return nil
	}

	// encode as many values into one as we can
	encoded, n, err := Encode(e.buf[e.h:e.t])
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint64(e.b, encoded)
	if e.bp+8 > len(e.bytes) {
		e.bytes = append(e.bytes, e.b...)
		e.bp = len(e.bytes)
	} else {
		copy(e.bytes[e.bp:e.bp+8], e.b)
		e.bp += 8
	}

	// Move the head forward since we encoded those values
	e.h += n

	// If we encoded them all, reset the head/tail pointers to the beginning
	if e.h == e.t {
		e.h = 0
		e.t = 0
	}

	return nil
}

func (e *Encoder) Bytes() ([]byte, error) {
	for e.t > 0 {
		if err := e.flush(); err != nil {
			return nil, err
		}
	}

	return e.bytes[:e.bp], nil
}

// Decoder converts a compressed byte slice to a stream of unsigned 64bit integers.
type Decoder struct {
	bytes []byte
	buf   [240]uint64
	i     int
	n     int
}

// NewDecoder returns a Decoder from a byte slice
func NewDecoder(b []byte) *Decoder {
	return &Decoder{
		bytes: b,
	}
}

// Next returns true if there are remaining values to be read.  Successive
// calls to Next advance the current element pointer.
func (d *Decoder) Next() bool {
	d.i += 1

	if d.i >= d.n {
		d.read()
	}

	return len(d.bytes) >= 8 || (d.i >= 0 && d.i < d.n)
}

func (d *Decoder) SetBytes(b []byte) {
	d.bytes = b
	d.i = 0
	d.n = 0
}

// Read returns the current value.  Successive calls to Read return the same
// value.
func (d *Decoder) Read() uint64 {
	v := d.buf[d.i]
	return v
}

func (d *Decoder) read() {
	if len(d.bytes) < 8 {
		return
	}

	v := binary.BigEndian.Uint64(d.bytes[:8])
	d.bytes = d.bytes[8:]
	d.n, _ = Decode(&d.buf, v)
	d.i = 0
}

type packing struct {
	n, bit int
	unpack func(uint64, *[240]uint64)
	pack   func([]uint64) uint64
}

var selector [16]packing = [16]packing{
	packing{240, 0, unpack240, pack240},
	packing{120, 0, unpack120, pack120},
	packing{60, 1, unpack60, pack60},
	packing{30, 2, unpack30, pack30},
	packing{20, 3, unpack20, pack20},
	packing{15, 4, unpack15, pack15},
	packing{12, 5, unpack12, pack12},
	packing{10, 6, unpack10, pack10},
	packing{8, 7, unpack8, pack8},
	packing{7, 8, unpack7, pack7},
	packing{6, 10, unpack6, pack6},
	packing{5, 12, unpack5, pack5},
	packing{4, 15, unpack4, pack4},
	packing{3, 20, unpack3, pack3},
	packing{2, 30, unpack2, pack2},
	packing{1, 60, unpack1, pack1},
}

// Count returns the number of integers encoded in the byte slice
func CountBytes(b []byte) (int, error) {
	var count int
	for len(b) >= 8 {
		v := binary.BigEndian.Uint64(b[:8])
		b = b[8:]
		n, err := Count(v)
		if err != nil {
			return 0, err
		}

		count += n
	}

	if len(b) > 0 {
		return 0, fmt.Errorf("invalid slice len remaining: %v", len(b))
	}
	return count, nil
}

// Count returns the number of integers encoded within an uint64
func Count(v uint64) (int, error) {
	sel := v >> 60
	if sel >= 16 {
		return 0, fmt.Errorf("invalid selector value: %v", sel)
	}
	return selector[sel].n, nil
}

// Encode packs as many values into a single uint64.  It returns the packed
// uint64, how many values from src were packed, or an error if the values exceed
// the maximum value range.
func Encode(src []uint64) (value uint64, n int, err error) {
	if canPack(src, 240, 0) {
		return uint64(0), 240, nil
	} else if canPack(src, 120, 0) {
		return 1 << 60, 120, nil
	} else if canPack(src, 60, 1) {
		return pack60(src[:60]), 60, nil
	} else if canPack(src, 30, 2) {
		return pack30(src[:30]), 30, nil
	} else if canPack(src, 20, 3) {
		return pack20(src[:20]), 20, nil
	} else if canPack(src, 15, 4) {
		return pack15(src[:15]), 15, nil
	} else if canPack(src, 12, 5) {
		return pack12(src[:12]), 12, nil
	} else if canPack(src, 10, 6) {
		return pack10(src[:10]), 10, nil
	} else if canPack(src, 8, 7) {
		return pack8(src[:8]), 8, nil
	} else if canPack(src, 7, 8) {
		return pack7(src[:7]), 7, nil
	} else if canPack(src, 6, 10) {
		return pack6(src[:6]), 6, nil
	} else if canPack(src, 5, 12) {
		return pack5(src[:5]), 5, nil
	} else if canPack(src, 4, 15) {
		return pack4(src[:4]), 4, nil
	} else if canPack(src, 3, 20) {
		return pack3(src[:3]), 3, nil
	} else if canPack(src, 2, 30) {
		return pack2(src[:2]), 2, nil
	} else if canPack(src, 1, 60) {
		return pack1(src[:1]), 1, nil
	} else {
		if len(src) > 0 {
			return 0, 0, fmt.Errorf("value out of bounds: %v", src)
		}
		return 0, 0, nil
	}
}

// Encode returns a packed slice of the values from src.  If a value is over
// 1 << 60, an error is returned.  The input src is modified to avoid extra
// allocations.  If you need to re-use, use a copy.
func EncodeAll(src []uint64) ([]uint64, error) {
	i := 0

	// Re-use the input slice and write encoded values back in place
	dst := src
	j := 0

	for {
		if i >= len(src) {
			break
		}
		remaining := src[i:]

		if canPack(remaining, 240, 0) {
			dst[j] = 0
			i += 240
		} else if canPack(remaining, 120, 0) {
			dst[j] = 1 << 60
			i += 120
		} else if canPack(remaining, 60, 1) {
			dst[j] = pack60(src[i : i+60])
			i += 60
		} else if canPack(remaining, 30, 2) {
			dst[j] = pack30(src[i : i+30])
			i += 30
		} else if canPack(remaining, 20, 3) {
			dst[j] = pack20(src[i : i+20])
			i += 20
		} else if canPack(remaining, 15, 4) {
			dst[j] = pack15(src[i : i+15])
			i += 15
		} else if canPack(remaining, 12, 5) {
			dst[j] = pack12(src[i : i+12])
			i += 12
		} else if canPack(remaining, 10, 6) {
			dst[j] = pack10(src[i : i+10])
			i += 10
		} else if canPack(remaining, 8, 7) {
			dst[j] = pack8(src[i : i+8])
			i += 8
		} else if canPack(remaining, 7, 8) {
			dst[j] = pack7(src[i : i+7])
			i += 7
		} else if canPack(remaining, 6, 10) {
			dst[j] = pack6(src[i : i+6])
			i += 6
		} else if canPack(remaining, 5, 12) {
			dst[j] = pack5(src[i : i+5])
			i += 5
		} else if canPack(remaining, 4, 15) {
			dst[j] = pack4(src[i : i+4])
			i += 4
		} else if canPack(remaining, 3, 20) {
			dst[j] = pack3(src[i : i+3])
			i += 3
		} else if canPack(remaining, 2, 30) {
			dst[j] = pack2(src[i : i+2])
			i += 2
		} else if canPack(remaining, 1, 60) {
			dst[j] = pack1(src[i : i+1])
			i += 1
		} else {
			return nil, fmt.Errorf("value out of bounds")
		}
		j += 1
	}
	return dst[:j], nil
}

func Decode(dst *[240]uint64, v uint64) (n int, err error) {
	sel := v >> 60
	if sel >= 16 {
		return 0, fmt.Errorf("invalid selector value: %b", sel)
	}
	selector[sel].unpack(v, dst)
	return selector[sel].n, nil
}

// Decode writes the uncompressed values from src to dst.  It returns the number
// of values written or an error.
func DecodeAll(dst, src []uint64) (value int, err error) {
	j := 0
	for _, v := range src {
		sel := v >> 60
		if sel >= 16 {
			return 0, fmt.Errorf("invalid selector value: %b", sel)
		}
		selector[sel].unpack(v, (*[240]uint64)(unsafe.Pointer(&dst[j])))
		j += selector[sel].n
	}
	return j, nil
}

// canPack returs true if n elements from in can be stored using bits per element
func canPack(src []uint64, n, bits int) bool {
	if len(src) < n {
		return false
	}

	end := len(src)
	if n < end {
		end = n
	}

	// Selector 0,1 are special and use 0 bits to encode runs of 1's
	if bits == 0 {
		for _, v := range src {
			if v != 1 {
				return false
			}
		}
		return true
	}

	max := uint64((1 << uint64(bits)) - 1)

	for i := 0; i < end; i++ {
		if src[i] > max {
			return false
		}
	}

	return true
}

// pack240 packs 240 ones from in using 1 bit each
func pack240(src []uint64) uint64 {
	return 0
}

// pack120 packs 120 ones from in using 1 bit each
func pack120(src []uint64) uint64 {
	return 0
}

// pack60 packs 60 values from in using 1 bit each
func pack60(src []uint64) uint64 {
	return 2<<60 |
		src[0] |
		src[1]<<1 |
		src[2]<<2 |
		src[3]<<3 |
		src[4]<<4 |
		src[5]<<5 |
		src[6]<<6 |
		src[7]<<7 |
		src[8]<<8 |
		src[9]<<9 |
		src[10]<<10 |
		src[11]<<11 |
		src[12]<<12 |
		src[13]<<13 |
		src[14]<<14 |
		src[15]<<15 |
		src[16]<<16 |
		src[17]<<17 |
		src[18]<<18 |
		src[19]<<19 |
		src[20]<<20 |
		src[21]<<21 |
		src[22]<<22 |
		src[23]<<23 |
		src[24]<<24 |
		src[25]<<25 |
		src[26]<<26 |
		src[27]<<27 |
		src[28]<<28 |
		src[29]<<29 |
		src[30]<<30 |
		src[31]<<31 |
		src[32]<<32 |
		src[33]<<33 |
		src[34]<<34 |
		src[35]<<35 |
		src[36]<<36 |
		src[37]<<37 |
		src[38]<<38 |
		src[39]<<39 |
		src[40]<<40 |
		src[41]<<41 |
		src[42]<<42 |
		src[43]<<43 |
		src[44]<<44 |
		src[45]<<45 |
		src[46]<<46 |
		src[47]<<47 |
		src[48]<<48 |
		src[49]<<49 |
		src[50]<<50 |
		src[51]<<51 |
		src[52]<<52 |
		src[53]<<53 |
		src[54]<<54 |
		src[55]<<55 |
		src[56]<<56 |
		src[57]<<57 |
		src[58]<<58 |
		src[59]<<59

}

// pack30 packs 30 values from in using 2 bits each
func pack30(src []uint64) uint64 {
	return 3<<60 |
		src[0] |
		src[1]<<2 |
		src[2]<<4 |
		src[3]<<6 |
		src[4]<<8 |
		src[5]<<10 |
		src[6]<<12 |
		src[7]<<14 |
		src[8]<<16 |
		src[9]<<18 |
		src[10]<<20 |
		src[11]<<22 |
		src[12]<<24 |
		src[13]<<26 |
		src[14]<<28 |
		src[15]<<30 |
		src[16]<<32 |
		src[17]<<34 |
		src[18]<<36 |
		src[19]<<38 |
		src[20]<<40 |
		src[21]<<42 |
		src[22]<<44 |
		src[23]<<46 |
		src[24]<<48 |
		src[25]<<50 |
		src[26]<<52 |
		src[27]<<54 |
		src[28]<<56 |
		src[29]<<58
}

// pack20 packs 20 values from in using 3 bits each
func pack20(src []uint64) uint64 {
	return 4<<60 |
		src[0] |
		src[1]<<3 |
		src[2]<<6 |
		src[3]<<9 |
		src[4]<<12 |
		src[5]<<15 |
		src[6]<<18 |
		src[7]<<21 |
		src[8]<<24 |
		src[9]<<27 |
		src[10]<<30 |
		src[11]<<33 |
		src[12]<<36 |
		src[13]<<39 |
		src[14]<<42 |
		src[15]<<45 |
		src[16]<<48 |
		src[17]<<51 |
		src[18]<<54 |
		src[19]<<57
}

// pack15 packs 15 values from in using 3 bits each
func pack15(src []uint64) uint64 {
	return 5<<60 |
		src[0] |
		src[1]<<4 |
		src[2]<<8 |
		src[3]<<12 |
		src[4]<<16 |
		src[5]<<20 |
		src[6]<<24 |
		src[7]<<28 |
		src[8]<<32 |
		src[9]<<36 |
		src[10]<<40 |
		src[11]<<44 |
		src[12]<<48 |
		src[13]<<52 |
		src[14]<<56
}

// pack12 packs 12 values from in using 5 bits each
func pack12(src []uint64) uint64 {
	return 6<<60 |
		src[0] |
		src[1]<<5 |
		src[2]<<10 |
		src[3]<<15 |
		src[4]<<20 |
		src[5]<<25 |
		src[6]<<30 |
		src[7]<<35 |
		src[8]<<40 |
		src[9]<<45 |
		src[10]<<50 |
		src[11]<<55
}

// pack10 packs 10 values from in using 6 bits each
func pack10(src []uint64) uint64 {
	return 7<<60 |
		src[0] |
		src[1]<<6 |
		src[2]<<12 |
		src[3]<<18 |
		src[4]<<24 |
		src[5]<<30 |
		src[6]<<36 |
		src[7]<<42 |
		src[8]<<48 |
		src[9]<<54
}

// pack8 packs 8 values from in using 7 bits each
func pack8(src []uint64) uint64 {
	return 8<<60 |
		src[0] |
		src[1]<<7 |
		src[2]<<14 |
		src[3]<<21 |
		src[4]<<28 |
		src[5]<<35 |
		src[6]<<42 |
		src[7]<<49
}

// pack7 packs 7 values from in using 8 bits each
func pack7(src []uint64) uint64 {
	return 9<<60 |
		src[0] |
		src[1]<<8 |
		src[2]<<16 |
		src[3]<<24 |
		src[4]<<32 |
		src[5]<<40 |
		src[6]<<48
}

// pack6 packs 6 values from in using 10 bits each
func pack6(src []uint64) uint64 {
	return 10<<60 |
		src[0] |
		src[1]<<10 |
		src[2]<<20 |
		src[3]<<30 |
		src[4]<<40 |
		src[5]<<50
}

// pack5 packs 5 values from in using 12 bits each
func pack5(src []uint64) uint64 {
	return 11<<60 |
		src[0] |
		src[1]<<12 |
		src[2]<<24 |
		src[3]<<36 |
		src[4]<<48
}

// pack4 packs 4 values from in using 15 bits each
func pack4(src []uint64) uint64 {
	return 12<<60 |
		src[0] |
		src[1]<<15 |
		src[2]<<30 |
		src[3]<<45
}

// pack3 packs 3 values from in using 20 bits each
func pack3(src []uint64) uint64 {
	return 13<<60 |
		src[0] |
		src[1]<<20 |
		src[2]<<40
}

// pack2 packs 2 values from in using 30 bits each
func pack2(src []uint64) uint64 {
	return 14<<60 |
		src[0] |
		src[1]<<30
}

// pack1 packs 1 values from in using 60 bits each
func pack1(src []uint64) uint64 {
	return 15<<60 |
		src[0]
}
