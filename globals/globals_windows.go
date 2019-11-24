package globals

import "github.com/richardwilkes/win32"

var (
	ModuleInstance win32.HINSTANCE
	UIThreadID     win32.DWORD
	InvokeMsgID    uint32
)

func osInitialize() {
	ModuleInstance = win32.GetModuleHandleS("")
	UIThreadID = win32.GetCurrentThreadID()
	InvokeMsgID = win32.RegisterWindowMessageS("invoke")
}
