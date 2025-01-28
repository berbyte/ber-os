# Product Owner Agent

A BER agent for helping with story refinement, acceptance criteria, and estimation.

## Features

- Acceptance Criteria Definition (Given-When-Then format)
- Story Point Estimation (using Fibonacci scale)
- Story Refinement and Clarification
- Component-level Analysis
- Risk and Dependency Assessment

## Prerequisites

- Access to project management tools
- Environment variables for necessary integrations

## Usage

```
go run . webhook --debug
```

The agent provides product management capabilities through the following skills:
```go
ProductManagement:
  - AcceptanceCriteria: Creates clear, testable acceptance criteria using Given-When-Then format
  - Estimate: Provides detailed story point estimates with complexity analysis
  - Refine: Identifies gaps and suggests improvements in user stories
```

## Validation

The agent performs several validations across its skills:

Acceptance Criteria:
- Test scenario completeness
- Given-When-Then format compliance
- Happy path and error scenarios
- Edge case coverage

Estimation:
- Component breakdown
- Complexity assessment (Low/Medium/High)
- Story point allocation (Fibonacci scale)
- Dependencies and risks identification

Refinement:
- Clarity and ambiguity checks
- Business value alignment
- Completeness verification
- Improvement suggestions
