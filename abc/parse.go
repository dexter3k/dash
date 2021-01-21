package abc

import (
	// "fmt"
	"math"
)

const (
	kUndefined                      = 0x00
	kUtf8                           = 0x01
	kFloat                          = 0x02
	kInt                            = 0x03
	kUint                           = 0x04
	kPrivateNamespace               = 0x05
	kDouble                         = 0x06
	kQualified                      = 0x07
	kNamespace                      = 0x08
	kMultiname                      = 0x09
	kFalse                          = 0x0A
	kTrue                           = 0x0B
	kNull                           = 0x0C
	kQualifiedAttribute             = 0x0D
	kMultinameAttribute             = 0x0E
	kRealTimeQualified              = 0x0F
	kRealTimeQualifiedAttribute     = 0x10
	kRealTimeQualifiedLink          = 0x11
	kRealTimeQualifiedLinkAttribute = 0x12
	kNameLink                       = 0x13
	kNameLinkAttribute              = 0x14
	kNamespaceSet                   = 0x15
	kPackageNamespace               = 0x16
	kPackageInternalNamespace       = 0x17
	kProtectedNamespace             = 0x18
	kExplicitNamespace              = 0x19
	kStaticProtectedNamespace       = 0x1A
	kMultinameLink                  = 0x1B
	kMultinameLinkAttribute         = 0x1C
	kTypeName                       = 0x1D

	kClassSealed             = 0x01
	kClassFinal              = 0x02
	kClassInterface          = 0x04
	kClassProtectedNamespace = 0x08

	kNeedArguments = 0x01
	kActivation    = 0x02
	kNeedRest      = 0x04
	kHasOptional   = 0x08
	kIgnoreRest    = 0x10
	kNative        = 0x20
	kSetDXNS       = 0x40
	kHasParamNames = 0x80

	kSlot   = 0x0
	kMethod = 0x1
	kGetter = 0x2
	kSetter = 0x3
	kClass  = 0x4
	kConst  = 0x6

	kIsFinal     = 0x1
	kIsOverride  = 0x2
	kHasMetadata = 0x4
)

func (abc *File) checkMagic() {
	abc.Magic = abc.stream.ReadU32()
	if abc.Magic < ((46 << 16) | 15) {
		panic("VerifyError")
	}
}

func (abc *File) parseConstantPool() {
	abc.parseNumericConstants()
	abc.parseStringConstants()
}

func (abc *File) parseNumericConstants() {
	var n int
	s := abc.stream

	// fmt.Println("Parsing abc")

	// int32
	n = int(s.ReadVarlenU30())
	if n == 0 {
		n = 1
	}
	abc.Ints = make([]int32, n, n)
	abc.Ints[0] = 0
	for i := 1; i < n; i++ {
		abc.Ints[i] = s.ReadVarlenS32()
	}

	// fmt.Println("Parsed", len(abc.Ints), "ints")

	// uint32
	n = int(s.ReadVarlenU30())
	if n == 0 {
		n = 1
	}
	abc.Uints = make([]uint32, n, n)
	abc.Uints[0] = 0
	for i := 1; i < n; i++ {
		abc.Uints[i] = s.ReadVarlenU32()
	}

	// float64
	n = int(s.ReadVarlenU30())
	if n == 0 {
		n = 1
	}
	abc.Doubles = make([]float64, n, n)
	abc.Doubles[0] = math.NaN()
	for i := 1; i < n; i++ {
		abc.Doubles[i] = math.Float64frombits(s.ReadU64())
	}
}

func (abc *File) parseStringConstants() {
	var n int
	s := abc.stream

	// strings
	n = int(s.ReadVarlenU30())
	if n == 0 {
		n = 1
	}
	abc.Strings = make([]string, n, n)
	abc.Strings[0] = ""
	for i := 1; i < n; i++ {
		abc.Strings[i] = s.ReadPrefixedString()
	}
}

// var (
// 	nsPrefixes = [...]string{
// 		"P", // Public
// 		"I", // package Internal
// 		"R", // pRotected
// 		"E", // Explicit
// 		"S", // Static protected
// 		"V", // priVate
// 	}
// )

