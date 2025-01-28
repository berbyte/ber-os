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
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/berbyte/ber-os/internal/logger"
	"github.com/berbyte/ber-os/internal/services/memory_store"
	"github.com/berbyte/ber-os/internal/template"
	llm "github.com/berbyte/ber-os/pkg/openai"
	"go.uber.org/zap"
)

// Default error template shown to users when workflow fails
const DefaultErrorTemplate = `{{if not .success}}## ❌ Error:

{{.error}}

Try rephrasing your request 🔄 {{end}}`

// WorkflowResult contains the final output of a workflow execution
type WorkflowResult struct {
	Response string // Final response to show user
	Error    error  // Any error that occurred
}

// WorkflowInput contains all necessary data to execute a workflow
type WorkflowInput struct {
	Message      string                                     // User's input message
	Agent        *BERAgent                                  // Agent instance handling the workflow
	ChatHistory  []llm.ChatMessage                          // Previous chat messages for context
	LLMClient    *llm.Client                                // Client for LLM interactions
	SkillMatcher func(string, *BERAgent) (BaseSkill, error) // Function to match skills to messages
	IsAction     bool                                       // Whether this is an action execution
	MemoryStore  *memory_store.MemoryStore                  // Store for persisting data between steps
	WorkflowId   string                                     // Unique ID for this workflow
}

// Common workflow error types
var (
	ErrInvalidInput    = errors.New("invalid workflow input")
	ErrSkillNotFound   = errors.New("no suitable skill found")
	ErrLLMFailure      = errors.New("LLM request failed")
	ErrTemplateFailure = errors.New("template rendering failed")
)

// Generates random string for memory keys
func generateRandomString() string {
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return hex.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	}
	return hex.EncodeToString(randomBytes)
}

// Checks if message is an action command
func IsAction(message string) bool {
	return strings.HasPrefix(message, "@ber approve")
}

// Main workflow execution function
func ExecuteWorkflow(ctx context.Context, input WorkflowInput) WorkflowResult {
	logger.Log.Debug("Starting workflow execution")
	input.WorkflowId = generateRandomString()
	logger.Log.Debug("Memory store data", zap.Any("data", input.MemoryStore.GetAllData()))

	input.IsAction = IsAction(input.Message)

	result := processWorkflow(ctx, input)
	formattedResponse, err := renderWorkflowResponse(result.templateData, result.templateToUse, result.baseSkill, result.executionError, input)

	if err != nil {
		return WorkflowResult{
			Response: fmt.Sprintf("Template rendering failed: %v", err),
			Error:    fmt.Errorf("%w: %v", ErrTemplateFailure, err),
		}
	}

	return WorkflowResult{
		Response: formattedResponse,
		Error:    result.executionError,
	}
}

// Internal result type for workflow processing
type workflowProcessResult struct {
	templateData   interface{} // Data to render in template
	templateToUse  string      // Template string to use
	baseSkill      BaseSkill   // Matched skill
	executionError error       // Any error during execution
}

