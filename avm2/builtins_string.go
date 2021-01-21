package avm2

import (
// "fmt"
)

type Class_String struct {
	Object_Class
}

// Constructed by Core when boxing String
type Object_String struct {
	Object_Object
	Value String
}

func (class *Class_String) InitializeNativeTraits(core *Core) {
	// ...
}

func (class *Class_String) Call(_ Any, args []Any) Any {
	panic(len(args))
	// if len(args) == 0 {

	// }
	// result := class.Construct(args)
	// fmt.Printf("String(): %#T\n", result)
	// fmt.Println([]Any{result})
	// return result
}
