package simple8b

func unpack240(v uint64, dst *[240]uint64) {
	for i := range dst {
		dst[i] = 1
	}
}

func unpack120(v uint64, dst *[240]uint64) {
	for i := range dst {
		dst[i] = 1
	}
}

func unpack60(v uint64, dst *[240]uint64) {
	dst[0] = v & 1
	dst[1] = (v >> 1) & 1
	dst[2] = (v >> 2) & 1
	dst[3] = (v >> 3) & 1
	dst[4] = (v >> 4) & 1
	dst[5] = (v >> 5) & 1
	dst[6] = (v >> 6) & 1
	dst[7] = (v >> 7) & 1
	dst[8] = (v >> 8) & 1
	dst[9] = (v >> 9) & 1
	dst[10] = (v >> 10) & 1
	dst[11] = (v >> 11) & 1
	dst[12] = (v >> 12) & 1
	dst[13] = (v >> 13) & 1
	dst[14] = (v >> 14) & 1
	dst[15] = (v >> 15) & 1
	dst[16] = (v >> 16) & 1
	dst[17] = (v >> 17) & 1
	dst[18] = (v >> 18) & 1
	dst[19] = (v >> 19) & 1
	dst[20] = (v >> 20) & 1
	dst[21] = (v >> 21) & 1
	dst[22] = (v >> 22) & 1
	dst[23] = (v >> 23) & 1
	dst[24] = (v >> 24) & 1
	dst[25] = (v >> 25) & 1
	dst[26] = (v >> 26) & 1
	dst[27] = (v >> 27) & 1
	dst[28] = (v >> 28) & 1
	dst[29] = (v >> 29) & 1
	dst[30] = (v >> 30) & 1
	dst[31] = (v >> 31) & 1
	dst[32] = (v >> 32) & 1
	dst[33] = (v >> 33) & 1
	dst[34] = (v >> 34) & 1
	dst[35] = (v >> 35) & 1
	dst[36] = (v >> 36) & 1
	dst[37] = (v >> 37) & 1
	dst[38] = (v >> 38) & 1
	dst[39] = (v >> 39) & 1
	dst[40] = (v >> 40) & 1
	dst[41] = (v >> 41) & 1
	dst[42] = (v >> 42) & 1
	dst[43] = (v >> 43) & 1
	dst[44] = (v >> 44) & 1
	dst[45] = (v >> 45) & 1
	dst[46] = (v >> 46) & 1
	dst[47] = (v >> 47) & 1
	dst[48] = (v >> 48) & 1
	dst[49] = (v >> 49) & 1
	dst[50] = (v >> 50) & 1
	dst[51] = (v >> 51) & 1
	dst[52] = (v >> 52) & 1
	dst[53] = (v >> 53) & 1
	dst[54] = (v >> 54) & 1
	dst[55] = (v >> 55) & 1
	dst[56] = (v >> 56) & 1
	dst[57] = (v >> 57) & 1
	dst[58] = (v >> 58) & 1
	dst[59] = (v >> 59) & 1
}

func unpack30(v uint64, dst *[240]uint64) {
	dst[0] = v & 3
	dst[1] = (v >> 2) & 3
	dst[2] = (v >> 4) & 3
	dst[3] = (v >> 6) & 3
	dst[4] = (v >> 8) & 3
	dst[5] = (v >> 10) & 3
	dst[6] = (v >> 12) & 3
	dst[7] = (v >> 14) & 3
	dst[8] = (v >> 16) & 3
	dst[9] = (v >> 18) & 3
	dst[10] = (v >> 20) & 3
	dst[11] = (v >> 22) & 3
	dst[12] = (v >> 24) & 3
	dst[13] = (v >> 26) & 3
	dst[14] = (v >> 28) & 3
	dst[15] = (v >> 30) & 3
	dst[16] = (v >> 32) & 3
	dst[17] = (v >> 34) & 3
	dst[18] = (v >> 36) & 3
	dst[19] = (v >> 38) & 3
	dst[20] = (v >> 40) & 3
	dst[21] = (v >> 42) & 3
	dst[22] = (v >> 44) & 3
	dst[23] = (v >> 46) & 3
	dst[24] = (v >> 48) & 3
	dst[25] = (v >> 50) & 3
	dst[26] = (v >> 52) & 3
	dst[27] = (v >> 54) & 3
	dst[28] = (v >> 56) & 3
	dst[29] = (v >> 58) & 3
}

