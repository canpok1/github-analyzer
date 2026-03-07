package gemini

import (
	"encoding/json"
	"testing"
)

func TestGeminiRequest_MarshalJSON(t *testing.T) {
	req := geminiRequest{
		Contents: []content{
			{
				Parts: []part{
					{Text: "Hello"},
				},
			},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	contents, ok := decoded["contents"].([]any)
	if !ok || len(contents) != 1 {
		t.Fatalf("expected 1 content, got %v", decoded["contents"])
	}
}

func TestGeminiResponse_UnmarshalJSON(t *testing.T) {
	body := `{
		"candidates": [
			{
				"content": {
					"parts": [
						{"text": "Analysis result"}
					]
				}
			}
		]
	}`

	var resp geminiResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if len(resp.Candidates) != 1 {
		t.Fatalf("expected 1 candidate, got %d", len(resp.Candidates))
	}
	if len(resp.Candidates[0].Content.Parts) != 1 {
		t.Fatalf("expected 1 part, got %d", len(resp.Candidates[0].Content.Parts))
	}
	if resp.Candidates[0].Content.Parts[0].Text != "Analysis result" {
		t.Errorf("text = %q, want %q", resp.Candidates[0].Content.Parts[0].Text, "Analysis result")
	}
}

func TestGeminiErrorResponse_UnmarshalJSON(t *testing.T) {
	body := `{
		"error": {
			"code": 429,
			"message": "Resource has been exhausted",
			"status": "RESOURCE_EXHAUSTED"
		}
	}`

	var resp geminiErrorResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if resp.Error.Code != 429 {
		t.Errorf("code = %d, want %d", resp.Error.Code, 429)
	}
	if resp.Error.Message != "Resource has been exhausted" {
		t.Errorf("message = %q, want %q", resp.Error.Message, "Resource has been exhausted")
	}
	if resp.Error.Status != "RESOURCE_EXHAUSTED" {
		t.Errorf("status = %q, want %q", resp.Error.Status, "RESOURCE_EXHAUSTED")
	}
}
