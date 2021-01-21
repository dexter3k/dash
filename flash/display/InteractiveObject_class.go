package display

import (
	"github.com/dexter3k/dash/avm2"
)

type Class_InteractiveObject struct {
	avm2.Object_Class
}

type Object_InteractiveObject struct {
	avm2.Object_Object
}

func CreateClass_InteractiveObject(core *avm2.Core, super avm2.Constructible) *Class_InteractiveObject {
	class := &Class_InteractiveObject{
		Object_Class: core.CreateNativeClass("InteractiveObject", "flash.display", super),
	}

	class.InstanceTraits.CreateInstance = func() avm2.Any {
		panic("ni")
	}

	return class
}
