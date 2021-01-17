package player

import (
	"bytes"
	"fmt"
	bin "encoding/binary"
	"io"

	"github.com/dexter3k/dash/swf"
	"github.com/dexter3k/dash/com"
)

const (
	TagEndId                = 0
	TagShowFrameId          = 1
	TagDefineShape1Id       = 2
	TagSetBackgroundColorId = 9
	TagDefineSoundId        = 14
	TagDefineShape2Id       = 22
	TagDefineShape3Id       = 32
	TagProductInfoId        = 41
	TagFrameLabelId         = 43
	TagScriptLimitsId       = 65
	TagFileAttributesId     = 69
	TagMetadataId           = 77
	TagDefineShape4Id       = 83
)

type tagParser func(tagType uint8, tagData []byte) (swf.Tag, error)

var (
	tagParsers = [256]tagParser{
		TagEndId:                endTagParser,
		TagShowFrameId:          showFrameParser,
		TagDefineShape1Id:       defineShape123TagParser,
		TagSetBackgroundColorId: setBackgroundColorParser,
		TagDefineSoundId:        defineSoundTagParser,
		TagDefineShape2Id:       defineShape123TagParser,
		TagDefineShape3Id:       defineShape123TagParser,
		TagProductInfoId:        productInfoTagParser,
		TagFrameLabelId:         frameLabelTagParser,
		TagScriptLimitsId:       scriptLimitsParser,
		TagFileAttributesId:     fileAttributesTagParser,
		TagMetadataId:           metadataTagParser,
	}
)

func parseFillStyles(r io.Reader, tagType uint8) ([]swf.FillStyle, error) {
	var fillStyleCount int
	{
		var fillStyleCount8 uint8
		if err := bin.Read(r, le, &fillStyleCount8); err != nil {
			return nil, err
		}
		fillStyleCount = int(fillStyleCount8)

		if tagType != TagDefineShape1Id && fillStyleCount == 0xff {
			var fillStyleCount16 uint16
			if err := bin.Read(r, le, &fillStyleCount16); err != nil {
				return nil, err
			}
			fillStyleCount = int(fillStyleCount16)
		}
	}

	fillStyles := make([]swf.FillStyle, fillStyleCount)
	for i := 0; i < fillStyleCount; i++ {
		var fillStyleType uint8
		if err := bin.Read(r, le, &fillStyleType); err != nil {
			return nil, err
		}

		switch fillStyleType {
		case 0x00: // solid fill style
			if tagType == TagDefineShape1Id || tagType == TagDefineShape2Id {
				var rgb com.RGB
				if err := bin.Read(r, le, &rgb); err != nil {
					return nil, err
				}
				fillStyles[i] = &swf.SolidFillStyle{
					Color: rgb.RGBA(),
				}
			} else {
				var rgba com.RGBA
				if err := bin.Read(r, le, &rgba); err != nil {
					return nil, err
				}
				fillStyles[i] = &swf.SolidFillStyle{
					Color: rgba,
				}
			}
		case 0x10, 0x12, 0x13: // linear/radial/focal gradient fill
			g := &swf.GradientFillStyle{
				Radial: fillStyleType != 0x10,
				Focal:  fillStyleType == 0x13,
			}
			fillStyles[i] = g

			bits := com.NewBitStream(r)
			if err := bits.ReadMatrix(&g.Matrix); err != nil {
				return fillStyles, err
			}

			bits.Reset()
			spread, err := bits.ReadUB(2)
			if err != nil { return fillStyles, err }
			g.Spread = swf.GradientSpread(spread)

			interp, err := bits.ReadUB(2)
			if err != nil { return fillStyles, err }
			g.LinearRGB = interp != 0

			points, err := bits.ReadUB(4)
			if err != nil { return fillStyles, err }

			g.Points = make([]*swf.GradientPoint, points)
			for i := 0; i < int(points); i++ {
				var ratio uint8
				if err := bin.Read(r, le, &ratio); err != nil {
					return fillStyles, err
				}
				if tagType == TagDefineShape1Id || tagType == TagDefineShape2Id {
					var rgb com.RGB
					if err := bin.Read(r, le, &rgb); err != nil {
						return nil, err
					}
					g.Points[i] = &swf.GradientPoint{
						Ratio: ratio,
						Color: rgb.RGBA(),
					}
				} else {
					var rgba com.RGBA
					if err := bin.Read(r, le, &rgba); err != nil {
						return nil, err
					}
					g.Points[i] = &swf.GradientPoint{
						Ratio: ratio,
						Color: rgba,
					}
				}
			}
		default:
			panic(fillStyleType)
		}
	}

	return fillStyles, nil
}

