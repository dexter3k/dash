package abc

import (
	"github.com/dexter3k/dash/abc/opcodes"
)

func writeS24(buf []byte, value int) int {
	if value > 8388607 || value < -8388608 {
		panic("writeS24 out of range")
	}
	buf[0] = uint8(value)
	buf[1] = uint8(value >> 8)
	buf[2] = uint8(value >> 16)
	return 3
}

func formatU30(value int) []byte {
	result := []byte{}
	x := uint32(value)
	for {
		next := uint8(x & 0x7f)
		x >>= 7

		if x == 0 {
			result = append(result, next)
			break
		} else {
			result = append(result, next|0x80)
		}
	}
	return result
}

func asm(instr *Instruction) []byte {
	result := []byte{instr.Opcode}
	op := opcodes.Opcodes[instr.Opcode]
	for i := 0; i < len(op.Args); i++ {
		dt := op.Args[i] & opcodes.DataTypeMask
		switch dt {
		case opcodes.U8:
			result = append(result, uint8(instr.Args[i]))
		case opcodes.S24:
			result = append(result, make([]byte, 3)...)
			writeS24(result[len(result)-3:], instr.Args[i])
		case opcodes.U30:
			result = append(result, formatU30(instr.Args[i])...)
		default:
			panic(uint(dt))
		}
	}
	return result
}
