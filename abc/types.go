package abc

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

type MethodArgument struct {
	Name    string
	Type    Name
	Default []DefaultValue
}

type Method struct {
	Name       string
	Arguments  []MethodArgument
	ReturnType Name

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
	Name   FullName
	SlotId uint32
}

type SlotTrait struct {
	BaseTrait

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

type MethodTrait struct {
	BaseTrait

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

	Int       []int32
	Uint      []uint32
	Number    []float64
	String    []string
	Space     []Space
	Set       [][]Space
	Name      []Name

	Method []*Method
	Class  []*Class
	Script []*Script
	Body   []*Body
}