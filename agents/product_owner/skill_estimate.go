// Copyright 2025 BER - ber.run
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package product_owner

import "github.com/berbyte/ber-os/internal/agent"

// EstimateSchema defines the structure for story point estimation responses
type EstimateSchema struct {
	// Estimates contains a list of component-level estimates with their details
	Estimates []struct {
		// Component identifies the part of the story being estimated
		Component string `json:"component" jsonschema:"required,description=The component or feature being estimated"`
		// Complexity indicates the implementation difficulty (Low/Medium/High)
		Complexity string `json:"complexity" jsonschema:"required,enum=Low,Medium,High,description=The complexity level of implementation"`
		// StoryPoints represents the estimated effort using Fibonacci scale
		StoryPoints float64 `json:"story_points" jsonschema:"required,description=Story point estimate using Fibonacci scale"`
		// Rationale explains the reasoning behind the estimate
		Rationale string `json:"rationale" jsonschema:"required,description=Explanation of how the estimate was determined"`
		// Assumptions lists any key assumptions made during estimation
		Assumptions string `json:"assumptions,omitempty" jsonschema:"required,description=Key assumptions that influenced the estimate"`
		// Dependencies identifies other components or systems this depends on
		Dependencies string `json:"dependencies,omitempty" jsonschema:"required,description=External dependencies that may impact implementation"`
		// Risks highlights potential challenges or uncertainties
		Risks string `json:"risks,omitempty" jsonschema:"required,description=Potential risks that could affect the estimate"`
	} `json:"estimates" jsonschema:"required,description=List of component estimates"`
}

type EstimateSkill struct {
	agent.Skill[EstimateSchema]
}

// Estimate defines the skill configuration for story point estimation
var Estimate = EstimateSkill{
	Skill: agent.Skill[EstimateSchema]{
		Name:        "Estimate",
		Tag:         "estimate",
		Description: "Provides detailed story point estimates for user stories or features with supporting rationale",
		Prompt: `As a Product Owner, provide story point estimates for the given user story or feature.
Follow these guidelines:
- Break down into components if needed
- Assess complexity (Low/Medium/High)
- Provide story points using Fibonacci scale
- Include rationale for estimates
- Document key assumptions
- Identify dependencies and risks`,
		Template: `# Story Point Estimates

{{- range .estimates}}
## Component: {{.component}}

### Complexity Level
{{.complexity}}

### Story Points
{{.story_points}}

### Rationale
{{.rationale}}

{{- if .assumptions}}
### Assumptions
{{.assumptions}}
{{- end}}

{{- if .dependencies}}
### Dependencies
{{.dependencies}}
{{- end}}

{{- if .risks}}
### Risks
{{.risks}}
{{- end}}

{{end}}`,
		LLMSchema: EstimateSchema{},
	},
}
