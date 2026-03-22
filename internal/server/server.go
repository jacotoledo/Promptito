package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jtg365/promptito/internal/models"
	"github.com/jtg365/promptito/internal/storage"
)

const (
	maxBodySize      = 1 << 20 // 1MB max request body
	maxSearchResults = 100     // Max results per search
	maxBundleSize    = 50      // Max items in bundle
	maxSlugLength    = 100     // Max slug length
	rateLimitWindow  = 1       // Window in seconds
	rateLimitMaxReqs = 100     // Max requests per window
)

var validSlugRe = regexp.MustCompile(`^[a-z0-9][-a-z0-9]*$`)

type Server struct {
	mux       *http.ServeMux
	storage   storage.Storage
	dir       string
	static    string
	rateLimit *RateLimiter
}

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	var validReqs []time.Time
	for _, t := range rl.requests[ip] {
		if t.After(windowStart) {
			validReqs = append(validReqs, t)
		}
	}

	if len(validReqs) >= rl.limit {
		rl.requests[ip] = validReqs
		return false
	}

	validReqs = append(validReqs, now)
	rl.requests[ip] = validReqs
	return true
}

type Option func(*Server)

func WithStorage(s storage.Storage) Option {
	return func(srv *Server) {
		srv.storage = s
	}
}

func WithStatic(dir string) Option {
	return func(srv *Server) {
		srv.static = dir
	}
}

func WithPromptDir(dir string) Option {
	return func(srv *Server) {
		srv.dir = dir
	}
}

func New(opts ...Option) (*Server, error) {
	srv := &Server{
		mux:       http.NewServeMux(),
		rateLimit: NewRateLimiter(rateLimitMaxReqs, time.Duration(rateLimitWindow)*time.Second),
	}

	for _, opt := range opts {
		opt(srv)
	}

	srv.registerRoutes()
	return srv, nil
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /api/skills", s.handleList)
	s.mux.HandleFunc("GET /api/skills/", s.handleGet)
	s.mux.HandleFunc("GET /api/search", s.handleSearch)
	s.mux.HandleFunc("GET /api/tags", s.handleTags)
	s.mux.HandleFunc("GET /api/categories", s.handleCategories)
	s.mux.HandleFunc("POST /api/bundle", s.handleBundle)
	s.mux.HandleFunc("GET /api/prompts", s.handlePromptsLegacy)
	s.mux.HandleFunc("GET /health", s.handleHealth)

	if s.static != "" {
		absStatic, err := filepath.Abs(s.static)
		if err == nil {
			s.mux.Handle("GET /", s.serveStatic(absStatic))
		}
	} else {
		s.mux.HandleFunc("GET /", s.handleIndex)
	}
}

