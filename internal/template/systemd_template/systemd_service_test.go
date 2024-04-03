package systemd_template

import (
	_ "embed"
	"testing"
)

func Test_serviceTypeAllowed(t *testing.T) {
	type args struct {
		st string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "service_type_not_allowed",
			args: args{
				st: "none",
			},
			want: false,
		},
		{
			name: "service_type_allowed",
			args: args{
				st: "fork",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serviceTypeAllowed(tt.args.st); got != tt.want {
				t.Errorf("serviceTypeAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}
