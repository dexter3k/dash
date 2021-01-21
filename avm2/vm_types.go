package avm2

import (
	"fmt"
	"math"
)

type nullType struct{}
type undefinedType struct{}

var Null = &nullType{}
var Undefined = &undefinedType{}

type Int int32
type Uint uint32
type Number float64
type String string
type Boolean bool
type Namespace string

var NotANumber = Number(math.NaN())

// Anything can be Any, but we assume its one of AVM2 builtin primitive types
type Any interface{}

// Todo: there is no fast way to convert Any to Value
// Even when we are sure that Any is a Value
// There are really just two ways of doing that efficiently
// Atoms like in AVM+ and atom-like fat pointers

// Any object (also boxed primitives)
type Value interface {
	Any

	GetTraits() *Traits

	GetSlot(index uint) Any
	SetSlot(index uint, value Any)

	// Cant return Value as we might return Null
	Coerce(source Any) Any

	CallProperty(spaces []string, name string, self Any, args []Any) Any

	HasTrait(spaces []string, name string) bool

	GetProperty(self Any, spaces []string, name string) Any
	SetProperty(self Any, spaces []string, name string, value Any)

	GetIndexedProperty(self Any, spaces []string, index uint) Any
	SetIndexedProperty(self Any, spaces []string, index uint, value Any)

	GetDebugName() string
}

type Objectable interface {
	GetEmbeddedObject() *Object_Object
}

// Classes and Functions implement this
type Constructible interface {
	Value

	InitializeNativeTraits(core *Core)

	Call(self Any, args []Any) Any
	Construct(args []Any) Any

	GetSuper() Constructible
	GetPrototype() Value
	GetInstanceTraits() *Traits
}

type PreferredType int

const (
	HintNoPreferred PreferredType = iota
	HintNumber
)

func IsNotANumber(x Number) bool {
	return math.IsNaN(float64(x))
}

func IsInfinite(x Number) bool {
	return math.IsInf(float64(x), 0)
}

func IsNull(x Any) Boolean {
	_, isNull := x.(*nullType)
	return Boolean(isNull)
}

func IsUndefined(x Any) Boolean {
	_, isUndef := x.(*undefinedType)
	return Boolean(isUndef)
}

func IsNullOrUndefined(x Any) Boolean {
	return IsNull(x) || IsUndefined(x)
}

func IsPrimitive(x Any) Boolean {
	switch x.(type) {
	case *undefinedType, *nullType, Boolean, Number, String, Int, Uint:
		return true
	default:
		return false
	}
}

func ToPrimitive(x Any, hint PreferredType) Any {
	if IsPrimitive(x) {
		return x
	}

	switch val := x.(type) {
	case Value:
		if hint == HintNumber {
			if val.GetProperty(val, []string{"P"}, "valueOf") != Undefined {
				result := val.CallProperty([]string{"P"}, "valueOf", val, nil)
				if IsPrimitive(result) {
					return result
				}
			}
			if val.GetProperty(val, []string{"P"}, "toString") != Undefined {
				result := val.CallProperty([]string{"P"}, "toString", val, nil)
				if IsPrimitive(result) {
					return result
				}
			}
		} else {
			if val.GetProperty(val, []string{"P"}, "toString") != Undefined {
				result := val.CallProperty([]string{"P"}, "toString", val, nil)
				if IsPrimitive(result) {
					return result
				}
			}
			if val.GetProperty(val, []string{"P"}, "valueOf") != Undefined {
				result := val.CallProperty([]string{"P"}, "valueOf", val, nil)
				if IsPrimitive(result) {
					return result
				}
			}
		}
		panic("TypeError: nothing usable found")
	default:
		fmt.Printf("%T\n", x)
		panic(x)
	}
}

func IsObject(x Any) bool {
	switch x.(type) {
	case Value:
		return true
	default:
		return false
	}
}

func IsAnyInteger(x Any) bool {
	switch x.(type) {
	case Uint:
		return true
	case Int:
		return true
	default:
		return false
	}
}

func IntegerCanBeUnsigned(x Any) bool {
	switch value := x.(type) {
	case Number:
		return IntegerCanBeUnsigned(ToInteger(x))
	case Uint:
		return true
	case Int:
		i := int64(value)
		u := uint32(value)
		return i == int64(u)
	default:
		return false
	}
}

func NumberToInteger(x Number) Int {
	if IsNotANumber(x) || IsInfinite(x) || x == 0.0 {
		return Int(0)
	}
	return Int(x) // todo
}

func ToInteger(x Any) Int {
	switch value := x.(type) {
	case Int:
		return value
	case Uint:
		return Int(value)
	case Boolean:
		if value {
			return 1
		} else {
			return 0
		}
	default:
		return NumberToInteger(ToNumber(x))
	}
}

func ToUint(x Any) Uint {
	switch value := x.(type) {
	case Int:
		return Uint(value)
	case Uint:
		return value
	default:
		y := ToNumber(x)
		if IsNotANumber(y) {
			return Uint(0)
		}
		return Uint(y)
	}
}

