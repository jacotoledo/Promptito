package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jtg365/promptito/internal/models"
	"github.com/jtg365/promptito/internal/parser"
)

var validSlugRe = regexp.MustCompile(`^[a-z0-9][-a-z0-9]*$`)

type Storage interface {
	List() ([]*models.Skill, error)
	Get(slug string) (*models.Skill, error)
	Search(query string) ([]*models.Skill, error)
	ByTag(tag string) ([]*models.Skill, error)
	ByCategory(category string) ([]*models.Skill, error)
	BySFIALevel(level int) ([]*models.Skill, error)
	Bundle(slugs []string) ([]*models.Skill, error)
	Reload() error
}

type FileStorage struct {
	mu     sync.RWMutex
	skills map[string]*models.Skill
	dir    string
	parser *parser.Parser
}

type Config struct {
	Directory string
}

func New(cfg Config) (*FileStorage, error) {
	if cfg.Directory == "" {
		cfg.Directory = "prompts"
	}

	absDir, err := filepath.Abs(cfg.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve directory")
	}

	info, err := os.Stat(absDir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(absDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory")
			}
		} else {
			return nil, fmt.Errorf("failed to access directory")
		}
	} else if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory")
	}

	fs := &FileStorage{
		dir:    absDir,
		parser: parser.New(),
		skills: make(map[string]*models.Skill),
	}

	if err := fs.load(); err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *FileStorage) load() error {
	entries, err := os.ReadDir(fs.dir)
	if err != nil {
		return fmt.Errorf("failed to read directory")
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		slug := entry.Name()

		if !isValidSlug(slug) {
			continue
		}

		skillPath := filepath.Join(fs.dir, slug, "SKILL.md")

		info, err := os.Stat(skillPath)
		if err != nil {
			continue
		}

		skill, err := fs.parser.ParseFile(skillPath)
		if err != nil {
			continue
		}

		if info.ModTime().After(skill.UpdatedAt) {
			skill.UpdatedAt = info.ModTime()
		}

		fs.skills[slug] = skill
	}

	return nil
}

func (fs *FileStorage) List() ([]*models.Skill, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	skills := make([]*models.Skill, 0, len(fs.skills))
	for _, skill := range fs.skills {
		skills = append(skills, skill)
	}

	return skills, nil
}

func (fs *FileStorage) Get(slug string) (*models.Skill, error) {
	if !isValidSlug(slug) {
		return nil, fmt.Errorf("invalid slug format")
	}

	fs.mu.RLock()
	defer fs.mu.RUnlock()

	skill, ok := fs.skills[slug]
	if !ok {
		return nil, fmt.Errorf("skill not found")
	}

	return skill, nil
}

func (fs *FileStorage) Search(query string) ([]*models.Skill, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	query = sanitizeSearchQuery(query)
	if query == "" {
		return []*models.Skill{}, nil
	}

	var results []*models.Skill

	for _, skill := range fs.skills {
		if contains(skill.Name, query) ||
			contains(skill.Description, query) ||
			containsAny(skill.Tags, query) ||
			contains(skill.Category, query) {
			results = append(results, skill)
		}
	}

	return results, nil
}

func (fs *FileStorage) ByTag(tag string) ([]*models.Skill, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var results []*models.Skill
	tag = sanitizeSearchQuery(tag)

	for _, skill := range fs.skills {
		for _, t := range skill.Tags {
			if strings.ToLower(t) == tag {
				results = append(results, skill)
				break
			}
		}
	}

	return results, nil
}

func (fs *FileStorage) ByCategory(category string) ([]*models.Skill, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var results []*models.Skill
	category = strings.ToLower(sanitizeSearchQuery(category))

	for _, skill := range fs.skills {
		if strings.ToLower(skill.Category) == category {
			results = append(results, skill)
		}
	}

	return results, nil
}

func (fs *FileStorage) BySFIALevel(level int) ([]*models.Skill, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var results []*models.Skill

	for _, skill := range fs.skills {
		if skill.SFIA != nil && skill.SFIA.Level == level {
			results = append(results, skill)
		}
	}

	return results, nil
}

func (fs *FileStorage) Bundle(slugs []string) ([]*models.Skill, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var results []*models.Skill
	seen := make(map[string]bool)

	for _, slug := range slugs {
		if !isValidSlug(slug) {
			continue
		}

		if seen[slug] {
			continue
		}
		seen[slug] = true

		skill, ok := fs.skills[slug]
		if !ok {
			continue
		}
		results = append(results, skill)
	}

	return results, nil
}

func (fs *FileStorage) Reload() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.skills = make(map[string]*models.Skill)
	return fs.load()
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func containsAny(ss []string, substr string) bool {
	for _, s := range ss {
		if contains(s, substr) {
			return true
		}
	}
	return false
}

func isValidSlug(slug string) bool {
	if len(slug) == 0 || len(slug) > 100 {
		return false
	}
	return validSlugRe.MatchString(slug)
}

func sanitizeSearchQuery(query string) string {
	query = strings.TrimSpace(query)
	if len(query) > 200 {
		query = query[:200]
	}

	var result strings.Builder
	for _, r := range query {
		if r == '%' || r == '_' || r == '\\' || r == '"' || r == '\'' {
			continue
		}
		if r < 32 || r == 127 {
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

type healthStatus struct {
	uptime       time.Time
	skillsLoaded int
}
