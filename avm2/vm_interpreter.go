package avm2

import (
	"fmt"

	"github.com/dexter3k/dash/abc"
	"github.com/dexter3k/dash/abc/opcodes"
)

type ExecutionFrame struct {
	// Entry conditions
	Core        *Core
	Pool        *abc.File
	Method      *abc.Method
	ParentScope *Scope
	Self        Any
	Args        []Any

	// Current state
	PC     int
	Code   []byte
	Stack  []Any
	Locals []Any
	Scope  *Scope
	Return Any
}

func Execute(core *Core, pool *abc.File, method *abc.Method, parentScope *Scope, self Any, args []Any) Any {
	if core == nil || pool == nil || method == nil || parentScope == nil || self == nil {
		panic("Bad execution frame")
	}

	if method.Body == nil || method.Body.Code == nil {
		// TODO: We might want to just return from this function?
		panic("Method has no body, but passed to Execute")
	}

	frame := ExecutionFrame{
		Core:        core,
		Pool:        pool,
		Method:      method,
		ParentScope: parentScope,
		Self:        self,
		Args:        args,
		PC:          0,
		Code:        method.Body.Code,
		Stack:       make([]Any, 0, method.Body.StackLimit),
		Locals:      make([]Any, method.Body.LocalLimit),
		Scope:       parentScope,
		Return:      nil,
	}
	frame.Locals[0] = self
	for i := 1; i < method.Body.LocalLimit; i++ {
		if i <= len(args) {
			frame.Locals[i] = args[i-1]
		} else {
			frame.Locals[i] = Undefined
		}
	}

	if core.debugger != nil {
		core.debugger.EnterMethod(&frame)
		defer core.debugger.LeaveMethod(&frame)
	}

	return frame.execute()
}

func (f *ExecutionFrame) execute() Any {
	for f.Return == nil {
		op := f.u8()
		cb := switchTable[op]
		if cb == nil {
			fmt.Printf("%04X: %02X (%q): %02X\n", f.PC, op, opcodes.Names[op], f.Code[f.PC:])
			panic("Unknown opcode reached")
		}
		cb(f)
	}
	return f.Return
}

func (f *ExecutionFrame) u8() uint8 {
	f.PC++
	return f.Code[f.PC-1]
}

func (f *ExecutionFrame) s24() int {
	result := int32(f.Code[f.PC]) << 8
	result |= int32(f.Code[f.PC+1]) << 16
	result |= int32(f.Code[f.PC+2]) << 24
	result >>= 8
	f.PC += 3
	return int(result)
}

func (f *ExecutionFrame) u30() int {
	result := int(f.Code[f.PC])
	f.PC++
	if (result & 0x80) == 0 {
		return result
	}
	result = (result & 0x7f) | (int(f.Code[f.PC]) << 7)
	f.PC++
	if (result & 0x4000) == 0 {
		return result
	}
	result = (result & 0x3fff) | (int(f.Code[f.PC]) << 14)
	f.PC++
	if (result & 0x200000) == 0 {
		return result
	}
	result = (result & 0x1fffff) | (int(f.Code[f.PC]) << 21)
	f.PC++
	return result
}

func (f *ExecutionFrame) push(x Any) {
	f.Stack = append(f.Stack, x)
}

func (f *ExecutionFrame) pop() Any {
	value := f.Stack[len(f.Stack)-1]
	f.Stack = f.Stack[:len(f.Stack)-1]
	return value
}

func (f *ExecutionFrame) top() Any {
	return f.Stack[len(f.Stack)-1]
}

func (f *ExecutionFrame) set_top(x Any) {
	f.Stack[len(f.Stack)-1] = x
}

func (f *ExecutionFrame) pop_args(n int) []Any {
	result := f.Stack[len(f.Stack)-n:]
	f.Stack = f.Stack[:len(f.Stack)-n]
	return result
}

func (f *ExecutionFrame) popMultiname(mn *abc.Multiname) abc.Multiname {
	var rn abc.Multiname
	rn.Name = mn.Name
	rn.Spaces = mn.Spaces
	if mn.IsNameRuntime {
		panic("ni")
	}
	if mn.IsSpaceRuntime {
		panic("ni")
	}
	return rn
}

func vm_op_nop(_ *ExecutionFrame) {
	// ...
}

func vm_op_debug(f *ExecutionFrame) {
	f.u8()
	f.u30()
	f.u30()
	f.u30()
}

