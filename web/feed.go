package web

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zema1/watchvuln/ent"
	"github.com/zema1/watchvuln/ent/vulninformation"
)

const feedDefaultLimit = 50

type feedItem struct {
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	GUID        string
}

type feedChannel struct {
	Title       string
	Link        string
	Description string
	FeedURL     string
	Items       []feedItem
}

func (s *Server) handleFeed(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/feed.xml" {
		http.NotFound(w, r)
		return
	}

	ctx := r.Context()
	rows, err := s.db.VulnInformation.Query().
		Where(vulninformation.Pushed(true)).
		Order(
			ent.Desc(vulninformation.FieldDisclosure),
			ent.Desc(vulninformation.FieldUpdateTime),
		).
		Limit(feedDefaultLimit).
		All(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	baseURL := requestBaseURL(r)
	channel := feedChannel{
		Title:       "WatchVuln 漏洞情报",
		Link:        baseURL + "/",
		Description: "高价值漏洞推送订阅（最近50条已推送）",
		FeedURL:     baseURL + "/feed.xml",
		Items:       make([]feedItem, 0, len(rows)),
	}
	for _, row := range rows {
		channel.Items = append(channel.Items, feedItem{
			Title:       feedItemTitle(row),
			Link:        feedItemLink(row),
			Description: feedItemDescription(row),
			PubDate:     feedItemPubDate(row),
			GUID:        row.Key,
		})
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	_, _ = w.Write(renderRSS(channel))
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = strings.TrimSpace(strings.Split(proto, ",")[0])
	}
	host := r.Host
	if fwd := r.Header.Get("X-Forwarded-Host"); fwd != "" {
		host = strings.TrimSpace(strings.Split(fwd, ",")[0])
	}
	return scheme + "://" + host
}

func feedItemTitle(row *ent.VulnInformation) string {
	title := strings.TrimSpace(row.Title)
	if row.Cve != "" {
		return fmt.Sprintf("[%s] %s", row.Cve, title)
	}
	return title
}

func feedItemLink(row *ent.VulnInformation) string {
	if link := strings.TrimSpace(row.From); link != "" {
		return link
	}
	return ""
}

func feedItemPubDate(row *ent.VulnInformation) time.Time {
	if !row.UpdateTime.IsZero() {
		return row.UpdateTime
	}
	if t, ok := parseDisclosureDate(row.Disclosure); ok {
		return t
	}
	if !row.CreateTime.IsZero() {
		return row.CreateTime
	}
	return time.Now()
}

func parseDisclosureDate(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	layouts := []string{
		"2006-01-02",
		"2006/01/02",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func feedItemDescription(row *ent.VulnInformation) string {
	var b strings.Builder
	if row.Severity != "" {
		b.WriteString("<p><strong>等级：</strong>")
		b.WriteString(xmlEscape(row.Severity))
		b.WriteString("</p>")
	}
	if row.Cve != "" {
		b.WriteString("<p><strong>CVE：</strong>")
		b.WriteString(xmlEscape(row.Cve))
		b.WriteString("</p>")
	}
	if row.Disclosure != "" {
		b.WriteString("<p><strong>披露日期：</strong>")
		b.WriteString(xmlEscape(row.Disclosure))
		b.WriteString("</p>")
	}
	if len(row.Tags) > 0 {
		b.WriteString("<p><strong>标签：</strong>")
		b.WriteString(xmlEscape(strings.Join(row.Tags, ", ")))
		b.WriteString("</p>")
	}
	if desc := strings.TrimSpace(row.Description); desc != "" {
		b.WriteString("<p>")
		b.WriteString(xmlEscape(truncateRunes(desc, 500)))
		b.WriteString("</p>")
	}
	if sol := strings.TrimSpace(row.Solutions); sol != "" {
		b.WriteString("<p><strong>修复建议：</strong>")
		b.WriteString(xmlEscape(truncateRunes(sol, 300)))
		b.WriteString("</p>")
	}
	return b.String()
}

func truncateRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}

func renderRSS(ch feedChannel) []byte {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	buf.WriteString(`<rss version="2.0">` + "\n")
	buf.WriteString("<channel>\n")
	writeXMLTextEl(&buf, "title", ch.Title)
	writeXMLTextEl(&buf, "link", ch.Link)
	writeXMLTextEl(&buf, "description", ch.Description)
	writeXMLTextEl(&buf, "generator", "WatchVuln")
	writeXMLTextEl(&buf, "lastBuildDate", time.Now().UTC().Format(time.RFC1123Z))
	for _, item := range ch.Items {
		buf.WriteString("<item>\n")
		writeXMLTextEl(&buf, "title", item.Title)
		if item.Link != "" {
			writeXMLTextEl(&buf, "link", item.Link)
		}
		buf.WriteString("<description><![CDATA[")
		buf.WriteString(item.Description)
		buf.WriteString("]]></description>\n")
		writeXMLTextEl(&buf, "pubDate", item.PubDate.UTC().Format(time.RFC1123Z))
		buf.WriteString(`<guid isPermaLink="false">`)
		buf.WriteString(xmlEscape(item.GUID))
		buf.WriteString("</guid>\n")
		buf.WriteString("</item>\n")
	}
	buf.WriteString("</channel>\n</rss>\n")
	return buf.Bytes()
}

func writeXMLTextEl(buf *bytes.Buffer, name, value string) {
	buf.WriteString("<")
	buf.WriteString(name)
	buf.WriteString(">")
	buf.WriteString(xmlEscape(value))
	buf.WriteString("</")
	buf.WriteString(name)
	buf.WriteString(">\n")
}

func xmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(s)
}
