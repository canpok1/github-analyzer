package templates

import (
	"strings"
	"testing"
)

func TestConfigTemplate_ContainsExpectedFields(t *testing.T) {
	if len(ConfigTemplate) == 0 {
		t.Fatal("ConfigTemplate is empty")
	}

	content := string(ConfigTemplate)
	expectedFields := []string{"repo:", "tone:", "default_prompt:", "model:"}
	for _, field := range expectedFields {
		if !strings.Contains(content, field) {
			t.Errorf("ConfigTemplate does not contain field %q", field)
		}
	}
}
