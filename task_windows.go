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
	"sync"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/ux/globals"
	"github.com/richardwilkes/win32"
)

var (
	dispatchLock            sync.Mutex
	dispatchID              uint64 = 1
	dispatchMap                    = make(map[uint64]func())
	dispatchRecoverCallback errs.RecoveryHandler
)

func osInvokeUITask(f func()) {
	dispatchLock.Lock()
	id := dispatchID
	dispatchID++
	dispatchMap[id] = f
	dispatchLock.Unlock()
	if err := win32.PostThreadMessage(globals.UIThreadID, globals.InvokeMsgID, win32.WPARAM((id>>32)&0xFFFFFFFF), win32.LPARAM(id&0xFFFFFFFF)); err != nil {
		jot.Error(err)
	}
}

func osSetInvokeRecoverCallback(recoveryHandler errs.RecoveryHandler) {
	dispatchLock.Lock()
	dispatchRecoverCallback = recoveryHandler
	dispatchLock.Unlock()
}

func dispatchTask(id uint64) {
	dispatchLock.Lock()
	callback, ok := dispatchMap[id]
	if ok {
		delete(dispatchMap, id)
	}
	recoverCallback := dispatchRecoverCallback
	dispatchLock.Unlock()
	if callback != nil {
		defer errs.Recovery(recoverCallback)
		callback()
	}
}
