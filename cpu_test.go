package main

import (
    "testing"
)

func NewTestCPU(ops ...uint16) *CPU {
    return NewCPU(build(ops...))
}

func Test_CLS(t *testing.T) {
    c := NewTestCPU(
        CLS(),
    )
    c.DisplayBuffer[0][0] = 1

    c.Cycle()

    if c.DisplayBuffer[0][0] != 0 {
        t.Error("display buffer was not cleared")
    }
}

func Test_CALL(t *testing.T) {
    c := NewTestCPU(
        CALL(0x204),
    )
    c.Cycle()

    if c.ProgramCounter != 0x204 {
        t.Error("program counter is wrong")
    }
    if c.Stack[c.StackPointer - 1] != 0x200 {
        t.Error("stack value is wrong")
    }
}

func Test_RET(t *testing.T) {
    c := NewTestCPU(
        CALL(0x204),
        NOP(),
        RET(),
    )
    c.Cycle()
    c.Cycle()

    if c.ProgramCounter != 0x202 {
        t.Errorf("program counter is wrong: %v", c.ProgramCounter)
    }
}

func Test_JP(t *testing.T) {
    c := NewTestCPU(
        JP(0x204),
    )
    c.Cycle()

    if c.ProgramCounter != 0x204 {
        t.Error("program counter is wrong")
    }
}

func Test_SE_equal(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x12),
        SE(0x1, 0x12),
    )
    c.Cycle()
    c.Cycle()

    if c.ProgramCounter != 0x206 {
        t.Error("program counter is wrong")
    }
}

func Test_SE_not_equal(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x12),
        SE(0x1, 0x13),
    )
    c.Cycle()
    c.Cycle()

    if c.ProgramCounter != 0x204 {
        t.Error("program counter is wrong")
    }
}

func Test_SNE_equal(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x12),
        SNE(0x1, 0x12),
    )
    c.Cycle()
    c.Cycle()

    if c.ProgramCounter != 0x204 {
        t.Error("program counter is wrong")
    }
}

func Test_SNE_not_equal(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x12),
        SNE(0x1, 0x13),
    )
    c.Cycle()
    c.Cycle()

    if c.ProgramCounter != 0x206 {
        t.Error("program counter is wrong")
    }
}

func Test_SE_VXY_equal(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x12),
        LD(0x2, 0x12),
        SE_R(0x1, 0x2),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.ProgramCounter != 0x208 {
        t.Error("program counter is wrong")
    }
}

func Test_SE_VXY_not_equal(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x12),
        LD(0x1, 0x13),
        SE_R(0x1,0x2),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.ProgramCounter != 0x206 {
        t.Error("program counter is wrong")
    }
}

func Test_LD(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x12),
    )
    c.Cycle()

    if c.Register[1] != 0x12 {
        t.Error("value wasn't set")
    }
}

func Test_ADD(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x12),
        ADD(0x1, 0x10),
    )
    c.Cycle()
    c.Cycle()

    if c.Register[1] != 0x22 {
        t.Error("value wasn't added")
    }
}

func Test_ADD_overflow(t *testing.T) {
    c := NewTestCPU(
        LD(0x1,0xFF),
        ADD(0x1, 0x10),
    )
    c.Cycle()
    c.Cycle()

    if c.Register[1] != 0x0F {
        t.Error("value wasn't added")
    }

    if c.Register[0xF] != 0 {
        t.Error("carry should be ignored")
    }
}

func Test_LD_R(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x20),
        LD_R(0x2, 0x1),
    )
    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != c.Register[0x2] || c.Register[0x1] != 0x20 {
        t.Error("registers not set properly")
    }
}

func Test_OR(t *testing.T) {
    c := NewTestCPU(
        LD(0x1,0x20),
        LD(0x2, 0x10),
        OR(0x1, 0x2),
        OR(0x1, 0x2),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0x30 {
        t.Error("unexpected value")
    }
}

func Test_AND(t *testing.T) {
    c := NewTestCPU(
        LD(0x1,0xFF),
        LD(0x2, 0x04),
        AND(0x1, 0x2),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0x04 {
        t.Error("unexpected value")
    }
}

func Test_XOR(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x10),
        LD(0x2, 0xFF),
        XOR(0x1, 0x2),
    )

    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0xEF {
        t.Error("unexpected value")
    }
}