// Core workflow processing logic
func processWorkflow(ctx context.Context, input WorkflowInput) workflowProcessResult {
	// Validate input parameters
	if err := validateWorkflowInput(&input); err != nil {
		logger.Log.Error("Input validation failed", zap.Error(err))
		return workflowProcessResult{
			templateData: map[string]interface{}{
				"error":   fmt.Errorf("%w: %v", ErrInvalidInput, err).Error(),
				"success": false,
			},
			templateToUse:  DefaultErrorTemplate,
			executionError: fmt.Errorf("%w: %v", ErrInvalidInput, err),
		}
	}
	logger.Log.Debug("Input validation successful")

	// Match skill to user message
	baseSkill, skillErr := input.SkillMatcher(input.Message, input.Agent)
	if skillErr != nil {
		logger.Log.Error("Failed to find matching skill", zap.Error(skillErr))
		return workflowProcessResult{
			templateData: map[string]interface{}{
				"error":   fmt.Errorf("%w: %v", ErrSkillNotFound, skillErr).Error(),
				"success": false,
			},
			templateToUse:  DefaultErrorTemplate,
			executionError: fmt.Errorf("%w: %v", ErrSkillNotFound, skillErr),
		}
	}
	logger.Log.Info("Found matching skill", zap.String("skill_type", fmt.Sprintf("%T", baseSkill)))

	// Execute workflow based on type
	logger.Log.Debug("Executing workflow", zap.Bool("is_action", input.IsAction))
	var templateData interface{}
	var execErr error

	if input.IsAction {
		logger.Log.Info("Executing action workflow")
		templateData, execErr = executeActionWorkflow(ctx, baseSkill, input)
	} else {
		logger.Log.Info("Executing LLM workflow")
		templateData, execErr = executeLLMWorkflow(ctx, baseSkill, input)
	}

	if execErr != nil {
		return workflowProcessResult{
			templateData: map[string]interface{}{
				"error":   execErr.Error(),
				"success": false,
			},
			templateToUse:  DefaultErrorTemplate,
			executionError: execErr,
		}
	}

	// Handle successful execution
	if input.IsAction {
		return workflowProcessResult{
			templateData: map[string]interface{}{
				"success": true,
				"message": "✅ Action executed successfully",
			},
			templateToUse: "{{if .success}}{{.message}}{{else}}{{.error}}{{end}}",
			baseSkill:     baseSkill,
		}
	}

	return workflowProcessResult{
		templateData:  templateData,
		templateToUse: baseSkill.GetTemplate(),
		baseSkill:     baseSkill,
	}
}

// Renders final response using template
func renderWorkflowResponse(templateData interface{}, templateToUse string, baseSkill BaseSkill, executionError error, input WorkflowInput) (string, error) {
	renderer := template.NewRenderer()

	// Add available actions to template if skill has any
	if baseSkill != nil && !input.IsAction {
		if actions := baseSkill.GetActions(); len(actions) > 0 {
			templateToUse += "\n\n---\n\n#### Available Actions\n"
			templateToUse += "{{if .success}}\nThe following actions are available:\n"
			templateToUse += "{{range $name, $_ := .actions}}\n"
			templateToUse += fmt.Sprintf("- Type `@ber {{$name}} %s` to execute this action\n", input.WorkflowId)
			templateToUse += "{{end}}{{end}}"

			templateMap := convertTemplateDataToMap(templateData, actions)
			templateData = templateMap
		}
	}

	return renderer.RenderAgentTemplate(templateToUse, templateData)
}

// Converts template data to map format
func convertTemplateDataToMap(templateData interface{}, actions map[string]func(context.Context, interface{}) error) map[string]interface{} {
	// Use existing map if possible
	if templateMap, ok := templateData.(map[string]interface{}); ok {
		if _, exists := templateMap["success"]; !exists {
			templateMap["success"] = true
		}

		actionNames := make(map[string]interface{})
		for name := range actions {
			actionNames[name] = struct{}{}
		}
		templateMap["actions"] = actionNames

		return templateMap
	}

	// Handle DNS record schema specially
	if dnsData, ok := templateData.(struct{ Records []interface{} }); ok {
		return map[string]interface{}{
			"success": true,
			"actions": createActionNames(actions),
			"records": dnsData.Records,
		}
	}

	// Convert other types via JSON marshaling
	jsonBytes, err := json.Marshal(templateData)
	if err != nil {
		logger.Log.Error("Failed to marshal template data", zap.Error(err))
		return map[string]interface{}{
			"success": false,
			"error":   "Failed to process template data",
			"actions": createActionNames(actions),
		}
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &dataMap); err != nil {
		logger.Log.Error("Failed to unmarshal template data", zap.Error(err))
		return map[string]interface{}{
			"success": false,
			"error":   "Failed to process template data",
			"actions": createActionNames(actions),
		}
	}

	dataMap["success"] = true
	dataMap["actions"] = createActionNames(actions)
	return dataMap
}

// Creates map of action names
func createActionNames(actions map[string]func(context.Context, interface{}) error) map[string]interface{} {
	actionNames := make(map[string]interface{})
	for name := range actions {
		actionNames[name] = struct{}{}
	}
	return actionNames
}

