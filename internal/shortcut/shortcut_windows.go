// Package shortcut provides global keyboard shortcuts without CGO.
// This implementation is specific to Windows using the user32 API.
package shortcut

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")
	procGetMessage       = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessage  = user32.NewProc("DispatchMessageW")
)

const (
	wmHotkey  = 0x0312
	modCtrl   = 0x0002
	modShift  = 0x0004
	vkKeyS    = 0x53
	hotkeyID  = 1
)

// MSG represents the Windows MSG structure.
type msg struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

// Register registers a global Ctrl+Shift+S hotkey and invokes callback on activation.
func Register(callback func()) error {
	r, _, err := procRegisterHotKey.Call(0, uintptr(hotkeyID), uintptr(modCtrl|modShift), uintptr(vkKeyS))
	if r == 0 {
		return fmt.Errorf("RegisterHotKey failed: %w", err)
	}

	go func() {
		var m msg
		for {
			r, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
			if int32(r) == -1 {
				break
			}
			if m.Message == wmHotkey && m.WParam == uintptr(hotkeyID) {
				callback()
			}
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
			procDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
		}
	}()

	return nil
}

// Unregister removes the registered global hotkey.
func Unregister() {
	procUnregisterHotKey.Call(0, uintptr(hotkeyID))
}
