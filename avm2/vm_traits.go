/*

This file provides essentional functionality related to AS3 compile-time typed
data traits

Traits can be divided into two categories
1) data traits, aka slots
 - value of these is unique to each instance and thus take space per-instance.
 - ie they have value -> thy need to store it somewhere
 - Can be a variable, const or a class
 - consts take per-instance space because they are initialized on runtime also
   per instance. Ie constant holds value during lifetime of object, but can be
   different on different objects
 - classes are there for global objects and other classes. this is not trivial
   for why classes are like variables, but this is how it is. trait describes
   only the name of the class and in which spaces you might find it, but
   these are resolved per-instance for some reason.
 - though note that in the end these work pretty well
 - type and metainfo on each slot are stored on Traits object created to
 - be read-only shared for all instances that are pretty much the same
2) method traits
 - getters, setters, methods
 - setters
 - methods
 - ie these are pieces of code that operate on specific set of slots (aka instances)
 - these are not actual data and have no value
 - you can take value of a getter which will triggers its execution
 - you cant take value of a setter (not found)
 - you can take value of a method, but this will create MethodClosure linking this method
 - to current instance and return instance of this MethodClosure which essentionally is a
 - Function object.

ABC file provides traits on many things, not all of them are actually ever used
- Instance: traits describe the actual data that belongs

*/
package avm2

import (
	"fmt"

	"github.com/dexter3k/dash/abc"
	// "github.com/davecgh/go-spew/spew"
)

type TraitKind int

// Method, getter and setter traits are methods and are never copied
// for each instance of trait-host
// Value traits are traits for which a physical memory point is allocated
// per instance
// Constants, variables and even classes are value traits meant to be placed
// in "slots"
// Also no such thing as FunctionTrait actually exist
// functions are allocated in script initializer with NewFunction opcode
// GetterSetter is used when there are both present at the same time
const (
	MethodTrait TraitKind = iota
	GetterTrait
	SetterTrait
	GetterSetterTrait
	ValueTrait
	ConstTrait
	ClassTrait
)

type Trait struct {
	Link link

	Kind TraitKind

	Slot uint // vtable id for methods

	// plus bad idea cuz we have no link to corresponding ABC Pool data
	Class *abc.Class // todo: do we really need to know this?
	Value Any        // default value for slot traits, not runtime
	// Method func(self Value, args []Any) Any
	Method *Object_Function
	Set    func(self Any, args []Any) Any // todo: refactor []Any to Any
	Get    func(self Any, args []Any) Any // todo: remove args completely

	Override bool
}

func (t *Trait) IsSlot() bool {
	return t.Kind == ValueTrait || t.Kind == ConstTrait || t.Kind == ClassTrait
}

func (t *Trait) IsMethod() bool {
	return t.Kind == MethodTrait || t.Kind == GetterTrait || t.Kind == SetterTrait || t.Kind == GetterSetterTrait
}

// List of traits linked to every object instance
// Does not contain actual data, but contains info
// *about* data. Ie types, default value, methods
type Traits struct {
	parent *Traits // todo: cache parent entries?

	Initializer    func(args []Any) Any
	RunConstructor func(self Any, args []Any)

	// Internal map of traits used for lookup
	// lookup[name][space] = trait
	lookup  map[string](map[string]*Trait)
	vtable  []*Trait // list of virtual method overrides for this layer
	slots   []Trait  // list of all slots, including parent
	methods []Trait  // list of all methods

	protected      map[string]*Trait
	protectedSpace string
}

func newTraits(parent *Traits) *Traits {
	t := &Traits{
		parent:  parent,
		lookup:  map[string](map[string]*Trait){},
		methods: []Trait{},
	}
	if parent != nil {
		t.vtable = make([]*Trait, len(parent.vtable))
		copy(t.vtable, parent.vtable)
		t.slots = make([]Trait, len(parent.slots))
		copy(t.slots, parent.slots)
	}
	return t
}