// Executes an action workflow
func executeActionWorkflow(ctx context.Context, skill BaseSkill, input WorkflowInput) (interface{}, error) {
	actionName, payload, err := GetActionDetailsFromMessage(input.Message, input.MemoryStore)
	if err != nil {
		return nil, err
	}

	typedResponse, err := ConvertToSchema[any](payload)
	if err != nil {
		return nil, fmt.Errorf("failed to convert payload: %w", err)
	}

	// Run pre-action hook
	if hook := skill.GetHooks().PreAction; hook != nil {
		var responseInterface interface{} = typedResponse
		if err := executeHook(ctx, hook, &responseInterface); err != nil {
			return nil, fmt.Errorf("pre-action hook failed: %w", err)
		}
		typedResponse = responseInterface
	}

	// Execute action
	action, exists := skill.GetActions()[actionName]
	if !exists {
		return nil, fmt.Errorf("action %s not found", actionName)
	}
	if err := executeAction(ctx, action, typedResponse); err != nil {
		return nil, fmt.Errorf("action %s failed: %w", actionName, err)
	}

	// Run post-action hook
	if hook := skill.GetHooks().PostAction; hook != nil {
		if err := executeHook(ctx, hook, &typedResponse); err != nil {
			return nil, fmt.Errorf("post-action hook failed: %w", err)
		}
	}

	return typedResponse, nil
}

// Executes an LLM workflow
func executeLLMWorkflow(ctx context.Context, skill BaseSkill, input WorkflowInput) (interface{}, error) {
	// Prepare messages for LLM request
	messages := append([]llm.ChatMessage{
		{
			Role:    llm.RoleSystem,
			Content: skill.GetPrompt(),
		},
	}, input.ChatHistory...)

	_ = skill.GetSkill()
	schema := skill.GetSchema()

	// Query LLM with schema
	response, err := llm.QueryWithSchema(ctx, input.LLMClient, messages, schema)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLLMFailure, err)
	}

	typedResponse, err := ConvertToSchema[any](response)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response: %w", err)
	}

	// Execute hooks in sequence
	hooks := []struct {
		name string
		hook func(context.Context, *interface{}) error
	}{
		{"PostLLMRequest", skill.GetHooks().PostLLMRequest},
		{"PreValidate", skill.GetHooks().PreValidate},
		{"PostValidate", skill.GetHooks().PostValidate},
	}

	for _, h := range hooks {
		if h.hook != nil {
			var responseInterface interface{} = typedResponse
			if err := executeHook(ctx, h.hook, &responseInterface); err != nil {
				return nil, fmt.Errorf("%s hook failed: %w", h.name, err)
			}
			typedResponse = responseInterface
		}
	}

	// Run validators
	for name, validator := range skill.GetValidators() {
		if err := executeValidator(ctx, validator, typedResponse); err != nil {
			return nil, fmt.Errorf("Validator %s failed:\n %w", name, err)
		}
	}

	// Store response if actions exist
	if input.MemoryStore != nil && len(skill.GetActions()) > 0 {
		for actionName := range skill.GetActions() {
			key := fmt.Sprintf("%s-%s", actionName, input.WorkflowId)
			input.MemoryStore.Set(key, typedResponse)
		}
	}

	return typedResponse, nil
}

// Default skill matching implementation using embeddings
func DefaultSkillMatcher(ctx context.Context, message string, agent *BERAgent, client *llm.Client) (BaseSkill, error) {
	userEmbedding, err := client.GetEmbedding(ctx, message)
	if err != nil {
		logger.Log.Error("Failed to get user message embedding", zap.Error(err))
		return nil, err
	}

	var bestSkill BaseSkill
	var highestSimilarity float32

	// Find best matching skill by comparing embeddings
	for _, skill := range agent.Skills {
		skillEmbedding, err := client.GetEmbedding(ctx, skill.GetPrompt())
		if err != nil {
			logger.Log.Error("Failed to get skill embedding", zap.Error(err))
			continue
		}

		similarity, err := llm.CosineSimilarity(userEmbedding, skillEmbedding)
		if err != nil {
			continue
		}

		if similarity > highestSimilarity {
			highestSimilarity = similarity
			bestSkill = skill
		}
	}

	if bestSkill == nil {
		logger.Log.Error("No suitable skill found")
		return nil, errors.New("no suitable skill found")
	}

	return bestSkill, nil
}

