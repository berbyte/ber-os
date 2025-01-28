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

func init() {
	agent.RegisterAgent(agent.BERAgent{
		Name:        "Mermaid Diagram",
		Tag:         "mermaid",
		Description: "Diagram Agent: helps create and modify various types of Mermaid diagrams",
		Skills: []agent.BaseSkill{
			FlowchartDiagram,
			PieChart,
			SequenceDiagram,
		},
	})
}
