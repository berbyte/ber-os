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

package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

// Renderer handles rendering templates for different agents
type Renderer struct{}

// NewTemplateRenderer creates a new template renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// RenderAgentTemplate renders a template from a BERAgent using the provided response data
func (tr *Renderer) RenderAgentTemplate(tpl string, response interface{}) (string, error) {
	// Convert response to map using JSON marshaling/unmarshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("[INT-TPL] Failed to marshal response: %w", err)
	}

	var templateData map[string]interface{}
	if err := json.Unmarshal(jsonData, &templateData); err != nil {
		return "", fmt.Errorf("[INT-TPL] Failed to unmarshal response: %w", err)
	}

	// Create and parse template
	tmpl, err := template.New("agent").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("[INT-TPL] Failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("[INT-TPL] Failed to execute template: %w", err)
	}

	return buf.String(), nil
}
