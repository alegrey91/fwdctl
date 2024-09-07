# Integration Tests

This directory is dedicated to host integration tests written with [`testscript`](https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript).

## Converting from test to trace file

As you can see, for each command we have a couple of files (eg. `apply.txtar`, `apply_trace.txtar`).

The `<command>.txtar` file is used for testing purposes and is run at the beginning of the pipeline to ensure the binary is behave like we expect.

The `<command>_trace.txtar`, on the other side, is used for tracing purposes. This means that we re-run the same commands of the previous file, tracing them with [`harpoon`](https://github.com/alegrey91/harpoon) under the hood.

The file content are quite similar, there are just few differences to follow:

* Each `exec` of the command under test (eg. `fwdctl apply`) have to be replaced with `exec_cmd`.
  
  `exec_cmd` is a custom testscript function to trace the command using `harpoon`.

* If the `exec` of the command under test had a negation (`!`), this should not be added in the command used with `exec_cmd`.
  
  This because in this case we don't care about the result of `harpoon` that will execute the real command.

Here's an example:

```txt
# command.txtar

# normal execution of command
exec command list -a
# -x flag doesn't exists, so this should handle the error
! exec command list -x
```

Should be converted into this:

```txt
# command_trace.txtar

exec_cmd command list -a
exec_cmd command list -x
```