package globals

import (
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/softref"
)

var (
	// Pool is the soft reference pool used by all soft references allocated
	// by the ux packages.
	Pool = softref.NewPool(&jot.Logger{})
)

// Initialize the globals.
func Initialize() {
	osInitialize()
}
