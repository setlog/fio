package fio

import (
	"io/fs"
	"os"
	"syscall"
	"time"

	"github.com/setlog/fio/fsi"
)

var fsApi fsi.FileSystem = &fileSystemImpl{}

type fileSystemImpl struct{}

func (fs *fileSystemImpl) OpenFile(name string, flag int, perm os.FileMode) (fsi.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (fs *fileSystemImpl) FcntlFlock(fd uintptr, cmd int, lk *syscall.Flock_t) error {
	return syscall.FcntlFlock(fd, cmd, lk)
}

func (fs *fileSystemImpl) Remove(name string) error {
	return os.Remove(name)
}

func (fs *fileSystemImpl) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

type fileInfoImpl struct {
	name string
	size int64
	mode fs.FileMode
}

func (fi *fileInfoImpl) Name() string {
	return fi.name
}

func (fi *fileInfoImpl) Size() int64 {
	return fi.size
}

func (fi *fileInfoImpl) Mode() fs.FileMode {
	return fi.mode
}

func (fi *fileInfoImpl) ModTime() time.Time {
	return time.Now()
}

func (fi *fileInfoImpl) IsDir() bool {
	return false
}

func (fi *fileInfoImpl) Sys() interface{} {
	return nil
}
