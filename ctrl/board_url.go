package ctrl

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// BoardPublicURL returns the URL appended to push messages for the vuln board.
// Port is always taken from web_addr so changing --web-addr port updates the link automatically.
func (c *WatchVulnAppConfig) BoardPublicURL() string {
	addr := strings.TrimSpace(c.WebAddr)
	if addr == "" {
		return ""
	}

	if raw := strings.TrimSpace(c.WebPublicURL); raw != "" {
		return normalizeBoardURL(raw)
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return ""
	}

	if pubHost := strings.TrimSpace(c.WebPublicHost); pubHost != "" {
		return buildBoardURL("http", pubHost, port, "")
	}

	switch host {
	case "0.0.0.0", "::":
		return ""
	}

	return buildBoardURL("http", host, port, "")
}

func buildBoardURL(scheme, host, port, path string) string {
	if path == "" {
		path = "/"
	} else if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(host, port),
		Path:   path,
	}
	return u.String()
}

func normalizeBoardURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return ""
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	if u.Path == "" {
		u.Path = "/"
	}
	return u.String()
}

// BoardPublicURLHint describes how the push link is resolved (for startup logs).
func (c *WatchVulnAppConfig) BoardPublicURLHint() string {
	if strings.TrimSpace(c.WebAddr) == "" {
		return ""
	}
	if u := c.BoardPublicURL(); u != "" {
		return fmt.Sprintf("push link: %s", u)
	}
	if strings.TrimSpace(c.WebPublicHost) == "" && isWildcardListenAddr(c.WebAddr) {
		return "push link: disabled (set web_public_host or web_public_url when listening on 0.0.0.0)"
	}
	return "push link: disabled (invalid web_addr)"
}

func isWildcardListenAddr(addr string) bool {
	host, _, err := net.SplitHostPort(strings.TrimSpace(addr))
	if err != nil {
		return false
	}
	return host == "0.0.0.0" || host == "::"
}
