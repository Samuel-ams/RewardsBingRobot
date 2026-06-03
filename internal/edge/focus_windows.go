//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                   = windows.NewLazySystemDLL("user32.dll")
	user32spi                = windows.NewLazySystemDLL("user32.dll")
	procEnumWindows          = user32.NewProc("EnumWindows")
	procGetWindowThreadPid   = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible      = user32.NewProc("IsWindowVisible")
	procGetClassName         = user32.NewProc("GetClassNameW")
	procGetForegroundWindow  = user32.NewProc("GetForegroundWindow")
	procAttachThreadInput    = user32.NewProc("AttachThreadInput")
	procSetForegroundWindow  = user32.NewProc("SetForegroundWindow")
	procShowWindow           = user32.NewProc("ShowWindow")
	procSetWindowPos         = user32.NewProc("SetWindowPos")
	procKeybdEvent           = user32.NewProc("keybd_event")
	procSystemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

const (
	swRestore           = 9
	hwndTopmost         = ^uintptr(0) // -1
	hwndNoTopmost       = ^uintptr(1) // -2
	swpNoSize           = 0x0001
	swpNoMove           = 0x0002
	vkMenu              = 0x12 // Alt key
	keybdEventKeyUp     = 0x0002
	spiSetFgLockTimeout = 0x2001
	spifSendChange      = 0x0002
)

// FocusPID brings the main visible Chromium window owned by the given PID
// to the OS foreground, bypassing Windows' foreground lock.
func FocusPID(pid uint32) {
	// Disable foreground lock timeout so SetForegroundWindow always works
	procSystemParametersInfo.Call(spiSetFgLockTimeout, 0, 0, spifSendChange)

	hwnd := findChromiumWindowByPID(pid)
	if hwnd == 0 {
		// fallback: any visible Chromium window
		hwnd = findChromiumWindow(0)
		if hwnd == 0 {
			return
		}
	}

	// The keybd_event(Alt) trick grants this process foreground permission
	procKeybdEvent.Call(vkMenu, 0, 0, 0)
	procKeybdEvent.Call(vkMenu, 0, keybdEventKeyUp, 0)

	fgHwnd, _, _ := procGetForegroundWindow.Call()
	var fgPid uint32
	fgTid, _, _ := procGetWindowThreadPid.Call(fgHwnd, uintptr(unsafe.Pointer(&fgPid)))
	curTid := uintptr(windows.GetCurrentThreadId())

	if fgTid != 0 && fgTid != curTid {
		procAttachThreadInput.Call(curTid, fgTid, 1)
		defer procAttachThreadInput.Call(curTid, fgTid, 0)
	}

	procShowWindow.Call(hwnd, swRestore)
	procSetWindowPos.Call(hwnd, hwndTopmost, 0, 0, 0, 0, swpNoSize|swpNoMove)
	procSetWindowPos.Call(hwnd, hwndNoTopmost, 0, 0, 0, 0, swpNoSize|swpNoMove)
	procSetForegroundWindow.Call(hwnd)
}

// findChromiumWindowByPID finds the main (largest) Chrome_WidgetWin_1 window for a given PID.
func findChromiumWindowByPID(pid uint32) uintptr {
	var found uintptr

	cb := syscall.NewCallback(func(hwnd, _ uintptr) uintptr {
		vis, _, _ := procIsWindowVisible.Call(hwnd)
		if vis == 0 {
			return 1
		}

		var winPid uint32
		procGetWindowThreadPid.Call(hwnd, uintptr(unsafe.Pointer(&winPid)))
		if winPid != pid {
			return 1
		}

		buf := make([]uint16, 256)
		procGetClassName.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), 256)
		if syscall.UTF16ToString(buf) == "Chrome_WidgetWin_1" {
			found = hwnd
			return 0 // stop on first match for this PID
		}
		return 1
	})

	procEnumWindows.Call(cb, 0)
	return found
}

// findChromiumWindow finds the first visible Chrome_WidgetWin_1 window (any PID).
func findChromiumWindow(_ uint32) uintptr {
	var found uintptr

	cb := syscall.NewCallback(func(hwnd, _ uintptr) uintptr {
		vis, _, _ := procIsWindowVisible.Call(hwnd)
		if vis == 0 {
			return 1
		}
		buf := make([]uint16, 256)
		procGetClassName.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), 256)
		if syscall.UTF16ToString(buf) == "Chrome_WidgetWin_1" {
			found = hwnd
			return 0
		}
		return 1
	})

	procEnumWindows.Call(cb, 0)
	return found
}
