package iptables

import (
	"errors"
	"testing"
)

func Test_validateForward(t *testing.T) {
	tests := []struct {
		name          string
		iface         string
		proto         string
		dport         int
		saddr         string
		sport         int
		expectedError error
		wantErr       bool
	}{
		{
			"should_not_fail",
			"lo",
			"tcp",
			9090,
			"127.0.0.1",
			80,
			nil,
			false,
		},
		{
			"should_fail_protocol",
			"lo",
			"tcps",
			9090,
			"127.0.0.1",
			80,
			errors.New("protocol: 'tcps' protocol name not allowed"),
			true,
		},
		{
			"should_fail_destination_port",
			"lo",
			"tcp",
			10202020,
			"127.0.0.1",
			80,
			errors.New("destination port: '10202020' port number not allowed"),
			true,
		},
		{
			"shoudl_fail_source_port",
			"lo",
			"tcp",
			9090,
			"127.0.0.1",
			800000000,
			errors.New("source port: '800000000' port number not allowed"),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateForward(tt.iface, tt.proto, tt.dport, tt.saddr, tt.sport); (err != nil) != tt.wantErr {
				t.Errorf("validateForward() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
