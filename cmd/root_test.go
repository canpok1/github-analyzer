package cmd

import (
	"bytes"
	"testing"
)

func TestExecute_NoArgs_ReturnsError(t *testing.T) {
	err := Execute()
	if err == nil {
		t.Fatal("Execute() without args should return error")
	}
}

func TestHelp_ContainsCommandName(t *testing.T) {
	rootCmd := makeRootCmd()

	if rootCmd.Use != "github-analyzer" {
		t.Errorf("expected Use to be 'github-analyzer', got: %s", rootCmd.Use)
	}
}

func TestHelp_ShowsHelpOutput(t *testing.T) {
	rootCmd := makeRootCmd()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned error: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("help output should not be empty")
	}
}
