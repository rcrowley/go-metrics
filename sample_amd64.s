// +build !gccgo

#include "textflag.h"

TEXT ·SampleSum(SB), NOSPLIT, $0-32
	CMPB ·x86HasSSE42(SB), $0
	JNE  HasAVX2
	CMPB ·x86HasAVX(SB), $0
	JNE  HasAVX
	CMPB ·x86HasSSE42(SB), $0
	JNE  HasSSE42
	JMP  ·sampleSum(SB)

HasAVX2:
	JMP sampleSumAVX2<>(SB)

HasAVX:
	JMP sampleSumAVX<>(SB)

HasSSE42:
	JMP sampleSumSSE42<>(SB)

TEXT ·SampleVariance(SB), NOSPLIT, $0-32
	CMPB ·x86HasAVX2(SB), $0
	JNE  HasAVX2
	CMPB ·x86HasAVX(SB), $0
	JNE  HasAVX
	CMPB ·x86HasSSE42(SB), $0
	JNE  HasSSE42
	JMP  ·sampleVariance(SB)

HasAVX2:
	JMP sampleVarianceAVX2<>(SB)

HasAVX:
	JMP sampleVarianceAVX<>(SB)

HasSSE42:
	JMP sampleVarianceSSE42<>(SB)

TEXT ·SampleMin(SB), NOSPLIT, $0-32
	CMPB ·x86HasAVX2(SB), $0
	JNE  HasAVX2
	CMPB ·x86HasAVX(SB), $0
	JNE  HasAVX
	CMPB ·x86HasSSE42(SB), $0
	JNE  HasSSE42
	JMP  ·sampleMin(SB)

HasAVX2:
	JMP sampleMinAVX2<>(SB)

HasAVX:
	JMP sampleMinAVX<>(SB)

HasSSE42:
	JMP sampleMinSSE42<>(SB)

TEXT ·SampleMax(SB), NOSPLIT, $0-32
	CMPB ·x86HasAVX2(SB), $0
	JNE  HasAVX2
	CMPB ·x86HasAVX(SB), $0
	JNE  HasAVX
	CMPB ·x86HasSSE42(SB), $0
	JNE  HasSSE42
	JMP  ·sampleMax(SB)

HasAVX2:
	JMP sampleMaxAVX2<>(SB)

HasAVX:
	JMP sampleMaxAVX<>(SB)

HasSSE42:
	JMP sampleMaxSSE42<>(SB)

TEXT sampleSumAVX2<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX

	WORD $0x8548; BYTE $0xf6 // test    rsi, rsi
	JLE  LBB0_1
	LONG $0x10fe8348         // cmp    rsi, 16
	JAE  LBB0_4
	WORD $0xc931             // xor    ecx, ecx
	WORD $0xc031             // xor    eax, eax
	JMP  LBB0_11

LBB0_1:
	WORD $0xc031      // xor    eax, eax
	MOVQ AX, x+24(FP)
	RET

LBB0_4:
	WORD $0x8948; BYTE $0xf1     // mov    rcx, rsi
	LONG $0xf0e18348             // and    rcx, -16
	LONG $0xf0518d48             // lea    rdx, [rcx - 16]
	WORD $0x8948; BYTE $0xd0     // mov    rax, rdx
	LONG $0x04e8c148             // shr    rax, 4
	LONG $0x01c08348             // add    rax, 1
	WORD $0x8941; BYTE $0xc0     // mov    r8d, eax
	LONG $0x01e08341             // and    r8d, 1
	WORD $0x8548; BYTE $0xd2     // test    rdx, rdx
	JE   LBB0_5
	LONG $0x000001ba; BYTE $0x00 // mov    edx, 1
	WORD $0x2948; BYTE $0xc2     // sub    rdx, rax
	LONG $0x10048d49             // lea    rax, [r8 + rdx]
	LONG $0xffc08348             // add    rax, -1
	LONG $0xc0eff9c5             // vpxor    xmm0, xmm0, xmm0
	WORD $0xd231                 // xor    edx, edx
	LONG $0xc9eff1c5             // vpxor    xmm1, xmm1, xmm1
	LONG $0xd2efe9c5             // vpxor    xmm2, xmm2, xmm2
	LONG $0xdbefe1c5             // vpxor    xmm3, xmm3, xmm3

LBB0_7:
	LONG $0x04d4fdc5; BYTE $0xd7         // vpaddq    ymm0, ymm0, yword [rdi + 8*rdx]
	LONG $0x4cd4f5c5; WORD $0x20d7       // vpaddq    ymm1, ymm1, yword [rdi + 8*rdx + 32]
	LONG $0x54d4edc5; WORD $0x40d7       // vpaddq    ymm2, ymm2, yword [rdi + 8*rdx + 64]
	LONG $0x5cd4e5c5; WORD $0x60d7       // vpaddq    ymm3, ymm3, yword [rdi + 8*rdx + 96]
	QUAD $0x000080d784d4fdc5; BYTE $0x00 // vpaddq    ymm0, ymm0, yword [rdi + 8*rdx + 128]
	QUAD $0x0000a0d78cd4f5c5; BYTE $0x00 // vpaddq    ymm1, ymm1, yword [rdi + 8*rdx + 160]
	QUAD $0x0000c0d794d4edc5; BYTE $0x00 // vpaddq    ymm2, ymm2, yword [rdi + 8*rdx + 192]
	QUAD $0x0000e0d79cd4e5c5; BYTE $0x00 // vpaddq    ymm3, ymm3, yword [rdi + 8*rdx + 224]
	LONG $0x20c28348                     // add    rdx, 32
	LONG $0x02c08348                     // add    rax, 2
	JNE  LBB0_7
	WORD $0x854d; BYTE $0xc0             // test    r8, r8
	JE   LBB0_10

LBB0_9:
	LONG $0x5cd4e5c5; WORD $0x60d7 // vpaddq    ymm3, ymm3, yword [rdi + 8*rdx + 96]
	LONG $0x54d4edc5; WORD $0x40d7 // vpaddq    ymm2, ymm2, yword [rdi + 8*rdx + 64]
	LONG $0x4cd4f5c5; WORD $0x20d7 // vpaddq    ymm1, ymm1, yword [rdi + 8*rdx + 32]
	LONG $0x04d4fdc5; BYTE $0xd7   // vpaddq    ymm0, ymm0, yword [rdi + 8*rdx]

LBB0_10:
	LONG $0xcbd4f5c5               // vpaddq    ymm1, ymm1, ymm3
	LONG $0xc2d4fdc5               // vpaddq    ymm0, ymm0, ymm2
	LONG $0xc1d4fdc5               // vpaddq    ymm0, ymm0, ymm1
	LONG $0x397de3c4; WORD $0x01c1 // vextracti128    xmm1, ymm0, 1
	LONG $0xc1d4fdc5               // vpaddq    ymm0, ymm0, ymm1
	LONG $0xc870f9c5; BYTE $0x4e   // vpshufd    xmm1, xmm0, 78
	LONG $0xc1d4fdc5               // vpaddq    ymm0, ymm0, ymm1
	LONG $0x7ef9e1c4; BYTE $0xc0   // vmovq    rax, xmm0
	WORD $0x3948; BYTE $0xf1       // cmp    rcx, rsi
	JE   LBB0_12

LBB0_11:
	LONG $0xcf040348         // add    rax, qword [rdi + 8*rcx]
	LONG $0x01c18348         // add    rcx, 1
	WORD $0x3948; BYTE $0xce // cmp    rsi, rcx
	JNE  LBB0_11

LBB0_12:
	BYTE $0xc5; BYTE $0xf8; BYTE $0x77 // VZEROUPPER
	MOVQ AX, x+24(FP)
	RET

LBB0_5:
	LONG $0xc0eff9c5         // vpxor    xmm0, xmm0, xmm0
	WORD $0xd231             // xor    edx, edx
	LONG $0xc9eff1c5         // vpxor    xmm1, xmm1, xmm1
	LONG $0xd2efe9c5         // vpxor    xmm2, xmm2, xmm2
	LONG $0xdbefe1c5         // vpxor    xmm3, xmm3, xmm3
	WORD $0x854d; BYTE $0xc0 // test    r8, r8
	JNE  LBB0_9
	JMP  LBB0_10

TEXT sampleSumAVX<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX

	WORD $0x8548; BYTE $0xf6 // test    rsi, rsi
	JLE  LBB0_1
	LONG $0x10fe8348         // cmp    rsi, 16
	JAE  LBB0_4
	WORD $0xc931             // xor    ecx, ecx
	WORD $0xc031             // xor    eax, eax
	JMP  LBB0_11

LBB0_1:
	WORD $0xc031      // xor    eax, eax
	MOVQ AX, x+24(FP)
	RET

LBB0_4:
	WORD $0x8948; BYTE $0xf1     // mov    rcx, rsi
	LONG $0xf0e18348             // and    rcx, -16
	LONG $0xf0518d48             // lea    rdx, [rcx - 16]
	WORD $0x8948; BYTE $0xd0     // mov    rax, rdx
	LONG $0x04e8c148             // shr    rax, 4
	LONG $0x01c08348             // add    rax, 1
	WORD $0x8941; BYTE $0xc0     // mov    r8d, eax
	LONG $0x01e08341             // and    r8d, 1
	WORD $0x8548; BYTE $0xd2     // test    rdx, rdx
	JE   LBB0_5
	LONG $0x000001ba; BYTE $0x00 // mov    edx, 1
	WORD $0x2948; BYTE $0xc2     // sub    rdx, rax
	LONG $0x10048d49             // lea    rax, [r8 + rdx]
	LONG $0xffc08348             // add    rax, -1
	LONG $0xef3941c4; BYTE $0xc0 // vpxor    xmm8, xmm8, xmm8
	WORD $0xd231                 // xor    edx, edx
	LONG $0xdbefe1c5             // vpxor    xmm3, xmm3, xmm3
	LONG $0xef3141c4; BYTE $0xc9 // vpxor    xmm9, xmm9, xmm9
	LONG $0xef2941c4; BYTE $0xd2 // vpxor    xmm10, xmm10, xmm10

LBB0_7:
	LONG $0x246ffec5; BYTE $0xd7         // vmovdqu    ymm4, yword [rdi + 8*rdx]
	LONG $0x6c6ffec5; WORD $0x20d7       // vmovdqu    ymm5, yword [rdi + 8*rdx + 32]
	LONG $0x746ffec5; WORD $0x40d7       // vmovdqu    ymm6, yword [rdi + 8*rdx + 64]
	LONG $0x7c6ffec5; WORD $0x60d7       // vmovdqu    ymm7, yword [rdi + 8*rdx + 96]
	LONG $0xd45941c4; BYTE $0xd8         // vpaddq    xmm11, xmm4, xmm8
	LONG $0x197de3c4; WORD $0x01e4       // vextractf128    xmm4, ymm4, 1
	LONG $0x197d63c4; WORD $0x01c1       // vextractf128    xmm1, ymm8, 1
	LONG $0xc9d4d9c5                     // vpaddq    xmm1, xmm4, xmm1
	LONG $0xebd451c5                     // vpaddq    xmm13, xmm5, xmm3
	LONG $0x197de3c4; WORD $0x01ed       // vextractf128    xmm5, ymm5, 1
	LONG $0x197de3c4; WORD $0x01db       // vextractf128    xmm3, ymm3, 1
	LONG $0xdbd4d1c5                     // vpaddq    xmm3, xmm5, xmm3
	LONG $0xd449c1c4; BYTE $0xe9         // vpaddq    xmm5, xmm6, xmm9
	LONG $0x197de3c4; WORD $0x01f6       // vextractf128    xmm6, ymm6, 1
	LONG $0x197d63c4; WORD $0x01ca       // vextractf128    xmm2, ymm9, 1
	LONG $0xd2d4c9c5                     // vpaddq    xmm2, xmm6, xmm2
	LONG $0xd441c1c4; BYTE $0xf2         // vpaddq    xmm6, xmm7, xmm10
	LONG $0x197de3c4; WORD $0x01ff       // vextractf128    xmm7, ymm7, 1
	LONG $0x197d63c4; WORD $0x01d0       // vextractf128    xmm0, ymm10, 1
	LONG $0xc0d4c1c5                     // vpaddq    xmm0, xmm7, xmm0
	QUAD $0x000080d7bc6ffec5; BYTE $0x00 // vmovdqu    ymm7, yword [rdi + 8*rdx + 128]
	QUAD $0x0000a0d78c6f7ec5; BYTE $0x00 // vmovdqu    ymm9, yword [rdi + 8*rdx + 160]
	QUAD $0x0000c0d7946f7ec5; BYTE $0x00 // vmovdqu    ymm10, yword [rdi + 8*rdx + 192]
	QUAD $0x0000e0d7a46f7ec5; BYTE $0x00 // vmovdqu    ymm12, yword [rdi + 8*rdx + 224]
	LONG $0x197de3c4; WORD $0x01fc       // vextractf128    xmm4, ymm7, 1
	LONG $0xc9d4d9c5                     // vpaddq    xmm1, xmm4, xmm1
	LONG $0xd441c1c4; BYTE $0xe3         // vpaddq    xmm4, xmm7, xmm11
	LONG $0x185d63c4; WORD $0x01c1       // vinsertf128    ymm8, ymm4, xmm1, 1
	LONG $0x197d63c4; WORD $0x01c9       // vextractf128    xmm1, ymm9, 1
	LONG $0xcbd4f1c5                     // vpaddq    xmm1, xmm1, xmm3
	LONG $0xd431c1c4; BYTE $0xdd         // vpaddq    xmm3, xmm9, xmm13
	LONG $0x1865e3c4; WORD $0x01d9       // vinsertf128    ymm3, ymm3, xmm1, 1
	LONG $0x197d63c4; WORD $0x01d1       // vextractf128    xmm1, ymm10, 1
	LONG $0xcad4f1c5                     // vpaddq    xmm1, xmm1, xmm2
	LONG $0xd5d4a9c5                     // vpaddq    xmm2, xmm10, xmm5
	LONG $0x186d63c4; WORD $0x01c9       // vinsertf128    ymm9, ymm2, xmm1, 1
	LONG $0x197d63c4; WORD $0x01e1       // vextractf128    xmm1, ymm12, 1
	LONG $0xc0d4f1c5                     // vpaddq    xmm0, xmm1, xmm0
	LONG $0xced499c5                     // vpaddq    xmm1, xmm12, xmm6
	LONG $0x187563c4; WORD $0x01d0       // vinsertf128    ymm10, ymm1, xmm0, 1
	LONG $0x20c28348                     // add    rdx, 32
	LONG $0x02c08348                     // add    rax, 2
	JNE  LBB0_7
	WORD $0x854d; BYTE $0xc0             // test    r8, r8
	JE   LBB0_10

LBB0_9:
	LONG $0x246ffec5; BYTE $0xd7   // vmovdqu    ymm4, yword [rdi + 8*rdx]
	LONG $0x446ffec5; WORD $0x20d7 // vmovdqu    ymm0, yword [rdi + 8*rdx + 32]
	LONG $0x4c6ffec5; WORD $0x40d7 // vmovdqu    ymm1, yword [rdi + 8*rdx + 64]
	LONG $0x546ffec5; WORD $0x60d7 // vmovdqu    ymm2, yword [rdi + 8*rdx + 96]
	LONG $0x197de3c4; WORD $0x01d5 // vextractf128    xmm5, ymm2, 1
	LONG $0x197d63c4; WORD $0x01d6 // vextractf128    xmm6, ymm10, 1
	LONG $0xeed4d1c5               // vpaddq    xmm5, xmm5, xmm6
	LONG $0xd469c1c4; BYTE $0xd2   // vpaddq    xmm2, xmm2, xmm10
	LONG $0x186d63c4; WORD $0x01d5 // vinsertf128    ymm10, ymm2, xmm5, 1
	LONG $0x197de3c4; WORD $0x01ca // vextractf128    xmm2, ymm1, 1
	LONG $0x197d63c4; WORD $0x01cd // vextractf128    xmm5, ymm9, 1
	LONG $0xd5d4e9c5               // vpaddq    xmm2, xmm2, xmm5
	LONG $0xd471c1c4; BYTE $0xc9   // vpaddq    xmm1, xmm1, xmm9
	LONG $0x187563c4; WORD $0x01ca // vinsertf128    ymm9, ymm1, xmm2, 1
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197de3c4; WORD $0x01da // vextractf128    xmm2, ymm3, 1
	LONG $0xcad4f1c5               // vpaddq    xmm1, xmm1, xmm2
	LONG $0xc3d4f9c5               // vpaddq    xmm0, xmm0, xmm3
	LONG $0x187de3c4; WORD $0x01d9 // vinsertf128    ymm3, ymm0, xmm1, 1
	LONG $0x197de3c4; WORD $0x01e0 // vextractf128    xmm0, ymm4, 1
	LONG $0x197d63c4; WORD $0x01c1 // vextractf128    xmm1, ymm8, 1
	LONG $0xc1d4f9c5               // vpaddq    xmm0, xmm0, xmm1
	LONG $0xd459c1c4; BYTE $0xc8   // vpaddq    xmm1, xmm4, xmm8
	LONG $0x187563c4; WORD $0x01c0 // vinsertf128    ymm8, ymm1, xmm0, 1

