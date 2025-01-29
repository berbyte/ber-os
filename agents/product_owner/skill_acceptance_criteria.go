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

// AcceptanceCriteriaSchema defines the structure for acceptance criteria responses
type AcceptanceCriteriaSchema struct {
	// AcceptanceCriteria contains a list of test scenarios with their acceptance criteria
	AcceptanceCriteria []struct {
		// Scenario describes the test case being defined
		Scenario string `json:"scenario" jsonschema:"required,description=The title/description of the test scenario"`
		// Given describes the initial context/state
		Given string `json:"given" jsonschema:"required,description=The preconditions and setup for the scenario"`
		// When describes the action being taken
		When string `json:"when" jsonschema:"required,description=The action or event that triggers the scenario"`
		// Then lists the expected outcomes
		Then []string `json:"then" jsonschema:"required,description=List of expected outcomes and postconditions"`
		// Notes contains optional additional context or considerations
		Notes string `json:"notes,omitempty" jsonschema:"required,description=Optional additional notes or context for the scenario"`
	} `json:"acceptance_criteria" jsonschema:"required,description=List of acceptance criteria scenarios"`
}

type AcceptanceCriteriaSkill struct {
	agent.Skill[AcceptanceCriteriaSchema]
}

var AcceptanceCriteria = AcceptanceCriteriaSkill{
	Skill: agent.Skill[AcceptanceCriteriaSchema]{
		Name:        "Acceptance Criteria",
		Tag:         "acceptance_criteria",
		Description: "Creates clear, testable acceptance criteria for user stories using the Given-When-Then format",
		Prompt: `As a Product Owner, create acceptance criteria for the provided user story or feature.
Follow these guidelines:
- Use Given-When-Then format
- Make criteria specific and testable
- Include both happy path and error scenarios
- Consider edge cases
- Focus on business value and user experience`,
		Template: `# Acceptance Criteria

{{- range .acceptance_criteria}}
## Scenario: {{.scenario}}

### Given
{{.given}}

### When
{{.when}}

### Then
{{- range .then}}
- {{.}}
{{- end}}

{{- if .notes}}
### Notes
{{.notes}}
{{- end}}

{{end}}`,
		LLMSchema: AcceptanceCriteriaSchema{},
	},
}
