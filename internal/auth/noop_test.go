package auth

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewNoOpAuthorizer(t *testing.T) {
	tests := []struct {
		name string
		want *NoOp
	}{
		{
			"Empty NoOp Authorizer",
			&NoOp{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNoOpAuthorizer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNoOpAuthorizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoOp_IsAuthorized(t *testing.T) {
	type args struct {
		in0 *http.Request
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Return always true",
			args{in0: nil},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NoOp{}
			if got := a.IsAuthorized(tt.args.in0); got != tt.want {
				t.Errorf("IsAuthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}
