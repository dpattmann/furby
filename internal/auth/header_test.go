package auth

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHeaderAuthorizer_IsAuthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Add("X-Furby-Header", "value")

	type fields struct {
		headerName   string
		headerValues []string
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Find correct header value",
			fields: fields{
				headerName:   "X-Furby-Header",
				headerValues: []string{"value"},
			},
			args: args{r: req},
			want: true,
		},
		{
			name: "Find correct header value in slice",
			fields: fields{
				headerName:   "X-Furby-Header",
				headerValues: []string{"someValue", "value"},
			},
			args: args{r: req},
			want: true,
		},
		{
			name: "Header name case insensitive",
			fields: fields{
				headerName:   "x-furby-header",
				headerValues: []string{"value"},
			},
			args: args{r: req},
			want: true,
		},
		{
			name: "Header not found",
			fields: fields{
				headerName:   "X-Not-Furby-Header",
				headerValues: []string{"value"},
			},
			args: args{r: req},
			want: false,
		},
		{
			name: "Value not found",
			fields: fields{
				headerName:   "X-Furby-Header",
				headerValues: []string{"noValue"},
			},
			args: args{r: req},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HeaderAuthorizer{
				headerName:   tt.fields.headerName,
				headerValues: tt.fields.headerValues,
			}
			if got := h.IsAuthorized(tt.args.r); got != tt.want {
				t.Errorf("IsAuthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHeaderAuthorizer(t *testing.T) {
	type args struct {
		headerName   string
		headerValues []string
	}
	tests := []struct {
		name string
		args args
		want *HeaderAuthorizer
	}{
		{
			name: "Foo",
			args: args{
				headerName:   "X-Furby-Header",
				headerValues: []string{"value"},
			},
			want: NewHeaderAuthorizer("X-Furby-Header", []string{"value"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHeaderAuthorizer(tt.args.headerName, tt.args.headerValues); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHeaderAuthorizer() = %v, want %v", got, tt.want)
			}
		})
	}
}
