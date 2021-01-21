package opcodes

// Describe how opcode affects stack
type StackEffect struct {
	Pop  int
	Push int
	// Number of argument to use for source
	// 0 for not used
	NeedsLink  int
	NeedsArgs  int
	ObjectArgs int // pop twice more args
}

var StackEffects = map[uint8]*StackEffect{
	Breakpoint:             &StackEffect{Pop: 0, Push: 0}, // unsupported
	NoOperation:            &StackEffect{Pop: 0, Push: 0}, // unsupported
	Throw:                  &StackEffect{Pop: 1, Push: 0},
	GetSuper:               &StackEffect{Pop: 1, Push: 1, NeedsLink: 1},
	SetSuper:               &StackEffect{Pop: 2, Push: 0, NeedsLink: 1},
	SetDefaultXmlSpace:     &StackEffect{Pop: 0, Push: 0},
	SetDefaultXmlSpaceLate: &StackEffect{Pop: 1, Push: 0},
	EraseLocal:             &StackEffect{Pop: 0, Push: 0}, // local effect
	Label:                  &StackEffect{Pop: 0, Push: 0},
	LoadFloat16:            &StackEffect{Pop: 0, Push: 0}, // unsupported
	StoreFloat16:           &StackEffect{Pop: 0, Push: 0}, // unsupported
	IfNotLessThan:          &StackEffect{Pop: 2, Push: 0},
	IfNotLessEqual:         &StackEffect{Pop: 2, Push: 0},
	IfNotGreaterThan:       &StackEffect{Pop: 2, Push: 0},
	IfNotGreaterEqual:      &StackEffect{Pop: 2, Push: 0},
	Jump:                   &StackEffect{Pop: 0, Push: 0},
	IfTrue:                 &StackEffect{Pop: 1, Push: 0},
	IfFalse:                &StackEffect{Pop: 1, Push: 0},
	IfEqual:                &StackEffect{Pop: 2, Push: 0},
	IfNotEqual:             &StackEffect{Pop: 2, Push: 0},
	IfLessThan:             &StackEffect{Pop: 2, Push: 0},
	IfLessEqual:            &StackEffect{Pop: 2, Push: 0},
	IfGreaterThan:          &StackEffect{Pop: 2, Push: 0},
	IfGreaterEqual:         &StackEffect{Pop: 2, Push: 0},
	IfStrictEqual:          &StackEffect{Pop: 2, Push: 0},
	IfStrictNotEqual:       &StackEffect{Pop: 2, Push: 0},
	LookupSwitch:           &StackEffect{Pop: 1, Push: 0},
	PushWithScope:          &StackEffect{Pop: 1, Push: 0},
	PopScope:               &StackEffect{Pop: 0, Push: 0},
	GetNameOfNextProperty:  &StackEffect{Pop: 2, Push: 1},
	GetNextPropertyIndex:   &StackEffect{Pop: 2, Push: 1},
	PushNull:               &StackEffect{Pop: 0, Push: 1},
	PushUndefined:          &StackEffect{Pop: 0, Push: 1},
	PushFloat32:            &StackEffect{Pop: 0, Push: 1},
	GetValueOfNextProperty: &StackEffect{Pop: 2, Push: 1},
	PushInt8:               &StackEffect{Pop: 0, Push: 1},
	PushInt16:              &StackEffect{Pop: 0, Push: 1},
	PushTrue:               &StackEffect{Pop: 0, Push: 1},
	PushFalse:              &StackEffect{Pop: 0, Push: 1},
	PushNotANumber:         &StackEffect{Pop: 0, Push: 1},
	Pop:                    &StackEffect{Pop: 1, Push: 0},
	Duplicate:              &StackEffect{Pop: 1, Push: 2}, // dup = {x=pop(), push(x), push(x)}
	Swap:                   &StackEffect{Pop: 2, Push: 2}, // swap = {x=pop(), y=pop(), push(x), push(y)}
	PushString:             &StackEffect{Pop: 0, Push: 1},
	PushInt:                &StackEffect{Pop: 0, Push: 1},
	PushUint:               &StackEffect{Pop: 0, Push: 1},
	PushNumber:             &StackEffect{Pop: 0, Push: 1},
	PushScope:              &StackEffect{Pop: 1, Push: 0},
	PushNamespace:          &StackEffect{Pop: 0, Push: 1},
	GetNextPropertyOnLocal: &StackEffect{Pop: 0, Push: 1}, // local effect
	LoadInt8:               &StackEffect{Pop: 0, Push: 0}, // unsupported
	LoadInt16:              &StackEffect{Pop: 0, Push: 0}, // unsupported
	LoadInt32:              &StackEffect{Pop: 0, Push: 0}, // unsupported
	LoadFloat32:            &StackEffect{Pop: 0, Push: 0}, // unsupported
	LoadFloat64:            &StackEffect{Pop: 0, Push: 0}, // unsupported
	StoreInt8:              &StackEffect{Pop: 0, Push: 0}, // unsupported
	StoreInt16:             &StackEffect{Pop: 0, Push: 0}, // unsupported
	StoreInt32:             &StackEffect{Pop: 0, Push: 0}, // unsupported
	StoreFloat32:           &StackEffect{Pop: 0, Push: 0}, // unsupported
	StoreFloat64:           &StackEffect{Pop: 0, Push: 0}, // unsupported
	NewFunction:            &StackEffect{Pop: 0, Push: 1},
	Call:                   &StackEffect{Pop: 2, Push: 1, NeedsArgs: 1},
	Construct:              &StackEffect{Pop: 1, Push: 1, NeedsArgs: 1},
	CallMethod:             &StackEffect{Pop: 1, Push: 1, NeedsArgs: 2},
	CallStatic:             &StackEffect{Pop: 1, Push: 1, NeedsArgs: 2},
	CallSuper:              &StackEffect{Pop: 1, Push: 1, NeedsArgs: 2, NeedsLink: 1},
	CallProperty:           &StackEffect{Pop: 1, Push: 1, NeedsArgs: 2, NeedsLink: 1},
	ReturnVoid:             &StackEffect{Pop: 0, Push: 0},
	ReturnValue:            &StackEffect{Pop: 1, Push: 0},
	ConstructSuper:         &StackEffect{Pop: 1, Push: 0, NeedsArgs: 1},
	ConstructProperty:      &StackEffect{Pop: 1, Push: 1, NeedsArgs: 2, NeedsLink: 1},
	FindAndCallProperty:    &StackEffect{Pop: 1, Push: 1, NeedsArgs: 2, NeedsLink: 1},
	CallSuperVoid:          &StackEffect{Pop: 1, Push: 0, NeedsArgs: 2, NeedsLink: 1},
	CallPropertyVoid:       &StackEffect{Pop: 1, Push: 0, NeedsArgs: 2, NeedsLink: 1},
	SignExtend1To32:        &StackEffect{Pop: 0, Push: 0}, // unsupported
	SignExtend8To32:        &StackEffect{Pop: 0, Push: 0}, // unsupported
	SignExtend16To32:       &StackEffect{Pop: 0, Push: 0}, // unsupported
	ApplyVectorType:        &StackEffect{Pop: 1, Push: 1, NeedsArgs: 1},
	PushFloat16:            &StackEffect{Pop: 0, Push: 1},
	NewObject:              &StackEffect{Pop: 0, Push: 1, ObjectArgs: 1},
	NewArray:               &StackEffect{Pop: 0, Push: 1, NeedsArgs: 1},
	NewActivation:          &StackEffect{Pop: 0, Push: 1},
	NewClass:               &StackEffect{Pop: 1, Push: 1},
	GetDescendants:         &StackEffect{Pop: 1, Push: 1, NeedsLink: 1},
	NewCatch:               &StackEffect{Pop: 0, Push: 1},
	FindPropertyStrict:     &StackEffect{Pop: 0, Push: 1, NeedsLink: 1},
	FindProperty:           &StackEffect{Pop: 0, Push: 1, NeedsLink: 1},
	FindScriptDefinition:   &StackEffect{Pop: 0, Push: 1}, // NeedsLink but not runtime, don't pop
	FindAndGetProperty:     &StackEffect{Pop: 0, Push: 1}, // NeedsLink but not runtime, don't pop
	SetProperty:            &StackEffect{Pop: 2, Push: 0, NeedsLink: 1},
	GetLocal:               &StackEffect{Pop: 0, Push: 1}, // local effect
	SetLocal:               &StackEffect{Pop: 1, Push: 0}, // local effect
	GetGlobalScope:         &StackEffect{Pop: 0, Push: 1},
	GetScopeObject:         &StackEffect{Pop: 0, Push: 1},
	GetProperty:            &StackEffect{Pop: 1, Push: 1, NeedsLink: 1},
	GetOuterScope:          &StackEffect{Pop: 0, Push: 1},
	InitProperty:           &StackEffect{Pop: 2, Push: 0, NeedsLink: 1},
	DeleteProperty:         &StackEffect{Pop: 1, Push: 1, NeedsLink: 1},
	GetSlot:                &StackEffect{Pop: 1, Push: 1},
	SetSlot:                &StackEffect{Pop: 2, Push: 0},
	GetGlobalSlot:          &StackEffect{Pop: 0, Push: 1},
	SetGlobalSlot:          &StackEffect{Pop: 1, Push: 0},
	ConvertToString:        &StackEffect{Pop: 1, Push: 1},
	EscapeXmlElement:       &StackEffect{Pop: 1, Push: 1},
	EscapeXmlAttribute:     &StackEffect{Pop: 1, Push: 1},
	ConvertToInteger:       &StackEffect{Pop: 1, Push: 1},
	ConvertToUnsigned:      &StackEffect{Pop: 1, Push: 1},
	ConvertToNumber:        &StackEffect{Pop: 1, Push: 1},
	ConvertToBoolean:       &StackEffect{Pop: 1, Push: 1},
	ConvertToObject:        &StackEffect{Pop: 0, Push: 0}, // really does not affect stack
	CheckFilter:            &StackEffect{Pop: 0, Push: 0}, // same
	ConvertToFloat32:       &StackEffect{Pop: 0, Push: 0}, // unsupported
	UnaryPlusOperator:      &StackEffect{Pop: 0, Push: 0}, // unsupported
	ConvertToFloat16:       &StackEffect{Pop: 0, Push: 0}, // unsupported
	CoerceToType:           &StackEffect{Pop: 1, Push: 1},
	CoerceToBoolean:        &StackEffect{Pop: 1, Push: 1},
	CoerceToAny:            &StackEffect{Pop: 0, Push: 0},
	CoerceToInteger:        &StackEffect{Pop: 1, Push: 1},
	CoerceToNumber:         &StackEffect{Pop: 1, Push: 1},
	CoerceToString:         &StackEffect{Pop: 1, Push: 1},
	CheckType:              &StackEffect{Pop: 1, Push: 1},
	CheckTypeLate:          &StackEffect{Pop: 2, Push: 1},
	CoerceToUnsigned:       &StackEffect{Pop: 1, Push: 1},
	CoerceToObject:         &StackEffect{Pop: 1, Push: 1},
	Negate:                 &StackEffect{Pop: 1, Push: 1},
	Increment:              &StackEffect{Pop: 1, Push: 1},
	IncrementLocal:         &StackEffect{Pop: 0, Push: 0}, // local effect
	Decrement:              &StackEffect{Pop: 1, Push: 1},
	DecrementLocal:         &StackEffect{Pop: 0, Push: 0}, // local effect
	GetTypeName:            &StackEffect{Pop: 1, Push: 1},
	BooleanNot:             &StackEffect{Pop: 1, Push: 1},
	BitwiseNot:             &StackEffect{Pop: 1, Push: 1},
	Add:                    &StackEffect{Pop: 2, Push: 1},
	Subtract:               &StackEffect{Pop: 2, Push: 1},
	Multiply:               &StackEffect{Pop: 2, Push: 1},
	Divide:                 &StackEffect{Pop: 2, Push: 1},
	Modulo:                 &StackEffect{Pop: 2, Push: 1},
	ShiftLeft:              &StackEffect{Pop: 2, Push: 1},
	ShiftRightSigned:       &StackEffect{Pop: 2, Push: 1},
	ShiftRightUnsigned:     &StackEffect{Pop: 2, Push: 1},
	BitwiseAnd:             &StackEffect{Pop: 2, Push: 1},
	BitwiseOr:              &StackEffect{Pop: 2, Push: 1},
	BitwiseXor:             &StackEffect{Pop: 2, Push: 1},
	CompareEquals:          &StackEffect{Pop: 2, Push: 1},
	CompareStrictEquals:    &StackEffect{Pop: 2, Push: 1},
	CompareLessThan:        &StackEffect{Pop: 2, Push: 1},
	CompareLessEquals:      &StackEffect{Pop: 2, Push: 1},
	CompareGreaterThan:     &StackEffect{Pop: 2, Push: 1},
	CompareGreaterEquals:   &StackEffect{Pop: 2, Push: 1},
	InstanceOf:             &StackEffect{Pop: 2, Push: 1},
	IsObjectOfType:         &StackEffect{Pop: 1, Push: 1},
	IsObjectOfTypeLate:     &StackEffect{Pop: 2, Push: 1},
	ObjectHasProperty:      &StackEffect{Pop: 2, Push: 1},
	IncrementInteger:       &StackEffect{Pop: 1, Push: 1},
	DecrementInteger:       &StackEffect{Pop: 1, Push: 1},
	IncrementLocalInteger:  &StackEffect{Pop: 0, Push: 0}, // local effect
	DecrementLocalInteger:  &StackEffect{Pop: 0, Push: 0}, // local effect
	NegateInteger:          &StackEffect{Pop: 1, Push: 1},
	AddInteger:             &StackEffect{Pop: 2, Push: 1},
	SubtractInteger:        &StackEffect{Pop: 2, Push: 1},
	MultiplyInteger:        &StackEffect{Pop: 2, Push: 1},
	GetLocal_0:             &StackEffect{Pop: 0, Push: 1}, // local effect
	GetLocal_1:             &StackEffect{Pop: 0, Push: 1}, // local effect
	GetLocal_2:             &StackEffect{Pop: 0, Push: 1}, // local effect
	GetLocal_3:             &StackEffect{Pop: 0, Push: 1}, // local effect
	SetLocal_0:             &StackEffect{Pop: 1, Push: 0}, // local effect
	SetLocal_1:             &StackEffect{Pop: 1, Push: 0}, // local effect
	SetLocal_2:             &StackEffect{Pop: 1, Push: 0}, // local effect
	SetLocal_3:             &StackEffect{Pop: 1, Push: 0}, // local effect
	Debug:                  &StackEffect{Pop: 0, Push: 0},
	DebugLine:              &StackEffect{Pop: 0, Push: 0},
	DebugFile:              &StackEffect{Pop: 0, Push: 0},
	BreakpointLine:         &StackEffect{Pop: 0, Push: 0}, // unsupported
	Timestamp:              &StackEffect{Pop: 0, Push: 0}, // unsupported
}
