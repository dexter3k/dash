package swf

import (
	"github.com/dexter3k/dash/com"
)

type Tag interface {
	// ...
}

type Unknown struct {
	Type uint8
	Data []byte
}

type End struct {
	// nothing
}

type FileAttributes struct {
	UseDirectBlit bool
	UseGPU        bool
	HasMetadata   bool
	ActionScript3 bool
	NoCache       bool
	UseNetwork    bool
}

type Metadata struct {
	// ignored...
}

type ProductInfo struct {
	// ignored...
}

type ShowFrame struct {
	// empty
}

type DefineFont3 struct {
	FontId uint16
	// other fields not implemented
}

type DefineFontAlignZones struct {
	FontId uint16
	// other fields not implemented
}

type DefineFontName struct {
	FontId uint16
	// other fields not implemented
}

type DefineSprite struct {
	SpriteId uint16
	Frames   uint16
	Controls []Tag
}

type DefineBitsJPEG2 struct {
	CharacterId uint16
	Data        []byte
}

type AssetExport struct {
	Id   uint16
	Name string
}

type ExportAssets struct {
	Assets []AssetExport
}

type DoABC struct {
	LazyInit bool
	Name     string
	Data     []byte
}

type SymbolClass struct {
	Assets []AssetExport
}

type GlyphEntry struct {
	Index   uint32
	Advance int32
}

type TextRecord struct {
	FontId  uint16
	Color   com.RGBA
	Offset  com.Vector2
	Height  uint16
	Entries []*GlyphEntry
}

type DefineText struct {
	TextId  uint16
	Bounds  com.Rect
	Matrix  com.Matrix
	Records []*TextRecord // not implemented
}


const (
	SoundUncompressedBigEndian    = 0
	SoundCompressedADPCM          = 1
	SoundCompressedMP3            = 2
	SoundUncompressedLittleEndian = 3
	SoundCompressedNellymoser1    = 4
	SoundCompressedNellymoser2    = 5
	SoundCompressedNellymoser3    = 6
	SoundCompressedSpeex          = 11
)

type DefineSound struct {
	Id      uint16
	Format  uint8
	RateHz  uint16 // 5500, 11000, 22000, 44000
	Is16Bit bool
	Stereo  bool
	Samples uint32
	Data    []byte
}

type ScriptLimits struct {
	Recursion uint16
	Timeout   uint16
}

type SetBackgroundColor struct {
	com.RGB
}

type FrameLabel struct {
	// ...
}
