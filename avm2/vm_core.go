package avm2

import (
	"fmt"

	"github.com/dexter3k/dash/abc"
)

const (
	AS3Namespace     = "Phttp://adobe.com/AS3/2006/builtin"
	FlexNamespace    = "Phttp://www.adobe.com/2006/flex/mx/internal"
	VectorNamespace  = "P__AS3__.vec"
	DefaultNamespace = "P"
)

var (
	DefaultSpaceSet = []string{DefaultNamespace}
)

type Core struct {
	Root     *Domain
	App      *Domain
	Builtins builtins

	debugger Debugger

	InsideExecute bool
	Exited        bool
	ExitError     error
}

func NewCore() *Core {
	c := &Core{}
	c.Root = NewDomain(c, nil)
	c.App = NewDomain(c, c.Root)
	c.initialize()
	// c.debugger = NewBasicDebugger(c)
	return c
}

func (c *Core) GetTraitsFor(x Any) *Traits {
	// If x is a primitive, get primitive class traits
	// else get objects own traits
	switch obj := x.(type) {
	case *Object_Object:
		return obj.Traits
	case *Global:
		return obj.Traits
	case Value:
		return obj.GetTraits()
	default:
		panic(x)
	}
}

func (c *Core) GetPrototypeFor(x Any) Value {
	panic("ni")
}

func (c *Core) WrapDefaultValue(scope *Scope, value interface{}, datatype *abc.Multiname) Any {
	// 1) find datatype
	if datatype.IsRuntime() {
		panic("Nope")
	}
	// g := scope.Parent.FindScopeProperty(datatype, true)
	fmt.Println(datatype.GetMangledName())
	g := scope.Global.Object.(*Global).Domain.FindProperty(datatype)
	if g == nil {
		panic(fmt.Errorf("WrapDefaultValue: Prop not found: %s", datatype.GetMangledName()))
	}
	dt := g.GetProperty(g, datatype.Spaces, datatype.Name)

	var v Any
	// 2) convert value to avm2 type
	switch x := value.(type) {
	case int32:
		v = Int(x)
	case string:
		v = String(x)
	default:
		// todo
		panic([]Any{value}[0].(Value))
	}
	// 3) coerce
	return dt.(Value).Coerce(v)
}

func (c *Core) BoxIfPrimitive(any Any) Value {
	switch value := any.(type) {
	case *nullType:
		panic("ni")
	case *undefinedType:
		panic("ni")
	case Boolean:
		panic("ni")
	case Int:
		panic("ni")
	case Uint:
		panic("ni")
	case Number:
		panic("ni")
	case String:
		return c.BoxString(value)
	case Namespace:
		panic("ni")
	default:
		// Todo: assertE2I is slow.
		return any.(Value)
	}
}

func (c *Core) BoxString(value String) Value {
	panic("ni")
}

func (c *Core) WrapMethod(m func(self Any, args []Any) Any) *Object_Function {
	f := &Object_Function{
		Execute: m,
	}
	f.InitTraits(c.Builtins.Function.Traits)
	f.DebugName = "Method"
	return f
}

func (c *Core) CreateArray(values []Any) Any {
	return ArrayFromSlice(c, values)
}

func (c *Core) CreateFunction(pool *abc.File, m *abc.Method, scope *Scope) Value {
	panic("ni")
	f := &Object_Function{
		Object_Class: Object_Class{
			// Object_Object: *c.Builtins.Function.instantiate(),
			// Prototype:     c.Builtins.Object.instantiate(),
		},
		Execute: func(self Any, args []Any) Any {
			return Execute(c, pool, m, scope, self, args)
		},
	}
	f.DebugName = "Function"
	f.Prototype.(*Object_Object).Constructor = f
	return f
}

func (self *Core) CreatePrototypeForClass(class Constructible) *Object_Object {
	proto := CreateObjectForTraits(self.Builtins.Object.Traits)
	proto.Constructor = class
	return proto
}

