package avm2

// public dynamic class Class
// - public static const length : int = 1 // slot 1
//
// - public final get prototype() : *

type Class_Class struct {
	Object_Class
}

type Object_Class struct {
	Object_Object

	Super          Constructible
	BasePrototype  Value
	InstanceTraits *Traits
}

func (class *Class_Class) InitializeNativeTraits(core *Core) {
	class.Traits.AddTrait(core.NewGetterTrait(
		"length", "P",
		func(_ Any) Any {
			return Int(1)
		},
	))

	inst := class.InstanceTraits
	inst.AddTrait(core.NewGetterTrait(
		"prototype", "P",
		func(self Any) Any {
			this := self.(Constructible)
			return this.GetPrototype()
		},
	))
}
