package avm2

func (core *Core) NewGetterTrait(name, space string, getter func(self Any) Any) Trait {
	return Trait{
		Link: Link(name, space),
		Kind: GetterTrait,
		Get: func(self Any, args []Any) Any {
			return getter(self)
		},
	}
}

func (core *Core) NewSetterTrait(name, space string, setter func(self Any, args []Any) Any) Trait {
	return Trait{
		Link: Link(name, space),
		Kind: SetterTrait,
		Set:  setter,
	}
}

func (core *Core) NewMethodTrait(name, space string, method func(self Any, args []Any) Any) Trait {
	return Trait{
		Link:   Link(name, space),
		Kind:   MethodTrait,
		Method: core.WrapMethod(method),
	}
}

func (core *Core) NewClassTrait(name, space string, class Constructible) Trait {
	return Trait{
		Link:  Link(name, space),
		Kind:  ClassTrait,
		Value: class,
	}
}

func (c *Core) CreateCustomObject() *Object_Object {
	obj := &Object_Object{
		DebugName: "custom",
	}
	obj.InitTraits(newTraits(c.Builtins.Object.GetInstanceTraits()))
	return obj
}

func (c *Core) CreateNativeScript() *Script {
	global := &Global{}
	global.InitTraits(newTraits(c.Builtins.Object.InstanceTraits))
	script := &Script{
		Abc:    nil,
		State:  Executed,
		Script: nil,
		Global: global,
		Scope:  NewScope(nil, global),
	}
	global.Script = script
	return script
}
