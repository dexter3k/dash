package player

import (
	"fmt"

	"github.com/dexter3k/dash/avm2"
	"github.com/dexter3k/dash/swf"
	"github.com/dexter3k/dash/gfx"
)

type Player struct {
	RootSwf   *SwfLoader
	Ended     bool
	SwfFiles  map[string]*SwfLoader
	Limits    swf.ScriptLimits
	Frame     int
	NextFrame int

	avm2 *avm2.State

	// Sounds map[uint16]*swf.DefineSound

	Gfx gfx.Graphics
}

func NewPlayer() *Player {
	p := &Player{
		SwfFiles: map[string]*SwfLoader{},
	}
	return p
}

func (p *Player) AddSwf(name string) *SwfLoader {
	if _, found := p.SwfFiles[name]; found {
		panic("alread has an swf file!")
	}
	l := NewSwfLoader(p)
	p.SwfFiles[name] = l
	if p.RootSwf == nil {
		p.RootSwf = l
	}
	return l
}

func (p *Player) Play() error {
	if p.RootSwf == nil {
		return nil
	}

	p.NextFrame = 1

	for !p.Ended {
		if p.NextFrame > p.Frame {
			err := p.NextTag()
			if err != nil {
				return err
			}
		} else if p.NextFrame < p.Frame {
			panic("not implemented")
		}
	}

	return nil
}

func (p *Player) initAvm2() {
	if p.avm2 != nil {
		return
	}
	p.avm2 = avm2.NewState()
}

func (p *Player) NextTag() error {
	for p.RootSwf.Head >= len(p.RootSwf.Tags) {
		_, err := p.RootSwf.ReadNextTag()
		if err != nil {
			return err
		}
	}

	anyTag := p.RootSwf.Tags[p.RootSwf.Head]
	switch tag := anyTag.(type) {
	case *swf.FileAttributes:
		if p.RootSwf.Head != 0 {
			fmt.Println("Skipping file attributes at unexpected position")
			break
		}
		p.RootSwf.Attributes = tag
		if p.RootSwf.Attributes.ActionScript3 {
			p.initAvm2()
		}
	case *swf.DefineSound:
		// define the sound...
	case *swf.End:
		p.Ended = true
	case *swf.Metadata:
		// nop
	case *swf.ProductInfo:
		// nop
	case *swf.ScriptLimits:
		p.Limits = *tag
	case *swf.SetBackgroundColor:
		// ...
	case *swf.FrameLabel:
		// ...
	case *swf.DefineShape:
		// ...
	case *swf.ShowFrame:
		p.Frame++
	case *swf.DefineFont3:
		// ...
	case *swf.DefineFontAlignZones:
		// ...
	case *swf.DefineFontName:
		// ...
	case *swf.DefineText:
		// ...
	case *swf.ExportAssets:
		fmt.Println("Exported assets:", tag.Assets)
	case *swf.DoABC:
		p.avm2.AddAbc(tag.Name, tag.Data, tag.LazyInit)
	case *swf.SymbolClass:
		fmt.Println(tag)
	case *swf.Unknown:
		fmt.Printf("Skipping unknown tag %3d (%d bytes)\n", tag.Type, len(tag.Data))
	default:
		fmt.Printf("Unhandled known tag: %T\n", tag)
	}

	p.RootSwf.Head++

	return nil
}

func (p *Player) Destroy() {
	p.RootSwf = nil
	for _, v := range p.SwfFiles {
		v.Destroy()
	}
	p.SwfFiles = nil
}
