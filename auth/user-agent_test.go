package auth

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserAgentAuthorizer_IsAuthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Add("User-Agent", "Simple Firefox User Agent")

	type fields struct {
		userAgents []string
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
			name:   "Find UserAgent",
			fields: fields{userAgents: []string{"firefox"}},
			args:   args{r: req},
			want:   true,
		},
		{
			name:   "Find no UserAgent",
			fields: fields{userAgents: []string{"chrome"}},
			args:   args{r: req},
			want:   false,
		},
		{
			name:   "Find one UserAgent",
			fields: fields{userAgents: []string{"edge", "firefox"}},
			args:   args{r: req},
			want:   true,
		},
		{
			name:   "Find without UserAgent header",
			fields: fields{userAgents: []string{"chrome", "firefox"}},
			args:   args{r: &http.Request{}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserAgentAuthorizer{
				userAgents: tt.fields.userAgents,
			}
			if got := u.IsAuthorized(tt.args.r); got != tt.want {
				t.Errorf("IsAuthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserAgentAuthorizer(t *testing.T) {
	type args struct {
		userAgents []string
	}
	tests := []struct {
		name string
		args args
		want *UserAgentAuthorizer
	}{
		{
			name: "Generate new user agent authorizer",
			args: args{userAgents: []string{"Firefox", "Chrome"}},
			want: &UserAgentAuthorizer{userAgents: []string{"firefox", "chrome"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserAgentAuthorizer(tt.args.userAgents); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserAgentAuthorizer() = %v, want %v", got, tt.want)
			}
		})
	}
}