// Apply traits from ABC data, also inheriting from parent traits
func ApplyAbcTraits(source []abc.Trait, parent *Traits, core *Core, pool *abc.File, scope *Scope, protectedSpace string, traitPrefix string) *Traits {
	t := newTraits(parent)
	if protectedSpace != "" {
		t.AddProtectedNamespace(protectedSpace)
	}

	for i := 0; i < len(source); i++ {
		raw := &source[i]
		if raw.Link.IsRuntime() {
			panic("Trait link is dynamic!?")
		}

		trait := Trait{
			Link:     Link(raw.Link.Name, raw.Link.Spaces[0]),
			Override: raw.IsOverride,
		}

		// Fix trait names for protected ones
		if trait.Link.IsProtected() {
			name := trait.Link.Name()
			old := t.FindProtected(name)
			if old != nil {
				trait.Link = Link(name, old.Link.Space())
			}
		}

		// Apply our custom trait types
		switch raw.Kind {
		case abc.TraitGetter:
			trait.Kind = GetterTrait
			trait.Get = func(self Any, args []Any) Any {
				return Execute(core, pool, raw.Method, scope, self, args)
			}
			if core.debugger != nil {
				core.debugger.RegisterMethodName(raw.Method, traitPrefix + "get " + trait.Link.Link())
			}
		case abc.TraitSetter:
			trait.Kind = SetterTrait
			trait.Set = func(self Any, args []Any) Any {
				return Execute(core, pool, raw.Method, scope, self, args)
			}
			if core.debugger != nil {
				core.debugger.RegisterMethodName(raw.Method, traitPrefix + "set " + trait.Link.Link())
			}
		case abc.TraitMethod:
			trait.Kind = MethodTrait
			trait.Method = core.WrapMethod(func(self Any, args []Any) Any {
				return Execute(core, pool, raw.Method, scope, self, args)
			})
			if core.debugger != nil {
				core.debugger.RegisterMethodName(raw.Method, traitPrefix + "meth " + trait.Link.Link())
			}
		case abc.TraitVariable:
			if raw.Value != nil {
				trait.Value = core.WrapDefaultValue(scope, raw.Value, raw.DataType)
			}
			trait.Kind = ValueTrait
		case abc.TraitConst:
			if raw.Value != nil {
				trait.Value = core.WrapDefaultValue(scope, raw.Value, raw.DataType)
			}
			trait.Kind = ConstTrait
		case abc.TraitClass:
			trait.Kind = ClassTrait
		default:
			panic("VerifyError: unknown trait kind")
		}

		t.AddTrait(trait)
	}
	return t
}

func (self *Traits) AddProtectedNamespace(space string) {
	if self.protectedSpace != "" {
		panic(self.protectedSpace)
	}

	if space != "" {
		self.protectedSpace = space
		self.protected = map[string]*Trait{}
	} else {
		panic(space)
	}
}

func (self *Traits) FindProtected(name string) *Trait {
	if self.protected != nil {
		if value, found := self.protected[name]; found {
			return value
		}
	}
	if self.parent != nil {
		return self.parent.FindProtected(name)
	}
	return nil
}

func (self *Traits) ApplyInterface(iface *Traits) {
	// spew.Dump(self)
	// spew.Dump(iface)
	// panic("Applying interface")

	// For each trait that IsMethod()==true, in iface we must have corresponding
	// trait in our type or in parent. Interfacing is recursive!
	if iface.parent != nil {
		panic("Nope i didnt test that sorry")
	}
	for _, trait := range iface.GetAllTraits() {
		if !trait.IsMethod() {
			continue
		}

		name := trait.Link.Name()

		ourTrait := self.FindTrait([]string{"P"}, name)
		if ourTrait == nil {
			panic("Broken interface implementation")
		}

		// Add interfaced shortcut to our lookup
		if mapping, found := self.lookup[name]; found {
			old := mapping[trait.Link.Space()]
			if old != nil {
				panic("the fyck?")
			}
			mapping[trait.Link.Space()] = ourTrait
		} else {
			self.lookup[name] = map[string]*Trait{
				trait.Link.Space(): ourTrait,
			}
		}
	}

	// panic("Henlo")
}