LBB0_10:
	LONG $0x197d63c4; WORD $0x01c0 // vextractf128    xmm0, ymm8, 1
	LONG $0x197de3c4; WORD $0x01d9 // vextractf128    xmm1, ymm3, 1
	LONG $0xc0d4f1c5               // vpaddq    xmm0, xmm1, xmm0
	LONG $0xd461c1c4; BYTE $0xc8   // vpaddq    xmm1, xmm3, xmm8
	LONG $0x197d63c4; WORD $0x01ca // vextractf128    xmm2, ymm9, 1
	LONG $0x197d63c4; WORD $0x01d3 // vextractf128    xmm3, ymm10, 1
	LONG $0xd3d4e9c5               // vpaddq    xmm2, xmm2, xmm3
	LONG $0xc2d4f9c5               // vpaddq    xmm0, xmm0, xmm2
	LONG $0xd431c1c4; BYTE $0xd2   // vpaddq    xmm2, xmm9, xmm10
	LONG $0xcad4f1c5               // vpaddq    xmm1, xmm1, xmm2
	LONG $0xc0d4f1c5               // vpaddq    xmm0, xmm1, xmm0
	LONG $0xc870f9c5; BYTE $0x4e   // vpshufd    xmm1, xmm0, 78
	LONG $0xc1d4f9c5               // vpaddq    xmm0, xmm0, xmm1
	LONG $0x7ef9e1c4; BYTE $0xc0   // vmovq    rax, xmm0
	WORD $0x3948; BYTE $0xf1       // cmp    rcx, rsi
	JE   LBB0_12

LBB0_11:
	LONG $0xcf040348         // add    rax, qword [rdi + 8*rcx]
	LONG $0x01c18348         // add    rcx, 1
	WORD $0x3948; BYTE $0xce // cmp    rsi, rcx
	JNE  LBB0_11

LBB0_12:
	BYTE $0xc5; BYTE $0xf8; BYTE $0x77 // VZEROUPPER
	MOVQ AX, x+24(FP)
	RET

LBB0_5:
	LONG $0xef3941c4; BYTE $0xc0 // vpxor    xmm8, xmm8, xmm8
	WORD $0xd231                 // xor    edx, edx
	LONG $0xdbefe1c5             // vpxor    xmm3, xmm3, xmm3
	LONG $0xef3141c4; BYTE $0xc9 // vpxor    xmm9, xmm9, xmm9
	LONG $0xef2941c4; BYTE $0xd2 // vpxor    xmm10, xmm10, xmm10
	WORD $0x854d; BYTE $0xc0     // test    r8, r8
	JNE  LBB0_9
	JMP  LBB0_10

TEXT sampleSumSSE42<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX

	WORD $0x8548; BYTE $0xf6 // test    rsi, rsi
	JLE  LBB0_1
	LONG $0x04fe8348         // cmp    rsi, 4
	JAE  LBB0_4
	WORD $0xc931             // xor    ecx, ecx
	WORD $0xc031             // xor    eax, eax
	JMP  LBB0_12

LBB0_1:
	WORD $0xc031 // xor    eax, eax
	JMP  LBB0_13

LBB0_4:
	WORD $0x8948; BYTE $0xf1 // mov    rcx, rsi
	LONG $0xfce18348         // and    rcx, -4
	LONG $0xfc518d48         // lea    rdx, [rcx - 4]
	WORD $0x8948; BYTE $0xd0 // mov    rax, rdx
	LONG $0x02e8c148         // shr    rax, 2
	LONG $0x01c08348         // add    rax, 1
	WORD $0x8941; BYTE $0xc0 // mov    r8d, eax
	LONG $0x03e08341         // and    r8d, 3
	LONG $0x0cfa8348         // cmp    rdx, 12
	JAE  LBB0_6
	LONG $0xc0ef0f66         // pxor    xmm0, xmm0
	WORD $0xd231             // xor    edx, edx
	LONG $0xc9ef0f66         // pxor    xmm1, xmm1
	WORD $0x854d; BYTE $0xc0 // test    r8, r8
	JNE  LBB0_9
	JMP  LBB0_11

LBB0_6:
	LONG $0x000001ba; BYTE $0x00 // mov    edx, 1
	WORD $0x2948; BYTE $0xc2     // sub    rdx, rax
	LONG $0x10048d49             // lea    rax, [r8 + rdx]
	LONG $0xffc08348             // add    rax, -1
	LONG $0xc0ef0f66             // pxor    xmm0, xmm0
	WORD $0xd231                 // xor    edx, edx
	LONG $0xc9ef0f66             // pxor    xmm1, xmm1

LBB0_7:
	LONG $0x146f0ff3; BYTE $0xd7   // movdqu    xmm2, oword [rdi + 8*rdx]
	LONG $0xd0d40f66               // paddq    xmm2, xmm0
	LONG $0x446f0ff3; WORD $0x10d7 // movdqu    xmm0, oword [rdi + 8*rdx + 16]
	LONG $0xc1d40f66               // paddq    xmm0, xmm1
	LONG $0x4c6f0ff3; WORD $0x20d7 // movdqu    xmm1, oword [rdi + 8*rdx + 32]
	LONG $0x5c6f0ff3; WORD $0x30d7 // movdqu    xmm3, oword [rdi + 8*rdx + 48]
	LONG $0x646f0ff3; WORD $0x40d7 // movdqu    xmm4, oword [rdi + 8*rdx + 64]
	LONG $0xe1d40f66               // paddq    xmm4, xmm1
	LONG $0xe2d40f66               // paddq    xmm4, xmm2
	LONG $0x546f0ff3; WORD $0x50d7 // movdqu    xmm2, oword [rdi + 8*rdx + 80]
	LONG $0xd3d40f66               // paddq    xmm2, xmm3
	LONG $0xd0d40f66               // paddq    xmm2, xmm0
	LONG $0x446f0ff3; WORD $0x60d7 // movdqu    xmm0, oword [rdi + 8*rdx + 96]
	LONG $0xc4d40f66               // paddq    xmm0, xmm4
	LONG $0x4c6f0ff3; WORD $0x70d7 // movdqu    xmm1, oword [rdi + 8*rdx + 112]
	LONG $0xcad40f66               // paddq    xmm1, xmm2
	LONG $0x10c28348               // add    rdx, 16
	LONG $0x04c08348               // add    rax, 4
	JNE  LBB0_7
	WORD $0x854d; BYTE $0xc0       // test    r8, r8
	JE   LBB0_11

LBB0_9:
	LONG $0xd7048d48         // lea    rax, [rdi + 8*rdx]
	LONG $0x10c08348         // add    rax, 16
	WORD $0xf749; BYTE $0xd8 // neg    r8

LBB0_10:
	LONG $0x506f0ff3; BYTE $0xf0 // movdqu    xmm2, oword [rax - 16]
	LONG $0xc2d40f66             // paddq    xmm0, xmm2
	LONG $0x106f0ff3             // movdqu    xmm2, oword [rax]
	LONG $0xcad40f66             // paddq    xmm1, xmm2
	LONG $0x20c08348             // add    rax, 32
	LONG $0x01c08349             // add    r8, 1
	JNE  LBB0_10

LBB0_11:
	LONG $0xc1d40f66             // paddq    xmm0, xmm1
	LONG $0xc8700f66; BYTE $0x4e // pshufd    xmm1, xmm0, 78
	LONG $0xc8d40f66             // paddq    xmm1, xmm0
	LONG $0x7e0f4866; BYTE $0xc8 // movq    rax, xmm1
	WORD $0x3948; BYTE $0xf1     // cmp    rcx, rsi
	JE   LBB0_13

LBB0_12:
	LONG $0xcf040348         // add    rax, qword [rdi + 8*rcx]
	LONG $0x01c18348         // add    rcx, 1
	WORD $0x3948; BYTE $0xce // cmp    rsi, rcx
	JNE  LBB0_12

LBB0_13:
	MOVQ AX, x+24(FP)
	RET

TEXT sampleVarianceAVX2<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX

	WORD $0x8548; BYTE $0xf6 // test    rsi, rsi
	JE   LBB3_1
	JLE  LBB3_23
	LONG $0x0ffe8348         // cmp    rsi, 15
	JA   LBB3_5
	WORD $0xc031             // xor    eax, eax
	WORD $0xc931             // xor    ecx, ecx
	JMP  LBB3_12

LBB3_1:
	LONG $0xc057f8c5  // vxorps    xmm0, xmm0, xmm0
	MOVQ X0, x+24(FP)
	RET

LBB3_23:
	LONG $0x2afb61c4; BYTE $0xde // vcvtsi2sd    xmm11, xmm0, rsi
	LONG $0xd257e9c5             // vxorpd    xmm2, xmm2, xmm2
	LONG $0x5e6bc1c4; BYTE $0xc3 // vdivsd    xmm0, xmm2, xmm11
	JMP  LBB3_22

LBB3_5:
	WORD $0x8948; BYTE $0xf0     // mov    rax, rsi
	LONG $0xf0e08348             // and    rax, -16
	LONG $0xf0508d48             // lea    rdx, [rax - 16]
	WORD $0x8948; BYTE $0xd1     // mov    rcx, rdx
	LONG $0x04e9c148             // shr    rcx, 4
	LONG $0x01c18348             // add    rcx, 1
	WORD $0x8941; BYTE $0xc8     // mov    r8d, ecx
	LONG $0x01e08341             // and    r8d, 1
	WORD $0x8548; BYTE $0xd2     // test    rdx, rdx
	JE   LBB3_6
	LONG $0x000001ba; BYTE $0x00 // mov    edx, 1
	WORD $0x2948; BYTE $0xca     // sub    rdx, rcx
	LONG $0x100c8d49             // lea    rcx, [r8 + rdx]
	LONG $0xffc18348             // add    rcx, -1
	LONG $0xc0eff9c5             // vpxor    xmm0, xmm0, xmm0
	WORD $0xd231                 // xor    edx, edx
	LONG $0xc9eff1c5             // vpxor    xmm1, xmm1, xmm1
	LONG $0xd2efe9c5             // vpxor    xmm2, xmm2, xmm2
	LONG $0xdbefe1c5             // vpxor    xmm3, xmm3, xmm3

LBB3_8:
	LONG $0x04d4fdc5; BYTE $0xd7         // vpaddq    ymm0, ymm0, yword [rdi + 8*rdx]
	LONG $0x4cd4f5c5; WORD $0x20d7       // vpaddq    ymm1, ymm1, yword [rdi + 8*rdx + 32]
	LONG $0x54d4edc5; WORD $0x40d7       // vpaddq    ymm2, ymm2, yword [rdi + 8*rdx + 64]
	LONG $0x5cd4e5c5; WORD $0x60d7       // vpaddq    ymm3, ymm3, yword [rdi + 8*rdx + 96]
	QUAD $0x000080d784d4fdc5; BYTE $0x00 // vpaddq    ymm0, ymm0, yword [rdi + 8*rdx + 128]
	QUAD $0x0000a0d78cd4f5c5; BYTE $0x00 // vpaddq    ymm1, ymm1, yword [rdi + 8*rdx + 160]
	QUAD $0x0000c0d794d4edc5; BYTE $0x00 // vpaddq    ymm2, ymm2, yword [rdi + 8*rdx + 192]
	QUAD $0x0000e0d79cd4e5c5; BYTE $0x00 // vpaddq    ymm3, ymm3, yword [rdi + 8*rdx + 224]
	LONG $0x20c28348                     // add    rdx, 32
	LONG $0x02c18348                     // add    rcx, 2
	JNE  LBB3_8
	WORD $0x854d; BYTE $0xc0             // test    r8, r8
	JE   LBB3_11

LBB3_10:
	LONG $0x5cd4e5c5; WORD $0x60d7 // vpaddq    ymm3, ymm3, yword [rdi + 8*rdx + 96]
	LONG $0x54d4edc5; WORD $0x40d7 // vpaddq    ymm2, ymm2, yword [rdi + 8*rdx + 64]
	LONG $0x4cd4f5c5; WORD $0x20d7 // vpaddq    ymm1, ymm1, yword [rdi + 8*rdx + 32]
	LONG $0x04d4fdc5; BYTE $0xd7   // vpaddq    ymm0, ymm0, yword [rdi + 8*rdx]

LBB3_11:
	LONG $0xcbd4f5c5               // vpaddq    ymm1, ymm1, ymm3
	LONG $0xc2d4fdc5               // vpaddq    ymm0, ymm0, ymm2
	LONG $0xc1d4fdc5               // vpaddq    ymm0, ymm0, ymm1
	LONG $0x397de3c4; WORD $0x01c1 // vextracti128    xmm1, ymm0, 1
	LONG $0xc1d4fdc5               // vpaddq    ymm0, ymm0, ymm1
	LONG $0xc870f9c5; BYTE $0x4e   // vpshufd    xmm1, xmm0, 78
	LONG $0xc1d4fdc5               // vpaddq    ymm0, ymm0, ymm1
	LONG $0x7ef9e1c4; BYTE $0xc1   // vmovq    rcx, xmm0
	WORD $0x3948; BYTE $0xf0       // cmp    rax, rsi
	JE   LBB3_13

LBB3_12:
	LONG $0xc70c0348         // add    rcx, qword [rdi + 8*rax]
	LONG $0x01c08348         // add    rax, 1
	WORD $0x3948; BYTE $0xc6 // cmp    rsi, rax
	JNE  LBB3_12

LBB3_13:
	LONG $0x2adb61c4; BYTE $0xde // vcvtsi2sd    xmm11, xmm4, rsi
	WORD $0x8548; BYTE $0xf6     // test    rsi, rsi
	JLE  LBB3_14
	LONG $0x2adbe1c4; BYTE $0xc9 // vcvtsi2sd    xmm1, xmm4, rcx
	LONG $0x5e7341c4; BYTE $0xe3 // vdivsd    xmm12, xmm1, xmm11
	LONG $0x10fe8348             // cmp    rsi, 16
	JAE  LBB3_17
	LONG $0xd2efe9c5             // vpxor    xmm2, xmm2, xmm2
	WORD $0xc031                 // xor    eax, eax
	JMP  LBB3_20

LBB3_14:
	LONG $0xd2efe9c5             // vpxor    xmm2, xmm2, xmm2
	LONG $0x5e6bc1c4; BYTE $0xc3 // vdivsd    xmm0, xmm2, xmm11
	JMP  LBB3_22

LBB3_17:
	WORD $0x8948; BYTE $0xf0     // mov    rax, rsi
	LONG $0xf0e08348             // and    rax, -16
	LONG $0x197d42c4; BYTE $0xec // vbroadcastsd    ymm13, xmm12
	LONG $0x570941c4; BYTE $0xf6 // vxorpd    xmm14, xmm14, xmm14
	WORD $0xc931                 // xor    ecx, ecx
	LONG $0xe457d9c5             // vxorpd    xmm4, xmm4, xmm4
	LONG $0xed57d1c5             // vxorpd    xmm5, xmm5, xmm5
	LONG $0xf657c9c5             // vxorpd    xmm6, xmm6, xmm6

