package com

type I8_8 int16

func (f I8_8) ToFloat() float32 {
	return float32(f) / 0x100
}