func (abc *File) parseNamespaces() {
	var n int
	s := abc.stream

	n = int(s.ReadVarlenU30())
	if n == 0 {
		n = 1
	}
	abc.Spaces = make([]string, n, n)
	abc.Spaces[0] = "P"
	for i := 1; i < n; i++ {
		kind := s.ReadU8()
		name := abc.Strings[s.ReadVarlenU30()]
		switch kind {
		case kNamespace, kPackageNamespace:
			abc.Spaces[i] = "P" + name
		case kPackageInternalNamespace:
			abc.Spaces[i] = "I" + name
		case kProtectedNamespace:
			abc.Spaces[i] = "R" + name
		case kExplicitNamespace:
			abc.Spaces[i] = "E" + name
		case kStaticProtectedNamespace:
			abc.Spaces[i] = "S" + name
		case kPrivateNamespace:
			abc.Spaces[i] = "V" + name
		default:
			panic("Unknown namespace kind")
		}
	}
}

func (abc *File) parseNamespaceSets() {
	var n int
	s := abc.stream

	n = int(s.ReadVarlenU30())
	if n == 0 {
		n = 1
	}
	abc.SpaceSets = make([][]string, n, n)
	abc.SpaceSets[0] = []string{}
	for i := 1; i < n; i++ {
		setSize := int(s.ReadVarlenU30())
		abc.SpaceSets[i] = make([]string, setSize, setSize)
		for j := 0; j < setSize; j++ {
			index := int(s.ReadVarlenU30())
			abc.SpaceSets[i][j] = abc.Spaces[index]
		}
	}
}

func (abc *File) parseMultinames() {
	var n int
	s := abc.stream

	n = int(s.ReadVarlenU30())
	if n == 0 {
		n = 1
	}
	abc.Links = make([]Multiname, n, n)
	abc.Links[0] = Multiname{
		Spaces: []string{
			abc.Spaces[0],
		},
	}

	typenameOffsets := []int{}
	typenameIndexes := []int{}
	for i := 1; i < n; i++ {
		kind := s.ReadU8()
		var mn Multiname
		// Catch attributes first as they have no direct effect on parsing
		switch kind {
		case kQualifiedAttribute:
			mn.IsAttribute = true
			kind = kQualified
		case kRealTimeQualifiedAttribute:
			mn.IsAttribute = true
			kind = kRealTimeQualified
		case kRealTimeQualifiedLinkAttribute:
			mn.IsAttribute = true
			kind = kRealTimeQualifiedLink
		case kMultinameAttribute:
			mn.IsAttribute = true
			kind = kMultiname
		case kMultinameLinkAttribute:
			mn.IsAttribute = true
			kind = kMultinameLink
		}

		switch kind {
		case kQualified:
			spaceIndex := s.ReadVarlenU30()
			mn.Spaces = []string{
				abc.Spaces[spaceIndex],
			}

			nameIndex := s.ReadVarlenU30()
			mn.Name = abc.Strings[nameIndex]
		case kRealTimeQualified:
			mn.IsSpaceRuntime = true

			nameIndex := s.ReadVarlenU30()
			mn.Name = abc.Strings[nameIndex]
		case kRealTimeQualifiedLink:
			mn.IsSpaceRuntime = true
			mn.IsNameRuntime = true
		case kMultiname:
			nameIndex := s.ReadVarlenU30()
			mn.Name = abc.Strings[nameIndex]

			setIndex := s.ReadVarlenU30()
			mn.Spaces = abc.SpaceSets[setIndex]
		case kMultinameLink:
			mn.IsNameRuntime = true

			setIndex := s.ReadVarlenU30()
			mn.Spaces = abc.SpaceSets[setIndex]
		case kTypeName:
			typenameOffsets = append(typenameOffsets, s.GetOffset())
			typenameIndexes = append(typenameIndexes, i)
			// in "Foo<Bar>" its the Foo
			index := int(s.ReadVarlenU30())
			// We have more
			if index >= n {
				panic("Nope")
			}
			// Save this so we can identify this mn later
			mn.typeName = index

			// We only support single-subtype names
			count := int(s.ReadVarlenU30())
			if count != 1 {
				panic("Nope")
			}

			// in "Foo<Bar>" this is the Bar
			// and this might be recursive
			// Like "Foo<Bar<Baz>>"
			subtype := int(s.ReadVarlenU30())
			if subtype >= n {
				panic("Nope")
			}
		}
		abc.Links[i] = mn
	}

	whereToGoBack := s.GetOffset()
	for len(typenameOffsets) != 0 {
		s.SetOffset(typenameOffsets[len(typenameOffsets)-1])
		typenameOffsets = typenameOffsets[:len(typenameOffsets)-1]

		index := typenameIndexes[len(typenameIndexes)-1]
		typenameIndexes = typenameIndexes[:len(typenameIndexes)-1]

		mn := &abc.Links[index]
		if mn.typeName == 0 {
			panic("Nope")
		}

		name := &abc.Links[mn.typeName]
		if name.typeName != 0 || name.SubType != nil {
			panic("Nope")
		}

		// Just copy it
		abc.Links[index] = *name
		// Make sure no other types are there, AVMPlus does not
		// enforce this, but also does not support anything else
		if abc.Links[index].Name != "Vector" {
			panic(abc.Links[index].Name)
		} else if abc.Links[index].Spaces[0] != "P__AS3__.vec" {
			panic(abc.Links[index].Spaces[0])
		}

		s.ReadVarlenU30() // name index
		s.ReadVarlenU30() // count value
		subType := int(s.ReadVarlenU30())
		abc.Links[index].SubType = &abc.Links[subType]
		if abc.Links[index].SubType.typeName != 0 {
			// panic("Nope")
			// too lazy, fuck it if we'll go into recursion
		}
	}

	// Now check for recursion in Vector.<Vector.<Vector.<...>>> like
	// insane stuff
	// for i := 0; i < n; i++ {
	// 	link := &abc.Links[i]
	// 	if link.SubType == nil {
	// 		continue
	// 	}
	// }
	s.SetOffset(whereToGoBack)
}