LBB3_18:
	LONG $0x046f7ec5; BYTE $0xcf   // vmovdqu    ymm8, yword [rdi + 8*rcx]
	LONG $0x546f7ec5; WORD $0x20cf // vmovdqu    ymm10, yword [rdi + 8*rcx + 32]
	LONG $0x4c6f7ec5; WORD $0x40cf // vmovdqu    ymm9, yword [rdi + 8*rcx + 64]
	LONG $0x397d63c4; WORD $0x01c0 // vextracti128    xmm0, ymm8, 1
	LONG $0x16f9e3c4; WORD $0x01c2 // vpextrq    rdx, xmm0, 1
	LONG $0x7c6ffec5; WORD $0x60cf // vmovdqu    ymm7, yword [rdi + 8*rcx + 96]
	LONG $0x2a83e1c4; BYTE $0xca   // vcvtsi2sd    xmm1, xmm15, rdx
	LONG $0x7ef9e1c4; BYTE $0xc2   // vmovq    rdx, xmm0
	LONG $0x2a83e1c4; BYTE $0xc2   // vcvtsi2sd    xmm0, xmm15, rdx
	LONG $0x16f963c4; WORD $0x01c2 // vpextrq    rdx, xmm8, 1
	LONG $0xc114f9c5               // vunpcklpd    xmm0, xmm0, xmm1
	LONG $0x2a83e1c4; BYTE $0xca   // vcvtsi2sd    xmm1, xmm15, rdx
	LONG $0x7ef961c4; BYTE $0xc2   // vmovq    rdx, xmm8
	LONG $0x2a83e1c4; BYTE $0xda   // vcvtsi2sd    xmm3, xmm15, rdx
	LONG $0xc914e1c5               // vunpcklpd    xmm1, xmm3, xmm1
	LONG $0x397d63c4; WORD $0x01d3 // vextracti128    xmm3, ymm10, 1
	LONG $0x16f9e3c4; WORD $0x01da // vpextrq    rdx, xmm3, 1
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0x187563c4; WORD $0x01c0 // vinsertf128    ymm8, ymm1, xmm0, 1
	LONG $0x7ef9e1c4; BYTE $0xda   // vmovq    rdx, xmm3
	LONG $0x2a83e1c4; BYTE $0xc2   // vcvtsi2sd    xmm0, xmm15, rdx
	LONG $0x16f963c4; WORD $0x01d2 // vpextrq    rdx, xmm10, 1
	LONG $0x2a83e1c4; BYTE $0xca   // vcvtsi2sd    xmm1, xmm15, rdx
	LONG $0xc214f9c5               // vunpcklpd    xmm0, xmm0, xmm2
	LONG $0x7ef961c4; BYTE $0xd2   // vmovq    rdx, xmm10
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0xc914e9c5               // vunpcklpd    xmm1, xmm2, xmm1
	LONG $0x397d63c4; WORD $0x01ca // vextracti128    xmm2, ymm9, 1
	LONG $0x16f9e3c4; WORD $0x01d2 // vpextrq    rdx, xmm2, 1
	LONG $0x187563c4; WORD $0x01d0 // vinsertf128    ymm10, ymm1, xmm0, 1
	LONG $0x2a83e1c4; BYTE $0xca   // vcvtsi2sd    xmm1, xmm15, rdx
	LONG $0x7ef9e1c4; BYTE $0xd2   // vmovq    rdx, xmm2
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0x16f963c4; WORD $0x01ca // vpextrq    rdx, xmm9, 1
	LONG $0xc914e9c5               // vunpcklpd    xmm1, xmm2, xmm1
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0x7ef961c4; BYTE $0xca   // vmovq    rdx, xmm9
	LONG $0x2a83e1c4; BYTE $0xda   // vcvtsi2sd    xmm3, xmm15, rdx
	LONG $0xd214e1c5               // vunpcklpd    xmm2, xmm3, xmm2
	LONG $0x397de3c4; WORD $0x01fb // vextracti128    xmm3, ymm7, 1
	LONG $0x16f9e3c4; WORD $0x01da // vpextrq    rdx, xmm3, 1
	LONG $0x2a83e1c4; BYTE $0xc2   // vcvtsi2sd    xmm0, xmm15, rdx
	LONG $0x186de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm2, xmm1, 1
	LONG $0x7ef9e1c4; BYTE $0xda   // vmovq    rdx, xmm3
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0x16f9e3c4; WORD $0x01fa // vpextrq    rdx, xmm7, 1
	LONG $0x2a83e1c4; BYTE $0xda   // vcvtsi2sd    xmm3, xmm15, rdx
	LONG $0xc014e9c5               // vunpcklpd    xmm0, xmm2, xmm0
	LONG $0x7ef9e1c4; BYTE $0xfa   // vmovq    rdx, xmm7
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0xd314e9c5               // vunpcklpd    xmm2, xmm2, xmm3
	LONG $0x186de3c4; WORD $0x01c0 // vinsertf128    ymm0, ymm2, xmm0, 1
	LONG $0x5c3dc1c4; BYTE $0xd5   // vsubpd    ymm2, ymm8, ymm13
	LONG $0x5c2dc1c4; BYTE $0xdd   // vsubpd    ymm3, ymm10, ymm13
	LONG $0x5c75c1c4; BYTE $0xcd   // vsubpd    ymm1, ymm1, ymm13
	LONG $0x5c7dc1c4; BYTE $0xc5   // vsubpd    ymm0, ymm0, ymm13
	LONG $0xd259edc5               // vmulpd    ymm2, ymm2, ymm2
	LONG $0x586d41c4; BYTE $0xf6   // vaddpd    ymm14, ymm2, ymm14
	LONG $0xd359e5c5               // vmulpd    ymm2, ymm3, ymm3
	LONG $0xe458edc5               // vaddpd    ymm4, ymm2, ymm4
	LONG $0xc959f5c5               // vmulpd    ymm1, ymm1, ymm1
	LONG $0xed58f5c5               // vaddpd    ymm5, ymm1, ymm5
	LONG $0xc059fdc5               // vmulpd    ymm0, ymm0, ymm0
	LONG $0xf658fdc5               // vaddpd    ymm6, ymm0, ymm6
	LONG $0x10c18348               // add    rcx, 16
	WORD $0x3948; BYTE $0xc8       // cmp    rax, rcx
	JNE  LBB3_18
	LONG $0x585dc1c4; BYTE $0xc6   // vaddpd    ymm0, ymm4, ymm14
	LONG $0xc058d5c5               // vaddpd    ymm0, ymm5, ymm0
	LONG $0xc058cdc5               // vaddpd    ymm0, ymm6, ymm0
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0xc158fdc5               // vaddpd    ymm0, ymm0, ymm1
	LONG $0xd07cfdc5               // vhaddpd    ymm2, ymm0, ymm0
	WORD $0x3948; BYTE $0xf0       // cmp    rax, rsi
	JE   LBB3_21

LBB3_20:
	LONG $0x2a83e1c4; WORD $0xc704 // vcvtsi2sd    xmm0, xmm15, qword [rdi + 8*rax]
	LONG $0x5c7bc1c4; BYTE $0xc4   // vsubsd    xmm0, xmm0, xmm12
	LONG $0xc059fbc5               // vmulsd    xmm0, xmm0, xmm0
	LONG $0xd258fbc5               // vaddsd    xmm2, xmm0, xmm2
	LONG $0x01c08348               // add    rax, 1
	WORD $0x3948; BYTE $0xc6       // cmp    rsi, rax
	JNE  LBB3_20

LBB3_21:
	LONG $0x5e6bc1c4; BYTE $0xc3 // vdivsd    xmm0, xmm2, xmm11

LBB3_22:
	BYTE $0xc5; BYTE $0xf8; BYTE $0x77 // VZEROUPPER
	MOVQ X0, x+24(FP)
	RET

LBB3_6:
	LONG $0xc0eff9c5         // vpxor    xmm0, xmm0, xmm0
	WORD $0xd231             // xor    edx, edx
	LONG $0xc9eff1c5         // vpxor    xmm1, xmm1, xmm1
	LONG $0xd2efe9c5         // vpxor    xmm2, xmm2, xmm2
	LONG $0xdbefe1c5         // vpxor    xmm3, xmm3, xmm3
	WORD $0x854d; BYTE $0xc0 // test    r8, r8
	JNE  LBB3_10
	JMP  LBB3_11

TEXT sampleVarianceAVX<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX

	WORD $0x8548; BYTE $0xf6 // test    rsi, rsi
	JE   LBB3_1
	JLE  LBB3_23
	LONG $0x0ffe8348         // cmp    rsi, 15
	JA   LBB3_5
	WORD $0xc031             // xor    eax, eax
	WORD $0xc931             // xor    ecx, ecx
	JMP  LBB3_12

LBB3_1:
	LONG $0xc057f8c5  // vxorps    xmm0, xmm0, xmm0
	MOVQ X0, x+24(FP)
	RET

LBB3_23:
	LONG $0x2afb61c4; BYTE $0xde // vcvtsi2sd    xmm11, xmm0, rsi
	LONG $0xd257e9c5             // vxorpd    xmm2, xmm2, xmm2
	LONG $0x5e6bc1c4; BYTE $0xc3 // vdivsd    xmm0, xmm2, xmm11
	JMP  LBB3_22

LBB3_5:
	WORD $0x8948; BYTE $0xf0     // mov    rax, rsi
	LONG $0xf0e08348             // and    rax, -16
	LONG $0xf0508d48             // lea    rdx, [rax - 16]
	WORD $0x8948; BYTE $0xd1     // mov    rcx, rdx
	LONG $0x04e9c148             // shr    rcx, 4
	LONG $0x01c18348             // add    rcx, 1
	WORD $0x8941; BYTE $0xc8     // mov    r8d, ecx
	LONG $0x01e08341             // and    r8d, 1
	WORD $0x8548; BYTE $0xd2     // test    rdx, rdx
	JE   LBB3_6
	LONG $0x000001ba; BYTE $0x00 // mov    edx, 1
	WORD $0x2948; BYTE $0xca     // sub    rdx, rcx
	LONG $0x100c8d49             // lea    rcx, [r8 + rdx]
	LONG $0xffc18348             // add    rcx, -1
	LONG $0xef3941c4; BYTE $0xc0 // vpxor    xmm8, xmm8, xmm8
	WORD $0xd231                 // xor    edx, edx
	LONG $0xdbefe1c5             // vpxor    xmm3, xmm3, xmm3
	LONG $0xef3141c4; BYTE $0xc9 // vpxor    xmm9, xmm9, xmm9
	LONG $0xef2941c4; BYTE $0xd2 // vpxor    xmm10, xmm10, xmm10

LBB3_8:
	LONG $0x246ffec5; BYTE $0xd7         // vmovdqu    ymm4, yword [rdi + 8*rdx]
	LONG $0x6c6ffec5; WORD $0x20d7       // vmovdqu    ymm5, yword [rdi + 8*rdx + 32]
	LONG $0x746ffec5; WORD $0x40d7       // vmovdqu    ymm6, yword [rdi + 8*rdx + 64]
	LONG $0x7c6ffec5; WORD $0x60d7       // vmovdqu    ymm7, yword [rdi + 8*rdx + 96]
	LONG $0xd45941c4; BYTE $0xd8         // vpaddq    xmm11, xmm4, xmm8
	LONG $0x197de3c4; WORD $0x01e4       // vextractf128    xmm4, ymm4, 1
	LONG $0x197d63c4; WORD $0x01c1       // vextractf128    xmm1, ymm8, 1
	LONG $0xc9d4d9c5                     // vpaddq    xmm1, xmm4, xmm1
	LONG $0xebd451c5                     // vpaddq    xmm13, xmm5, xmm3
	LONG $0x197de3c4; WORD $0x01ed       // vextractf128    xmm5, ymm5, 1
	LONG $0x197de3c4; WORD $0x01db       // vextractf128    xmm3, ymm3, 1
	LONG $0xdbd4d1c5                     // vpaddq    xmm3, xmm5, xmm3
	LONG $0xd449c1c4; BYTE $0xe9         // vpaddq    xmm5, xmm6, xmm9
	LONG $0x197de3c4; WORD $0x01f6       // vextractf128    xmm6, ymm6, 1
	LONG $0x197d63c4; WORD $0x01ca       // vextractf128    xmm2, ymm9, 1
	LONG $0xd2d4c9c5                     // vpaddq    xmm2, xmm6, xmm2
	LONG $0xd441c1c4; BYTE $0xf2         // vpaddq    xmm6, xmm7, xmm10
	LONG $0x197de3c4; WORD $0x01ff       // vextractf128    xmm7, ymm7, 1
	LONG $0x197d63c4; WORD $0x01d0       // vextractf128    xmm0, ymm10, 1
	LONG $0xc0d4c1c5                     // vpaddq    xmm0, xmm7, xmm0
	QUAD $0x000080d7bc6ffec5; BYTE $0x00 // vmovdqu    ymm7, yword [rdi + 8*rdx + 128]
	QUAD $0x0000a0d78c6f7ec5; BYTE $0x00 // vmovdqu    ymm9, yword [rdi + 8*rdx + 160]
	QUAD $0x0000c0d7946f7ec5; BYTE $0x00 // vmovdqu    ymm10, yword [rdi + 8*rdx + 192]
	QUAD $0x0000e0d7a46f7ec5; BYTE $0x00 // vmovdqu    ymm12, yword [rdi + 8*rdx + 224]
	LONG $0x197de3c4; WORD $0x01fc       // vextractf128    xmm4, ymm7, 1
	LONG $0xc9d4d9c5                     // vpaddq    xmm1, xmm4, xmm1
	LONG $0xd441c1c4; BYTE $0xe3         // vpaddq    xmm4, xmm7, xmm11
	LONG $0x185d63c4; WORD $0x01c1       // vinsertf128    ymm8, ymm4, xmm1, 1
	LONG $0x197d63c4; WORD $0x01c9       // vextractf128    xmm1, ymm9, 1
	LONG $0xcbd4f1c5                     // vpaddq    xmm1, xmm1, xmm3
	LONG $0xd431c1c4; BYTE $0xdd         // vpaddq    xmm3, xmm9, xmm13
	LONG $0x1865e3c4; WORD $0x01d9       // vinsertf128    ymm3, ymm3, xmm1, 1
	LONG $0x197d63c4; WORD $0x01d1       // vextractf128    xmm1, ymm10, 1
	LONG $0xcad4f1c5                     // vpaddq    xmm1, xmm1, xmm2
	LONG $0xd5d4a9c5                     // vpaddq    xmm2, xmm10, xmm5
	LONG $0x186d63c4; WORD $0x01c9       // vinsertf128    ymm9, ymm2, xmm1, 1
	LONG $0x197d63c4; WORD $0x01e1       // vextractf128    xmm1, ymm12, 1
	LONG $0xc0d4f1c5                     // vpaddq    xmm0, xmm1, xmm0
	LONG $0xced499c5                     // vpaddq    xmm1, xmm12, xmm6
	LONG $0x187563c4; WORD $0x01d0       // vinsertf128    ymm10, ymm1, xmm0, 1
	LONG $0x20c28348                     // add    rdx, 32
	LONG $0x02c18348                     // add    rcx, 2
	JNE  LBB3_8
	WORD $0x854d; BYTE $0xc0             // test    r8, r8
	JE   LBB3_11

LBB3_10:
	LONG $0x246ffec5; BYTE $0xd7   // vmovdqu    ymm4, yword [rdi + 8*rdx]
	LONG $0x446ffec5; WORD $0x20d7 // vmovdqu    ymm0, yword [rdi + 8*rdx + 32]
	LONG $0x4c6ffec5; WORD $0x40d7 // vmovdqu    ymm1, yword [rdi + 8*rdx + 64]
	LONG $0x546ffec5; WORD $0x60d7 // vmovdqu    ymm2, yword [rdi + 8*rdx + 96]
	LONG $0x197de3c4; WORD $0x01d5 // vextractf128    xmm5, ymm2, 1
	LONG $0x197d63c4; WORD $0x01d6 // vextractf128    xmm6, ymm10, 1
	LONG $0xeed4d1c5               // vpaddq    xmm5, xmm5, xmm6
	LONG $0xd469c1c4; BYTE $0xd2   // vpaddq    xmm2, xmm2, xmm10
	LONG $0x186d63c4; WORD $0x01d5 // vinsertf128    ymm10, ymm2, xmm5, 1
	LONG $0x197de3c4; WORD $0x01ca // vextractf128    xmm2, ymm1, 1
	LONG $0x197d63c4; WORD $0x01cd // vextractf128    xmm5, ymm9, 1
	LONG $0xd5d4e9c5               // vpaddq    xmm2, xmm2, xmm5
	LONG $0xd471c1c4; BYTE $0xc9   // vpaddq    xmm1, xmm1, xmm9
	LONG $0x187563c4; WORD $0x01ca // vinsertf128    ymm9, ymm1, xmm2, 1
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197de3c4; WORD $0x01da // vextractf128    xmm2, ymm3, 1
	LONG $0xcad4f1c5               // vpaddq    xmm1, xmm1, xmm2
	LONG $0xc3d4f9c5               // vpaddq    xmm0, xmm0, xmm3
	LONG $0x187de3c4; WORD $0x01d9 // vinsertf128    ymm3, ymm0, xmm1, 1
	LONG $0x197de3c4; WORD $0x01e0 // vextractf128    xmm0, ymm4, 1
	LONG $0x197d63c4; WORD $0x01c1 // vextractf128    xmm1, ymm8, 1
	LONG $0xc1d4f9c5               // vpaddq    xmm0, xmm0, xmm1
	LONG $0xd459c1c4; BYTE $0xc8   // vpaddq    xmm1, xmm4, xmm8
	LONG $0x187563c4; WORD $0x01c0 // vinsertf128    ymm8, ymm1, xmm0, 1

LBB3_11:
	LONG $0x197d63c4; WORD $0x01c0 // vextractf128    xmm0, ymm8, 1
	LONG $0x197de3c4; WORD $0x01d9 // vextractf128    xmm1, ymm3, 1
	LONG $0xc0d4f1c5               // vpaddq    xmm0, xmm1, xmm0
	LONG $0xd461c1c4; BYTE $0xc8   // vpaddq    xmm1, xmm3, xmm8
	LONG $0x197d63c4; WORD $0x01ca // vextractf128    xmm2, ymm9, 1
	LONG $0x197d63c4; WORD $0x01d3 // vextractf128    xmm3, ymm10, 1
	LONG $0xd3d4e9c5               // vpaddq    xmm2, xmm2, xmm3
	LONG $0xc2d4f9c5               // vpaddq    xmm0, xmm0, xmm2
	LONG $0xd431c1c4; BYTE $0xd2   // vpaddq    xmm2, xmm9, xmm10
	LONG $0xcad4f1c5               // vpaddq    xmm1, xmm1, xmm2
	LONG $0xc0d4f1c5               // vpaddq    xmm0, xmm1, xmm0
	LONG $0xc870f9c5; BYTE $0x4e   // vpshufd    xmm1, xmm0, 78
	LONG $0xc1d4f9c5               // vpaddq    xmm0, xmm0, xmm1
	LONG $0x7ef9e1c4; BYTE $0xc1   // vmovq    rcx, xmm0
	WORD $0x3948; BYTE $0xf0       // cmp    rax, rsi
	JE   LBB3_13

