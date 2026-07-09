package rules

import (
	"strings"
	"testing"
)

// TestHooksParsed verifies that on-create / on-delete hooks are read
// from the rules file into the Rule struct.
func TestHooksParsed(t *testing.T) {
	file := `
rules:
- dport: 3000
  saddr: 127.0.0.1
  sport: 80
  iface: lo
  proto: tcp
  hooks:
    on-create: "/bin/on-create.sh --port 3000"
    on-delete: "/bin/on-delete.sh --port 3000"
`
	rs, err := NewRuleSetFromFile(strings.NewReader(file))
	if err != nil {
		t.Fatalf("NewRuleSetFromFile() error = %v", err)
	}
	if len(rs.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rs.Rules))
	}
	for _, rule := range rs.Rules {
		if rule.Hooks == nil {
			t.Fatal("expected Hooks to be set, got nil")
		}
		if got, want := rule.Hooks.OnCreate, "/bin/on-create.sh --port 3000"; got != want {
			t.Errorf("OnCreate = %q, want %q", got, want)
		}
		if got, want := rule.Hooks.OnDelete, "/bin/on-delete.sh --port 3000"; got != want {
			t.Errorf("OnDelete = %q, want %q", got, want)
		}
	}
}

// TestHooksExcludedFromIdentity verifies that two rules differing only by
// their hooks hash to the same key. This is what keeps the daemon from
// churning a forward (delete + recreate) when only a hook string changes.
func TestHooksExcludedFromIdentity(t *testing.T) {
	withHooks := `
rules:
- dport: 3000
  saddr: 127.0.0.1
  sport: 80
  iface: lo
  proto: tcp
  hooks:
    on-create: "/bin/a.sh"
    on-delete: "/bin/b.sh"
`
	withoutHooks := `
rules:
- dport: 3000
  saddr: 127.0.0.1
  sport: 80
  iface: lo
  proto: tcp
`
	differentHooks := `
rules:
- dport: 3000
  saddr: 127.0.0.1
  sport: 80
  iface: lo
  proto: tcp
  hooks:
    on-create: "/bin/completely-different.sh"
`
	rsA, err := NewRuleSetFromFile(strings.NewReader(withHooks))
	if err != nil {
		t.Fatal(err)
	}
	rsB, err := NewRuleSetFromFile(strings.NewReader(withoutHooks))
	if err != nil {
		t.Fatal(err)
	}
	rsC, err := NewRuleSetFromFile(strings.NewReader(differentHooks))
	if err != nil {
		t.Fatal(err)
	}

	keyOf := func(rs *RuleSet) string {
		for k := range rs.Rules {
			return k
		}
		return ""
	}
	kA, kB, kC := keyOf(rsA), keyOf(rsB), keyOf(rsC)
	if kA != kB {
		t.Errorf("hook presence changed identity: %q != %q", kA, kB)
	}
	if kA != kC {
		t.Errorf("hook content changed identity: %q != %q", kA, kC)
	}

	// A Diff between a ruleset with hooks and one without must be empty:
	// nothing to add, nothing to remove.
	diff := Diff(rsA, rsC)
	if len(diff.ToAdd) != 0 || len(diff.ToRemove) != 0 {
		t.Errorf("hook-only change produced a non-empty diff: ToAdd=%d ToRemove=%d",
			len(diff.ToAdd), len(diff.ToRemove))
	}
}
