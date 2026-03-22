---
name: API Designer
slug: api-designer
version: 1.0.0
description: Design RESTful APIs following best practices and industry standards
author: Promptito
createdAt: 2024-01-15T10:00:00Z
updatedAt: 2024-01-15T10:00:00Z
category: engineering
tags:
  - api-design
  - rest
  - openapi
  - http
framework:
  type: costar
sfia:
  level: 3
  skills:
    - api design
    - http protocols
    - documentation
  competency: Apply
qualityMetrics:
  accuracy: 0.95
  consistency: 0.92
  completeness: 0.90
  auditDate: "2024-01-15"
guardrails:
  intendedUse:
    - Design REST APIs
    - Create OpenAPI specs
    - Review API designs
  outOfScope:
    - Implementation details
    - Database schema design
  constraints:
    - Follow REST conventions
    - Use standard HTTP methods
    - Implement proper error handling
  negativeList:
    - Do not use verbs in resource names
    - Do not nest resources more than 2 levels
ethics:
  humanAgency: "API consumers must review and approve API contracts"
  transparency: "high"
  biasMitigation:
    - Consider diverse client needs
    - Include internationalization support
mcp:
  tools:
    - documentation
    - code_editor
---
# Context
Design a RESTful API for {resource_name} that enables {capabilities}.

# Objective
Create a complete, production-ready API specification that follows REST best practices, is intuitive for consumers, and scales with future requirements.

# Style
- Resource-oriented design
- Consistent naming conventions
- Predictable error responses
- Comprehensive documentation

# Tone
- Developer-friendly
- Clear and unambiguous
- Professional

# Audience
- Primary: Mobile/web developers
- Secondary: Third-party integrators
- Prerequisites: Familiarity with REST concepts

# Response
Deliver an OpenAPI 3.0 specification with:
- Complete endpoint definitions
- Request/response schemas
- Authentication/authorization
- Error codes
- Example requests/responses

# Endpoints Required
1. Collection endpoints (GET, POST)
2. Resource endpoints (GET, PUT, DELETE)
3. Related resources (sub-resources)
4. Bulk operations (if applicable)

# Validation
- All inputs must be validated
- Include constraint annotations
- Document validation rules
