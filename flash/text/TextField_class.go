package text

import (
	"github.com/dexter3k/dash/avm2"
)

type Class_TextField struct {
	avm2.Object_Class
}

type Object_TextField struct {
	avm2.Object_Object
}

func CreateClass_TextField(core *avm2.Core, super avm2.Constructible) *Class_TextField {
	class := &Class_TextField{
		Object_Class: core.CreateNativeClass("TextField", "flash.text", super),
	}

	class.InstanceTraits.Initializer = func(args []avm2.Any) avm2.Any {
		panic("ni")
	}

	return class
}
