package fio

//go:generate mockgen -source=fsi/interface.go -destination=mock/fio_gen_mock.go -package=mock

import (
	"bytes"
	"os"
	"syscall"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/setlog/fio/mock"
)

const testSourceFileName = "foo"
const testDestinationFileName = "bar"
const testData = "Hello World"

var nextFd uintptr = 1

func TestOpenFileRead(t *testing.T) {
	ctrl, fsMock := prepareFileSystemMock(t)
	fileMock := mock.NewMockFile(ctrl)

	expectOpen(fsMock, fileMock, testSourceFileName, os.O_RDONLY)

	openFile(testSourceFileName, os.O_RDONLY, 0660)
}

func TestOpenFileWrite(t *testing.T) {
	ctrl, fsMock := prepareFileSystemMock(t)
	fileMock := mock.NewMockFile(ctrl)

	expectOpen(fsMock, fileMock, testSourceFileName, os.O_WRONLY)

	openFile(testSourceFileName, os.O_WRONLY, 0660)
}

func TestReadFile(t *testing.T) {
	ctrl, fsMock := prepareFileSystemMock(t)
	fileMock := mock.NewMockFile(ctrl)

	openCall := expectOpen(fsMock, fileMock, testSourceFileName, os.O_RDONLY)
	readCall := expectRead(fsMock, fileMock, []byte(testData)).After(openCall)
	fileMock.EXPECT().Close().Times(1).After(readCall)

	readFile(testSourceFileName)
}

func TestCopyFile(t *testing.T) {
	ctrl, fsMock := prepareFileSystemMock(t)
	srcFileMock := mock.NewMockFile(ctrl)
	dstFileMock := mock.NewMockFile(ctrl)

	statCall := expectStat(fsMock, srcFileMock)
	srcOpenCall := expectOpen(fsMock, srcFileMock, testSourceFileName, os.O_RDONLY).After(statCall)
	dstOpenCall := expectOpen(fsMock, dstFileMock, testDestinationFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE).After(srcOpenCall)
	readCall := expectRead(fsMock, srcFileMock, []byte(testData)).After(dstOpenCall)
	writeCall := expectWrite(fsMock, dstFileMock, []byte(testData)).After(dstOpenCall)
	dstCloseCall := dstFileMock.EXPECT().Close().Times(1).After(writeCall).After(readCall)
	srcFileMock.EXPECT().Close().Times(1).After(dstCloseCall)

	copyFile(testSourceFileName, testDestinationFileName)
}

func TestMoveFile(t *testing.T) {
	ctrl, fsMock := prepareFileSystemMock(t)
	srcFileMock := mock.NewMockFile(ctrl)
	dstFileMock := mock.NewMockFile(ctrl)

	statCall := expectStat(fsMock, srcFileMock)
	srcOpenCall := expectOpen(fsMock, srcFileMock, testSourceFileName, os.O_RDONLY).After(statCall)
	dstOpenCall := expectOpen(fsMock, dstFileMock, testDestinationFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE).After(srcOpenCall)
	readCall := expectRead(fsMock, srcFileMock, []byte(testData)).After(dstOpenCall)
	writeCall := expectWrite(fsMock, dstFileMock, []byte(testData)).After(dstOpenCall)
	dstCloseCall := dstFileMock.EXPECT().Close().Times(1).After(writeCall).After(readCall)
	srcRemoveCall := fsMock.EXPECT().Remove(testSourceFileName).After(dstCloseCall)
	srcFileMock.EXPECT().Close().Times(1).After(srcRemoveCall)

	moveFile(testSourceFileName, testDestinationFileName)
}

func prepareFileSystemMock(t *testing.T) (*gomock.Controller, *mock.MockFileSystem) {
	ctrl := gomock.NewController(t)
	fsMock := mock.NewMockFileSystem(ctrl)
	fsApi = fsMock
	return ctrl, fsMock
}

func expectStat(fsMock *mock.MockFileSystem, fileMock *mock.MockFile) *gomock.Call {
	return fsMock.EXPECT().Stat(testSourceFileName).Times(1).Return(&fileInfoImpl{
		name: testSourceFileName, size: int64(len(testData)), mode: 0660,
	}, nil)
}

func expectOpen(fsMock *mock.MockFileSystem, fileMock *mock.MockFile, name string, flag int) *gomock.Call {
	fd := nextFd
	nextFd++
	openCall := fsMock.EXPECT().OpenFile(name, flag, os.FileMode(0660)).Times(1).Return(fileMock, nil)
	fdCall := fileMock.EXPECT().Fd().Return(fd).After(openCall)
	var lk *syscall.Flock_t
	if (flag & (os.O_RDONLY | os.O_WRONLY | os.O_RDWR)) == os.O_RDONLY {
		lk = rdLock()
	} else {
		lk = wrLock()
	}
	return fsMock.EXPECT().FcntlFlock(fd, syscall.F_SETLK, gomock.Eq(lk)).Times(1).Return(nil).After(fdCall)
}

func expectRead(fsMock *mock.MockFileSystem, fileMock *mock.MockFile, data []byte) *gomock.Call {
	textBuffer := bytes.NewBuffer([]byte(testData))
	return fileMock.EXPECT().Read(gomock.Any()).MinTimes(1).DoAndReturn(func(p []byte) (int, error) {
		return textBuffer.Read(p)
	})
}

func expectWrite(fsMock *mock.MockFileSystem, fileMock *mock.MockFile, data []byte) *gomock.Call {
	return fileMock.EXPECT().Write(gomock.Eq(data)).MinTimes(1).DoAndReturn(func(p []byte) (int, error) {
		return len(p), nil
	})
}
