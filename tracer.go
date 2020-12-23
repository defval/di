package di

import "log"

var tracer Tracer = &nopTracer{}

// SetTracer sets global tracer.
func SetTracer(t Tracer) {
	tracer = t
}

// Tracer traces dependency injection cycle.
type Tracer interface {
	// Trace prints library logs.
	Trace(format string, args ...interface{})
}

// StdTracer traces dependency injection cycle to stdout.
type StdTracer struct {
}

func (s StdTracer) Trace(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// default nop tracer
type nopTracer struct {
}

func (n nopTracer) Trace(format string, args ...interface{}) {
}
