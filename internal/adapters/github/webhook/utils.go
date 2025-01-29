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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"

	"github.com/berbyte/ber-os/internal/logger"
	g "github.com/google/go-github/v64/github"
	"go.uber.org/zap"
)

// verifySignature validates that a webhook payload came from GitHub by checking
// the HMAC signature. It computes an HMAC using SHA-256 with the webhook secret
// and compares it to the signature provided in the X-Hub-Signature-256 header.
//
// Parameters:
//   - payload: The raw webhook payload bytes
//   - signature: The signature from the X-Hub-Signature-256 header
//
// Returns:
//   - bool: true if signatures match, false otherwise
func verifySignature(payload []byte, signature string) bool {
	// Create HMAC hasher with webhook secret
	key := hmac.New(sha256.New, []byte(os.Getenv("GH_WEBHOOK_SECRET")))
	// Write payload bytes to hasher
	key.Write([]byte(string(payload)))
	// Compute hex-encoded SHA-256 HMAC signature with "sha256=" prefix
	computedSignature := "sha256=" + hex.EncodeToString(key.Sum(nil))
	// Log computed signature for debugging
	logger.Log.Debug("computed signature", zap.String("signature", computedSignature), zap.String("tag", "github-webhook-utils"))

	// Compare computed signature with provided signature
	return computedSignature == signature
}

// checkMessageForBer checks if the webhook payload contains a message addressed to @ber.
// It parses the payload based on the event type and looks for "@ber" mentions.
// The function also validates that the message comes from a real GitHub user, not a bot.
//
// For IssueComment events, it checks the comment body.
// For Issues events, it checks the issue body.
// Other event types always return false.
//
// Parameters:
//   - payload: The raw webhook payload bytes to parse
//   - event: The type of GitHub webhook event (IssueComment or Issues)
//
// Returns:
//   - bool: true if the message is from a user and mentions @ber, false otherwise
func checkMessageForBer(payload []byte, event Event) bool {
	switch event {
	case IssueComment:
		// Parse the issue comment event payload
		var commentEvent g.IssueCommentEvent
		if err := json.Unmarshal(payload, &commentEvent); err != nil {
			logger.Log.Error("Failed to Unmarshal Issue Comment Event", zap.String("tag", "github-webhook-utils"), zap.Error(err))
			return false
		}

		// Validate comment, user and user type exist
		if commentEvent.Comment == nil || commentEvent.Comment.User == nil || commentEvent.Comment.User.Type == nil {
			return false
		}

		// Skip if not from a real GitHub user
		if string(*commentEvent.Comment.User.Type) != "User" {
			logger.Log.Info("Skipping Comment from Bot/App", zap.String("tag", "github-webhook-utils"))
			return false
		}

		// Check if comment body exists and mentions @ber
		return commentEvent.Comment.Body != nil &&
			strings.Contains(*commentEvent.Comment.Body, "@ber")

	case Issues:
		// Parse the issues event payload
		var issuesEvent g.IssuesEvent
		if err := json.Unmarshal(payload, &issuesEvent); err != nil {
			logger.Log.Error("Failed to Unmarshal Event", zap.String("tag", "github-webhook-utils"), zap.Error(err))
			return false
		}

		// We are only interested in issue creation events
		if issuesEvent.Action != nil && *issuesEvent.Action != "opened" {
			return false
		}

		// Validate issue, user and user type exist
		if issuesEvent.Issue == nil || issuesEvent.Issue.User == nil || issuesEvent.Issue.User.Type == nil {
			return false
		}

		// Skip if not from a real GitHub user
		if string(*issuesEvent.Issue.User.Type) != "User" {
			logger.Log.Info("Skipping Comment from Bot/App", zap.String("tag", "github-webhook-utils"))
			return false
		}

		// Check if issue body exists and mentions @ber
		return issuesEvent.Issue.Body != nil &&
			strings.Contains(*issuesEvent.Issue.Body, "@ber")

	default:
		return false
	}
}