func Test_ADD_R(t *testing.T) {
    c := NewTestCPU(
        LD(0x1,0x12),
        LD(0x2, 0x15),
        ADD_R(0x1,0x2),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0x27 {
        t.Error("unexpected value")
    }

    if c.Register[0xF] != 0x0 {
        t.Error("unexpected carry")
    }
}

func Test_ADD_R_overflow(t *testing.T) {
    c := NewTestCPU(
        LD(0x1,0xFF),
        LD(0x2, 0x15),
        ADD_R(0x1,0x2),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0x14 {
        t.Errorf("unexpected value: %v", c.Register[0x1])
    }

    if c.Register[0xF] != 0x1 {
        t.Errorf("unexpected carry: %v", c.Register[0xF])
    }
}

func Test_SUB(t *testing.T) {
    c := NewTestCPU(
        LD(0x1,0x12),
        LD(0x2, 0x10),
        SUB(0x1,0x2),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0x2 {
        t.Error("unexpected value")
    }

    if c.Register[0xF] != 0x1 {
        t.Error("unexpected borrow value")
    }
}

func Test_SUB_borrow(t *testing.T) {
    c := NewTestCPU(
        LD(0x1,0x10),
        LD(0x2, 0x12),
        SUB(0x1,0x2),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0xFE {
        t.Error("unexpected value")
    }

    if c.Register[0xF] != 0x0 {
        t.Error("unexpected borrow value")
    }
}

func Test_SHR(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x4),
        SHR(0x1),
    )

    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0x2 {
        t.Error("unexpected value")
    }

    if c.Register[0xF] != 0 {
        t.Error("unexpected value")
    }
}

func Test_SHR_carry(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x5),
        SHR(0x1),
    )

    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0x2 {
        t.Error("unexpected value")
    }

    if c.Register[0xF] != 1 {
        t.Error("unexpected value")
    }
}

func Test_SUBN(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x10),
        LD(0x2, 0x12),
        SUBN(0x1,0x2),
    )

    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0x2 {
        t.Error("unexpected value")
    }
    if c.Register[0xF] != 0x1 {
        t.Error("unexpected borrow")
    }
}

func Test_SUBN_borrow(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0xFF),
        LD(0x2, 0x12),
        SUBN(0x1,0x2),
    )

    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0x13 {
        t.Error("unexpected value")
    }
    if c.Register[0xF] != 0x0 {
        t.Error("unexpected borrow")
    }
}

func Test_SHL(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x7E),
        SHL(0x1),
    )

    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0xFC {
        t.Error("unexpected value")
    }
    if c.Register[0xF] != 0x0 {
        t.Error("unexpected carry")
    }
}

func Test_SHL_carry(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x80),
        SHL(0x1),
    )

    c.Cycle()
    c.Cycle()

    if c.Register[0x1] != 0x00 {
        t.Error("unexpected value")
    }
    if c.Register[0xF] != 0x1 {
        t.Error("unexpected carry")
    }
}

func Test_SNE_R_equal(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x12),
        LD(0x2, 0x12),
        SNE_R(0x1,0x2),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.ProgramCounter != 0x206 {
        t.Error("unexpected pc")
    }
}

func Test_SNE_R_not_equal(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x12),
        LD(0x2, 0x13),
        SNE_R(0x1,0x2),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.ProgramCounter != 0x208 {
        t.Error("unexpected pc")
    }
}

func Test_LDI(t *testing.T) {
    c := NewTestCPU(
        LDI(0x120),
    )

    c.Cycle()

    if c.Index != 0x120 {
        t.Error("unexpected i")
    }
}

func Test_JP_R(t *testing.T) {
    c := NewTestCPU(
        LD(0, 0x12),
        JP_R(0x200),
    )
    c.Cycle()
    c.Cycle()

    if c.ProgramCounter != 0x212 {
        t.Errorf("unexpected pc: %x", c.ProgramCounter)
    }
}

func Test_RND(t *testing.T) {
    c := NewTestCPU(RND(1, 0x13))
    c.Cycle()

    // no assert, just making sure it doesn't crash
}

func Test_DRW(t *testing.T) {
    c := NewTestCPU(
        LD(0x1, 0x2),
        LD(0x2, 0x3),
        LD(0x3, 0b00010000),

        DRW(0x1, 0x2, 1),
    )

    c.Cycle()
    c.Cycle()
    c.Cycle()
    c.Cycle()

}

func Test_SKP(t *testing.T) {
    t.Fail()
}

func Test_SKNP(t *testing.T) {
    t.Fail()
}

func Test_LD_R_DT(t *testing.T) {
    t.Fail()
}

func Test_LDK(t *testing.T) {
    t.Fail()
}

    func Test_LD_DT_R(t *testing.T) {
        t.Fail()

}

func Test_LD_ST_R(t *testing.T) {
    t.Fail()

}

func Test_ADDI(t *testing.T) {
    t.Fail()

}

func Test_LDF(t *testing.T) {
    t.Fail()

}

func Test_LDB(t *testing.T) {
    c := NewTestCPU(
        LDI(0x400),
        LD(0x1, 174),
        LDB(0x1),
    )
    c.Cycle()
    c.Cycle()
    c.Cycle()

    if c.Memory[0x400] != 1 {
        t.Errorf("unexpected value for hundreds: %v", c.Memory[0x400])
    }
    if c.Memory[0x401] != 7 {
        t.Errorf("unexpected value for tens: %v", c.Memory[0x401])
    }
    if c.Memory[0x402] != 4 {
        t.Errorf("unexpected value for ones: %v", c.Memory[0x402])
    }
}

