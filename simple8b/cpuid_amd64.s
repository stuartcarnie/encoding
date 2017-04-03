#include "textflag.h"

#define cpuid_ecx R8

TEXT ·cpu_info(SB),NOSPLIT,$0
	// find out information about the processor we're on
	MOVQ	$0, AX
	CPUID
	MOVQ	AX, SI
	CMPQ	AX, $0
	JE	done

	// Load EAX=1 cpuid flags
	MOVQ	$1, AX
	CPUID
	MOVL	CX, cpuid_ecx

	// Load EAX=7/ECX=0 cpuid flags
	CMPQ	SI, $7
	XORQ	BX, BX
	JLT	no7
	MOVL	$7, AX
	MOVL	$0, CX
	CPUID
no7:
	// Detect AVX and AVX2 as per 14.7.1  Detection of AVX2 chapter of [1]
	// [1] 64-ia-32-architectures-software-developer-manual-325462.pdf
	// http://www.intel.com/content/dam/www/public/us/en/documents/manuals/64-ia-32-architectures-software-developer-manual-325462.pdf
	ANDL    $0x18000000, cpuid_ecx // check for OSXSAVE and AVX bits
	CMPL    cpuid_ecx, $0x18000000
	JNE     noavx2
	MOVL    $0, CX
	// For XGETBV, OSXSAVE bit is required and sufficient
	XGETBV
	ANDL    $6, AX
	CMPL    AX, $6 // Check for OS support of YMM registers
	JNE     noavx2
	TESTL   $(1<<5), BX // check for AVX2 bit
	JEQ     noavx2
	MOVB    $1, ·support_avx2(SB)
	JMP     done
noavx2:
	MOVB    $0, ·support_avx2(SB)
done:
    RET