LBB3_12:
	LONG $0xc70c0348         // add    rcx, qword [rdi + 8*rax]
	LONG $0x01c08348         // add    rax, 1
	WORD $0x3948; BYTE $0xc6 // cmp    rsi, rax
	JNE  LBB3_12

LBB3_13:
	LONG $0x2a8b61c4; BYTE $0xde // vcvtsi2sd    xmm11, xmm14, rsi
	WORD $0x8548; BYTE $0xf6     // test    rsi, rsi
	JLE  LBB3_14
	LONG $0x2a8be1c4; BYTE $0xc9 // vcvtsi2sd    xmm1, xmm14, rcx
	LONG $0x5e7341c4; BYTE $0xe3 // vdivsd    xmm12, xmm1, xmm11
	LONG $0x10fe8348             // cmp    rsi, 16
	JAE  LBB3_17
	LONG $0xd2efe9c5             // vpxor    xmm2, xmm2, xmm2
	WORD $0xc031                 // xor    eax, eax
	JMP  LBB3_20

LBB3_14:
	LONG $0xd2efe9c5             // vpxor    xmm2, xmm2, xmm2
	LONG $0x5e6bc1c4; BYTE $0xc3 // vdivsd    xmm0, xmm2, xmm11
	JMP  LBB3_22

LBB3_17:
	WORD $0x8948; BYTE $0xf0       // mov    rax, rsi
	LONG $0xf0e08348               // and    rax, -16
	LONG $0x127bc1c4; BYTE $0xd4   // vmovddup    xmm2, xmm12
	LONG $0x186d63c4; WORD $0x01ea // vinsertf128    ymm13, ymm2, xmm2, 1
	LONG $0x570941c4; BYTE $0xf6   // vxorpd    xmm14, xmm14, xmm14
	WORD $0xc931                   // xor    ecx, ecx
	LONG $0xe4efd9c5               // vpxor    xmm4, xmm4, xmm4
	LONG $0xedefd1c5               // vpxor    xmm5, xmm5, xmm5
	LONG $0xf6efc9c5               // vpxor    xmm6, xmm6, xmm6

LBB3_18:
	LONG $0x046f7ec5; BYTE $0xcf   // vmovdqu    ymm8, yword [rdi + 8*rcx]
	LONG $0x546f7ec5; WORD $0x20cf // vmovdqu    ymm10, yword [rdi + 8*rcx + 32]
	LONG $0x4c6f7ec5; WORD $0x40cf // vmovdqu    ymm9, yword [rdi + 8*rcx + 64]
	LONG $0x197d63c4; WORD $0x01c0 // vextractf128    xmm0, ymm8, 1
	LONG $0x16f9e3c4; WORD $0x01c2 // vpextrq    rdx, xmm0, 1
	LONG $0x7c6ffec5; WORD $0x60cf // vmovdqu    ymm7, yword [rdi + 8*rcx + 96]
	LONG $0x2a83e1c4; BYTE $0xca   // vcvtsi2sd    xmm1, xmm15, rdx
	LONG $0x7ef9e1c4; BYTE $0xc2   // vmovq    rdx, xmm0
	LONG $0x2a83e1c4; BYTE $0xc2   // vcvtsi2sd    xmm0, xmm15, rdx
	LONG $0x16f963c4; WORD $0x01c2 // vpextrq    rdx, xmm8, 1
	LONG $0xc116f8c5               // vmovlhps    xmm0, xmm0, xmm1
	LONG $0x2a83e1c4; BYTE $0xca   // vcvtsi2sd    xmm1, xmm15, rdx
	LONG $0x7ef961c4; BYTE $0xc2   // vmovq    rdx, xmm8
	LONG $0x2a83e1c4; BYTE $0xda   // vcvtsi2sd    xmm3, xmm15, rdx
	LONG $0xc916e0c5               // vmovlhps    xmm1, xmm3, xmm1
	LONG $0x197d63c4; WORD $0x01d3 // vextractf128    xmm3, ymm10, 1
	LONG $0x16f9e3c4; WORD $0x01da // vpextrq    rdx, xmm3, 1
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0x187563c4; WORD $0x01c0 // vinsertf128    ymm8, ymm1, xmm0, 1
	LONG $0x7ef9e1c4; BYTE $0xda   // vmovq    rdx, xmm3
	LONG $0x2a83e1c4; BYTE $0xc2   // vcvtsi2sd    xmm0, xmm15, rdx
	LONG $0x16f963c4; WORD $0x01d2 // vpextrq    rdx, xmm10, 1
	LONG $0x2a83e1c4; BYTE $0xca   // vcvtsi2sd    xmm1, xmm15, rdx
	LONG $0xc216f8c5               // vmovlhps    xmm0, xmm0, xmm2
	LONG $0x7ef961c4; BYTE $0xd2   // vmovq    rdx, xmm10
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0xc916e8c5               // vmovlhps    xmm1, xmm2, xmm1
	LONG $0x197d63c4; WORD $0x01ca // vextractf128    xmm2, ymm9, 1
	LONG $0x16f9e3c4; WORD $0x01d2 // vpextrq    rdx, xmm2, 1
	LONG $0x187563c4; WORD $0x01d0 // vinsertf128    ymm10, ymm1, xmm0, 1
	LONG $0x2a83e1c4; BYTE $0xca   // vcvtsi2sd    xmm1, xmm15, rdx
	LONG $0x7ef9e1c4; BYTE $0xd2   // vmovq    rdx, xmm2
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0x16f963c4; WORD $0x01ca // vpextrq    rdx, xmm9, 1
	LONG $0xc916e8c5               // vmovlhps    xmm1, xmm2, xmm1
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0x7ef961c4; BYTE $0xca   // vmovq    rdx, xmm9
	LONG $0x2a83e1c4; BYTE $0xda   // vcvtsi2sd    xmm3, xmm15, rdx
	LONG $0xd216e0c5               // vmovlhps    xmm2, xmm3, xmm2
	LONG $0x197de3c4; WORD $0x01fb // vextractf128    xmm3, ymm7, 1
	LONG $0x16f9e3c4; WORD $0x01da // vpextrq    rdx, xmm3, 1
	LONG $0x2a83e1c4; BYTE $0xc2   // vcvtsi2sd    xmm0, xmm15, rdx
	LONG $0x186de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm2, xmm1, 1
	LONG $0x7ef9e1c4; BYTE $0xda   // vmovq    rdx, xmm3
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0x16f9e3c4; WORD $0x01fa // vpextrq    rdx, xmm7, 1
	LONG $0x2a83e1c4; BYTE $0xda   // vcvtsi2sd    xmm3, xmm15, rdx
	LONG $0xc016e8c5               // vmovlhps    xmm0, xmm2, xmm0
	LONG $0x7ef9e1c4; BYTE $0xfa   // vmovq    rdx, xmm7
	LONG $0x2a83e1c4; BYTE $0xd2   // vcvtsi2sd    xmm2, xmm15, rdx
	LONG $0xd316e8c5               // vmovlhps    xmm2, xmm2, xmm3
	LONG $0x186de3c4; WORD $0x01c0 // vinsertf128    ymm0, ymm2, xmm0, 1
	LONG $0x5c3dc1c4; BYTE $0xd5   // vsubpd    ymm2, ymm8, ymm13
	LONG $0x5c2dc1c4; BYTE $0xdd   // vsubpd    ymm3, ymm10, ymm13
	LONG $0x5c75c1c4; BYTE $0xcd   // vsubpd    ymm1, ymm1, ymm13
	LONG $0x5c7dc1c4; BYTE $0xc5   // vsubpd    ymm0, ymm0, ymm13
	LONG $0xd259edc5               // vmulpd    ymm2, ymm2, ymm2
	LONG $0x586d41c4; BYTE $0xf6   // vaddpd    ymm14, ymm2, ymm14
	LONG $0xd359e5c5               // vmulpd    ymm2, ymm3, ymm3
	LONG $0xe458edc5               // vaddpd    ymm4, ymm2, ymm4
	LONG $0xc959f5c5               // vmulpd    ymm1, ymm1, ymm1
	LONG $0xed58f5c5               // vaddpd    ymm5, ymm1, ymm5
	LONG $0xc059fdc5               // vmulpd    ymm0, ymm0, ymm0
	LONG $0xf658fdc5               // vaddpd    ymm6, ymm0, ymm6
	LONG $0x10c18348               // add    rcx, 16
	WORD $0x3948; BYTE $0xc8       // cmp    rax, rcx
	JNE  LBB3_18
	LONG $0x585dc1c4; BYTE $0xc6   // vaddpd    ymm0, ymm4, ymm14
	LONG $0xc058d5c5               // vaddpd    ymm0, ymm5, ymm0
	LONG $0xc058cdc5               // vaddpd    ymm0, ymm6, ymm0
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0xc158fdc5               // vaddpd    ymm0, ymm0, ymm1
	LONG $0xd07cfdc5               // vhaddpd    ymm2, ymm0, ymm0
	WORD $0x3948; BYTE $0xf0       // cmp    rax, rsi
	JE   LBB3_21

LBB3_20:
	LONG $0x2a83e1c4; WORD $0xc704 // vcvtsi2sd    xmm0, xmm15, qword [rdi + 8*rax]
	LONG $0x5c7bc1c4; BYTE $0xc4   // vsubsd    xmm0, xmm0, xmm12
	LONG $0xc059fbc5               // vmulsd    xmm0, xmm0, xmm0
	LONG $0xd258fbc5               // vaddsd    xmm2, xmm0, xmm2
	LONG $0x01c08348               // add    rax, 1
	WORD $0x3948; BYTE $0xc6       // cmp    rsi, rax
	JNE  LBB3_20

LBB3_21:
	LONG $0x5e6bc1c4; BYTE $0xc3 // vdivsd    xmm0, xmm2, xmm11

LBB3_22:
	BYTE $0xc5; BYTE $0xf8; BYTE $0x77 // VZEROUPPER
	MOVQ X0, x+24(FP)
	RET

LBB3_6:
	LONG $0xef3941c4; BYTE $0xc0 // vpxor    xmm8, xmm8, xmm8
	WORD $0xd231                 // xor    edx, edx
	LONG $0xdbefe1c5             // vpxor    xmm3, xmm3, xmm3
	LONG $0xef3141c4; BYTE $0xc9 // vpxor    xmm9, xmm9, xmm9
	LONG $0xef2941c4; BYTE $0xd2 // vpxor    xmm10, xmm10, xmm10
	WORD $0x854d; BYTE $0xc0     // test    r8, r8
	JNE  LBB3_10
	JMP  LBB3_11

TEXT sampleVarianceSSE42<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX

	WORD $0x8548; BYTE $0xf6 // test    rsi, rsi
	JE   LBB3_1
	JLE  LBB3_25
	LONG $0x03fe8348         // cmp    rsi, 3
	JA   LBB3_5
	WORD $0xc031             // xor    eax, eax
	WORD $0xc931             // xor    ecx, ecx
	JMP  LBB3_13

LBB3_1:
	WORD $0x570f; BYTE $0xc0 // xorps    xmm0, xmm0
	JMP  LBB3_24

LBB3_25:
	LONG $0x2a0f48f2; BYTE $0xce // cvtsi2sd    xmm1, rsi
	LONG $0xc0570f66             // xorpd    xmm0, xmm0
	LONG $0xc15e0ff2             // divsd    xmm0, xmm1
	JMP  LBB3_24

LBB3_5:
	WORD $0x8948; BYTE $0xf0 // mov    rax, rsi
	LONG $0xfce08348         // and    rax, -4
	LONG $0xfc508d48         // lea    rdx, [rax - 4]
	WORD $0x8948; BYTE $0xd1 // mov    rcx, rdx
	LONG $0x02e9c148         // shr    rcx, 2
	LONG $0x01c18348         // add    rcx, 1
	WORD $0x8941; BYTE $0xc8 // mov    r8d, ecx
	LONG $0x03e08341         // and    r8d, 3
	LONG $0x0cfa8348         // cmp    rdx, 12
	JAE  LBB3_7
	LONG $0xc0ef0f66         // pxor    xmm0, xmm0
	WORD $0xd231             // xor    edx, edx
	LONG $0xc9ef0f66         // pxor    xmm1, xmm1
	WORD $0x854d; BYTE $0xc0 // test    r8, r8
	JNE  LBB3_10
	JMP  LBB3_12

LBB3_7:
	LONG $0x000001ba; BYTE $0x00 // mov    edx, 1
	WORD $0x2948; BYTE $0xca     // sub    rdx, rcx
	LONG $0x100c8d49             // lea    rcx, [r8 + rdx]
	LONG $0xffc18348             // add    rcx, -1
	LONG $0xc0ef0f66             // pxor    xmm0, xmm0
	WORD $0xd231                 // xor    edx, edx
	LONG $0xc9ef0f66             // pxor    xmm1, xmm1

LBB3_8:
	LONG $0x146f0ff3; BYTE $0xd7   // movdqu    xmm2, oword [rdi + 8*rdx]
	LONG $0xd0d40f66               // paddq    xmm2, xmm0
	LONG $0x446f0ff3; WORD $0x10d7 // movdqu    xmm0, oword [rdi + 8*rdx + 16]
	LONG $0xc1d40f66               // paddq    xmm0, xmm1
	LONG $0x4c6f0ff3; WORD $0x20d7 // movdqu    xmm1, oword [rdi + 8*rdx + 32]
	LONG $0x5c6f0ff3; WORD $0x30d7 // movdqu    xmm3, oword [rdi + 8*rdx + 48]
	LONG $0x646f0ff3; WORD $0x40d7 // movdqu    xmm4, oword [rdi + 8*rdx + 64]
	LONG $0xe1d40f66               // paddq    xmm4, xmm1
	LONG $0xe2d40f66               // paddq    xmm4, xmm2
	LONG $0x546f0ff3; WORD $0x50d7 // movdqu    xmm2, oword [rdi + 8*rdx + 80]
	LONG $0xd3d40f66               // paddq    xmm2, xmm3
	LONG $0xd0d40f66               // paddq    xmm2, xmm0
	LONG $0x446f0ff3; WORD $0x60d7 // movdqu    xmm0, oword [rdi + 8*rdx + 96]
	LONG $0xc4d40f66               // paddq    xmm0, xmm4
	LONG $0x4c6f0ff3; WORD $0x70d7 // movdqu    xmm1, oword [rdi + 8*rdx + 112]
	LONG $0xcad40f66               // paddq    xmm1, xmm2
	LONG $0x10c28348               // add    rdx, 16
	LONG $0x04c18348               // add    rcx, 4
	JNE  LBB3_8
	WORD $0x854d; BYTE $0xc0       // test    r8, r8
	JE   LBB3_12

LBB3_10:
	LONG $0xd70c8d48         // lea    rcx, [rdi + 8*rdx]
	LONG $0x10c18348         // add    rcx, 16
	WORD $0xf749; BYTE $0xd8 // neg    r8

LBB3_11:
	LONG $0x516f0ff3; BYTE $0xf0 // movdqu    xmm2, oword [rcx - 16]
	LONG $0xc2d40f66             // paddq    xmm0, xmm2
	LONG $0x116f0ff3             // movdqu    xmm2, oword [rcx]
	LONG $0xcad40f66             // paddq    xmm1, xmm2
	LONG $0x20c18348             // add    rcx, 32
	LONG $0x01c08349             // add    r8, 1
	JNE  LBB3_11

LBB3_12:
	LONG $0xc1d40f66             // paddq    xmm0, xmm1
	LONG $0xc8700f66; BYTE $0x4e // pshufd    xmm1, xmm0, 78
	LONG $0xc8d40f66             // paddq    xmm1, xmm0
	LONG $0x7e0f4866; BYTE $0xc9 // movq    rcx, xmm1
	WORD $0x3948; BYTE $0xf0     // cmp    rax, rsi
	JE   LBB3_14

LBB3_13:
	LONG $0xc70c0348         // add    rcx, qword [rdi + 8*rax]
	LONG $0x01c08348         // add    rax, 1
	WORD $0x3948; BYTE $0xc6 // cmp    rsi, rax
	JNE  LBB3_13

LBB3_14:
	WORD $0x570f; BYTE $0xc9     // xorps    xmm1, xmm1
	LONG $0x2a0f48f2; BYTE $0xce // cvtsi2sd    xmm1, rsi
	WORD $0x8548; BYTE $0xf6     // test    rsi, rsi
	JLE  LBB3_15
	WORD $0x570f; BYTE $0xd2     // xorps    xmm2, xmm2
	LONG $0x2a0f48f2; BYTE $0xd1 // cvtsi2sd    xmm2, rcx
	LONG $0xd15e0ff2             // divsd    xmm2, xmm1
	LONG $0xff4e8d48             // lea    rcx, [rsi - 1]
	WORD $0xf089                 // mov    eax, esi
	WORD $0xe083; BYTE $0x03     // and    eax, 3
	LONG $0x03f98348             // cmp    rcx, 3
	JAE  LBB3_18
	LONG $0xc0ef0f66             // pxor    xmm0, xmm0
	WORD $0xc931                 // xor    ecx, ecx
	WORD $0x8548; BYTE $0xc0     // test    rax, rax
	JNE  LBB3_21
	JMP  LBB3_23

