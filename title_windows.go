//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

func setConsoleTitle(title string) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleTitleProc := kernel32.NewProc("SetConsoleTitleW")
	titlePtr, err := syscall.UTF16PtrFromString(title)
	if err == nil {
		_, _, _ = setConsoleTitleProc.Call(uintptr(unsafe.Pointer(titlePtr)))
	}
}