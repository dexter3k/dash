package avm2

type Class_Vector struct {
	Object_Class
}

func (class *Class_Vector) Call(_ Any, args []Any) Any {
	return class.Construct(args)
}

func (class *Class_Vector) Construct(_ []Any) Any {
	panic("Vector cannot be constructed")
}

func (class *Class_Vector) ApplyType(core *Core, x Any) Constructible {
	switch x.(type) {
	case *nullType:
		panic("Vector<any>")
	case *Class_Int:
		panic("Vector<int>")
	case *Class_Uint:
		return core.Builtins.Vector_uint
	case *Class_Number:
		panic("Vector<Number>")
	default:
		panic(x)
	}
}

func (class *Class_Vector) InitializeNativeTraits(core *Core) {
	class.InstanceTraits.Initializer = func(args []Any) Any {
		panic("No really ever possible to happen cuz we override Construct")
	}
}