LBB3_15:
	LONG $0xc0ef0f66 // pxor    xmm0, xmm0
	LONG $0xc15e0ff2 // divsd    xmm0, xmm1
	JMP  LBB3_24

LBB3_18:
	WORD $0x2948; BYTE $0xc6 // sub    rsi, rax
	LONG $0xc0ef0f66         // pxor    xmm0, xmm0
	WORD $0xc931             // xor    ecx, ecx

LBB3_19:
	WORD $0x570f; BYTE $0xdb                   // xorps    xmm3, xmm3
	LONG $0x2a0f48f2; WORD $0xcf1c             // cvtsi2sd    xmm3, qword [rdi + 8*rcx]
	LONG $0xda5c0ff2                           // subsd    xmm3, xmm2
	LONG $0xdb590ff2                           // mulsd    xmm3, xmm3
	WORD $0x570f; BYTE $0xe4                   // xorps    xmm4, xmm4
	LONG $0x2a0f48f2; WORD $0xcf64; BYTE $0x08 // cvtsi2sd    xmm4, qword [rdi + 8*rcx + 8]
	LONG $0xd8580ff2                           // addsd    xmm3, xmm0
	LONG $0xe25c0ff2                           // subsd    xmm4, xmm2
	WORD $0x570f; BYTE $0xed                   // xorps    xmm5, xmm5
	LONG $0x2a0f48f2; WORD $0xcf6c; BYTE $0x10 // cvtsi2sd    xmm5, qword [rdi + 8*rcx + 16]
	LONG $0xe4590ff2                           // mulsd    xmm4, xmm4
	LONG $0xea5c0ff2                           // subsd    xmm5, xmm2
	LONG $0xed590ff2                           // mulsd    xmm5, xmm5
	LONG $0xec580ff2                           // addsd    xmm5, xmm4
	LONG $0xeb580ff2                           // addsd    xmm5, xmm3
	WORD $0x570f; BYTE $0xc0                   // xorps    xmm0, xmm0
	LONG $0x2a0f48f2; WORD $0xcf44; BYTE $0x18 // cvtsi2sd    xmm0, qword [rdi + 8*rcx + 24]
	LONG $0xc25c0ff2                           // subsd    xmm0, xmm2
	LONG $0xc0590ff2                           // mulsd    xmm0, xmm0
	LONG $0xc5580ff2                           // addsd    xmm0, xmm5
	LONG $0x04c18348                           // add    rcx, 4
	WORD $0x3948; BYTE $0xce                   // cmp    rsi, rcx
	JNE  LBB3_19
	WORD $0x8548; BYTE $0xc0                   // test    rax, rax
	JE   LBB3_23

LBB3_21:
	LONG $0xcf0c8d48 // lea    rcx, [rdi + 8*rcx]
	WORD $0xd231     // xor    edx, edx

LBB3_22:
	WORD $0x570f; BYTE $0xdb       // xorps    xmm3, xmm3
	LONG $0x2a0f48f2; WORD $0xd11c // cvtsi2sd    xmm3, qword [rcx + 8*rdx]
	LONG $0xda5c0ff2               // subsd    xmm3, xmm2
	LONG $0xdb590ff2               // mulsd    xmm3, xmm3
	LONG $0xc3580ff2               // addsd    xmm0, xmm3
	LONG $0x01c28348               // add    rdx, 1
	WORD $0x3948; BYTE $0xd0       // cmp    rax, rdx
	JNE  LBB3_22

LBB3_23:
	LONG $0xc15e0ff2 // divsd    xmm0, xmm1

LBB3_24:
	MOVQ X0, x+24(FP)
	RET

DATA LCDATA1<>+0x000(SB)/8, $0x8000000000000000
DATA LCDATA1<>+0x008(SB)/8, $0x8000000000000000
DATA LCDATA1<>+0x010(SB)/8, $0x8000000000000000
DATA LCDATA1<>+0x018(SB)/8, $0x8000000000000000
GLOBL LCDATA1<>(SB), 8, $32

TEXT sampleMaxAVX2<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX
	LEAQ LCDATA1<>(SB), BP

	WORD $0x8548; BYTE $0xf6               // test    rsi, rsi
	JE   LBB1_1
	QUAD $0x000000000000b848; WORD $0x8000 // mov    rax, -9223372036854775808
	JLE  LBB1_13
	LONG $0x10fe8348                       // cmp    rsi, 16
	JAE  LBB1_5
	WORD $0xc931                           // xor    ecx, ecx
	JMP  LBB1_12

LBB1_1:
	WORD $0xc031      // xor    eax, eax
	MOVQ AX, x+24(FP)
	RET

LBB1_5:
	WORD $0x8948; BYTE $0xf1       // mov    rcx, rsi
	LONG $0xf0e18348               // and    rcx, -16
	LONG $0xf0518d48               // lea    rdx, [rcx - 16]
	WORD $0x8948; BYTE $0xd0       // mov    rax, rdx
	LONG $0x04e8c148               // shr    rax, 4
	LONG $0x01c08348               // add    rax, 1
	WORD $0x8941; BYTE $0xc0       // mov    r8d, eax
	LONG $0x01e08341               // and    r8d, 1
	WORD $0x8548; BYTE $0xd2       // test    rdx, rdx
	JE   LBB1_6
	LONG $0x000001ba; BYTE $0x00   // mov    edx, 1
	WORD $0x2948; BYTE $0xc2       // sub    rdx, rax
	WORD $0x014c; BYTE $0xc2       // add    rdx, r8
	LONG $0xffc28348               // add    rdx, -1
	LONG $0x597de2c4; WORD $0x0045 // vpbroadcastq    ymm0, qword 0[rbp] /* [rip + .LCPI1_0] */
	WORD $0xc031                   // xor    eax, eax
	LONG $0xc86ffdc5               // vmovdqa    ymm1, ymm0
	LONG $0xd06ffdc5               // vmovdqa    ymm2, ymm0
	LONG $0xd86ffdc5               // vmovdqa    ymm3, ymm0

LBB1_8:
	LONG $0x246ffec5; BYTE $0xc7         // vmovdqu    ymm4, yword [rdi + 8*rax]
	LONG $0x6c6ffec5; WORD $0x20c7       // vmovdqu    ymm5, yword [rdi + 8*rax + 32]
	LONG $0x746ffec5; WORD $0x40c7       // vmovdqu    ymm6, yword [rdi + 8*rax + 64]
	LONG $0x375de2c4; BYTE $0xf8         // vpcmpgtq    ymm7, ymm4, ymm0
	LONG $0x4b7de3c4; WORD $0x70c4       // vblendvpd    ymm0, ymm0, ymm4, ymm7
	LONG $0x646ffec5; WORD $0x60c7       // vmovdqu    ymm4, yword [rdi + 8*rax + 96]
	LONG $0x3755e2c4; BYTE $0xf9         // vpcmpgtq    ymm7, ymm5, ymm1
	LONG $0x4b75e3c4; WORD $0x70cd       // vblendvpd    ymm1, ymm1, ymm5, ymm7
	LONG $0x374de2c4; BYTE $0xea         // vpcmpgtq    ymm5, ymm6, ymm2
	LONG $0x4b6de3c4; WORD $0x50d6       // vblendvpd    ymm2, ymm2, ymm6, ymm5
	LONG $0x375de2c4; BYTE $0xeb         // vpcmpgtq    ymm5, ymm4, ymm3
	LONG $0x4b65e3c4; WORD $0x50dc       // vblendvpd    ymm3, ymm3, ymm4, ymm5
	QUAD $0x000080c7a46ffec5; BYTE $0x00 // vmovdqu    ymm4, yword [rdi + 8*rax + 128]
	QUAD $0x0000a0c7ac6ffec5; BYTE $0x00 // vmovdqu    ymm5, yword [rdi + 8*rax + 160]
	QUAD $0x0000c0c7b46ffec5; BYTE $0x00 // vmovdqu    ymm6, yword [rdi + 8*rax + 192]
	LONG $0x375de2c4; BYTE $0xf8         // vpcmpgtq    ymm7, ymm4, ymm0
	LONG $0x4b7de3c4; WORD $0x70c4       // vblendvpd    ymm0, ymm0, ymm4, ymm7
	QUAD $0x0000e0c7a46ffec5; BYTE $0x00 // vmovdqu    ymm4, yword [rdi + 8*rax + 224]
	LONG $0x3755e2c4; BYTE $0xf9         // vpcmpgtq    ymm7, ymm5, ymm1
	LONG $0x4b75e3c4; WORD $0x70cd       // vblendvpd    ymm1, ymm1, ymm5, ymm7
	LONG $0x374de2c4; BYTE $0xea         // vpcmpgtq    ymm5, ymm6, ymm2
	LONG $0x4b6de3c4; WORD $0x50d6       // vblendvpd    ymm2, ymm2, ymm6, ymm5
	LONG $0x375de2c4; BYTE $0xeb         // vpcmpgtq    ymm5, ymm4, ymm3
	LONG $0x4b65e3c4; WORD $0x50dc       // vblendvpd    ymm3, ymm3, ymm4, ymm5
	LONG $0x20c08348                     // add    rax, 32
	LONG $0x02c28348                     // add    rdx, 2
	JNE  LBB1_8
	WORD $0x854d; BYTE $0xc0             // test    r8, r8
	JE   LBB1_11

LBB1_10:
	LONG $0x646ffec5; WORD $0x60c7 // vmovdqu    ymm4, yword [rdi + 8*rax + 96]
	LONG $0x375de2c4; BYTE $0xeb   // vpcmpgtq    ymm5, ymm4, ymm3
	LONG $0x4b65e3c4; WORD $0x50dc // vblendvpd    ymm3, ymm3, ymm4, ymm5
	LONG $0x646ffec5; WORD $0x40c7 // vmovdqu    ymm4, yword [rdi + 8*rax + 64]
	LONG $0x375de2c4; BYTE $0xea   // vpcmpgtq    ymm5, ymm4, ymm2
	LONG $0x4b6de3c4; WORD $0x50d4 // vblendvpd    ymm2, ymm2, ymm4, ymm5
	LONG $0x646ffec5; WORD $0x20c7 // vmovdqu    ymm4, yword [rdi + 8*rax + 32]
	LONG $0x375de2c4; BYTE $0xe9   // vpcmpgtq    ymm5, ymm4, ymm1
	LONG $0x4b75e3c4; WORD $0x50cc // vblendvpd    ymm1, ymm1, ymm4, ymm5
	LONG $0x246ffec5; BYTE $0xc7   // vmovdqu    ymm4, yword [rdi + 8*rax]
	LONG $0x375de2c4; BYTE $0xe8   // vpcmpgtq    ymm5, ymm4, ymm0
	LONG $0x4b7de3c4; WORD $0x50c4 // vblendvpd    ymm0, ymm0, ymm4, ymm5

LBB1_11:
	LONG $0x377de2c4; BYTE $0xe1   // vpcmpgtq    ymm4, ymm0, ymm1
	LONG $0x4b75e3c4; WORD $0x40c0 // vblendvpd    ymm0, ymm1, ymm0, ymm4
	LONG $0x377de2c4; BYTE $0xca   // vpcmpgtq    ymm1, ymm0, ymm2
	LONG $0x4b6de3c4; WORD $0x10c0 // vblendvpd    ymm0, ymm2, ymm0, ymm1
	LONG $0x377de2c4; BYTE $0xcb   // vpcmpgtq    ymm1, ymm0, ymm3
	LONG $0x4b65e3c4; WORD $0x10c0 // vblendvpd    ymm0, ymm3, ymm0, ymm1
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x377de2c4; BYTE $0xd1   // vpcmpgtq    ymm2, ymm0, ymm1
	LONG $0x4b75e3c4; WORD $0x20c0 // vblendvpd    ymm0, ymm1, ymm0, ymm2
	LONG $0x0479e3c4; WORD $0x4ec8 // vpermilps    xmm1, xmm0, 78
	LONG $0x377de2c4; BYTE $0xd1   // vpcmpgtq    ymm2, ymm0, ymm1
	LONG $0x4b75e3c4; WORD $0x20c0 // vblendvpd    ymm0, ymm1, ymm0, ymm2
	LONG $0x7ef9e1c4; BYTE $0xc0   // vmovq    rax, xmm0
	WORD $0x3948; BYTE $0xf1       // cmp    rcx, rsi
	JE   LBB1_13

LBB1_12:
	LONG $0xcf148b48         // mov    rdx, qword [rdi + 8*rcx]
	WORD $0x3948; BYTE $0xc2 // cmp    rdx, rax
	LONG $0xc24d0f48         // cmovge    rax, rdx
	LONG $0x01c18348         // add    rcx, 1
	WORD $0x3948; BYTE $0xce // cmp    rsi, rcx
	JNE  LBB1_12

LBB1_13:
	BYTE $0xc5; BYTE $0xf8; BYTE $0x77 // VZEROUPPER
	MOVQ AX, x+24(FP)
	RET

LBB1_6:
	LONG $0x597de2c4; WORD $0x0045 // vpbroadcastq    ymm0, qword 0[rbp] /* [rip + .LCPI1_0] */
	WORD $0xc031                   // xor    eax, eax
	LONG $0xc86ffdc5               // vmovdqa    ymm1, ymm0
	LONG $0xd06ffdc5               // vmovdqa    ymm2, ymm0
	LONG $0xd86ffdc5               // vmovdqa    ymm3, ymm0
	WORD $0x854d; BYTE $0xc0       // test    r8, r8
	JNE  LBB1_10
	JMP  LBB1_11

TEXT sampleMaxAVX<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX
	LEAQ LCDATA1<>(SB), BP

	WORD $0x8548; BYTE $0xf6               // test    rsi, rsi
	JE   LBB1_1
	QUAD $0x000000000000b848; WORD $0x8000 // mov    rax, -9223372036854775808
	JLE  LBB1_13
	LONG $0x10fe8348                       // cmp    rsi, 16
	JAE  LBB1_5
	WORD $0xc931                           // xor    ecx, ecx
	JMP  LBB1_12

LBB1_1:
	WORD $0xc031      // xor    eax, eax
	MOVQ AX, x+24(FP)
	RET

LBB1_5:
	WORD $0x8948; BYTE $0xf1     // mov    rcx, rsi
	LONG $0xf0e18348             // and    rcx, -16
	LONG $0xf0518d48             // lea    rdx, [rcx - 16]
	WORD $0x8948; BYTE $0xd0     // mov    rax, rdx
	LONG $0x04e8c148             // shr    rax, 4
	LONG $0x01c08348             // add    rax, 1
	WORD $0x8941; BYTE $0xc0     // mov    r8d, eax
	LONG $0x01e08341             // and    r8d, 1
	WORD $0x8548; BYTE $0xd2     // test    rdx, rdx
	JE   LBB1_6
	LONG $0x000001ba; BYTE $0x00 // mov    edx, 1
	WORD $0x2948; BYTE $0xc2     // sub    rdx, rax
	WORD $0x014c; BYTE $0xc2     // add    rdx, r8
	LONG $0xffc28348             // add    rdx, -1
	LONG $0x4d6f7dc5; BYTE $0x00 // vmovdqa    ymm9, yword 0[rbp] /* [rip + .LCPI1_0] */
	WORD $0xc031                 // xor    eax, eax
	LONG $0x6f7dc1c4; BYTE $0xd9 // vmovdqa    ymm3, ymm9
	LONG $0x6f7dc1c4; BYTE $0xd1 // vmovdqa    ymm2, ymm9
	LONG $0x6f7d41c4; BYTE $0xc1 // vmovdqa    ymm8, ymm9

