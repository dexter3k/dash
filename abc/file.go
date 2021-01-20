package abc

import (
	"bytes"
	bin "encoding/binary"
	"fmt"
	"math"
	"io"

	// "github.com/davecgh/go-spew/spew"
)

// TODO: error handling, currently we panic in most occasions
// there should be 0 panics, unchecked array accesses, etc during parsing

var le = bin.LittleEndian

func Parse(data []byte) (*File, error) {
	p := newParser(data)
	return p.parseFile(), p.err()
}

type parser struct {
	f     *File
	r     *bytes.Buffer
	error error
}

func newParser(data []byte) *parser {
	return &parser{
		r: bytes.NewBuffer(data),
	}
}

func (p *parser) err() error {
	return p.error
}

func (p *parser) parseFile() *File {
	if p.err() != nil {
		return nil
	}

	p.f = &File{}

	p.error = bin.Read(p.r, le, &p.f.Minor)
	p.error = bin.Read(p.r, le, &p.f.Major)

	p.parseConstantPool()
	if p.err() != nil {
		return nil
	}

	p.readMethodPool()
	if p.err() != nil {
		return nil
	}

	p.skipMetadataPool()
	if p.err() != nil {
		return nil
	}

	p.readClassPool()
	if p.err() != nil {
		return nil
	}

	p.readScriptPool()
	if p.err() != nil {
		return nil
	}

	p.readBodyPool()
	if p.err() != nil {
		return nil
	}

	return p.f
}

func (p *parser) parseConstantPool() {
	{
		intCount := p.u30()
		if intCount < 1 {
			intCount = 1
		}
		p.f.Int = make([]int32, intCount)
		for i := 1; i < len(p.f.Int); i++ {
			p.f.Int[i] = p.s32()
		}

		if p.err() != nil {
			return
		}
	}

	{
		uintCount := p.u30()
		if uintCount < 1 {
			uintCount = 1
		}
		p.f.Uint = make([]uint32, uintCount)
		for i := 1; i < len(p.f.Uint); i++ {
			p.f.Uint[i] = p.u32()
		}

		if p.err() != nil {
			return
		}
	}

	{
		numberCount := p.u30()
		if numberCount < 1 {
			numberCount = 1
		}
		p.f.Number = make([]float64, numberCount)
		p.error = bin.Read(p.r, le, p.f.Number[1:])
		p.f.Number[0] = math.NaN()

		if p.err() != nil {
			return
		}
	}

	{
		stringCount := p.u30()
		if stringCount < 1 {
			stringCount = 1
		}
		p.f.String = make([]string, stringCount)
		for i := 1; i < len(p.f.String); i++ {
			stringSize := p.u30()
			stringData := p.r.Next(int(stringSize))
			if len(stringData) < int(stringSize) {
				p.error = io.EOF
				return
			}
			p.f.String[i] = string(stringData)
		}

		if p.err() != nil {
			return
		}
	}

	p.readNamespacePool()
	if p.err() != nil {
		return
	}

	p.readSetPool()
	if p.err() != nil {
		return
	}

	p.readNamePool()
	if p.err() != nil {
		return
	}
}

func (p *parser) readNamespacePool() {
	spaceCount := int(p.u30())
	if spaceCount < 1 {
		spaceCount = 1
	}
	p.f.Space = make([]Space, spaceCount)
	for i := 1; i < len(p.f.Space); i++ {
		var kind uint8
		kind, p.error = p.r.ReadByte()
		if p.err() != nil {
			return
		}

		index := int(p.u30()) // a zero index is an empty string, not "*"
		switch kind {
		case 0x08, 0x16: // no idea why there are 2 distinct namespace types
			p.f.Space[i] = Space(PublicNamespacePrefix + p.f.String[index])
		case 0x17:
			p.f.Space[i] = Space(PackageInternalNamespacePrefix + p.f.String[index])
		case 0x18:
			p.f.Space[i] = Space(ProtectedNamespacePrefix + p.f.String[index])
		case 0x19:
			p.f.Space[i] = Space(ExplicitNamespacePrefix + p.f.String[index])
		case 0x1a:
			p.f.Space[i] = Space(StaticProtectedNamespacePrefix + p.f.String[index])
		case 0x05:
			p.f.Space[i] = Space(PrivateNamespacePrefix + p.f.String[index])
		default:
			panic("Unknown namespace type")
		}

		// Note: we don't handle versioned spaces at all here
		// These should only ever be inside avmplus builtins, never in user code
		// So if we'll handle all user code that happens to have versioned urls
		// the same way as normal code, we'll be okay!
	}
	p.f.Space[0] = Space(PublicNamespacePrefix + "*") // lets assume its public any
}

