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

package rest

import (
	"fmt"

	"github.com/berbyte/ber-os/internal/adapters/github/client"
	"github.com/berbyte/ber-os/internal/agent"
	"github.com/berbyte/ber-os/internal/logger"
	"github.com/google/go-github/v64/github"
	"go.uber.org/zap"
)

// CreateLabels creates all the labels defined in agent configs for a repository
func CreateLabels(repoOwner string, repoName string) error {
	client, ctx, err := client.NewClient(repoOwner)
	if err != nil {
		return err
	}

	// Iterate through all agents and create their labels
	for _, agent := range agent.GetRegisteredAgents() {
		logger.Log.Info("Creating Label", zap.String("issue-tag", agent.Tag), zap.String("tag", "github-rest-labels"))

		label := &github.Label{
			Name:        github.String(fmt.Sprintf("BER:%s", agent.Tag)),
			Color:       github.String("1F618D"),
			Description: github.String(agent.Description),
		}

		_, _, err := client.Issues.CreateLabel(ctx, repoOwner, repoName, label)
		if err != nil {
			// If label already exists, continue silently
			if _, ok := err.(*github.ErrorResponse); ok {
				continue
			}
			logger.Log.Error("Failed to Create Label", zap.String("issue-tag", agent.Tag), zap.String("tag", "github-rest-labels"), zap.Error(err))
			return fmt.Errorf("Failed to create label %s: %w", agent.Tag, err)
		}
		logger.Log.Info("Label Created", zap.String("issue-tag", agent.Tag), zap.String("tag", "github-rest-labels"))
	}

	return nil
}