func (self *Traits) AddTrait(trait Trait) {
	var ptr *Trait
	if trait.IsSlot() {
		// A slot (var, const, class)
		// Can extend parent slots, but not override
		slot := trait.Slot
		if slot == 0 || slot-1 == uint(len(self.slots)) {
			self.slots = append(self.slots, trait)
			ptr = &self.slots[len(self.slots)-1]
			ptr.Slot = uint(len(self.slots))
		} else {
			panic("VerifyError: assigning slots out-of-order")
		}
	} else if trait.IsMethod() {
		// A method (, get, set)
		// Add to method list and apply to vtable if needed
		self.methods = append(self.methods, trait)
		ptr = &self.methods[len(self.methods)-1]

		// Apply to vtable only if such optimisation is enabled
		vslot := ptr.Slot
		if vslot != 0 {
			if vslot-1 == uint(len(self.vtable)) {
				// New vtable slot added in order
				self.vtable = append(self.vtable, ptr)
			} else if vslot-1 < uint(len(self.vtable)) {
				// Overriding vtable slot
				self.vtable[vslot-1] = ptr
			} else {
				panic("Insane virtual method")
			}
		}
	} else {
		panic("Fixme")
	}

	// Now that we have a pointer we can fix out protected
	if ptr.Link.IsProtected() {
		if self.protected == nil {
			fmt.Println(ptr.Link.Link())
			panic("wtf")
		}
		self.protected[ptr.Link.Name()] = ptr
	}

	// Add new trait to current layer, and find old one if present
	var old *Trait
	link := ptr.Link
	if mapping, found := self.lookup[link.Name()]; found {
		old = mapping[link.Space()]
		mapping[link.Space()] = ptr
	} else {
		self.lookup[link.Name()] = map[string]*Trait{
			link.Space(): ptr,
		}
	}

	// No old value in current trait list found, search for parent
	// Why? We might have eg setter on parent and getter on child
	// Without lookup getter will shadow the setter
	if old == nil && self.parent != nil {
		old = self.parent.FindTrait([]string{link.Space()}, link.Name())
	}

	// If old is eg getter and new is setter, then we need to make
	// the new one gettersetter
	if old != nil {
		if ptr.Kind == GetterTrait && old.Set != nil {
			ptr.Kind = GetterSetterTrait
			ptr.Set = old.Set
		} else if ptr.Kind == SetterTrait && old.Get != nil {
			ptr.Kind = GetterSetterTrait
			ptr.Get = old.Get
		} else {
			// fmt.Println("Traits override: ", ptr.Link.Link(), old.Link.Link())
			if !ptr.Override {
				panic(int(ptr.Kind))
			}
		}
	} else if ptr.Override {
		fmt.Println(ptr.Link.Link())
		// spew.Dump(self)
		panic("Override without override? wtf")
	}
}

func (self *Traits) GetSlot(index int) *Trait {
	if index == 0 || index > len(self.slots) {
		return nil
	}
	return &self.slots[index-1]
}

func (self *Traits) GetVirtual(index int) *Trait {
	if index == 0 || index > len(self.vtable) {
		return nil
	}
	return self.vtable[index-1]
}

func (self *Traits) GetMethod(index int) *Trait {
	if index == 0 || index > len(self.methods) {
		return nil
	}
	return &self.methods[index-1]
}

func (self *Traits) FindTrait(spaces []string, name string) *Trait {
	// TODO: caching if appropriate
	if mapping, found := self.lookup[name]; found {
		for i := 0; i < len(spaces); i++ {
			if trait, found := mapping[spaces[i]]; found {
				return trait
			}
		}
	}
	if self.parent != nil {
		// TODO: caching if appropriate
		return self.parent.FindTrait(spaces, name)
	}
	return nil
}

func (self *Traits) GetAllTraits() []*Trait {
	if self.parent != nil {
		return append(self.parent.GetAllTraits(), self.GetOwnTraits()...)
	}
	return self.GetOwnTraits()
}

func (self *Traits) GetOwnTraits() []*Trait {
	parentSlots := 0
	if self.parent != nil {
		parentSlots = len(self.parent.slots)
	}
	resultCount := len(self.slots) - parentSlots + len(self.methods)
	result := make([]*Trait, 0, resultCount)

	for i := parentSlots; i < len(self.slots); i++ {
		result = append(result, &self.slots[i])
	}
	for i := 0; i < len(self.methods); i++ {
		result = append(result, &self.methods[i])
	}

	return result
}

func (self *Traits) CreateSlots() []Any {
	values := make([]Any, 0, len(self.slots))
	for i := 0; i < len(self.slots); i++ {
		values = append(values, self.slots[i].Value)
	}
	return values
}
