// Package hooks runs the optional commands associated with a rule's
// create/delete lifecycle events. Commands are executed through "sh -c",
// so they support arguments, pipes and shell expansion. No fwdctl-specific
// environment is injected: the hook runs with fwdctl's own environment.
package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alegrey91/fwdctl/pkg/iptables"
)

// Run executes command through "sh -c", inheriting fwdctl's stdout/stderr.
// An empty command is a no-op and returns nil. A non-nil error is returned
// only when the command was run and exited non-zero.
func Run(command string) error {
	if strings.TrimSpace(command) == "" {
		return nil
	}
	// #nosec G204 -- hooks come from the root-owned rules file and are
	// intentionally executed as a shell command with arguments.
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("hook %q failed: %w", command, err)
	}
	return nil
}

// RunOnCreate runs the rule's on-create hook, if any.
func RunOnCreate(rule *iptables.Rule) error {
	if rule == nil || rule.Hooks == nil {
		return nil
	}
	return Run(rule.Hooks.OnCreate)
}

// RunOnDelete runs the rule's on-delete hook, if any.
func RunOnDelete(rule *iptables.Rule) error {
	if rule == nil || rule.Hooks == nil {
		return nil
	}
	return Run(rule.Hooks.OnDelete)
}
