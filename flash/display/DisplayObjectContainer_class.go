package display

import (
	"github.com/dexter3k/dash/avm2"
)

type Class_DisplayObjectContainer struct {
	avm2.Object_Class
}

type Object_DisplayObjectContainer struct {
	avm2.Object_Object
}

func CreateClass_DisplayObjectContainer(core *avm2.Core, super avm2.Constructible) *Class_DisplayObjectContainer {
	class := &Class_DisplayObjectContainer{
		Object_Class: core.CreateNativeClass("DisplayObjectContainer", "flash.display", super),
	}

	class.InstanceTraits.CreateInstance = func() avm2.Any {
		panic("ni")
	}

	return class
}
