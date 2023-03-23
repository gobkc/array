package array

import (
	"reflect"
	"testing"
)

func TestIds(t *testing.T) {
	type args struct {
		dest any
	}
	type testCase[T interface{ int | int32 | int64 }] struct {
		name string
		args args
		want []T
	}
	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ids(tt.args.dest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ids() = %v, want %v", got, tt.want)
			}
		})
	}
}
