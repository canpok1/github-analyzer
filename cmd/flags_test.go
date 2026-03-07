package cmd

import (
	"testing"
)

func TestValidation_TodayAndSinceConflict(t *testing.T) {
	cmd := makeRootCmd()
	cmd.SetArgs([]string{"--today", "--since", "7d"})
	_ = cmd.ParseFlags([]string{"--today", "--since", "7d"})
	err := validateFlags(cmd)
	if err == nil {
		t.Error("expected error when --today and --since are both specified")
	}
}

func TestValidation_PRAndIssueConflict(t *testing.T) {
	cmd := makeRootCmd()
	_ = cmd.ParseFlags([]string{"--pr", "123", "--issue", "456"})
	err := validateFlags(cmd)
	if err == nil {
		t.Error("expected error when --pr and --issue are both specified")
	}
}

func TestValidation_NoTargetSpecified(t *testing.T) {
	cmd := makeRootCmd()
	_ = cmd.ParseFlags([]string{})
	err := validateFlags(cmd)
	if err == nil {
		t.Error("expected error when no target is specified")
	}
}

func TestValidation_TodayOnly(t *testing.T) {
	cmd := makeRootCmd()
	_ = cmd.ParseFlags([]string{"--today"})
	err := validateFlags(cmd)
	if err != nil {
		t.Errorf("--today only should be valid, got error: %v", err)
	}
}

func TestValidation_SinceOnly(t *testing.T) {
	cmd := makeRootCmd()
	_ = cmd.ParseFlags([]string{"--since", "7d"})
	err := validateFlags(cmd)
	if err != nil {
		t.Errorf("--since only should be valid, got error: %v", err)
	}
}

func TestValidation_PROnly(t *testing.T) {
	cmd := makeRootCmd()
	_ = cmd.ParseFlags([]string{"--pr", "123"})
	err := validateFlags(cmd)
	if err != nil {
		t.Errorf("--pr only should be valid, got error: %v", err)
	}
}

func TestValidation_IssueOnly(t *testing.T) {
	cmd := makeRootCmd()
	_ = cmd.ParseFlags([]string{"--issue", "456"})
	err := validateFlags(cmd)
	if err != nil {
		t.Errorf("--issue only should be valid, got error: %v", err)
	}
}

func TestValidation_TodayAndPR(t *testing.T) {
	cmd := makeRootCmd()
	_ = cmd.ParseFlags([]string{"--today", "--pr", "123"})
	err := validateFlags(cmd)
	if err != nil {
		t.Errorf("--today and --pr should be valid, got error: %v", err)
	}
}

func TestValidation_SinceAndStatus(t *testing.T) {
	cmd := makeRootCmd()
	_ = cmd.ParseFlags([]string{"--since", "7d", "--status", "open"})
	err := validateFlags(cmd)
	if err != nil {
		t.Errorf("--since and --status should be valid, got error: %v", err)
	}
}

func TestValidation_SinceInvalidValue(t *testing.T) {
	cmd := makeRootCmd()
	_ = cmd.ParseFlags([]string{"--since", "abc"})
	err := validateFlags(cmd)
	if err == nil {
		t.Error("expected error for invalid --since value")
	}
}

func TestFlags_Defined(t *testing.T) {
	cmd := makeRootCmd()

	flags := []string{"today", "since", "pr", "issue", "status", "prompt", "repo", "output"}
	for _, name := range flags {
		t.Run(name, func(t *testing.T) {
			f := cmd.Flags().Lookup(name)
			if f == nil {
				t.Errorf("flag --%s is not defined", name)
			}
		})
	}
}
