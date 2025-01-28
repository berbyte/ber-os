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
	"io"
	"net/http"

	"github.com/berbyte/ber-os/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Event string

const (
	Install      Event = "installation"
	Issues       Event = "issues"
	IssueComment Event = "issue_comment"
)

// Map of supported events and their handlers
var eventHandlers = map[Event]func([]byte) error{
	Install:      consumeInstallEvent,
	Issues:       consumeIssuesEvent,
	IssueComment: consumeIssueCommentEvent,
}

func ConsumeEvent(c *gin.Context) {
	log := logger.GetLogger()
	log.Debug("starting to consume event")

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Debug("failed to read request body", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"reason": "Failed to read request body"})
		return
	}
	log.Debug("read request body", zap.Int("payload_size", len(payload)))

	event := Event(c.GetHeader("X-GitHub-Event"))
	log.Debug("received github event", zap.String("event_type", string(event)))

	signature := c.GetHeader("X-Hub-Signature-256")
	log.Debug("checking signature", zap.String("signature", signature))
	if !verifySignature(payload, signature) {
		log.Debug("signature verification failed")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "Signature mismatch"})
		return
	}
	log.Debug("signature verified successfully")

	// Only process events that mention @ber or are installation events
	if event != Install && !checkMessageForBer(payload, event) {
		log.Debug("message not addressed to @ber, skipping processing",
			zap.String("event", string(event)))
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	log.Debug("message is relevant for processing")

	// Handle event if supported
	if handler, exists := eventHandlers[event]; exists {
		log.Info("consuming event", zap.String("event", string(event)))
		log.Debug("found handler for event", zap.String("event", string(event)))

		if err := handler(payload); err != nil {
			log.Error("couldn't consume event",
				zap.String("event", string(event)),
				zap.Error(err))
			log.Debug("handler failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"reason": err})
			return
		}

		log.Info("consumed event", zap.String("event", string(event)))
		log.Debug("handler completed successfully")
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	log.Warn("unsupported event", zap.String("event", string(event)))
	log.Debug("no handler found for event type")
	c.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{"reason": "Unsupported event: " + string(event)})
}
