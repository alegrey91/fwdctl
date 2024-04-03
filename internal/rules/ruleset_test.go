package rules

import (
	"reflect"
	"strings"
	"testing"

	"github.com/alegrey91/fwdctl/pkg/iptables"
)

func TestDiff(t *testing.T) {
	type args struct {
		newRS string
		oldRS string
	}
	tests := []struct {
		name string
		args args
		want *RuleSetDiff
	}{
		// TODO: Add test cases.
		{
			name: "should_return_empty_RuleSetDiff",
			args: args{
				newRS: "",
				oldRS: "",
			},
			want: &RuleSetDiff{},
		},
		{
			name: "should_return_two_Rules_to_be_added",
			args: args{
				newRS: `
rules:
- dport: 3000
  saddr: 127.0.0.1
  sport: 80
  iface: lo
  proto: tcp
- dport: 2000
  saddr: 127.0.0.1
  sport: 22 
  iface: lo
  proto: tcp
`,
				oldRS: "",
			},
			want: &RuleSetDiff{
				ToAdd: []*iptables.Rule{
					{
						Iface: "lo",
						Proto: "tcp",
						Dport: 3000,
						Saddr: "127.0.0.1",
						Sport: 80,
					},
					{
						Iface: "lo",
						Proto: "tcp",
						Dport: 2000,
						Saddr: "127.0.0.1",
						Sport: 22,
					},
				},
			},
		},
		{
			name: "should_return_two_Rules_to_be_removed",
			args: args{
				newRS: "",
				oldRS: `
rules:
- dport: 3000
  saddr: 127.0.0.1
  sport: 80
  iface: lo
  proto: tcp
- dport: 2000
  saddr: 127.0.0.1
  sport: 22 
  iface: lo
  proto: tcp
`,
			},
			want: &RuleSetDiff{
				ToRemove: []*iptables.Rule{
					{
						Iface: "lo",
						Proto: "tcp",
						Dport: 3000,
						Saddr: "127.0.0.1",
						Sport: 80,
					},
					{
						Iface: "lo",
						Proto: "tcp",
						Dport: 2000,
						Saddr: "127.0.0.1",
						Sport: 22,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newRuleSet, err := NewRuleSetFromFile(strings.NewReader(tt.args.newRS))
			if err != nil {
				t.Error(err)
			}
			oldRuleSet, err := NewRuleSetFromFile(strings.NewReader(tt.args.oldRS))
			if err != nil {
				t.Error(err)
			}
			if got := Diff(oldRuleSet, newRuleSet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}
