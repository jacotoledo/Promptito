package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := New(Config{Directory: tmpDir})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
}

func TestNewCreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "new-prompt-dir")

	_, err := New(Config{Directory: newDir})
	if err != nil {
		t.Fatalf("New() failed to create directory: %v", err)
	}

	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("directory was not created")
	}
}

func TestNewNotDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")

	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := New(Config{Directory: filePath})
	if err == nil {
		t.Error("expected error for non-directory path")
	}
}

func TestStorageList(t *testing.T) {
	tmpDir := t.TempDir()

	skillDir := filepath.Join(tmpDir, "test-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}

	skillContent := `---
name: Test Skill
version: 1.0.0
description: Test description
category: testing
tags:
  - test
---
# Role
Test role.
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0644); err != nil {
		t.Fatal(err)
	}

	store, err := New(Config{Directory: tmpDir})
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	skills, err := store.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skills))
	}
}

func TestStorageGet(t *testing.T) {
	tmpDir := t.TempDir()

	skillDir := filepath.Join(tmpDir, "my-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}

	skillContent := `---
name: My Skill
version: 1.0.0
description: Test
category: test
tags:
  - test
---
# Content
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0644); err != nil {
		t.Fatal(err)
	}

	store, err := New(Config{Directory: tmpDir})
	if err != nil {
		t.Fatal(err)
	}

	skill, err := store.Get("my-skill")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if skill.Name != "My Skill" {
		t.Errorf("expected name 'My Skill', got '%s'", skill.Name)
	}
}

func TestStorageGetInvalidSlug(t *testing.T) {
	tmpDir := t.TempDir()

	store, err := New(Config{Directory: tmpDir})
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.Get("../etc/passwd")
	if err == nil {
		t.Error("expected error for path traversal attempt")
	}

	_, err = store.Get("")
	if err == nil {
		t.Error("expected error for empty slug")
	}
}

func TestStorageSearch(t *testing.T) {
	tmpDir := t.TempDir()

	for _, name := range []string{"golang-skill", "rust-skill", "python-skill"} {
		skillDir := filepath.Join(tmpDir, name)
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			t.Fatal(err)
		}

		content := `---
name: ` + name + `
version: 1.0.0
description: A programming skill
category: dev
tags:
  - programming
  - ` + name + `
---
# Content
`
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	store, err := New(Config{Directory: tmpDir})
	if err != nil {
		t.Fatal(err)
	}

	results, err := store.Search("golang")
	if err != nil {
		t.Fatalf("Search() failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestStorageSearchEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	store, err := New(Config{Directory: tmpDir})
	if err != nil {
		t.Fatal(err)
	}

	results, err := store.Search("")
	if err != nil {
		t.Fatalf("Search() failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 results for empty query, got %d", len(results))
	}
}

func TestStorageBundle(t *testing.T) {
	tmpDir := t.TempDir()

	for _, name := range []string{"skill-one", "skill-two", "skill-three"} {
		skillDir := filepath.Join(tmpDir, name)
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			t.Fatal(err)
		}

		content := `---
name: ` + name + `
version: 1.0.0
description: Test
category: test
tags:
  - test
---
# Content
`
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	store, err := New(Config{Directory: tmpDir})
	if err != nil {
		t.Fatal(err)
	}

	results, err := store.Bundle([]string{"skill-one", "skill-two"})
	if err != nil {
		t.Fatalf("Bundle() failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestStorageBundleDeduplicates(t *testing.T) {
	tmpDir := t.TempDir()

	skillDir := filepath.Join(tmpDir, "skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := `---
name: skill
version: 1.0.0
description: Test
category: test
tags:
  - test
---
# Content
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	store, err := New(Config{Directory: tmpDir})
	if err != nil {
		t.Fatal(err)
	}

	results, err := store.Bundle([]string{"skill", "skill", "skill"})
	if err != nil {
		t.Fatalf("Bundle() failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result (deduplicated), got %d", len(results))
	}
}

func TestSanitizeSearchQuery(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal query", "normal query"},
		{"query with space", "query with space"},
		{"", ""},
		{"query%", "query"},
		{"query_", "query"},
		{"query\\", "query"},
		{"query\"", "query"},
		{"query'", "query"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeSearchQuery(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeSearchQuery(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeSearchQueryLength(t *testing.T) {
	longQuery := make([]byte, 500)
	for i := range longQuery {
		longQuery[i] = 'a'
	}

	result := sanitizeSearchQuery(string(longQuery))
	if len(result) > 200 {
		t.Errorf("sanitized query too long: %d chars", len(result))
	}
}
