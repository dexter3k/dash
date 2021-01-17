package avm2

import (
	"github.com/dexter3k/dash/abc"
)

type State struct {
	abcs []*abc.File
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
	s.abcs = append(s.abcs, file)
	if !lazy {
		s.initAbc(file)
	}
	return nil
}

func (s *State) initAbc(f *abc.File) {
	// ...
}
