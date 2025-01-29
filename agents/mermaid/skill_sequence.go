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

type SequenceSkill struct {
	agent.Skill[MermaidSchema]
}

// SequenceDiagram defines the skill configuration for creating sequence diagrams
var SequenceDiagram = SequenceSkill{
	Skill: agent.Skill[MermaidSchema]{
		Name:        "Sequence",
		Tag:         "sequence",
		Description: "Creates sequence diagrams showing interactions between components",
		// Prompt provides instructions to the LLM for generating sequence diagrams
		Prompt: `As a Mermaid diagram expert, create a sequence diagram based on the description.
Follow these guidelines:
- Define clear participant names/actors
- Show message flows with appropriate arrows (->>, -->, ->, -)
- Include activations and deactivations where relevant
- Use notes for additional context when needed
- Add loops, alt, opt, and par blocks where appropriate
- Keep the diagram clean and readable
- Ensure valid Mermaid sequence diagram syntax
- Don't include markdown code blocks in the response
`,
		// Template defines how the response should be formatted
		Template: `### Sequence Diagram

{{.diagram_code}}

### Explanation

{{.explanation}}`,
		LLMSchema: MermaidSchema{}, // Schema for validating LLM responses
		Hooks: agent.Hooks[MermaidSchema]{
			PostLLMRequest: postLLMRequest,
		},
	},
}
