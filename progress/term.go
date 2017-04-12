package progress

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

type window struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func TerminalWidth() (int, error) {
	w := new(window)
	tio := syscall.TIOCGWINSZ
	if runtime.GOOS == "darwin" {
		tio = 0x40087468
	}
	tty, err := os.Open("dev/tty")
	if err != nil {
		tty = os.Stdin
	}
	res, _, err := syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), uintptr(tio), uintptr(unsafe.Pointer(w)))
	if int(res) == -1 {
		return 0, err
	}
	return int(w.Col), nil
}
