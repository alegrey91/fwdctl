package iptables

import (
	"reflect"
	"testing"
)

func TestExtractRuleInfo(t *testing.T) {
	tests := []struct {
		name    string
		rule    string
		want    *Rule
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "example_1",
			rule: "-A PREROUTING -i lo -p tcp -m tcp --dport 3001 -m comment --comment fwdctl -j DNAT --to-destination 127.0.0.1:80",
			want: &Rule{
				Iface: "lo",
				Proto: "tcp",
				Dport: 3001,
				Saddr: "127.0.0.1",
				Sport: 80,
			},
		},
		{
			name: "example_2",
			rule: "-A PREROUTING -i eth0 -p tcp -m tcp --dport 3001 -m comment --comment fwdctl -j DNAT --to-destination 127.0.0.1:80",
			want: &Rule{
				Iface: "eth0",
				Proto: "tcp",
				Dport: 3001,
				Saddr: "127.0.0.1",
				Sport: 80,
			},
		},
		{
			name:    "example_3",
			rule:    "-A PREROUTING -i eth0 -p tcp -m tcp --dport 3001 -m comment --comment fwdctl -j DNAT",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractRuleInfo(tt.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractRuleInfo() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractRuleInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRule_String(t *testing.T) {
	type fields struct {
		Iface   string
		Proto   string
		Dport   int
		Saddr   string
		Sport   int
		Comment string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
		{
			name: "should_return_consistent_value",
			fields: fields{
				Iface: "eth0",
				Proto: "tcp",
				Dport: 8080,
				Saddr: "127.0.0.1",
				Sport: 80,
			},
			want: []string{"-i", "eth0", "-p", "tcp", "-m", "tcp", "--dport", "8080", "-m", "comment", "--comment", "fwdctl", "-j", "DNAT", "--to-destination", "127.0.0.1:80"},
		},
		{
			name: "should_return_consistent_value",
			fields: fields{
				Iface: "lo",
				Proto: "tcp",
				Dport: 8080,
				Saddr: "127.0.0.1",
				Sport: 90,
			},
			want: []string{"-i", "lo", "-p", "tcp", "-m", "tcp", "--dport", "8080", "-m", "comment", "--comment", "fwdctl", "-j", "DNAT", "--to-destination", "127.0.0.1:90"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &Rule{
				Iface:   tt.fields.Iface,
				Proto:   tt.fields.Proto,
				Dport:   tt.fields.Dport,
				Saddr:   tt.fields.Saddr,
				Sport:   tt.fields.Sport,
				Comment: tt.fields.Comment,
			}
			if got := rule.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rule.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