func (p *parser) readSetPool() {
	setCount := int(p.u30())
	if setCount < 1 {
		setCount = 1
	}
	p.f.Set = make([][]Space, setCount)
	for i := 1; i < setCount; i++ {
		spaceCount := int(p.u30())
		p.f.Set[i] = make([]Space, spaceCount)
		for j := 0; j < spaceCount; j++ {
			index := p.u30()
			if index == 0 {
				panic("Entry in a namespace set cannot be zero")
			}
			p.f.Set[i][j] = p.f.Space[index]
		}
	}
}

func (p *parser) readNamePool() {
	nameCount := int(p.u30())
	if nameCount < 1 {
		nameCount = 1
	}
	p.f.Name = make([]Name, nameCount)
	for i := 1; i < nameCount; i++ {
		var kind uint8
		kind, p.error = p.r.ReadByte()
		if p.err() != nil {
			return
		}

		switch kind {
		case 0x07, 0x0d: // both name and space known, exact name
			ns := int(p.u30())   // zero is any (default)
			index := int(p.u30()) // zero is any (not default)
			name := "*"
			if index != 0 {
				name = p.f.String[index]
			}

			p.f.Name[i] = &FullName{
				BaseName: BaseName{kind == 0x0d},
				Name:     name,
				Space:    p.f.Space[ns],
			}
		case 0x0f, 0x10: // space specified by bytecode (only bytecode can refer to this?)
			index := int(p.u30()) // zero is any (not default)
			name := "*"
			if index != 0 {
				name = p.f.String[index]
			}

			p.f.Name[i] = &RuntimeSpaceName{
				BaseName: BaseName{kind == 0x10},
				Name:     name,
			}
		case 0x11, 0x12: // space and name specified by bytecode (why have multiple of this?)
			p.f.Name[i] = &RuntimeName{
				BaseName: BaseName{kind == 0x12},
			}
		case 0x09, 0x0e: // name is searched in a space set (much have only 1 match!)
			index := int(p.u30()) // zero is any (not default)
			name := "*"
			if index != 0 {
				name = p.f.String[index]
			}
			set := int(p.u30())
			if set == 0 {
				panic("Indexing a zero namespace set")
			}

			p.f.Name[i] = &NameSpaceQuery{
				BaseName: BaseName{kind == 0x0e},
				Name:     name,
				Set:      p.f.Set[set],
			}
		case 0x1b, 0x1c: // name specified by bytecode, searched in the set
			set := int(p.u30())
			if set == 0 {
				panic("Indexing a zero namespace set")
			}

			p.f.Name[i] = &RuntimeNameSpaceQuery{
				BaseName: BaseName{kind == 0x1c},
				Set:      p.f.Set[set],
			}
		default:
			panic(fmt.Errorf("Unknown namespace kind: %02X", kind))
		}
	}
}

