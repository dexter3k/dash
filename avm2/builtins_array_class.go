package avm2

// public dynamic class Array
// - public static get CASEINSENSITIVE() : uint
// - public static get DESCENDING() : uint
// - public static get UNIQUESORT() : uint
// - public static get RETURNINDEXEDARRAY() : uint
// - public static get NUMBERIC() : uint
// - public static get length() : int
//
// - public get length() : uint
// - public set length(value : uint) : *
// - AS3 func insertAt(index : int, arg2 : *) : void
// - AS3 func removeAt(index : int) : *
// - AS3 func pop() : *
// - AS3 func push(... rest) : uint
// - AS3 func shift() : *
// - AS3 func unshift(... rest) : Array
// - AS3 func reverse() : Array
// - AS3 func concat(... args) : Array
// - AS3 func slice(A : * = 0, B : * = 4294967295) : Array
// - AS3 func join(sep : * = undefined) : String
// - AS3 func splice(... args) : *
// - AS3 func sort(... args) : *
// - AS3 func sortOn(names : *, options: * = 0, ...ignored) : *
// - AS3 func indexOf(element : *, fromIndex : * = 0) : int
// - AS3 func lastIndexOf(element : *, fromIndex : * = 2147483647) : int
// - AS3 func every(cb : Function, this : * = null) : Boolean
// - AS3 func filter(cb : Function, this : * = null) : Array
// - AS3 func forEach(cb : Function, this : * = null) : void
// - AS3 func map(cb : Function, this : * = null) : Array
// - AS3 func some(cb : Function, this : * = null) : Boolean
//
// - nonenum proto func join(sep : * = undefined) : String
// - nonenum proto func pop() : *
// - nonenum proto func toString() : String
// - nonenum proto func toLocaleString() : String
// - nonenum proto func push(... rest) : uint
// - nonenum proto func reverse() : Array
// - nonenum proto func concat(... args) : Array
// - nonenum proto func shift() : *
// - nonenum proto func slice(A : * = 0, B : * = 4294967295) : Array
// - nonenum proto func unshift(... rest) : Array
// - nonenum proto func splice(... args) : *
// - nonenum proto func sort(... args) : *
// - nonenum proto func sortOn(names : *, options: * = 0, ...ignored) : *
// - nonenum proto func indexOf(element : *, fromIndex : * = 0) : int
// - nonenum proto func lastIndexOf(element : *, fromIndex : * = 2147483647) : int
// - nonenum proto func every(cb : Function, this : * = null) : Boolean
// - nonenum proto func filter(cb : Function, this : * = null) : Array
// - nonenum proto func forEach(cb : Function, this : * = null) : void
// - nonenum proto func map(cb : Function, this : * = null) : Array
// - nonenum proto func some(cb : Function, this : * = null) : Boolean
type Class_Array struct {
	Object_Class
}

type Object_Array struct {
	Object_Object

	values []Any
}

func ArrayFromSlice(core *Core, values []Any) *Object_Array {
	array := &Object_Array{
		values:        make([]Any, len(values)),
	}
	array.Object_Object.InitTraits(core.Builtins.Array.InstanceTraits)
	copy(array.values, values)
	array.DebugName = "Array"
	return array
}

func Object_Array_Construct(core *Core, args []Any) *Object_Array {
	this := &Object_Array{}
	this.InitTraits(core.Builtins.Array.InstanceTraits)
	this.DebugName = "Array"

	if len(args) == 1 {
		length := ToUint(args[0])
		this.values = make([]Any, length)
	} else if len(args) != 0 {
		this.values = make([]Any, len(args))
		copy(this.values, args)
	} else {
		this.values = make([]Any, 0)
	}

	return this
}

func (class *Class_Array) Call(_ Any, args []Any) Any {
	return class.Construct(args)
}

func (class *Class_Array) InitializeNativeTraits(core *Core) {
	class.Traits.AddTrait(core.NewGetterTrait(
		"CASEINSENSITIVE", "P", func(_ Any) Any {
			return Uint(1)
		},
	))
	class.Traits.AddTrait(core.NewGetterTrait(
		"DESCENDING", "P", func(_ Any) Any {
			return Uint(2)
		},
	))
	class.Traits.AddTrait(core.NewGetterTrait(
		"UNIQUESORT", "P", func(_ Any) Any {
			return Uint(4)
		},
	))
	class.Traits.AddTrait(core.NewGetterTrait(
		"RETURNINDEXEDARRAY", "P", func(_ Any) Any {
			return Uint(8)
		},
	))
	class.Traits.AddTrait(core.NewGetterTrait(
		"NUMBERIC", "P", func(_ Any) Any {
			return Uint(16)
		},
	))
	class.Traits.AddTrait(core.NewGetterTrait(
		"length", "P", func(_ Any) Any {
			return Int(1)
		},
	))

	inst := class.InstanceTraits
	inst.Initializer = func(args []Any) Any {
		return Object_Array_Construct(core, args)
	}
	inst.AddTrait(core.NewGetterTrait(
		"length", "P", func(self Any) Any {
			this := self.(*Object_Array)
			return Uint(len(this.values))
		},
	))
	inst.AddTrait(core.NewSetterTrait(
		"length", "P", func(self Any, args []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"insertAt", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"removeAt", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"pop", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"push", AS3Namespace, func(self Any, args []Any) Any {
			this := self.(*Object_Array)
			return this.avm2_as3_func_push(args)
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"shift", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"unshift", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"reverse", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"concat", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"slice", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"join", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"splice", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"sort", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"sortOn", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"indexOf", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"lastIndexOf", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"every", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"filter", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"forEach", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"map", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"some", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))

	proto := class.BasePrototype
	proto.SetProperty(proto, DefaultSpaceSet, "join", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "pop", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "toString", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "toLocaleString", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "push", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "reverse", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "concat", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "shift", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "slice", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "unshift", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "splice", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "sort", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "sortOn", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "indexOf", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "lastIndexOf", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "every", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "filter", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "forEach", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "map", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "some", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
}
