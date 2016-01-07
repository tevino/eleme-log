package log

import "syscall"

type Termios syscall.Termios

const ioctlReadTermios = syscall.TIOCGETA
