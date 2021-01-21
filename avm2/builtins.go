package avm2

var kDefaultInitializer = func() Any {
	panic("ni")
	return Undefined
}

type builtins struct {
	Object        *Class_Object
	Class         *Class_Class
	Function      *Class_Function
	Date          *Class_Date
	Math          *Class_Math
	Int           *Class_Int
	Uint          *Class_Uint
	Number        *Class_Number
	String        *Class_String
	Namespace     *Class_Namespace
	Array         *Class_Array
	Vector        *Class_Vector
	Vector_any    *Class_Vector_any
	Vector_int    *Class_Vector_int
	Vector_number *Class_Vector_number
	Vector_uint   *Class_Vector_uint

	Script *Script
}

func (b *builtins) initialize(core *Core) {
	b.initializeClasses()
	b.initializeClassTraits(core)
	b.populateScript(core)
}

func (b *builtins) initializeClasses() {
	b.Object = &Class_Object{
		Object_Class: Object_Class{
			Object_Object: Object_Object{
				DebugName: "Object",
			},
		},
	}

	b.Class = &Class_Class{
		Object_Class: Object_Class{
			Object_Object: Object_Object{
				DebugName: "Class",
			},
			Super: &b.Object.Object_Class,
		},
	}
	b.Class.Constructor = b.Class
	b.Object.Constructor = b.Class

	// Object instance traits
	b.Object.InstanceTraits = newTraits(nil)
	b.Object.InstanceTraits.CreateInstance = kDefaultInitializer

	// Class instance traits inherit object instance
	b.Class.InstanceTraits = newTraits(b.Object.InstanceTraits)

	// All static traits inherit class instance traits
	b.Object.Traits = newTraits(b.Class.InstanceTraits)
	b.Class.Traits = newTraits(b.Class.InstanceTraits)

	createDefaultProto := func(ctor Constructible) *Object_Object {
		proto := CreateObjectForTraits(b.Object.Traits)
		proto.Constructor = ctor
		return proto
	}

	b.Object.BasePrototype = createDefaultProto(b.Object)
	b.Class.BasePrototype = createDefaultProto(b.Class)

	createClass := func(name string, parent Constructible) Object_Class {
		class := Object_Class{
			Object_Object: Object_Object{
				Constructor: b.Class,
				DebugName:   name,
			},
			Super: parent,
		}
		class.Traits = newTraits(b.Class.InstanceTraits)
		class.InstanceTraits = newTraits(class.Super.GetInstanceTraits())
		return class
	}

	b.Function = &Class_Function{Object_Class: createClass("Function", b.Object)}
	b.Date = &Class_Date{Object_Class: createClass("Date", b.Object)}
	b.Math = &Class_Math{Object_Class: createClass("Math", b.Object)}
	b.Int = &Class_Int{Object_Class: createClass("Int", b.Object)}
	b.Uint = &Class_Uint{Object_Class: createClass("Uint", b.Object)}
	b.Number = &Class_Number{Object_Class: createClass("Number", b.Object)}
	b.String = &Class_String{Object_Class: createClass("String", b.Object)}
	b.Namespace = &Class_Namespace{Object_Class: createClass("Namespace", b.Object)}
	b.Array = &Class_Array{Object_Class: createClass("Array", b.Object)}
	b.Vector = &Class_Vector{Object_Class: createClass("Vector", b.Object)}
	b.Vector_any = &Class_Vector_any{Object_Class: createClass("Vector_any", b.Object)}
	b.Vector_int = &Class_Vector_int{Object_Class: createClass("Vector_int", b.Object)}
	b.Vector_number = &Class_Vector_number{Object_Class: createClass("Vector_number", b.Object)}
	b.Vector_uint = &Class_Vector_uint{Object_Class: createClass("Vector_uint", b.Object)}

	b.Function.BasePrototype = createDefaultProto(b.Function)
	b.Date.BasePrototype = createDefaultProto(b.Date)
	b.Math.BasePrototype = createDefaultProto(b.Math)
	b.Int.BasePrototype = createDefaultProto(b.Int)
	b.Uint.BasePrototype = createDefaultProto(b.Uint)
	b.Number.BasePrototype = createDefaultProto(b.Number)
	b.String.BasePrototype = createDefaultProto(b.String)
	b.Namespace.BasePrototype = createDefaultProto(b.Namespace)
	b.Array.BasePrototype = createDefaultProto(b.Array)
	b.Vector.BasePrototype = createDefaultProto(b.Vector)
	b.Vector_any.BasePrototype = createDefaultProto(b.Vector_any)
	b.Vector_int.BasePrototype = createDefaultProto(b.Vector_int)
	b.Vector_number.BasePrototype = createDefaultProto(b.Vector_number)
	b.Vector_uint.BasePrototype = createDefaultProto(b.Vector_uint)
}

func (b *builtins) initializeClassTraits(core *Core) {
	b.Object.InitializeNativeTraits(core)
	b.Class.InitializeNativeTraits(core)
	b.Function.InitializeNativeTraits(core)
	b.Date.InitializeNativeTraits(core)
	b.Math.InitializeNativeTraits(core)
	b.Int.InitializeNativeTraits(core)
	b.Uint.InitializeNativeTraits(core)
	b.Number.InitializeNativeTraits(core)
	b.String.InitializeNativeTraits(core)
	b.Namespace.InitializeNativeTraits(core)
	b.Array.InitializeNativeTraits(core)
	b.Vector.InitializeNativeTraits(core)
	b.Vector_any.InitializeNativeTraits(core)
	b.Vector_int.InitializeNativeTraits(core)
	b.Vector_number.InitializeNativeTraits(core)
	b.Vector_uint.InitializeNativeTraits(core)
}

func (b *builtins) populateScript(core *Core) {
	b.Script = &Script{
		Abc:    nil,
		State:  Executed,
		Script: nil,
	}
	g := &Global{
		Script: b.Script,
	}
	g.Object_Object.InitTraits(b.Object.Traits)
	g.Script.Global = g
	g.Script.Scope = NewScope(nil, g)
	g.Traits = newTraits(g.Traits)

	addClass := func(name, space string, class Constructible) {
		g.Traits.AddTrait(Trait{
			Link:  Link(name, space),
			Kind:  ClassTrait,
			Value: class,
		})
	}

	// addMethod := func(name, space string, method func(Value, []Any) Any) {
	// 	g.Traits.AddTrait(Trait{
	// 		Link:   Link(name, space),
	// 		Kind:   MethodTrait,
	// 		Method: core.WrapMethod(method),
	// 	})
	// }

	addClass("Object", DefaultNamespace, b.Object)
	addClass("Array", DefaultNamespace, b.Array)
	addClass("int", DefaultNamespace, b.Int)
	addClass("uint", DefaultNamespace, b.Uint)
	addClass("String", DefaultNamespace, b.String)
	addClass("Vector", VectorNamespace, b.Vector)
	// addClass("Date", "P", b.Date)
	// addClass("Math", "P", b.Math)

	// addMethod("print", "P", func(_ Value, args []Any) Any {
	// 	// log.Println("Print:", string(args[0].(String)))
	// 	log.Println("Print:", string(ToString(args[0])))
	// 	return Undefined
	// })
	// addMethod("trace", "P", func(_ Value, args []Any) Any {
	// 	log.Println("Trace:", string(args[0].(String)))
	// 	return Undefined
	// })

	/*Remember to reset slots to defaults after we're done*/
	g.Fixed = g.Traits.CreateSlots()
}
