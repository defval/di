package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestId_String(t *testing.T) {
	i := id{
		Type: reflect.TypeOf(int(0)),
	}
	require.Equal(t, "int", i.String())
	i.Name = "foo"
	require.Equal(t, "int[foo]", i.String())
}
