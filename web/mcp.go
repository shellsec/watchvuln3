package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zema1/watchvuln/ent"
	"github.com/zema1/watchvuln/grab"
)

const (
	mcpProtocol2025_06 = "2025-06-18"
	mcpProtocol2025_03 = "2025-03-26"
	mcpProtocol2024_11 = "2024-11-05"
	mcpJSONRPCVersion  = "2.0"
)

var mcpSupportedProtocols = []string{mcpProtocol2025_06, mcpProtocol2025_03, mcpProtocol2024_11}

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

func (s *Server) handleMCP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleMCPGet(w, r)
	case http.MethodPost:
		s.handleMCPPost(w, r)
	case http.MethodDelete:
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleMCPGet(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	if accept != "" && !strings.Contains(accept, "text/event-stream") {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("MCP-Protocol-Version", mcpProtocol2025_03)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	<-r.Context().Done()
}

func (s *Server) handleMCPPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeRPCError(w, r, nil, -32700, "parse error")
		return
	}
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		writeRPCError(w, r, nil, -32600, "invalid request")
		return
	}

	if body[0] == '[' {
		var batch []rpcRequest
		if err := json.Unmarshal(body, &batch); err != nil {
			writeRPCError(w, r, nil, -32700, "parse error")
			return
		}
		out := make([]rpcResponse, 0, len(batch))
		for _, req := range batch {
			if resp, ok := s.dispatchMCP(r.Context(), req); ok {
				out = append(out, resp)
			}
		}
		if len(out) == 0 {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		writeRPC(w, r, out)
		return
	}

	var req rpcRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeRPCError(w, r, nil, -32700, "parse error")
		return
	}
	resp, ok := s.dispatchMCP(r.Context(), req)
	if !ok {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	writeRPC(w, r, resp)
}

func (s *Server) dispatchMCP(ctx context.Context, req rpcRequest) (rpcResponse, bool) {
	if req.Method == "" {
		return rpcResponse{JSONRPC: mcpJSONRPCVersion, ID: req.ID, Error: &rpcError{Code: -32600, Message: "invalid request"}}, hasRPCID(req.ID)
	}
	if strings.HasPrefix(req.Method, "notifications/") {
		return rpcResponse{}, false
	}
	if !hasRPCID(req.ID) {
		return rpcResponse{}, false
	}

	switch req.Method {
	case "initialize":
		return rpcResponse{JSONRPC: mcpJSONRPCVersion, ID: req.ID, Result: s.mcpInitialize(req.Params)}, true
	case "ping":
		return rpcResponse{JSONRPC: mcpJSONRPCVersion, ID: req.ID, Result: map[string]any{}}, true
	case "tools/list":
		return rpcResponse{JSONRPC: mcpJSONRPCVersion, ID: req.ID, Result: map[string]any{"tools": mcpTools()}}, true
	case "tools/call":
		result, rpcErr := s.mcpCallTool(ctx, req.Params)
		if rpcErr != nil {
			return rpcResponse{JSONRPC: mcpJSONRPCVersion, ID: req.ID, Error: rpcErr}, true
		}
		return rpcResponse{JSONRPC: mcpJSONRPCVersion, ID: req.ID, Result: result}, true
	default:
		return rpcResponse{JSONRPC: mcpJSONRPCVersion, ID: req.ID, Error: &rpcError{Code: -32601, Message: "method not found: " + req.Method}}, true
	}
}

func (s *Server) mcpInitialize(params json.RawMessage) map[string]any {
	protocol := mcpProtocol2025_03
	var in struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if err := json.Unmarshal(params, &in); err == nil {
		for _, p := range mcpSupportedProtocols {
			if in.ProtocolVersion == p {
				protocol = p
				break
			}
		}
	}
	return map[string]any{
		"protocolVersion": protocol,
		"capabilities": map[string]any{
			"tools": map[string]any{"listChanged": false},
		},
		"serverInfo": map[string]any{
			"name":    "watchvuln",
			"version": s.versionLabel(),
		},
		"instructions": "WatchVuln 本地漏洞情报。用 search_vulns 检索，get_vuln 看详情，list_sources / get_stats 查看来源与统计。服务跟看板同一地址，换 IP 后用新 Host 访问 /mcp 即可。",
	}
}

func (s *Server) versionLabel() string {
	if s.version != "" {
		return s.version
	}
	return "dev"
}

