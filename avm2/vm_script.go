package avm2

import (
	"github.com/dexter3k/dash/abc"
)

type ScriptState int

const (
	NotExecuted ScriptState = iota
	Executing
	Executed
)

type Script struct {
	Abc    *abc.File
	State  ScriptState
	Script *abc.Script
	Global *Global
	Scope  *Scope
}

func NewScript(file *abc.File, script *abc.Script) *Script {
	return &Script{
		Abc:    file,
		State:  NotExecuted,
		Script: script,
	}
}
