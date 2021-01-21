package flash

import (
	"github.com/dexter3k/dash/avm2"

	// "github.com/dexter3k/dash/flash/accessibility"
	// "github.com/dexter3k/dash/flash/automation"
	// "github.com/dexter3k/dash/flash/concurrent"
	// "github.com/dexter3k/dash/flash/crypto"
	// "github.com/dexter3k/dash/flash/debugger"
	// "github.com/dexter3k/dash/flash/desktop"
	"github.com/dexter3k/dash/flash/display"
	// "github.com/dexter3k/dash/flash/display3D"
	// "github.com/dexter3k/dash/flash/errors"
	"github.com/dexter3k/dash/flash/events"
	// "github.com/dexter3k/dash/flash/external"
	// "github.com/dexter3k/dash/flash/filters"
	// "github.com/dexter3k/dash/flash/geom"
	// "github.com/dexter3k/dash/flash/globalization"
	"github.com/dexter3k/dash/flash/media"
	// "github.com/dexter3k/dash/flash/net"
	// "github.com/dexter3k/dash/flash/printing"
	// "github.com/dexter3k/dash/flash/profiler"
	// "github.com/dexter3k/dash/flash/sampler"
	// "github.com/dexter3k/dash/flash/security"
	// "github.com/dexter3k/dash/flash/sensors"
	// "github.com/dexter3k/dash/flash/system"
	"github.com/dexter3k/dash/flash/text"
	// "github.com/dexter3k/dash/flash/trace"
	// "github.com/dexter3k/dash/flash/ui"
	// "github.com/dexter3k/dash/flash/utils"
	// "github.com/dexter3k/dash/flash/xml"
)

type Flash struct {
	// Accessibility accessibility.Accessibility
	// Automation    automation.Automation
	// Concurrent    concurrent.Concurrent
	// Crypto        crypto.Crypto
	// Debugger      debugger.Debugger
	// Desktop       desktop.Desktop
	Display       display.Display
	// Display3D     display3D.Display3D
	// Errors        errors.Errors
	Events        events.Events
	// External      external.External
	// Filters       filters.Filters
	// Geom          geom.Geom
	// Globalization globalization.Globalization
	Media         media.Media
	// Net           net.Net
	// Printing      printing.Printing
	// Profiler      profiler.Profiler
	// Sampler       sampler.Sampler
	// Security      security.Security
	// Sensors       sensors.Sensors
	// System        system.System
	Text          text.Text
	// Trace         trace.Trace
	// Ui            ui.Ui
	// Utils         utils.Utils
	// Xml           xml.Xml
}

func InitBuiltins(core *avm2.Core) (*Flash, *avm2.Script) {
	flash := &Flash{}
	script := core.CreateNativeScript()

	flash.Events.EventDispatcher = events.CreateClass_EventDispatcher(core)
	script.Global.Traits.AddTrait(core.NewClassTrait("EventDispatcher", "Pflash.events", flash.Events.EventDispatcher))

	flash.Display.DisplayObject = display.CreateClass_DisplayObject(core, &flash.Events.EventDispatcher.Object_Class)
	script.Global.Traits.AddTrait(core.NewClassTrait("DisplayObject", "Pflash.display", flash.Display.DisplayObject))

	flash.Display.Bitmap = display.CreateClass_Bitmap(core, &flash.Display.DisplayObject.Object_Class)
	script.Global.Traits.AddTrait(core.NewClassTrait("Bitmap", "Pflash.display", flash.Display.Bitmap))

	flash.Display.InteractiveObject = display.CreateClass_InteractiveObject(core, &flash.Display.DisplayObject.Object_Class)
	script.Global.Traits.AddTrait(core.NewClassTrait("InteractiveObject", "Pflash.display", flash.Display.InteractiveObject))

	flash.Display.DisplayObjectContainer = display.CreateClass_DisplayObjectContainer(core, &flash.Display.InteractiveObject.Object_Class)
	script.Global.Traits.AddTrait(core.NewClassTrait("DisplayObjectContainer", "Pflash.display", flash.Display.DisplayObjectContainer))

	flash.Display.Sprite = display.CreateClass_Sprite(core, &flash.Display.DisplayObjectContainer.Object_Class)
	script.Global.Traits.AddTrait(core.NewClassTrait("Sprite", "Pflash.display", flash.Display.Sprite))

	flash.Display.MovieClip = display.CreateClass_MovieClip(core, &flash.Display.Sprite.Object_Class)
	script.Global.Traits.AddTrait(core.NewClassTrait("MovieClip", "Pflash.display", flash.Display.MovieClip))

	flash.Text.TextField = text.CreateClass_TextField(core, &flash.Display.InteractiveObject.Object_Class)
	script.Global.Traits.AddTrait(core.NewClassTrait("TextField", "Pflash.text", flash.Text.TextField))

	return flash, script
}
