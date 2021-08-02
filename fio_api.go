package fio

import (
	"bytes"
	"io"
	"io/fs"
	"os"

	"github.com/setlog/panik"
)

// OpenFile opens the file at filePath in the same manner as os.OpenFile, but also
// claims an advisory lock matching your access flags (r/w/rw) which will be released
// when closing the file.
//
// Note that opening a file and getting an advisory lock are not (and cannot be) an atomic operation.
//
// Unlike the other functions in this package, this function will never log, because
// it will not know when you close the file descriptor.
//
// Errors result in panics created with panik.
func OpenFile(filePath string, flag int, perm fs.FileMode) *os.File {
	file, err := fsApi.OpenFile(filePath, flag, perm)
	panik.OnError(err)
	return file.(*os.File)
}

// ReadFile opens the file at filePath, claims an advisory read lock, reads all
// of its contents, closes the file, logs on success and returns the read contents.
//
// Note that opening a file and getting an advisory lock are not (and cannot be) an atomic operation.
//
// Errors result in panics created with panik.
func ReadFile(filePath string) []byte {
	data := readFile(filePath)
	if log := logger(); log != nil {
		log.Printf("Read '%s'.", filePath)
	}
	return data
}

// MoveFile creates a file at toFilePath, truncating it if it already exists,
// writes to it all data read from the file at fromFilePath, removes the file
// at fromFilePath, logs on success and returns the amount of bytes moved.
//
// Explicitly creating the target file effectively allows for it to be moved between mounts,
// which is not possible when using os.Rename().
//
// For the operation, an advisory read lock is claimed for the file at fromFilePath
// and an advisory write lock is claimed for the file at toFilePath.
//
// Note that opening a file and getting an advisory lock are not (and cannot be) an atomic operation.
//
// Errors result in panics created with panik.
func MoveFile(fromFilePath, toFilePath string) int64 {
	n, err := moveFile(fromFilePath, toFilePath)
	if err != nil {
		panik.Panicf("move '%s' to '%s': %w", fromFilePath, toFilePath, err)
	}
	if log := logger(); log != nil {
		log.Printf("Moved '%s' to '%s'.", fromFilePath, toFilePath)
	}
	return n
}

// RenameFile is a shorthand for panik.OnError(os.Rename(fromFilePath, toFilePath)).
//
// If fromFilePath and toFilePath are on different mounts, consider using MoveFile() instead.
func RenameFile(fromFilePath, toFilePath string) {
	panik.OnError(os.Rename(fromFilePath, toFilePath))
}

// CopyFile creates a file at toFilePath, truncating it if it already exists,
// writes to it all data read from the file at fromFilePath,  logs on success
// and returns the amount of bytes copied.
//
// For the operation, an advisory read lock is claimed for the file at fromFilePath
// and an advisory write lock is claimed for the file at toFilePath.
//
// Note that opening a file and getting an advisory lock are not (and cannot be) an atomic operation.
//
// Errors result in panics created with panik.
func CopyFile(fromFilePath, toFilePath string) int64 {
	n, err := copyFile(fromFilePath, toFilePath)
	panik.OnError(err)
	if log := logger(); log != nil {
		log.Printf("Copied '%s' to '%s'.", fromFilePath, toFilePath)
	}
	return n
}

// WriteFile creates a file at filePath, truncating it if it already exists,
// writes data to it and logs on success.
//
// For the operation, an advisory write lock is claimed for the file.
//
// Note that opening a file and getting an advisory lock are not (and cannot be) an atomic operation.
//
// Errors result in panics created with panik.
func WriteFile(filePath string, data []byte) {
	_, err := writeFile(filePath, bytes.NewReader(data), 0660)
	panik.OnError(err)
	if log := logger(); log != nil {
		log.Printf("Wrote '%s'.", filePath)
	}
}

// WriteFilePerm creates a file with permissions perm at filePath, truncating it
// if it already exists, writes data to it and logs on success.
//
// For the operation, an advisory write lock is claimed for the file.
//
// Note that opening a file and getting an advisory lock are not (and cannot be) an atomic operation.
//
// Errors result in panics created with panik.
func WriteFilePerm(filePath string, data []byte, perm fs.FileMode) {
	_, err := writeFile(filePath, bytes.NewReader(data), perm)
	panik.OnError(err)
	if log := logger(); log != nil {
		log.Printf("Wrote '%s'.", filePath)
	}
}

// WriteFileWithReader creates a file at filePath, truncating it if it already exists,
// writes to it all data read from reader, logs on success and returns the amount of bytes written.
//
// For the operation, an advisory write lock is claimed for the file.
//
// Note that opening a file and getting an advisory lock are not (and cannot be) an atomic operation.
//
// Errors result in panics created with panik.
func WriteFileWithReader(filePath string, reader io.Reader) int64 {
	n, err := writeFile(filePath, reader, 0660)
	panik.OnError(err)
	if log := logger(); log != nil {
		log.Printf("Wrote '%s'.", filePath)
	}
	return n
}

// WriteFileWithReaderPerm creates a file with permissions perm at filePath, truncating it
// if it already exists, writes to it all data read from reader, logs on success and returns
// the amount of bytes written.
//
// For the operation, an advisory write lock is claimed for the file.
//
// Note that opening a file and getting an advisory lock are not (and cannot be) an atomic operation.
//
// Errors result in panics created with panik.
func WriteFileWithReaderPerm(filePath string, reader io.Reader, perm os.FileMode) int64 {
	n, err := writeFile(filePath, reader, perm)
	panik.OnError(err)
	if log := logger(); log != nil {
		log.Printf("Wrote '%s'.", filePath)
	}
	return n
}

// RemoveFile removes the file at filePath if it exists, logs this
// and returns true on success. Returns false if the file did not exist.
//
// Errors result in panics created with panik.
func RemoveFile(filePath string) bool {
	err := os.Remove(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			panik.OnError(err)
		}
		return false
	}
	if log := logger(); log != nil {
		log.Printf("Removed '%s'.", filePath)
	}
	return true
}
