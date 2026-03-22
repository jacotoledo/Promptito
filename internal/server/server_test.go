package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jtg365/promptito/internal/models"
	"github.com/jtg365/promptito/internal/storage"
)

func setupTestServer(t *testing.T) (*Server, func()) {
	tmpDir := t.TempDir()

	skillDir := filepath.Join(tmpDir, "test-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}

	skillContent := `---
name: Test Skill
slug: test-skill
version: 1.0.0
description: A test skill
category: testing
tags:
  - test
  - unit
sfia:
  level: 3
---
# Role
You are a test skill.
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0644); err != nil {
		t.Fatal(err)
	}

	secondDir := filepath.Join(tmpDir, "another-skill")
	if err := os.MkdirAll(secondDir, 0755); err != nil {
		t.Fatal(err)
	}

	secondContent := `---
name: Another Skill
slug: another-skill
version: 1.0.0
description: Another test skill
category: dev
tags:
  - test
---
# Role
Another test.
`
	if err := os.WriteFile(filepath.Join(secondDir, "SKILL.md"), []byte(secondContent), 0644); err != nil {
		t.Fatal(err)
	}

	store, err := storage.New(storage.Config{Directory: tmpDir})
	if err != nil {
		t.Fatal(err)
	}

	srv, err := New(WithStorage(store), WithStatic(tmpDir))
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return srv, cleanup
}

func TestHandleHealth(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp models.APIResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}
}

func TestHandleList(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/skills", nil)
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp models.APIResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}

	skills, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatal("expected data to be array")
	}

	if len(skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(skills))
	}
}

func TestHandleGet(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name       string
		slug       string
		wantStatus int
		wantName   string
	}{
		{"valid slug", "test-skill", http.StatusOK, "Test Skill"},
		{"another valid slug", "another-skill", http.StatusOK, "Another Skill"},
		{"nonexistent slug", "does-not-exist", http.StatusNotFound, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/skills/"+tt.slug, nil)
			rr := httptest.NewRecorder()

			srv.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}

			if tt.wantName != "" {
				var resp models.APIResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}

				data, ok := resp.Data.(map[string]interface{})
				if !ok {
					t.Fatal("expected data to be object")
				}

				if data["name"] != tt.wantName {
					t.Errorf("expected name %s, got %v", tt.wantName, data["name"])
				}

				if data["slug"] != tt.slug {
					t.Errorf("expected slug %s, got %v", tt.slug, data["slug"])
				}
			}
		})
	}
}

func TestHandleSearch(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name       string
		query      string
		wantCount  int
		wantStatus int
	}{
		{"search by name", "test", 2, http.StatusOK},
		{"search by tag", "unit", 1, http.StatusOK},
		{"empty query", "", 0, http.StatusBadRequest},
		{"no results", "xyz123", 0, http.StatusOK},
		{"search with special chars", "../", 0, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/search?q="+tt.query, nil)
			rr := httptest.NewRecorder()

			srv.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}

			if tt.wantStatus == http.StatusOK && tt.wantCount > 0 {
				var resp models.APIResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}

				skills, ok := resp.Data.([]interface{})
				if !ok {
					t.Fatal("expected data to be array")
				}

				if len(skills) != tt.wantCount {
					t.Errorf("expected %d results, got %d", tt.wantCount, len(skills))
				}
			}
		})
	}
}

func TestHandleSearchQueryLengthLimit(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	longQuery := strings.Repeat("a", 300)
	req := httptest.NewRequest(http.MethodGet, "/api/search?q="+longQuery, nil)
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestHandleBundle(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := `{"slugs": ["test-skill", "another-skill"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/bundle", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp models.APIResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	skills, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatal("expected data to be array")
	}

	if len(skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(skills))
	}
}

func TestHandleBundleLimit(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	slugs := make([]string, 100)
	for i := range slugs {
		slugs[i] = "test-skill"
	}

	body, _ := json.Marshal(map[string]interface{}{"slugs": slugs})
	req := httptest.NewRequest(http.MethodPost, "/api/bundle", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestHandleBundleInvalidSlug(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := `{"slugs": ["../etc/passwd"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/bundle", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHandleBundleWrongMethod(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/bundle", nil)
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404 (GET not registered for /api/bundle), got %d", rr.Code)
	}
}

func TestHandleBundleLargeBody(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	largeBody := strings.Repeat("a", 2<<20)
	req := httptest.NewRequest(http.MethodPost, "/api/bundle", strings.NewReader(largeBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected status 413, got %d", rr.Code)
	}
}

func TestHandleTags(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp models.APIResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	tags, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatal("expected data to be array")
	}

	if len(tags) == 0 {
		t.Error("expected at least one tag")
	}
}

func TestHandleCategories(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp models.APIResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	categories, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatal("expected data to be array")
	}

	if len(categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(categories))
	}
}

func TestSecurityHeaders(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/skills", nil)
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	headers := []string{
		"X-Content-Type-Options",
		"Content-Type",
	}

	for _, header := range headers {
		if rr.Header().Get(header) == "" {
			t.Errorf("expected header %s to be set", header)
		}
	}
}

func TestPromptsLegacyEndpoint(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/prompts", nil)
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var prompts []map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &prompts); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(prompts) != 2 {
		t.Errorf("expected 2 prompts, got %d", len(prompts))
	}
}

func TestInvalidJSON(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := `{"slugs": [invalid]}`
	req := httptest.NewRequest(http.MethodPost, "/api/bundle", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestEmptyJSON(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/bundle", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestReadPastEOF(t *testing.T) {
	srv, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/api/bundle", nil)
	req.Body = io.NopCloser(strings.NewReader(""))
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}
