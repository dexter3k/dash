package abc

type undefinedType struct{}
type nullType struct{}

var Undefined = new(undefinedType)
var Null = new(nullType)

// Only used as a namespace type wrapper for File.GetConstant
type Namespace string

// Note: we handle all namespaces as prefixed strings, with prefix
// being a single-character tag marking ns type
// This is ok, because we don't usually care about specific space
// types, but we do have to differentite them

// For data of known length use slices of values, not pointers
// to keep allocator fragmentation to minimum and ease work for GC
// strings and other golang-primitives are all used by value in
// sub-data like namespace, even though a pointer would be more
// efficient for strings (kinda, since strings are immutable, there
// is surely optimisations done about this)
type File struct {
	Magic uint32

	Ints    []int32
	Uints   []uint32
	Doubles []float64
	Strings []string

	Spaces    []string
	SpaceSets [][]string
	Links     []Multiname

	Metadata []Metadata
	Methods  []Method
	Bodies   []Body
	Classes  []Class
	Scripts  []Script

	Raw []byte

	stream *abcStream
}

func Parse(buffer []byte) *File {
	abc := &File{
		Raw:    buffer,
		stream: NewAbcStream(buffer),
	}

	abc.checkMagic()

	abc.parseConstantPool()
	abc.parseNamespaces()
	abc.parseNamespaceSets()
	abc.parseMultinames()

	abc.parseMethodInfos()
	abc.parseMetaData()
	abc.parseInstanceAndClassInfos()
	abc.parseScriptInfos()
	abc.parseMethodBodyInfos()

	abc.stream = nil

	return abc
}

// Kind: constant kind from AVM2 Overview paper
// Index: index into the namespace pool
// Return types: int32, uint32, float64, string, bool,
// *nullKind, *undefinedKind, *Namespace
// Panics for unknown type
func (abc *File) GetConstant(kind, index int) interface{} {
	switch kind {
	case kInt:
		return abc.Ints[index]
	case kUint:
		return abc.Uints[index]
	case kDouble:
		return abc.Doubles[index]
	case kUtf8:
		return abc.Strings[index]
	case kTrue:
		return true
	case kFalse:
		return false
	case kNull:
		return Null
	case kUndefined:
		return Undefined
	case kNamespace, kPackageNamespace, kPackageInternalNamespace,
		kProtectedNamespace, kExplicitNamespace, kStaticProtectedNamespace,
		kPrivateNamespace:
		return (*Namespace)(&abc.Spaces[index])
	default:
		panic(kind)
	}
}

// SubType field is for TypeName multiname (i.e. Vector<T>)
// Currently only one subtype is possible (no stuff like Vector<X, Y, Z>)
// So we don't even bother using slice for now (for... ever?)
// Imagine including 16 bytes into struct just to have some fucking vector there
type Multiname struct {
	Name    string   // "" if name is runtime
	Spaces  []string // nil if space(s) is/are runtime
	SubType *Multiname

	IsAttribute    bool
	IsNameRuntime  bool
	IsSpaceRuntime bool

	typeName    int
	mangledName *string
}

func QualifiedName(space, name string) *Multiname {
	return &Multiname{
		Name:   name,
		Spaces: []string{space},
	}
}

func (mn *Multiname) IsRuntime() bool {
	return mn.IsNameRuntime || mn.IsSpaceRuntime
}

func (mn *Multiname) GetMangledName() string {
	if mn.mangledName == nil {
		return mn.mangleName()
	}
	return *mn.mangledName
}

// We don't store typed traits in the same map as dynamics properties
// So we only care about accidental (or forged) name and space collision
func (mn *Multiname) mangleName() string {
	if len(mn.Spaces) != 1 {
		panic("Tried to mangle non-qname")
	}

	mangledName := mn.Spaces[0] + ":" + mn.Name
	if !mn.IsNameRuntime && !mn.IsSpaceRuntime {
		mn.mangledName = &mangledName
	}
	return mangledName
}

type MethodParam struct {
	Name     string
	Type     *Multiname
	Optional interface{}
}

type Method struct {
	Type   *Multiname
	Params []MethodParam
	Flags  uint8
	Name   string

	// can be nil
	Body *Body
}

type MetadataEntry struct {
	Key   string
	Value string
}

type Metadata struct {
	Name   string
	Values []MetadataEntry
}

type TraitKind int

const (
	TraitClass    TraitKind = iota // 0
	TraitMethod                    // 1
	TraitGetter                    // 2
	TraitSetter                    // 3
	TraitVariable                  // 4
	TraitConst                     // 5
)

type Trait struct {
	Link *Multiname // always qualified
	Kind TraitKind

	Slot int // For Method, Getter and Setter its a DispId for VTable!

	// Only for TraitClass
	Class *Class
	// TraitMethod, TraitGetter, TraitSetter
	Method *Method
	// TraitConst, TraitVariable
	DataType *Multiname
	Value    interface{}

	IsFinal    bool
	IsOverride bool

	Metadata []*Metadata
}

type Class struct {
	Name           *Multiname // always qualified
	Super          *Multiname // empty for no base name
	ProtectedSpace string
	Implements     []*Multiname

	IsSealed    bool
	IsFinal     bool
	IsInterface bool

	Constructor       *Method
	StaticConstructor *Method

	Traits       []Trait
	StaticTraits []Trait
}

type Script struct {
	Script *Method
	Traits []Trait
}

type Exception struct {
	From   int // First byte of exception area
	To     int // Last byte of exception area
	Target int // Index at which jump is performed in case of exception
	// Daym, avm2overview.pdf is a bug on a bug
	// VarType string // Type of the exception
	// VarName string // Name of the exception
	VarType *Multiname
	VarName *Multiname // only if not 46<<16|15
}

type Instruction struct {
	Offset int
	Raw    []byte

	Opcode uint8
	Args   []int

	// Liveness/reachability analysis
	Ins  []int
	Outs []int
	Dead bool
}

type Body struct {
	Method     *Method
	Code       []byte
	Exceptions []Exception

	StackLimit    int
	LocalLimit    int
	MinScopeLimit int
	MaxScopeLimit int

	// What traits are defined for bodies?
	Traits []Trait

	// Not used by avm2, but by accompanying toolset
	// eg debugger
	disassembled []Instruction
	// This in practice either means we have a bug
	// or this body won't even pass first stage of
	// verifier
	disassemblyFailed bool
}

func (b *Body) OnChanged() {
	b.disassembled = nil
	b.disassemblyFailed = false
}

func Assemble(instr *Instruction) []byte {
	return asm(instr)
}

// func AssemblyMany(ops []Instruction) []byte {
// 	return asm_many(ops)
// }

func (b *Body) Disassemble() []Instruction {
	if b.disassembled != nil {
		return b.disassembled
	}
	if b.disassemblyFailed {
		return nil
	}

	result := unasm(b.Code, b.Exceptions)
	if result == nil {
		b.disassemblyFailed = true
		return nil
	}
	b.disassembled = result
	return b.disassembled
}
