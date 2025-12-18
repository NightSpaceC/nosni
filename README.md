# nosni
This is a HTTPS proxy to remove the ServerName extension from your TLS Client Hello.

In order to pretend to be a client who doesn't support this extension, I use [gomonkey](https://github.com/agiledragon/gomonkey) to hook the standard library of Go.

Although [gomonkey](https://github.com/agiledragon/gomonkey) support most of common platforms, depending on your hardware, it may be unusable.

You can't just use `go build` to build it because Go 1.23 doesn't allow to use `//go:linkname` to link the standard library in the Pull mode.

The correct command is
```
go build -ldflags="-checklinkname=0"
```
or
```
make
```