// ECMA-262 7.1.2 ToBoolean
// Extended for AVM2 compat
func ToBoolean(x Any) Boolean {
	switch value := x.(type) {
	case *undefinedType, *nullType:
		return false
	case Boolean:
		return value
	case Number:
		return !Boolean(value == 0)
	case Int:
		return !Boolean(value == 0)
	case Uint:
		return !Boolean(value == 0)
	case Value:
		return true
	default:
		panic(x)
	}
}

// integer   ->
// undefined -> NaN
func ToNumber(x Any) Number {
	switch value := x.(type) {
	case *undefinedType:
		return Number(math.NaN())
	case *nullType:
		return 0.0
	case Boolean:
		if value == true {
			return 1.0
		} else {
			return 0.0
		}
	case Number:
		return value
	case Int:
		return Number(value)
	case Uint:
		return Number(value)
	case String:
		panic("is a string, not implemented")
	case Value:
		return ToNumber(ToPrimitive(x, HintNumber))
	default:
		panic("Should never happen")
	}
}

func ToString(x Any) String {
	switch val := x.(type) {
	case *undefinedType:
		return "undefined"
	case *nullType:
		return "null"
	case Boolean:
		if val {
			return "true"
		} else {
			return "false"
		}
	case String:
		return val
	case Number:
		return String(fmt.Sprintf("%v", val))
	case Int:
		return String(fmt.Sprintf("%v", val))
	case Uint:
		return String(fmt.Sprintf("%v", val))
	default:
		panic(x)
	}
}

func TypeToString(x Any) String {
	switch x.(type) {
	case *undefinedType:
		return "undefined"
	case *nullType:
		return "null"
	case Boolean:
		return "boolean"
	case Number:
		return "number"
	case Int:
		return "int"
	case Uint:
		return "uint"
	case String:
		return "string"
	case Value:
		return "object"
	default:
		panic(x)
	}
}

// ECMA-262 11.6.1 + (in future) ECMA-357 11.4
func Add(left Any, right Any) Any {
	x := ToPrimitive(left, HintNoPreferred)
	y := ToPrimitive(right, HintNoPreferred)

	if IsAnyInteger(x) && IsAnyInteger(y) {
		return ToUint(x) + ToUint(y)
	}

	_, isXStr := x.(String)
	_, isYStr := y.(String)
	if isXStr || isYStr {
		return ToString(x) + ToString(y)
		// panic(string("Not imeplemted " + xStr + " " + yStr))
	}

	return ToNumber(x) + ToNumber(y)
}

func Subtract(left Any, right Any) Number {
	return ToNumber(left) - ToNumber(right)
}

func Multiply(left Any, right Any) Number {
	return ToNumber(left) * ToNumber(right)
}

func Divide(left Any, right Any) Number {
	return ToNumber(left) / ToNumber(right)
}

func Modulo(left Any, right Any) Number {
	return Number(math.Mod(float64(ToNumber(left)), float64(ToNumber(right))))
}

func CompareEquals(left Any, right Any) Boolean {
	lt := TypeToString(left)
	rt := TypeToString(right)
	if lt != rt {
		if lt == "undefined" && rt == "null" {
			return true
		}
		if rt == "undefined" && lt == "null" {
			return true
		}
		if lt == "number" && rt == "string" {
			panic("yeah")
		}
		if lt == "string" && rt == "number" {
			panic("beah")
		}
		if lt == "boolean" || lt == "int" || lt == "uint" {
			return CompareEquals(ToNumber(left), right)
		}
		if rt == "boolean" || rt == "int" || rt == "uint" {
			return CompareEquals(left, ToNumber(right))
		}
		if (lt == "string" || lt == "number") && rt == "object" {
			return CompareEquals(left, ToPrimitive(right, HintNoPreferred))
		}
		if (rt == "string" || rt == "number") && lt == "object" {
			return CompareEquals(ToPrimitive(left, HintNoPreferred), right)
		}
		return false
	}
	switch lt {
	case "undefined":
		return true
	case "null":
		return true
	case "number":
		return left.(Number) == right.(Number)
	case "int":
		return left.(Int) == right.(Int)
	case "string":
		return left.(String) == right.(String)
	default:
		panic(string(lt))
	}
}

// As per ECMA-262 11.8.5, extended to AS3
// Returns bool or Undefined
func CompareLessThan(left Any, right Any) Any {
	x := ToPrimitive(left, HintNumber)
	y := ToPrimitive(right, HintNumber)
	xStr, isXStr := x.(String)
	yStr, isYStr := y.(String)
	if isXStr && isYStr {
		panic("Is string" + xStr + " " + yStr)
	}
	xn := ToNumber(x)
	if IsNotANumber(xn) {
		return Undefined
	}
	yn := ToNumber(y)
	if IsNotANumber(yn) {
		return Undefined
	}
	return Boolean(xn < yn)
}
