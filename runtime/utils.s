TEXT ·add(SB), $0-24
    MOVQ x+0(FP), BX
    MOVQ y+8(FP), BP
    ADDQ BP, BX
    MOVQ BX, ret+16(FP)
    RET

TEXT ·sub(SB), $0-24
    MOVQ x+0(FP), BX
    MOVQ y+8(FP), BP
    SUBQ BP, BX
    MOVQ BX, ret+16(FP)
    RET

TEXT ·mul(SB), $0-24
    MOVQ x+0(FP), BX
    MOVQ y+8(FP), BP
    IMULQ BP, BX  // use IMUL due to int64
    MOVQ BX, ret+16(FP)
    RET
