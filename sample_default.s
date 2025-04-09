// +build !amd64,!gccgo

#include "textflag.h"

TEXT ·SampleSum(SB), NOSPLIT, $0-32
	JMP ·sampleSum(SB)

TEXT ·SampleVariance(SB), NOSPLIT, $0-32
	JMP ·sampleVariance(SB)

TEXT ·SampleMin(SB), NOSPLIT, $0-32
	JMP ·sampleMin(SB)

TEXT ·SampleMax(SB), NOSPLIT, $0-32
	JMP ·sampleMax(SB)
