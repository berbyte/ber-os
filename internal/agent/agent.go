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

package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/berbyte/ber-os/internal/logger"
	"go.uber.org/zap"
)

// BaseSkill defines common interface for all agent skills
type BaseSkill interface {
	GetName() string
	GetTag() string
	GetDescription() string
	GetPrompt() string
	GetTemplate() string
	GetSkill() any  // Returns concrete skill implementation
	GetSchema() any // Returns skill's data schema
	GetHooks() Hooks[any]
	GetActions() map[string]func(context.Context, any) error
	GetValidators() map[string]func(context.Context, any) error
}

// BERAgent represents an AI agent with its capabilities
type BERAgent struct {
	Name        string      `yaml:"name" json:"name"`               // Display name
	Tag         string      `yaml:"tag" json:"tag"`                 // Unique ID
	Description string      `yaml:"description" json:"description"` // Purpose description
	Skills      []BaseSkill `yaml:"skills" json:"skills"`           // Available skills
}

// Skill represents a specific AI capability with generic type T
type Skill[T any] struct {
	Name        string                                     `yaml:"name" json:"name"`
	Tag         string                                     `yaml:"tag" json:"tag"`
	Description string                                     `yaml:"description" json:"description"`
	Prompt      string                                     `yaml:"prompt" json:"prompt"`         // LLM instruction prompt
	Template    string                                     `yaml:"template" json:"template"`     // Output template
	LLMSchema   T                                          `yaml:"llm_schema" json:"llm_schema"` // Response schema
	Actions     map[string]func(context.Context, *T) error `yaml:"actions" json:"actions"`       // Post-processing actions
	Validators  map[string]func(context.Context, *T) error `yaml:"validators" json:"validators"` // Validation functions
	Hooks       Hooks[T]                                   `yaml:"hooks" json:"hooks"`           // Lifecycle hooks
}

// BaseSkill interface implementation
func (s Skill[T]) GetName() string        { return s.Name }
func (s Skill[T]) GetTag() string         { return s.Tag }
func (s Skill[T]) GetDescription() string { return s.Description }
func (s Skill[T]) GetPrompt() string      { return s.Prompt }
func (s Skill[T]) GetTemplate() string    { return s.Template }
func (s Skill[T]) GetSkill() any          { return s }
func (s Skill[T]) GetSchema() any         { return s.LLMSchema }
func (s Skill[T]) GetHooks() Hooks[any] {
	return Hooks[any]{
		PreLLMRequest:  convertHookToAny(s.Hooks.PreLLMRequest),
		PostLLMRequest: convertHookToAny(s.Hooks.PostLLMRequest),
		PreValidate:    convertHookToAny(s.Hooks.PreValidate),
		PostValidate:   convertHookToAny(s.Hooks.PostValidate),
		PreAction:      convertHookToAny(s.Hooks.PreAction),
		PostAction:     convertHookToAny(s.Hooks.PostAction),
	}
}

// Convert typed actions to interface{} actions
func (s Skill[T]) GetActions() map[string]func(context.Context, any) error {
	actions := make(map[string]func(context.Context, any) error)
	for k, v := range s.Actions {
		actions[k] = convertFunc(v)
	}
	return actions
}

// Convert typed validators to interface{} validators
func (s Skill[T]) GetValidators() map[string]func(context.Context, any) error {
	validators := make(map[string]func(context.Context, any) error)
	for k, v := range s.Validators {
		validators[k] = convertFunc(v)
	}
	return validators
}

// Hooks defines lifecycle hooks for skill execution
type Hooks[T any] struct {
	PreLLMRequest  func(context.Context, *T) error `yaml:"pre_llm_request" json:"pre_llm_request"`   // Before LLM call
	PostLLMRequest func(context.Context, *T) error `yaml:"post_llm_request" json:"post_llm_request"` // After LLM call
	PreValidate    func(context.Context, *T) error `yaml:"pre_validate" json:"pre_validate"`         // Before validation
	PostValidate   func(context.Context, *T) error `yaml:"post_validate" json:"post_validate"`       // After validation
	PreAction      func(context.Context, *T) error `yaml:"pre_action" json:"pre_action"`             // Before actions
	PostAction     func(context.Context, *T) error `yaml:"post_action" json:"post_action"`           // After actions
}

// ConvertToSchema converts interface{} to specific type T
func ConvertToSchema[T any](data interface{}) (T, error) {
	var result T

	// Direct type conversion if possible
	if typed, ok := data.(T); ok {
		return typed, nil
	}

	// JSON conversion fallback
	jsonData, err := json.Marshal(data)
	if err != nil {
		return result, fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := json.Unmarshal(jsonData, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return result, nil
}

// Convert typed function to interface{} function
func convertFunc[T any](f func(context.Context, *T) error) func(context.Context, any) error {
	return func(ctx context.Context, a any) error {
		if ptr, ok := a.(*any); ok {
			a = *ptr
		}

		t, err := ConvertToSchema[T](a)
		if err != nil {
			return fmt.Errorf("failed to convert type: %w", err)
		}
		return f(ctx, &t)
	}
}

// Convert typed hook to interface{} hook
func convertHookToAny[T any](f func(context.Context, *T) error) func(context.Context, *any) error {
	if f == nil {
		return nil
	}
	return func(ctx context.Context, a *any) error {
		t, err := ConvertToSchema[T](*a)
		if err != nil {
			return fmt.Errorf("failed to convert type: %w", err)
		}

		err = f(ctx, &t)
		if err != nil {
			return err
		}

		*a = t // Update original data
		return nil
	}
}

// AgentRegistry manages registered agents
type AgentRegistry struct {
	agents []BERAgent
}

// NewAgentRegistry creates empty registry
func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents: []BERAgent{},
	}
}

// Global registry instance
var globalRegistry = NewAgentRegistry()

// RegisterAgent adds agent to global registry
func RegisterAgent(agent BERAgent) {
	log := logger.GetLogger()
	log.Info("registering agent",
		zap.String("name", agent.Name),
		zap.String("tag", agent.Tag),
	)
	globalRegistry.agents = append(globalRegistry.agents, agent)
}

// GetRegisteredAgents returns all registered agents
func GetRegisteredAgents() []BERAgent {
	return globalRegistry.agents
}

// GetAgentTags returns list of all agent tags
func GetAgentTags() []string {
	tags := []string{}
	for _, agent := range globalRegistry.agents {
		tags = append(tags, agent.Tag)
	}
	return tags
}

// GetAgentByTag finds agent by tag ID
func GetAgentByTag(tag string) *BERAgent {
	for _, agent := range globalRegistry.agents {
		if agent.Tag == tag {
			return &agent
		}
	}
	return nil
}
