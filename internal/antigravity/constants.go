package antigravity

import (
	"fmt"
	"net/http"
	"runtime"
)

const (
	endpointDaily    = "https://daily-cloudcode-pa.googleapis.com"
	userAgentVersion = "1.0.5"
	RequestUserAgent = "antigravity"
	RequestTypeAgent = "agent"
)

var Endpoints = []string{
	endpointDaily,
}

func platformUserAgent() string {
	return fmt.Sprintf("antigravity/cli/%s %s/%s", userAgentVersion, runtime.GOOS, runtime.GOARCH)
}

func ApplyHeaders(header http.Header, token string, accept string) {
	header.Set("Authorization", "Bearer "+token)
	header.Set("Content-Type", "application/json")
	header.Set("User-Agent", platformUserAgent())
}
