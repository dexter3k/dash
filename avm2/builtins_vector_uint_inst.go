package avm2

func (this *Object_Vector_uint) GetIndexedProperty(self Any, spaces []string, index uint) Any {
	if index < uint(len(this.data)) {
		return this.data[index]
	}
	return this.Object_Object.GetIndexedProperty(self, spaces, index)
}

func (this *Object_Vector_uint) SetIndexedProperty(self Any, spaces []string, index uint, value Any) {
	v := ToUint(value)
	if index < uint(len(this.data)) {
		this.data[index] = v
	} else if index == uint(len(this.data)) {
		this.data = append(this.data, v)
	} else if (index - uint(len(this.data))) <= 16 {
		// Adding value somewhere close, but not in order
		missing := 1 + index - uint(len(this.data))
		this.data = append(this.data, make([]Uint, missing)...)
		this.data[index] = v
	} else {
		panic(index)
	}
}
