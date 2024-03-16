package main

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestFwdctl(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:                 "tests",
		Cmds:                customCommands(),
		RequireExplicitExec: true,
	})
}
