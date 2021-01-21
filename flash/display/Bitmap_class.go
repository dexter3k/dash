package display

import (
	"github.com/dexter3k/dash/avm2"
)

type Class_Bitmap struct {
	avm2.Object_Class
}

type Object_Bitmap struct {
	avm2.Object_Object
}

func CreateClass_Bitmap(core *avm2.Core, super avm2.Constructible) *Class_Bitmap {
	class := &Class_Bitmap{
		Object_Class: core.CreateNativeClass("Bitmap", "flash.display", super),
	}

	class.InstanceTraits.Initializer = func(args []avm2.Any) avm2.Any {
		panic("ni")
	}

	return class
}
