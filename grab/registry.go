package grab

import (
	"fmt"
	"strings"
)

// DefaultSources comma-separated default for CLI / env.
const DefaultSources = "avd,chaitin,nox,oscs,threatbook,seebug,struts2,kev,venustech"

// DefaultSourceList is the default enabled source ids.
var DefaultSourceList = []string{"avd", "chaitin", "nox", "oscs", "threatbook", "seebug", "struts2", "kev", "venustech"}

// SourceMeta describes a vulnerability data source for CLI / web UI.
type SourceMeta struct {
	ID          string   `json:"id"`
	Aliases     []string `json:"aliases,omitempty"`
	Name        string   `json:"name"`
	Link        string   `json:"link"`
	InDefault   bool     `json:"in_default"`
	HostMatches []string `json:"-"`
}

// ListSources returns metadata for all built-in grab sources.
func ListSources() []SourceMeta {
	defaultSet := make(map[string]struct{})
	for _, s := range DefaultSourceList {
		defaultSet[s] = struct{}{}
	}
	return []SourceMeta{
		sourceMeta("avd", nil, NewAVDCrawler(), defaultSet),
		sourceMeta("chaitin", nil, NewChaitinCrawler(), defaultSet),
		sourceMeta("nox", []string{"ti"}, NewTiCrawler(), defaultSet),
		sourceMeta("oscs", nil, NewOSCSCrawler(), defaultSet),
		sourceMeta("threatbook", nil, NewThreatBookCrawler(), defaultSet),
		sourceMeta("seebug", nil, NewSeebugCrawler(), defaultSet),
		sourceMeta("struts2", []string{"structs2"}, NewStruts2Crawler(), defaultSet),
		sourceMeta("kev", nil, NewKEVCrawler(), defaultSet),
		sourceMeta("venustech", nil, NewVenustechCrawler(), defaultSet),
	}
}

func sourceMeta(id string, aliases []string, g Grabber, defaultSet map[string]struct{}) SourceMeta {
	p := g.ProviderInfo()
	_, inDefault := defaultSet[id]
	return SourceMeta{
		ID:          id,
		Aliases:     aliases,
		Name:        p.DisplayName,
		Link:        p.Link,
		InDefault:   inDefault,
		HostMatches: sourceHostMatches(id, p.Link),
	}
}

func sourceHostMatches(id, link string) []string {
	switch id {
	case "avd":
		return []string{"avd.aliyun.com"}
	case "chaitin":
		return []string{"stack.chaitin.com"}
	case "nox":
		return []string{"ti.qianxin.com"}
	case "oscs":
		return []string{"oscs1024.com"}
	case "threatbook":
		return []string{"threatbook.com"}
	case "seebug":
		return []string{"seebug.org"}
	case "struts2":
		return []string{"cwiki.apache.org", "struts.apache.org"}
	case "kev":
		return []string{"cisa.gov", "_KEV"}
	case "venustech":
		return []string{"venustech.com.cn"}
	default:
		if link != "" {
			return []string{strings.TrimPrefix(strings.TrimPrefix(link, "https://"), "http://")}
		}
		return nil
	}
}

// SourceMetaByID finds source metadata by id or alias.
func SourceMetaByID(id string) (*SourceMeta, bool) {
	id = strings.ToLower(strings.TrimSpace(id))
	for _, s := range ListSources() {
		if s.ID == id {
			return &s, true
		}
		for _, a := range s.Aliases {
			if a == id {
				return &s, true
			}
		}
	}
	return nil, false
}

// BuildGrabbers creates grabbers from configured source ids.
func BuildGrabbers(sourceIDs []string) ([]Grabber, error) {
	var grabs []Grabber
	for _, part := range sourceIDs {
		part = strings.ToLower(strings.TrimSpace(part))
		if part == "" {
			continue
		}
		switch part {
		case "chaitin":
			grabs = append(grabs, NewChaitinCrawler())
		case "avd":
			grabs = append(grabs, NewAVDCrawler())
		case "nox", "ti":
			grabs = append(grabs, NewTiCrawler())
		case "oscs":
			grabs = append(grabs, NewOSCSCrawler())
		case "seebug":
			grabs = append(grabs, NewSeebugCrawler())
		case "threatbook":
			grabs = append(grabs, NewThreatBookCrawler())
		case "struts2", "structs2":
			grabs = append(grabs, NewStruts2Crawler())
		case "kev":
			grabs = append(grabs, NewKEVCrawler())
		case "venustech":
			grabs = append(grabs, NewVenustechCrawler())
		default:
			return nil, fmt.Errorf("invalid grab source %s (run 'watchvuln list-sources' for ids)", part)
		}
	}
	return grabs, nil
}

// FormatSourcesTable returns a human-readable table of sources.
func FormatSourcesTable() string {
	var b strings.Builder
	b.WriteString("ID       | 默认 | 名称                   | 链接\n")
	b.WriteString("---------+------+------------------------+------------------------------------------\n")
	for _, s := range ListSources() {
		def := "否"
		if s.InDefault {
			def = "是"
		}
		alias := ""
		if len(s.Aliases) > 0 {
			alias = fmt.Sprintf(" (别名: %s)", strings.Join(s.Aliases, ", "))
		}
		b.WriteString(fmt.Sprintf("%-8s | %s  | %-22s | %s%s\n", s.ID, def, s.Name, s.Link, alias))
	}
	return b.String()
}
