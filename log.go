package di

// log
var _log = func(format string, v ...interface{}) {}

// SetLogFunc sets log function for debug purposes.
func SetLogFunc(fn func(format string, v ...interface{})) {
	if fn == nil {
		panic("log function should not be nil")
	}
	_log = fn
}