func vm_op_debug_line(f *ExecutionFrame) {
	f.u30()
}

func vm_op_debug_file(f *ExecutionFrame) {
	f.u30()
}

func vm_op_return_void(f *ExecutionFrame) {
	f.Return = Undefined
}

func vm_op_return_value(f *ExecutionFrame) {
	f.Return = f.pop()
}

func vm_op_push_null(f *ExecutionFrame) {
	f.push(Null)
}

func vm_op_push_undefined(f *ExecutionFrame) {
	panic("vm_op_push_undefined")
}

func vm_op_push_int8(f *ExecutionFrame) {
	f.push(Int(int8(f.u8())))
}

func vm_op_push_int16(f *ExecutionFrame) {
	f.push(Int(int16(f.u30())))
}

func vm_op_push_true(f *ExecutionFrame) {
	panic("vm_op_push_true")
}

func vm_op_push_false(f *ExecutionFrame) {
	panic("vm_op_push_false")
}

func vm_op_push_nan(f *ExecutionFrame) {
	panic("vm_op_push_nan")
}

func vm_op_push_string(f *ExecutionFrame) {
	f.push(String(f.Pool.Strings[f.u30()]))
}

func vm_op_push_int32(f *ExecutionFrame) {
	panic("vm_op_push_int32")
}

func vm_op_push_uint32(f *ExecutionFrame) {
	f.push(Uint(f.Pool.Uints[f.u30()]))
}

func vm_op_push_float64(f *ExecutionFrame) {
	f.push(Number(f.Pool.Doubles[f.u30()]))
}

func vm_op_pop(f *ExecutionFrame) {
	f.pop()
}

func vm_op_dup(f *ExecutionFrame) {
	f.push(f.top())
}

func vm_op_swap(f *ExecutionFrame) {
	panic("vm_op_swap")
}

func vm_op_get_local(f *ExecutionFrame) {
	f.push(f.Locals[f.u30()])
}

func vm_op_get_local_0(f *ExecutionFrame) {
	f.push(f.Locals[0])
}

func vm_op_get_local_1(f *ExecutionFrame) {
	f.push(f.Locals[1])
}

func vm_op_get_local_2(f *ExecutionFrame) {
	f.push(f.Locals[2])
}

func vm_op_get_local_3(f *ExecutionFrame) {
	f.push(f.Locals[3])
}

func vm_op_set_local(f *ExecutionFrame) {
	f.Locals[f.u30()] = f.pop()
}

func vm_op_set_local_0(f *ExecutionFrame) {
	f.Locals[0] = f.pop()
}

func vm_op_set_local_1(f *ExecutionFrame) {
	f.Locals[1] = f.pop()
}

func vm_op_set_local_2(f *ExecutionFrame) {
	f.Locals[2] = f.pop()
}

func vm_op_set_local_3(f *ExecutionFrame) {
	f.Locals[3] = f.pop()
}

func vm_op_erase_local(f *ExecutionFrame) {
	f.Locals[f.u30()] = Undefined
}

func vm_op_push_scope(f *ExecutionFrame) {
	f.Scope = NewScope(f.Scope, f.Core.BoxIfPrimitive(f.pop()))
}

func vm_op_pop_scope(f *ExecutionFrame) {
	f.Scope = f.Scope.Parent
}

func vm_op_get_scope_object(f *ExecutionFrame) {
	f.push(f.Scope.GetAtHeight(int(f.u8()) + 1).Object)
}

func vm_op_add(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	f.push(Add(first, second))
}

func vm_op_subtract(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	f.push(Subtract(first, second))
}

func vm_op_multiply(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	f.push(Multiply(first, second))
}

func vm_op_divide(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	f.push(Divide(first, second))
}

func vm_op_modulo(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	f.push(Modulo(first, second))
}

func vm_op_negate(f *ExecutionFrame) {
	f.set_top(-ToNumber(f.top()))
}

func vm_op_negate_int(f *ExecutionFrame) {
	f.set_top(-ToInteger(f.top()))
}

func vm_op_boolean_not(f *ExecutionFrame) {
	f.set_top(!ToBoolean(f.top()))
}

func vm_op_bitwise_and(f *ExecutionFrame) {
	f.push(ToInteger(f.pop()) & ToInteger(f.pop()))
}

func vm_op_bitwise_or(f *ExecutionFrame) {
	f.push(ToInteger(f.pop()) | ToInteger(f.pop()))
}

