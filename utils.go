package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	uuid := getRandomString()
	workDir, err := os.Getwd()
	if err != nil {
		ts.Fatalf("unable to find work dir: %v", err)
	}
	customCommand := []string{
		"/usr/sbin/harpoon",
		"capture",
		"-f",
		"main.main",
		"--save",
		"--directory",
		fmt.Sprintf("%s/integration-test-syscalls", workDir),
		"--include-cmd-stdout",
		"--include-cmd-stderr",
		"--name",
		fmt.Sprintf("main_main_%s", uuid),
		"--",
	}

	// find binary path for primary command
	cmdPath, err := exec.LookPath(args[0])
	if err != nil {
		ts.Fatalf("unable to find binary path for %s: %v", args[0], err)
	}
	args[0] = cmdPath
	customCommand = append(customCommand, args...)

	ts.Logf("executing tracing command: %s", strings.Join(customCommand, " "))
	// check if command has '&' as last char to be ran in background
	if len(args) > 0 && backgroundSpecifier.MatchString(args[len(args)-1]) {
		_, err = execBackground(ts, customCommand[0], customCommand[1:len(args)-1]...)
	} else {
		err = ts.Exec(customCommand[0], customCommand[1:]...)
	}
	if err != nil {
		ts.Logf("[%v]\n", err)
		if !neg {
			ts.Fatalf("unexpected go command failure")
		}
	} else {
		if neg {
			ts.Fatalf("unexpected go command success")
		}
	}
}

func execBackground(ts *testscript.TestScript, command string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(command, args...)
	path := ts.MkAbs(".")
	dir, _ := filepath.Split(path)

	var stdoutBuf, stderrBuf strings.Builder
	cmd.Dir = dir
	cmd.Env = append(cmd.Env, "PWD="+dir)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	return cmd, cmd.Start()
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
