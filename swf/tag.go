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