func (abc *File) parseMethodInfos() {
	var n int
	s := abc.stream

	n = int(s.ReadVarlenU30())
	abc.Methods = make([]Method, n, n)
	for i := 0; i < n; i++ {
		paramCount := int(s.ReadVarlenU30())
		returnType := int(s.ReadVarlenU30())
		abc.Methods[i].Type = &abc.Links[returnType]

		abc.Methods[i].Params = make([]MethodParam, paramCount, paramCount)
		for j := 0; j < paramCount; j++ {
			paramType := int(s.ReadVarlenU30())
			abc.Methods[i].Params[j].Type = &abc.Links[paramType]
		}

		nameIndex := int(s.ReadVarlenU30())
		abc.Methods[i].Name = abc.Strings[nameIndex]

		flags := s.ReadU8()
		abc.Methods[i].Flags = flags

		if (flags & kHasOptional) != 0 {
			count := int(s.ReadVarlenU30())
			for j := paramCount - count; j < paramCount; j++ {
				poolIndex := s.ReadVarlenU30()
				abc.Methods[i].Params[j].Optional = abc.GetConstant(
					int(s.ReadU8()), int(poolIndex))
			}
		}

		if (flags & kHasParamNames) != 0 {
			for j := 0; j < paramCount; j++ {
				index := int(s.ReadVarlenU30())
				if index < len(abc.Strings[index]) {
					abc.Methods[i].Params[j].Name = abc.Strings[index]
				}
			}
		}
	}
}

func (abc *File) parseMetaData() {
	var n int
	s := abc.stream

	n = int(s.ReadVarlenU30())
	abc.Metadata = make([]Metadata, n, n)
	for i := 0; i < n; i++ {
		m := &abc.Metadata[i]
		m.Name = abc.Strings[int(s.ReadVarlenU30())]
		count := int(s.ReadVarlenU30())
		m.Values = make([]MetadataEntry, count, count)
		for j := 0; j < count; j++ {
			m.Values[j].Key = abc.Strings[int(s.ReadVarlenU30())]
		}
		for j := 0; j < count; j++ {
			m.Values[j].Value = abc.Strings[int(s.ReadVarlenU30())]
		}
	}
}

func (abc *File) parseTraits() []Trait {
	var n int
	s := abc.stream

	n = int(s.ReadVarlenU30())
	t := make([]Trait, n, n)
	for i := 0; i < n; i++ {
		nameIndex := int(s.ReadVarlenU30())
		t[i].Link = &abc.Links[nameIndex]

		kindAndFlags := s.ReadU8()
		kind := kindAndFlags & 0x0F
		flags := kindAndFlags >> 4

		t[i].Slot = int(s.ReadVarlenU30())

		switch kind {
		case kSlot, kConst:
			dataType := int(s.ReadVarlenU30())
			t[i].DataType = &abc.Links[dataType]
			poolIndex := int(s.ReadVarlenU30())
			if poolIndex != 0 {
				t[i].Value = abc.GetConstant(int(s.ReadU8()), poolIndex)
			}
			if kind == kSlot {
				t[i].Kind = TraitVariable
			} else {
				t[i].Kind = TraitConst
			}
		case kClass:
			classIndex := int(s.ReadVarlenU30())
			t[i].Class = &abc.Classes[classIndex]

			t[i].Kind = TraitClass
		case kMethod, kGetter, kSetter:
			methodIndex := int(s.ReadVarlenU30())
			t[i].Method = &abc.Methods[methodIndex]

			switch kind {
			case kMethod:
				t[i].Kind = TraitMethod
			case kGetter:
				t[i].Kind = TraitGetter
			case kSetter:
				t[i].Kind = TraitSetter
			default:
				panic("VerifyError: internal error")
			}
		default:
			panic("VerifyError: unknown trait kind")
		}

		t[i].IsFinal = (flags & kIsFinal) != 0
		t[i].IsOverride = (flags & kIsOverride) != 0

		metaCount := 0
		if (flags & kHasMetadata) != 0 {
			metaCount = int(s.ReadVarlenU30())
		}
		t[i].Metadata = make([]*Metadata, metaCount, metaCount)
		for j := 0; j < metaCount; j++ {
			index := int(s.ReadVarlenU30())
			t[i].Metadata[j] = &abc.Metadata[index]
		}
	}
	return t
}

