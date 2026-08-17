package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/zema1/watchvuln/ent"
	"github.com/zema1/watchvuln/grab"
)

func (s *Server) handleAPIIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api" && r.URL.Path != "/api/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	base := requestBaseURL(r)
	writeJSON(w, apiCatalog(base))
}

func apiCatalog(base string) map[string]any {
	base = strings.TrimRight(base, "/")
	return map[string]any{
		"name":        "WatchVuln",
		"description": "高价值漏洞情报 REST API 与 MCP。地址跟随当前访问 Host，换 IP 后用新地址打开即可。",
		"base_url":    base,
		"mcp": map[string]any{
			"url":       base + "/mcp",
			"transport": "streamable-http",
			"tools":     []string{"search_vulns", "get_vuln", "list_sources", "get_stats"},
		},
		"endpoints": []map[string]any{
			{
				"method":      "GET",
				"path":        "/api",
				"description": "本目录（含当前 base_url）",
			},
			{
				"method":      "GET",
				"path":        "/api/stats",
				"description": "漏洞总数与等级分布",
			},
			{
				"method":      "GET",
				"path":        "/api/sources",
				"description": "可用数据源列表",
			},
			{
				"method":      "GET",
				"path":        "/api/vulns",
				"description": "分页检索漏洞",
				"query": map[string]string{
					"q":        "标题 / CVE / 描述关键词",
					"severity": "严重 | 高危 | 中危 | 低危",
					"source":   "数据源 ID，如 avd、chaitin、kev",
					"sort":     "disclosure（默认）| update",
					"page":     "页码，从 1 起",
					"limit":    "每页条数，默认 30，最大 100",
				},
			},
			{
				"method":      "GET",
				"path":        "/api/vuln",
				"description": "按 id / cve / key 取单条详情（至少提供一个）",
				"query": map[string]string{
					"id":  "数字 ID",
					"cve": "CVE 编号，大小写不敏感",
					"key": "库内唯一 key",
				},
			},
		},
		"examples": map[string]string{
			"list":  "curl -s '" + base + "/api/vulns?q=CVE-2024&severity=严重&limit=5'",
			"get":   "curl -s '" + base + "/api/vuln?cve=CVE-2024-0001'",
			"stats": "curl -s '" + base + "/api/stats'",
			"mcp":   "curl -s -X POST '" + base + "/mcp' -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"initialize\",\"params\":{\"protocolVersion\":\"2025-03-26\",\"capabilities\":{},\"clientInfo\":{\"name\":\"curl\",\"version\":\"1.0\"}}}'",
		},
	}
}

func (s *Server) handleAPISources(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, grab.ListSources())
}

func (s *Server) handleAPIStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.vulnStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, stats)
}

func (s *Server) handleAPIVulns(w http.ResponseWriter, r *http.Request) {
	p := parseVulnListParams(r.URL.Query(), 30, 100)
	result, err := s.listVulns(r.Context(), p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, result)
}

func (s *Server) handleAPIVuln(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id, _ := strconv.Atoi(strings.TrimSpace(q.Get("id")))
	cve := strings.TrimSpace(q.Get("cve"))
	key := strings.TrimSpace(q.Get("key"))
	if id <= 0 && cve == "" && key == "" {
		writeError(w, http.StatusBadRequest, errQuery("id, cve or key is required"))
		return
	}
	item, err := s.findVuln(r.Context(), id, cve, key)
	if err != nil {
		if ent.IsNotFound(err) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, item)
}

type queryError string

func (e queryError) Error() string { return string(e) }

func errQuery(msg string) error { return queryError(msg) }

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func writeError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	writeJSON(w, map[string]string{"error": err.Error()})
}
