package antigravity

import (
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestPrepareAntigravityRequestMatchesCLIShape(t *testing.T) {
	req := &GenerateContentRequest{
		Project: "test-project",
		Model:   "gemini-pro-agent",
		Request: GeminiInternalRequest{
			Contents: []Content{{
				Role:  "user",
				Parts: []ContentPart{{Text: "hello"}},
			}},
			SystemInstruction: &SystemInstruction{
				Role:  "user",
				Parts: []ContentPart{{Text: "client system"}},
			},
		},
	}

	prepareAntigravityRequest(req)

	if req.UserAgent != RequestUserAgent {
		t.Fatalf("UserAgent = %q, want %q", req.UserAgent, RequestUserAgent)
	}
	if req.RequestType != RequestTypeAgent {
		t.Fatalf("RequestType = %q, want %q", req.RequestType, RequestTypeAgent)
	}

	requestIDPattern := regexp.MustCompile(`^agent/[0-9a-f-]{36}/[0-9]{13}/[0-9a-f-]{36}/1$`)
	if !requestIDPattern.MatchString(req.RequestID) {
		t.Fatalf("RequestID = %q, want Antigravity CLI shape", req.RequestID)
	}

	if req.Request.SystemInstruction == nil {
		t.Fatal("SystemInstruction is nil")
	}
	if req.Request.SystemInstruction.Role != "user" {
		t.Fatalf("SystemInstruction.Role = %q, want user", req.Request.SystemInstruction.Role)
	}
	parts := req.Request.SystemInstruction.Parts
	if len(parts) != 2 {
		t.Fatalf("SystemInstruction parts = %d, want 2", len(parts))
	}
	if !strings.Contains(parts[0].Text, "<identity>") || !strings.Contains(parts[0].Text, "You are Antigravity") {
		t.Fatalf("first system part does not contain Antigravity identity")
	}
	if strings.Contains(parts[0].Text, "Please ignore the following [ignore]") {
		t.Fatalf("first system part contains legacy ignore injection")
	}
	if parts[1].Text != "client system" {
		t.Fatalf("existing system part = %q, want client system", parts[1].Text)
	}
}

func TestApplyHeadersMatchesAntigravityCLI(t *testing.T) {
	header := http.Header{}
	ApplyHeaders(header, "token", "application/json")

	if got := header.Get("Authorization"); got != "Bearer token" {
		t.Fatalf("Authorization = %q", got)
	}
	if got := header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q", got)
	}
	if got := header.Get("User-Agent"); !strings.HasPrefix(got, "antigravity/cli/1.0.5 ") {
		t.Fatalf("User-Agent = %q", got)
	}
	if got := header.Get("X-Goog-Api-Client"); got != "" {
		t.Fatalf("X-Goog-Api-Client = %q, want empty", got)
	}
	if got := header.Get("Client-Metadata"); got != "" {
		t.Fatalf("Client-Metadata = %q, want empty", got)
	}
	if got := header.Get("Accept"); got != "" {
		t.Fatalf("Accept = %q, want empty", got)
	}
}
