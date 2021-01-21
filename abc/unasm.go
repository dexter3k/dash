package abc

import (
	"fmt"

	"github.com/dexter3k/dash/abc/opcodes"

	"github.com/davecgh/go-spew/spew"
)

// Attempt disasm body code

func unasm(code []byte, eh []Exception) []Instruction {
	result := []Instruction{}

	var pc int
	u8 := func() uint8 {
		value := code[pc]
		pc++
		return value
	}
	u30 := func() int {
		result := uint32(u8())
		if (result & 0x80) != 0 {
			result = (result & 0x7f) | (uint32(u8()) << 7)
			if (result & 0x4000) != 0 {
				result = (result & 0x3fff) | (uint32(u8()) << 14)
				if (result & 0x200000) != 0 {
					result = (result & 0x1fffff) | (uint32(u8()) << 21)
					if (result & 0x10000000) != 0 {
						result = (result & 0xfffffff) | (uint32(u8()) << 28)
						result &= 0xffffffff
					}
				}
			}
		}
		return int(result)
	}
	s24 := func() int {
		result := uint32(u8()) | (uint32(u8()) << 8) | (uint32(u8()) << 16)
		if (result & 0x800000) != 0 {
			result |= 0xFF000000
		}
		return int(int32(result))
	}
	for pc < len(code) {
		startOffset := pc
		opcode := u8()
		args, found := opcodes.Arguments[opcode]
		if !found {
			for i := 0; i < len(result); i++ {
				kek := result[i]
				name, found := opcodes.Names[kek.Opcode]
				if !found {
					name = "???"
				}
				fmt.Printf("%04X: %02X %s\n", kek.Offset, kek.Raw, name)
			}
			name, found := opcodes.Names[opcode]
			if !found {
				name = "???"
			}
			fmt.Printf("Unknown opcode: [%02X] %s\n", opcode, name)
			// spew.Dump(result)
			panic("Unknown opcode")
		}

		opers := make([]int, len(args), len(args))
		for i := 0; i < len(args); i++ {
			dt := args[i] & opcodes.DataTypeMask
			switch dt {
			case opcodes.U8:
				opers[i] = int(u8())
			case opcodes.S24:
				opers[i] = s24()
			case opcodes.U30:
				opers[i] = u30()
			case opcodes.OffsetLookup:
				count := u30() + 1
				if count*3 > len(code)-pc {
					// fuck no, we can't do that.
					fmt.Println("Failed lookup switch count check", len(code), pc, count)
					panic("foobar")
					return nil
				}

				if i+1 != len(args) {
					panic("OffsetLookup should be last")
				}

				opers = append(opers, make([]int, count-1)...)
				for j := 0; j < count; j++ {
					opers[i+j] = s24()
				}
			default:
				fmt.Println(dt)
				panic("unknown arg")
			}
		}

		result = append(result, Instruction{
			Offset: startOffset,
			Raw:    code[startOffset:pc],
			Opcode: opcode,
			Args:   opers,
		})

		if false {
			spew.Dump(result)
		}
	}

	return result
}
