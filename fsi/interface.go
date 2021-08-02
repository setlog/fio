package fsi

import (
	"io"
	"os"
	"syscall"
)

type File interface {
	io.ReadWriteCloser
	Fd() uintptr
	Name() string
}

type FileSystem interface {
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	FcntlFlock(fd uintptr, cmd int, lk *syscall.Flock_t) error
	Remove(name string) error
	Stat(name string) (os.FileInfo, error)
}
