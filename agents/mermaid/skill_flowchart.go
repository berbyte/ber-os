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

type FlowchartSkill struct {
	agent.Skill[MermaidSchema]
}

var FlowchartDiagram = FlowchartSkill{
	Skill: agent.Skill[MermaidSchema]{
		Name:        "Mermaid",                                                // Display name of the skill
		Tag:         "mermaid",                                                // Unique identifier for the skill
		Description: "Creates simple flowchart diagrams using Mermaid syntax", // Brief description of functionality
		// Prompt provides instructions to the LLM for generating flowcharts
		Prompt: `As a Mermaid diagram expert, create a flowchart based on the description.
Follow these guidelines:
- Use clear and concise node labels
- Create logical connections between nodes
- Use appropriate arrow types for relationships
- Keep the diagram clean and readable
- Ensure valid Mermaid flowchart syntax which can redered in markdown
- Don't include markdown code blocks in the response
`,
		// Template defines how the response should be formatted
		Template: `### Flowchart Diagram

{{.diagram_code}}

### Explanation

{{.explanation}}`,
		LLMSchema: MermaidSchema{}, // Schema for validating LLM responses
		Hooks: agent.Hooks[MermaidSchema]{
			PostLLMRequest: postLLMRequest,
		},
	},
}
