package display

import (
	"github.com/dexter3k/dash/avm2"
)

type Class_Sprite struct {
	avm2.Object_Class
}

type Object_Sprite struct {
	avm2.Object_Object
}

func CreateClass_Sprite(core *avm2.Core, super avm2.Constructible) *Class_Sprite {
	class := &Class_Sprite{
		Object_Class: core.CreateNativeClass("Sprite", "flash.display", super),
	}

	class.InstanceTraits.Initializer = func(args []avm2.Any) avm2.Any {
		panic("ni")
	}

	return class
}
