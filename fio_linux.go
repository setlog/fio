package fio

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/setlog/fio/fsi"
	"github.com/setlog/panik"
)

func openFile(filePath string, flag int, perm fs.FileMode) (fsi.File, error) {
	haveLock := false
	file, err := fsApi.OpenFile(filePath, flag, perm)
	if err != nil {
		return nil, err
	}
	defer func() {
		if !haveLock {
			file.Close()
		}
	}()
	err = lockForFlag(file.Fd(), flag)
	if err != nil {
		return nil, fmt.Errorf("open '%s': %w", filePath, err)
	}
	haveLock = true
	return file, nil
}

func readFile(filePath string) []byte {
	file, err := fsApi.OpenFile(filePath, os.O_RDONLY, 0660)
	panik.OnError(err)
	defer file.Close()
	if err := lockForFlag(file.Fd(), os.O_RDONLY); err != nil {
		panik.Panicf("open '%s': %w", filePath, err)
	}
	data, err := ioutil.ReadAll(file)
	panik.OnError(err)
	return data
}

func copyFile(fromFilePath, toFilePath string) (int64, error) {
	fileInfo, err := fsApi.Stat(fromFilePath)
	if err != nil {
		return 0, fmt.Errorf("copy '%s' to '%s': stat source: %w", fromFilePath, toFilePath, err)
	}
	src, err := openFile(fromFilePath, os.O_RDONLY, 0660)
	if err != nil {
		return 0, fmt.Errorf("copy '%s' to '%s': open source: %w", fromFilePath, toFilePath, err)
	}
	defer src.Close()
	var n int64
	if n, err = writeFile(toFilePath, src, fileInfo.Mode().Perm()); err != nil {
		return n, fmt.Errorf("copy '%s' to '%s': open destination: %w", fromFilePath, toFilePath, err)
	}
	return n, nil
}

func moveFile(fromFilePath, toFilePath string) (int64, error) {
	fileInfo, err := fsApi.Stat(fromFilePath)
	if err != nil {
		return 0, fmt.Errorf("move '%s' to '%s': stat source: %w", fromFilePath, toFilePath, err)
	}
	src, err := openFile(fromFilePath, os.O_RDONLY, 0660)
	if err != nil {
		return 0, fmt.Errorf("move '%s' to '%s': open source: %w", fromFilePath, toFilePath, err)
	}
	defer src.Close()
	var n int64
	if n, err = writeFile(toFilePath, src, fileInfo.Mode().Perm()); err != nil {
		return n, fmt.Errorf("move '%s' to '%s': open destination: %w", fromFilePath, toFilePath, err)
	}
	if err = fsApi.Remove(fromFilePath); err != nil {
		return n, fmt.Errorf("move '%s' to '%s': Remove source: %w", fromFilePath, toFilePath, err)
	}
	return n, nil
}

func writeFile(filePath string, reader io.Reader, perm fs.FileMode) (n int64, retErr error) {
	finishedWriting := false
	dst, err := openFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, perm)
	if err != nil {
		return 0, err
	}
	defer func() {
		dst.Close()
		if !finishedWriting {
			if remErr := fsApi.Remove(filePath); remErr != nil && !os.IsNotExist(remErr) {
				if err != nil {
					retErr = fmt.Errorf("%w. Then: %v", err, remErr)
				} else {
					retErr = remErr
				}
			}
		}
	}()
	n, err = io.Copy(dst, reader)
	if err != nil {
		return n, err
	}
	finishedWriting = true
	return n, nil
}

func lockForFlag(fd uintptr, flag int) (err error) {
	const mask = os.O_RDONLY | os.O_WRONLY | os.O_RDWR
	accessMode := flag & mask
	if accessMode == os.O_RDWR || accessMode == os.O_WRONLY {
		err = fsApi.FcntlFlock(fd, syscall.F_SETLK, wrLock())
		if err != nil {
			err = fmt.Errorf("acquire write-lock: %w", err)
		}
		return err
	} else if accessMode == os.O_RDONLY {
		err = fsApi.FcntlFlock(fd, syscall.F_SETLK, rdLock())
		if err != nil {
			err = fmt.Errorf("acquire read-lock: %w", err)
		}
		return err
	}
	return fmt.Errorf("acquire lock: bad access mode %d for flag %d", accessMode, flag)
}

func rdLock() *syscall.Flock_t {
	return lockWithType(syscall.F_RDLCK)
}

func wrLock() *syscall.Flock_t {
	return lockWithType(syscall.F_WRLCK)
}

// func unlock() *syscall.Flock_t {
// 	return lockWithType(syscall.F_UNLCK)
// }

func lockWithType(typ int16) *syscall.Flock_t {
	return &syscall.Flock_t{
		Type:   typ,
		Whence: io.SeekStart,
		Start:  0,
		Len:    0,
	}
}