func (s *Server) serveStatic(absDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		cleanPath := filepath.Clean(r.URL.Path)
		if strings.HasPrefix(cleanPath, "..") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		absPath := filepath.Join(absDir, cleanPath)
		if cleanPath == "/" || cleanPath == "" {
			absPath = filepath.Join(absDir, "index.html")
		}

		http.ServeFile(w, r, absPath)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.rateLimit.Allow(r.RemoteAddr) {
		http.Error(w, `{"success":false,"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
		return
	}

	s.setSecurityHeaders(w)
	s.mux.ServeHTTP(w, r)
}

func (s *Server) setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; script-src 'self'; object-src 'none'")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
	w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.json(w, http.StatusOK, models.APIResponse{Success: true})
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	skills, err := s.storage.List()
	if err != nil {
		s.error(w, http.StatusInternalServerError, err)
		return
	}

	if len(skills) > maxSearchResults {
		skills = skills[:maxSearchResults]
	}

	s.json(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    skills,
		Meta:    &models.Meta{Total: len(skills)},
	})
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	slug = strings.TrimSuffix(slug, "/")

	if !isValidSlug(slug) {
		http.Error(w, "Invalid slug format", http.StatusBadRequest)
		return
	}

	skill, err := s.storage.Get(slug)
	if err != nil {
		http.Error(w, "Skill not found", http.StatusNotFound)
		return
	}

	s.json(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    skill,
	})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, `{"success":false,"error":"query parameter 'q' is required"}`, http.StatusBadRequest)
		return
	}

	if len(query) > 200 {
		query = query[:200]
	}

	skills, err := s.storage.Search(query)
	if err != nil {
		s.error(w, http.StatusInternalServerError, err)
		return
	}

	if len(skills) > maxSearchResults {
		skills = skills[:maxSearchResults]
	}

	s.json(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    skills,
		Meta:    &models.Meta{Total: len(skills)},
	})
}

func (s *Server) handleTags(w http.ResponseWriter, r *http.Request) {
	skills, err := s.storage.List()
	if err != nil {
		s.error(w, http.StatusInternalServerError, err)
		return
	}

	tagCount := make(map[string]int)
	for _, skill := range skills {
		for _, tag := range skill.Tags {
			tagCount[tag]++
		}
	}

	tags := make([]struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}, 0, len(tagCount))
	for name, count := range tagCount {
		tags = append(tags, struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		}{name, count})
	}

	s.json(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    tags,
	})
}

func (s *Server) handleCategories(w http.ResponseWriter, r *http.Request) {
	skills, err := s.storage.List()
	if err != nil {
		s.error(w, http.StatusInternalServerError, err)
		return
	}

	catCount := make(map[string]int)
	for _, skill := range skills {
		catCount[skill.Category]++
	}

	categories := make([]struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}, 0, len(catCount))
	for name, count := range catCount {
		categories = append(categories, struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		}{name, count})
	}

	s.json(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    categories,
	})
}

func (s *Server) handleBundle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"success":false,"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodySize))
	if err != nil {
		http.Error(w, `{"success":false,"error":"request too large"}`, http.StatusRequestEntityTooLarge)
		return
	}
	defer r.Body.Close()

	var req models.BundleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"success":false,"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	for _, slug := range req.Slugs {
		if !isValidSlug(slug) {
			http.Error(w, `{"success":false,"error":"invalid slug in bundle"}`, http.StatusBadRequest)
			return
		}
	}

	if len(req.Slugs) > maxBundleSize {
		req.Slugs = req.Slugs[:maxBundleSize]
	}

	skills, err := s.storage.Bundle(req.Slugs)
	if err != nil {
		s.error(w, http.StatusInternalServerError, err)
		return
	}

	s.json(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    skills,
		Meta:    &models.Meta{Total: len(skills)},
	})
}

func (s *Server) handlePromptsLegacy(w http.ResponseWriter, r *http.Request) {
	skills, err := s.storage.List()
	if err != nil {
		s.error(w, http.StatusInternalServerError, err)
		return
	}

	if len(skills) > maxSearchResults {
		skills = skills[:maxSearchResults]
	}

	prompts := make([]map[string]interface{}, 0, len(skills))
	for _, skill := range skills {
		prompts = append(prompts, map[string]interface{}{
			"slug":           skill.Slug,
			"name":           skill.Name,
			"description":    skill.Description,
			"promptTemplate": skill.PromptTemplate,
		})
	}

	s.json(w, http.StatusOK, prompts)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join(s.static, "index.html")
	http.ServeFile(w, r, indexPath)
}

func (s *Server) json(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("json encode error")
	}
}

func (s *Server) error(w http.ResponseWriter, status int, err error) {
	http.Error(w, fmt.Sprintf(`{"success":false,"error":"%s"}`, http.StatusText(status)), status)
}

func isValidSlug(slug string) bool {
	if len(slug) == 0 || len(slug) > maxSlugLength {
		return false
	}
	if slug[len(slug)-1] == '-' {
		return false
	}
	return validSlugRe.MatchString(slug)
}

type TimeoutHandler struct {
	timeout time.Duration
	handler http.Handler
}

func (h *TimeoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()
	h.handler.ServeHTTP(w, r.WithContext(ctx))
}
