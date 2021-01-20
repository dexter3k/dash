package avm2

// naming convention is as follows:
// Class_*name*  -- static implementation of a class
// Object_*name* -- instance of a class

// Anything inside avm2 inherits and object
// in class Object, the class itself is an object which is an instance of a class Class
// class Class is an instance of itself

type Class_Object struct {}

type Object_Object struct {}