func vm_op_bitwise_not(f *ExecutionFrame) {
	f.push(^ToInteger(f.pop()))
}

func vm_op_bitwise_xor(f *ExecutionFrame) {
	f.push(ToInteger(f.pop()) ^ ToInteger(f.pop()))
}

func vm_op_shift_left(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	f.push(ToUint(first) << (ToUint(second) & 0x1F))
}

func vm_op_shift_right_signed(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	f.push(ToInteger(first) >> Int(ToUint(second)&0x1F))
}

func vm_op_shift_right_unsigned(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	f.push(ToUint(first) >> (ToUint(second) & 0x1F))
}

func Increment(x Any) Any {
	return ToInteger(x) + 1
}

func Decrement(x Any) Any {
	return ToInteger(x) - 1
}

func IncrementInteger(x Any) Any {
	return Increment(x)
}

func vm_op_increment(f *ExecutionFrame) {
	f.set_top(Increment(f.top()))
}

func vm_op_decrement(f *ExecutionFrame) {
	f.set_top(Decrement(f.top()))
}

func vm_op_inc_local_int(f *ExecutionFrame) {
	index := f.u30()
	f.Locals[index] = IncrementInteger(f.Locals[index])
}

func vm_op_convert_float64(f *ExecutionFrame) {
	f.set_top(ToNumber(f.top()))
}

func vm_op_convert_uint32(f *ExecutionFrame) {
	f.set_top(ToUint(f.top()))
}

func vm_op_convert_int32(f *ExecutionFrame) {
	f.set_top(ToInteger(f.top()))
}

func vm_op_coerce_type(f *ExecutionFrame) {
	rn := f.popMultiname(&f.Pool.Links[f.u30()])
	value := f.top()
	obj := f.Scope.FindScopeProperty(&rn, true)
	receiver := obj.GetProperty(obj, rn.Spaces, rn.Name)
	f.set_top(f.Core.BoxIfPrimitive(receiver).Coerce(value))
}

func vm_op_coerce_string(f *ExecutionFrame) {
	x := f.top()
	if IsNullOrUndefined(x) {
		f.set_top(Null)
	} else {
		f.set_top(ToString(x))
	}
}

func vm_op_cmp_equals(f *ExecutionFrame) {
	second := f.pop()
	first := f.top()
	f.set_top(CompareEquals(first, second))
}

func vm_op_jump(f *ExecutionFrame) {
	f.PC += f.s24()
}

func vm_op_if_true(f *ExecutionFrame) {
	if ToBoolean(f.pop()) == true {
		offset := f.s24()
		f.PC += offset
	} else {
		f.PC += 3
	}
}

func vm_op_if_false(f *ExecutionFrame) {
	if ToBoolean(f.pop()) == false {
		offset := f.s24()
		f.PC += offset
	} else {
		f.PC += 3
	}
}

func vm_op_if_equal(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	if CompareEquals(first, second) == true {
		offset := f.s24()
		f.PC += offset
	} else {
		f.PC += 3
	}
}

func vm_op_if_less_than(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	if value, ok := CompareLessThan(first, second).(Boolean); ok && value == true {
		offset := f.s24()
		f.PC += offset
	} else {
		f.PC += 3
	}
}

func vm_op_if_not_equal(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	if CompareEquals(first, second) == false {
		offset := f.s24()
		f.PC += offset
	} else {
		f.PC += 3
	}
}

func vm_op_if_not_less_than(f *ExecutionFrame) {
	second := f.pop()
	first := f.pop()
	if value, ok := CompareLessThan(first, second).(Boolean); ok && value == false {
		offset := f.s24()
		f.PC += offset
	} else {
		f.PC += 3
	}
}

func vm_op_new_array(f *ExecutionFrame) {
	argCount := f.u30()
	args := f.Stack[len(f.Stack)-argCount:]
	f.Stack = f.Stack[:len(f.Stack)-argCount]

	f.push(f.Core.CreateArray(args))
}

func vm_op_new_class(f *ExecutionFrame) {
	super := f.top()
	class := &f.Pool.Classes[f.u30()]
	if class.IsInterface {
		f.set_top(f.Core.CreateInterface(f.Pool, class))
	} else {
		f.set_top(f.Core.CreateClass(f.Pool, class, f.Scope, super.(Constructible)))
	}
}

