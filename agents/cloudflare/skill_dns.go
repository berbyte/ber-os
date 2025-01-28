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

package cloudflare

import (
	"context"

	"github.com/berbyte/ber-os/internal/agent"
)

type DNSSkill struct {
	agent.Skill[DNSRecordSchema]
}

type DNSRecordSchema struct {
	Records []struct {
		Type    string   `json:"type" jsonschema:"required,description=DNS record type: A, CNAME, MX, TXT, etc."`
		Name    string   `json:"name" jsonschema:"required,description=The DNS record name (hostname)"`
		Content []string `json:"content" jsonschema:"required,description=List of the DNS record values"`
		TTL     int      `json:"ttl" jsonschema:"required,description=Time to live in seconds"`
		Proxied bool     `json:"proxied" jsonschema:"required,description=Whether the record is proxied through Cloudflare"`
	} `json:"records" jsonschema:"required,description=List of DNS records to manage"`
}

var DNSManagement = DNSSkill{
	Skill: agent.Skill[DNSRecordSchema]{
		Name:        "DNS Management",
		Tag:         "dns",
		Description: "Manages DNS records and provides configuration validation",
		Prompt: `As a Cloudflare expert, help manage DNS records based on the provided requirements.
Follow these guidelines:
- Validate DNS record syntax.
- Use the correct record type based on the input:
  - If input specifies "CNAME" create a CNAME record.
  - If input specifies "A" create an A record.
  - If input specifies "MX" create an MX record.
  - If input specifies "TXT" create a TXT record.
- Content must match the specified type:
  - A: A list of IP addresses.
  - CNAME: A single hostname.
  - MX: A list of hostnames with priorities.
  - TXT: A list of strings.
- Include the TTL value and proxied status for eligible records.
- Always respect the record type specified in the input.`,
		Template: `# DNS Configuration

{{- range .records}}
## Record: {{.name}}
Type: {{.type}}
Content: {{.content}}
TTL: {{.ttl}}
{{- if .priority}}
Priority: {{.priority}}
{{- end}}
Proxied: {{.proxied}}

{{end}}`,
		LLMSchema: DNSRecordSchema{},
		Actions: map[string]func(context.Context, *DNSRecordSchema) error{
			"approve": createRecords,
		},
		Validators: map[string]func(context.Context, *DNSRecordSchema) error{
			"Valid DNS Record": validateDNSRecord,
			"Allowed Domains":  validateAllowedDomains,
		},
		Hooks: agent.Hooks[DNSRecordSchema]{
			PreLLMRequest:  preLLMRequest,
			PostLLMRequest: postLLMRequest,
			PreValidate:    preValidate,
			PostValidate:   postValidate,
			PreAction:      preAction,
			PostAction:     postAction,
		},
	},
}
