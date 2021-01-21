package opcodes

// Keep theese jpexs compatible for clarity
type Arg uint32

const (
	NoArgument Arg = 0x0000

	U30          = 0x100
	U8           = 0x200
	S24          = 0x300
	OffsetLookup = 0x400

	ByteValue      = U8 | 0x00
	ScopeByteIndex = U8 | 0x0A
	ShortValue     = U30 | 0x00
	LinkIndex      = U30 | 0x01
	ArgCount       = U30 | 0x02
	MethodIndex    = U30 | 0x03
	StringIndex    = U30 | 0x04
	LocalIndex     = U30 | 0x08
	SlotIndex      = U30 | 0x09
	DoubleIndex    = U30 | 0x10
	ScopeIndex     = U30 | 0x0A
	RelativeOffset = S24 | 0x0B
	ExceptionIndex = U30 | 0x0C
	ClassIndex     = U30 | 0x0D
	IntIndex       = U30 | 0x0E
	UintIndex      = U30 | 0x0F

	DataTypeMask = 0xff00
)

var Arguments = map[uint8][]Arg{
	NoOperation: []Arg{},
	CoerceToAny: []Arg{},
	Label:       []Arg{},
	Debug:       []Arg{ByteValue, StringIndex, LocalIndex, U30},
	DebugLine:   []Arg{U30},
	DebugFile:   []Arg{StringIndex},

	Throw:       []Arg{},
	ReturnVoid:  []Arg{},
	ReturnValue: []Arg{},

	PushNull:       []Arg{},
	PushUndefined:  []Arg{},
	PushNotANumber: []Arg{},
	PushTrue:       []Arg{},
	PushFalse:      []Arg{},
	PushInt8:       []Arg{ByteValue},
	PushInt16:      []Arg{ShortValue},
	PushInt:        []Arg{IntIndex},
	PushUint:       []Arg{UintIndex},
	PushNumber:     []Arg{DoubleIndex},
	PushString:     []Arg{StringIndex},

	Pop:       []Arg{},
	Swap:      []Arg{},
	Duplicate: []Arg{},

	GetLocal_0: []Arg{},
	GetLocal_1: []Arg{},
	GetLocal_2: []Arg{},
	GetLocal_3: []Arg{},
	SetLocal_0: []Arg{},
	SetLocal_1: []Arg{},
	SetLocal_2: []Arg{},
	SetLocal_3: []Arg{},
	GetLocal:   []Arg{LocalIndex},
	SetLocal:   []Arg{LocalIndex},
	EraseLocal: []Arg{LocalIndex},

	PushScope:      []Arg{},
	PopScope:       []Arg{},
	PushWithScope:  []Arg{},
	GetGlobalScope: []Arg{},
	GetScopeObject: []Arg{ScopeByteIndex},

	Add:             []Arg{},
	AddInteger:      []Arg{},
	Subtract:        []Arg{},
	SubtractInteger: []Arg{},
	Multiply:        []Arg{},
	MultiplyInteger: []Arg{},
	Divide:          []Arg{},
	Modulo:          []Arg{},
	Negate:          []Arg{},
	NegateInteger:   []Arg{},
	BooleanNot:      []Arg{},

	BitwiseAnd:         []Arg{},
	BitwiseOr:          []Arg{},
	BitwiseNot:         []Arg{},
	BitwiseXor:         []Arg{},
	ShiftLeft:          []Arg{},
	ShiftRightSigned:   []Arg{},
	ShiftRightUnsigned: []Arg{},

	Increment:             []Arg{},
	Decrement:             []Arg{},
	IncrementLocal:        []Arg{LocalIndex},
	DecrementLocal:        []Arg{LocalIndex},
	IncrementInteger:      []Arg{},
	DecrementInteger:      []Arg{},
	IncrementLocalInteger: []Arg{LocalIndex},
	DecrementLocalInteger: []Arg{LocalIndex},

	ConvertToObject:   []Arg{},
	ConvertToString:   []Arg{},
	ConvertToNumber:   []Arg{},
	ConvertToBoolean:  []Arg{},
	ConvertToUnsigned: []Arg{},
	ConvertToInteger:  []Arg{},
	CoerceToType:      []Arg{LinkIndex},
	CoerceToString:    []Arg{},

	CompareEquals:        []Arg{},
	CompareStrictEquals:  []Arg{},
	CompareLessThan:      []Arg{},
	CompareLessEquals:    []Arg{},
	CompareGreaterEquals: []Arg{},
	CompareGreaterThan:   []Arg{},

	Jump:              []Arg{RelativeOffset},
	IfTrue:            []Arg{RelativeOffset},
	IfFalse:           []Arg{RelativeOffset},
	IfEqual:           []Arg{RelativeOffset},
	IfNotEqual:        []Arg{RelativeOffset},
	IfStrictEqual:     []Arg{RelativeOffset},
	IfStrictNotEqual:  []Arg{RelativeOffset},
	IfLessThan:        []Arg{RelativeOffset},
	IfNotLessThan:     []Arg{RelativeOffset},
	IfLessEqual:       []Arg{RelativeOffset},
	IfNotLessEqual:    []Arg{RelativeOffset},
	IfGreaterThan:     []Arg{RelativeOffset},
	IfNotGreaterThan:  []Arg{RelativeOffset},
	IfGreaterEqual:    []Arg{RelativeOffset},
	IfNotGreaterEqual: []Arg{RelativeOffset},
	LookupSwitch:      []Arg{RelativeOffset, OffsetLookup},

	NewFunction:   []Arg{MethodIndex},
	NewObject:     []Arg{ArgCount},
	NewClass:      []Arg{ClassIndex},
	NewArray:      []Arg{ArgCount},
	NewActivation: []Arg{},
	NewCatch:      []Arg{ExceptionIndex},

	ApplyVectorType:        []Arg{ArgCount},
	CheckType:              []Arg{ArgCount},
	CheckTypeLate:          []Arg{},
	IsObjectOfType:         []Arg{LinkIndex},
	IsObjectOfTypeLate:     []Arg{},
	CheckFilter:            []Arg{},
	ObjectHasProperty:      []Arg{},
	InstanceOf:             []Arg{},
	GetTypeName:            []Arg{},
	GetNextPropertyIndex:   []Arg{},
	GetNextPropertyOnLocal: []Arg{LocalIndex, LocalIndex},
	GetNameOfNextProperty:  []Arg{},
	GetValueOfNextProperty: []Arg{},
	SetDefaultXmlSpace:     []Arg{StringIndex},
	SetDefaultXmlSpaceLate: []Arg{},

	SetSlot:            []Arg{SlotIndex},
	GetSlot:            []Arg{SlotIndex},
	GetSuper:           []Arg{LinkIndex},
	SetSuper:           []Arg{LinkIndex},
	SetProperty:        []Arg{LinkIndex},
	GetProperty:        []Arg{LinkIndex},
	InitProperty:       []Arg{LinkIndex},
	DeleteProperty:     []Arg{LinkIndex},
	FindProperty:       []Arg{LinkIndex},
	FindPropertyStrict: []Arg{LinkIndex},
	FindAndGetProperty: []Arg{LinkIndex},
	GetDescendants:     []Arg{LinkIndex},
	// FindScriptDefinition: []Arg{} // todo

	Call:                []Arg{ArgCount},
	Construct:           []Arg{ArgCount},
	ConstructSuper:      []Arg{ArgCount},
	ConstructProperty:   []Arg{LinkIndex, ArgCount},
	CallPropertyVoid:    []Arg{LinkIndex, ArgCount},
	CallProperty:        []Arg{LinkIndex, ArgCount},
	FindAndCallProperty: []Arg{LinkIndex, ArgCount},
	CallSuperVoid:       []Arg{LinkIndex, ArgCount},
	CallSuper:           []Arg{LinkIndex, ArgCount},
}
