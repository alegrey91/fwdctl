package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func strace(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) < 1 {
		ts.Fatalf("syntax: strace needs at least one argument")
	}

	outputDir := "syscalls"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		ts.Fatalf("failed to create output directory: %v", err)
	}

	// Prepare the strace output file
	timestamp := time.Now().Format("20060102-150405")
	outputFile := filepath.Join(outputDir, fmt.Sprintf("strace-%s.log", timestamp))

	// Prepare the command to execute with strace, filtering only system calls
	straceArgs := []string{"-f", "-e", "trace=all"}
	straceArgs = append(straceArgs, args...)
	strace := exec.Command("strace", straceArgs...)

	var stdoutBuf bytes.Buffer
	strace.Stdout = &stdoutBuf
	var stderrBuf bytes.Buffer
	strace.Stderr = &stderrBuf

	// Run the command
	if err := strace.Run(); err != nil {
		if !neg {
			ts.Fatalf("command failed: %v", err)
		}
	} else {
		if neg {
			ts.Fatalf("command succeeded when it should have failed")
		}
	}

	fmt.Fprintf(ts.Stdout(), "%s", stdoutBuf.String())

	fmt.Fprintf(ts.Stdout(), "%s", stderrBuf.String())
	syscalls := processStraceOutput(stderrBuf.String())
	if err := os.WriteFile(outputFile, []byte(syscalls), 0644); err != nil {
		ts.Fatalf("error saving strace output to %s", outputFile)
	}
	ts.Logf("strace output saved to %s", outputFile)
}

func processStraceOutput(output string) string {
	lines := strings.Split(output, "\n")

	// Use a map to store unique system call names
	syscalls := make(map[string]struct{})

	// Iterate over each line and extract the system call name
	for _, line := range lines {
		// Extract the system call name before the first '('
		if idx := strings.Index(line, "("); idx != -1 {
			syscall := strings.TrimSpace(line[:idx])
			if syscall != "" {
				syscalls[syscall] = struct{}{}
			}
		}
	}

	// Convert the map keys to a slice to get unique system call names
	var syscallList []string
	for syscall := range syscalls {
		syscallList = append(syscallList, syscall)
	}

	// Join the unique system call names into a single string, one per line
	return strings.Join(syscallList, "\n")
}

func customCommands() map[string]func(ts *testscript.TestScript, neg bool, args []string) {
	return map[string]func(ts *testscript.TestScript, neg bool, args []string){

		// fwd_exists check that the given forward exists
		// invoke as "fwd_exists iface proto dest_port src_addr src_port"
		"fwd_exists": fwdExists,
		"strace":     strace,
	}
}
