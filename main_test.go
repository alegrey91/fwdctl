package main

import (
	"flag"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestFwdctl(t *testing.T) {
	flag.Parse()
	testscript.Run(t, testscript.Params{
		Dir:                 "tests",
		Cmds:                customCommands(),
		RequireExplicitExec: true,
		Setup: func(env *testscript.Env) error {
			env.Setenv("GOCOVERDIR", "/tmp/integration")
			return nil
		},
	})
}
