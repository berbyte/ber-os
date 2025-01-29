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

package mermaid

import "github.com/berbyte/ber-os/internal/agent"

type PieChartSkill struct {
	agent.Skill[MermaidSchema]
}

// PieChart defines the skill configuration for creating pie charts
var PieChart = PieChartSkill{
	Skill: agent.Skill[MermaidSchema]{
		Name:        "Pie Chart",                               // Display name of the skill
		Tag:         "pie",                                     // Unique identifier for the skill
		Description: "Creates pie charts using Mermaid syntax", // Brief description of functionality
		// Prompt provides instructions to the LLM for generating pie charts
		Prompt: `As a Mermaid diagram expert, create a pie chart based on the provided data.
Follow these guidelines:
- Use clear section labels
- Include percentage or value for each section
- Ensure sections add up to 100%
- Keep the chart simple and readable
- Use valid Mermaid pie chart syntax which can redered in markdown
- Don't include markdown code blocks in the response`,
		// Template defines how the response should be formatted
		Template: `### Pie Chart

{{.diagram_code}}

### Explanation

{{.explanation}}`,
		LLMSchema: MermaidSchema{}, // Schema for validating LLM responses
		Hooks: agent.Hooks[MermaidSchema]{
			PostLLMRequest: postLLMRequest,
		},
	},
}
