package hooks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alegrey91/fwdctl/pkg/iptables"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{name: "empty command is a no-op", command: "", wantErr: false},
		{name: "blank command is a no-op", command: "   ", wantErr: false},
		{name: "successful command", command: "true", wantErr: false},
		{name: "command with arguments", command: "echo hello world", wantErr: false},
		{name: "failing command", command: "exit 1", wantErr: true},
		{name: "non-existent command", command: "fwdctl_no_such_binary_xyz", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.command, nil); (err != nil) != tt.wantErr {
				t.Errorf("Run(%q) error = %v, wantErr %v", tt.command, err, tt.wantErr)
			}
		})
	}
}

func TestRunSideEffect(t *testing.T) {
	dir := t.TempDir()
	flag := filepath.Join(dir, "flag")
	if err := Run("echo created > "+flag, nil); err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}
	content, err := os.ReadFile(flag)
	if err != nil {
		t.Fatalf("reading flag file: %v", err)
	}
	if got := string(content); got != "created\n" {
		t.Errorf("hook side effect = %q, want %q", got, "created\n")
	}
}

func TestRunOnCreate(t *testing.T) {
	dir := t.TempDir()
	flag := filepath.Join(dir, "created")

	// nil Hooks must be a no-op.
	if err := RunOnCreate(&iptables.Rule{}); err != nil {
		t.Errorf("RunOnCreate() with nil Hooks error = %v, want nil", err)
	}
	// nil rule must be a no-op.
	if err := RunOnCreate(nil); err != nil {
		t.Errorf("RunOnCreate(nil) error = %v, want nil", err)
	}
	// on-create hook runs.
	rule := &iptables.Rule{Hooks: &iptables.Hooks{OnCreate: "touch " + flag}}
	if err := RunOnCreate(rule); err != nil {
		t.Fatalf("RunOnCreate() error = %v", err)
	}
	if _, err := os.Stat(flag); err != nil {
		t.Errorf("on-create hook did not run: %v", err)
	}
	// on-delete must not have run.
	if err := RunOnDelete(rule); err != nil {
		t.Errorf("RunOnDelete() with empty OnDelete error = %v, want nil", err)
	}
}

func TestRunOnDelete(t *testing.T) {
	dir := t.TempDir()
	flag := filepath.Join(dir, "deleted")

	if err := RunOnDelete(&iptables.Rule{}); err != nil {
		t.Errorf("RunOnDelete() with nil Hooks error = %v, want nil", err)
	}
	rule := &iptables.Rule{Hooks: &iptables.Hooks{OnDelete: "touch " + flag}}
	if err := RunOnDelete(rule); err != nil {
		t.Fatalf("RunOnDelete() error = %v", err)
	}
	if _, err := os.Stat(flag); err != nil {
		t.Errorf("on-delete hook did not run: %v", err)
	}
}

func TestRunOnCreateFailure(t *testing.T) {
	rule := &iptables.Rule{Hooks: &iptables.Hooks{OnCreate: "exit 1"}}
	if err := RunOnCreate(rule); err == nil {
		t.Error("RunOnCreate() with failing hook error = nil, want error")
	}
}

func TestRunOnCreateInjectsEnv(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env")
	rule := &iptables.Rule{
		Iface: "eth0",
		Proto: "tcp",
		Dport: 9090,
		Saddr: "127.0.0.1",
		Sport: 8080,
		Hooks: &iptables.Hooks{
			OnCreate: "printf '%s\\n' \"$FWDCTL_IFACE\" \"$FWDCTL_PROTO\" \"$FWDCTL_DPORT\" \"$FWDCTL_SADDR\" \"$FWDCTL_SPORT\" > " + out,
		},
	}
	if err := RunOnCreate(rule); err != nil {
		t.Fatalf("RunOnCreate() error = %v", err)
	}
	content, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading env file: %v", err)
	}
	want := "eth0\ntcp\n9090\n127.0.0.1\n8080\n"
	if got := string(content); got != want {
		t.Errorf("injected env = %q, want %q", got, want)
	}
}