func vm_op_apply_vector(f *ExecutionFrame) {
	count := f.u30()
	if count != 1 {
		panic("nope")
	}

	arg := f.pop()

	if vector, isVector := f.top().(*Class_Vector); isVector {
		f.set_top(vector.ApplyType(f.Core, arg))
	} else {
		f.pop()
		panic("nope")
	}
}

func vm_op_find_property_strict(f *ExecutionFrame) {
	rn := f.popMultiname(&f.Pool.Links[f.u30()])
	f.push(f.Scope.FindScopeProperty(&rn, true))
}

func vm_op_find_property(f *ExecutionFrame) {
	rn := f.popMultiname(&f.Pool.Links[f.u30()])
	f.push(f.Scope.FindScopeProperty(&rn, false))
}

func vm_op_find_and_get_property(f *ExecutionFrame) {
	mn := &f.Pool.Links[f.u30()]
	if mn.IsRuntime() {
		panic("Runtime mn for getlex?")
	}
	obj := f.Scope.FindScopeProperty(mn, true)
	f.push(obj.GetProperty(obj, mn.Spaces, mn.Name))
}

func vm_op_set_property(f *ExecutionFrame) {
	value := f.pop()
	mn := &f.Pool.Links[f.u30()]
	if !mn.IsRuntime() {
		obj := f.Core.BoxIfPrimitive(f.pop())
		obj.SetProperty(obj, mn.Spaces, mn.Name, value)
	} else if !mn.IsSpaceRuntime && IntegerCanBeUnsigned(f.Stack[len(f.Stack)-1]) && IsObject(f.Stack[len(f.Stack)-2]) {
		index := uint(ToUint(f.pop()))
		obj := f.Core.BoxIfPrimitive(f.pop())
		obj.SetIndexedProperty(obj, mn.Spaces, index, value)
	} else if !mn.IsSpaceRuntime && IsObject(f.Stack[len(f.Stack)-2]) {
		name := string(ToString(f.pop()))
		obj := f.Core.BoxIfPrimitive(f.pop())
		obj.SetProperty(obj, mn.Spaces, name, value)
	} else {
		// fmt.Println(stack[len(stack)-2:])
		// fmt.Println(mn)
		// fmt.Println(mn.IsSpaceRuntime)
		// fmt.Println(IntegerCanBeUnsigned(stack[len(stack)-1]))
		// fmt.Println(stack[len(stack)-1].(String))
		// fmt.Println(IsObject(stack[len(stack)-2]))
		panic("Not implemented")
	}
}

func vm_op_get_property(f *ExecutionFrame) {
	mn := &f.Pool.Links[f.u30()]
	if !mn.IsRuntime() {
		obj := f.Core.BoxIfPrimitive(f.pop())
		f.push(obj.GetProperty(obj, mn.Spaces, mn.Name))
	} else if !mn.IsSpaceRuntime && IntegerCanBeUnsigned(f.Stack[len(f.Stack)-1]) && IsObject(f.Stack[len(f.Stack)-2]) {
		index := uint(ToUint(f.pop()))
		obj := f.Core.BoxIfPrimitive(f.pop())
		f.push(obj.GetIndexedProperty(obj, mn.Spaces, index))
	} else if !mn.IsSpaceRuntime && IsObject(f.Stack[len(f.Stack)-2]) {
		panic("boo")
		// name := string(ToString(pop()))
		// obj := core.BoxIfPrimitive(pop())
		// push(obj.GetProperty(obj, mn.Spaces, name))
	} else {
		// fmt.Printf("Code: %02X\n", code)
		// fmt.Println(stack[len(stack)-2:])
		// fmt.Println(mn)
		// fmt.Println(mn.IsSpaceRuntime)
		// fmt.Println(IntegerCanBeUnsigned(stack[len(stack)-1]))
		// fmt.Println(stack[len(stack)-1].(Uint))
		// fmt.Println(IsObject(stack[len(stack)-2]))
		panic("Not implemented")
	}
}

func vm_op_init_property(f *ExecutionFrame) {
	value := f.pop()
	rn := f.popMultiname(&f.Pool.Links[f.u30()])
	obj := f.Core.BoxIfPrimitive(f.pop())
	obj.SetProperty(obj, rn.Spaces, rn.Name, value)
}

func vm_op_construct(f *ExecutionFrame) {
	args := f.pop_args(f.u30())
	if obj, isConstructor := f.top().(Constructible); isConstructor {
		f.set_top(obj.Construct(args))
	} else {
		f.pop()
		panic("not a constructible")
	}
}

