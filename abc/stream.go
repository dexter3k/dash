package abc

import (
	"encoding/binary"
)

type abcStream struct {
	buffer []byte
	offset int
}

func NewAbcStream(buffer []byte) *abcStream {
	return &abcStream{
		buffer: buffer,
		offset: 0,
	}
}

func (s *abcStream) GetOffset() int {
	return s.offset
}

func (s *abcStream) SetOffset(offset int) {
	s.offset = offset
}

func (s *abcStream) ReadU8() uint8 {
	s.offset++
	return s.buffer[s.offset-1]
}

func (s *abcStream) ReadU32() uint32 {
	value := binary.LittleEndian.Uint32(s.buffer[s.offset : s.offset+4])
	s.offset += 4
	return value
}

func (s *abcStream) ReadU64() uint64 {
	value := binary.LittleEndian.Uint64(s.buffer[s.offset : s.offset+8])
	s.offset += 8
	return value
}

func (s *abcStream) ReadVarlenU32() uint32 {
	result := uint32(s.ReadU8())
	if (result & 0x80) != 0 {
		result = (result & 0x7f) | (uint32(s.ReadU8()) << 7)
		if (result & 0x4000) != 0 {
			result = (result & 0x3fff) | (uint32(s.ReadU8()) << 14)
			if (result & 0x200000) != 0 {
				result = (result & 0x1fffff) | (uint32(s.ReadU8()) << 21)
				if (result & 0x10000000) != 0 {
					result = (result & 0xfffffff) | (uint32(s.ReadU8()) << 28)
					result &= 0xffffffff
				}
			}
		}
	}
	return result
}

func (s *abcStream) ReadVarlenS32() int32 {
	return int32(s.ReadVarlenU32())
}

func (s *abcStream) ReadVarlenU30() uint32 {
	value := s.ReadVarlenU32()
	if (value & 0xc0000000) != 0 {
		panic("nope")
	}
	return value
}

func (s *abcStream) ReadPrefixedString() string {
	length := int(s.ReadVarlenU30())
	result := string(s.buffer[s.offset : s.offset+length])
	s.offset += length
	return result
}

func (s *abcStream) ReadByteSlice(length int) []byte {
	out := make([]byte, length, length)
	copy(out, s.buffer[s.offset:s.offset+length])
	s.offset += length
	return out
}
