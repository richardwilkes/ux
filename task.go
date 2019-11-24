package ux

import (
	"time"

	"github.com/richardwilkes/toolbox/errs"
)

// Invoke a task on the UI thread. The task is put into the system event queue
// and will be run at the next opportunity.
func Invoke(taskFunction func()) {
	osInvokeUITask(taskFunction)
}

// InvokeAfter schedules a task to be run on the UI thread after waiting for
// the specified duration.
func InvokeAfter(taskFunction func(), after time.Duration) {
	time.AfterFunc(after, func() { osInvokeUITask(taskFunction) })
}

// SetInvokeRecoverCallback sets a callback that will be called should an
// invoked task panic. If no recover callback is set, the panic will be
// silently swallowed.
func SetInvokeRecoverCallback(recoveryHandler errs.RecoveryHandler) {
	osSetInvokeRecoverCallback(recoveryHandler)
}
