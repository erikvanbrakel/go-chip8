package main

import (
    "encoding/binary"
    "fmt"
)

type CPU struct {
    Memory [4096]byte
    Register [16]byte
    Index uint16
    DelayTimer byte
    SoundTimer byte
    ProgramCounter uint16
    StackPointer byte
    DisplayBuffer [64][32]byte
    Stack [16]uint16
}

var fontSet = []byte {
    0xF0,0x90,0x90,0x90,0xF0, // "0"
    0x20,0x60,0x20,0x20,0x70, // "1"
    0xF0,0x10,0xF0,0x80,0xF0, // "2"
    0xF0,0x10,0xF0,0x10,0xF0, // "3"
    0x90,0x90,0xF0,0x10,0x10, // "4"
    0xF0,0x80,0xF0,0x10,0xF0, // "5"
    0xF0,0x80,0xF0,0x90,0xF0, // "6"
    0xF0,0x10,0x20,0x40,0x40, // "7"
    0xF0,0x90,0xF0,0x90,0xF0, // "8"
    0xF0,0x90,0xF0,0x10,0xF0, // "9"
    0xF0,0x90,0xF0,0x90,0x90, // "A"
    0xE0,0x90,0xE0,0x90,0xE0, // "B"
    0xF0,0x80,0x80,0x80,0xF0, // "C"
    0xE0,0x90,0x90,0x90,0xE0, // "D"
    0xF0,0x80,0xF0,0x80,0xF0, // "E"
    0xF0,0x80,0xF0,0x80,0x80, // "F"
}

func NewCPU(programData []byte) *CPU {
    cpu := &CPU{}
    cpu.Initialize()
    cpu.LoadProgram(programData)
    return cpu
}

func (c *CPU) Initialize() {
    c.DelayTimer = 0
    c.SoundTimer = 0
    c.ProgramCounter = 0x200
    c.Index = 0
    c.StackPointer = 0
    c.DisplayBuffer = [64][32]byte{}
    c.Memory = [4096]byte{}
    c.Stack = [16]uint16{}
    c.Register = [16]byte {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
    for i := 0;i< len(fontSet);i++ {
        c.Memory[i] = fontSet[i]
    }
}

func (c *CPU) Cycle() {
    opCode := binary.BigEndian.Uint16(c.Memory[c.ProgramCounter:c.ProgramCounter+2])
    vX := opCode & 0x0F00 >> 8
    vY := opCode & 0x00F0 >> 4

    switch opCode & 0xF000 {
    case 0x0000: // SYS
        if opCode == 0x00E0 {
            c.DisplayBuffer = [64][32]byte{}
        }
        if opCode == 0x00EE {
            c.StackPointer--
            c.ProgramCounter = c.Stack[c.StackPointer]
        }
    case 0x1000: // JMP
        c.ProgramCounter = opCode & 0x0FFF - 2
    case 0x2000: // CALL
        c.Stack[c.StackPointer] = c.ProgramCounter
        c.StackPointer++
        c.ProgramCounter = opCode & 0x0FFF - 2
    case 0x3000: // SKPE
        if c.Register[vX] == (byte)(opCode & 0x00FF) {
            c.ProgramCounter += 2
        }
    case 0x4000: // SKPNE
        if c.Register[vX] != (byte)(opCode & 0x00FF) {
            c.ProgramCounter += 2
        }
    case 0x5000: // SKIPRE
        if c.Register[vX] == c.Register[vY] {
            c.ProgramCounter += 2
        }
    case 0x6000: // SET
        c.Register[vX] = (byte)(opCode & 0x00FF)
    case 0x7000:
        c.Register[vX] = byte((uint16(c.Register[vX]) + opCode & 0x00FF) & 0xFF)
    case 0x8000:
        op := opCode & 0x000F
        switch op {
            case 0x0000: c.Register[vX] = c.Register[vY]
            case 0x0001: c.Register[vX] |= c.Register[vY]
            case 0x0002: c.Register[vX] &= c.Register[vY]
            case 0x0003: c.Register[vX] ^= c.Register[vY]
            case 0x0004:
                if uint16(c.Register[vX]) + uint16(c.Register[vY]) > 255 { c.Register[0xF] = 1 } else { c.Register[0xF] = 0 }
                c.Register[vX] = (c.Register[vX] + c.Register[vY]) & 0xFF
            case 0x0005:
                if c.Register[vX] > c.Register[vY] { c.Register[0xF] = 1 } else { c.Register[0xF] = 0 }
                c.Register[vX] = (c.Register[vX] - c.Register[vY]) & 0xFF
            case 0x0006:
                c.Register[0xF] = c.Register[vX] & 0x1
                c.Register[vX] >>= 1
            case 0x0007:
                if c.Register[vY] > c.Register[vX] { c.Register[0xF] = 1 } else { c.Register[0xF] = 0 }
                c.Register[vX] = (c.Register[vY] - c.Register[vX]) & 0xFF
            case 0x000E :
                c.Register[0xF] = c.Register[vX] & 0x80 >> 7
                c.Register[vX] <<= 1
            }

    case 0x9000:
        if c.Register[vX] != c.Register[vY] {
            c.ProgramCounter += 2
        }
    case 0xA000: // set I
        c.Index = opCode & 0x0FFF
    case 0xB000:
        c.ProgramCounter = uint16(c.Register[0]) + opCode & 0x0FFF - 2
    case 0xC000:
        break
    case 0xD000:
        rows := opCode & 0x000F
        x := uint16(c.Register[vX])
        y := uint16(c.Register[vY])
        c.Register[0xF] = 0
        for row:=uint16(0);row<rows;row++ {
            for b := uint16(0); b < 8; b++ {
                v := c.Memory[row + c.Index] & (1 << (8-b)) >> (8-b)
                if c.DisplayBuffer[x + b][row + y] == 1 {
                    c.Register[0xF] = 1
                }
                c.DisplayBuffer[x + b][row + y] ^= v
            }
        }
    case 0xF000:
        op := opCode & 0x00FF
        switch op {
            case 0x0015:
                c.DelayTimer = c.Register[vX]
        case 0x0018:
            c.SoundTimer = c.Register[vX]
        case 0x001E:
            c.Index += uint16(c.Register[vX])
        case 0x0029:
            c.Index = vX * 5
        case 0x0055:
            for i:=uint16(0);i<=vX;i++ {
                c.Memory[c.Index + i] = c.Register[i]
            }
        case 0x0065:
            for i:=uint16(0);i<=vX;i++ {
                 c.Register[i] = c.Memory[c.Index + i]
            }
        case 0x0033:
            hundreds := c.Register[vX] / 100
            tens := (c.Register[vX] - hundreds * 100) / 10
            ones := c.Register[vX] - hundreds * 100 - tens * 10
            c.Memory[c.Index] = hundreds
            c.Memory[c.Index + 1] = tens
            c.Memory[c.Index + 2] = ones

        default:
            panic(fmt.Sprintf("%x: %x", op, opCode))
        }
    default:
        panic(fmt.Sprintf("unknown opcode: %X", opCode))
    }
    c.ProgramCounter += 2
    if c.SoundTimer > 0 { c.SoundTimer-- }
    if c.DelayTimer > 0 { c.DelayTimer-- }
}

func (c *CPU) LoadProgram(data []byte) {
    for i := 0; i<len(data);i++ {
        c.Memory[i + 0x200] = data[i]
    }
}
