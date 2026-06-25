package web

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zema1/watchvuln/ent"
)

func TestRenderRSS(t *testing.T) {
	out := string(renderRSS(feedChannel{
		Title:       "WatchVuln 漏洞情报",
		Link:        "http://192.168.1.100:8766/",
		Description: "高价值漏洞推送订阅（最近50条已推送）",
		FeedURL:     "http://192.168.1.100:8766/feed.xml",
		Items: []feedItem{
			{
				Title:       "[CVE-2024-0001] Test Vuln",
				Link:        "https://example.com/vuln/1",
				Description: "<p><strong>等级：</strong>高危</p>",
				PubDate:     time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC),
				GUID:        "avd:CVE-2024-0001",
			},
		},
	}))

	assert.Contains(t, out, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, out, "<rss version=\"2.0\">")
	assert.Contains(t, out, "<title>WatchVuln 漏洞情报</title>")
	assert.Contains(t, out, "<link>https://example.com/vuln/1</link>")
	assert.Contains(t, out, "<guid isPermaLink=\"false\">avd:CVE-2024-0001</guid>")
	assert.Contains(t, out, "<![CDATA[<p><strong>等级：</strong>高危</p>]]>")
}

func TestXMLEscape(t *testing.T) {
	assert.Equal(t, "a &amp; b &lt; c", xmlEscape("a & b < c"))
}

func TestFeedItemTitle(t *testing.T) {
	row := &ent.VulnInformation{Title: "RCE", Cve: "CVE-2024-1"}
	assert.Equal(t, "[CVE-2024-1] RCE", feedItemTitle(row))

	row.Cve = ""
	assert.Equal(t, "RCE", feedItemTitle(row))
}

func TestFeedItemDescription(t *testing.T) {
	row := &ent.VulnInformation{
		Severity:    "严重",
		Cve:         "CVE-2024-1",
		Disclosure:  "2024-01-02",
		Tags:        []string{"POC公开"},
		Description: "test & <desc>",
		Solutions:   "upgrade",
	}
	desc := feedItemDescription(row)
	assert.Contains(t, desc, "严重")
	assert.Contains(t, desc, "POC公开")
	assert.Contains(t, desc, "test &amp; &lt;desc&gt;")
}

func TestParseDisclosureDate(t *testing.T) {
	tm, ok := parseDisclosureDate("2024-03-15")
	require.True(t, ok)
	assert.Equal(t, 2024, tm.Year())
	assert.Equal(t, time.March, tm.Month())
	assert.Equal(t, 15, tm.Day())

	_, ok = parseDisclosureDate("")
	assert.False(t, ok)
}

func TestTruncateRunes(t *testing.T) {
	assert.Equal(t, "abc", truncateRunes("abc", 10))
	assert.True(t, strings.HasSuffix(truncateRunes("一二三四五六", 3), "..."))
}
