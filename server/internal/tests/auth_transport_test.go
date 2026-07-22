package tests

import (
	"net/http"
	"strings"

	"github.com/ddc-111/agentGame/server/internal/network"
)

const skipTestAuthHeader = "X-AgentGame-Test-Skip-Auth"

type managementAuthTransport struct {
	base  http.RoundTripper
	token string
}

func (t managementAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get(skipTestAuthHeader) != "" {
		clone := req.Clone(req.Context())
		clone.Header = req.Header.Clone()
		clone.Header.Del(skipTestAuthHeader)
		return t.base.RoundTrip(clone)
	}
	if !isManagementTestPath(req.URL.Path) || req.Header.Get("Authorization") != "" {
		return t.base.RoundTrip(req)
	}

	clone := req.Clone(req.Context())
	clone.Header = req.Header.Clone()
	clone.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(clone)
}

func isManagementTestPath(path string) bool {
	if path == "/api/gm/login" {
		return false
	}
	if path == "/mcp" || strings.HasPrefix(path, "/mcp/") {
		return true
	}

	prefixes := []string{
		"/api/gm/",
		"/api/generator/",
		"/api/mcp/",
		"/api/scenes",
		"/api/npcs",
		"/api/agents",
		"/api/llm/",
		"/api/prompts",
		"/api/shops",
		"/api/items",
		"/api/tasks",
		"/api/flows",
		"/api/players",
		"/api/conversations",
		"/api/config",
		"/api/export",
		"/api/import",
	}
	for _, prefix := range prefixes {
		if strings.HasSuffix(prefix, "/") {
			if strings.HasPrefix(path, prefix) {
				return true
			}
			continue
		}
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}

func init() {
	token, err := network.GenerateJWT("change-me-in-production", "admin", "gm", 1)
	if err != nil {
		panic(err)
	}
	base := http.DefaultTransport
	http.DefaultTransport = managementAuthTransport{base: base, token: token}
}
