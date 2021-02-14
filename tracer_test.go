package di

import "testing"

func TestNopTracer_Trace(t *testing.T) {
	tracer := nopTracer{}
	tracer.Trace("test")
}
