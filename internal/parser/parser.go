package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jtg365/promptito/internal/models"
)

var frontmatterRE = regexp.MustCompile(`(?s)^---\n(.*?)\n---`)
var kvRE = regexp.MustCompile(`^(\w+):\s*(.*)$`)
var nestedKeyRE = regexp.MustCompile(`^(\s+)(\w+):\s*(.*)$`)
var listItemRE = regexp.MustCompile(`^\s+-\s+(.*)$`)
var validSlugRe = regexp.MustCompile(`^[a-z0-9][-a-z0-9]*$`)

func isValidSlug(slug string) bool {
	if len(slug) == 0 || len(slug) > 100 {
		return false
	}
	if slug[len(slug)-1] == '-' {
		return false
	}
	return validSlugRe.MatchString(slug)
}

type ParseError struct {
	Path string
	Err  error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error in %s: %v", e.Path, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFile(path string) (*models.Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ParseError{Path: path, Err: err}
	}

	return p.Parse(string(data), path)
}

func (p *Parser) Parse(content string, path string) (*models.Skill, error) {
	matches := frontmatterRE.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil, &ParseError{Path: path, Err: fmt.Errorf("missing frontmatter")}
	}

	frontmatter := matches[1]
	skill := &models.Skill{}

	tokens := p.tokenize(frontmatter)
	if err := p.parseIntoSkill(tokens, skill); err != nil {
		return nil, &ParseError{Path: path, Err: err}
	}

	skill.PromptTemplate = strings.TrimSpace(frontmatterRE.ReplaceAllString(content, ""))

	if skill.Slug == "" {
		skill.Slug = strings.TrimSuffix(filepath.Base(filepath.Dir(path)), ".md")
	}

	if skill.CreatedAt.IsZero() {
		skill.CreatedAt = time.Now()
	}
	if skill.UpdatedAt.IsZero() {
		skill.UpdatedAt = time.Now()
	}

	return skill, nil
}

type token struct {
	key      string
	value    string
	indent   int
	isNested bool
	isList   bool
	listVal  string
}

func (p *Parser) tokenize(input string) []token {
	var tokens []token
	scanner := bufio.NewScanner(strings.NewReader(input))
	var currentKey string
	var currentIndent int

	for scanner.Scan() {
		line := scanner.Text()

		if listItemRE.MatchString(line) {
			matches := listItemRE.FindStringSubmatch(line)
			tokens = append(tokens, token{
				key:     currentKey,
				isList:  true,
				listVal: strings.Trim(matches[1], "\""),
			})
			continue
		}

		if matches := nestedKeyRE.FindStringSubmatch(line); matches != nil {
			indent := len(matches[1])
			key := matches[2]
			value := strings.Trim(matches[3], "\"")

			if currentKey != "" && indent > currentIndent {
				tokens = append(tokens, token{
					key:      currentKey + "." + key,
					value:    value,
					indent:   indent,
					isNested: true,
				})
			} else {
				tokens = append(tokens, token{key: key, value: value, indent: indent})
			}
			currentKey = key
			currentIndent = indent
			continue
		}

		if matches := kvRE.FindStringSubmatch(line); matches != nil {
			key := matches[1]
			value := strings.Trim(matches[2], "\"")

			tokens = append(tokens, token{key: key, value: value})
			currentKey = key
			currentIndent = 0
		}
	}

	return tokens
}

func (p *Parser) parseIntoSkill(tokens []token, skill *models.Skill) error {
	var listKey string
	var listValues []string

	for i, tok := range tokens {
		if tok.isList {
			listValues = append(listValues, tok.listVal)
			continue
		}

		if len(listValues) > 0 {
			p.assignListValue(skill, listKey, listValues)
			listValues = nil
		}

		listKey = tok.key
		if tok.value != "" {
			listValues = []string{tok.value}
		}

		if strings.Contains(tok.key, ".") {
			parts := strings.Split(tok.key, ".")
			if len(parts) == 2 {
				p.parseNestedField(skill, parts[0], parts[1], tok.value)
			}
			continue
		}

		if err := p.assignField(skill, tok.key, tok.value); err != nil {
			if i+1 < len(tokens) && !tokens[i+1].isList {
				continue
			}
		}
	}

	if len(listValues) > 0 {
		p.assignListValue(skill, listKey, listValues)
	}

	return nil
}

