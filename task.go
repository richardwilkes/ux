// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

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