func vm_op_call_property_impl(f *ExecutionFrame) Any {
	mn := &f.Pool.Links[f.u30()]

	args := f.pop_args(f.u30())

	rn := f.popMultiname(mn)

	obj := f.pop()
	traits := f.Core.GetTraitsFor(obj)
	trait := traits.FindTrait(rn.Spaces, rn.Name)

	if trait == nil {
		if value, isValue := obj.(Value); isValue {
			return value.CallProperty(rn.Spaces, rn.Name, obj, args)
		} else {
			value := f.Core.GetPrototypeFor(obj)
			x := value.GetProperty(value, rn.Spaces, rn.Name)
			if ctor, isCtor := x.(Constructible); isCtor {
				return ctor.Call(obj, args)
			} else {
				panic("ni")
			}
		}
	}

	switch trait.Kind {
	case MethodTrait:
		return trait.Method.Call(obj, args)
	case ValueTrait, ConstTrait, ClassTrait:
		value := f.Core.BoxIfPrimitive(obj)
		x := value.GetSlot(trait.Slot)
		if ctor, isCtor := x.(Constructible); isCtor {
			return ctor.Call(obj, args)
		} else {
			panic("ni")
		}
	case GetterTrait, GetterSetterTrait:
		panic("ni")
	case SetterTrait:
		panic("error: write only")
	default:
		panic("unknown trait kind")
	}
}

func vm_op_call_property(f *ExecutionFrame) {
	f.push(vm_op_call_property_impl(f))
}

func vm_op_call_property_void(f *ExecutionFrame) {
	vm_op_call_property_impl(f)
}

func vm_op_construct_super(f *ExecutionFrame) {
	args := f.pop_args(f.u30())
	obj := f.pop()
	super := f.ParentScope.Object.(Constructible).GetSuper()
	ctor := super.GetInstanceTraits().RunConstructor
	if ctor == nil {
		panic("Ctor missing on " + super.GetDebugName())
	}
	ctor(obj, args)
}

func vm_op_construct_property(f *ExecutionFrame) {
	mn := &f.Pool.Links[f.u30()]

	args := f.pop_args(f.u30())

	rn := f.popMultiname(mn)

	obj := f.pop()
	traits := f.Core.GetTraitsFor(obj)
	trait := traits.FindTrait(rn.Spaces, rn.Name)
	switch trait.Kind {
	case MethodTrait:
		panic("error: not constructor")
	case ValueTrait, ConstTrait, ClassTrait:
		value := f.Core.BoxIfPrimitive(obj)
		x := value.GetSlot(trait.Slot)
		if ctor, isCtor := x.(Constructible); isCtor {
			f.push(ctor.Construct(args))
		} else {
			panic("ni")
		}
	case GetterTrait, GetterSetterTrait:
		panic("ni")
	case SetterTrait:
		panic("error: write only")
	default:
		value, isValue := obj.(Value)
		if !isValue {
			value = f.Core.GetPrototypeFor(obj)
		}
		x := value.GetProperty(value, rn.Spaces, rn.Name)
		if ctor, isCtor := x.(Constructible); isCtor {
			f.push(ctor.Construct(args))
		} else {
			panic("ni")
		}
	}
}

func init() {
	// Fix golang dumbness
	switchTable[opcodes.NewClass] = vm_op_new_class
	switchTable[opcodes.FindPropertyStrict] = vm_op_find_property_strict
	switchTable[opcodes.FindProperty] = vm_op_find_property
	switchTable[opcodes.FindAndGetProperty] = vm_op_find_and_get_property
	switchTable[opcodes.CoerceToType] = vm_op_coerce_type
}

