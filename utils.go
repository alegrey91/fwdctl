package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/alegrey91/fwdctl/pkg/iptables"
	goiptables "github.com/coreos/go-iptables/iptables"
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/u-root/u-root/pkg/strace"
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

func straceCmd(ts *testscript.TestScript, neg bool, args []string) {
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
	tracedCmd := exec.Command(args[0], args[1:]...)

	var stdoutBuf bytes.Buffer
	tracedCmd.Stdout = &stdoutBuf
	var stderrBuf bytes.Buffer
	tracedCmd.Stderr = &stderrBuf

	// Run the command
	writer := new(bytes.Buffer)
	if err := strace.Strace(tracedCmd, writer); err != nil {
		if !neg {
			ts.Fatalf("command failed: %v", err)
		}
	} else {
		if neg {
			ts.Fatalf("command succeeded when it should have failed")
		}
	}

	fmt.Fprintf(ts.Stdout(), "%s", stdoutBuf.String())

	syscalls, err := processStraceOutput(writer.String())
	if err != nil {
		ts.Fatalf("error processing strace output: %v", err)
	}

	err = writeToFile(outputFile, syscalls)
	if err != nil {
		ts.Fatalf("error creating file: %v", err)
	}

	ts.Logf("strace output saved to %s", outputFile)
}

func writeToFile(outputFile string, syscalls []string) error {
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range syscalls {
		_, _ = datawriter.WriteString(data + "\n")
	}
	datawriter.Flush()
	file.Close()
	return nil
}

func processStraceOutput(straceOutput string) ([]string, error) {
	// Regular expression to match system calls in the format "[pid xxx] E syscall("
	re := regexp.MustCompile(`\[pid \d+\] E (\w+)\(`)

	// Find all matches of the pattern in the strace output
	matches := re.FindAllStringSubmatch(straceOutput, -1)

	// Use a map to keep track of unique system calls
	systemCalls := make(map[string]struct{})

	for _, match := range matches {
		if len(match) > 1 {
			systemCalls[match[1]] = struct{}{}
		}
	}

	// Convert map keys to a slice
	var result []string
	for call := range systemCalls {
		result = append(result, call)
	}

	return result, nil
}

func customCommands() map[string]func(ts *testscript.TestScript, neg bool, args []string) {
	return map[string]func(ts *testscript.TestScript, neg bool, args []string){

		// fwd_exists check that the given forward exists
		// invoke as "fwd_exists iface proto dest_port src_addr src_port"
		"fwd_exists": fwdExists,
		"strace":     straceCmd,
	}
}
