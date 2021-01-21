package avm2

import (
	"fmt"
	"strconv"
)

func CreateObjectForTraits(traits *Traits) *Object_Object {
	obj := &Object_Object{
		Traits:  traits,
		Dynamic: map[string]Any{},
	}
	obj.Fixed = traits.CreateSlots()
	return obj
}

func (this *Object_Object) InitTraits(traits *Traits) {
	if this.Traits != nil && this.Traits != traits {
		panic("reapplication of traits")
	}
	this.Dynamic = map[string]Any{}
	this.Traits = traits
	this.Fixed = traits.CreateSlots()
}

func (self *Object_Object) GetTraits() *Traits {
	if self.Traits == nil {
		panic("GetTraits: traits object missing on " + self.GetDebugName())
	}
	return self.Traits
}

func (o *Object_Object) GetSlot(index uint) Any {
	if o.Traits == nil || o.Fixed == nil {
		panic("GetSlot: traits object missing")
	}
	if uint(len(o.Fixed)) < index {
		panic("GetSlot: out of bounds")
	}
	return o.Fixed[index-1]
}

func (o *Object_Object) SetSlot(index uint, value Any) {
	if o.Traits == nil || o.Fixed == nil {
		panic("SetSlot: traits object missing")
	}
	if uint(len(o.Fixed)) < index {
		panic("SetSlot: out of bounds")
	}
	o.Fixed[index-1] = value
}

func (this *Object_Object) Coerce(any Any) Any {
	if IsNullOrUndefined(any) {
		return Null
	}
	return any
}

func (o *Object_Object) CallProperty(spaces []string, name string, self Any, args []Any) Any {
	if o.Traits == nil {
		panic("CallProperty: traits object missing")
	}

	var value Any
	t := o.Traits.FindTrait(spaces, name)
	if t != nil {
		if t.Method != nil {
			return t.Method.Execute(self, args)
		}
		if t.Slot == 0 {
			panic("nope")
		}
		value = o.Fixed[t.Slot-1]
	} else {
		var searchDynamic bool
		for _, space := range spaces {
			if space == "P" {
				searchDynamic = true
			}
		}
		if !searchDynamic {
			panic(spaces[0])
			panic("Not found")
		}

		if prop, found := o.Dynamic[name]; found {
			value = prop
		} else {
			if o.Prototype == nil {
				panic("ni")
			} else {
				return o.Prototype.CallProperty(spaces, name, self, args)
			}
		}
	}

	if prop, ctible := value.(Constructible); ctible {
		return prop.Call(self, args)
	} else {
		panic("ni")
	}
}

func (o *Object_Object) HasTrait(spaces []string, name string) bool {
	if o.Traits == nil {
		fmt.Println(o, o.Fixed == nil)
		panic("HasTrait: traits object missing on " + o.GetDebugName())
	}
	return o.Traits.FindTrait(spaces, name) != nil
}

func (o *Object_Object) GetProperty(self Any, spaces []string, name string) Any {
	if o.Traits == nil {
		panic("missing traits on object")
	}
	t := o.Traits.FindTrait(spaces, name)
	if t != nil {
		if t.Get != nil {
			return t.Get(self, nil)
		}
		if t.Method != nil {
			return t.Method
		}
		// if t.Class != nil {
		// 	return t.Class
		// }
		if t.Slot == 0 {
			panic("nope")
		}
		return o.Fixed[t.Slot-1]
	} else {
		// Check we have a "P" namespace in spaces
		var searchDynamic bool
		for _, space := range spaces {
			if space == "P" {
				searchDynamic = true
				break
			}
		}
		if !searchDynamic {
			panic("ni")
			return Undefined
		}

		if prop, found := o.Dynamic[name]; found {
			return prop
		} else {
			if o.Prototype == nil {
				panic("ni")
			} else {
				return o.Prototype.GetProperty(self, spaces, name)
			}
		}
	}
}

func (o *Object_Object) SetProperty(self Any, spaces []string, name string, value Any) {
	if o.Traits != nil {
		t := o.Traits.FindTrait(spaces, name)
		if t != nil {
			if t.Set != nil {
				t.Set(self, []Any{value})
				return
			}
			if t.Slot == 0 {
				panic("nope")
			}
			o.Fixed[t.Slot-1] = value
			return
		}
	}

	// Check we have a "P" namespace in spaces
	var searchDynamic bool
	for _, space := range spaces {
		if space == "P" {
			searchDynamic = true
			break
		}
	}
	if !searchDynamic {
		panic("Not found")
	}

	o.Dynamic[name] = value
}

func (o *Object_Object) GetIndexedProperty(self Any, spaces []string, index uint) Any {
	return o.GetProperty(self, spaces, strconv.FormatUint(uint64(index), 10))
}

func (o *Object_Object) SetIndexedProperty(self Any, spaces []string, index uint, value Any) {
	o.SetProperty(self, spaces, strconv.FormatUint(uint64(index), 10), value)
}

func (this *Object_Object) GetDebugName() string {
	return "[object " + this.DebugName + "]"
}
