package avm2

// class Function is simply a normal class
type Class_Function struct {
	Object_Class
}

func (class *Class_Function) InitializeNativeTraits(core *Core) {
	// ...
}

// instance of class Function is based on instance of a class
type Object_Function struct {
	Object_Class

	Execute func(self Any, args []Any) Any
}

func (f *Object_Function) GetDebugName() string {
	return "[function " + f.DebugName + "]"
}

func (f *Object_Function) Construct(args []Any) Any {
	o := &Object_Object{
		Constructor: f,
		Traits:      newTraits(nil),
		Dynamic:     map[string]Any{},
		DebugName:   "constructed_" + f.DebugName,
	}
	result := f.Execute(o, args)
	if _, isNull := result.(*nullType); isNull {
		return Null
	}
	if _, isObject := result.(Value); isObject {
		return result
	}
	return o
}

func (f *Object_Function) Call(self Any, args []Any) Any {
	return f.Execute(self, args)
}