func (p *parser) readMethodPool() {
	methodCount := int(p.u30())
	p.f.Method = make([]*Method, methodCount)
	for i := 0; i < methodCount; i++ {
		p.f.Method[i] = &Method{}
		m := p.f.Method[i]

		paramCount := int(p.u30())
		m.ReturnType = p.f.Name[int(p.u30())]
		m.Params = make([]MethodParam, paramCount)
		for j := 0; j < paramCount; j++ {
			m.Params[j].Type = p.f.Name[int(p.u30())]
		}

		m.Name = p.f.String[p.u30()]

		flags, _ := p.r.ReadByte()
		m.NeedsArguments = (flags & 0x01) != 0
		m.UsesActivation = (flags & 0x02) != 0
		m.NeedsRest = (flags & 0x04) != 0
		m.UsesDxns = (flags & 40) != 0

		if (flags & 0x08) != 0 { // optional values for arguments
			optionCount := int(p.u30())

			// options are for last arguments, left to right
			for j := paramCount-optionCount; j < paramCount; j++ {
				m.Params[j].Default.Index = p.u30()
				kind, err := p.r.ReadByte()
				if err != nil {
					p.error = err
				}
				m.Params[j].Default.Kind = kind
				if kind == 0x08 || kind == 0x16 || kind == 0x17 || kind == 0x18 || kind == 0x19 || kind == 0x1a || kind == 0x05 {
					fmt.Printf("%02x\n", kind)
					panic("todo: ns usage in options?")
				}
			}
		}

		if (flags & 0x80) != 0 { // param names specified
			for j := 0; j < paramCount; j++ {
				m.Params[j].Name = p.f.String[p.u30()]
			}
		}
	}
}

func (p *parser) skipMetadataPool() {
	metaCount := int(p.u30())
	for i := 0; i < metaCount; i++ {
		panic("panik!")
	}
}

func (p *parser) readClassPool() {
	classCount := int(p.u30())
	p.f.Class = make([]*Class, classCount)
	for i := 0; i < classCount; i++ {
		c := &Class{}
		p.f.Class[i] = c

		c.Name = p.f.Name[p.u30()].(*FullName)
		super := p.u30()
		if super != 0 {
			c.Super = p.f.Name[super]
		}

		flags, _ := p.r.ReadByte()
		c.Sealed = (flags & 0x01) != 0
		c.Final = (flags & 0x02) != 0
		c.Interface = (flags & 0x04) != 0
		if (flags & 0x08) != 0 {
			c.ProtectedSpace = p.f.Space[p.u30()]
		}

		interfaceCount := int(p.u30())
		c.Implements = make([]Name, interfaceCount)
		for j := 0; j < interfaceCount; j++ {
			c.Implements[j] = p.f.Name[p.u30()]
		}

		c.Init = p.f.Method[p.u30()]
		c.Traits = p.readTraits()
	}

	for i := 0; i < classCount; i++ {
		p.f.Class[i].StaticInit = p.f.Method[p.u30()]
		p.f.Class[i].StaticTraits = p.readTraits()
	}
}

func (p *parser) readTraits() []Trait {
	traitCount := int(p.u30())
	traits := make([]Trait, traitCount)
	for i := 0; i < traitCount; i++ {
		var base BaseTrait
		base.Name = p.f.Name[p.u30()].(*FullName)
		kind, _ := p.r.ReadByte()
		base.SlotId = p.u30()

		switch (kind & 0x0f) {
		case 0x00, 0x06: // variable, const
			slot := &SlotTrait{
				BaseTrait: base,
				Const:     (kind & 0x0f) == 0x06,
				Type:      p.f.Name[p.u30()],
				Index:     p.u30(),
			}
			if slot.Index != 0 {
				slot.ValueKind, _ = p.r.ReadByte()
				kind := slot.ValueKind
				if kind == 0x08 || kind == 0x16 || kind == 0x17 || kind == 0x18 || kind == 0x19 || kind == 0x1a || kind == 0x05 {
					if kind == 0x08 && base.Name.Name == "mx_internal" {
						// ok
					} else {
						fmt.Println(base.Name)
						fmt.Printf("%02x\n", kind)
						panic("todo: ns usage in slot?")
					}
				}
			}
			traits[i] = slot
		case 0x01, 0x02, 0x03: // method, getter, setter
			traits[i] = &MethodTrait{
				BaseTrait: base,
				Type:      MethodType((kind & 0x0f) - 1),
				Method:    p.f.Method[p.u30()],
				Final:     (kind & 0x10) != 0,
				Override:  (kind & 0x20) != 0,
			}
		case 0x04: // class
			traits[i] = &ClassTrait{
				BaseTrait: base,
				Class:     p.f.Class[p.u30()],
			}
		default:
			panic(fmt.Errorf("Unknown trait kind: %02X", kind & 0x0f))
		}

		if (kind & 0x40) != 0 {
			panic("trait metadata not implemented")
		}
	}
	return traits
}

