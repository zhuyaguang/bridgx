package utils

import (
	"testing"
)

func TestSshCheck(t *testing.T) {
	type args struct {
		ip       string
		username string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ssh check test",
			args: args{
				ip:       "127.0.0.1",
				username: "root",
				password: "123456",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := SshCheck(tt.args.ip, tt.args.username, tt.args.password); res != tt.want {
				t.Errorf("SshCheck() error = %v, want %v", res, tt.want)
			}
		})
	}
}
