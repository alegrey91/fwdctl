package iptables

import (
	"fmt"
	"testing"
)

func TestValidateForward(t *testing.T) {
	testCases := []struct {
		id             int
		iface          string
		proto          string
		dport          int
		saddr          string
		sport          int
		expectedResult bool
	}{
		{
			1,
			"lo",
			"tcp",
			9090,
			"127.0.0.1",
			80,
			true,
		},
		{
			2,
			"lo",
			"tcps",
			9090,
			"127.0.0.1",
			80,
			false,
		},
		{
			3,
			"lo",
			"tcp",
			10202020,
			"127.0.0.1",
			80,
			false,
		},
		{
			4,
			"lo",
			"tcp",
			9090,
			"127.0.0.1",
			800000000,
			false,
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("Checking rule with id %d", tt.id), func(t *testing.T) {
			testResult, testErr := ValidateForward(tt.iface, tt.proto, tt.dport, tt.saddr, tt.sport)
			if testErr != nil {
				t.Logf("%v", testErr)
			}
			if testResult != tt.expectedResult {
				t.Fatal("Test failed")
			}
		})
	}
}
