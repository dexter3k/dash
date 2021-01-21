package avm2

func (this *Object_Class) Construct(args []Any) Any {
	if this.InstanceTraits == nil {
		panic("Construct: no instance traits")
	}
	if this.InstanceTraits.Initializer == nil {
		panic("Construct: no constructor")
	}

	return this.InstanceTraits.Initializer(args)
}

func (this *Object_Class) Call(self Any, args []Any) Any {
	panic("ni")
}

func (this *Object_Class) GetSuper() Constructible {
	return this.Super
}

func (this *Object_Class) GetPrototype() Value {
	return this.BasePrototype
}

func (this *Object_Class) GetInstanceTraits() *Traits {
	return this.InstanceTraits
}

func (this *Object_Class) GetDebugName() string {
	return "[class " + this.DebugName + "]"
}

func (this *Object_Class) InitializeNativeTraits(_ *Core) {
	panic("Please override InitializeNativeTraits on native class")
}
