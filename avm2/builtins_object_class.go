package avm2

// public dynamic class Object
// - $1 public static const length() : int = 1
// - #9 pkgint static func init() : *
//
// - #1 AS3 isPrototypeOf(V : * = undefined) : Boolean
// - #2 AS3 hasOwnProperty(V : * = undefined) : Boolean
// - #3 AS3 propertyIsEnumerable(V : * = undefined) : Boolean
//
// - nonenum proto func hasOwnProperty(V : * = undefined) : Boolean
// - nonenum proto func propertyIsEnumerable(V : * = undefined) : *
// - nonenum proto func setPropertyIsEnumerable(name : String, value : Boolean) : void
// - nonenum proto func isPrototypeOf(V : * = undefined) : Boolean
// - nonenum proto func toString() : String
// - nonenum proto func toLocaleString() : String
// - nonenum proto func valueOf() : *

type Class_Object struct {
	Object_Class
}

type Object_Object struct {
	Traits    *Traits
	Fixed     []Any
	Dynamic   map[string]Any
	Prototype Value

	Constructor  Constructible

	DebugName string
}

func (obj *Object_Object) GetEmbeddedObject() *Object_Object {
	return obj
}

func (class *Class_Object) InitializeNativeTraits(core *Core) {
	// class.Traits.AddTrait(core.NewConstTrait(
	// 	1, "length", "P", "int", "P", Int(1),
	// ))
	class.Traits.AddTrait(core.NewMethodTrait(
		/*9, */"init", "I", func(_ Any, _ []Any) Any {
			return Undefined
		},
	))

	inst := class.InstanceTraits
	inst.RunConstructor = func(Any, []Any) {
		// Nothing here
	}

	inst.AddTrait(core.NewMethodTrait(
		/*1, */"isPrototypeOf", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("ni")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		/*2, */"hasOwnProperty", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("ni")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		/*3, */"propertyIsEnumerable", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("ni")
		},
	))

	proto := class.BasePrototype
	proto.SetProperty(proto, DefaultSpaceSet, "hasOwnProperty", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("ni")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "propertyIsEnumerable", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("ni")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "setPropertyIsEnumerable", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("ni")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "isPrototypeOf", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("ni")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "toString", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("ni")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "toLocaleString", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("ni")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "valueOf", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("ni")
		},
	))
}
