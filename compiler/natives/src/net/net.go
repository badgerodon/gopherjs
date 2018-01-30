// +build js

package net

import (
	"errors"
	"syscall"

	"github.com/gopherjs/gopherjs/js"
)

// extendable functions
var (
	ListenFunc = func(net, laddr string) (Listener, error) {
		panic(errors.New("network access is not supported by GopherJS"))
	}
	DialFunc = func(network, address string) (Conn, error) {
		panic(errors.New("network access is not supported by GopherJS"))
	}
)

func Listen(net, laddr string) (Listener, error) {
	return ListenFunc(net, laddr)
}

func (d *Dialer) Dial(network, address string) (Conn, error) {
	return DialFunc(network, address)
}

func sysInit() {
}

func probeIPv4Stack() bool {
	return false
}

func probeIPv6Stack() (supportsIPv6, supportsIPv4map bool) {
	return false, false
}

func probeWindowsIPStack() (supportsVistaIP bool) {
	return false
}

func maxListenerBacklog() int {
	return syscall.SOMAXCONN
}

// Copy of strings.IndexByte.
func byteIndex(s string, c byte) int {
	return js.InternalObject(s).Call("indexOf", js.Global.Get("String").Call("fromCharCode", c)).Int()
}

// Copy of bytes.Equal.
func bytesEqual(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}
	for i, b := range x {
		if b != y[i] {
			return false
		}
	}
	return true
}

// Copy of bytes.IndexByte.
func bytesIndexByte(s []byte, c byte) int {
	for i, b := range s {
		if b == c {
			return i
		}
	}
	return -1
}