func (p *Parser) assignField(skill *models.Skill, key, value string) error {
	switch key {
	case "name":
		skill.Name = value
	case "version":
		skill.Version = value
	case "description":
		skill.Description = value
	case "author":
		skill.Author = value
	case "authorUrl":
		skill.AuthorURL = value
	case "category":
		skill.Category = value
	case "createdAt":
		t, err := time.Parse(time.RFC3339, value)
		if err == nil {
			skill.CreatedAt = t
		}
	case "updatedAt":
		t, err := time.Parse(time.RFC3339, value)
		if err == nil {
			skill.UpdatedAt = t
		}
	case "tags":
		return fmt.Errorf("tags is a list")
	default:
		return fmt.Errorf("unknown field: %s", key)
	}
	return nil
}

func (p *Parser) assignListValue(skill *models.Skill, key string, values []string) {
	switch key {
	case "tags":
		skill.Tags = values
	case "variables.name":
		for i := 0; i < len(values); i++ {
			if i >= len(skill.Variables) {
				skill.Variables = append(skill.Variables, models.Variable{})
			}
			skill.Variables[i].Name = values[i]
		}
	case "variables.description":
		for i := 0; i < len(values); i++ {
			if i >= len(skill.Variables) {
				skill.Variables = append(skill.Variables, models.Variable{})
			}
			skill.Variables[i].Description = values[i]
		}
	case "variables.default":
		for i := 0; i < len(values); i++ {
			if i >= len(skill.Variables) {
				skill.Variables = append(skill.Variables, models.Variable{})
			}
			skill.Variables[i].Default = values[i]
		}
	case "examples.title":
		for i := 0; i < len(values); i++ {
			if i >= len(skill.Examples) {
				skill.Examples = append(skill.Examples, models.Example{})
			}
			skill.Examples[i].Title = values[i]
		}
	case "qualityMetrics.accuracy":
		if v, err := strconv.ParseFloat(values[0], 64); err == nil {
			if skill.QualityMetrics == nil {
				skill.QualityMetrics = &models.QualityMetrics{}
			}
			skill.QualityMetrics.Accuracy = v
		}
	case "qualityMetrics.consistency":
		if v, err := strconv.ParseFloat(values[0], 64); err == nil {
			if skill.QualityMetrics == nil {
				skill.QualityMetrics = &models.QualityMetrics{}
			}
			skill.QualityMetrics.Consistency = v
		}
	case "qualityMetrics.completeness":
		if v, err := strconv.ParseFloat(values[0], 64); err == nil {
			if skill.QualityMetrics == nil {
				skill.QualityMetrics = &models.QualityMetrics{}
			}
			skill.QualityMetrics.Completeness = v
		}
	case "qualityMetrics.auditDate":
		if skill.QualityMetrics == nil {
			skill.QualityMetrics = &models.QualityMetrics{}
		}
		skill.QualityMetrics.AuditDate = values[0]
	case "guardrails.intendedUse":
		if skill.Guardrails == nil {
			skill.Guardrails = &models.Guardrails{}
		}
		skill.Guardrails.IntendedUse = values
	case "guardrails.outOfScope":
		if skill.Guardrails == nil {
			skill.Guardrails = &models.Guardrails{}
		}
		skill.Guardrails.OutOfScope = values
	case "guardrails.constraints":
		if skill.Guardrails == nil {
			skill.Guardrails = &models.Guardrails{}
		}
		skill.Guardrails.Constraints = values
	case "guardrails.negativeList":
		if skill.Guardrails == nil {
			skill.Guardrails = &models.Guardrails{}
		}
		skill.Guardrails.NegativeList = values
	case "ethics.humanAgency":
		if skill.Ethics == nil {
			skill.Ethics = &models.Ethics{}
		}
		skill.Ethics.HumanAgency = values[0]
	case "ethics.transparency":
		if skill.Ethics == nil {
			skill.Ethics = &models.Ethics{}
		}
		skill.Ethics.Transparency = values[0]
	case "ethics.biasMitigation":
		if skill.Ethics == nil {
			skill.Ethics = &models.Ethics{}
		}
		skill.Ethics.BiasMitigation = values
	case "sfia.level":
		if v, err := strconv.Atoi(values[0]); err == nil {
			if skill.SFIA == nil {
				skill.SFIA = &models.SFIA{}
			}
			skill.SFIA.Level = v
		}
	case "sfia.skills":
		if skill.SFIA == nil {
			skill.SFIA = &models.SFIA{}
		}
		skill.SFIA.Skills = values
	case "sfia.competency":
		if skill.SFIA == nil {
			skill.SFIA = &models.SFIA{}
		}
		skill.SFIA.Competency = values[0]
	case "framework.type":
		if skill.Framework == nil {
			skill.Framework = &models.Framework{}
		}
		skill.Framework.Type = values[0]
	case "mcp.tools":
		if skill.MCP == nil {
			skill.MCP = &models.MCP{}
		}
		skill.MCP.Tools = values
	case "mcp.resources":
		if skill.MCP == nil {
			skill.MCP = &models.MCP{}
		}
		skill.MCP.Resources = values
	case "mcp.servers":
		if skill.MCP == nil {
			skill.MCP = &models.MCP{}
		}
		skill.MCP.Servers = values
	case "iptc.aiPromptWriterName":
		if skill.IPTC == nil {
			skill.IPTC = &models.IPTC{}
		}
		skill.IPTC.PromptWriterName = values[0]
	case "iptc.aiSystemUsed":
		if skill.IPTC == nil {
			skill.IPTC = &models.IPTC{}
		}
		skill.IPTC.SystemUsed = values[0]
	}
}

