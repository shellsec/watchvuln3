package ctrl

import (
	"fmt"
	"strings"

	"github.com/kataras/golog"
)

// LogStartupSummary prints a concise configuration checklist before running.
func LogStartupSummary(config *WatchVulnAppConfig) {
	log := golog.Child("[config]")
	log.Infof("WatchVuln %s", config.Version)
	log.Infof("database: %s", config.DBConn)
	log.Infof("interval: %s", config.Interval)
	log.Infof("sources (%d): %s", len(config.Sources), strings.Join(config.Sources, ", "))
	log.Infof("filters: no_filter=%v cve_filter=%v whitelist=%d blacklist=%d",
		config.NoFilter, config.EnableCVEFilter != nil && *config.EnableCVEFilter,
		len(config.WhiteKeywords), len(config.BlackKeywords))
	if config.NoStartMessage != nil && *config.NoStartMessage {
		log.Infof("startup push message: disabled")
	} else {
		log.Infof("startup push message: enabled (use --no-start-message / NO_START_MESSAGE=true to disable)")
	}
	if config.WebAddr != "" {
		log.Infof("vuln board: http://%s/", config.WebAddr)
	} else {
		log.Infof("vuln board: disabled (use --web-addr 127.0.0.1:8765 or `watchvuln board`)")
	}
	summary := DescribePushers(config.Pusher)
	if summary == "" {
		log.Warn("pushers: none configured yet")
	} else {
		log.Infof("pushers:\n%s", summary)
	}
	if len(config.Pusher) > 1 {
		log.Infof("multi-pusher: enabled via config (same channel type supported, e.g. multiple dingding)")
	}
}

// DescribePushers returns a short summary of configured push channels.
func DescribePushers(pushers []map[string]string) string {
	if len(pushers) == 0 {
		return ""
	}
	counts := map[string]int{}
	for _, p := range pushers {
		t := strings.ToLower(strings.TrimSpace(p["type"]))
		if t == "" {
			t = "unknown"
		}
		counts[t]++
	}
	var parts []string
	for t, n := range counts {
		if n > 1 {
			parts = append(parts, fmt.Sprintf("%s x%d", t, n))
		} else {
			parts = append(parts, t)
		}
	}
	return "  - " + strings.Join(parts, "\n  - ")
}
