package display

import (
	"github.com/dexter3k/dash/avm2"
)

type Class_DisplayObject struct {
	avm2.Object_Class
}

type Object_DisplayObject struct {
	avm2.Object_Object
}

func CreateClass_DisplayObject(core *avm2.Core, super avm2.Constructible) *Class_DisplayObject {
	class := &Class_DisplayObject{
		Object_Class: core.CreateNativeClass("DisplayObject", "flash.display", super),
	}

	class.InstanceTraits.CreateInstance = func() avm2.Any {
		panic("ni")
	}

	// public get visible() : Boolean
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"visible", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get visible: not implemented")
		},
	))

	// public get rotation() : Number
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"rotation", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get rotation: not implemented")
		},
	))

	// public localToGlobal(param1:Point) : Point
	class.InstanceTraits.AddTrait(core.NewMethodTrait(
		"localToGlobal", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("method localToGlobal: not implemented")
		},
	))
	// public get name() : String
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"name", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get name: not implemented")
		},
	))

	// public set width(param1:Number) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"width", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set width: not implemented")
		},
	))

	// public get blendMode() : String
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"blendMode", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get blendMode: not implemented")
		},
	))

	// public get scale9Grid() : Rectangle
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"scale9Grid", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get scale9Grid: not implemented")
		},
	))

	// public set name(param1:String) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"name", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set name: not implemented")
		},
	))

	// public set scaleX(param1:Number) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"scaleX", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set scaleX: not implemented")
		},
	))

	// public set scaleY(param1:Number) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"scaleY", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set scaleY: not implemented")
		},
	))

	// public get accessibilityProperties() : AccessibilityProperties
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"accessibilityProperties", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get accessibilityProperties: not implemented")
		},
	))

	// public set scrollRect(param1:Rectangle) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"scrollRect", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set scrollRect: not implemented")
		},
	))

	// public get cacheAsBitmap() : Boolean
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"cacheAsBitmap", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get cacheAsBitmap: not implemented")
		},
	))

	// public globalToLocal(param1:Point) : Point
	class.InstanceTraits.AddTrait(core.NewMethodTrait(
		"globalToLocal", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("method globalToLocal: not implemented")
		},
	))
	// public get height() : Number
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"height", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get height: not implemented")
		},
	))

	// public set blendMode(param1:String) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"blendMode", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set blendMode: not implemented")
		},
	))

	// public get parent() : DisplayObjectContainer
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"parent", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get parent: not implemented")
		},
	))

	// public getBounds(param1:DisplayObject) : Rectangle
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"getBounds", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get getBounds: not implemented")
		},
	))

	// public get opaqueBackground() : Object
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"opaqueBackground", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get opaqueBackground: not implemented")
		},
	))

	// public set scale9Grid(param1:Rectangle) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"scale9Grid", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set scale9Grid: not implemented")
		},
	))

	// public set alpha(param1:Number) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"alpha", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set alpha: not implemented")
		},
	))

	// public set accessibilityProperties(param1:AccessibilityProperties) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"accessibilityProperties", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set accessibilityProperties: not implemented")
		},
	))

	// public get width() : Number
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"width", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get width: not implemented")
		},
	))

	// public hitTestPoint(param1:Number, param2:Number, param3:Boolean = false) : Boolean
	class.InstanceTraits.AddTrait(core.NewMethodTrait(
		"hitTestPoint", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("method hitTestPoint: not implemented")
		},
	))
	// public set cacheAsBitmap(param1:Boolean) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"cacheAsBitmap", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set cacheAsBitmap: not implemented")
		},
	))

	// public get scaleX() : Number
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"scaleX", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get scaleX: not implemented")
		},
	))

	// public get scaleY() : Number
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"scaleY", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get scaleY: not implemented")
		},
	))

	// public get scrollRect() : Rectangle
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"scrollRect", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get scrollRect: not implemented")
		},
	))

	// public get mouseX() : Number
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"mouseX", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get mouseX: not implemented")
		},
	))

	// public get mouseY() : Number
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"mouseY", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get mouseY: not implemented")
		},
	))

	// public set height(param1:Number) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"height", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set height: not implemented")
		},
	))

	// public set mask(param1:DisplayObject) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"mask", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set mask: not implemented")
		},
	))

	// public getRect(param1:DisplayObject) : Rectangle
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"getRect", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get getRect: not implemented")
		},
	))

	// public get alpha() : Number
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"alpha", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get alpha: not implemented")
		},
	))

	// public set transform(param1:Transform) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"transform", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set transform: not implemented")
		},
	))

	// public get loaderInfo() : LoaderInfo
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"loaderInfo", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get loaderInfo: not implemented")
		},
	))

	// public get root() : DisplayObject
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"root", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get root: not implemented")
		},
	))

	// public hitTestObject(param1:DisplayObject) : Boolean
	class.InstanceTraits.AddTrait(core.NewMethodTrait(
		"hitTestObject", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("method hitTestObject: not implemented")
		},
	))
	// public set opaqueBackground(param1:Object) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"opaqueBackground", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set opaqueBackground: not implemented")
		},
	))

	// public set visible(param1:Boolean) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"visible", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set visible: not implemented")
		},
	))

	// public get mask() : DisplayObject
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"mask", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get mask: not implemented")
		},
	))

	// public set x(param1:Number) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"x", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set x: not implemented")
		},
	))

	// public set y(param1:Number) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"y", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set y: not implemented")
		},
	))

	// public get transform() : Transform
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"transform", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get transform: not implemented")
		},
	))

	// public set filters(param1:Array) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"filters", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set filters: not implemented")
		},
	))

	// public get x() : Number
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"x", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get x: not implemented")
		},
	))

	// public get y() : Number
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"y", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get y: not implemented")
		},
	))

	// public get filters() : Array
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"filters", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get filters: not implemented")
		},
	))

	// public set rotation(param1:Number) : void
	class.InstanceTraits.AddTrait(core.NewSetterTrait(
		"rotation", avm2.DefaultNamespace,
		func(self avm2.Any, args []avm2.Any) avm2.Any {
			panic("set rotation: not implemented")
		},
	))

	// public get stage() : Stage
	class.InstanceTraits.AddTrait(core.NewGetterTrait(
		"stage", avm2.DefaultNamespace,
		func(self avm2.Any) avm2.Any {
			panic("get stage: not implemented")
		},
	))


	return class
}
