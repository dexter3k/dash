package abc

const (
	Undefined = 0x00
	String    = 0x01
	Int       = 0x03
	Uint      = 0x04
	Double    = 0x06
	False     = 0x0a
	True      = 0x0b
	Null      = 0x0c
	// why do we have to list these? why?
	PublicNamespace          = 0x08
	PackageNamespace         = 0x15
	InternalNamespace        = 0x17
	ProtectedNamespace       = 0x18
	ExplicitNamespace        = 0x19
	StaticProtectedNamespace = 0x1a
	PrivateNamespace         = 0x05
)

const (
	PublicNamespacePrefix          = "P"
	PackageInternalNamespacePrefix = "I"
	ProtectedNamespacePrefix       = "R"
	ExplicitNamespacePrefix        = "E"
	StaticProtectedNamespacePrefix = "S"
	PrivateNamespacePrefix         = "V"
)


type Space string

type Name interface{}

type BaseName struct {
	Attribute bool // wtf?
}

// aka QName: Foo.Bar
type FullName struct {
	BaseName

	Name      string
	Space     Space
}

// aka RTQName: [getNamespace()].Bar
type RuntimeSpaceName struct {
	BaseName

	Name      string
}

// aka RTQNameL: [getNamespace()]["Bar"]
type RuntimeName struct {
	BaseName
}

// aka Multiname: [{current visibility}].Bar
type NameSpaceQuery struct {
	BaseName

	Name      string
	Set       []Space
}

// aka MultinameL: [{current visibility}]["Bar"]
type RuntimeNameSpaceQuery struct {
	BaseName

	Set       []Space
}

type DefaultValue struct {
	Kind  uint8
	Index uint32
}

type MethodParam struct {
	Name    string
	Type    Name
	Default DefaultValue
}

type Method struct {
	Name       string
	Params     []MethodParam
	ReturnType Name
	Body       *Body

	NeedsArguments bool
	UsesActivation bool
	NeedsRest      bool
	UsesDxns       bool
}

type Meta struct {
	// ignored by vm
}

type Trait interface{}

type BaseTrait struct {
	Name   *FullName
	SlotId uint32
}

type SlotTrait struct {
	BaseTrait

	Const     bool
	Type      Name
	Index     uint32
	ValueKind uint8
}

type ClassTrait struct {
	BaseTrait

	Class *Class
}

type FunctionTrait struct {
	BaseTrait

	Function *Method
}

type MethodType uint8
const (
	NormalMethod MethodType = iota
	GetterMethod
	SetterMethod
	GetterSetterMethod
)

type MethodTrait struct {
	BaseTrait

	Type     MethodType
	Method   *Method
	Final    bool
	Override bool
}

type Class struct {
	StaticInit   *Method
	StaticTraits []Trait

	Name           *FullName
	Super          Name
	Sealed         bool
	Final          bool
	Interface      bool
	ProtectedSpace Space
	Implements     []Name
	Init           *Method
	Traits         []Trait
}

type Script struct {
	Init   *Method
	Traits []Trait
}

type Exception struct {
	From    uint32
	To      uint32
	Target  uint32
	Type    string
	VarName string
}

type Body struct {
	Method     *Method
	StackSize  uint32
	Locals     uint32
	MinScope   uint32
	MaxScope   uint32
	Code       []byte
	Exceptions []Exception
	Traits     []Trait
}

type File struct {
	Minor, Major uint16

	Int       []int32   // [0] = 0
	Uint      []uint32  // [0] = 0
	Number    []float64 // [0] = NaN
	String    []string  // [0] = ""
	Space     []Space   // [0] = Public(any)
	Set       [][]Space // [0] = invalid!
	Name      []Name    // [0] = invalid!

	Method []*Method
	Class  []*Class
	Script []*Script
	Body   []*Body
}
