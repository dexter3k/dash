package abc

/* Full name in format namespace.class */
func (f *File) FindClass(fullName string) *Class {
	fullName = "P" + fullName // we're searching the public namespace
	for i := 0; i < len(f.Classes); i++ {
		class := &f.Classes[i]
		link := class.Name
		if fullName == (link.Spaces[0] + "." + link.Name) {
			return class
		}
	}
	return nil
}
