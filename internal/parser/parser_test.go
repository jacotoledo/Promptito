package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jtg365/promptito/internal/models"
)

func TestParseFile(t *testing.T) {
	tmpDir := t.TempDir()
	skillPath := filepath.Join(tmpDir, "test-skill", "SKILL.md")

	skillContent := `---
name: Test Skill
version: 1.0.0
description: A test skill for unit testing
category: testing
tags:
  - test
  - unit
sfia:
  level: 3
---
# Role
You are a test skill.

# Instructions
This is a test.
`

	if err := os.MkdirAll(filepath.Dir(skillPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skillPath, []byte(skillContent), 0644); err != nil {
		t.Fatal(err)
	}

	p := New()
	skill, err := p.ParseFile(skillPath)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if skill.Name != "Test Skill" {
		t.Errorf("expected name 'Test Skill', got '%s'", skill.Name)
	}
	if skill.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", skill.Version)
	}
	if len(skill.Tags) != 2 {
		t.Logf("tags: %v", skill.Tags)
		t.Errorf("expected 2 tags, got %d", len(skill.Tags))
	}
	if skill.SFIA == nil || skill.SFIA.Level != 3 {
		t.Errorf("expected SFIA level 3")
	}
	if skill.PromptTemplate == "" {
		t.Error("expected non-empty prompt template")
	}
}

func TestParseInvalidFrontmatter(t *testing.T) {
	p := New()

	_, err := p.Parse("No frontmatter here", "test.md")
	if err == nil {
		t.Error("expected error for missing frontmatter")
	}
}

func TestParseEmptyContent(t *testing.T) {
	p := New()

	_, err := p.Parse("", "test.md")
	if err == nil {
		t.Error("expected error for empty content")
	}
}

func TestTokenize(t *testing.T) {
	p := New()

	input := `name: Test
tags:
  - tag1
  - tag2
nested:
  key: value`

	tokens := p.tokenize(input)

	if len(tokens) == 0 {
		t.Error("expected tokens")
	}
}

func TestAssignField(t *testing.T) {
	p := New()

	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{"valid name", "name", "Test", false},
		{"valid version", "version", "1.0.0", false},
		{"valid description", "description", "Test desc", false},
		{"valid category", "category", "testing", false},
		{"tags is list", "tags", "value", true},
		{"unknown field", "unknown", "value", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.assignField(&models.Skill{}, tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("assignField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidSlug(t *testing.T) {
	tests := []struct {
		slug  string
		valid bool
	}{
		{"valid-slug", true},
		{"slug123", true},
		{"a", true},
		{"123", true},
		{"", false},
		{"-invalid", false},
		{"invalid-", false},
		{"has space", false},
		{"has/slash", false},
		{"UPPERCASE", false},
		{"mixed-Case", false},
	}

	for _, tt := range tests {
		t.Run(tt.slug, func(t *testing.T) {
			if got := isValidSlug(tt.slug); got != tt.valid {
				t.Errorf("isValidSlug(%q) = %v, want %v", tt.slug, got, tt.valid)
			}
		})
	}
}
