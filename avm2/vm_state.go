package avm2

import (
	"fmt"

	"github.com/dexter3k/dash/abc"
)

type Script struct {
	Initialized  bool
	Initializing bool

	Scope  *Scope // maybe?
	Source *abc.Script
}

type ScriptCollection struct {
	File    *abc.File
	Scripts []*Script

	Defines map[abc.Space]map[string]*Script
}

func (sc *ScriptCollection) preInitScripts() {
	if sc.Scripts != nil {
		panic("wait what")
	}

	sc.Defines = map[abc.Space]map[string]*Script{}

	f := sc.File
	sc.Scripts = make([]*Script, len(f.Script))
	for i := 0; i < len(f.Script); i++ {
		s := &Script{
			Source: f.Script[i],
		}
		sc.Scripts[i] = s

		for _, anyTrait := range s.Source.Traits {
			link := anyTrait.GetFullName()
			name := link.Name
			space := link.Space

			names := sc.Defines[space]
			if names == nil {
				names = map[string]*Script{}
				sc.Defines[space] = names
			}
			names[name] = s

			switch trait := anyTrait.(type) {
			case *abc.ClassTrait:
				// fmt.Println("Class:", trait.Name)
			case *abc.SlotTrait:
				// fmt.Println("Slot:", trait.Name, trait.Const, trait.Type, trait.Index, trait.ValueKind)
			default:
				panic(trait)
			}
		}
	}

	if false {
		for space, names := range sc.Defines {
			fmt.Println("from space", space)
			for name, _ := range names {
				fmt.Println(" *", name)
			}
		}
	}
}

func (sc *ScriptCollection) initLastScript() {
	panic("todo")
}

func NewScriptCollection(f *abc.File) *ScriptCollection {
	sc := &ScriptCollection{
		File:    f,
	}
	sc.preInitScripts()
	return sc
}

type State struct {
	Apps []*ScriptCollection
}

func NewState() *State {
	s := &State{}
	return s
}

func (s *State) AddAbc(name string, data []byte, lazy bool) error {
	file, err := abc.Parse(data)
	if err != nil {
		return err
	}
	s.Apps = append(s.Apps, NewScriptCollection(file))
	if !lazy {
		s.Apps[len(s.Apps)-1].initLastScript()
	}
	return nil
}
