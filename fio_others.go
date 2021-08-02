//+build !linux

package fio

import (
	"io"
	"io/fs"
	"os"
)

// This file only exists so developers on non-Linux operating systems can work
// on the code with gopls while it does not yet support working on multi-platform projects.

const errorMessage = "this is only implemented for Linux"

func openFile(filePath string, flag int, perm gofs.FileMode) (fs.File, error) {
	panic(errorMessage)
}

func readFile(filePath string) []byte {
	panic(errorMessage)
}

func copyFile(fromFilePath, toFilePath string) (int64, error) {
	panic(errorMessage)
}

func moveFile(fromFilePath, toFilePath string) (int64, error) {
	panic(errorMessage)
}

func writeFile(filePath string, reader io.Reader, perm gofs.FileMode) (n int64, retErr error) {
	panic(errorMessage)
}

func lockForFlag(fd uintptr, flag int) (err error) {
	panic(errorMessage)
}
