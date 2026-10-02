//go:build !windows

package main

func setConsoleTitle(title string) {
	// No-op on non-Windows platforms
}