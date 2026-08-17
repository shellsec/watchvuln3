package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPICatalogFollowsHost(t *testing.T) {
	s := NewServer(nil, "0.0.0.0:8765").SetVersion("vtest")
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api", nil)
	require.NoError(t, err)
	req.Host = "192.168.9.9:8765"
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "http://192.168.9.9:8765", body["base_url"])
	mcp, ok := body["mcp"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "http://192.168.9.9:8765/mcp", mcp["url"])
	assert.Equal(t, "streamable-http", mcp["transport"])
}

func TestAPICatalogOnNewIP(t *testing.T) {
	s := NewServer(nil, "0.0.0.0:8765")
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api", nil)
	require.NoError(t, err)
	req.Host = "10.0.0.8:9000"
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "http://10.0.0.8:9000", body["base_url"])
}

func TestMCPInitializeAndToolsList(t *testing.T) {
	s := NewServer(nil, "127.0.0.1:8765").SetVersion("vtest")
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	initBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2025-03-26",
			"capabilities":    map[string]any{},
			"clientInfo":      map[string]string{"name": "test", "version": "1"},
		},
	}
	raw, _ := json.Marshal(initBody)
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/mcp", bytes.NewReader(raw))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var rpc map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&rpc))
	result, ok := rpc["result"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "2025-03-26", result["protocolVersion"])
	info := result["serverInfo"].(map[string]any)
	assert.Equal(t, "watchvuln", info["name"])
	assert.Equal(t, "vtest", info["version"])

	listBody := []byte(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`)
	req, err = http.NewRequest(http.MethodPost, ts.URL+"/mcp", bytes.NewReader(listBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&rpc))
	result, ok = rpc["result"].(map[string]any)
	require.True(t, ok)
	tools, ok := result["tools"].([]any)
	require.True(t, ok)
	names := map[string]bool{}
	for _, item := range tools {
		names[item.(map[string]any)["name"].(string)] = true
	}
	assert.True(t, names["search_vulns"])
	assert.True(t, names["get_vuln"])
	assert.True(t, names["list_sources"])
	assert.True(t, names["get_stats"])
}

func TestMCPInitializedNotification(t *testing.T) {
	s := NewServer(nil, "127.0.0.1:8765")
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/mcp", bytes.NewReader([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}`)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Empty(t, bytes.TrimSpace(body))
}

func TestMCPListSourcesTool(t *testing.T) {
	s := NewServer(nil, "127.0.0.1:8765")
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	raw := []byte(`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_sources","arguments":{}}}`)
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/mcp", bytes.NewReader(raw))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var rpc map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&rpc))
	result := rpc["result"].(map[string]any)
	assert.Equal(t, false, result["isError"])
	content := result["content"].([]any)
	require.NotEmpty(t, content)
	text := content[0].(map[string]any)["text"].(string)
	assert.Contains(t, text, `"id":"avd"`)
}

func TestMCPUnknownMethod(t *testing.T) {
	s := NewServer(nil, "127.0.0.1:8765")
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/mcp", bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":9,"method":"nope"}`)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var rpc map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&rpc))
	errObj := rpc["error"].(map[string]any)
	assert.Equal(t, float64(-32601), errObj["code"])
}

func TestAPISourcesRouteNotCapturedByIndex(t *testing.T) {
	s := NewServer(nil, "127.0.0.1:8765")
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/sources")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var list []map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&list))
	require.NotEmpty(t, list)
}

func TestAPIVulnRequiresQuery(t *testing.T) {
	s := NewServer(nil, "127.0.0.1:8765")
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/vuln")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDashboardContainsAPIMCPEntry(t *testing.T) {
	assert.Contains(t, dashboardHTML, `id="apiMcpBtn"`)
	assert.Contains(t, dashboardHTML, "window.location.origin")
	assert.Contains(t, dashboardHTML, "/mcp")
	assert.Contains(t, dashboardHTML, "GET /api/vulns")
}

func TestParseVulnListParams(t *testing.T) {
	p := parseVulnListParams(url.Values{
		"page":     {"2"},
		"limit":    {"999"},
		"q":        {" CVE-1 "},
		"severity": {"严重"},
		"sort":     {"update"},
	}, 10, 50)
	assert.Equal(t, 2, p.Page)
	assert.Equal(t, 50, p.Limit)
	assert.Equal(t, "CVE-1", p.Q)
	assert.Equal(t, "严重", p.Severity)
	assert.Equal(t, "update", p.Sort)
}

func TestRequestBaseURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.local/api", nil)
	req.Host = "203.0.113.10:8765"
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "vuln.example.com")
	assert.Equal(t, "https://vuln.example.com", requestBaseURL(req))
}

func TestCORSPreflight(t *testing.T) {
	s := NewServer(nil, "127.0.0.1:8765")
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodOptions, ts.URL+"/mcp", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "https://cursor.sh")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, "https://cursor.sh", resp.Header.Get("Access-Control-Allow-Origin"))
}
