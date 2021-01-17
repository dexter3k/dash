package player

import (
	"fmt"
	"bytes"
	"io"
	"compress/zlib"
	bin "encoding/binary"

	"github.com/dexter3k/dash/com"
	"github.com/dexter3k/dash/swf"
)

var (
	le = bin.LittleEndian
)

type swfHeader struct {
	Compression uint8
	Magic       uint16
	Version     uint8
	Size        uint32
}

type SwfLoader struct {
	Player *Player
	Data   []byte
	stream io.Reader

	header      swfHeader
	DisplayRect com.Rect
	FrameRate   com.I8_8
	FrameCount  uint16

	Tags []swf.Tag
	Head int

	Attributes *swf.FileAttributes

	bits *com.BitStream
}

func NewSwfLoader(p *Player) *SwfLoader {
	l := &SwfLoader{
		Player: p,
	}
	return l
}

func (l *SwfLoader) LoadRawBytes(data []byte) error {
	l.Data = make([]byte, len(data))
	copy(l.Data, data)
	l.stream = bytes.NewBuffer(l.Data)
	return l.parseHeader()
}

func (l *SwfLoader) parseHeader() error {
	if err := bin.Read(l.stream, le, &l.header); err != nil {
		return err
	}
	if l.header.Magic != 0x5357 {
		return fmt.Errorf("Corrupted header")
	}
	if l.header.Compression == 0x46 {
		// cool and good, no compression
	} else if l.header.Version >= 6 && l.header.Compression == 0x43 {
		// squeeeezed with zlib
		rc, err := zlib.NewReader(l.stream)
		if err != nil {
			return err
		}
		l.stream = rc
	} else if l.header.Version >= 13 && l.header.Compression == 0x5a {
		// lzma :D
		panic("ni")
	} else {
		return fmt.Errorf("Unknown compression format: %#v\n", string(rune(l.header.Compression)))
	}

	l.bits = com.NewBitStream(l.stream)

	if err := l.bits.ReadRect(&l.DisplayRect); err != nil {
		return err
	}
	if err := bin.Read(l.stream, le, &l.FrameRate); err != nil {
		return err
	}
	if err := bin.Read(l.stream, le, &l.FrameCount); err != nil {
		return err
	}

	return nil
}

func (l *SwfLoader) ReadNextTag() (swf.Tag, error) {
	tag, err := readNextTag(l.stream)
	if err == nil {
		l.Tags = append(l.Tags, tag)
	}
	return tag, err
}

func readNextTag(stream io.Reader) (swf.Tag, error) {
	if stream == nil {
		panic("not ready for reading")
	}

	var codeAndLength uint16
	if err := bin.Read(stream, le, &codeAndLength); err != nil {
		return nil, err
	}
	tagType := uint8(codeAndLength >> 6)
	tagSize := uint32(codeAndLength & 0x3f)
	if tagSize == 0x3f {
		if err := bin.Read(stream, le, &tagSize); err != nil {
			return nil, err
		}
	}

	tagData := make([]byte, tagSize)
	if err := bin.Read(stream, le, &tagData); err != nil {
		return nil, err
	}
	parser := tagParsers[tagType]
	if parser == nil {
		parser = unknownTagParser
	}
	tag, err := parser(tagType, tagData)
	return tag, err
}

func (l *SwfLoader) Destroy() {
	if rc, isrc := l.stream.(io.ReadCloser); isrc {
		rc.Close()
	}
	l.stream = nil
	l.Data = nil
}
