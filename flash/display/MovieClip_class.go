package display

import (
	"github.com/dexter3k/dash/avm2"
)

type Class_MovieClip struct {
	avm2.Object_Class
}

type Object_MovieClip struct {
	avm2.Object_Object
}

func CreateClass_MovieClip(core *avm2.Core, super avm2.Constructible) *Class_MovieClip {
	class := &Class_MovieClip{
		Object_Class: core.CreateNativeClass("MovieClip", "flash.display", super),
	}

	class.InstanceTraits.Initializer = func(args []avm2.Any) avm2.Any {
		panic("ni")
	}

	class.InstanceTraits.RunConstructor = func(self avm2.Any, args []avm2.Any) {
		// panic("ni")
	}

	// public stop() : void
	class.InstanceTraits.AddTrait(core.NewMethodTrait(
		"stop", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			return avm2.Undefined
		},
	))

	return class
}
