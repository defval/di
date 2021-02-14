package di

import (
	"reflect"
	"testing"
)

func Test_stacktrace(t *testing.T) {
	type args struct {
		skip int
	}
	tests := []struct {
		name      string
		args      args
		wantFrame callerFrame
	}{
		{
			name: "incorrect skip",
			args: args{skip: 10},
			wantFrame: callerFrame{
				function: "",
				file:     "",
				line:     0,
			},
		},
		{
			name: "incorrect skip",
			args: args{skip: 10},
			wantFrame: callerFrame{
				function: "",
				file:     "",
				line:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFrame := stacktrace(tt.args.skip); !reflect.DeepEqual(gotFrame, tt.wantFrame) {
				t.Errorf("stacktrace() = %v, want %v", gotFrame, tt.wantFrame)
			}
		})
	}
}
