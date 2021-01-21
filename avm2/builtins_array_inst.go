package avm2

func (array *Object_Array) GetIndexedProperty(self Any, spaces []string, index uint) Any {
	// fmt.Printf("READ Array[%d]\n", index)
	if index < uint(len(array.values)) {
		value := array.values[index]
		if value == nil {
			return Undefined
		}
		return value
	}
	return array.Object_Object.GetIndexedProperty(self, spaces, index)
}

func (array *Object_Array) SetIndexedProperty(self Any, spaces []string, index uint, value Any) {
	// fmt.Printf("WRITE Array[%d] = %08X\n", index, ToUint(value))
	if index < uint(len(array.values)) {
		array.values[index] = value
	} else if index == uint(len(array.values)) {
		array.values = append(array.values, value)
		// fmt.Println("Array changed length to", len(array.values))
		// fmt.Println(array)
	} else if (index - uint(len(array.values))) <= 16 {
		// Adding value somewhere close, but not in order
		missing := 1 + index - uint(len(array.values))
		array.values = append(array.values, make([]Any, missing)...)
		array.values[index] = value
	} else {
		panic(index)
	}
	// panic(index)
	// self.SetProperty(self, spaces, strconv.FormatUint(uint64(index), 10), value)
}

func (this *Object_Array) avm2_as3_func_push(rest []Any) Uint {
	this.values = append(this.values, rest...)
	return Uint(len(this.values))
}
