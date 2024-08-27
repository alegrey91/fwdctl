package main

import (
	"flag"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

var trace bool

func init() {
	flag.BoolVar(&trace, "tracing", false, "Trace syscalls from integration-tests")
}

func TestFwdctl(t *testing.T) {
	flag.Parse()
	testscript.Run(t, testscript.Params{
		Dir:                 "tests",
		Cmds:                customCommands(),
		RequireExplicitExec: true,
		Setup: func(env *testscript.Env) error {
			env.Setenv("GOCOVERDIR", "/tmp/integration")
			if trace {
				env.Setenv("TRACING", "true")
			}
			return nil
		},
	})
}