func unpack20(v uint64, dst *[240]uint64) {
	dst[0] = v & 7
	dst[1] = (v >> 3) & 7
	dst[2] = (v >> 6) & 7
	dst[3] = (v >> 9) & 7
	dst[4] = (v >> 12) & 7
	dst[5] = (v >> 15) & 7
	dst[6] = (v >> 18) & 7
	dst[7] = (v >> 21) & 7
	dst[8] = (v >> 24) & 7
	dst[9] = (v >> 27) & 7
	dst[10] = (v >> 30) & 7
	dst[11] = (v >> 33) & 7
	dst[12] = (v >> 36) & 7
	dst[13] = (v >> 39) & 7
	dst[14] = (v >> 42) & 7
	dst[15] = (v >> 45) & 7
	dst[16] = (v >> 48) & 7
	dst[17] = (v >> 51) & 7
	dst[18] = (v >> 54) & 7
	dst[19] = (v >> 57) & 7
}

func unpack15(v uint64, dst *[240]uint64) {
	dst[0] = v & 15
	dst[1] = (v >> 4) & 15
	dst[2] = (v >> 8) & 15
	dst[3] = (v >> 12) & 15
	dst[4] = (v >> 16) & 15
	dst[5] = (v >> 20) & 15
	dst[6] = (v >> 24) & 15
	dst[7] = (v >> 28) & 15
	dst[8] = (v >> 32) & 15
	dst[9] = (v >> 36) & 15
	dst[10] = (v >> 40) & 15
	dst[11] = (v >> 44) & 15
	dst[12] = (v >> 48) & 15
	dst[13] = (v >> 52) & 15
	dst[14] = (v >> 56) & 15
}

func unpack12(v uint64, dst *[240]uint64) {
	dst[0] = v & 31
	dst[1] = (v >> 5) & 31
	dst[2] = (v >> 10) & 31
	dst[3] = (v >> 15) & 31
	dst[4] = (v >> 20) & 31
	dst[5] = (v >> 25) & 31
	dst[6] = (v >> 30) & 31
	dst[7] = (v >> 35) & 31
	dst[8] = (v >> 40) & 31
	dst[9] = (v >> 45) & 31
	dst[10] = (v >> 50) & 31
	dst[11] = (v >> 55) & 31
}

func unpack10(v uint64, dst *[240]uint64) {
	dst[0] = v & 63
	dst[1] = (v >> 6) & 63
	dst[2] = (v >> 12) & 63
	dst[3] = (v >> 18) & 63
	dst[4] = (v >> 24) & 63
	dst[5] = (v >> 30) & 63
	dst[6] = (v >> 36) & 63
	dst[7] = (v >> 42) & 63
	dst[8] = (v >> 48) & 63
	dst[9] = (v >> 54) & 63
}

func unpack8(v uint64, dst *[240]uint64) {
	dst[0] = v & 127
	dst[1] = (v >> 7) & 127
	dst[2] = (v >> 14) & 127
	dst[3] = (v >> 21) & 127
	dst[4] = (v >> 28) & 127
	dst[5] = (v >> 35) & 127
	dst[6] = (v >> 42) & 127
	dst[7] = (v >> 49) & 127
}

func unpack7(v uint64, dst *[240]uint64) {
	dst[0] = v & 255
	dst[1] = (v >> 8) & 255
	dst[2] = (v >> 16) & 255
	dst[3] = (v >> 24) & 255
	dst[4] = (v >> 32) & 255
	dst[5] = (v >> 40) & 255
	dst[6] = (v >> 48) & 255
}

func unpack6(v uint64, dst *[240]uint64) {
	dst[0] = v & 1023
	dst[1] = (v >> 10) & 1023
	dst[2] = (v >> 20) & 1023
	dst[3] = (v >> 30) & 1023
	dst[4] = (v >> 40) & 1023
	dst[5] = (v >> 50) & 1023
}

func unpack5(v uint64, dst *[240]uint64) {
	dst[0] = v & 4095
	dst[1] = (v >> 12) & 4095
	dst[2] = (v >> 24) & 4095
	dst[3] = (v >> 36) & 4095
	dst[4] = (v >> 48) & 4095
}

func unpack4(v uint64, dst *[240]uint64) {
	dst[0] = v & 32767
	dst[1] = (v >> 15) & 32767
	dst[2] = (v >> 30) & 32767
	dst[3] = (v >> 45) & 32767
}

func unpack3(v uint64, dst *[240]uint64) {
	dst[0] = v & 1048575
	dst[1] = (v >> 20) & 1048575
	dst[2] = (v >> 40) & 1048575
}

func unpack2(v uint64, dst *[240]uint64) {
	dst[0] = v & 1073741823
	dst[1] = (v >> 30) & 1073741823
}

func unpack1(v uint64, dst *[240]uint64) {
	dst[0] = v & 1152921504606846975
}
