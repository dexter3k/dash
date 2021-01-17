package com

type RGB struct {
	R, G, B uint8
}

func (rgb RGB) RGBA() RGBA {
	return RGBA{
		RGB: rgb,
		A:   0xff,
	}
}

type RGBA struct {
	RGB
	A uint8
}
