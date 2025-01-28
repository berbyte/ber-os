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

// RefineSchema defines the structure for story refinement responses
type RefineSchema struct {
	// Refinements contains a list of identified issues and their refinements
	Refinements []struct {
		// Issue describes the unclear or problematic aspect of the story
		Issue string `json:"issue" jsonschema:"required,description=The identified issue or unclear aspect"`
		// Clarification provides detailed explanation to address the issue
		Clarification string `json:"clarification" jsonschema:"required,description=Detailed explanation to resolve the issue"`
		// Questions lists clarifying questions that need answers
		Questions []string `json:"questions,omitempty" jsonschema:"required,description=List of questions that need to be answered"`
		// Suggestions provides recommendations for improving the story
		Suggestions []string `json:"suggestions" jsonschema:"required,description=List of suggested improvements"`
	} `json:"refinements" jsonschema:"required,description=List of refinement items for the story"`
}

type RefineSkill struct {
	agent.Skill[RefineSchema]
}

// Refine defines the skill configuration for story refinement
var Refine = RefineSkill{
	Skill: agent.Skill[RefineSchema]{
		Name:        "Refine",
		Tag:         "refine",
		Description: "Helps refine and clarify user stories by identifying gaps, asking questions, and suggesting improvements",
		Prompt: `As a Product Owner, help refine and clarify the given user story or feature.
Follow these guidelines:
- Identify any unclear or ambiguous aspects
- Ask clarifying questions
- Suggest improvements for clarity and completeness
- Focus on business value and user needs`,
		Template: `# Story Refinement

{{- range .refinements}}
## Issue
{{.issue}}

## Clarification
{{.clarification}}

{{- if .questions}}
## Questions to Address
{{- range .questions}}
- {{.}}
{{- end}}
{{- end}}

## Suggestions for Improvement
{{- range .suggestions}}
- {{.}}
{{- end}}

{{end}}`,
		LLMSchema: RefineSchema{},
	},
}