var switchTable = [256]func(f *ExecutionFrame){
	// Ops that we leave ignored
	opcodes.NoOperation: vm_op_nop,
	opcodes.CoerceToAny: vm_op_nop,
	opcodes.Label:       vm_op_nop,
	opcodes.Debug:       vm_op_debug,
	opcodes.DebugLine:   vm_op_debug_line,
	opcodes.DebugFile:   vm_op_debug_file,
	// Control flow terminators
	opcodes.Throw:       nil,
	opcodes.ReturnVoid:  vm_op_return_void,
	opcodes.ReturnValue: vm_op_return_value,
	// Push value
	opcodes.PushNull:       vm_op_push_null,
	opcodes.PushUndefined:  vm_op_push_undefined,
	opcodes.PushInt8:       vm_op_push_int8,
	opcodes.PushInt16:      vm_op_push_int16,
	opcodes.PushTrue:       vm_op_push_true,
	opcodes.PushFalse:      vm_op_push_false,
	opcodes.PushNotANumber: vm_op_push_nan,
	opcodes.PushString:     vm_op_push_string,
	opcodes.PushInt:        vm_op_push_int32,
	opcodes.PushUint:       vm_op_push_uint32,
	opcodes.PushNumber:     vm_op_push_float64,
	opcodes.PushNamespace:  nil,
	// Normal stack ops
	opcodes.Pop:       vm_op_pop,
	opcodes.Duplicate: vm_op_dup,
	opcodes.Swap:      vm_op_swap,
	// Local ops
	opcodes.GetLocal:   vm_op_get_local,
	opcodes.GetLocal_0: vm_op_get_local_0,
	opcodes.GetLocal_1: vm_op_get_local_1,
	opcodes.GetLocal_2: vm_op_get_local_2,
	opcodes.GetLocal_3: vm_op_get_local_3,
	opcodes.SetLocal:   vm_op_set_local,
	opcodes.SetLocal_0: vm_op_set_local_0,
	opcodes.SetLocal_1: vm_op_set_local_1,
	opcodes.SetLocal_2: vm_op_set_local_2,
	opcodes.SetLocal_3: vm_op_set_local_3,
	opcodes.EraseLocal: vm_op_erase_local,
	// Scope stack ops
	opcodes.PushScope:      vm_op_push_scope,
	opcodes.PopScope:       vm_op_pop_scope,
	opcodes.PushWithScope:  nil,
	opcodes.GetGlobalScope: nil,
	opcodes.GetScopeObject: vm_op_get_scope_object,
	opcodes.GetOuterScope:  nil,
	// Mathematical operators
	opcodes.Add:             vm_op_add,
	opcodes.AddInteger:      nil,
	opcodes.Subtract:        vm_op_subtract,
	opcodes.SubtractInteger: nil,
	opcodes.Multiply:        vm_op_multiply,
	opcodes.MultiplyInteger: nil,
	opcodes.Divide:          vm_op_divide,
	opcodes.Modulo:          vm_op_modulo,
	opcodes.Negate:          vm_op_negate,
	opcodes.NegateInteger:   vm_op_negate_int,
	opcodes.BooleanNot:      vm_op_boolean_not,
	// Bitwise operations
	opcodes.BitwiseAnd:         vm_op_bitwise_and,
	opcodes.BitwiseOr:          vm_op_bitwise_or,
	opcodes.BitwiseNot:         vm_op_bitwise_not,
	opcodes.BitwiseXor:         vm_op_bitwise_xor,
	opcodes.ShiftLeft:          vm_op_shift_left,
	opcodes.ShiftRightSigned:   vm_op_shift_right_signed,
	opcodes.ShiftRightUnsigned: vm_op_shift_right_unsigned,
	// Increment/decrement opcodes
	opcodes.Increment:             vm_op_increment,
	opcodes.Decrement:             vm_op_decrement,
	opcodes.IncrementLocal:        nil,
	opcodes.DecrementLocal:        nil,
	opcodes.IncrementInteger:      nil,
	opcodes.DecrementInteger:      nil,
	opcodes.IncrementLocalInteger: vm_op_inc_local_int,
	opcodes.DecrementLocalInteger: nil,
	// Type operations
	opcodes.ConvertToObject:   nil,
	opcodes.ConvertToString:   nil,
	opcodes.ConvertToNumber:   vm_op_convert_float64,
	opcodes.ConvertToBoolean:  nil,
	opcodes.ConvertToUnsigned: vm_op_convert_uint32,
	opcodes.ConvertToInteger:  vm_op_convert_int32,
	opcodes.CoerceToType:      nil, // see init
	opcodes.CoerceToBoolean:   nil,
	opcodes.CoerceToInteger:   nil,
	opcodes.CoerceToNumber:    nil,
	opcodes.CoerceToString:    vm_op_coerce_string,
	opcodes.CoerceToUnsigned:  nil,
	opcodes.CoerceToObject:    nil,
	// Comparison
	opcodes.CompareEquals:        vm_op_cmp_equals,
	opcodes.CompareStrictEquals:  nil,
	opcodes.CompareLessThan:      nil,
	opcodes.CompareLessEquals:    nil,
	opcodes.CompareGreaterThan:   nil,
	opcodes.CompareGreaterEquals: nil,
	// Jump, branches, switch
	opcodes.Jump:              vm_op_jump,
	opcodes.IfTrue:            vm_op_if_true,
	opcodes.IfFalse:           vm_op_if_false,
	opcodes.IfEqual:           vm_op_if_equal,
	opcodes.IfLessThan:        vm_op_if_less_than,
	opcodes.IfLessEqual:       nil,
	opcodes.IfGreaterThan:     nil,
	opcodes.IfGreaterEqual:    nil,
	opcodes.IfStrictEqual:     nil,
	opcodes.IfNotEqual:        vm_op_if_not_equal,
	opcodes.IfNotLessThan:     vm_op_if_not_less_than,
	opcodes.IfNotLessEqual:    nil,
	opcodes.IfNotGreaterThan:  nil,
	opcodes.IfNotGreaterEqual: nil,
	opcodes.IfStrictNotEqual:  nil,
	opcodes.LookupSwitch:      nil,
	// *new* opcodes
	opcodes.NewFunction:   nil,
	opcodes.NewObject:     nil,
	opcodes.NewArray:      vm_op_new_array,
	opcodes.NewActivation: nil,
	opcodes.NewClass:      nil, // see init
	opcodes.NewCatch:      nil,
	// Type- related opcodes
	opcodes.ApplyVectorType:    vm_op_apply_vector,
	opcodes.CheckFilter:        nil,
	opcodes.CheckType:          nil,
	opcodes.CheckTypeLate:      nil,
	opcodes.GetTypeName:        nil,
	opcodes.InstanceOf:         nil,
	opcodes.IsObjectOfType:     nil,
	opcodes.IsObjectOfTypeLate: nil,
	// Property access
	opcodes.FindPropertyStrict:     nil, // see init
	opcodes.FindProperty:           nil, // see init
	opcodes.FindScriptDefinition:   nil,
	opcodes.FindAndGetProperty:     nil, // see init
	opcodes.SetProperty:            vm_op_set_property,
	opcodes.GetProperty:            vm_op_get_property,
	opcodes.InitProperty:           vm_op_init_property,
	opcodes.DeleteProperty:         nil,
	opcodes.GetSlot:                nil,
	opcodes.SetSlot:                nil,
	opcodes.GetGlobalSlot:          nil,
	opcodes.SetGlobalSlot:          nil,
	opcodes.GetSuper:               nil,
	opcodes.SetSuper:               nil,
	opcodes.ObjectHasProperty:      nil,
	opcodes.GetNameOfNextProperty:  nil,
	opcodes.GetNextPropertyIndex:   nil,
	opcodes.GetValueOfNextProperty: nil,
	opcodes.GetNextPropertyOnLocal: nil,
	opcodes.GetDescendants:         nil,
	// Calling, construction
	opcodes.Call:                nil,
	opcodes.Construct:           vm_op_construct,
	opcodes.CallMethod:          nil,
	opcodes.CallStatic:          nil,
	opcodes.CallSuper:           nil,
	opcodes.CallProperty:        vm_op_call_property,
	opcodes.CallPropertyVoid:    vm_op_call_property_void,
	opcodes.ConstructSuper:      vm_op_construct_super,
	opcodes.ConstructProperty:   vm_op_construct_property,
	opcodes.FindAndCallProperty: nil,
	opcodes.CallSuperVoid:       nil,
	// Memory access extension
	opcodes.LoadInt8:         nil,
	opcodes.LoadInt16:        nil,
	opcodes.LoadInt32:        nil,
	opcodes.LoadFloat32:      nil,
	opcodes.LoadFloat64:      nil,
	opcodes.StoreInt8:        nil,
	opcodes.StoreInt16:       nil,
	opcodes.StoreInt32:       nil,
	opcodes.StoreFloat32:     nil,
	opcodes.StoreFloat64:     nil,
	opcodes.SignExtend1To32:  nil,
	opcodes.SignExtend8To32:  nil,
	opcodes.SignExtend16To32: nil,
	// XML table extension
	opcodes.SetDefaultXmlSpace:     nil,
	opcodes.SetDefaultXmlSpaceLate: nil,
	opcodes.EscapeXmlElement:       nil,
	opcodes.EscapeXmlAttribute:     nil,
}
