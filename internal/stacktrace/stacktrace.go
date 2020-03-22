package stacktrace

import (
	"fmt"
	"runtime"
	"strings"
)

// CallerFrame
func CallerFrame(skip int) (frame Frame) {
	pc, file, line, ok := runtime.Caller(skip + 2)
	if !ok {
		return Frame{}
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return Frame{}
	}
	return Frame{
		function: shortFuncName(f),
		file:     file,
		line:     line,
	}
}

// Frame
type Frame struct {
	function string
	file     string
	line     int
}

// Format
func (f Frame) Format(s fmt.State, c rune) {
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
