package events

import (
	"github.com/dexter3k/dash/avm2"
)

type Class_EventDispatcher struct {
	avm2.Object_Class
}

type Object_EventDispatcher struct {
	avm2.Object_Object
}

func CreateClass_EventDispatcher(core *avm2.Core) *Class_EventDispatcher {
	class := &Class_EventDispatcher{
		Object_Class: core.CreateNativeClass(
			"EventDispatcher", "flash.events",
			&core.Builtins.Object.Object_Class,
		),
	}

	inst := class.InstanceTraits
	inst.CreateInstance = func() avm2.Any {
		panic("ni")
		// return Object_EventDispatcher_Construct(core, obj, args)
	}

	return class
}
