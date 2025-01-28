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

// init registers the Product Owner agent with its skills and capabilities
func init() {
	agent.RegisterAgent(agent.BERAgent{
		Name:        "Product Owner",                                                                         // Display name shown to users
		Tag:         "product_owner",                                                                         // Unique identifier for this agent
		Description: "Product Owner Agent: helps with story refinement, acceptance criteria, and estimation", // Purpose description
		Skills: []agent.BaseSkill{
			AcceptanceCriteria, // Skill for defining acceptance criteria
			Estimate,           // Skill for story estimation
			Refine,             // Skill for story refinement
		},
	})
}