func parseLineStyles(r io.Reader, tagType uint8) ([]swf.LineStyle, error) {
	var lineStyleCount int
	{
		var lineStyleCount8 uint8
		if err := bin.Read(r, le, &lineStyleCount8); err != nil {
			return nil, err
		}
		lineStyleCount = int(lineStyleCount8)

		if lineStyleCount == 0xff {
			var lineStyleCount16 uint16
			if err := bin.Read(r, le, &lineStyleCount16); err != nil {
				return nil, err
			}
			lineStyleCount = int(lineStyleCount16)
		}
	}

	lineStyles := make([]swf.LineStyle, lineStyleCount)
	for i := 0; i < lineStyleCount; i++ {
		if tagType == TagDefineShape4Id {
			panic("panik!")
		}

		var width uint16
		if err := bin.Read(r, le, &width); err != nil {
			return nil, err
		}

		if tagType == TagDefineShape1Id || tagType == TagDefineShape2Id {
			var rgb com.RGB
			if err := bin.Read(r, le, &rgb); err != nil {
				return nil, err
			}
			lineStyles[i] = &swf.LineStyleSimple{
				Width: width,
				Color: rgb.RGBA(),
			}
		} else {
			var rgba com.RGBA
			if err := bin.Read(r, le, &rgba); err != nil {
				return nil, err
			}
			lineStyles[i] = &swf.LineStyleSimple{
				Width: width,
				Color: rgba,
			}
		}
	}

	return lineStyles, nil
}

func defineShape123TagParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	r := bytes.NewBuffer(tagData)
	bits := com.NewBitStream(r)
	t := &swf.DefineShape{}

	if err := bin.Read(r, le, &t.Id); err != nil {
		return nil, err
	}
	if err := bits.ReadRect(&t.Bounds); err != nil {
		return nil, err
	}

	fillStyles, err := parseFillStyles(r, tagType)
	if err != nil {
		return t, err
	}

	lineStyles, err := parseLineStyles(r, tagType)
	if err != nil {
		return t, err
	}

	bits.Reset()
	numFillBits, err := bits.ReadUB(4)
	if err != nil {
		return t, err
	}
	numLineBits, err := bits.ReadUB(4)
	if err != nil {
		return t, err
	}

	var currentShape *swf.Shape
	var position com.Vector2

	t.Shapes = []*swf.Shape{}

	// Extract shape records
	for {
		if typeFlag, err := bits.ReadUB(1); err != nil {
			return t, err
		} else if typeFlag == 0 { // non-edge
			newStyles, err := bits.ReadUB(1)
			if err != nil { return t, err }
			lineStyle, err := bits.ReadUB(1)
			if err != nil { return t, err }
			fillStyle1, err := bits.ReadUB(1)
			if err != nil { return t, err }
			fillStyle0, err := bits.ReadUB(1)
			if err != nil { return t, err }
			moveTo, err := bits.ReadUB(1)
			if err != nil { return t, err }

			if currentShape != nil {
				if len(currentShape.Edges) > 0 {
					t.Shapes = append(t.Shapes, currentShape)
				}
				currentShape = nil
			}

			if newStyles == 0 && lineStyle == 0 && fillStyle1 == 0 && fillStyle0 == 0 && moveTo == 0 {
				break
			}

			currentShape = &swf.Shape{
				Edges:    []*swf.Edge{},
			}

			if moveTo != 0 {
				moveBits, err := bits.ReadUB(5)
				if err != nil { return t, err }
				deltaX, err := bits.ReadSB(int(moveBits))
				if err != nil { return t, err }
				deltaY, err := bits.ReadSB(int(moveBits))
				if err != nil { return t, err }

				position.X += int32(deltaX)
				position.Y += int32(deltaY)
			}
			currentShape.Position = position

			if fillStyle0 != 0 {
				index, err := bits.ReadUB(int(numFillBits))
				if err != nil { return t, err }
				if index != 0 {
					currentShape.FillStyle0 = fillStyles[index-1]
				}
			}

			if fillStyle1 != 0 {
				index, err := bits.ReadUB(int(numFillBits))
				if err != nil { return t, err }
				if index != 0 {
					currentShape.FillStyle1 = fillStyles[index-1]
				}
			}

			if lineStyle != 0 {
				index, err := bits.ReadUB(int(numLineBits))
				if err != nil { return t, err }
				if index != 0 {
					currentShape.LineStyle = lineStyles[index-1]
				}
			}

			if newStyles != 0 {
				fillStyles, err = parseFillStyles(r, tagType)
				if err != nil {
					return t, err
				}

				lineStyles, err = parseLineStyles(r, tagType)
				if err != nil {
					return t, err
				}

				bits.Reset()
				numFillBits, err = bits.ReadUB(4)
				if err != nil {
					return t, err
				}
				numLineBits, err = bits.ReadUB(4)
				if err != nil {
					return t, err
				}
			}
		} else { // edge
			if currentShape == nil {
				panic("No current shape!")
			}

			straight, err := bits.ReadUB(1)
			if err != nil { return t, err }

			tmp, err := bits.ReadUB(4)
			if err != nil { return t, err }
			nBits := int(tmp) + 2

			if straight == 1 {
				hasDeltaX := uint32(1)
				hasDeltaY := uint32(1)

				if general, err := bits.ReadUB(1); err != nil {
					return t, err
				} else if general == 0 {
					if vertical, err := bits.ReadUB(1); err != nil {
						return t, err
					} else if vertical == 1 {
						hasDeltaX = 0
					} else {
						hasDeltaY = 0
					}
				}

				var delta com.Vector2
				if hasDeltaX != 0 {
					deltaX, err := bits.ReadSB(nBits)
					if err != nil { return t, err }
					delta.X = int32(deltaX)
				}
				if hasDeltaY != 0 {
					deltaY, err := bits.ReadSB(nBits)
					if err != nil { return t, err }
					delta.Y = int32(deltaY)
				}

				currentShape.Edges = append(currentShape.Edges, &swf.Edge{
					Target: delta,
				})
			} else {
				// quadratic curve
				var control, anchor com.Vector2
				tmp, err := bits.ReadSB(nBits)
				if err != nil { return t, err }
				control.X = int32(tmp)
				tmp, err = bits.ReadSB(nBits)
				if err != nil { return t, err }
				control.Y = int32(tmp)
				tmp, err = bits.ReadSB(nBits)
				if err != nil { return t, err }
				anchor.X = int32(tmp)
				tmp, err = bits.ReadSB(nBits)
				if err != nil { return t, err }
				anchor.Y = int32(tmp)

				currentShape.Edges = append(currentShape.Edges, &swf.Edge{
					Target:  anchor,
					Control: control,
				})
			}
		}
	}

	return t, nil
}

func unknownTagParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	t := &swf.Unknown{
		Type: tagType,
		Data: tagData,
	}
	return t, nil
}

func setBackgroundColorParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	s := &swf.SetBackgroundColor{}
	r := bytes.NewBuffer(tagData)
	if err := bin.Read(r, le, s); err != nil {
		return nil, err
	}
	return s, nil
}

func scriptLimitsParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	s := &swf.ScriptLimits{}
	r := bytes.NewBuffer(tagData)
	if err := bin.Read(r, le, s); err != nil {
		return nil, err
	}
	return s, nil
}

func showFrameParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	return &swf.ShowFrame{}, nil
}

func frameLabelTagParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	return &swf.FrameLabel{}, nil
}

func metadataTagParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	return &swf.Metadata{}, nil
}

func productInfoTagParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	return &swf.ProductInfo{}, nil
}

func defineSoundTagParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	r := bytes.NewBuffer(tagData)
	s := &swf.DefineSound{}
	if err := bin.Read(r, le, &s.Id); err != nil {
		return nil, err
	}
	var flags uint8
	if err := bin.Read(r, le, &flags); err != nil {
		return nil, err
	}
	s.Format = flags >> 4
	s.RateHz = []uint16{5500, 11000, 22000, 44000}[((flags >> 2) & 0x3)]
	if s.Format == swf.SoundCompressedMP3 && s.RateHz == 5500 {
		return s, fmt.Errorf("5.5kHz is not allowed for MP3")
	}

	s.Is16Bit = (flags & 0x2) != 0
	if !s.Is16Bit && s.Format == swf.SoundUncompressedLittleEndian {
		s.Format = swf.SoundUncompressedBigEndian
	}
	if s.Format != swf.SoundUncompressedBigEndian {
		s.Is16Bit = true
	}

	s.Stereo = (flags & 0x1) != 0
	if err := bin.Read(r, le, &s.Samples); err != nil {
		return nil, err
	}
	s.Data = r.Bytes()
	return s, nil
}

func endTagParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	return &swf.End{}, nil
}

func fileAttributesTagParser(tagType uint8, tagData []byte) (swf.Tag, error) {
	var flags uint32
	r := bytes.NewBuffer(tagData)
	if err := bin.Read(r, le, &flags); err != nil {
		return nil, err
	}
	t := &swf.FileAttributes{
		UseDirectBlit: (flags & 0x00000040) != 0,
		UseGPU:        (flags & 0x00000020) != 0,
		HasMetadata:   (flags & 0x00000010) != 0,
		ActionScript3: (flags & 0x00000008) != 0,
		NoCache:       (flags & 0x00000004) != 0,
		UseNetwork:    (flags & 0x00000001) != 0,
	}
	return t, nil
}
