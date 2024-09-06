package main

import (
	"crypto/rand"
	"encoding/base64"
	"os/exec"
	"regexp"
	"strings"

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

//nolint:all
func execCmd(ts *testscript.TestScript, neg bool, args []string) {
	var backgroundSpecifier = regexp.MustCompile(`^&([a-zA-Z_0-9]+&)?$`)
	//uuid := getRandomString()
	customCommand := []string{
		"/usr/local/bin/harpoon",
		"capture",
		"-f",
		"main.main",
		"--save",
		"--directory",
		"integration-test-syscalls",
		"--include-cmd-stdout",
		"--include-cmd-stderr",
		//"--name",
		//fmt.Sprintf("main_main_%s", uuid),
		"--",
	}
	customCommand = append(customCommand, args...)
	ts.Logf("executing tracing command: %s", strings.Join(customCommand, " "))
	if backgroundSpecifier.MatchString(args[len(args)-1]) {
		execBackground(args[0], args[1:]...)
		return
	}
	err := ts.Exec(customCommand[0], customCommand[1:]...)

	if err != nil {
		if neg {
			ts.Logf("expected error: %v", err)
			return
		}
		ts.Fatalf("error: %v", err)
	}
}

func execBackground(command string, args ...string) {
	cmd := exec.Command(command, args...)
	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	cmd.Start()
}

//nolint:all
func getRandomString() string {
	b := make([]byte, 4) // 4 bytes will give us 6 base64 characters
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	randomString := base64.URLEncoding.EncodeToString(b)[:6]
	return randomString
}

func customCommands() map[string]func(ts *testscript.TestScript, neg bool, args []string) {
	return map[string]func(ts *testscript.TestScript, neg bool, args []string){

		// fwd_exists check that the given forward exists
		// invoke as "fwd_exists iface proto dest_port src_addr src_port"
		"fwd_exists": fwdExists,
		"exec_cmd":   execCmd,
	}
}
