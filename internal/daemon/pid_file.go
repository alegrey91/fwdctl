package daemon

import (
	"os"
	"strconv"
)

var (
	pidFilePath = "/tmp/fwdctl.pid"
)

// Create PID file
func createPidFile() error {
	pid := []byte(strconv.Itoa(os.Getpid()))
	err := os.WriteFile(pidFilePath, pid, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Retrieve PID by reading file content
func readPidFile() (int, error) {
	pidB, err := os.ReadFile(pidFilePath)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(string(pidB))
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// Remove PID file
func removePidFile() error {
	err := os.Remove(pidFilePath)
	if err != nil {
		return err
	}
	return nil
}
