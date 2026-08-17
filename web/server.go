package web

import (
	"context"
	"net/http"
	"time"

	"github.com/kataras/golog"
	"github.com/zema1/watchvuln/ent"
)

type Server struct {
	db      *ent.Client
	log     *golog.Logger
	addr    string
	version string
	server  *http.Server
}

func NewServer(db *ent.Client, addr string) *Server {
	return &Server{
		db:   db,
		log:  golog.Child("[board]"),
		addr: addr,
	}
}

func (s *Server) SetVersion(v string) *Server {
	s.version = v
	return s
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/feed.xml", s.handleFeed)
	mux.HandleFunc("/api", s.handleAPIIndex)
	mux.HandleFunc("/api/", s.handleAPIIndex)
	mux.HandleFunc("/api/vulns", s.handleAPIVulns)
	mux.HandleFunc("/api/vuln", s.handleAPIVuln)
	mux.HandleFunc("/api/stats", s.handleAPIStats)
	mux.HandleFunc("/api/sources", s.handleAPISources)
	mux.HandleFunc("/mcp", s.handleMCP)
	mux.HandleFunc("/mcp/", s.handleMCP)
	return withCORS(mux)
}

func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{
		Addr:              s.addr,
		Handler:           s.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}
	s.log.Infof("vuln intelligence board listening on http://%s/", s.addr)
	s.log.Infof("rss feed: http://%s/feed.xml", s.addr)
	s.log.Infof("rest api: http://%s/api", s.addr)
	s.log.Infof("mcp: http://%s/mcp (streamable-http, follows request host)", s.addr)
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
