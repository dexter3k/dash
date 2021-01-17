package com

type I8_8 int16

func (f I8_8) ToFloat() float32 {
	return float32(f) / 0x100
}


type I16_16 int32

func (f I16_16) ToFloat() float32 {
	return float32(f) / 0x10000
}
