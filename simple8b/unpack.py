from peachpy import *
from peachpy.x86_64 import *


def make_unpack240_sse():
	v = Argument(uint64_t)
	dst_base = Argument(ptr())
	dst_len = Argument(size_t)
	dst_cap = Argument(size_t)

	with Function("unpack240SSE", (v, dst_base, dst_len, dst_cap), target=uarch.default + isa.sse4_1) as function:
		reg_dst_base = GeneralPurposeRegister64()
		reg_dst_len = GeneralPurposeRegister64()

		LOAD.ARGUMENT(reg_dst_base, dst_base)
		LOAD.ARGUMENT(reg_dst_len, dst_len)
		SHR(reg_dst_len, 3)

		x1 = XMMRegister()
		tmp = GeneralPurposeRegister64()
		MOV(tmp, 1)
		MOVQ(x1, tmp)
		MOV(tmp, 1)
		PINSRQ(x1, tmp, 1)

		with Loop() as loop:
			MOVDQU([reg_dst_base], x1)
			MOVDQU([reg_dst_base + 16], x1)
			MOVDQU([reg_dst_base + 32], x1)
			MOVDQU([reg_dst_base + 48], x1)
			ADD(reg_dst_base, 64)
			SUB(reg_dst_len, 1)
			JNZ(loop.begin)

		RETURN()


def make_unpack240_avx2():
	v = Argument(uint64_t)
	dst_base = Argument(ptr())
	dst_len = Argument(size_t)
	dst_cap = Argument(size_t)

	with Function("unpack240AVX2", (v, dst_base, dst_len, dst_cap), target=uarch.default + isa.avx2):
		reg_dst_base = GeneralPurposeRegister64()
		reg_dst_len = GeneralPurposeRegister64()

		LOAD.ARGUMENT(reg_dst_base, dst_base)
		LOAD.ARGUMENT(reg_dst_len, dst_len)
		SHR(reg_dst_len, 3)

		x1 = XMMRegister()
		tmp = GeneralPurposeRegister64()
		MOV(tmp, 1)
		MOVQ(x1, tmp)

		r_mask = YMMRegister()
		VPBROADCASTQ(r_mask, x1)

		with Loop() as loop:
			VMOVDQU([reg_dst_base], r_mask)
			VMOVDQU([reg_dst_base + 32], r_mask)
			ADD(reg_dst_base, 64)
			SUB(reg_dst_len, 1)
			JNZ(loop.begin)

		RETURN()


class Unpack:
	def __init__(self, size, bits):
		self.name = "unpack%dAVX2" % size
		self.size = size
		self.count = size >> 2
		self.rem = size & 3
		self.bits = bits
		self.shift = 4 * bits
		self.mask = (1 << self.bits) - 1

	def make_mask(self, n):
		if n == 0:
			return None

		tmp = GeneralPurposeRegister64()
		MOV(tmp, 0x8000000000000000)

		x0 = XMMRegister()
		PXOR(x0, x0)

		MOVQ(x0, tmp)
		if n == 2:
			PINSRQ(x0, tmp, 1)

		x1 = XMMRegister()
		if n == 3:
			PXOR(x1, x1)
			MOVQ(x1, tmp)

		y0 = YMMRegister()
		VINSERTI128(y0, y0, x0, 0)
		if n == 3:
			VINSERTI128(y0, y0, x1, 1)

		return y0

	def generate(self):
		v = Argument(uint64_t)
		dst_base = Argument(ptr())
		dst_len = Argument(size_t)
		dst_cap = Argument(size_t)

		with Function(self.name, (v, dst_base, dst_len, dst_cap), target=uarch.default + isa.avx2) as function:
			reg_v = GeneralPurposeRegister64()
			reg_dst_base = GeneralPurposeRegister64()

			LOAD.ARGUMENT(reg_v, v)
			LOAD.ARGUMENT(reg_dst_base, dst_base)

			x0 = XMMRegister()
			x1 = XMMRegister()

			MOVQ(x0, reg_v)
			SHR(reg_v, self.bits)
			PINSRQ(x0, reg_v, 1)
			SHR(reg_v, self.bits)
			MOVQ(x1, reg_v)
			SHR(reg_v, self.bits)
			PINSRQ(x1, reg_v, 1)

			tmp = GeneralPurposeRegister64()

			x2 = XMMRegister()
			MOV(tmp, self.mask)
			MOVQ(x2, tmp)

			mask = self.make_mask(self.rem)

			r_mask = YMMRegister()
			VPBROADCASTQ(r_mask, x2)

			y0 = YMMRegister()
			y1 = YMMRegister()
			VINSERTI128(y0, y0, x0, 0)
			VINSERTI128(y0, y0, x1, 1)

			ofs = 0
			for i in range(self.count):
				VPAND(y1, y0, r_mask)
				VMOVDQU([reg_dst_base + ofs], y1)
				VPSRLQ(y1, y0, self.shift)
				y1, y0 = y0, y1
				ofs += 32
				if ofs == 128:
					ofs = 0
					ADD(reg_dst_base, 128)

			if mask is not None:
				VPAND(y1, y0, r_mask)
				VPMASKMOVQ([reg_dst_base + ofs], mask, y1)

			RETURN()

make_unpack240_sse()
make_unpack240_avx2()
Unpack(60, 1).generate()
Unpack(30, 2).generate()
Unpack(20, 3).generate()
