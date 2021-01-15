package player

import (
	"fmt"
	"bytes"
	"io"
	"compress/zlib"
	bin "encoding/binary"

	"github.com/dexter3k/dash/com"
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

	bitBuffer uint8
	bitCount  int
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

	if err := l.readRect(&l.DisplayRect); err != nil {
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

func (l *SwfLoader) Destroy() {
	if rc, isrc := l.stream.(io.ReadCloser); isrc {
		rc.Close()
	}
	l.stream = nil
	l.Data = nil
}

func (l *SwfLoader) readRect(r *com.Rect) error {
	l.resetBitStream()

	n, err := l.readUB(5)
	if err != nil {
		return err
	}
	minX, err := l.readSB(int(n))
	if err != nil {
		return err
	}
	minY, err := l.readSB(int(n))
	if err != nil {
		return err
	}
	maxX, err := l.readSB(int(n))
	if err != nil {
		return err
	}
	maxY, err := l.readSB(int(n))
	if err != nil {
		return err
	}
	r.MinX, r.MinY, r.MaxX, r.MaxY = int32(minX), int32(minY), int32(maxX), int32(maxY)
	return nil
}

func (l *SwfLoader) resetBitStream() {
	l.bitBuffer = 0
	l.bitCount = 0
}

func (l *SwfLoader) readUB(n int) (uint, error) {
	var res uint
	for n != 0 {
		if l.bitCount == 0 {
			if err := bin.Read(l.stream, le, &l.bitBuffer); err != nil {
				return 0, err
			}
			l.bitCount = 8
		}

		m := n
		if m > l.bitCount {
			m = l.bitCount
		}
		res <<= m
		res |= uint(l.bitBuffer >> (8 - m))

		n -= m
		l.bitCount -= m
		l.bitBuffer <<= m
	}
	return res, nil
}

func (l *SwfLoader) readSB(n int) (int, error) {
	y, err := l.readUB(n)
	x := int(y)
	if (x >> (n-1)) != 0 && x > 0 {
		x = -x
	}
	return x, err
}