func mcpTools() []map[string]any {
	return []map[string]any{
		{
			"name":        "search_vulns",
			"description": "检索本地漏洞情报库。可按关键词、等级、数据源筛选，返回摘要列表（不含长描述）。",
			"inputSchema": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"q":        map[string]any{"type": "string", "description": "标题 / CVE / 描述关键词"},
					"severity": map[string]any{"type": "string", "description": "严重 | 高危 | 中危 | 低危"},
					"source":   map[string]any{"type": "string", "description": "数据源 ID，如 avd、chaitin、kev、seebug"},
					"sort":     map[string]any{"type": "string", "description": "disclosure（默认，按披露日期）或 update（按入库更新）"},
					"page":     map[string]any{"type": "integer", "description": "页码，从 1 起", "minimum": 1},
					"limit":    map[string]any{"type": "integer", "description": "每页条数，默认 10，最大 50", "minimum": 1, "maximum": 50},
				},
			},
		},
		{
			"name":        "get_vuln",
			"description": "按数字 id、CVE 或库内 key 获取单条漏洞完整详情（描述、修复建议、参考链接等）。id / cve / key 至少提供一个。",
			"inputSchema": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"id":  map[string]any{"type": "integer", "description": "漏洞数字 ID"},
					"cve": map[string]any{"type": "string", "description": "CVE 编号，大小写不敏感"},
					"key": map[string]any{"type": "string", "description": "库内唯一 key"},
				},
			},
		},
		{
			"name":        "list_sources",
			"description": "列出 WatchVuln 内置漏洞数据源（ID、名称、链接）。",
			"inputSchema": map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{}},
		},
		{
			"name":        "get_stats",
			"description": "返回本地库漏洞总数及按等级分布。",
			"inputSchema": map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{}},
		},
	}
}

func (s *Server) mcpCallTool(ctx context.Context, params json.RawMessage) (map[string]any, *rpcError) {
	var call struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &call); err != nil || call.Name == "" {
		return nil, &rpcError{Code: -32602, Message: "invalid tools/call params"}
	}
	args := map[string]any{}
	if len(bytes.TrimSpace(call.Arguments)) > 0 && string(call.Arguments) != "null" {
		if err := json.Unmarshal(call.Arguments, &args); err != nil {
			return nil, &rpcError{Code: -32602, Message: "invalid tool arguments"}
		}
	}

	var (
		payload any
		err     error
	)
	switch call.Name {
	case "search_vulns":
		payload, err = s.mcpSearchVulns(ctx, args)
	case "get_vuln":
		payload, err = s.mcpGetVuln(ctx, args)
	case "list_sources":
		payload = grab.ListSources()
	case "get_stats":
		payload, err = s.vulnStats(ctx)
	default:
		return toolError("unknown tool: " + call.Name), nil
	}
	if err != nil {
		if ent.IsNotFound(err) {
			return toolError("vulnerability not found"), nil
		}
		return toolError(err.Error()), nil
	}
	raw, _ := json.Marshal(payload)
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": string(raw)}},
		"isError": false,
	}, nil
}

func (s *Server) mcpSearchVulns(ctx context.Context, args map[string]any) (any, error) {
	p := parseVulnListParams(urlValuesFromArgs(args), 10, 50)
	result, err := s.listVulns(ctx, p)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, compactVulnItem(item))
	}
	return map[string]any{
		"total": result.Total,
		"page":  result.Page,
		"limit": result.Limit,
		"sort":  result.Sort,
		"items": items,
	}, nil
}

func (s *Server) mcpGetVuln(ctx context.Context, args map[string]any) (any, error) {
	id := intFromArg(args["id"])
	cve := stringFromArg(args["cve"])
	key := stringFromArg(args["key"])
	if id <= 0 && cve == "" && key == "" {
		return nil, fmt.Errorf("id, cve or key is required")
	}
	return s.findVuln(ctx, id, cve, key)
}

func toolError(msg string) map[string]any {
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": msg}},
		"isError": true,
	}
}

func urlValuesFromArgs(args map[string]any) url.Values {
	values := url.Values{}
	for _, k := range []string{"q", "severity", "source", "sort", "page", "limit"} {
		if v, ok := args[k]; ok && v != nil {
			values.Set(k, fmt.Sprint(v))
		}
	}
	return values
}

func stringFromArg(v any) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

func intFromArg(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	case string:
		i, _ := strconv.Atoi(strings.TrimSpace(n))
		return i
	default:
		return 0
	}
}

func hasRPCID(id json.RawMessage) bool {
	s := strings.TrimSpace(string(id))
	return s != "" && s != "null"
}

func writeRPCError(w http.ResponseWriter, r *http.Request, id json.RawMessage, code int, msg string) {
	writeRPC(w, r, rpcResponse{
		JSONRPC: mcpJSONRPCVersion,
		ID:      id,
		Error:   &rpcError{Code: code, Message: msg},
	})
}

func writeRPC(w http.ResponseWriter, r *http.Request, v any) {
	w.Header().Set("MCP-Protocol-Version", mcpProtocol2025_03)
	accept := r.Header.Get("Accept")
	raw, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if strings.Contains(accept, "text/event-stream") && !strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = fmt.Fprintf(w, "event: message\ndata: %s\n\n", raw)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(raw)
	_, _ = w.Write([]byte("\n"))
}
