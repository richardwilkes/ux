package ux

import (
	"github.com/richardwilkes/macos/dispatch"
	"github.com/richardwilkes/toolbox/errs"
)

func osInvokeUITask(f func()) {
	dispatch.AsyncFunctionOnMainQueue(f)
}

func osSetInvokeRecoverCallback(recoveryHandler errs.RecoveryHandler) {
	dispatch.SetDispatchRecoverCallback(recoveryHandler)
}
