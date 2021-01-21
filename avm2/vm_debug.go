package avm2

import (
	// "bufio"
	"fmt"
	// "log"
	// "os"
	"time"
	"sort"

	"github.com/dexter3k/dash/abc"
	// "github.com/davecgh/go-spew/spew"
)

// Basic support for code debugging

type Debugger interface {
	RegisterMethodName(m *abc.Method, name string)

	PrintProfilingInfo()

	EnterMethod(frame *ExecutionFrame) bool
	SteppingMode(frame *ExecutionFrame) bool
	LeaveMethod(frame *ExecutionFrame)
}

type profileInfo struct {
	EnterCount uint64
	TotalTime  time.Duration

	name string
}

type BasicDebugger struct {
	core *Core

	methodNames map[*abc.Method]string
	methodStack []*abc.Method

	profiler  map[*abc.Method]*profileInfo
	enterTime time.Time
}

func NewBasicDebugger(core *Core) *BasicDebugger {
	d := &BasicDebugger{
		core:        core,
		methodNames: map[*abc.Method]string{},
		methodStack: []*abc.Method{},
		profiler:    map[*abc.Method]*profileInfo{},
	}
	return d
}

func (d *BasicDebugger) PrintProfilingInfo() {
	fmt.Println("Profiling info:", len(d.profiler))
	var total time.Duration
	top := make([]*profileInfo, 0, len(d.profiler))
	for method, profile := range d.profiler {
		if name, found := d.methodNames[method]; found {
			profile.name = name
		} else {
			profile.name = "???"
		}
		top = append(top, profile)
		total += profile.TotalTime
	}
	sort.Slice(top, func(i, j int) bool {
		return top[i].TotalTime > top[j].TotalTime
	})

	fmt.Println("Total time:", total)
	for i := 0; i < len(top) && i < 10; i++ {
		fmt.Printf("%2d: %-12s %6.2f%% %10d %q\n", i+1, top[i].TotalTime.String(), float64(top[i].TotalTime)/float64(total)*100.0, top[i].EnterCount, top[i].name)
	}
}

func (d *BasicDebugger) RegisterMethodName(m *abc.Method, name string) {
	d.methodNames[m] = name
}

func (d *BasicDebugger) SteppingMode(frame *ExecutionFrame) bool {
	return false
}

func (d *BasicDebugger) EnterMethod(frame *ExecutionFrame) bool {
	if len(d.methodStack) != 0 {
		passed := time.Since(d.enterTime)
		method := d.methodStack[len(d.methodStack)-1]
		profile := d.profiler[method]
		if profile == nil {
			profile = &profileInfo{}
			d.profiler[method] = profile
		}
		profile.TotalTime += passed
	}

	d.methodStack = append(d.methodStack, frame.Method)
	d.enterTime = time.Now()
	return false
}

func (d *BasicDebugger) LeaveMethod(frame *ExecutionFrame) {
	method := d.methodStack[len(d.methodStack)-1]
	if frame.Method != method {
		panic("nope")
	}
	passed := time.Since(d.enterTime)
	profile := d.profiler[method]
	if profile == nil {
		profile = &profileInfo{}
		d.profiler[method] = profile
	}
	profile.EnterCount++
	profile.TotalTime += passed
	d.methodStack = d.methodStack[:len(d.methodStack)-1]
	d.enterTime = time.Now()
}
