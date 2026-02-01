package main

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpCommands(t *testing.T) {
	longs := []string{}
	for _, cmd := range rootCmd.Commands() {
		longs = append(longs, cmd.Long)
	}

	assert.Contains(t, longs, HTTP_COMMAND)
}

func TestMainCommand(t *testing.T) {
	// Run the crashing code when FLAG is set
	if os.Getenv("TestMainCommand") == "1" {
		main()

		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestMainCommand", "invalid")
	cmd.Env = append(os.Environ(), "TestMainCommand=1")
	err := cmd.Run()
	time.Sleep(1 * time.Second)
	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	assert.Equal(t, true, ok)
	assert.NotNil(t, e)
}