func (p *Parser) parseNestedField(skill *models.Skill, parent, field, value string) {
	switch parent {
	case "sfia":
		if skill.SFIA == nil {
			skill.SFIA = &models.SFIA{}
		}
		switch field {
		case "level":
			if v, err := strconv.Atoi(value); err == nil {
				skill.SFIA.Level = v
			}
		case "competency":
			skill.SFIA.Competency = value
		}
	case "framework":
		if skill.Framework == nil {
			skill.Framework = &models.Framework{}
		}
		if field == "type" {
			skill.Framework.Type = value
		}
	case "qualityMetrics":
		if skill.QualityMetrics == nil {
			skill.QualityMetrics = &models.QualityMetrics{}
		}
		switch field {
		case "accuracy", "consistency", "completeness":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				switch field {
				case "accuracy":
					skill.QualityMetrics.Accuracy = v
				case "consistency":
					skill.QualityMetrics.Consistency = v
				case "completeness":
					skill.QualityMetrics.Completeness = v
				}
			}
		case "auditDate":
			skill.QualityMetrics.AuditDate = value
		case "auditNote":
			skill.QualityMetrics.AuditNote = value
		}
	case "guardrails":
		if skill.Guardrails == nil {
			skill.Guardrails = &models.Guardrails{}
		}
		switch field {
		case "intendedUse", "outOfScope", "constraints", "negativeList":
		}
	case "ethics":
		if skill.Ethics == nil {
			skill.Ethics = &models.Ethics{}
		}
		switch field {
		case "humanAgency":
			skill.Ethics.HumanAgency = value
		case "transparency":
			skill.Ethics.Transparency = value
		}
	case "mcp":
		if skill.MCP == nil {
			skill.MCP = &models.MCP{}
		}
		switch field {
		case "tools", "resources", "servers":
		}
	case "iptc":
		if skill.IPTC == nil {
			skill.IPTC = &models.IPTC{}
		}
		switch field {
		case "aiPromptWriterName":
			skill.IPTC.PromptWriterName = value
		case "aiSystemUsed":
			skill.IPTC.SystemUsed = value
		}
	}
}