// Executes a hook with timeout and panic recovery
func executeHook(ctx context.Context, hook func(context.Context, *interface{}) error, data *interface{}) error {
	logger.Log.Debug("Starting hook execution")
	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error("Hook execution panicked", zap.Any("panic", r))
				done <- fmt.Errorf("hook panicked: %v", r)
			}
		}()
		err := hook(ctx, data)
		if err != nil {
			logger.Log.Error("Hook execution failed", zap.Error(err))
		} else {
			logger.Log.Debug("Hook execution completed successfully")
		}
		done <- err
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		logger.Log.Error("Hook execution timed out")
		return errors.New("hook timed out")
	}
}

// Executes a validator with timeout and panic recovery
func executeValidator(ctx context.Context, validator func(context.Context, interface{}) error, data interface{}) error {
	logger.Log.Debug("Starting validator execution")
	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error("Validator execution panicked", zap.Any("panic", r))
				done <- fmt.Errorf("validator panicked: %v", r)
			}
		}()
		err := validator(ctx, data)
		if err != nil {
			logger.Log.Error("Validator execution failed", zap.Error(err))
		} else {
			logger.Log.Debug("Validator execution completed successfully")
		}
		done <- err
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		logger.Log.Error("Validator execution timed out")
		return errors.New("validator timed out")
	}
}

// Executes an action with timeout and panic recovery
func executeAction(ctx context.Context, action func(context.Context, interface{}) error, data interface{}) error {
	logger.Log.Debug("Starting action execution")
	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error("Action execution panicked", zap.Any("panic", r))
				done <- fmt.Errorf("action panicked: %v", r)
			}
		}()
		err := action(ctx, data)
		if err != nil {
			logger.Log.Error("Action execution failed", zap.Error(err))
		} else {
			logger.Log.Debug("Action execution completed successfully")
		}
		done <- err
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		logger.Log.Error("Action execution timed out")
		return errors.New("action timed out")
	}
}

// Validates workflow input parameters
func validateWorkflowInput(input *WorkflowInput) error {
	if input.Agent == nil {
		return errors.New("agent is required")
	}
	if input.LLMClient == nil {
		return errors.New("LLM client is required")
	}
	if input.SkillMatcher == nil {
		return errors.New("skill matcher is required")
	}
	if input.Message == "" {
		return errors.New("message is required")
	}
	if input.MemoryStore == nil {
		return errors.New("memory store is required")
	}
	return nil
}

// Extracts action details from user message
func GetActionDetailsFromMessage(message string, memoryStore *memory_store.MemoryStore) (string, interface{}, error) {
	message = strings.TrimSpace(strings.TrimPrefix(message, "@ber "))
	logger.Log.Debug("Parsed message after trimming prefix", zap.String("message", message))

	parts := strings.Fields(message)
	if len(parts) < 2 {
		logger.Log.Error("Not enough parts in message", zap.Int("parts_count", len(parts)))
		return "", nil, errors.New("not enough parts in message, expected at least 2")
	}

	requestedAction := parts[0]
	workflowId := parts[1]

	if memoryStore == nil {
		logger.Log.Error("No memory store provided")
		return "", nil, errors.New("no memory store provided")
	}

	value, exists := memoryStore.Get(fmt.Sprintf("%s-%s", requestedAction, workflowId))
	if exists {
		logger.Log.Debug("Found matching key",
			zap.Any("value", value),
			zap.String("workflow_id", workflowId))
		return requestedAction, value, nil
	}

	logger.Log.Error("No matching key found",
		zap.String("action", requestedAction),
		zap.String("workflow_id", workflowId))
	return "", nil, errors.New("no matching key found for action and random")
}
