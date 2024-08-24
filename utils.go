package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/alegrey91/fwdctl/pkg/iptables"
	goiptables "github.com/coreos/go-iptables/iptables"
	"github.com/rogpeppe/go-internal/testscript"
)

func fwdExists(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) < 5 {
		ts.Fatalf("syntax: fwd_exists iface proto dest_port src_addr src_port")
	}

	ipt, err := goiptables.New()
	if err != nil {
		ts.Fatalf("error creating iptables instance: %q", err)
	}

	ruleSpec := []string{
		"-i", args[0], // interface
		"-p", args[1], // protocol
		"-m", args[1], // protocol
		"--dport", args[2], // destination-port
		"-m", "comment", "--comment", "fwdctl",
		"-j", iptables.FwdTarget, // target (DNAT)
		"--to-destination", args[3] + ":" + args[4], // source-address / source-port
	}

	exists, err := ipt.Exists(iptables.FwdTable, iptables.FwdChain, ruleSpec...)
	if err != nil {
		ts.Fatalf("error checking rule: %v", err)
	}
	if neg && !exists {
		ts.Logf("forward doesn't exist")
		return
	}
	if !exists {
		ts.Fatalf("forward doesn't exist")
	}
}

func cmdTraceSyscalls(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) == 0 {
		ts.Fatalf("usage: tracesyscalls <command> [args...]")
	}

	command := args[0]
	commandArgs := args[1:]

	// Directory to store syscall logs
	syscallDir := filepath.Join("syscalls")
	if err := os.MkdirAll(syscallDir, 0755); err != nil {
		ts.Fatalf("failed to create syscall logs directory: %v", err)
	}

	// Generate a filename based on the command and arguments
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s.log", timestamp)
	filepath := filepath.Join(syscallDir, filename)

	// Capture stdout, stderr, and the exit code
	var stdout, stderr string
	stdout, stderr, err := TraceSyscalls(filepath, command, commandArgs...)

	// Write stdout and stderr to files for comparison
	if stdout != "" {
		ts.Logf("[stdout]\n%s", stdout)
	}
	if stderr != "" {
		ts.Logf("[stderr]\n%s", stderr)
	}

	// Write the exit code to a file
	if neg && err == nil {
		ts.Fatalf("expected command to fail, but it succeeded")
	} else if !neg && err != nil {
		ts.Fatalf("tracesyscalls failed: %v", err)
	}
}

func TraceSyscalls(outputFile string, command string, args ...string) (stdout, stderr string, err error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Start the command without waiting for it to finish
	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Ptrace: true}

	// Bind stdout and stderr
	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		return "", "", fmt.Errorf("failed to start command: %v", err)
	}

	// Attach to the process using Ptrace
	pid := cmd.Process.Pid
	if err := syscall.PtraceAttach(pid); err != nil {
		return "", "", fmt.Errorf("failed to attach ptrace to the process: %v", err)
	}

	// Wait for the child to stop due to the PTRACE_ATTACH
	ws := syscall.WaitStatus(0)
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		return "", "", fmt.Errorf("wait4 failed: %v", err)
	}

	// Open the output file
	file, err := os.Create(outputFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Continue the child process and trace its system calls
	for {
		// Continue the process until the next system call
		if err := syscall.PtraceSyscall(pid, 0); err != nil {
			return "", "", fmt.Errorf("ptrace syscall failed: %v", err)
		}

		// Wait for the child to stop again at the next system call
		_, err := syscall.Wait4(pid, &ws, 0, nil)
		if err != nil {
			return "", "", fmt.Errorf("wait4 failed: %v", err)
		}

		// Check if the process has exited
		if ws.Exited() || ws.Signaled() {
			break
		}

		// Get the system call number
		var regs syscall.PtraceRegs
		if err := syscall.PtraceGetRegs(pid, &regs); err != nil {
			return "", "", fmt.Errorf("ptrace get regs failed: %v", err)
		}

		// Write the syscall number to the file
		_, err = fmt.Fprintf(file, "%d\n", regs.Orig_rax)
		if err != nil {
			return "", "", fmt.Errorf("failed to write syscall to file: %v", err)
		}
	}

	// Wait for the command to complete
	err = cmd.Wait()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return "", "", fmt.Errorf("%d", exitErr.ExitCode())
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}

func customCommands() map[string]func(ts *testscript.TestScript, neg bool, args []string) {
	return map[string]func(ts *testscript.TestScript, neg bool, args []string){

		// fwd_exists check that the given forward exists
		// invoke as "fwd_exists iface proto dest_port src_addr src_port"
		"fwd_exists":    fwdExists,
		"tracesyscalls": cmdTraceSyscalls,
	}
}
