package di

import (
	"fmt"
	"runtime"
	"strings"
)

// stacktrace returns stacktrace call frame with skip.
func stacktrace(skip int) (frame callerFrame) {
	pc, file, line, ok := runtime.Caller(skip + 2)
	if !ok {
		return callerFrame{}
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return callerFrame{}
	}
	return callerFrame{
		function: shortFuncName(f),
		file:     file,
		line:     line,
	}
}

// callerFrame represents stacktrace frame.
type callerFrame struct {
	function string
	file     string
	line     int
}

// Format formats stacktrace frame.
func (f callerFrame) Format(s fmt.State, c rune) {
	_, _ = fmt.Fprintf(s, "%s:%d", f.file, f.line)
}

func shortFuncName(f *runtime.Func) string {
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}
