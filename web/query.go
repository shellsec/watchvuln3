package web

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zema1/watchvuln/ent"
	"github.com/zema1/watchvuln/ent/predicate"
	"github.com/zema1/watchvuln/ent/vulninformation"
	"github.com/zema1/watchvuln/grab"
)

type vulnItem struct {
	ID           int      `json:"id"`
	Key          string   `json:"key"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Severity     string   `json:"severity"`
	CVE          string   `json:"cve"`
	Disclosure   string   `json:"disclosure"`
	Solutions    string   `json:"solutions"`
	References   []string `json:"references"`
	Tags         []string `json:"tags"`
	GithubSearch []string `json:"github_search"`
	From         string   `json:"from"`
	Pushed       bool     `json:"pushed"`
	CreateTime   string   `json:"create_time"`
	UpdateTime   string   `json:"update_time"`
}

type vulnListParams struct {
	Page     int
	Limit    int
	Q        string
	Severity string
	Source   string
	Sort     string
}

type vulnListResult struct {
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Sort  string     `json:"sort"`
	Items []vulnItem `json:"items"`
}

func parseVulnListParams(values url.Values, defaultLimit, maxLimit int) vulnListParams {
	page, _ := strconv.Atoi(values.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(values.Get("limit"))
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	sortBy := strings.TrimSpace(values.Get("sort"))
	if sortBy != "update" {
		sortBy = "disclosure"
	}
	return vulnListParams{
		Page:     page,
		Limit:    limit,
		Q:        strings.TrimSpace(values.Get("q")),
		Severity: strings.TrimSpace(values.Get("severity")),
		Source:   strings.TrimSpace(values.Get("source")),
		Sort:     sortBy,
	}
}

func (s *Server) listVulns(ctx context.Context, p vulnListParams) (*vulnListResult, error) {
	query := s.db.VulnInformation.Query()
	if p.Severity != "" {
		query = query.Where(vulninformation.SeverityEQ(p.Severity))
	}
	if p.Q != "" {
		query = query.Where(vulninformation.Or(
			vulninformation.TitleContainsFold(p.Q),
			vulninformation.CveContainsFold(p.Q),
			vulninformation.DescriptionContainsFold(p.Q),
		))
	}
	if p.Source != "" {
		if meta, ok := grab.SourceMetaByID(p.Source); ok && len(meta.HostMatches) > 0 {
			var preds []predicate.VulnInformation
			for _, host := range meta.HostMatches {
				if host == "_KEV" {
					preds = append(preds, vulninformation.KeyContains("_KEV"))
				} else {
					preds = append(preds, vulninformation.FromContains(host))
				}
			}
			if len(preds) == 1 {
				query = query.Where(preds[0])
			} else if len(preds) > 1 {
				query = query.Where(vulninformation.Or(preds...))
			}
		}
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	switch p.Sort {
	case "update":
		query = query.Order(ent.Desc(vulninformation.FieldUpdateTime))
	default:
		query = query.Order(
			ent.Desc(vulninformation.FieldDisclosure),
			ent.Desc(vulninformation.FieldUpdateTime),
		)
	}
	rows, err := query.
		Offset((p.Page - 1) * p.Limit).
		Limit(p.Limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]vulnItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, toVulnItem(row))
	}
	return &vulnListResult{
		Total: total,
		Page:  p.Page,
		Limit: p.Limit,
		Sort:  p.Sort,
		Items: items,
	}, nil
}

func (s *Server) findVuln(ctx context.Context, id int, cve, key string) (*vulnItem, error) {
	if id > 0 {
		row, err := s.db.VulnInformation.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		item := toVulnItem(row)
		return &item, nil
	}
	query := s.db.VulnInformation.Query()
	switch {
	case key != "":
		query = query.Where(vulninformation.KeyEQ(key))
	case cve != "":
		query = query.Where(vulninformation.CveEqualFold(cve))
	default:
		return nil, fmt.Errorf("id, cve or key is required")
	}
	row, err := query.Order(ent.Desc(vulninformation.FieldUpdateTime)).First(ctx)
	if err != nil {
		return nil, err
	}
	item := toVulnItem(row)
	return &item, nil
}

func (s *Server) vulnStats(ctx context.Context) (map[string]any, error) {
	total, err := s.db.VulnInformation.Query().Count(ctx)
	if err != nil {
		return nil, err
	}
	bySeverity := map[string]int{}
	for _, sev := range []string{"严重", "高危", "中危", "低危"} {
		n, err := s.db.VulnInformation.Query().Where(vulninformation.SeverityEQ(sev)).Count(ctx)
		if err != nil {
			return nil, err
		}
		if n > 0 {
			bySeverity[sev] = n
		}
	}
	return map[string]any{
		"total":       total,
		"by_severity": bySeverity,
	}, nil
}

func toVulnItem(row *ent.VulnInformation) vulnItem {
	return vulnItem{
		ID:           row.ID,
		Key:          row.Key,
		Title:        row.Title,
		Description:  row.Description,
		Severity:     row.Severity,
		CVE:          row.Cve,
		Disclosure:   row.Disclosure,
		Solutions:    row.Solutions,
		References:   row.References,
		Tags:         row.Tags,
		GithubSearch: row.GithubSearch,
		From:         row.From,
		Pushed:       row.Pushed,
		CreateTime:   row.CreateTime.Format(time.RFC3339),
		UpdateTime:   row.UpdateTime.Format(time.RFC3339),
	}
}

func compactVulnItem(v vulnItem) map[string]any {
	return map[string]any{
		"id":          v.ID,
		"key":         v.Key,
		"title":       v.Title,
		"severity":    v.Severity,
		"cve":         v.CVE,
		"disclosure":  v.Disclosure,
		"tags":        v.Tags,
		"from":        v.From,
		"pushed":      v.Pushed,
		"update_time": v.UpdateTime,
	}
}
