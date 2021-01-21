package avm2

import (
	"github.com/dexter3k/dash/abc"
)

type Domain struct {
	Core   *Core
	Parent *Domain

	// Defines[name][space] = script
	Defines map[string](map[string]*Script)
}

func NewDomain(core *Core, parent *Domain) *Domain {
	return &Domain{
		Core:    core,
		Parent:  parent,
		Defines: map[string](map[string]*Script){},
	}
}

func (d *Domain) LoadLazyAbcData(name string, data []byte) {
	file := abc.Parse(data)
	d.LoadLazyAbc(file)
}

func (d *Domain) LoadAndExecuteAbcData(name string, data []byte) {
	file := abc.Parse(data)
	d.LoadAndExecuteAbc(file)
}

// Load abc, but don't run its last script yet
func (d *Domain) LoadLazyAbc(file *abc.File) {
	d.loadAbc(file)
}

func (d *Domain) LoadAndExecuteAbc(file *abc.File) {
	lastScript := d.loadAbc(file)
	if lastScript == nil {
		panic("No scripts in ABC?")
	}

	d.executeScript(lastScript)
}

// Bug?: before searching domains we search traits on global
// object. Is this ok? This does work for working AS3 code
// but we can compile a program that breaks this.
//
func (d *Domain) FindProperty(mn *abc.Multiname) Value {
	if d.Parent != nil {
		if result := d.Parent.FindProperty(mn); result != nil {
			return result
		}
	}
	if mapping, found := d.Defines[mn.Name]; found {
		for i := 0; i < len(mn.Spaces); i++ {
			if script, found := mapping[mn.Spaces[i]]; found {
				if script.State == NotExecuted {
					d.executeScript(script)
				}
				return script.Global
			}
		}
	}
	return nil
}

// Add users builtins script to this domain
func (d *Domain) ApplyUserScript(avmScript *Script) {
	avmScript.Global.Fixed = avmScript.Global.Traits.CreateSlots()
	traits := avmScript.Global.Traits.GetOwnTraits()
	for _, trait := range traits {
		d.addScriptToDefines(trait.Link.Name(), trait.Link.Space(), avmScript)
	}
}

// For each script abc defines creates our own
// script object and updates list of own Defines.
// Last script processed is returned
// Returns nil if no scripts defined in file
func (d *Domain) loadAbc(file *abc.File) *Script {
	var lastScript *Script
	for i := 0; i < len(file.Scripts); i++ {
		abcScript := &file.Scripts[i]
		avmScript := NewScript(file, abcScript)
		if i == len(file.Scripts)-1 {
			lastScript = avmScript
		}

		for _, t := range abcScript.Traits {
			if t.Link.IsRuntime() {
				panic("Script trait name is runtime!?")
			}

			d.addScriptToDefines(t.Link.Name, t.Link.Spaces[0], avmScript)
		}
	}
	return lastScript
}

// Expects script to be marked as executed and to have Global
// Defines are cloned from traits found on Global object
func (d *Domain) loadNativeScript(script *Script) {
	if script.State != Executed {
		panic("Native script is not executed")
	}

	g := script.Global
	if g == nil {
		panic("Global is missing from native script")
	}

	if g.Domain != nil {
		panic("Global already has a domain")
	}
	g.Domain = d

	traits := g.Traits
	if traits == nil {
		panic("Global has no traits initialized!")
	}

	list := traits.GetOwnTraits()
	if list == nil {
		panic("Cannot happen")
	}

	for i := 0; i < len(list); i++ {
		// fmt.Println(i, list[i])
		link := list[i].Link
		d.addScriptToDefines(link.Name(), link.Space(), script)
	}
}

func (d *Domain) addScriptToDefines(name, space string, script *Script) {
	if mapping, found := d.Defines[name]; found {
		if _, found := mapping[space]; found {
			panic("Name collision in domain: " + name + " from " + space)
		} else {
			mapping[space] = script
		}
	} else {
		d.Defines[name] = map[string]*Script{
			space: script,
		}
	}
}

func (d *Domain) executeScript(s *Script) {
	if s.State != NotExecuted {
		panic("Recursive script execution!")
	}

	// At this point we should have
	// 1) Working builtins
	// 2) Loaded ABC
	// We also need to construct a working
	// global for this script
	d.Core.CreateGlobalOnScript(s, d)

	// And then we can execute the actual code
	// passing global as this
	s.State = Executing
	// We don't need to provide any self value because core
	// has already set it up

	s.Global.Traits.Initializer(nil)
	s.State = Executed
}