func (abc *File) parseInstanceAndClassInfos() {
	s := abc.stream
	n := int(s.ReadVarlenU30())
	abc.Classes = make([]Class, n, n)
	for i := 0; i < n; i++ {
		c := &abc.Classes[i]
		c.Name = &abc.Links[s.ReadVarlenU30()]
		c.Super = &abc.Links[s.ReadVarlenU30()]
		flags := s.ReadU8()
		c.IsSealed = (flags & 0x1) != 0
		c.IsFinal = (flags & 0x2) != 0
		c.IsInterface = (flags & 0x4) != 0
		// c.HasProtectedSpace = (flags & 0x8) != 0
		if (flags & 0x8) != 0 {
			c.ProtectedSpace = abc.Spaces[s.ReadVarlenU30()]
		}
		ifn := int(s.ReadVarlenU30())
		c.Implements = make([]*Multiname, ifn, ifn)
		for i := 0; i < ifn; i++ {
			c.Implements[i] = &abc.Links[s.ReadVarlenU30()]
		}
		c.Constructor = &abc.Methods[s.ReadVarlenU30()]
		c.Traits = abc.parseTraits()
	}
	for i := 0; i < n; i++ {
		abc.Classes[i].StaticConstructor = &abc.Methods[s.ReadVarlenU30()]
		abc.Classes[i].StaticTraits = abc.parseTraits()
	}
}

func (abc *File) parseScriptInfos() {
	s := abc.stream
	n := int(s.ReadVarlenU30())
	abc.Scripts = make([]Script, n, n)
	for i := 0; i < n; i++ {
		scriptIndex := int(s.ReadVarlenU30())
		abc.Scripts[i].Script = &abc.Methods[scriptIndex]
		abc.Scripts[i].Traits = abc.parseTraits()
	}
}

func (abc *File) parseExceptionInfo() []Exception {
	s := abc.stream
	n := int(s.ReadVarlenU30())
	e := make([]Exception, n, n)
	for i := 0; i < n; i++ {
		e[i].From = int(s.ReadVarlenU30())
		e[i].To = int(s.ReadVarlenU30())
		e[i].Target = int(s.ReadVarlenU30())
		varType := int(s.ReadVarlenU30())
		e[i].VarType = &abc.Links[varType]

		if abc.Magic != 46<<16|15 {
			varName := int(s.ReadVarlenU30())
			e[i].VarName = &abc.Links[varName]
		}
	}
	return e
}

func (abc *File) parseMethodBodyInfos() {
	s := abc.stream
	n := int(s.ReadVarlenU30())
	abc.Bodies = make([]Body, n, n)
	for i := 0; i < n; i++ {
		methodIndex := int(s.ReadVarlenU30())
		abc.Bodies[i].Method = &abc.Methods[methodIndex]
		abc.Methods[methodIndex].Body = &abc.Bodies[i]

		abc.Bodies[i].StackLimit = int(s.ReadVarlenU30())
		abc.Bodies[i].LocalLimit = int(s.ReadVarlenU30())
		abc.Bodies[i].MinScopeLimit = int(s.ReadVarlenU30())
		abc.Bodies[i].MaxScopeLimit = int(s.ReadVarlenU30())

		codeLength := int(s.ReadVarlenU30())
		abc.Bodies[i].Code = s.ReadByteSlice(codeLength)
		abc.Bodies[i].Exceptions = abc.parseExceptionInfo()
		abc.Bodies[i].Traits = abc.parseTraits()
	}
}