func (p *parser) readScriptPool() {
	scriptCount := int(p.u30())
	p.f.Script = make([]*Script, scriptCount)
	for i := 0; i < scriptCount; i++ {
		p.f.Script[i] = &Script{
			Init:   p.f.Method[p.u30()],
			Traits: p.readTraits(),
		}
	}
}

func (p *parser) readBodyPool() {
	bodyCount := int(p.u30())
	p.f.Body = make([]*Body, bodyCount)
	for i := 0; i < bodyCount; i++ {
		p.f.Body[i] = &Body{
			Method:     p.f.Method[p.u30()],
			StackSize:  p.u30(),
			Locals:     p.u30(),
			MinScope:   p.u30(),
			MaxScope:   p.u30(),
			Code:       p.readPrefixedByteArray(),
			Exceptions: p.readBodyExceptions(),
			Traits:     p.readTraits(),
		}
		if p.f.Body[i].Method.Body != nil {
			panic("method got a duplicate body?")
		}
		p.f.Body[i].Method.Body = p.f.Body[i]
	}
}

func (p *parser) readBodyExceptions() []Exception {
	exceptionCount := int(p.u30())
	exceptions := make([]Exception, exceptionCount)
	for i := 0; i < exceptionCount; i++ {
		exceptions[i].From = p.u30()
		exceptions[i].To = p.u30()
		exceptions[i].Target = p.u30()
		index := p.u30()
		if index != 0 {
			exceptions[i].Type = p.f.String[index]
		} else {
			exceptions[i].Type = "*"
		}
		exceptions[i].VarName = p.f.String[p.u30()]
	}
	return exceptions
}

func (p *parser) readPrefixedByteArray() []byte {
	size := int(p.u30())
	data := make([]byte, size)
	_, err := p.r.Read(data)
	if err != nil {
		p.error = err
	}
	return data
}

func (p *parser) u30() uint32 {
	var result uint32
	// 7 bits
	{
		x, err := p.r.ReadByte()
		result = uint32(x & 0x7f)
		if (x & 0x80) == 0 {
			p.error = err
			return result
		}
	}

	// 14 bits
	{
		x, err := p.r.ReadByte()
		result |= uint32(x & 0x7f) << 7
		if (x & 0x80) == 0 {
			p.error = err
			return result
		}
	}

	// 21 bits
	{
		x, err := p.r.ReadByte()
		result |= uint32(x & 0x7f) << 14
		if (x & 0x80) == 0 {
			p.error = err
			return result
		}
	}

	// 28 bits
	{
		x, err := p.r.ReadByte()
		result |= uint32(x & 0x7f) << 21
		if (x & 0x80) == 0 {
			p.error = err
			return result
		}
	}

	// 30 bits
	{
		x, err := p.r.ReadByte()
		result |= uint32(x & 0x7f) << 28
		if (x & 0xfc) != 0 {
			panic("unexpected data after 30 bits")
		}
		p.error = err
		return result
	}
}

func (p *parser) s32() int32 {
	return int32(p.u32())
}

func (p *parser) u32() uint32 {
	var result uint32
	// 7 bits
	{
		x, err := p.r.ReadByte()
		result = uint32(x & 0x7f)
		if (x & 0x80) == 0 {
			p.error = err
			return result
		}
	}

	// 14 bits
	{
		x, err := p.r.ReadByte()
		result |= uint32(x & 0x7f) << 7
		if (x & 0x80) == 0 {
			p.error = err
			return result
		}
	}

	// 21 bits
	{
		x, err := p.r.ReadByte()
		result |= uint32(x & 0x7f) << 14
		if (x & 0x80) == 0 {
			p.error = err
			return result
		}
	}

	// 28 bits
	{
		x, err := p.r.ReadByte()
		result |= uint32(x & 0x7f) << 21
		if (x & 0x80) == 0 {
			p.error = err
			return result
		}
	}

	// 32 bits
	{
		x, err := p.r.ReadByte()
		result |= uint32(x & 0x7f) << 28
		if (x & 0xf0) != 0 {
			panic("unexpected data after 32 bits")
		}
		p.error = err
		return result
	}
}
