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
	"github.com/berbyte/ber-os/internal/adapters/github/client"

	llm "github.com/berbyte/ber-os/pkg/openai"

	"github.com/google/go-github/v64/github"
)

// NewIssueComment creates a new comment on a GitHub issue
func NewIssueComment(repoOwner string, repoName string, issueNumber int, message string) error {
	// Initialize GitHub client with repository owner's token
	client, ctx, err := client.NewClient(repoOwner)
	if err != nil {
		return err
	}

	// Create a new issue comment struct with the message
	comment := &github.IssueComment{
		Body: github.String(message), // Convert message string to GitHub string pointer
	}

	// Post the comment to the specified issue using GitHub API
	_, _, err = client.Issues.CreateComment(ctx, repoOwner, repoName, issueNumber, comment)
	if err != nil {
		return err
	}

	return nil
}

// GetIssueComments retrieves all comments from a GitHub issue and converts them to chat messages
func GetIssueComments(repoOwner string, repoName string, issueNumber int) ([]llm.ChatMessage, error) {
	// Initialize GitHub client with repository owner's token
	client, ctx, err := client.NewClient(repoOwner)
	if err != nil {
		return nil, err
	}

	// First fetch the issue itself to get the initial body
	issue, _, err := client.Issues.Get(ctx, repoOwner, repoName, issueNumber)
	if err != nil {
		return nil, err
	}

	// Fetch all comments for the issue from GitHub API
	comments, _, err := client.Issues.ListComments(ctx, repoOwner, repoName, issueNumber, nil)
	if err != nil {
		return nil, err
	}

	// Initialize slice to store converted chat messages, +1 for the issue body
	messages := make([]llm.ChatMessage, len(comments)+1)

	// Add the issue body as the first message
	if issue.Body != nil {
		messages[0] = llm.ChatMessage{
			Role:    llm.RoleUser,
			Content: *issue.Body,
		}
	}

	// Convert comments to chat messages
	for i, comment := range comments {
		// Skip comments with nil body
		if comment.Body == nil {
			continue
		}

		// Determine message role based on comment author type
		if comment.User != nil && comment.User.Type != nil && *comment.User.Type == "Bot" {
			// Bot comments are treated as assistant messages
			messages[i+1] = llm.ChatMessage{
				Role:    llm.RoleAssistant,
				Content: *comment.Body,
			}
		} else {
			// Human user comments are treated as user messages
			messages[i+1] = llm.ChatMessage{
				Role:    llm.RoleUser,
				Content: *comment.Body,
			}
		}
	}

	return messages, nil
}
