package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/golog"
	"github.com/zema1/watchvuln/ent"
	"github.com/zema1/watchvuln/grab"
	"github.com/zema1/watchvuln/ent/predicate"
	"github.com/zema1/watchvuln/ent/vulninformation"
)

type Server struct {
	db     *ent.Client
	log    *golog.Logger
	addr   string
	server *http.Server
}

func NewServer(db *ent.Client, addr string) *Server {
	return &Server{
		db:   db,
		log:  golog.Child("[board]"),
		addr: addr,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/vulns", s.handleAPIVulns)
	mux.HandleFunc("/api/stats", s.handleAPIStats)
	mux.HandleFunc("/api/sources", s.handleAPISources)

	s.server = &http.Server{
		Addr:              s.addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	s.log.Infof("vuln intelligence board listening on http://%s/", s.addr)
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.server.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.server.Shutdown(shutdownCtx)
		return ctx.Err()
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(dashboardHTML))
}

func (s *Server) handleAPISources(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, grab.ListSources())
}

func (s *Server) handleAPIStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	total, err := s.db.VulnInformation.Query().Count(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	bySeverity := map[string]int{}
	for _, sev := range []string{"严重", "高危", "中危", "低危"} {
		n, err := s.db.VulnInformation.Query().Where(vulninformation.SeverityEQ(sev)).Count(ctx)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		if n > 0 {
			bySeverity[sev] = n
		}
	}
	writeJSON(w, map[string]any{
		"total":       total,
		"by_severity": bySeverity,
	})
}

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

func (s *Server) handleAPIVulns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	severity := strings.TrimSpace(r.URL.Query().Get("severity"))
	source := strings.TrimSpace(r.URL.Query().Get("source"))
	sortBy := strings.TrimSpace(r.URL.Query().Get("sort"))
	if sortBy != "update" {
		sortBy = "disclosure"
	}

	query := s.db.VulnInformation.Query()
	if severity != "" {
		query = query.Where(vulninformation.SeverityEQ(severity))
	}
	if q != "" {
		query = query.Where(vulninformation.Or(
			vulninformation.TitleContainsFold(q),
			vulninformation.CveContainsFold(q),
			vulninformation.DescriptionContainsFold(q),
		))
	}
	if source != "" {
		if meta, ok := grab.SourceMetaByID(source); ok && len(meta.HostMatches) > 0 {
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
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	switch sortBy {
	case "update":
		query = query.Order(ent.Desc(vulninformation.FieldUpdateTime))
	default:
		query = query.Order(
			ent.Desc(vulninformation.FieldDisclosure),
			ent.Desc(vulninformation.FieldUpdateTime),
		)
	}
	rows, err := query.
		Offset((page - 1) * limit).
		Limit(limit).
		All(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	items := make([]vulnItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, vulnItem{
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
		})
	}
	writeJSON(w, map[string]any{
		"total": total,
		"page":  page,
		"limit": limit,
		"sort":  sortBy,
		"items": items,
	})
}

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
