// +build linux

package applog

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

// These constants are declared here, rather than importing
// them from the syscall package as some syscall packages, even
// on linux, for example gccgo, do not declare them.
const ioctlReadTermios = 0x5401  // syscall.TCGETS
const ioctlWriteTermios = 0x5402 // syscall.TCSETS

func isTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

// WriterIsTerminal returns true if the given io.Writer converts to
// an *os.File and the file's fd is a terminal.
func WriterIsTerminal(writer io.Writer) bool {
	file, ok := writer.(*os.File)
	return ok && isTerminal(file.Fd())
}

// IsTerminal returns true if the given file descriptor is a terminal.
// Swiped from golang.org/x/crypto/ssh/terminal
func IsTerminal(fd int) bool {
	return isTerminal(uintptr(fd))
}
