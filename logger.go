package di

// Logger logs internal container actions. By default it omits logs.
// You can set logger by Option LogFunc(). Also, you can provide you own logger to
// container as di.Logger interface. Then container use it for internal logs.
type Logger interface {
	Logf(format string, values ...interface{})
}

// is a default logger that discard logs
type nopLogger struct {
}

func (n nopLogger) Logf(_ string, _ ...interface{}) {}
