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

package actions

import (
	"context"
	"errors"
	"time"

	"github.com/berbyte/ber-os/internal/adapters/github/rest"
	"github.com/berbyte/ber-os/internal/agent"
	"github.com/berbyte/ber-os/internal/logger"
	"github.com/berbyte/ber-os/internal/services/memory_store"
	llm "github.com/berbyte/ber-os/pkg/openai"
	"go.uber.org/zap"
)

func IssueComment(owner, repo string, number int, message string) (string, error) {
	log := logger.GetLogger()

	log.Debug("handling issue comment",
		zap.String("owner", owner),
		zap.String("repo", repo),
		zap.Int("issue_number", number),
		zap.String("message", message))

	// Get issue labels and check for BER label
	selectedAgentLabels, err := handleLabels(owner, repo, number)
	if err != nil {
		log.Error("failed to handle labels",
			zap.String("owner", owner),
			zap.String("repo", repo),
			zap.Int("issue_number", number),
			zap.Error(err))
		return "", err
	}

	// If no agent was selected, return early
	if len(selectedAgentLabels) == 0 {
		return "‼️ Please **select a label** in the top-right corner to choose your BER Agent, then send a new message to @ber.", nil
	}

	if len(selectedAgentLabels) > 1 {
		return "‼️ Please **select a single label** in the top-right corner to choose your BER Agent, then send a new message to @ber.", nil
	}

	// Get the agent by tag
	selectedAgent := agent.GetAgentByTag(selectedAgentLabels[0])
	if selectedAgent == nil {
		log.Error("selected agent not found",
			zap.String("agent_tag", selectedAgentLabels[0]))
		return "", errors.New("selected agent not found")
	}

	// Get conversation history
	messages, err := rest.GetIssueComments(owner, repo, number)
	if err != nil {
		log.Error("failed to get issue comments",
			zap.String("owner", owner),
			zap.String("repo", repo),
			zap.Int("issue_number", number),
			zap.Error(err))
		return "", err
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create LLM client
	llmclient := llm.NewClient()

	// Execute workflow
	memoryStore := memory_store.NewMemoryStore()
	result := agent.ExecuteWorkflow(ctx, agent.WorkflowInput{
		Message:     message,
		Agent:       selectedAgent,
		ChatHistory: messages,
		LLMClient:   llmclient,
		SkillMatcher: func(msg string, a *agent.BERAgent) (agent.BaseSkill, error) {
			return agent.DefaultSkillMatcher(ctx, msg, a, llmclient)
		},
		MemoryStore: memoryStore,
	})

	if result.Error != nil {
		log.Error("workflow execution failed",
			zap.String("agent_tag", selectedAgentLabels[0]),
			zap.Error(result.Error))
	}

	return result.Response, nil
}

func handleLabels(owner, repo string, number int) ([]string, error) {
	log := logger.GetLogger()

	// Get issue labels
	labels, err := rest.GetIssueLabels(owner, repo, number)
	if err != nil {
		log.Error("failed to get issue labels",
			zap.String("owner", owner),
			zap.String("repo", repo),
			zap.Int("issue_number", number),
			zap.Error(err))
		return nil, err
	}

	// Check for BER labels
	var berLabels []string
	for _, label := range labels {
		if label != nil && label.Name != nil && len(*label.Name) >= 4 && (*label.Name)[:4] == "BER:" {
			berLabels = append(berLabels, (*label.Name)[4:])
		}
	}

	if len(berLabels) > 0 {
		return berLabels, nil
	}

	// Create labels if none found
	if err := rest.CreateLabels(owner, repo); err != nil {
		log.Error("failed to create labels",
			zap.String("owner", owner),
			zap.String("repo", repo),
			zap.Error(err))
		return nil, err
	}

	return nil, nil
}
