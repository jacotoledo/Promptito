package models

import "time"

type Skill struct {
	Slug           string          `json:"slug"`
	Version        string          `json:"version"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Author         string          `json:"author,omitempty"`
	AuthorURL      string          `json:"authorUrl,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Tags           []string        `json:"tags"`
	Category       string          `json:"category"`
	PromptTemplate string          `json:"promptTemplate"`
	Variables      []Variable      `json:"variables,omitempty"`
	Examples       []Example       `json:"examples,omitempty"`
	QualityMetrics *QualityMetrics `json:"qualityMetrics,omitempty"`
	Guardrails     *Guardrails     `json:"guardrails,omitempty"`
	Ethics         *Ethics         `json:"ethics,omitempty"`
	SFIA           *SFIA           `json:"sfia,omitempty"`
	Framework      *Framework      `json:"framework,omitempty"`
	MCP            *MCP            `json:"mcp,omitempty"`
	IPTC           *IPTC           `json:"iptc,omitempty"`
}

type Variable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

type Example struct {
	Title  string            `json:"title"`
	Input  map[string]string `json:"input"`
	Output string            `json:"output,omitempty"`
}

type QualityMetrics struct {
	Accuracy     float64 `json:"accuracy,omitempty"`
	Consistency  float64 `json:"consistency,omitempty"`
	Completeness float64 `json:"completeness,omitempty"`
	AuditDate    string  `json:"auditDate,omitempty"`
	AuditNote    string  `json:"auditNote,omitempty"`
}

type Guardrails struct {
	IntendedUse  []string `json:"intendedUse,omitempty"`
	OutOfScope   []string `json:"outOfScope,omitempty"`
	Constraints  []string `json:"constraints,omitempty"`
	NegativeList []string `json:"negativeList,omitempty"`
}

type Ethics struct {
	HumanAgency    string   `json:"humanAgency,omitempty"`
	Transparency   string   `json:"transparency,omitempty"`
	BiasMitigation []string `json:"biasMitigation,omitempty"`
}

type SFIA struct {
	Level      int      `json:"level"`
	Skills     []string `json:"skills,omitempty"`
	Competency string   `json:"competency,omitempty"`
}

type Framework struct {
	Type string `json:"type"`
}

type MCP struct {
	Tools     []string `json:"tools,omitempty"`
	Resources []string `json:"resources,omitempty"`
	Servers   []string `json:"servers,omitempty"`
}

type IPTC struct {
	PromptWriterName string `json:"aiPromptWriterName,omitempty"`
	SystemUsed       string `json:"aiSystemUsed,omitempty"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

type BundleRequest struct {
	Slugs []string `json:"slugs"`
}
