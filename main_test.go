package main

import (
	"flag"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

var tracing *bool

func init() {
	tracing = flag.Bool("tracing", false, "Trace syscalls during tests")
}

func TestFwdctl(t *testing.T) {
	flag.Parse()
	testscript.Run(t, testscript.Params{
		Dir:                 "tests",
		Cmds:                customCommands(),
		RequireExplicitExec: true,
		Setup: func(env *testscript.Env) error {
			env.Setenv("GOCOVERDIR", "/tmp/integration")
			if *tracing {
				env.Setenv("TRACING", "true")
			}
			return nil
		},
	})
}
