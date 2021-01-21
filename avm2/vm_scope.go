package avm2

import (
	// "fmt"

	"github.com/dexter3k/dash/abc"
)

type Scope struct {
	Parent *Scope
	Global *Scope
	Object Value
	IsWith bool

	height int
}

func NewScope(parent *Scope, object Value) *Scope {
	s := &Scope{
		Parent: parent,
		Object: object,
	}
	if parent == nil {
		s.Global = s
	} else {
		s.Global = parent.Global
	}
	return s
}

func (s *Scope) GetHeight() int {
	if s.height == 0 {
		if s.Parent == nil {
			s.height = 1
		} else {
			s.height = s.Parent.GetHeight() + 1
		}
	}
	return s.height
}

func (s *Scope) GetAtHeight(height int) *Scope {
	own := s.GetHeight()
	if height > own {
		panic(height)
	} else if height == own {
		return s
	} else if s.Parent != nil {
		return s.Parent.GetAtHeight(height)
	} else {
		panic(height)
		return nil
	}
}

func (s *Scope) FindScopeProperty(mn *abc.Multiname, strict bool) Value {
	if s.Object != nil {
		if s.IsWith {
			panic("With scopes are not yet implemented")
		} else if s.Object.HasTrait(mn.Spaces, mn.Name) {
			return s.Object
		}
	}

	if s.Parent != nil {
		return s.Parent.FindScopeProperty(mn, strict)
	}

	if mn.IsAttribute {
		panic("Attributes cannot be on global")
	}

	g := s.Global.Object.(*Global)
	if object := g.Domain.FindProperty(mn); object != nil {
		return object
	}
	if strict {
		panic("Property " + mn.Name + " is not found")
	}
	return g
}
