package com

import (
	"io"
	bin "encoding/binary"
)

var (
	le = bin.LittleEndian
)

type BitStream struct {
	r io.Reader

	bitBuffer uint8
	bitCount  int
}

func NewBitStream(r io.Reader) *BitStream {
	return &BitStream{
		r: r,
	}
}

func (l *BitStream) Reset() {
	l.bitBuffer = 0
	l.bitCount = 0
}

/*

{01100 100}{01111100}{0 0111000}{01000110}{11010100}{00010010}{11000000}



*/

func (l *BitStream) ReadRect(r *Rect) error {
	l.Reset()

	nn, err := l.ReadUB(5)
	if err != nil {
		return err
	}
	n := int(nn)
	minX, err := l.ReadSB(n)
	if err != nil {
		return err
	}
	minY, err := l.ReadSB(n)
	if err != nil {
		return err
	}
	maxX, err := l.ReadSB(n)
	if err != nil {
		return err
	}
	maxY, err := l.ReadSB(n)
	if err != nil {
		return err
	}
	r.MinX, r.MinY, r.MaxX, r.MaxY = int32(minX), int32(minY), int32(maxX), int32(maxY)
	return nil
}

func (l *BitStream) ReadMatrix(m *Matrix) error {
	l.Reset()

	if hasScale, err := l.ReadUB(1); err != nil {
		return err
	} else if hasScale != 0 {
		nBits, err := l.ReadUB(5)
		if err != nil {
			return err
		}
		if err := l.ReadVector2F(&m.Scale, int(nBits)); err != nil {
			return err
		}
	}

	if hasRotate, err := l.ReadUB(1); err != nil {
		return err
	} else if hasRotate != 0 {
		nBits, err := l.ReadUB(5)
		if err != nil {
			return err
		}
		if err := l.ReadVector2F(&m.Skew, int(nBits)); err != nil {
			return err
		}
	}

	nBits, err := l.ReadUB(5)
	if err != nil {
		return err
	}
	if err := l.ReadVector2(&m.Translate, int(nBits)); err != nil {
		return err
	}

	return nil
}

func (l *BitStream) ReadVector2(v *Vector2, n int) error {
	x, err := l.ReadSB(n)
	if err != nil {
		return err
	}

	y, err := l.ReadSB(n)
	if err != nil {
		return err
	}

	v.X, v.Y = int32(x), int32(y)

	return nil
}

func (l *BitStream) ReadVector2F(v *Vector2F, n int) error {
	x, err := l.ReadSB(n)
	if err != nil {
		return err
	}

	y, err := l.ReadSB(n)
	if err != nil {
		return err
	}

	v.X, v.Y = I16_16(x), I16_16(y)

	return nil
}

func (l *BitStream) ReadUB(n int) (uint64, error) {
	var res uint64
	for n != 0 {
		if l.bitCount == 0 {
			if err := bin.Read(l.r, le, &l.bitBuffer); err != nil {
				return 0, err
			}
			l.bitCount = 8
		}

		m := n
		if m > l.bitCount {
			m = l.bitCount
		}
		res <<= m
		res |= uint64(l.bitBuffer >> (8 - m))

		n -= m
		l.bitCount -= m
		l.bitBuffer <<= m
	}
	return res, nil
}

func (l *BitStream) ReadSB(n int) (int64, error) {
	y, err := l.ReadUB(n)
	x := int64(y << (64 - n)) >> (64 - n)
	return x, err
}
