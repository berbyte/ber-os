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

	"github.com/berbyte/ber-os/internal/logger"
	"go.uber.org/zap"
)

func preLLMRequest(ctx context.Context, data *DNSRecordSchema) error {
	logger.Log.Info("Pre LLM Request", zap.String("tag", "cloudflare-hooks"))
	return nil
}

func postLLMRequest(ctx context.Context, data *DNSRecordSchema) error {
	logger.Log.Info("Post LLM Request", zap.String("tag", "cloudflare-hooks"))
	return nil
}

func preValidate(ctx context.Context, data *DNSRecordSchema) error {
	logger.Log.Info("Pre Validate", zap.String("tag", "cloudflare-hooks"))
	return nil
}

func postValidate(ctx context.Context, data *DNSRecordSchema) error {
	logger.Log.Info("Post Validate", zap.String("tag", "cloudflare-hooks"))
	return nil
}

func preAction(ctx context.Context, data *DNSRecordSchema) error {
	logger.Log.Info("Pre Action", zap.String("tag", "cloudflare-hooks"))
	return nil
}

func postAction(ctx context.Context, data *DNSRecordSchema) error {
	logger.Log.Info("Post Action", zap.String("tag", "cloudflare-hooks"))
	return nil
}
