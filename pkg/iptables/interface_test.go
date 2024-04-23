package iptables

import "testing"

func Test_interfaceExists(t *testing.T) {
	type args struct {
		iface string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "should_exist",
			args: args{
				iface: "lo",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "should_not_exist",
			args: args{
				iface: "xyz0",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "should_exist",
			args: args{
				iface: "lo",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := interfaceExists(tt.args.iface)
			if (err != nil) != tt.wantErr {
				t.Errorf("interfaceExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("interfaceExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
