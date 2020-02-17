package di

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetLogFunc(t *testing.T) {
	require.NotNil(t, _log)
	called := false
	logFn := func(format string, v ...interface{}) {
		called = true
	}
	SetLogFunc(logFn)
	_log("format %s", "value")
	require.True(t, called)
	require.PanicsWithValue(t, "log function should not be nil", func() {
		SetLogFunc(nil)
	})
}