LBB1_8:
	LONG $0x246ffec5; BYTE $0xc7         // vmovdqu    ymm4, yword [rdi + 8*rax]
	LONG $0x6c6ffec5; WORD $0x20c7       // vmovdqu    ymm5, yword [rdi + 8*rax + 32]
	LONG $0x746ffec5; WORD $0x40c7       // vmovdqu    ymm6, yword [rdi + 8*rax + 64]
	LONG $0x197de3c4; WORD $0x01e7       // vextractf128    xmm7, ymm4, 1
	LONG $0x197d63c4; WORD $0x01c9       // vextractf128    xmm1, ymm9, 1
	LONG $0x3741e2c4; BYTE $0xc9         // vpcmpgtq    xmm1, xmm7, xmm1
	LONG $0x3759c2c4; BYTE $0xf9         // vpcmpgtq    xmm7, xmm4, xmm9
	LONG $0x1845e3c4; WORD $0x01c9       // vinsertf128    ymm1, ymm7, xmm1, 1
	LONG $0x4b3563c4; WORD $0x10cc       // vblendvpd    ymm9, ymm9, ymm4, ymm1
	LONG $0x4c6ffec5; WORD $0x60c7       // vmovdqu    ymm1, yword [rdi + 8*rax + 96]
	LONG $0x197de3c4; WORD $0x01ec       // vextractf128    xmm4, ymm5, 1
	LONG $0x197de3c4; WORD $0x01df       // vextractf128    xmm7, ymm3, 1
	LONG $0x3759e2c4; BYTE $0xe7         // vpcmpgtq    xmm4, xmm4, xmm7
	LONG $0x3751e2c4; BYTE $0xfb         // vpcmpgtq    xmm7, xmm5, xmm3
	LONG $0x1845e3c4; WORD $0x01e4       // vinsertf128    ymm4, ymm7, xmm4, 1
	LONG $0x4b65e3c4; WORD $0x40dd       // vblendvpd    ymm3, ymm3, ymm5, ymm4
	LONG $0x197de3c4; WORD $0x01f4       // vextractf128    xmm4, ymm6, 1
	LONG $0x197de3c4; WORD $0x01d5       // vextractf128    xmm5, ymm2, 1
	LONG $0x3759e2c4; BYTE $0xe5         // vpcmpgtq    xmm4, xmm4, xmm5
	LONG $0x3749e2c4; BYTE $0xea         // vpcmpgtq    xmm5, xmm6, xmm2
	LONG $0x1855e3c4; WORD $0x01e4       // vinsertf128    ymm4, ymm5, xmm4, 1
	LONG $0x4b6de3c4; WORD $0x40d6       // vblendvpd    ymm2, ymm2, ymm6, ymm4
	LONG $0x197de3c4; WORD $0x01cc       // vextractf128    xmm4, ymm1, 1
	LONG $0x197d63c4; WORD $0x01c5       // vextractf128    xmm5, ymm8, 1
	LONG $0x3759e2c4; BYTE $0xe5         // vpcmpgtq    xmm4, xmm4, xmm5
	LONG $0x3771c2c4; BYTE $0xe8         // vpcmpgtq    xmm5, xmm1, xmm8
	LONG $0x1855e3c4; WORD $0x01e4       // vinsertf128    ymm4, ymm5, xmm4, 1
	LONG $0x4b3de3c4; WORD $0x40c9       // vblendvpd    ymm1, ymm8, ymm1, ymm4
	QUAD $0x000080c7a46ffec5; BYTE $0x00 // vmovdqu    ymm4, yword [rdi + 8*rax + 128]
	QUAD $0x0000a0c7ac6ffec5; BYTE $0x00 // vmovdqu    ymm5, yword [rdi + 8*rax + 160]
	QUAD $0x0000c0c7b46ffec5; BYTE $0x00 // vmovdqu    ymm6, yword [rdi + 8*rax + 192]
	LONG $0x197d63c4; WORD $0x01cf       // vextractf128    xmm7, ymm9, 1
	LONG $0x197de3c4; WORD $0x01e0       // vextractf128    xmm0, ymm4, 1
	LONG $0x3779e2c4; BYTE $0xc7         // vpcmpgtq    xmm0, xmm0, xmm7
	LONG $0x3759c2c4; BYTE $0xf9         // vpcmpgtq    xmm7, xmm4, xmm9
	LONG $0x1845e3c4; WORD $0x01c0       // vinsertf128    ymm0, ymm7, xmm0, 1
	LONG $0x4b3563c4; WORD $0x00cc       // vblendvpd    ymm9, ymm9, ymm4, ymm0
	QUAD $0x0000e0c7a46ffec5; BYTE $0x00 // vmovdqu    ymm4, yword [rdi + 8*rax + 224]
	LONG $0x197de3c4; WORD $0x01df       // vextractf128    xmm7, ymm3, 1
	LONG $0x197de3c4; WORD $0x01e8       // vextractf128    xmm0, ymm5, 1
	LONG $0x3779e2c4; BYTE $0xc7         // vpcmpgtq    xmm0, xmm0, xmm7
	LONG $0x3751e2c4; BYTE $0xfb         // vpcmpgtq    xmm7, xmm5, xmm3
	LONG $0x1845e3c4; WORD $0x01c0       // vinsertf128    ymm0, ymm7, xmm0, 1
	LONG $0x4b65e3c4; WORD $0x00dd       // vblendvpd    ymm3, ymm3, ymm5, ymm0
	LONG $0x197de3c4; WORD $0x01d0       // vextractf128    xmm0, ymm2, 1
	LONG $0x197de3c4; WORD $0x01f5       // vextractf128    xmm5, ymm6, 1
	LONG $0x3751e2c4; BYTE $0xc0         // vpcmpgtq    xmm0, xmm5, xmm0
	LONG $0x3749e2c4; BYTE $0xea         // vpcmpgtq    xmm5, xmm6, xmm2
	LONG $0x1855e3c4; WORD $0x01c0       // vinsertf128    ymm0, ymm5, xmm0, 1
	LONG $0x4b6de3c4; WORD $0x00d6       // vblendvpd    ymm2, ymm2, ymm6, ymm0
	LONG $0x197de3c4; WORD $0x01c8       // vextractf128    xmm0, ymm1, 1
	LONG $0x197de3c4; WORD $0x01e5       // vextractf128    xmm5, ymm4, 1
	LONG $0x3751e2c4; BYTE $0xc0         // vpcmpgtq    xmm0, xmm5, xmm0
	LONG $0x3759e2c4; BYTE $0xe9         // vpcmpgtq    xmm5, xmm4, xmm1
	LONG $0x1855e3c4; WORD $0x01c0       // vinsertf128    ymm0, ymm5, xmm0, 1
	LONG $0x4b7563c4; WORD $0x00c4       // vblendvpd    ymm8, ymm1, ymm4, ymm0
	LONG $0x20c08348                     // add    rax, 32
	LONG $0x02c28348                     // add    rdx, 2
	JNE  LBB1_8
	WORD $0x854d; BYTE $0xc0             // test    r8, r8
	JE   LBB1_11

LBB1_10:
	LONG $0x446ffec5; WORD $0x60c7 // vmovdqu    ymm0, yword [rdi + 8*rax + 96]
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197d63c4; WORD $0x01c4 // vextractf128    xmm4, ymm8, 1
	LONG $0x3771e2c4; BYTE $0xcc   // vpcmpgtq    xmm1, xmm1, xmm4
	LONG $0x3779c2c4; BYTE $0xe0   // vpcmpgtq    xmm4, xmm0, xmm8
	LONG $0x185de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm4, xmm1, 1
	LONG $0x4b3d63c4; WORD $0x10c0 // vblendvpd    ymm8, ymm8, ymm0, ymm1
	LONG $0x446ffec5; WORD $0x40c7 // vmovdqu    ymm0, yword [rdi + 8*rax + 64]
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197de3c4; WORD $0x01d4 // vextractf128    xmm4, ymm2, 1
	LONG $0x3771e2c4; BYTE $0xcc   // vpcmpgtq    xmm1, xmm1, xmm4
	LONG $0x3779e2c4; BYTE $0xe2   // vpcmpgtq    xmm4, xmm0, xmm2
	LONG $0x185de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm4, xmm1, 1
	LONG $0x4b6de3c4; WORD $0x10d0 // vblendvpd    ymm2, ymm2, ymm0, ymm1
	LONG $0x446ffec5; WORD $0x20c7 // vmovdqu    ymm0, yword [rdi + 8*rax + 32]
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197de3c4; WORD $0x01dc // vextractf128    xmm4, ymm3, 1
	LONG $0x3771e2c4; BYTE $0xcc   // vpcmpgtq    xmm1, xmm1, xmm4
	LONG $0x3779e2c4; BYTE $0xe3   // vpcmpgtq    xmm4, xmm0, xmm3
	LONG $0x185de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm4, xmm1, 1
	LONG $0x4b65e3c4; WORD $0x10d8 // vblendvpd    ymm3, ymm3, ymm0, ymm1
	LONG $0x046ffec5; BYTE $0xc7   // vmovdqu    ymm0, yword [rdi + 8*rax]
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197d63c4; WORD $0x01cc // vextractf128    xmm4, ymm9, 1
	LONG $0x3771e2c4; BYTE $0xcc   // vpcmpgtq    xmm1, xmm1, xmm4
	LONG $0x3779c2c4; BYTE $0xe1   // vpcmpgtq    xmm4, xmm0, xmm9
	LONG $0x185de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm4, xmm1, 1
	LONG $0x4b3563c4; WORD $0x10c8 // vblendvpd    ymm9, ymm9, ymm0, ymm1

LBB1_11:
	LONG $0x197de3c4; WORD $0x01d8 // vextractf128    xmm0, ymm3, 1
	LONG $0x197d63c4; WORD $0x01c9 // vextractf128    xmm1, ymm9, 1
	LONG $0x3771e2c4; BYTE $0xc0   // vpcmpgtq    xmm0, xmm1, xmm0
	LONG $0x3731e2c4; BYTE $0xcb   // vpcmpgtq    xmm1, xmm9, xmm3
	LONG $0x1875e3c4; WORD $0x01c0 // vinsertf128    ymm0, ymm1, xmm0, 1
	LONG $0x4b65c3c4; WORD $0x00c1 // vblendvpd    ymm0, ymm3, ymm9, ymm0
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197de3c4; WORD $0x01d3 // vextractf128    xmm3, ymm2, 1
	LONG $0x3771e2c4; BYTE $0xcb   // vpcmpgtq    xmm1, xmm1, xmm3
	LONG $0x3779e2c4; BYTE $0xda   // vpcmpgtq    xmm3, xmm0, xmm2
	LONG $0x1865e3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm3, xmm1, 1
	LONG $0x4b6de3c4; WORD $0x10c0 // vblendvpd    ymm0, ymm2, ymm0, ymm1
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197d63c4; WORD $0x01c2 // vextractf128    xmm2, ymm8, 1
	LONG $0x3771e2c4; BYTE $0xca   // vpcmpgtq    xmm1, xmm1, xmm2
	LONG $0x3779c2c4; BYTE $0xd0   // vpcmpgtq    xmm2, xmm0, xmm8
	LONG $0x186de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm2, xmm1, 1
	LONG $0x4b3de3c4; WORD $0x10c0 // vblendvpd    ymm0, ymm8, ymm0, ymm1
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x3779e2c4; BYTE $0xd1   // vpcmpgtq    xmm2, xmm0, xmm1
	LONG $0x3771e2c4; BYTE $0xd8   // vpcmpgtq    xmm3, xmm1, xmm0
	LONG $0x186de3c4; WORD $0x01d3 // vinsertf128    ymm2, ymm2, xmm3, 1
	LONG $0x4b75e3c4; WORD $0x20c0 // vblendvpd    ymm0, ymm1, ymm0, ymm2
	LONG $0x0479e3c4; WORD $0x4ec8 // vpermilps    xmm1, xmm0, 78
	LONG $0x3779e2c4; BYTE $0xd1   // vpcmpgtq    xmm2, xmm0, xmm1
	LONG $0x197de3c4; WORD $0x01c3 // vextractf128    xmm3, ymm0, 1
	LONG $0x3761e2c4; BYTE $0xd8   // vpcmpgtq    xmm3, xmm3, xmm0
	LONG $0x186de3c4; WORD $0x01d3 // vinsertf128    ymm2, ymm2, xmm3, 1
	LONG $0x4b75e3c4; WORD $0x20c0 // vblendvpd    ymm0, ymm1, ymm0, ymm2
	LONG $0x7ef9e1c4; BYTE $0xc0   // vmovq    rax, xmm0
	WORD $0x3948; BYTE $0xf1       // cmp    rcx, rsi
	JE   LBB1_13

LBB1_12:
	LONG $0xcf148b48         // mov    rdx, qword [rdi + 8*rcx]
	WORD $0x3948; BYTE $0xc2 // cmp    rdx, rax
	LONG $0xc24d0f48         // cmovge    rax, rdx
	LONG $0x01c18348         // add    rcx, 1
	WORD $0x3948; BYTE $0xce // cmp    rsi, rcx
	JNE  LBB1_12

LBB1_13:
	BYTE $0xc5; BYTE $0xf8; BYTE $0x77 // VZEROUPPER
	MOVQ AX, x+24(FP)
	RET

LBB1_6:
	LONG $0x4d6f7dc5; BYTE $0x00 // vmovdqa    ymm9, yword 0[rbp] /* [rip + .LCPI1_0] */
	WORD $0xc031                 // xor    eax, eax
	LONG $0x6f7dc1c4; BYTE $0xd9 // vmovdqa    ymm3, ymm9
	LONG $0x6f7dc1c4; BYTE $0xd1 // vmovdqa    ymm2, ymm9
	LONG $0x6f7d41c4; BYTE $0xc1 // vmovdqa    ymm8, ymm9
	WORD $0x854d; BYTE $0xc0     // test    r8, r8
	JNE  LBB1_10
	JMP  LBB1_11

TEXT sampleMaxSSE42<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX
	LEAQ LCDATA1<>(SB), BP

	WORD $0x8548; BYTE $0xf6               // test    rsi, rsi
	JE   LBB1_1
	QUAD $0x000000000000b848; WORD $0x8000 // mov    rax, -9223372036854775808
	JLE  LBB1_13
	LONG $0x04fe8348                       // cmp    rsi, 4
	JAE  LBB1_5
	WORD $0xc931                           // xor    ecx, ecx
	JMP  LBB1_12

LBB1_1:
	WORD $0xc031 // xor    eax, eax
	JMP  LBB1_13

LBB1_5:
	WORD $0x8948; BYTE $0xf1     // mov    rcx, rsi
	LONG $0xfce18348             // and    rcx, -4
	LONG $0xfc518d48             // lea    rdx, [rcx - 4]
	WORD $0x8948; BYTE $0xd0     // mov    rax, rdx
	LONG $0x02e8c148             // shr    rax, 2
	LONG $0x01c08348             // add    rax, 1
	WORD $0x8941; BYTE $0xc0     // mov    r8d, eax
	LONG $0x01e08341             // and    r8d, 1
	WORD $0x8548; BYTE $0xd2     // test    rdx, rdx
	JE   LBB1_6
	LONG $0x000001ba; BYTE $0x00 // mov    edx, 1
	WORD $0x2948; BYTE $0xc2     // sub    rdx, rax
	LONG $0x10048d49             // lea    rax, [r8 + rdx]
	LONG $0xffc08348             // add    rax, -1
	LONG $0x4d6f0f66; BYTE $0x00 // movdqa    xmm1, oword 0[rbp] /* [rip + .LCPI1_0] */
	WORD $0xd231                 // xor    edx, edx
	LONG $0xd16f0f66             // movdqa    xmm2, xmm1

LBB1_8:
	LONG $0x1c6f0ff3; BYTE $0xd7   // movdqu    xmm3, oword [rdi + 8*rdx]
	LONG $0x646f0ff3; WORD $0x10d7 // movdqu    xmm4, oword [rdi + 8*rdx + 16]
	LONG $0x6c6f0ff3; WORD $0x20d7 // movdqu    xmm5, oword [rdi + 8*rdx + 32]
	LONG $0xc36f0f66               // movdqa    xmm0, xmm3
	LONG $0x37380f66; BYTE $0xc1   // pcmpgtq    xmm0, xmm1
	LONG $0x15380f66; BYTE $0xcb   // blendvpd    xmm1, xmm3, xmm0
	LONG $0x5c6f0ff3; WORD $0x30d7 // movdqu    xmm3, oword [rdi + 8*rdx + 48]
	LONG $0xc46f0f66               // movdqa    xmm0, xmm4
	LONG $0x37380f66; BYTE $0xc2   // pcmpgtq    xmm0, xmm2
	LONG $0x15380f66; BYTE $0xd4   // blendvpd    xmm2, xmm4, xmm0
	LONG $0xc56f0f66               // movdqa    xmm0, xmm5
	LONG $0x37380f66; BYTE $0xc1   // pcmpgtq    xmm0, xmm1
	LONG $0x15380f66; BYTE $0xcd   // blendvpd    xmm1, xmm5, xmm0
	LONG $0xc36f0f66               // movdqa    xmm0, xmm3
	LONG $0x37380f66; BYTE $0xc2   // pcmpgtq    xmm0, xmm2
	LONG $0x15380f66; BYTE $0xd3   // blendvpd    xmm2, xmm3, xmm0
	LONG $0x08c28348               // add    rdx, 8
	LONG $0x02c08348               // add    rax, 2
	JNE  LBB1_8
	WORD $0x854d; BYTE $0xc0       // test    r8, r8
	JE   LBB1_11

LBB1_10:
	LONG $0x5c6f0ff3; WORD $0x10d7 // movdqu    xmm3, oword [rdi + 8*rdx + 16]
	LONG $0xc36f0f66               // movdqa    xmm0, xmm3
	LONG $0x37380f66; BYTE $0xc2   // pcmpgtq    xmm0, xmm2
	LONG $0x15380f66; BYTE $0xd3   // blendvpd    xmm2, xmm3, xmm0
	LONG $0x1c6f0ff3; BYTE $0xd7   // movdqu    xmm3, oword [rdi + 8*rdx]
	LONG $0xc36f0f66               // movdqa    xmm0, xmm3
	LONG $0x37380f66; BYTE $0xc1   // pcmpgtq    xmm0, xmm1
	LONG $0x15380f66; BYTE $0xcb   // blendvpd    xmm1, xmm3, xmm0

