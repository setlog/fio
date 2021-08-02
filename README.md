## fio (File Input/Output) ![](https://github.com/setlog/fio/workflows/Tests/badge.svg)

This is a Linux-only package of convenience file operation functions which utilize advisory file locks for interaction with applications which support them; e.g. most FTP servers. These functions do not return `error` values; they use [panik](https://github.com/setlog/panik#the-problem).

See `fio_api.go` for available functions.

### Development

The following needs to be run before working on tests locally:

```bash
go install github.com/golang/mock/mockgen@v1.6.0
go generate
```
