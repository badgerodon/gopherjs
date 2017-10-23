// +build js

package syscall

import (
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

var warningPrinted = false
var lineBuffer []byte

func init() {
	js.Global.Set("$flushConsole", js.InternalObject(func() {
		if len(lineBuffer) != 0 {
			js.Global.Get("console").Call("log", string(lineBuffer))
			lineBuffer = nil
		}
	}))
}

func printWarning() {
	if !warningPrinted {
		js.Global.Get("console").Call("error", "warning: system calls not available, see https://github.com/gopherjs/gopherjs/blob/master/doc/syscalls.md")
	}
	warningPrinted = true
}

func printToConsole(b []byte) {
	goPrintToConsole := js.Global.Get("goPrintToConsole")
	if goPrintToConsole != js.Undefined {
		goPrintToConsole.Invoke(js.InternalObject(b))
		return
	}

	lineBuffer = append(lineBuffer, b...)
	for {
		i := indexByte(lineBuffer, '\n')
		if i == -1 {
			break
		}
		js.Global.Get("console").Call("log", string(lineBuffer[:i])) // don't use println, since it does not externalize multibyte characters
		lineBuffer = lineBuffer[i+1:]
	}
}

func use(p unsafe.Pointer) {
	// no-op
}

// indexByte is copied from bytes package to avoid importing it (since the real syscall package doesn't).
func indexByte(s []byte, c byte) int {
	for i, b := range s {
		if b == c {
			return i
		}
	}
	return -1
}

func StartProcess(argv0 string, argv []string, attr *ProcAttr) (pid int, handle uintptr, err error) {
	return DefaultStartProcessFunction(argv0, argv, attr)
}

// default write and read functions to allow libraries to overwrite the behavior
var (
	DefaultReadFunction = func(fd uintptr, data []byte) (int, error) {
		printWarning()
		return -1, EACCES
	}
	DefaultStartProcessFunction = func(argv0 string, argv []string, attr *ProcAttr) (pid int, handle uintptr, err error) {
		panic("starting processes is not supported in gopherjs")
	}
)

func Open(path string, mode int, perm uint32) (fd int, err error) { 
	return DefaultOpenFunction(path, mode, perm)
}

var (
	DefaultOpenFunction = func(path string, mode int, perm uint32) (fd int, err error) {
		printWarning()
		return -1, EACCES
	}
)

func CloseOnExec(fd int) {
	// do nothing
}


func Write(fd int, p []byte) (n int, err error) {
	return DefaultWriteFunction(fd, p)
}

var DefaultWriteFunction = func(fd int, p []byte) (int, error) {
	if fd == 1 || fd == 2 {
		printToConsole(p)
		return len(p), nil
	}
	printWarning()
	return -1, EACCES
}

func Close(fd int) (err error) {
	return DefaultCloseFunction(fd)
}

var DefaultCloseFunction = func(fd int) (err error) {
	printWarning()
	return EACCES
}

func fcntl(fd int, cmd int, arg int) (val int, err error) {
	return DefaultFCNTLFunction(fd, cmd, arg)
}

var DefaultFCNTLFunction = func(fd int, cmd int, arg int) (val int, err error) {
	printWarning()
	return -1, EACCES
}