func Test_LD_I_R(t *testing.T) {
    t.Fail()

}

func Test_LD_R_I(t *testing.T) {
    c := NewTestCPU(
        LDI(0x400),
        LD_VX_I(0x3),
    )
    c.Memory[0x400] = 4
    c.Memory[0x401] = 3
    c.Memory[0x402] = 10
    c.Memory[0x403] = 100
    c.Memory[0x404] = 123

    c.Cycle()
    c.Cycle()

    if c.Register[0x0] != 4 {
        t.Errorf("unexpected v[0] value: %v", c.Register[0x0])
    }
    if c.Register[0x1] != 3 {
        t.Errorf("unexpected v[1] value: %v", c.Register[0x1])
    }
    if c.Register[0x2] != 10 {
        t.Errorf("unexpected v[2] value%v", c.Register[0x2])
    }
    if c.Register[0x3] != 100 {
        t.Errorf("unexpected v[3] value: %v", c.Register[0x3])
    }
    if c.Register[0x4] != 0 {
        t.Errorf("unexpected v[4] value: %v", c.Register[0x4])
    }
}

func build(ops ...uint16) []byte {
    data := []byte {}
    for _, v := range ops {
        data = append(data, byte(v >> 8))
        data = append(data, byte(v))
    }
    return data
}

// instruction builders
func CLS() uint16 {
    return 0x00E0
}

func CALL(addr uint16) uint16 {
    return 0x2000 | addr
}

func NOP() uint16 {
    return 0x0000
}

func RET() uint16 {
    return 0x00EE
}

func JP(addr uint16) uint16 {
    return 0x1000 | addr
}

func SE(index uint16, value uint16) uint16 {
    return 0x3000 | (index << 8) | (value)
}

func SNE(index uint16, value uint16) uint16 {
    return 0x4000 | (index << 8) | (value)
}

func SE_R(x, y uint16) uint16 {
    return 0x5000 | (x << 8) | (y << 4)
}

func LD(index uint16, value uint16) uint16 {
    return 0x6000 | (index << 8) | (value)
}

func LD_R(x, y uint16) uint16 {
    return 0x8000 | (x << 8) | (y << 4)
}

func ADD(index uint16, value uint16) uint16 {
    return 0x7000 | (index << 8) | (value)
}

func OR(x, y uint16) uint16 {
    return 0x8001 | (x << 8) | (y << 4)
}

func AND(x,y uint16) uint16 {
    return 0x8002 | (x << 8) | (y << 4)
}

func XOR(x, y uint16) uint16 {
    return 0x8003 | (x << 8) | (y << 4)
}

func ADD_R(x,y uint16) uint16 {
    return 0x8004 | (x << 8) | (y << 4)
}

func SUB(x,y uint16) uint16 {
    return 0x8005 | (x << 8) | (y << 4)
}

func SHR(x uint16) uint16 {
    return 0x8006 | (x << 8)
}

func SUBN(x,y uint16) uint16 {
    return 0x8007 | (x << 8) | (y << 4)
}

func SHL(x uint16) uint16 {
    return 0x800E | (x << 8)
}

func SNE_R(x,y uint16) uint16 {
    return 0x9000 | (x << 8) | (y << 4)
}

func LDI(value uint16) uint16 {
    return 0xA000 | value
}

func JP_R(value uint16) uint16 {
    return 0xB000 | value
}

func RND(x, value uint16) uint16 {
    return 0xC000 | (x << 8) | value
}

func DRW(x,y,n uint16) uint16 {
    return 0xD000 | (x << 8) | (y << 4) | n
}

func SKP(x uint16) uint16 {
    return 0xE09E | (x << 8)
}

func SKNP(x uint16) uint16 {
    return 0xE0A1 | (x << 8)
}

func LD_VX_DT(x uint16) uint16 {
    return 0xF007 | (x << 8)
}

func LD_VX_K(x uint16) uint16 {
    return 0xF00A | (x << 8)
}

func LD_DT_VX(x uint16) uint16 {
    return 0xF015 | (x << 8)
}

func LD_ST_VX(x uint16) uint16 {
    return 0xF018 | (x << 8)
}

func ADD_I(x uint16) uint16 {
    return 0xF01E | (x << 8)
}

func LDF(x uint16) uint16 {
    return 0xF029 | (x << 8)
}

func LDB(x uint16) uint16 {
    return 0xF033 | (x << 8)
}

func LD_I_VX(x uint16) uint16 {
    return 0xF055 | (x << 8)
}

func LD_VX_I(x uint16) uint16 {
    return 0xF065 | (x << 8)
}


