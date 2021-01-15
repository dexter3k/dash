package player

type Player struct {
	RootSwf  *SwfLoader
	SwfFiles map[string]*SwfLoader
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
	return nil
}

func (p *Player) Destroy() {
	p.RootSwf = nil
	for _, v := range p.SwfFiles {
		v.Destroy()
	}
	p.SwfFiles = nil
}
