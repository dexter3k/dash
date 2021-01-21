package avm2

// public final class uint
// - public static get MIN_VALUE() : uint
// - public static get MAX_VALUE() : uint
// - public static get length() : int
//
// - AS3 func toString(radix : * = 10) : String
// - AS3 func valueOf() : uint
// - AS3 func toExponential(p : * = 0) : String
// - AS3 func toPrecision(p : * = 0) : String
// - AS3 func toFixed(p : * = 0) : String
//
// - nonenum proto func toLocaleString(radix : * = 10) : String
// - nonenum proto func toString(radix : * = 10) : String
// - nonenum proto func valueOf() : uint
// - nonenum proto func toExponential(p : * = 0) : String
// - nonenum proto func toPrecision(p : * = 0) : String
// - nonenum proto func toFixed(p : * = 0) : String
type Class_Uint struct {
	Object_Class
}

func (class *Class_Uint) InitializeNativeTraits(core *Core) {
	class.Traits.AddTrait(core.NewGetterTrait(
		"MIN_VALUE", "P", func(_ Any) Any {
			return Uint(0)
		},
	))
	class.Traits.AddTrait(core.NewGetterTrait(
		"MAX_VALUE", "P", func(_ Any) Any {
			return Uint(4294967295)
		},
	))
	class.Traits.AddTrait(core.NewGetterTrait(
		"length", "P", func(_ Any) Any {
			return Int(1)
		},
	))

	inst := class.InstanceTraits
	inst.CreateInstance = func() Any {
		panic("Not implemented")
	}
	inst.AddTrait(core.NewMethodTrait(
		"toString", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"valueOf", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"toExponential", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"toPrecision", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	inst.AddTrait(core.NewMethodTrait(
		"toFixed", AS3Namespace, func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))

	proto := class.BasePrototype
	proto.SetProperty(proto, DefaultSpaceSet, "toString", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "valueOf", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "toExponential", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "toPrecision", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
	proto.SetProperty(proto, DefaultSpaceSet, "toFixed", core.WrapMethod(
		func(_ Any, _ []Any) Any {
			panic("Not implemented")
		},
	))
}
