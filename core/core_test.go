package core

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

var serverCmd *exec.Cmd

func TestMain(m *testing.M) {
	// Start the server before running tests
	serverCmd = exec.Command("bao", "server", "-dev", "-dev-root-token-id=dev")

	// Optionally capture server output
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	err := serverCmd.Start()
	if err != nil {
		panic(err)
	}

	// Give the server a moment to start up
	time.Sleep(2 * time.Second)

	// Run the tests
	code := m.Run()

	// Cleanup: Stop the server after tests
	if serverCmd.Process != nil {
		serverCmd.Process.Kill()
	}

	os.Exit(code)
}

func TestAnalyze(t *testing.T) {
	got := 3
	if got != 1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}