func (self *Core) CreateNativeClass(name, space string, super Constructible) Object_Class {
	link := Link(name, space)
	class := Object_Class{
		Object_Object: Object_Object{
			DebugName: link.Link(),
		},
		Super:          super,
		InstanceTraits: newTraits(super.GetInstanceTraits()),
	}
	class.InitTraits(self.Builtins.Class.Traits)
	return class
}

// A tricky one
func (self *Core) CreateInterface(pool *abc.File, c *abc.Class) Value {
	iface := &Object_Class{}
	iface.DebugName = "iface " + c.Name.Name
	iface.InitTraits(self.Builtins.Class.Traits)
	iface.InstanceTraits = ApplyAbcTraits(c.Traits, nil, self, pool, nil, "", "iface ")
	return iface
}

func (core *Core) CreateClass(pool *abc.File, c *abc.Class, parentScope *Scope, base Constructible) Value {
	class := &Object_Class{
		// Object_Object: *core.Builtins.Class.instantiate(),
		Super: base,
	}
	class.DebugName = c.Name.Name
	proto := core.CreatePrototypeForClass(class)
	proto.Prototype = base.GetPrototype()
	class.BasePrototype = proto

	// Apply static and instance traits on the class
	staticTraitPrefix := ""
	if core.debugger != nil {
		staticTraitPrefix = c.Name.GetMangledName() + " static "
	}
	scope := NewScope(parentScope, class)
	class.InitTraits(ApplyAbcTraits(
		c.StaticTraits, core.Builtins.Class.Traits,
		core, pool, scope, "", staticTraitPrefix,
	))

	instTraitPrefix := ""
	if core.debugger != nil {
		instTraitPrefix = c.Name.GetMangledName() + " "
	}
	class.InstanceTraits = ApplyAbcTraits(c.Traits, base.GetInstanceTraits(), core, pool, scope, c.ProtectedSpace, instTraitPrefix)
	class.InstanceTraits.RunConstructor = func(self Any, args []Any) {
		Execute(core, pool, c.Constructor, scope, self, args)
	}
	class.InstanceTraits.Initializer = func(args []Any) Any {
		// TODO: objects inheriting native objects?
		self := &Object_Object{}
		self.InitTraits(class.InstanceTraits)
		self.Prototype = class.BasePrototype
		class.InstanceTraits.RunConstructor(self, args)
		return self
	}

	// Now apply all interfaces this class implements
	for _, mn := range c.Implements {
		fmt.Println("Implementing interface", mn.GetMangledName())
		// resolve interface class
		result := scope.FindScopeProperty(mn, true)
		iface := result.GetProperty(result, mn.Spaces, mn.Name).(*Object_Class)
		class.InstanceTraits.ApplyInterface(iface.InstanceTraits)
	}

	if core.debugger != nil {
		core.debugger.RegisterMethodName(c.StaticConstructor, c.Name.GetMangledName()+"(static)")
		core.debugger.RegisterMethodName(c.Constructor, c.Name.GetMangledName()+"(ctor)")
	}

	// Call static constructor
	Execute(core, pool, c.StaticConstructor, scope, class, nil)

	return class
}

func (core *Core) CreateGlobalOnScript(script *Script, d *Domain) {
	if script.Global != nil || script.State != NotExecuted {
		panic("Creating global for already existing script")
	}

	global := &Global{
		Script: script,
		Domain: d,
	}
	script.Global = global
	script.Scope = NewScope(nil, script.Global)
	global.DebugName = "global"
	global.InitTraits(ApplyAbcTraits(
		script.Script.Traits, core.Builtins.Object.InstanceTraits,
		core, script.Abc, script.Scope, "", "global ",
	))
	

	// Now add init function for that trait
	init := script.Script.Script
	if init.Body == nil {
		panic("Body is missing from script init method")
	}
	global.Traits.Initializer = func(_ []Any) Any {
		Execute(core, script.Abc, init, script.Scope, script.Global, nil)
		return Undefined
	}
	if core.debugger != nil {
		core.debugger.RegisterMethodName(init, "global init")
	}
}

func (c *Core) initialize() {
	c.Builtins.initialize(c)
	c.Root.loadNativeScript(c.Builtins.Script)
}
