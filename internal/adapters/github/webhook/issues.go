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

package webhook

import (
	"encoding/json"

	"github.com/berbyte/ber-os/internal/adapters/github/actions"
	"github.com/berbyte/ber-os/internal/adapters/github/rest"
	"github.com/berbyte/ber-os/internal/logger"
	g "github.com/google/go-github/v64/github"
	"go.uber.org/zap"
)

func consumeIssuesEvent(payload []byte) error {
	var p g.IssueEvent
	err := json.Unmarshal(payload, &p)
	if err != nil {
		logger.Log.Error("Failed to Unmarshal Issue Event", zap.String("tag", "github-webhook-issues"), zap.Error(err))
		return err
	}

	// Add a reaction to the comment to indicate that the bot has processed it
	err = rest.NewIssueReaction(
		*p.Repository.Owner.Login, // Repository owner username
		*p.Repository.Name,        // Repository name
		int(*p.Issue.Number),      // Number of the comment to react to
		"eyes",                    // Reaction emoji to add
	)
	if err != nil {
		logger.Log.Error("Failed to Post Reaction on Comment", zap.String("tag", "github-webhook-issues"), zap.Error(err))
		return nil
	}

	tpl, err := actions.IssueComment(
		*p.Repository.Owner.Login,
		*p.Repository.Name,
		*p.Issue.Number,
		*p.Issue.Body,
	)
	if err != nil {
		return err
	}

	// Post the formatted response back to GitHub as a new comment
	err = rest.NewIssueComment(
		*p.Repository.Owner.Login, // Repository owner username
		*p.Repository.Name,        // Repository name
		*p.Issue.Number,           // Issue number
		tpl,                       // Formatted response content
	)
	if err != nil {
		logger.Log.Error("Failed to Post Comment", zap.String("tag", "github-webhook-issues"), zap.Error(err))
	}

	return nil
}
