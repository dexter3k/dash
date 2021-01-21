package avm2

// final dynamic class Vector$uint
// - public get length() : uint
// - public set length(value : uint) : *
// - public get fixed() : Boolean
// - public set fixed(value : Boolean) : *
// - AS3 func toString() : String
// - AS3 func toLocaleString() : String
// - AS3 func join(separator : String = ",") : String
// - AS3 func every(checked : Function, self : Object = null) : Boolean
// - AS3 func forEach(eacher : Function, self : Object = null) : void
// - AS3 func map(mapper : Function, self : Object = null) : *
// - AS3 func push(... rest) : uint
// - AS3 func some(checker : *, self : Object = null) : Boolean
// - AS3 func unshift(... rest) : uint
// - AS3 func concat(... items) : Vector$uint
// - AS3 func filter(checker : Function, self : Object = null) : Vector$uint
// - AS3 func pop() : uint
// - AS3 func reverse() : Vector$uint
// - AS3 func shift() : uint
// - AS3 func slice(start : Number = 0, end : Number = 2147483647) : Array
// - AS3 func sort(comparefn : *) : Vector$uint
// - AS3 func splice(start : Number, delCount : Number, ...items) : Vector$uint
// - AS3 func insertAt(index : int, arg2 : uint) : void
// - AS3 func removeAt(index : uint) : *
// - AS3 func indexOf(value : uint, from : Number = 0) : Number
// - AS3 func lastIndexOf(value : uint, from : Number = 2147483647) : int
//
// - nonenum proto func push(... rest) : uint
// - nonenum proto func toString() : String
// - nonenum proto func toLocaleString() : String
// - nonenum proto func join(separator : String = ",") : String
// - nonenum proto func concat(... items) : Vector$uint
// - nonenum proto func every(checked : Function, self : Object = null) : Boolean
// - nonenum proto func filter(checker : Function, self : Object = null) : Vector$uint
// - nonenum proto func forEach(eacher : Function, self : Object = null) : void
// - nonenum proto func indexOf(value : uint, from : Number = 0) : Number
// - nonenum proto func lastIndexOf(value : uint, from : Number = 2147483647) : int
// - nonenum proto func map(mapper : Function, self : Object = null) : *
// - nonenum proto func pop() : uint
// - nonenum proto func reverse() : Vector$uint
// - nonenum proto func shift() : uint
// - nonenum proto func slice(start : Number = 0, end : Number = 2147483647) : Array
// - nonenum proto func some(checker : *, self : Object = null) : Boolean
// - nonenum proto func sort(comparefn : *) : Vector$uint
// - nonenum proto func splice(start : Number, delCount : Number, ...items) : Vector$uint
// - nonenum proto func unshift(... rest) : uint
type Class_Vector_uint struct {
	Object_Class
}

type Object_Vector_uint struct {
	Object_Object

	isFixed bool
	data    []Uint
}

func Class_Vector_uint_Construct(core *Core, args []Any) *Object_Vector_uint {
	var length uint
	if len(args) >= 1 {
		length = uint(ToUint(args[0]))
	}

	var fixed bool
	if len(args) >= 2 {
		fixed = bool(ToBoolean(args[1]))
	}

	this := &Object_Vector_uint{
		isFixed:       fixed,
		data:          make([]Uint, length, length),
	}
	this.InitTraits(core.Builtins.Vector_uint.InstanceTraits)
	this.DebugName = "Vector$uint"

	return this
}

func (class *Class_Vector_uint) InitializeNativeTraits(core *Core) {
	inst := class.InstanceTraits
	inst.Initializer = func(args []Any) Any {
		return Class_Vector_uint_Construct(core, args)
	}
	inst.AddTrait(core.NewGetterTrait(
		"length", "P", func(self Any) Any {
			this := self.(*Object_Vector_uint)
			return Uint(len(this.data))
		},
	))
	inst.AddTrait(core.NewSetterTrait(
		"length", "P", func(self Any, args []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"toString", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"toLocaleString", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"join", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"every", AS3Namespace, func(_ Any, _ []Any) Any {
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
		"push", AS3Namespace, func(self Any, args []Any) Any {
			this := self.(*Object_Vector_uint)
			for i := 0; i < len(args); i++ {
				this.data = append(this.data, ToUint(args[i]))
			}
			return Uint(len(this.data))
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"some", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"unshift", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"concat", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"filter", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"pop", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"reverse", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"shift", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"slice", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"sort", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"splice", AS3Namespace, func(_ Any, _ []Any) Any {
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
		"indexOf", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"lastIndexOf", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))

	proto := class.BasePrototype
	proto.SetProperty(proto, DefaultSpaceSet, "push", core.WrapMethod(
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
	proto.SetProperty(proto, DefaultSpaceSet, "join", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "concat", core.WrapMethod(
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
	proto.SetProperty(proto, DefaultSpaceSet, "map", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "pop", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "reverse", core.WrapMethod(
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
	proto.SetProperty(proto, DefaultSpaceSet, "some", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "sort", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "splice", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "unshift", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
}
