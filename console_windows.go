//go:build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"
)

func initConsole() {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	if err := windows.GetConsoleMode(stdout, &originalMode); err == nil {
		mode := originalMode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		windows.SetConsoleMode(stdout, mode)
	}
}