LBB1_11:
	LONG $0xc16f0f66             // movdqa    xmm0, xmm1
	LONG $0x37380f66; BYTE $0xc2 // pcmpgtq    xmm0, xmm2
	LONG $0x15380f66; BYTE $0xd1 // blendvpd    xmm2, xmm1, xmm0
	LONG $0xca700f66; BYTE $0x4e // pshufd    xmm1, xmm2, 78
	LONG $0xc26f0f66             // movdqa    xmm0, xmm2
	LONG $0x37380f66; BYTE $0xc1 // pcmpgtq    xmm0, xmm1
	LONG $0x15380f66; BYTE $0xca // blendvpd    xmm1, xmm2, xmm0
	LONG $0x7e0f4866; BYTE $0xc8 // movq    rax, xmm1
	WORD $0x3948; BYTE $0xf1     // cmp    rcx, rsi
	JE   LBB1_13

LBB1_12:
	LONG $0xcf148b48         // mov    rdx, qword [rdi + 8*rcx]
	WORD $0x3948; BYTE $0xc2 // cmp    rdx, rax
	LONG $0xc24d0f48         // cmovge    rax, rdx
	LONG $0x01c18348         // add    rcx, 1
	WORD $0x3948; BYTE $0xce // cmp    rsi, rcx
	JNE  LBB1_12

LBB1_13:
	MOVQ AX, x+24(FP)
	RET

LBB1_6:
	LONG $0x4d6f0f66; BYTE $0x00 // movdqa    xmm1, oword 0[rbp] /* [rip + .LCPI1_0] */
	WORD $0xd231                 // xor    edx, edx
	LONG $0xd16f0f66             // movdqa    xmm2, xmm1
	WORD $0x854d; BYTE $0xc0     // test    r8, r8
	JNE  LBB1_10
	JMP  LBB1_11

DATA LCDATA2<>+0x000(SB)/8, $0x7fffffffffffffff
DATA LCDATA2<>+0x008(SB)/8, $0x7fffffffffffffff
DATA LCDATA2<>+0x010(SB)/8, $0x7fffffffffffffff
DATA LCDATA2<>+0x018(SB)/8, $0x7fffffffffffffff
GLOBL LCDATA2<>(SB), 8, $32

TEXT sampleMinAVX2<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX
	LEAQ LCDATA2<>(SB), BP

	WORD $0x8548; BYTE $0xf6               // test    rsi, rsi
	JE   LBB2_1
	QUAD $0xffffffffffffb848; WORD $0x7fff // mov    rax, 9223372036854775807
	JLE  LBB2_13
	LONG $0x10fe8348                       // cmp    rsi, 16
	JAE  LBB2_5
	WORD $0xc931                           // xor    ecx, ecx
	JMP  LBB2_12

LBB2_1:
	WORD $0xc031      // xor    eax, eax
	MOVQ AX, x+24(FP)
	RET

LBB2_5:
	WORD $0x8948; BYTE $0xf1       // mov    rcx, rsi
	LONG $0xf0e18348               // and    rcx, -16
	LONG $0xf0518d48               // lea    rdx, [rcx - 16]
	WORD $0x8948; BYTE $0xd0       // mov    rax, rdx
	LONG $0x04e8c148               // shr    rax, 4
	LONG $0x01c08348               // add    rax, 1
	WORD $0x8941; BYTE $0xc0       // mov    r8d, eax
	LONG $0x01e08341               // and    r8d, 1
	WORD $0x8548; BYTE $0xd2       // test    rdx, rdx
	JE   LBB2_6
	LONG $0x000001ba; BYTE $0x00   // mov    edx, 1
	WORD $0x2948; BYTE $0xc2       // sub    rdx, rax
	WORD $0x014c; BYTE $0xc2       // add    rdx, r8
	LONG $0xffc28348               // add    rdx, -1
	LONG $0x597de2c4; WORD $0x0045 // vpbroadcastq    ymm0, qword 0[rbp] /* [rip + .LCPI2_0] */
	WORD $0xc031                   // xor    eax, eax
	LONG $0xc86ffdc5               // vmovdqa    ymm1, ymm0
	LONG $0xd06ffdc5               // vmovdqa    ymm2, ymm0
	LONG $0xd86ffdc5               // vmovdqa    ymm3, ymm0

LBB2_8:
	LONG $0x246ffec5; BYTE $0xc7         // vmovdqu    ymm4, yword [rdi + 8*rax]
	LONG $0x6c6ffec5; WORD $0x20c7       // vmovdqu    ymm5, yword [rdi + 8*rax + 32]
	LONG $0x746ffec5; WORD $0x40c7       // vmovdqu    ymm6, yword [rdi + 8*rax + 64]
	LONG $0x377de2c4; BYTE $0xfc         // vpcmpgtq    ymm7, ymm0, ymm4
	LONG $0x4b7de3c4; WORD $0x70c4       // vblendvpd    ymm0, ymm0, ymm4, ymm7
	LONG $0x646ffec5; WORD $0x60c7       // vmovdqu    ymm4, yword [rdi + 8*rax + 96]
	LONG $0x3775e2c4; BYTE $0xfd         // vpcmpgtq    ymm7, ymm1, ymm5
	LONG $0x4b75e3c4; WORD $0x70cd       // vblendvpd    ymm1, ymm1, ymm5, ymm7
	LONG $0x376de2c4; BYTE $0xee         // vpcmpgtq    ymm5, ymm2, ymm6
	LONG $0x4b6de3c4; WORD $0x50d6       // vblendvpd    ymm2, ymm2, ymm6, ymm5
	LONG $0x3765e2c4; BYTE $0xec         // vpcmpgtq    ymm5, ymm3, ymm4
	LONG $0x4b65e3c4; WORD $0x50dc       // vblendvpd    ymm3, ymm3, ymm4, ymm5
	QUAD $0x000080c7a46ffec5; BYTE $0x00 // vmovdqu    ymm4, yword [rdi + 8*rax + 128]
	QUAD $0x0000a0c7ac6ffec5; BYTE $0x00 // vmovdqu    ymm5, yword [rdi + 8*rax + 160]
	QUAD $0x0000c0c7b46ffec5; BYTE $0x00 // vmovdqu    ymm6, yword [rdi + 8*rax + 192]
	LONG $0x377de2c4; BYTE $0xfc         // vpcmpgtq    ymm7, ymm0, ymm4
	LONG $0x4b7de3c4; WORD $0x70c4       // vblendvpd    ymm0, ymm0, ymm4, ymm7
	QUAD $0x0000e0c7a46ffec5; BYTE $0x00 // vmovdqu    ymm4, yword [rdi + 8*rax + 224]
	LONG $0x3775e2c4; BYTE $0xfd         // vpcmpgtq    ymm7, ymm1, ymm5
	LONG $0x4b75e3c4; WORD $0x70cd       // vblendvpd    ymm1, ymm1, ymm5, ymm7
	LONG $0x376de2c4; BYTE $0xee         // vpcmpgtq    ymm5, ymm2, ymm6
	LONG $0x4b6de3c4; WORD $0x50d6       // vblendvpd    ymm2, ymm2, ymm6, ymm5
	LONG $0x3765e2c4; BYTE $0xec         // vpcmpgtq    ymm5, ymm3, ymm4
	LONG $0x4b65e3c4; WORD $0x50dc       // vblendvpd    ymm3, ymm3, ymm4, ymm5
	LONG $0x20c08348                     // add    rax, 32
	LONG $0x02c28348                     // add    rdx, 2
	JNE  LBB2_8
	WORD $0x854d; BYTE $0xc0             // test    r8, r8
	JE   LBB2_11

LBB2_10:
	LONG $0x646ffec5; WORD $0x60c7 // vmovdqu    ymm4, yword [rdi + 8*rax + 96]
	LONG $0x3765e2c4; BYTE $0xec   // vpcmpgtq    ymm5, ymm3, ymm4
	LONG $0x4b65e3c4; WORD $0x50dc // vblendvpd    ymm3, ymm3, ymm4, ymm5
	LONG $0x646ffec5; WORD $0x40c7 // vmovdqu    ymm4, yword [rdi + 8*rax + 64]
	LONG $0x376de2c4; BYTE $0xec   // vpcmpgtq    ymm5, ymm2, ymm4
	LONG $0x4b6de3c4; WORD $0x50d4 // vblendvpd    ymm2, ymm2, ymm4, ymm5
	LONG $0x646ffec5; WORD $0x20c7 // vmovdqu    ymm4, yword [rdi + 8*rax + 32]
	LONG $0x3775e2c4; BYTE $0xec   // vpcmpgtq    ymm5, ymm1, ymm4
	LONG $0x4b75e3c4; WORD $0x50cc // vblendvpd    ymm1, ymm1, ymm4, ymm5
	LONG $0x246ffec5; BYTE $0xc7   // vmovdqu    ymm4, yword [rdi + 8*rax]
	LONG $0x377de2c4; BYTE $0xec   // vpcmpgtq    ymm5, ymm0, ymm4
	LONG $0x4b7de3c4; WORD $0x50c4 // vblendvpd    ymm0, ymm0, ymm4, ymm5

LBB2_11:
	LONG $0x3775e2c4; BYTE $0xe0   // vpcmpgtq    ymm4, ymm1, ymm0
	LONG $0x4b75e3c4; WORD $0x40c0 // vblendvpd    ymm0, ymm1, ymm0, ymm4
	LONG $0x376de2c4; BYTE $0xc8   // vpcmpgtq    ymm1, ymm2, ymm0
	LONG $0x4b6de3c4; WORD $0x10c0 // vblendvpd    ymm0, ymm2, ymm0, ymm1
	LONG $0x3765e2c4; BYTE $0xc8   // vpcmpgtq    ymm1, ymm3, ymm0
	LONG $0x4b65e3c4; WORD $0x10c0 // vblendvpd    ymm0, ymm3, ymm0, ymm1
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x3775e2c4; BYTE $0xd0   // vpcmpgtq    ymm2, ymm1, ymm0
	LONG $0x4b75e3c4; WORD $0x20c0 // vblendvpd    ymm0, ymm1, ymm0, ymm2
	LONG $0x0479e3c4; WORD $0x4ec8 // vpermilps    xmm1, xmm0, 78
	LONG $0x3775e2c4; BYTE $0xd0   // vpcmpgtq    ymm2, ymm1, ymm0
	LONG $0x4b75e3c4; WORD $0x20c0 // vblendvpd    ymm0, ymm1, ymm0, ymm2
	LONG $0x7ef9e1c4; BYTE $0xc0   // vmovq    rax, xmm0
	WORD $0x3948; BYTE $0xf1       // cmp    rcx, rsi
	JE   LBB2_13

LBB2_12:
	LONG $0xcf148b48         // mov    rdx, qword [rdi + 8*rcx]
	WORD $0x3948; BYTE $0xc2 // cmp    rdx, rax
	LONG $0xc24e0f48         // cmovle    rax, rdx
	LONG $0x01c18348         // add    rcx, 1
	WORD $0x3948; BYTE $0xce // cmp    rsi, rcx
	JNE  LBB2_12

LBB2_13:
	BYTE $0xc5; BYTE $0xf8; BYTE $0x77 // VZEROUPPER
	MOVQ AX, x+24(FP)
	RET

LBB2_6:
	LONG $0x597de2c4; WORD $0x0045 // vpbroadcastq    ymm0, qword 0[rbp] /* [rip + .LCPI2_0] */
	WORD $0xc031                   // xor    eax, eax
	LONG $0xc86ffdc5               // vmovdqa    ymm1, ymm0
	LONG $0xd06ffdc5               // vmovdqa    ymm2, ymm0
	LONG $0xd86ffdc5               // vmovdqa    ymm3, ymm0
	WORD $0x854d; BYTE $0xc0       // test    r8, r8
	JNE  LBB2_10
	JMP  LBB2_11

TEXT sampleMinAVX<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX
	LEAQ LCDATA2<>(SB), BP

	WORD $0x8548; BYTE $0xf6               // test    rsi, rsi
	JE   LBB2_1
	QUAD $0xffffffffffffb848; WORD $0x7fff // mov    rax, 9223372036854775807
	JLE  LBB2_13
	LONG $0x10fe8348                       // cmp    rsi, 16
	JAE  LBB2_5
	WORD $0xc931                           // xor    ecx, ecx
	JMP  LBB2_12

LBB2_1:
	WORD $0xc031      // xor    eax, eax
	MOVQ AX, x+24(FP)
	RET

LBB2_5:
	WORD $0x8948; BYTE $0xf1     // mov    rcx, rsi
	LONG $0xf0e18348             // and    rcx, -16
	LONG $0xf0518d48             // lea    rdx, [rcx - 16]
	WORD $0x8948; BYTE $0xd0     // mov    rax, rdx
	LONG $0x04e8c148             // shr    rax, 4
	LONG $0x01c08348             // add    rax, 1
	WORD $0x8941; BYTE $0xc0     // mov    r8d, eax
	LONG $0x01e08341             // and    r8d, 1
	WORD $0x8548; BYTE $0xd2     // test    rdx, rdx
	JE   LBB2_6
	LONG $0x000001ba; BYTE $0x00 // mov    edx, 1
	WORD $0x2948; BYTE $0xc2     // sub    rdx, rax
	WORD $0x014c; BYTE $0xc2     // add    rdx, r8
	LONG $0xffc28348             // add    rdx, -1
	LONG $0x4d6f7dc5; BYTE $0x00 // vmovdqa    ymm9, yword 0[rbp] /* [rip + .LCPI2_0] */
	WORD $0xc031                 // xor    eax, eax
	LONG $0x6f7dc1c4; BYTE $0xd9 // vmovdqa    ymm3, ymm9
	LONG $0x6f7dc1c4; BYTE $0xd1 // vmovdqa    ymm2, ymm9
	LONG $0x6f7d41c4; BYTE $0xc1 // vmovdqa    ymm8, ymm9

LBB2_8:
	LONG $0x246ffec5; BYTE $0xc7         // vmovdqu    ymm4, yword [rdi + 8*rax]
	LONG $0x6c6ffec5; WORD $0x20c7       // vmovdqu    ymm5, yword [rdi + 8*rax + 32]
	LONG $0x746ffec5; WORD $0x40c7       // vmovdqu    ymm6, yword [rdi + 8*rax + 64]
	LONG $0x197de3c4; WORD $0x01e7       // vextractf128    xmm7, ymm4, 1
	LONG $0x197d63c4; WORD $0x01c9       // vextractf128    xmm1, ymm9, 1
	LONG $0x3771e2c4; BYTE $0xcf         // vpcmpgtq    xmm1, xmm1, xmm7
	LONG $0x3731e2c4; BYTE $0xfc         // vpcmpgtq    xmm7, xmm9, xmm4
	LONG $0x1845e3c4; WORD $0x01c9       // vinsertf128    ymm1, ymm7, xmm1, 1
	LONG $0x4b3563c4; WORD $0x10cc       // vblendvpd    ymm9, ymm9, ymm4, ymm1
	LONG $0x4c6ffec5; WORD $0x60c7       // vmovdqu    ymm1, yword [rdi + 8*rax + 96]
	LONG $0x197de3c4; WORD $0x01ec       // vextractf128    xmm4, ymm5, 1
	LONG $0x197de3c4; WORD $0x01df       // vextractf128    xmm7, ymm3, 1
	LONG $0x3741e2c4; BYTE $0xe4         // vpcmpgtq    xmm4, xmm7, xmm4
	LONG $0x3761e2c4; BYTE $0xfd         // vpcmpgtq    xmm7, xmm3, xmm5
	LONG $0x1845e3c4; WORD $0x01e4       // vinsertf128    ymm4, ymm7, xmm4, 1
	LONG $0x4b65e3c4; WORD $0x40dd       // vblendvpd    ymm3, ymm3, ymm5, ymm4
	LONG $0x197de3c4; WORD $0x01f4       // vextractf128    xmm4, ymm6, 1
	LONG $0x197de3c4; WORD $0x01d5       // vextractf128    xmm5, ymm2, 1
	LONG $0x3751e2c4; BYTE $0xe4         // vpcmpgtq    xmm4, xmm5, xmm4
	LONG $0x3769e2c4; BYTE $0xee         // vpcmpgtq    xmm5, xmm2, xmm6
	LONG $0x1855e3c4; WORD $0x01e4       // vinsertf128    ymm4, ymm5, xmm4, 1
	LONG $0x4b6de3c4; WORD $0x40d6       // vblendvpd    ymm2, ymm2, ymm6, ymm4
	LONG $0x197de3c4; WORD $0x01cc       // vextractf128    xmm4, ymm1, 1
	LONG $0x197d63c4; WORD $0x01c5       // vextractf128    xmm5, ymm8, 1
	LONG $0x3751e2c4; BYTE $0xe4         // vpcmpgtq    xmm4, xmm5, xmm4
	LONG $0x3739e2c4; BYTE $0xe9         // vpcmpgtq    xmm5, xmm8, xmm1
	LONG $0x1855e3c4; WORD $0x01e4       // vinsertf128    ymm4, ymm5, xmm4, 1
	LONG $0x4b3de3c4; WORD $0x40c9       // vblendvpd    ymm1, ymm8, ymm1, ymm4
	QUAD $0x000080c7a46ffec5; BYTE $0x00 // vmovdqu    ymm4, yword [rdi + 8*rax + 128]
	QUAD $0x0000a0c7ac6ffec5; BYTE $0x00 // vmovdqu    ymm5, yword [rdi + 8*rax + 160]
	QUAD $0x0000c0c7b46ffec5; BYTE $0x00 // vmovdqu    ymm6, yword [rdi + 8*rax + 192]
	LONG $0x197d63c4; WORD $0x01cf       // vextractf128    xmm7, ymm9, 1
	LONG $0x197de3c4; WORD $0x01e0       // vextractf128    xmm0, ymm4, 1
	LONG $0x3741e2c4; BYTE $0xc0         // vpcmpgtq    xmm0, xmm7, xmm0
	LONG $0x3731e2c4; BYTE $0xfc         // vpcmpgtq    xmm7, xmm9, xmm4
	LONG $0x1845e3c4; WORD $0x01c0       // vinsertf128    ymm0, ymm7, xmm0, 1
	LONG $0x4b3563c4; WORD $0x00cc       // vblendvpd    ymm9, ymm9, ymm4, ymm0
	QUAD $0x0000e0c7a46ffec5; BYTE $0x00 // vmovdqu    ymm4, yword [rdi + 8*rax + 224]
	LONG $0x197de3c4; WORD $0x01df       // vextractf128    xmm7, ymm3, 1
	LONG $0x197de3c4; WORD $0x01e8       // vextractf128    xmm0, ymm5, 1
	LONG $0x3741e2c4; BYTE $0xc0         // vpcmpgtq    xmm0, xmm7, xmm0
	LONG $0x3761e2c4; BYTE $0xfd         // vpcmpgtq    xmm7, xmm3, xmm5
	LONG $0x1845e3c4; WORD $0x01c0       // vinsertf128    ymm0, ymm7, xmm0, 1
	LONG $0x4b65e3c4; WORD $0x00dd       // vblendvpd    ymm3, ymm3, ymm5, ymm0
	LONG $0x197de3c4; WORD $0x01d0       // vextractf128    xmm0, ymm2, 1
	LONG $0x197de3c4; WORD $0x01f5       // vextractf128    xmm5, ymm6, 1
	LONG $0x3779e2c4; BYTE $0xc5         // vpcmpgtq    xmm0, xmm0, xmm5
	LONG $0x3769e2c4; BYTE $0xee         // vpcmpgtq    xmm5, xmm2, xmm6
	LONG $0x1855e3c4; WORD $0x01c0       // vinsertf128    ymm0, ymm5, xmm0, 1
	LONG $0x4b6de3c4; WORD $0x00d6       // vblendvpd    ymm2, ymm2, ymm6, ymm0
	LONG $0x197de3c4; WORD $0x01c8       // vextractf128    xmm0, ymm1, 1
	LONG $0x197de3c4; WORD $0x01e5       // vextractf128    xmm5, ymm4, 1
	LONG $0x3779e2c4; BYTE $0xc5         // vpcmpgtq    xmm0, xmm0, xmm5
	LONG $0x3771e2c4; BYTE $0xec         // vpcmpgtq    xmm5, xmm1, xmm4
	LONG $0x1855e3c4; WORD $0x01c0       // vinsertf128    ymm0, ymm5, xmm0, 1
	LONG $0x4b7563c4; WORD $0x00c4       // vblendvpd    ymm8, ymm1, ymm4, ymm0
	LONG $0x20c08348                     // add    rax, 32
	LONG $0x02c28348                     // add    rdx, 2
	JNE  LBB2_8
	WORD $0x854d; BYTE $0xc0             // test    r8, r8
	JE   LBB2_11

LBB2_10:
	LONG $0x446ffec5; WORD $0x60c7 // vmovdqu    ymm0, yword [rdi + 8*rax + 96]
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197d63c4; WORD $0x01c4 // vextractf128    xmm4, ymm8, 1
	LONG $0x3759e2c4; BYTE $0xc9   // vpcmpgtq    xmm1, xmm4, xmm1
	LONG $0x3739e2c4; BYTE $0xe0   // vpcmpgtq    xmm4, xmm8, xmm0
	LONG $0x185de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm4, xmm1, 1
	LONG $0x4b3d63c4; WORD $0x10c0 // vblendvpd    ymm8, ymm8, ymm0, ymm1
	LONG $0x446ffec5; WORD $0x40c7 // vmovdqu    ymm0, yword [rdi + 8*rax + 64]
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197de3c4; WORD $0x01d4 // vextractf128    xmm4, ymm2, 1
	LONG $0x3759e2c4; BYTE $0xc9   // vpcmpgtq    xmm1, xmm4, xmm1
	LONG $0x3769e2c4; BYTE $0xe0   // vpcmpgtq    xmm4, xmm2, xmm0
	LONG $0x185de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm4, xmm1, 1
	LONG $0x4b6de3c4; WORD $0x10d0 // vblendvpd    ymm2, ymm2, ymm0, ymm1
	LONG $0x446ffec5; WORD $0x20c7 // vmovdqu    ymm0, yword [rdi + 8*rax + 32]
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197de3c4; WORD $0x01dc // vextractf128    xmm4, ymm3, 1
	LONG $0x3759e2c4; BYTE $0xc9   // vpcmpgtq    xmm1, xmm4, xmm1
	LONG $0x3761e2c4; BYTE $0xe0   // vpcmpgtq    xmm4, xmm3, xmm0
	LONG $0x185de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm4, xmm1, 1
	LONG $0x4b65e3c4; WORD $0x10d8 // vblendvpd    ymm3, ymm3, ymm0, ymm1
	LONG $0x046ffec5; BYTE $0xc7   // vmovdqu    ymm0, yword [rdi + 8*rax]
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197d63c4; WORD $0x01cc // vextractf128    xmm4, ymm9, 1
	LONG $0x3759e2c4; BYTE $0xc9   // vpcmpgtq    xmm1, xmm4, xmm1
	LONG $0x3731e2c4; BYTE $0xe0   // vpcmpgtq    xmm4, xmm9, xmm0
	LONG $0x185de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm4, xmm1, 1
	LONG $0x4b3563c4; WORD $0x10c8 // vblendvpd    ymm9, ymm9, ymm0, ymm1

LBB2_11:
	LONG $0x197d63c4; WORD $0x01c8 // vextractf128    xmm0, ymm9, 1
	LONG $0x197de3c4; WORD $0x01d9 // vextractf128    xmm1, ymm3, 1
	LONG $0x3771e2c4; BYTE $0xc0   // vpcmpgtq    xmm0, xmm1, xmm0
	LONG $0x3761c2c4; BYTE $0xc9   // vpcmpgtq    xmm1, xmm3, xmm9
	LONG $0x1875e3c4; WORD $0x01c0 // vinsertf128    ymm0, ymm1, xmm0, 1
	LONG $0x4b65c3c4; WORD $0x00c1 // vblendvpd    ymm0, ymm3, ymm9, ymm0
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197de3c4; WORD $0x01d3 // vextractf128    xmm3, ymm2, 1
	LONG $0x3761e2c4; BYTE $0xc9   // vpcmpgtq    xmm1, xmm3, xmm1
	LONG $0x3769e2c4; BYTE $0xd8   // vpcmpgtq    xmm3, xmm2, xmm0
	LONG $0x1865e3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm3, xmm1, 1
	LONG $0x4b6de3c4; WORD $0x10c0 // vblendvpd    ymm0, ymm2, ymm0, ymm1
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x197d63c4; WORD $0x01c2 // vextractf128    xmm2, ymm8, 1
	LONG $0x3769e2c4; BYTE $0xc9   // vpcmpgtq    xmm1, xmm2, xmm1
	LONG $0x3739e2c4; BYTE $0xd0   // vpcmpgtq    xmm2, xmm8, xmm0
	LONG $0x186de3c4; WORD $0x01c9 // vinsertf128    ymm1, ymm2, xmm1, 1
	LONG $0x4b3de3c4; WORD $0x10c0 // vblendvpd    ymm0, ymm8, ymm0, ymm1
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128    xmm1, ymm0, 1
	LONG $0x3771e2c4; BYTE $0xd0   // vpcmpgtq    xmm2, xmm1, xmm0
	LONG $0x3779e2c4; BYTE $0xd9   // vpcmpgtq    xmm3, xmm0, xmm1
	LONG $0x186de3c4; WORD $0x01d3 // vinsertf128    ymm2, ymm2, xmm3, 1
	LONG $0x4b75e3c4; WORD $0x20c0 // vblendvpd    ymm0, ymm1, ymm0, ymm2
	LONG $0x0479e3c4; WORD $0x4ec8 // vpermilps    xmm1, xmm0, 78
	LONG $0x3771e2c4; BYTE $0xd0   // vpcmpgtq    xmm2, xmm1, xmm0
	LONG $0x197de3c4; WORD $0x01c3 // vextractf128    xmm3, ymm0, 1
	LONG $0x3779e2c4; BYTE $0xdb   // vpcmpgtq    xmm3, xmm0, xmm3
	LONG $0x186de3c4; WORD $0x01d3 // vinsertf128    ymm2, ymm2, xmm3, 1
	LONG $0x4b75e3c4; WORD $0x20c0 // vblendvpd    ymm0, ymm1, ymm0, ymm2
	LONG $0x7ef9e1c4; BYTE $0xc0   // vmovq    rax, xmm0
	WORD $0x3948; BYTE $0xf1       // cmp    rcx, rsi
	JE   LBB2_13

LBB2_12:
	LONG $0xcf148b48         // mov    rdx, qword [rdi + 8*rcx]
	WORD $0x3948; BYTE $0xc2 // cmp    rdx, rax
	LONG $0xc24e0f48         // cmovle    rax, rdx
	LONG $0x01c18348         // add    rcx, 1
	WORD $0x3948; BYTE $0xce // cmp    rsi, rcx
	JNE  LBB2_12

LBB2_13:
	BYTE $0xc5; BYTE $0xf8; BYTE $0x77 // VZEROUPPER
	MOVQ AX, x+24(FP)
	RET

LBB2_6:
	LONG $0x4d6f7dc5; BYTE $0x00 // vmovdqa    ymm9, yword 0[rbp] /* [rip + .LCPI2_0] */
	WORD $0xc031                 // xor    eax, eax
	LONG $0x6f7dc1c4; BYTE $0xd9 // vmovdqa    ymm3, ymm9
	LONG $0x6f7dc1c4; BYTE $0xd1 // vmovdqa    ymm2, ymm9
	LONG $0x6f7d41c4; BYTE $0xc1 // vmovdqa    ymm8, ymm9
	WORD $0x854d; BYTE $0xc0     // test    r8, r8
	JNE  LBB2_10
	JMP  LBB2_11

TEXT sampleMinSSE42<>(SB), NOSPLIT, $0-32

	MOVQ addr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ cap+16(FP), DX
	LEAQ LCDATA2<>(SB), BP

	WORD $0x8548; BYTE $0xf6               // test    rsi, rsi
	JE   LBB2_1
	QUAD $0xffffffffffffb848; WORD $0x7fff // mov    rax, 9223372036854775807
	JLE  LBB2_13
	LONG $0x04fe8348                       // cmp    rsi, 4
	JAE  LBB2_5
	WORD $0xc931                           // xor    ecx, ecx
	JMP  LBB2_12

LBB2_1:
	WORD $0xc031 // xor    eax, eax
	JMP  LBB2_13

LBB2_5:
	WORD $0x8948; BYTE $0xf1     // mov    rcx, rsi
	LONG $0xfce18348             // and    rcx, -4
	LONG $0xfc518d48             // lea    rdx, [rcx - 4]
	WORD $0x8948; BYTE $0xd0     // mov    rax, rdx
	LONG $0x02e8c148             // shr    rax, 2
	LONG $0x01c08348             // add    rax, 1
	WORD $0x8941; BYTE $0xc0     // mov    r8d, eax
	LONG $0x01e08341             // and    r8d, 1
	WORD $0x8548; BYTE $0xd2     // test    rdx, rdx
	JE   LBB2_6
	LONG $0x000001ba; BYTE $0x00 // mov    edx, 1
	WORD $0x2948; BYTE $0xc2     // sub    rdx, rax
	LONG $0x10048d49             // lea    rax, [r8 + rdx]
	LONG $0xffc08348             // add    rax, -1
	LONG $0x4d6f0f66; BYTE $0x00 // movdqa    xmm1, oword 0[rbp] /* [rip + .LCPI2_0] */
	WORD $0xd231                 // xor    edx, edx
	LONG $0xd16f0f66             // movdqa    xmm2, xmm1

LBB2_8:
	LONG $0x1c6f0ff3; BYTE $0xd7   // movdqu    xmm3, oword [rdi + 8*rdx]
	LONG $0x646f0ff3; WORD $0x10d7 // movdqu    xmm4, oword [rdi + 8*rdx + 16]
	LONG $0x6c6f0ff3; WORD $0x20d7 // movdqu    xmm5, oword [rdi + 8*rdx + 32]
	LONG $0xc16f0f66               // movdqa    xmm0, xmm1
	LONG $0x37380f66; BYTE $0xc3   // pcmpgtq    xmm0, xmm3
	LONG $0x15380f66; BYTE $0xcb   // blendvpd    xmm1, xmm3, xmm0
	LONG $0x5c6f0ff3; WORD $0x30d7 // movdqu    xmm3, oword [rdi + 8*rdx + 48]
	LONG $0xc26f0f66               // movdqa    xmm0, xmm2
	LONG $0x37380f66; BYTE $0xc4   // pcmpgtq    xmm0, xmm4
	LONG $0x15380f66; BYTE $0xd4   // blendvpd    xmm2, xmm4, xmm0
	LONG $0xc1280f66               // movapd    xmm0, xmm1
	LONG $0x37380f66; BYTE $0xc5   // pcmpgtq    xmm0, xmm5
	LONG $0x15380f66; BYTE $0xcd   // blendvpd    xmm1, xmm5, xmm0
	LONG $0xc2280f66               // movapd    xmm0, xmm2
	LONG $0x37380f66; BYTE $0xc3   // pcmpgtq    xmm0, xmm3
	LONG $0x15380f66; BYTE $0xd3   // blendvpd    xmm2, xmm3, xmm0
	LONG $0x08c28348               // add    rdx, 8
	LONG $0x02c08348               // add    rax, 2
	JNE  LBB2_8
	WORD $0x854d; BYTE $0xc0       // test    r8, r8
	JE   LBB2_11

LBB2_10:
	LONG $0x5c6f0ff3; WORD $0x10d7 // movdqu    xmm3, oword [rdi + 8*rdx + 16]
	LONG $0xc26f0f66               // movdqa    xmm0, xmm2
	LONG $0x37380f66; BYTE $0xc3   // pcmpgtq    xmm0, xmm3
	LONG $0x15380f66; BYTE $0xd3   // blendvpd    xmm2, xmm3, xmm0
	LONG $0x1c6f0ff3; BYTE $0xd7   // movdqu    xmm3, oword [rdi + 8*rdx]
	LONG $0xc16f0f66               // movdqa    xmm0, xmm1
	LONG $0x37380f66; BYTE $0xc3   // pcmpgtq    xmm0, xmm3
	LONG $0x15380f66; BYTE $0xcb   // blendvpd    xmm1, xmm3, xmm0

LBB2_11:
	LONG $0xc26f0f66             // movdqa    xmm0, xmm2
	LONG $0x37380f66; BYTE $0xc1 // pcmpgtq    xmm0, xmm1
	LONG $0x15380f66; BYTE $0xd1 // blendvpd    xmm2, xmm1, xmm0
	LONG $0xca700f66; BYTE $0x4e // pshufd    xmm1, xmm2, 78
	LONG $0xc16f0f66             // movdqa    xmm0, xmm1
	LONG $0x37380f66; BYTE $0xc2 // pcmpgtq    xmm0, xmm2
	LONG $0x15380f66; BYTE $0xca // blendvpd    xmm1, xmm2, xmm0
	LONG $0x7e0f4866; BYTE $0xc8 // movq    rax, xmm1
	WORD $0x3948; BYTE $0xf1     // cmp    rcx, rsi
	JE   LBB2_13

LBB2_12:
	LONG $0xcf148b48         // mov    rdx, qword [rdi + 8*rcx]
	WORD $0x3948; BYTE $0xc2 // cmp    rdx, rax
	LONG $0xc24e0f48         // cmovle    rax, rdx
	LONG $0x01c18348         // add    rcx, 1
	WORD $0x3948; BYTE $0xce // cmp    rsi, rcx
	JNE  LBB2_12

LBB2_13:
	MOVQ AX, x+24(FP)
	RET

LBB2_6:
	LONG $0x4d6f0f66; BYTE $0x00 // movdqa    xmm1, oword 0[rbp] /* [rip + .LCPI2_0] */
	WORD $0xd231                 // xor    edx, edx
	LONG $0xd16f0f66             // movdqa    xmm2, xmm1
	WORD $0x854d; BYTE $0xc0     // test    r8, r8
	JNE  LBB2_10
	JMP  LBB2_11
